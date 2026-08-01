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
