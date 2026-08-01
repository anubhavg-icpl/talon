# Talon product baseline (long-lived E2E)

This is the **shipped product surface**. It is intentionally stable: CI-tested,
documented, and operator-ready for authorized engagements. Features that need
separate product/infra decisions stay **out of scope** until explicitly adopted.

## What ships end-to-end

| Layer | Capability |
|-------|------------|
| **Pipeline** | recon → exploit → post-exploit → Forge (codegen) → judge |
| **HITL** | `nmap_scan` approve / reject / edit |
| **Tools** | Arsenal MCP (nmap, nuclei, …) + Strike MCP (MSF 12 tools) |
| **Skills** | CyberStrike catalog (~7.6k) + builtins; agents use `skill_search` / `skill_get` |
| **Agents** | Modes: full / recon / web / network / exploit / post; A2A via `delegate_*` |
| **Findings** | 3-gate evidence, triage, global registry, kill-chain + methodology |
| **Reports** | Multi-section markdown; export; AI analyze |
| **Control plane** | Auth sessions, runs, config, MCP list, health probes |
| **Platform ops** | Scope/ROE, targets, schedules, webhooks, credentials, evidence, budget |
| **Ops UX** | Batch start, compare, timeline, notes, retest, HTML→PDF |
| **Surfaces** | `talon` CLI + Next.js dashboard + HTTP/WS/SSE API |
| **Lab** | `vuln-target` real vsftpd 2.3.4 E2E path |
| **Brand** | All `assets/*.webp` + Showcase reel + globe |

## Operator paths (durable)

```
talon CLI  ──HTTP──►  talon-core  ──stdio MCP──►  talon-arsenal ──► arsenal-engine
     │                     │                 └──►  talon-strike  ──► msfrpcd
     │                     └── optional talon-relay (RabbitMQ)
     └── dashboard (:3000) proxies /api/talon/* → core
```

**Day-to-day:** `talon status` → `talon run start … --watch` → findings/report.  
**UI:** login → Overview / Runs / Findings / Skills / Agents / Ops / Settings.  
**Persistence:** Postgres when configured; else `TALON_DATA_DIR/runs.json` + `platform.json`.

## Longevity rules

1. **Main stays green** — `go test ./...` and core/control/cli packages in CI.
2. **No half-wired modules** on `main` (uncompiled stubs that break the package).
3. **Auth by default** when Postgres + `TALON_ADMIN_PASSWORD` are set.
4. **Authorized use only** — lab profile is intentional malware for local validation.
5. **Portable patterns over monorepo clones** — no full CyberStrike/HackBrowser port unless product decides.

## Explicitly deferred (product / infra)

| Item | Today | Needs |
|------|--------|--------|
| Full multi-team RBAC | Single admin session auth | Tenancy model, persistence, audit |
| Interactive Meterpreter UI | Session tools via agents/CLI | Live console UX + safety gates |
| Burp-class browser/proxy | Skills + Arsenal tools | Proxy capture store, certs, traffic UI |
| Token $ invoices | Token/char counters | Provider usage APIs + pricing SKUs |
| True video generation | Showcase stills + optional MP4 slots | Video API availability / budget |

Do **not** implement these as thin stubs without a product decision.

## Verify local health

```bash
# Go
go test ./...
go build -o bin/talon ./cmd/talon
go build -o bin/talon-core ./cmd/talon-core

# Stack (compose)
docker compose up -d --build
./bin/talon status

# Dashboard types (from web/)
pnpm check-types   # or: npx tsc --noEmit
```

## Related docs

- [README.md](../README.md) — install, API, lab, visual tour  
- [FEATURE_MAP.md](FEATURE_MAP.md) — source extraction waves  
- `.env.example` — required secrets and LLM providers  
