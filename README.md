# Talon

![Talon](assets/hero-banner.webp)

**AI-driven penetration-testing orchestration, built natively in Go.**

Point Talon at a target (IP + optional CVE / service). It runs a full
validation pipeline: **recon → exploit → post-exploit → (optional) codegen
fallback → judge**. Operators drive it with the **`talon` CLI** or the HTTP
control plane. Only use this against systems you own or have **written
authorization** to test.

![Authorized use only](assets/security-warning-banner.webp)

> Offensive security tool. Unauthorized access is illegal. Read
> [Security & responsible use](#security--responsible-use).

---

## Table of contents

1. [What it does](#what-it-does)
2. [Components](#components)
3. [Quick start](#quick-start)
4. [Operator CLI (`talon`)](#operator-cli-talon)
5. [Web dashboard](#web-dashboard)
6. [Local E2E lab (CVE-2011-2523)](#local-e2e-lab-cve-2011-2523)
7. [Configuration](#configuration)
8. [HTTP API](#http-api)
9. [Architecture notes](#architecture-notes)
10. [Tools (Arsenal / Strike / Forge)](#tools-arsenal--strike--forge)
11. [Development](#development)
12. [Troubleshooting](#troubleshooting)
13. [Security & responsible use](#security--responsible-use)
14. [License](#license)

---

## What it does

```json
{
  "ip": "127.0.0.1",
  "cve_id": "CVE-2011-2523",
  "service_name": "vsftpd 2.3.4",
  "lhost": "127.0.0.1",
  "lport": 4446
}
```

| Stage | What happens |
|--------|----------------|
| **Recon** | Focused nmap (HITL gate), optional nuclei/smbmap |
| **Exploit** | Metasploit module search + `run_exploit` with session poll |
| **Post-exploit** | Shell commands on the session (proof of compromise) |
| **Forge** | If modules fail: LLM writes Python, runs in Docker sandbox |
| **Judge** | Second model returns whether compromise was real |

Human-in-the-loop: **`nmap_scan` is gated**. Approve, reject, or edit args
before the scan runs (CLI: `talon run approve|reject|edit`).

---

## Components

| Binary / service | Role |
|------------------|------|
| **`talon`** | Operator CLI → HTTP control plane |
| **`talon-core`** | Orchestrator + HTTP API (`:8000`) |
| **`talon-relay`** | Same orchestrator as AMQP worker |
| **`talon-arsenal`** | MCP stdio → Arsenal Engine (nmap, nuclei, …) |
| **`talon-strike`** | MCP stdio → Metasploit RPC (12 tools) |
| **Talon Forge** | Codegen sandbox (Docker) inside the orchestrator |
| **arsenal-engine** | Flask tool runner (Kali-based image) |
| **msf_rpc** | `msfrpcd` (Kali image) |
| **rabbitmq** | Broker for relay |
| **postgres** | Run history + dashboard users/sessions (auth) + runtime config |
| **redis** | Cache (health probes, AI analysis, sessions) — core runs uncached without it |
| **dashboard** | Web ops console (`:3000`) — Next.js UI over the core API |
| **vuln-target** | Lab only (`--profile vuln`) — real vsftpd 2.3.4 by default |

Core and relay spawn arsenal + strike as **local MCP stdio children**
(`HEXSTRIKE_MCP_PATH` / `METASPLOIT_MCP_PATH`).

---

## Quick start

### Docker layout (five Dockerfiles)

| Path | Image / service | Notes |
|------|-----------------|--------|
| **`Dockerfile`** | `talon:latest` → **talon-core** + **talon-relay** | All Go bins; one image, two commands |
| **`kali-msf/Dockerfile`** | **metasploit** (`msf_rpc`) | `msfrpcd` |
| **`arsenal-engine/Dockerfile`** | **arsenal-engine** | Kali tool runner |
| **`vuln-target/Dockerfile`** | **vuln-target** (profile `vuln`) | Targets: **`real`** (default) \| **`mimic`** |
| **`web/Dockerfile`** | **dashboard** | Next.js standalone build |

One compose file: **`docker-compose.yml`**. No overlays.

### Prerequisites

- Go **1.25+** (host builds / CLI)
- Docker + Compose
- An LLM backend: **Bedrock**, **OpenAI-compatible** (e.g. z.ai GLM), or **Ollama**
- For MSF sessions: working `msfrpcd` (compose service `metasploit`)

### 1. Configure secrets

```bash
cp .env.example .env
# Edit at least:
#   MSF_PASSWORD, RABBITMQ_PASSWORD, AMQP_URL
#   TALON_PG_PASSWORD, TALON_ADMIN_PASSWORD   (dashboard login + run persistence)
#   + your LLM keys (OPENAI_API_KEY or AWS_*, etc.)
```

### 2. Bring up the stack

```bash
# Core stack (MSF + Arsenal Engine + RabbitMQ + core + relay)
docker compose up -d --build

# Lab target — real infected vsftpd (MSF sessions work)
docker compose --profile vuln up -d --build vuln-target

# Optional: Python mimic instead of real vsftpd
# VULN_TARGET=mimic docker compose --profile vuln up -d --build vuln-target

# Optional local LLM
docker compose --profile ollama up -d
```

All services use **`network_mode: host`** so reverse shells and RPC stay simple.

### 3. Build the operator CLI (host)

```bash
go build -o bin/talon ./cmd/talon
go build -o bin/talon-core ./cmd/talon-core
go build -o bin/talon-strike ./cmd/talon-strike
go build -o bin/talon-arsenal ./cmd/talon-arsenal
go build -o bin/talon-relay ./cmd/talon-relay
# or: docker build -t talon:latest .
```

### 4. Health check

```bash
./bin/talon status
# overall=healthy — core, arsenal-engine, metasploit-rpc, rabbitmq
curl -s http://127.0.0.1:8000/health
# {"service":"talon-core","status":"ok"}
```

### Host-side core (dev / fast iteration)

Often easier than rebuilding images while hacking strike/core:

```bash
set -a; source .env; set +a
export MSF_SERVER=localhost
export LHOST=127.0.0.1
export HEXSTRIKE_MCP_PATH=$PWD/bin/talon-arsenal
export METASPLOIT_MCP_PATH=$PWD/bin/talon-strike
export HEXSTRIKE_SERVER_URL=http://localhost:8888
# keep your LLM_* / OPENAI_* from .env

./bin/talon-core   # listens on :8000
```

Compose still runs `msf_rpc`, `arsenal_engine`, and optionally `vuln-target`.

---

## Operator CLI (`talon`)

Config precedence: **flags → env → `~/.config/talon/config.yaml` → defaults**.

| Env / config | Meaning |
|--------------|---------|
| `TALON_CORE_URL` | Control plane (default `http://localhost:8000`) |
| `TALON_PROJECT_DIR` | Compose project root (for `talon logs`) |
| `TALON_OUTPUT` | `table` \| `json` \| `raw` |
| `TALON_CONFIG` | Config file path |

### Commands

```bash
# Stack health
talon status
talon status -o json

# Start a validation run
talon run start \
  --ip 127.0.0.1 \
  --cve CVE-2011-2523 \
  --service "vsftpd 2.3.4" \
  --lhost 127.0.0.1 \
  --lport 4446 \
  --watch \
  --auto-approve          # lab only: auto-approve nmap HITL

# Without --watch: manage the run yourself
talon run status <run_id>
talon run watch <run_id> --auto-approve --interval 4s
talon run approve <run_id>
talon run reject <run_id>
talon run edit <run_id> --args '{"target":"127.0.0.1","ports":"21,6200","scan_type":"-sT -Pn"}'
talon run tools <run_id> -o json
talon run traces <run_id>

# Logs (needs docker compose + project dir)
talon logs core --tail 100
talon logs arsenal -f
talon logs msf --tail 50
talon logs vuln

talon config
talon version
```

**Exit codes:** `0` ok / completed · non-zero for usage errors, failed runs,
or unreachable core (`talon run status` maps terminal run status to exit codes).

---

## Web dashboard

A full ops-console UI (**Talon Ops Console**) lives in `web/` — Next.js +
shadcn/ui, dark hacker theme, real-time everything over the core API.

```bash
docker compose up -d --build dashboard   # or just: docker compose up -d --build
# → http://localhost:3100 (or your DASHBOARD_PORT)
```

Login with `TALON_ADMIN_USERNAME` / `TALON_ADMIN_PASSWORD` from `.env`
(seeded into Postgres on first boot). The same credentials work for the CLI:
`talon auth login`.

| Page | What you get |
|------|--------------|
| `/overview` | Fleet stats, run-activity chart, verdicts donut, live active operations |
| `/runs` | Filterable/sortable table of every run (persisted across restarts) |
| `/runs/new` | Launch a run: target IP, CVE, service, LHOST/LPORT |
| `/runs/{id}` | Live terminal feed (WebSocket), **HITL approve/reject/edit gate**, report, **AI analysis**, traces |
| `/settings` | Live service health (7 probes), **config editor** (LLM/attacker/features, DB-backed), **MCP servers** panel |

Architecture: the browser only talks to the dashboard; a server-side proxy
route (`/api/talon/*`) forwards to `TALON_CORE_URL` (default
`http://localhost:8000`), so no CORS or exposed-core concerns. Live run
streams use **WebSocket first**, degrading to SSE then 3s polling. The
**AI ANALYSIS** tab on a run calls `POST /analyze/{run_id}`, which reuses the
LLM already configured on talon-core (no extra keys needed). Port override:
`DASHBOARD_PORT=3100 docker compose up -d dashboard`.

Local dev without Docker:

```bash
cd web && pnpm install && pnpm dev   # http://localhost:3000, needs core on :8000
```

---

## Local E2E lab (CVE-2011-2523)

Authorized **local** validation of the full pipeline (recon → MSF session →
root shell → report).

### Lab target (one Dockerfile, two targets)

```bash
# real (default) — infected vsftpd 2.3.4 + restart loop
docker compose --profile vuln up -d --build vuln-target
# FTP: 220 (vsFTPd 2.3.4); after USER *:) + PASS → shell on :6200

# mimic — Python only (recon/forge; no real MSF session)
VULN_TARGET=mimic docker compose --profile vuln up -d --build vuln-target
```

### Payload that works

For this CVE on current Metasploit:

| Setting | Value |
|---------|--------|
| Module | `exploit/unix/ftp/vsftpd_234_backdoor` |
| Payload | **`cmd/unix/reverse_bash`** (not `interact`, not reverse meterpreter) |
| Options | `RHOSTS`, `RPORT=21`, `LHOST`, `LPORT` |

Strike injects `cmd/unix/reverse_bash` + `LHOST`/`LPORT` when the agent omits
them. Session IDs from `session.list` are integer msgpack keys; the client
normalizes them so polling sees real shells.

### Full CLI E2E (proven)

```bash
# Stack: MSF + arsenal + real vuln + core (compose or host-side)
./bin/talon status

./bin/talon run start \
  --ip 127.0.0.1 \
  --cve CVE-2011-2523 \
  --service "vsftpd 2.3.4" \
  --lhost 127.0.0.1 \
  --lport 4446 \
  --watch --auto-approve -o json
```

Expected outcomes:

1. HITL nmap auto-approved (or `approve` / `edit`)
2. `run_exploit` → **`Session N created`** (seconds, not a 60s empty poll)
3. `send_session_command` → proof (`uid=0(root)` on the lab host/container)
4. Run **`status: completed`** with a final report

Inspect:

```bash
talon run tools <run_id> -o json   # look for run_exploit SESSION + shell output
```

### Manual MSF smoke (no agent)

```bash
# HTTPS msgpack RPC on MSF_PORT (default 5554), password = MSF_PASSWORD
# module.execute exploit unix/ftp/vsftpd_234_backdoor with
#   PAYLOAD=cmd/unix/reverse_bash LHOST=127.0.0.1 LPORT=<free>
# then session.list → shell session
```

---

## Configuration

Everything important is env vars (`.env` for compose). See **`.env.example`**.

### Required

| Variable | Used by |
|----------|---------|
| `MSF_PASSWORD` | `msfrpcd`, `talon-strike` |
| `RABBITMQ_PASSWORD` + `AMQP_URL` | RabbitMQ, `talon-relay` |
| LLM credentials | Depends on `LLM_PROVIDER` |

### Important defaults

| Variable | Default | Notes |
|----------|---------|--------|
| `MSF_PORT` | `5554` | SSL if `MSF_SSL=true` |
| `MSF_SSL` | `true` | Self-signed; strike skips verify |
| `LHOST` / `LPORT` | `127.0.0.1` / `4444` | Reverse listeners; core + strike |
| `LLM_PROVIDER` | `bedrock` | Or `openai` / `ollama` |
| `TALON_RUN_TIMEOUT` | `20m` | Wall clock per start/resume segment |
| `OPENAI_HTTP_TIMEOUT` | (client default) | e.g. `120s` for slow hosted APIs |
| `HEXSTRIKE_MCP_PATH` | sibling `talon-arsenal` | Core/relay spawn path |
| `METASPLOIT_MCP_PATH` | sibling `talon-strike` | Core/relay spawn path |

### LLM providers

```bash
# Hosted OpenAI-compatible (example: z.ai GLM)
LLM_PROVIDER=openai
OPENAI_BASE_URL=https://api.z.ai/api/coding/paas/v4
OPENAI_API_KEY=...
OPENAI_MAIN_MODEL=glm-5.2
OPENAI_JUDGE_MODEL=glm-5.2
OPENAI_CODE_MODEL=glm-5.2

# Local Ollama
LLM_PROVIDER=ollama
# docker compose --profile ollama up -d
# ollama create talon -f models/Modelfile

# Mix: orchestrator on hosted, codegen local
LLM_PROVIDER=openai
LLM_CODE_PROVIDER=ollama
```

---

## HTTP API

Base URL: **`http://localhost:8000`** (host networking).

| Method | Path | Purpose |
|--------|------|---------|
| `GET` | `/health` | Liveness (`{"status":"ok","service":"talon-core"}`) |
| `GET` | `/health/services` | Live probes of every dependency (postgres, arsenal, msfrpcd RPC auth, rabbitmq, ollama) with latency |
| `POST` | `/auth/login` | `{username,password}` → session token + `talon_session` cookie |
| `POST` | `/auth/logout` | Invalidate session, clear cookie |
| `GET` | `/auth/me` | Current session user |
| `POST` | `/input/start` | Start run (body: `ip`, `cve_id`, `service_name`, `description`, `lhost`, `lport`) |
| `GET` | `/runs` | Paginated run list (`?limit=&offset=` → `{runs,total,limit,offset}`, max 500) |
| `GET` | `/runs/summary` | Aggregate counts (total/active/compromised/awaiting/completed/errored) |
| `GET` | `/config` | Effective config (DB overrides → env → default), secrets masked |
| `PUT` | `/config` | Save whitelisted config keys to Postgres (FEATURE_* hot; LLM applies on restart) |
| `GET` | `/mcp/servers` | Connected MCP servers + their tool lists |
| `GET` | `/output/status/{run_id}` | `running` \| `awaiting_approval` \| `completed` \| … + optional `interrupt` / `output` |
| `POST` | `/output/resume/{run_id}` | HITL: `{"decision":"approve"\|"reject"\|"edit","edited_args":{...}}` |
| `POST` | `/analyze/{run_id}` | One-shot LLM analyst briefing (report + tool log → summary/kill-chain/findings/remediation) |
| `GET` | `/monitor/tools?run_id=` | Tool call log |
| `GET` | `/monitor/traces/{run_id}` | Message history |
| `GET` | `/monitor/stream/{run_id}` | **SSE** live stream: `tool` + `status` events, closes on terminal state |
| `GET` | `/monitor/ws/{run_id}` | **WebSocket** live stream: `{"type":"tool"\|"status","data":...}` messages |

The dashboard streams runs over **WebSocket first**, degrading to SSE and
then 3s polling if a tier is unavailable — a broken stream never breaks the
view. Override the WS target with `NEXT_PUBLIC_TALON_WS_URL`
(default `ws://<dashboard-host>:8000`).

CORS is permissive (`*`) for local tooling. **All routes except `/health*` and
`/auth/login` require a session** (cookie `talon_session` for the dashboard,
`Authorization: Bearer <token>` for the CLI) whenever Postgres is configured
and `TALON_ADMIN_PASSWORD` is set; set `TALON_AUTH_DISABLED=true` to run
without auth (old behavior).

Runs persist to **Postgres** (`TALON_DATABASE_URL`; users/sessions/run history
with queryable summary columns + indexes, so the registry paginates in SQL at
any size) and survive restarts; without a database core falls back to
`TALON_DATA_DIR/runs.json`, importing it into Postgres on first connect.
Runs interrupted mid-flight by a shutdown are marked `error` on reload.

**Redis** (`REDIS_URL`, compose service on :6380) caches service-health probes
(5s), AI analyses (24h per run), and session lookups (60s). Every cached
endpoint falls back to its uncached path when Redis is down.

**Runtime config** (`GET`/`PUT /config`) is editable from the dashboard
Settings page and stored in Postgres (whitelisted keys, secrets write-only).
`FEATURE_*` flags apply instantly; LLM/connection changes apply on restart.

Prefer the **CLI** or the **web dashboard** for day-to-day ops; the API is
what they wrap.

### AMQP (relay)

Relay consumes **`execute_agent_task`** and publishes **`agent.pentest.output`**.
Queue-driven runs auto-approve nmap today (no AMQP HITL yet).

---

## Architecture notes

- **MCP is the tool boundary** — arsenal and strike are stdio MCP servers, not
  in-process plugins.
- **HITL only on `nmap_scan`** by design; other tools run under the agent.
- **Recon budget** — limited nmap/nuclei turns; prefer `-sT -Pn` / `-sV -Pn`,
  avoid long `-sC` packs unless needed.
- **Exploit budget** — turn caps + LLM HTTP timeouts so runs finish.
- **Sessions** — `module.execute` may return `uuid` without `job_id`; strike
  still polls `session.list` (integer keys decoded correctly).
- **Forge** — needs Docker socket on core/relay for sandbox runs.

```
Operator (talon CLI)
        │ HTTP
   talon-core ──stdio MCP──► talon-arsenal ──HTTP──► arsenal-engine
        │              └──► talon-strike ──HTTPS msgpack──► msfrpcd
        └── optional talon-relay via RabbitMQ
```

---

## Tools (Arsenal / Strike / Forge)

### Arsenal (via Arsenal Engine)

Large catalog of offensive tools (nmap, nuclei, web/cloud/bug-bounty helpers,
…). Source of truth: `internal/arsenal/tools_generated.go` +
`tools_manual.go` and the engine at `arsenal-engine/`.

Default nmap: **`-sT -Pn`** with **`-T4 --host-timeout 60s`** (safe for HITL).

### Strike (12 Metasploit tools)

Including: `list_exploits`, `list_payloads`, `run_exploit`,
`run_post_module`, `run_auxiliary_module`, `list_active_sessions`,
`send_session_command`, `terminate_session`, `start_listener`, job helpers,
payload generation.

### Forge

LLM-generated Python exploits in a short-lived Docker sandbox with install
retry on missing packages.

---

## Development

```bash
go test ./...
go vet ./...
gofmt -l .

# Strike msgpack session-list regression
go test ./internal/strike/ -run 'TestDecodeMSFMap|TestAsStringKeyedMap' -count=1
```

### Layout

```
cmd/talon            operator CLI
cmd/talon-core       HTTP control plane
cmd/talon-relay      AMQP worker
cmd/talon-arsenal    Arsenal MCP server
cmd/talon-strike     Metasploit MCP server
internal/cli         CLI implementation
internal/control     HTTP routes + run store
internal/core        orchestrator / subagents / HITL
internal/strike      msfrpcd client + tools
internal/arsenal     engine proxy + tool tables
internal/forge       codegen sandbox
internal/llm         bedrock | openai | ollama
vuln-target/         lab images (real + mimic)
arsenal-engine/      Kali tool runner
kali-msf/            msfrpcd image
docker-compose.yml   single compose file (profiles: vuln, ollama)
```

### Images

| Build | Purpose |
|-------|---------|
| `docker compose build` | Platform + MSF + arsenal (and vuln if profile set) |
| `vuln-target` target `real` | Infected vsftpd + `entrypoint.sh` restart loop |
| `vuln-target` target `mimic` | `server.py` only |

One compose file; switch lab mode with **`VULN_TARGET=real|mimic`**.

---

## Troubleshooting

| Symptom | Fix |
|---------|-----|
| `MSF_PASSWORD is required` | Set in `.env`; must match `msfrpcd` |
| Core: MCP transport closed on metasploit | Strike failed to start — check MSF password/host/port/SSL |
| `awaiting_approval` forever | `talon run approve <id>` or `--auto-approve` in lab |
| `run_exploit` ~60s, “No session” | Rebuild strike with integer-key decoder; use **real** vuln image + `reverse_bash` + free `LPORT` |
| Port already in use (LPORT / 4444) | Pick another `--lport`; stop leftover handlers/sessions |
| FTP dead after one exploit | Real target entrypoint restarts vsftpd; `docker restart talon_vuln_target` if stuck |
| Judge says success, no session | Trust `talon run tools` (`run_exploit` / session tools), not the prose alone |
| LLM hangs | Set `OPENAI_HTTP_TIMEOUT`, `TALON_RUN_TIMEOUT`; check provider status |
| Forge fails | Mount Docker socket; host `docker` must work |

**Session list note:** empty `session.list` decodes fine; non-empty maps use
**integer keys**. Older clients that only decoded `map[string]any` always
missed sessions even when MSF had shells.

---

## Security & responsible use

- Only targets you own or have **written authorization** to test.
- Lab profile `vuln` is intentional malware/backdoor software for local use.
- No hardcoded production secrets — missing credentials fail fast.
- Host networking simplifies labs; harden networking for shared environments.
- HITL on nmap is a safety valve, not a full authorization system.

---

## License

No license file yet — all rights reserved by the repository owner unless an
explicit license is added. Contact the owner before redistribution.

---

*Talon: Go-native pentest orchestration. Operator path: `talon` CLI →
`talon-core` → arsenal/strike MCP → arsenal-engine / msfrpcd. Lab E2E:
real vsftpd 2.3.4 + `cmd/unix/reverse_bash` + session poll.*
