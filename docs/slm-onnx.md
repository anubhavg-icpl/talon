# SLM / ONNX end-to-end (UI → Go → runtime)

This is the **aligned** path from the dashboard to local SmolLM (or any
OpenAI-compatible provider) with **curated codebase tools**.

## Pipeline

```
Browser  /assist
   │  same-origin proxy  /api/talon/*
   ▼
Next.js  web  (route.ts streams SSE unbuffered)
   │  cookie session forwarded
   ▼
talon-core :8000
   │  POST /llm/assist   (tools + tokens)
   │  POST /llm/stream   (tokens only)
   │  GET  /llm/tools    (catalog)
   │  GET  /llm/info
   ▼
LLM_PROVIDER=
  onnx  →  onnx-slm :8090/v1  (compose profile slm)
  ollama →  ollama :11434     (compose profile ollama)
  openai →  hosted / vLLM
  bedrock → AWS

Tool execution (assist only) stays IN Go against:
  store (runs/findings) · core skills · platform targets · health probes · MCP names
```

| Layer | Path / package | Role |
|-------|----------------|------|
| UI | `web/src/views/assist/AssistView.tsx` | Chat + tool cards |
| Client | `web/src/lib/api.ts` → `streamLLMAssist`, `listLLMTools` | SSE parse |
| Proxy | `web/src/app/api/talon/[...path]/route.ts` | Stream-through to core |
| Nav | `web/src/configs/navConfig.tsx` → **SLM Assist** | `/assist` |
| API | `internal/control/slm_assist.go` | Tool loop + SSE |
| Tools | `internal/control/slm_tools.go` | Curated safe catalog |
| Stream | `internal/llm/stream.go` | OpenAI SSE token client |
| Factory | `internal/llm/factory.go` | `LLM_PROVIDER=onnx` |
| Runtime | `onnx-slm/` | SmolLM2 / optional ORT |
| Compose | `docker-compose.yml` profile **`slm`** | `onnx-slm` service |

## Quick start (aligned)

```bash
# 1) Core stack
docker compose up -d --build

# 2) Local SLM runtime (optional profile — first boot downloads weights)
docker compose --profile slm up -d --build onnx-slm

# 3) Point core at onnx (restart core after env change)
# In .env:
#   LLM_PROVIDER=onnx
#   ONNX_BASE_URL=http://localhost:8090/v1
#   ONNX_MAIN_MODEL=smollm
docker compose up -d talon-core dashboard

# 4) Health
curl -s http://localhost:8000/health
curl -s http://localhost:8090/health | jq .
curl -s http://localhost:8000/health/services | jq '.services[]|select(.name=="onnx-slm")'
curl -s http://localhost:8000/llm/tools | jq '.count,.tools[].name'
curl -s http://localhost:8000/llm/info | jq .

# 5) Assist (auth cookie if enabled)
curl -N -X POST http://localhost:8000/llm/assist \
  -H 'Content-Type: application/json' \
  -d '{"messages":[{"role":"user","content":"List runs and check service health."}]}'

# 6) UI
# open http://localhost:3000/assist  (login if auth on)
```

## Tools (UI-safe, codebase-backed)

SmolLM **does not** run MSF/nmap from chat. Assist tools are read-only:

| Tool | Backing code |
|------|----------------|
| `list_runs` / `runs_summary` | `Store.PaginatedList` / `Summary` |
| `get_run_status` / `get_run_tools` / `get_findings` | `Store.Get` / `ToolLog` |
| `search_skills` / `get_skill` | `core.QuerySkills` / `GetSkill` |
| `list_agents` / `list_playbooks` | `core.ListAgents` / `ListPlaybooks` |
| `list_targets` | `Platform.ListTargets` |
| `service_health` | same probes as `/health/services` |
| `list_mcp_tools` | `mcpclient.Multi.Servers` (names only) |
| `intel_feed` | `Store.IntelFeed` |

**Protocol (onnx / SmolLM text):**

```text
TOOL_CALL {"name":"list_runs","arguments":{"limit":5}}
```

OpenAI / Ollama use native function calling against the **same** catalog.

SSE events from `/llm/assist`:

`meta` → (`round` → `token`* → `tool_start` → `tool_result`)* → `done` | `error`

## Provider matrix

| `LLM_PROVIDER` | Agent tool loop (runs) | Assist tools (`/llm/assist`) | Plain stream |
|----------------|------------------------|------------------------------|--------------|
| `bedrock` | yes | via Converse one-shot | one-shot |
| `openai` | yes | native tools | SSE |
| `ollama` | yes | native tools | one-shot today |
| `onnx` | no (chat/assist only) | **text TOOL_CALL** | SSE (fast) |

Hybrid: orchestrator on `openai`/`bedrock`, local assist on `onnx`.

## Compose profiles

| Profile | Service | When |
|---------|---------|------|
| *(default)* | core, dashboard, msf, arsenal, … | hosted LLM |
| `ollama` | `ollama` | GGUF local |
| `slm` | `onnx-slm` | SmolLM / ONNX Runtime |
| `vuln` | lab target | CVE lab |

```bash
docker compose up -d
docker compose --profile slm up -d --build onnx-slm
docker compose --profile ollama up -d
```

## Env (see `.env.example`)

| Variable | Default | Meaning |
|----------|---------|---------|
| `LLM_PROVIDER` | `bedrock` | set `onnx` for local SLM |
| `ONNX_BASE_URL` | `http://localhost:8090/v1` | onnx-slm OpenAI base |
| `ONNX_API_KEY` | `talon-local` | optional local |
| `ONNX_MAIN_MODEL` | `smollm` | alias |
| `SLM_MODEL_ID` | `HuggingFaceTB/SmolLM2-360M-Instruct` | HF weights |
| `SLM_BACKEND` | `auto` | `transformers` \| `onnx` |

## E2E checklist

- [ ] `GET /health` → ok
- [ ] `GET /llm/tools` → count ≥ 10, all `safe: true`
- [ ] `GET /llm/info` → `assist_path`, `tools_path`, `stream_path`
- [ ] `POST /llm/assist` with “list runs” → `tool_start`/`tool_result` or prose
- [ ] Dashboard **SLM Assist** (`/assist`) streams tokens + tool cards
- [ ] With `--profile slm`: `onnx-slm` online in `/health/services`
- [ ] Proxy: browser Network tab shows `/api/talon/llm/assist` `text/event-stream`

## Optional browser WASM

`web/src/lib/slm-wasm.ts` — Transformers.js / ORT-Web for offline 135M demos.
Production path remains Go assist (tools + store).
