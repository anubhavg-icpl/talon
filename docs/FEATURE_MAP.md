# Feature map: agentic-os-personal + CyberStrike → Talon

This document records what was extracted from the source projects and what
was implemented into Talon (wave 1). Full monorepo ports are out of scope;
Talon keeps its Go orchestration model and absorbs portable patterns.

## Source inventories

### agentic-os-personal (`content-os`)

| Feature | Fit for Talon | Notes |
|---------|---------------|--------|
| Job runner + status machine | Partial | Talon already has run lifecycle + store |
| SSE live updates | Already present | `/monitor/stream` + WS |
| AI ranker / OpenRouter | Partial | Talon multi-provider LLM factory already covers this |
| RSS / Firecrawl collectors | **No** | Content publishing domain |
| Zernio social publish | **No** | Not pentest |
| SQLite freshness window | **No** | Domain mismatch |
| Studio caption versioning | **No** | Domain mismatch |

**Takeaway:** agentic-os is a **content OS**, not offensive security. Portable
ideas: structured job results, SSE progress, multi-model config UI (already
in Talon settings).

### CyberStrike

| Feature | Fit for Talon | Wave |
|---------|---------------|------|
| Multi-stage agent pipeline | Already present | recon→exploit→post→codegen→judge |
| Specialized agent prompts | **Implemented** | Skills injection per stage |
| 7,600+ signed skills | **Lite port** | Builtin skill pack + `/skills` API |
| 3-gate confirmation protocol | **Implemented** | `GateEvidence` + validation |
| Structured findings (`report_vulnerability`) | **Implemented** | Extract from tool log + API |
| Severity triage rules | **Implemented** | In report skill + findings severity |
| Multi-section `generate_report` | **Implemented** | `StructuredReport` markdown |
| 13+ domain agents (web/cloud/mobile) | Future | Would need more tools/MCP |
| Bolt remote tool execution | Future | Talon MCP multi already exists |
| HackBrowser / proxy testers | Future | Large TS surface |
| Post-exploit hooks (eBPF, win/mac) | Future | Data/scripts only if needed |
| MCP ecosystem (176+ tools) | Partial | Arsenal + Strike MCP today |
| TUI / multi-provider catalog | Partial | CLI + OpenAI/Bedrock/Ollama |

## Implemented in Talon (wave 1)

### Backend (`internal/core`)

| File | Role |
|------|------|
| `findings.go` | Finding model, 3-gate validation, extract from tool log, dedup, summary |
| `report.go` | Multi-section report (exec summary, findings, methodology, timeline, validation) |
| `skills.go` | Builtin methodology skills + `InjectSkills` into subagent prompts |
| `orchestrator.go` | Skills on subagents; `finalizeResult` attaches findings+report |
| `prompts.go` | Report/post/codegen prompts reference 3-gate + severity |

### Control plane (`internal/control`)

| Endpoint | Purpose |
|----------|---------|
| `GET /runs/{id}/findings` | Structured findings + severity roll-up |
| `GET /runs/{id}/report` | Full structured report (lazy-build for legacy runs) |
| `GET /skills` | Skill catalog (`?brief=1` for meta only) |
| `GET /output/status/{id}` | Adds `findings_summary`, `findings_count`, `has_report` |

Session persistence (JSON + Postgres jsonb) stores `Findings` and `Report`.

### Dashboard (`web`)

- New **Findings** tab with severity badges and 3-gate evidence cards
- **Report** tab prefers structured markdown from `/runs/{id}/report`
- API client: `getFindings`, `getReport`, `getSkills`

## Not ported (intentional)

- Content-OS social publishing / scrapers
- CyberStrike HackBrowser + 8 proxy testers as live HTTP agents
- 7,600 skill files (signing, lazy load) — replaced by small in-process pack
- Bolt multi-host remote execution fleet
- Full VRT / methodology DB from CyberStrike

## Wave 5 (completed) — agents *use* CyberStrike skills + MCP/A2A clarity

| Feature | Status |
|---------|--------|
| `skill_search` / `skill_get` tools on all subagents | Done |
| hybridExec routes skills + findings + MCP | Done |
| Prompts instruct agents to load skills on demand | Done |
| `GET /mcp/servers` includes talon-core virtual tools + A2A notes | Done |
| Settings: MCP + Agent-to-Agent panel | Done |
| Agents page: communication model | Done |

**Agent-to-agent model:** orchestrator → subagents via `delegate_*` only (hub-and-spoke).  
**MCP:** hexstrike (arsenal) + metasploit (strike) stdio.  
**Skills:** lazy load from 7.6k CyberStrike pack via tools + UI.

## Wave 4 (completed) — full CyberStrike skill catalog in UI

| Feature | Status |
|---------|--------|
| **7,656** CyberStrike `SKILL.md` files → `skills/cyberstrike/` | Done (~38MB) |
| Categories: WEB, mitre_attack, CIS, NIST, attack-*, postexploit, … | Done |
| `GET /skills` paginated + `category`/`q`/`stage`/`limit`/`offset` | Done |
| `GET /skills/{id}` full body | Done |
| Skills UI: sidebar categories, search, pagination, detail pane, recent | Done |
| CLI `talon skills --category WEB --q ssrf` | Done |
| Bounded agent injection (catalog full in UI only) | Done |

## Wave 3 (completed) — skills pack + live findings wiring

| Feature | Status |
|---------|--------|
| Disk skills + builtins | Done |
| Bounded prompt injection (12 disk + full builtins) | Done |
| Mid-run findings progress → store | Done |
| SSE/WS `findings` events | Done |
| Dashboard live findings on RunDetail | Done |
| Overview FINDINGS + SKILLS LOADED cards | Done |

## Wave 2 (completed)

| Feature | Status |
|---------|--------|
| Live `report_finding` / `triage_finding` tools | Done — hybrid exec on all subagents |
| Agent modes (full/recon/web/network/exploit/post) | Done — API + New Operation UI |
| Kill chain analysis | Done — API + Kill Chain tab |
| Methodology coverage | Done — API + Kill Chain tab |
| Disk skills (`skills/*.md`, `TALON_SKILLS_DIR`) | Done |
| Global findings registry UI | Done — `/findings` |
| Skills / Agents dashboard pages | Done |
| Report export (.md) | Done |
| Finding triage UI + API | Done |
| Runs list findings_count + agent_mode | Done |

## Still not ported (by design)

- Content-OS social publishing
- CyberStrike HackBrowser + 8 live proxy testers
- 7,600 signed skill files + Bolt multi-host fleet
- Full VRT methodology database

## Product longevity (shipped baseline)

Canonical long-lived E2E surface: **[PRODUCT.md](PRODUCT.md)** (routes, theme, WebGL, verify steps).

Shipped waves above are **done**. Also aligned on `main`:

| Surface | Status |
|---------|--------|
| Operator cyan UI theme | Done — void + electric cyan |
| Three.js E2E (globe, starfield, SkeletonUtils + Soldier.glb) | Done — Showcase ExamplesStage |
| Dynamic WebGL import + tab pause | Done |
| Nav: Engagements `/ops`, Showcase, Intel, Playbooks… | Done |
| README hero | `talon-mark-red.webp` only |

The following need **separate product / infra decisions** — no thin stubs:

| Item | Reality today |
|------|----------------|
| Full RBAC multi-team | Single admin session auth |
| Interactive Meterpreter UI | Session commands via agents / tools |
| Real browser/proxy testers | Skills + Arsenal tools, not Burp-class |
| Token $ cost from LLM provider | Counters only (not invoices) |
| True video generation | Showcase WebP stills + optional MP4 slots |

**Rule:** keep `go test ./...` and `web` `next build` green; never commit
unfinished control-plane files that reference missing Server fields.

---

## Pentest agent port (wave 2)

Source: Python pentest agent (Typer+FastAPI, v0.3.7)
Target: Talon `internal/core/` Go packages + `skills/vulnclaw/` knowledge tree.

### What was ported

| Pentest agent feature | Talon file | Wave | Notes |
|-----------------|-----------|------|-------|
| Skill catalog (7 core + 50 specialized + 2 warstories) | `skills/vulnclaw/*.md` | 1 | 59 skill files, Chinese translated to English (~90%) |
| Evidence store (append-only, dedup, preview) | `internal/core/evidence.go` | 2 | EvidenceStore + Record/Get/List/Search + bounded previews |
| Evidence agent tools | `internal/core/evidence_tools.go` | 2 | `evidence_list`, `evidence_view`, `evidence_search` |
| Anti-hallucination gate + correction layer | `internal/core/correction.go` | 3 | Duplicate detection, degraded health, stall detection, completion gate |
| Crypto toolkit (29 operations) | `internal/core/crypto_tools.go` | 4 | base64/32/58, hex, URL, HTML, Unicode, ROT13/Caesar/Morse, MD5/SHA/AES/DES, JWT, auto_decode |
| HTTP probe batch | `internal/core/probe_tools.go` | 4 | `http_probe_batch` — multiple request variants in one call |
| Web vulnerability tools | `internal/core/probe_tools.go` | 4 | `web_headers_audit`, `js_endpoint_extract` |
| Target-state persistence | `internal/core/targetstate.go` | 5 | Per-target store, snapshots, deterministic resume-plan builder |
| Traffic evidence store | `internal/core/traffic.go` | 6 | Per-run HTTP traffic recording, search, JSONL persist |
| Deterministic recap report | `internal/core/recap.go` | 7 | LLM-free run recap with solve path, evidence, reproduction |
| Headless CI presets | `internal/core/recap.go` | 7 | `quick`/`standard`/`deep` presets with exit codes |

### Explicitly not ported (Talon already has equivalent or out of scope)

| Pentest agent feature | Reason |
|-----------------|--------|
| i18n/zh-default | Talon is English-only |
| ChatGPT OAuth bridge | Out of scope |
| Original FastAPI web + React frontend | Talon has Next.js web |
| Typer CLI/TUI | Talon has Go CLI |
| ChromaDB KB | Superseded by Talon's 7.7k skill catalog |
| Team/parallel-agent orchestration | Talon hub-spoke delegates cover it |
| MCP registry/lifecycle | Talon internal/mcpclient covers it |
| mitmproxy capture backend | Deferred (traffic store records but doesn't capture) |

---

## SOC analysis port (wave 3)

Source: SOC detection skill library (50 skills)
Target: Talon `internal/core/` + `skills/detection/` + `internal/control/` + web UI.

### What was ported

| SOC analysis feature | Talon file | Notes |
|-------------------------|-----------|-------|
| 50 skills (triage/investigation/tuning) | `skills/detection/{triage,investigation,tuning}/*.md` | Talon-native format with `# stage:` / `# category:` headers |
| Triage engine (check-based, majority rule) | `internal/core/detection.go` | `ApplyTriageDecision()` — escalate/dismiss from N checks |
| Investigation engine (multi-signal verdict) | `internal/core/detection.go` | `ApplyInvestigationVerdict()` — malicious/suspicious/benign from A/B/C/D signals |
| Tuning engine (propose detection changes) | `internal/core/detection.go` | `TuningState` — exclude/include/modify/fork/none |
| Pipeline flow (triage→investigation→tuning) | `internal/core/detection.go` | `DetermineNextStage()` routes based on verdict |
| Case management | `internal/core/detection.go` | `CaseStore` — create, get, list, update triage/investigation/tuning |
| Agent tools (5 tools) | `internal/core/detection_tools.go` | `detection_create_case`, `detection_triage`, `detection_investigate`, `detection_tune`, `detection_list_cases` |
| Control plane routes (3 routes) | `internal/control/detection_handlers.go` | `GET /detection/cases`, `GET /detection/skills`, `GET /detection/skills/{type}` |
| Alert Triage UI page | `web/src/views/detection/DetectionView.tsx` | Browse 50 playbooks by type (triage/investigation/tuning) with detail panel |
| API client | `web/src/lib/api.ts` | `getDetectionSkills`, `getDetectionSkillsByType`, `getDetectionCases` |

### Pipeline architecture

```
Alert → detection_create_case
  → detection_triage (triage skill, N checks → majority verdict)
    ├── escalate → detection_investigate (investigation skill, signals → verdict)
    │   ├── malicious/suspicious → detection_tune (propose detection change)
    │   └── benign/inconclusive → done
    └── dismiss → detection_tune (review false positive for tuning)
```
