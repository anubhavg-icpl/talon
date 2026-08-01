# Talon skills (disk pack)

Full **CyberStrike skill catalog** (~7,600 methodology files) lives under
`skills/cyberstrike/`, plus optional flat files in `skills/*.md`.

## Dashboard

Open **Skills** in the UI to:

- Browse by **category** (WEB, mitre_attack, CIS_benchmarks, NIST, attack-*, …)
- **Search** across name / id / body
- Filter by **stage**
- **Paginate** (metadata list is brief by default)
- Open **full methodology** in the detail pane
- **Revisit** recently opened skills (stored in browser localStorage)

API:

```
GET /skills?brief=1&category=WEB&q=ssrf&limit=40&offset=0
GET /skills/{id}          # full body
```

## Agent use (runtime tools)

Every subagent can call:

| Tool | Purpose |
|------|---------|
| `skill_search` | Query CyberStrike pack (`q`, optional `category`/`stage`/`limit`) |
| `skill_get` | Load full methodology by `id` from search hits |

Example agent flow: `skill_search q="ssrf"` → `skill_get id="disk-cyberstrike-attack-ssrf"` → apply steps → `report_finding`.

Builtins are also **injected** into prompts. Disk skills are **sampled** (max 12)
so context stays bounded — agents pull the rest on demand via tools.

## Agent-to-agent

Orchestrator coordinates specialists with `delegate_*` tools (not peer-to-peer):

```
orchestrator ──delegate_recon──► recon subagent ──► return text
     │
     ├──delegate_exploit──► exploit (MCP arsenal/strike + skills)
     ├──delegate_post_exploit──► post
     ├──delegate_codegen──► forge sandbox
     └──delegate_report──► findings/report
```

MCP stdio servers: **hexstrike** (talon-arsenal) + **metasploit** (talon-strike).

## Headers in skill files

```markdown
# stage: exploit
# category: WEB

## Body…
```

## Override directory

```bash
export TALON_SKILLS_DIR=/path/to/skills
```

Search order: `TALON_SKILLS_DIR` → `./skills` → `/app/skills`.
