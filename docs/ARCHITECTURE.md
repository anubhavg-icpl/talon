# Talon Architecture (Scalable System Design)

Production architecture reference for the Talon platform. Grounded in the
actual codebase (Go modules, Postgres schema, Redis cache, RabbitMQ broker,
MCP tool boundary). Cross-references the [System Design Primer notes](system-design-primer-notes.md)
where the general theory lives.

---

## 1. System overview

Talon is an **AI-driven penetration-testing orchestration platform**. It
coordinates a multi-stage validation pipeline (recon, exploit, post-exploit,
codegen fallback, judge) against authorized targets, with human-in-the-loop
gating on invasive tools (nmap). Operators interact via a CLI or a Next.js
dashboard; the control plane is a Go HTTP server.

```
                          Operator
                              │
               ┌──────────────┼──────────────┐
               │              │              │
          talon CLI      Dashboard        Webhook
          (HTTP)         (Next.js)        (future)
               │              │              │
               │      /api/talon/* proxy    │
               │              │              │
               ▼              ▼              ▼
        ┌──────────────────────────────────────────┐
        │              talon-core (:8000)            │
        │  HTTP API + orchestrator + HITL gates      │
        │                                            │
        │  ┌─────────┐  ┌──────────┐  ┌───────────┐  │
        │  │ Arsenal │  │  Strike  │  │   Forge   │  │
        │  │  MCP    │  │   MCP    │  │ (Docker)  │  │
        │  └────┬────┘  └────┬─────┘  └───────────┘  │
        │       │ stdio      │ stdio                  │
        └───────┼────────────┼───────────────────────┘
                │            │
                ▼            ▼
         arsenal-engine   msfrpcd
         (BlackArch)      (Kali)
         :8888 / :7681    :5554

  ┌─────────────────────────────────────────────────────┐
  │                   Data layer                         │
  │  Postgres (:5432)  Redis (:6380)  RabbitMQ (:5672)   │
  └─────────────────────────────────────────────────────┘

         ┌─────────────────────────────┐
         │       talon-relay           │  (optional AMQP worker;
         │  consumes execute_agent_task │   scales horizontally)
         └─────────────────────────────┘
```

---

## 2. Service inventory (every component, what it does)

| Service | Binary / image | Port(s) | Role | Scales how |
|---------|----------------|---------|------|------------|
| **talon-core** | `talon:latest` | 8000 | HTTP API + orchestrator | Vertical (stateless API); horizontal for API-only load |
| **talon-relay** | `talon:latest` | (none) | AMQP worker, same orchestrator | **Horizontal** — N replicas behind one RabbitMQ queue |
| **arsenal-engine** | `talon-arsenal-engine` | 8888, 7681 | BlackArch tool runner (nmap, nuclei, ...) | One per host (privileged, host network) |
| **metasploit** | `talon-msf-rpc` | 5554 | `msfrpcd` HTTPS msgpack RPC | Single instance per deployment |
| **dashboard** | `talon-dashboard` | 3000 | Next.js ops console | Horizontal behind a load balancer |
| **postgres** | `postgres:17-alpine` | 5432 | Run history, auth sessions, settings | Vertical + read replicas for read-heavy dashboards |
| **redis** | `redis:7-alpine` | 6380 | Cache: health probes, AI analysis, sessions | Single (cache only; noop fallback if down) |
| **rabbitmq** | `rabbitmq:3-management-alpine` | 5672, 15672 | Broker for relay workers | Cluster of 3 for HA |
| **ollama** (profile) | `ollama/ollama` | 11434 | Local LLM runtime | GPU-bound; one per GPU node |
| **onnx-slm** (profile) | `talon-onnx-slm` | 8090 | SmolLM / ONNX SLM runtime | CPU-bound; horizontal for concurrent chat |
| **vuln-target** (profile) | `talon-vuln-target` | 21, 6200 | Lab target (CVE-2011-2523) | Lab only, never in prod |

---

## 3. Data flow (request lifecycle)

### 3.1 Start a run (operator → result)

```
Operator ──POST /input/start──► talon-core
  talon-core:
    1. Store.New(runID) — in-memory session created
    2. Persist to Postgres (upsertRun) or runs.json fallback
    3. Orchestrator.Run(ctx, input) — goroutine:
       a. Recon agent (MCP: arsenal nmap_scan — HITL gated)
       b. Exploit agent (MCP: strike run_exploit)
       c. Post-exploit agent (MCP: strike send_session_command)
       d. Forge fallback (Docker sandbox codegen) if modules fail
       e. Judge agent (separate model evaluates success)
    4. Each stage writes ToolLog + findings to store
    5. Finalize: structured findings (3-gate) + markdown report
    6. Persist final state to Postgres
  Operator streams via WS/SSE: /monitor/ws/{runID}
```

### 3.2 Dashboard streaming (resilient degradation)

```
Browser ──► /api/talon/* proxy ──► talon-core
  WebSocket:  /monitor/ws/{runID}     ← preferred (real-time push)
       │
       ├─ fallback → SSE: /monitor/stream/{runID}  (if WS fails)
       │
       └─ fallback → polling GET /output/status/{runID} every 3s
```

A broken stream tier never breaks the view. The dashboard auto-degrades.

### 3.3 Cache strategy (Redis)

| Key pattern | TTL | Fallback (if Redis down) |
|-------------|-----|--------------------------|
| `health:{service}` | 5s | Live probe on every request |
| `analysis:{runID}` | 24h | Re-run LLM analysis (idempotent) |
| `session:{token}` | 60s | Postgres lookup on every request |

Redis is **pure cache**. Losing it degrades performance, not correctness.

---

## 4. Scalability profile

### What scales horizontally today

| Component | Mechanism | Constraint |
|-----------|-----------|------------|
| **talon-relay** | N replicas consume from one RabbitMQ queue | Each run is a single message; no partitioning needed |
| **dashboard** | Stateless Next.js; sessions in Postgres/Redis | Put behind a load balancer |
| **talon-core (API)** | HTTP routes are stateless; runs persist to Postgres | Orchestrator goroutines are per-instance (no cluster coordination yet) |

### What does NOT scale horizontally yet

| Component | Why | Mitigation |
|-----------|-----|------------|
| **talon-core (orchestrator)** | Run state is in-memory per instance; no distributed lock | Keep one active orchestrator instance; use relay for fan-out |
| **arsenal-engine** | Privileged, host network, one per host | Deploy per-host; route via consistent hashing if multi-host |
| **metasploit** | Single `msfrpcd` instance | Shard by engagement if throughput-bound |
| **Postgres** | Single primary (no replication configured) | Add read replica; route dashboard reads to replica |

### Resource budget (from docker-compose limits)

| Service | Memory limit | CPU limit | Notes |
|---------|-------------|-----------|-------|
| talon-core | 2 GB | 2.0 | Orchestrator + HTTP + LLM client |
| talon-relay | 2 GB | 2.0 | Per worker; add workers for throughput |
| arsenal-engine | 4 GB | 3.0 | Security tools are memory-hungry |
| metasploit | 2 GB | 2.0 | Ruby + msfrpcd |
| postgres | 1 GB | 1.5 | `shared_buffers` tunable via env |
| redis | 256 MB | 0.5 | Pure cache, no persistence |
| rabbitmq | 1 GB | 1.5 | Broker only |
| dashboard | 512 MB | 1.0 | Next.js standalone |
| ollama | 8 GB | 4.0 | GPU model weights in RAM/VRAM |
| onnx-slm | 4 GB | 3.0 | CPU inference |

**Single-host minimum**: ~12 GB RAM for the core stack (without LLM runtimes).
**With local LLM**: ~20 GB RAM (ollama) or ~16 GB (onnx-slm).

---

## 5. Availability and reliability

### Failure modes and design responses

| Failure | System behavior | Design principle |
|---------|----------------|------------------|
| **Redis down** | All endpoints fall back to uncached path; noop cache | Cache is a speed optimization, not a dependency |
| **Postgres down** | Store falls back to `runs.json`; auth disabled | Graceful degradation (not suitable for prod auth) |
| **RabbitMQ down** | Relay cannot consume; core still serves API + runs directly | Core is independently useful |
| **arsenal-engine down** | Health probe marks it offline; recon tools fail with MCP error | Operator sees degraded health in dashboard |
| **metasploit down** | Strike MCP fails to connect; exploit tools unavailable | Health probe surfaces this |
| **LLM provider down/timeout** | Run errors after `OPENAI_HTTP_TIMEOUT` / `TALON_RUN_TIMEOUT` | Timeouts are configurable and enforced |
| **talon-core crash** | Non-terminal runs marked `error` on reload (crash recovery) | `markInterruptedRunsError()` in db.go |
| **Dashboard crash** | Compose `restart: unless-stopped`; CLI still works | Operator always has CLI fallback |

### What is NOT HA today

- **Single talon-core instance** (orchestrator state is in-memory).
- **Single Postgres primary** (no streaming replication).
- **Single RabbitMQ node** (no cluster quorum).
- **No health-based routing** (single-host `network_mode: host`).

To achieve multi-9 availability, see [DEPLOYMENT.md](DEPLOYMENT.md) for the
HA topology (reverse proxy + replicated workers + Postgres replica +
RabbitMQ quorum).

---

## 6. Security architecture

### Trust model

```
[Untrusted]          [Trusted network]              [Privileged]
  Browser     ──►    Dashboard (:3000)     ──►     talon-core (:8000)
                                                          │
                                         ┌────────────────┤
                                         │                │
                                  arsenal-engine     msfrpcd
                                  (privileged)       (Kali)
```

- **Browser** only talks to the dashboard (same-origin proxy to core). Core is
  never directly exposed to the browser.
- **Dashboard → core** is server-side via `TALON_CORE_URL` (localhost by
  default). No CORS exposure of the core.
- **Core → MCP children** is stdio (arsenal, strike), not network. No port
  exposure for MCP servers.

### Authentication and authorization

- **Session-based auth**: `POST /auth/login` with `{username, password}`.
- **Password hashing**: bcrypt-equivalent stored in Postgres `users` table.
- **Session storage**: Postgres `sessions` table with 30-day TTL,
  Redis-cached (60s) for lookup speed.
- **Session transport**: HTTP-only cookie `talon_session` (browser) or
  `Authorization: Bearer <token>` (CLI).
- **Route gating**: all routes except `/health*` and `/auth/login` require a
  valid session when Postgres is configured and `TALON_ADMIN_PASSWORD` is set.
- **Auth disabled mode**: `TALON_AUTH_DISABLED=true` reverts to open access
  (dev/CI only — never in prod).

### Current limitations (prod hardening needed)

| Gap | Risk | Mitigation |
|-----|------|------------|
| No TLS on core/dashboard | Session cookie in cleartext | Reverse proxy with TLS (see [DEPLOYMENT.md](DEPLOYMENT.md)) |
| `sslmode=disable` on DB URL | DB traffic unencrypted | `sslmode=require` + provision server cert |
| CORS `*` on core API | Any origin can call core (if exposed) | Restrict to dashboard origin behind proxy |
| Single admin user | No role separation | Implement RBAC before multi-operator use |
| docker.sock mount on core/relay | Container escape risk | Restrict socket access; use rootless Docker or gVisor |
| `privileged: true` on arsenal | Full host access | Required for security tools; isolate on dedicated host |
| No rate limiting | DoS via API | Add rate limiter at proxy layer |

### Secret management

- **Required secrets** (compose fails fast without them): `MSF_PASSWORD`,
  `RABBITMQ_PASSWORD`, `TALON_PG_PASSWORD`, `TALON_ADMIN_PASSWORD`.
- **No hardcoded secrets** in any image or compose file.
- **LLM credentials**: `OPENAI_API_KEY`, `AWS_ACCESS_KEY_ID`,
  `AWS_SECRET_ACCESS_KEY` — passed via env, never persisted to disk by Talon.
- **For prod**: inject secrets via Docker secrets, HashiCorp Vault, or cloud
  secret manager — not plaintext `.env`.

---

## 7. Persistence design

### Postgres schema (from `internal/control/db.go`)

```sql
-- Auth
users     (id UUID PK, username CITEXT UNIQUE, password_hash TEXT, created_at)
sessions  (token TEXT PK, user_id UUID FK→users, created_at, expires_at)

-- Run history (jsonb document + queryable summary columns)
runs      (run_id TEXT PK, data JSONB, updated_at,
           target TEXT, cve_id TEXT, service_name TEXT,
           status TEXT, judge_verdict BOOLEAN, tool_calls INT,
           started_at TIMESTAMPTZ)
          INDEX: runs_started_at_idx (started_at DESC)
          INDEX: runs_status_idx (status)

-- Runtime config
settings  (key TEXT PK, value JSONB, updated_at)
```

**Design decisions:**

1. **JSONB + summary columns**: full run state is a JSONB document (flexible),
   but pagination uses indexed summary columns (fast at any table size).
   This is the "wide column for queries, JSONB for schema flexibility" pattern.
2. **Idempotent migrations**: `CREATE TABLE IF NOT EXISTS` + `ALTER TABLE ADD
   COLUMN IF NOT EXISTS` — safe to run on every boot.
3. **Crash recovery**: `markInterruptedRunsError()` scans for non-terminal runs
   on startup and marks them `error` (orchestrator state is in-memory and
   cannot resume after a crash).
4. **`ON CONFLICT DO UPDATE`** (upsert) for idempotent run writes — safe under
   retries and network partitions.

### Cache invalidation

Redis keys are TTL-based (no explicit invalidation except logout, which
`DEL`s the session key). This is acceptable because:
- Health probes re-check every 5s regardless.
- AI analysis is idempotent (same input, same output).
- Session lookups hit Postgres on cache miss.

---

## 8. Tool boundary (MCP)

The Model Context Protocol (MCP) is the **single tool boundary** in Talon.
Agents never call tools directly; they call MCP servers over stdio.

```
talon-core (orchestrator)
  ├── talon-arsenal (MCP stdio server)
  │     └── HTTP → arsenal-engine (:8888)
  │           └── nmap, nuclei, gobuster, sqlmap, ... (BlackArch tools)
  │
  ├── talon-strike (MCP stdio server)
  │     └── HTTPS msgpack → msfrpcd (:5554)
  │           └── Metasploit: exploit, post, auxiliary, sessions
  │
  └── Forge (in-process, not MCP)
        └── Docker sandbox → Python exploit codegen
```

**Why MCP over stdio:**
- Process isolation (crash in a tool doesn't crash core).
- Clean contract (JSON schema in, JSON out).
- Language-agnostic (arsenal-engine is Python; strike is Go).

---

## 9. LLM provider abstraction

Talon supports four LLM backends behind a unified factory
(`internal/llm/factory.go`):

| Provider | Use case | Streaming | Tool calling |
|----------|----------|-----------|--------------|
| `bedrock` | Enterprise (AWS) | Converse API | Native |
| `openai` | Hosted or self-hosted (vLLM, z.ai, LiteLLM) | SSE | Native function calling |
| `ollama` | Local GGUF models | One-shot today | Native tools |
| `onnx` | Local SmolLM (millisecond tokens) | SSE (fast) | Text `TOOL_CALL` protocol |

**Role separation**: the orchestrator uses one model per role (`main`, `judge`,
`code`), all resolved through `config.ResolveModel()`. This lets you run a
strong model for the main agent and a cheap model for the judge.

**Hybrid mode**: `LLM_PROVIDER=openai` for the orchestrator + `LLM_CODE_PROVIDER=ollama`
for codegen only.

---

## 10. Related documents

| Document | Scope |
|----------|-------|
| [DEPLOYMENT.md](DEPLOYMENT.md) | Production deployment, TLS, backups, monitoring, HA topology |
| [PRODUCT.md](PRODUCT.md) | Shipped feature surface, UI routes, theme |
| [FEATURE_MAP.md](FEATURE_MAP.md) | Wave-by-wave feature history |
| [system-design-primer-notes.md](system-design-primer-notes.md) | General system design theory reference |
| [slm-onnx.md](slm-onnx.md) | Local SLM/ONNX runtime pipeline |
| [README.md](../README.md) | Getting started, CLI reference, API reference |
