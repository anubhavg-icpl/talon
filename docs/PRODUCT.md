# Talon product baseline (long-lived E2E)

Shipped, aligned surface — CI-tested Go core + production Next.js operator UI.
Features that need separate product/infra decisions stay **out of scope**.

## What ships end-to-end

| Layer | Capability |
|-------|------------|
| **Pipeline** | recon → exploit → post-exploit → Forge (codegen) → judge |
| **HITL** | `nmap_scan` approve / reject / edit |
| **Tools** | Arsenal MCP + Strike MCP (MSF) |
| **Skills** | CyberStrike catalog (~7.6k) + pentest agent skills (59) + SOC analysis skills (50) + `skill_search` / `skill_get` |
| **Agents** | full / recon / web / network / exploit / post · A2A `delegate_*` |
| **Findings** | 3-gate evidence, triage, registry, kill-chain, methodology |
| **Agent tools** | Evidence store, crypto toolkit (29 ops), HTTP probe batch, web headers audit, JS endpoint extract |
| **Anti-hallucination** | Duplicate detection, degraded tool health, stall detection, completion gate |
| **Target state** | Per-target persistence, snapshots, deterministic resume-plan builder |
| **Traffic store** | Per-run HTTP traffic recording, search, JSONL persist |
| **Recap** | LLM-free run recap (solve path, evidence, reproduction) + CI presets |
| **Reports** | Structured markdown, export, HTML print, AI analyze |
| **Control plane** | Auth, runs, config, MCP list, health |
| **Engagements** | Scope/ROE, targets, schedules, webhooks, credentials, budget, batch |
| **Run ops** | Compare, timeline, notes, retest, export |
| **CLI** | `talon` → core API |
| **Dashboard** | Next.js operator shell (cyan void theme) |
| **WebGL** | Dynamic Three.js: ambient starfield, TalonGlobe, SkeletonUtils + Soldier.glb |
| **Lab** | `vuln-target` real vsftpd 2.3.4 |
| **Brand** | README: `assets/talon-mark-red.webp` only · UI: electric cyan chrome |

## Operator path (aligned)

```
Operator
  ├─ talon CLI  ──HTTP──►  talon-core (:8000)
  │                           ├─stdio MCP──► talon-arsenal ──► arsenal-engine
  │                           └─stdio MCP──► talon-strike  ──► msfrpcd
  │                           └─ optional talon-relay (RabbitMQ)
  └─ dashboard (:3000)
        └─ /api/talon/* proxy ──► talon-core
```

### UI routes ↔ nav

| Nav | Route | Primary APIs |
|-----|-------|--------------|
| Overview | `/overview` | `/runs`, `/runs/summary`, findings, skills, globe |
| Showcase | `/showcase` | Three.js ExamplesStage (SkeletonUtils + globe) + still reel |
| Runs | `/runs`, `/runs/[id]`, `/runs/new` | start, status, tools, WS/SSE, findings, report, notes, retest |
| Findings | `/findings` | `/findings` |
| Compare | `/compare` | `/runs/compare` |
| Engagements | `/ops` | scope, targets, schedules, notify, credentials, budget, batch |
| Agents | `/agents` | `/agents` |
| Skills | `/skills` | `/skills` |
| Playbooks | `/playbooks` | `/playbooks` |
| Intel | `/intel` | `/intel` |
| Settings | `/settings` | `/config`, `/health/services`, `/mcp/servers` |
| Login | `/login` | `/auth/login` |

### Theme (aligned)

- **UI chrome:** void black + electric cyan (not brand-red)
- **Destructive red:** failures only
- **WebGL:** pause when tab hidden; client-only dynamic import
- **README hero:** red mark only (`talon-mark-red.webp`)

## Longevity rules

1. `go test ./...` and `next build` stay green  
2. No half-wired control stubs on `main`  
3. Auth when Postgres + `TALON_ADMIN_PASSWORD`  
4. Authorized targets only  
5. Deferred table needs product decision — no thin stubs  

## Explicitly deferred

| Item | Today |
|------|--------|
| Full multi-team RBAC | Single admin session |
| Interactive Meterpreter UI | Tools via agents/CLI |
| Burp-class proxy | Skills + Arsenal |
| Token $ invoices | Counters only |
| True video gen | Stills + optional MP4 slots |

## Verify alignment

```bash
# Backend
go test ./...
go build -o bin/talon ./cmd/talon
go build -o bin/talon-core ./cmd/talon-core

# Dashboard
cd web && ./node_modules/.bin/tsc --noEmit
./node_modules/.bin/next build

# Stack smoke
docker compose up -d --build
./bin/talon status
# UI: http://localhost:3000 (or DASHBOARD_PORT) → login → overview → showcase
```

## Related

- [README.md](../README.md)  
- [FEATURE_MAP.md](FEATURE_MAP.md)  
- `.env.example`  
