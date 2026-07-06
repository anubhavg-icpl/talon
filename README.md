# Talon

![Talon](assets/hero-banner.webp)

<img src="assets/logo-wordmark.webp" width="72" alt="Talon logo"> <img src="assets/favicon.webp" width="48" alt="Talon favicon">

**An AI-driven, autonomous penetration-testing orchestration platform, built natively in Go.**

Talon takes a target (an IP, a suspected CVE, a service fingerprint) and drives a
full validation workflow end to end: recon confirms the service is really
there, exploit searches and fires known modules against it, post-exploit
proves compromise by interacting with the resulting session, and if nothing
prebuilt works, a code-generation fallback writes and runs a custom exploit
in a sandboxed container. A judge model checks the final transcript and
tells you, in plain `true`/`false`, whether the target was actually popped —
no hand-waving, no "the scan completed successfully" theater.

It is built from four small, single-purpose binaries — **Talon Core**, **Talon
Relay**, **Talon Arsenal**, and **Talon Strike** — plus a **Talon Forge**
code-generation sandbox, wired together over the Model Context Protocol
(MCP) and an LLM tool-calling loop. Every component ships as a single
static Go binary: no interpreter, no virtualenv, no dependency resolution
at deploy time — `docker build`, done.

![Authorized use only](assets/security-warning-banner.webp)

> **This is an offensive security tool.** It runs Nmap, SQLMap, Hydra,
> Metasploit exploits, and an LLM-driven custom-exploit generator against
> whatever IP you hand it. Read [Security Model & Responsible
> Use](#security-model--responsible-use) before you point it at anything.
> Only use Talon against systems you own or have explicit, documented,
> written authorization to test. Unauthorized access to computer systems is
> a crime in essentially every jurisdiction on Earth.

---

## Table of Contents

1. [What Talon Actually Does](#what-talon-actually-does)
2. [Design Philosophy](#design-philosophy)
3. [Architecture](#architecture)
4. [The Talon Product Family](#the-talon-product-family)
   - [Talon Core](#talon-core)
   - [Talon Relay](#talon-relay)
   - [Talon Arsenal](#talon-arsenal)
   - [Talon Strike](#talon-strike)
   - [Talon Forge](#talon-forge)
   - [Shared Infrastructure Packages](#shared-infrastructure-packages)
5. [The Request Lifecycle, Step by Step](#the-request-lifecycle-step-by-step)
6. [Human-in-the-Loop: The `nmap_scan` Gate](#human-in-the-loop-the-nmap_scan-gate)
7. [Installation & Prerequisites](#installation--prerequisites)
8. [Configuration Reference](#configuration-reference)
9. [Deploying with Docker Compose](#deploying-with-docker-compose)
10. [Bare-Metal / systemd Deployment](#bare-metal--systemd-deployment)
11. [HTTP API Reference (Talon Core)](#http-api-reference-talon-core)
12. [AMQP Worker Contract (Talon Relay)](#amqp-worker-contract-talon-relay)
13. [Talon Arsenal — Full Tool Catalog (151 Tools)](#talon-arsenal--full-tool-catalog-151-tools)
14. [Talon Strike — Metasploit Tool Catalog (12 Tools)](#talon-strike--metasploit-tool-catalog-12-tools)
15. [Talon Forge — The Codegen Sandbox](#talon-forge--the-codegen-sandbox)
16. [Subagent System Prompts](#subagent-system-prompts)
17. [Security Model & Responsible Use](#security-model--responsible-use)
18. [Known Limitations & Upgrade Paths](#known-limitations--upgrade-paths)
19. [Development Guide](#development-guide)
20. [Troubleshooting / FAQ](#troubleshooting--faq)
21. [Roadmap](#roadmap)
22. [Contributing](#contributing)
23. [License](#license)

---

## What Talon Actually Does

![The validation pipeline](assets/header-collage.webp)

Feed Talon a target:

```json
{
  "ip": "192.168.122.250",
  "cve_id": "CVE-2011-3556",
  "service_name": "Java RMI",
  "description": "RMI registry misconfiguration allowing remote class loading",
  "lhost": "192.168.122.176",
  "lport": 4444
}
```

and it runs a five-stage validation pipeline, entirely autonomously except
for one deliberate human checkpoint on the initial network scan:

1. **Recon** confirms the service/vulnerability is actually present on the
   target using Nmap, SMB enumeration, or Nuclei — whatever the orchestrator
   decides is appropriate. This is the one place a human has to click
   "approve" by default (see [Human-in-the-Loop](#human-in-the-loop-the-nmap_scan-gate)).
2. **Exploit** searches Metasploit's module database for something that
   matches, configures it with your `LHOST`/`LPORT`, and fires it.
3. **Post-Exploit** — if a session came back — interacts with it
   (`sysinfo`, `whoami`, `hostname`) to produce hard proof of compromise,
   not just "a session object exists."
4. **Forge (codegen) fallback** — if every prebuilt module failed, an LLM
   writes a custom Python exploit from scratch, runs it inside a sandboxed
   Docker container with a hard timeout, watches for `ModuleNotFoundError`
   and auto-installs the missing package, and retries up to twice before
   giving up.
5. **Judge** reads the entire final transcript and returns a single boolean:
   did this actually work? This exists specifically so the platform can't
   lie to you by declaring victory on a scan that merely *completed*.

The whole thing is exposed two ways: a synchronous-ish HTTP API
(`talon-core`, poll for status, resume on interrupt) for interactive use,
and an AMQP worker (`talon-relay`) for fire-and-forget integration into a
larger orchestration platform that already speaks RabbitMQ.

---

## Design Philosophy

Talon is built Go-native, end to end, on purpose:

- **A single static binary per component.** No interpreter, no venv, no
  "works on my machine" dependency resolution. `docker build`, done.
- **A type system that makes the tool-calling contract explicit.** Every
  MCP tool's schema, every AMQP message shape, every HTTP request/response
  body is a real Go struct with `json` tags, not a dict that might or might
  not have the key you expect at 2am when something's on fire.
- **No hardcoded credentials, anywhere.** Every secret (`MSF_PASSWORD`,
  `AMQP_URL`, AWS credentials) is env-var-only with no fallback default —
  the relevant binary refuses to start rather than silently running with a
  guessable secret.
- **Correctness over convenience in concurrent state.** The HTTP session
  store is mutex-protected; resume/interrupt handling is a real state
  machine, not a decision that's silently dropped on the floor.
- **MCP as the tool boundary, not a framework dependency** (see below) —
  Talon Arsenal and Talon Strike are independent MCP servers any MCP
  client can drive directly, with or without the orchestrator in the loop.

---

## Architecture

```
                                   ┌─────────────────────────────┐
                                   │      Operator / Platform     │
                                   │   (human, or upstream app)   │
                                   └──────────────┬───────────────┘
                                                   │
                        ┌──────────────────────────┴──────────────────────────┐
                        │                                                      │
                        ▼                                                      ▼
              ┌───────────────────┐                                 ┌───────────────────┐
              │    Talon Core      │   HTTP (start/status/resume)   │    Talon Relay     │
              │  (cmd/talon-core)  │◄────────── OR ─────────────────►│ (cmd/talon-relay) │
              │  net/http :8000    │                                 │  AMQP consumer     │
              └─────────┬──────────┘                                 └─────────┬──────────┘
                        │                                                      │
                        │              both build one of these:               │
                        └───────────────────────┬──────────────────────────────┘
                                                 ▼
                                    ┌─────────────────────────┐
                                    │   internal/core          │
                                    │   The Orchestrator        │
                                    │  (LLM tool-calling loop)  │
                                    └────────────┬──────────────┘
                                                 │ delegates to five synthetic
                                                 │ "delegate_*" subagent tools
                     ┌───────────────┬────────────┼────────────┬───────────────┐
                     ▼               ▼            ▼            ▼               ▼
               ┌──────────┐   ┌──────────┐  ┌───────────┐ ┌──────────┐  ┌──────────┐
               │  recon    │   │ exploit  │  │post_exploit│ │ codegen  │  │  report   │
               │ subagent  │   │ subagent │  │ subagent   │ │ subagent │  │ subagent  │
               └─────┬─────┘   └────┬─────┘  └─────┬──────┘ └────┬─────┘  └──────────┘
                     │              │              │              │
                     │   MCP stdio  │   MCP stdio  │   MCP stdio  │  direct Go call
                     ▼              ▼              ▼              ▼
          ┌────────────────────────────────────┐      ┌──────────────────────────┐
          │         Talon Arsenal                │      │      Talon Forge          │
          │      (cmd/talon-arsenal)             │      │   (internal/forge)        │
          │  MCP server, 151 tools, proxies to    │      │  Docker-sandboxed Python  │
          │  the Talon Arsenal Engine HTTP API    │      │  code generator/executor  │
          └──────────────────┬───────────────────┘      └──────────────────────────┘
                              │ HTTP
                              ▼
                 ┌─────────────────────────┐          ┌────────────────────────────┐
                 │  Talon Arsenal Engine    │          │      Talon Strike           │
                 │ (arsenal-engine/, self-  │          │   (cmd/talon-strike)        │
                 │  built tool backend)     │          │  MCP server, 12 tools,      │
                 └─────────────────────────┘          │  msgpack-RPC to msfrpcd     │
                                                        └──────────────┬─────────────┘
                     Exploit/post_exploit                              │ msgpack-RPC
                     subagents also call Talon                        │
                     Strike's 12 tools over MCP ───────────────────────┘
                                                                        ▼
                                                          ┌──────────────────────────┐
                                                          │   msfrpcd (kali-msf/)     │
                                                          │   real Metasploit Framework│
                                                          └──────────────────────────┘
```

![Talon architecture](assets/architecture-poster.webp)

Four Go binaries, one shared Docker image, two backing tool-execution
engines (Talon Arsenal Engine and the Metasploit RPC daemon), and one
piece of infrastructure you probably already have (RabbitMQ, only if you
run Talon Relay).

### Design principle: MCP as the tool boundary, not a framework dependency

Talon Arsenal and Talon Strike are *independent* MCP stdio servers. Talon
Core and Talon Relay spawn them as child processes and talk to them over
the Model Context Protocol — the same protocol Claude Desktop, Claude Code,
and any other MCP-aware client speaks. This means:

- You can point *any* MCP client at `talon-arsenal` or `talon-strike`
  directly and drive them by hand, with no orchestrator in the loop at all.
  They don't know or care that Talon Core exists.
- Swapping the LLM tool-calling loop for a different orchestration
  framework in the future only touches `internal/core` — the tool servers
  are untouched.
- Testing a single tool doesn't require standing up the whole platform.

---

## The Talon Product Family

### Talon Core

![Talon Core](assets/talon-core-badge.webp)

**Binary:** `cmd/talon-core` · **Package:** `internal/control` (HTTP layer) + `internal/core` (orchestrator)

The control plane. An HTTP server (stdlib `net/http`, Go 1.22+
pattern-based `ServeMux`, zero router dependencies) sitting in front of the
orchestrator. Start a run, poll its status, resume it past a human-in-the-loop
gate, pull the tool-call log.

Talon Core owns:
- The five HTTP routes (see [HTTP API Reference](#http-api-reference-talon-core))
- Spawning `talon-arsenal` and `talon-strike` as MCP stdio subprocesses
- Constructing three separate model clients (main orchestrator model,
  judge model, code-generation model)
- An in-memory, mutex-protected session store keyed by run ID

### Talon Relay

![Talon Relay](assets/talon-relay-badge.webp)

**Binary:** `cmd/talon-relay` · **Package:** `internal/relay`

The queue worker. Consumes `execute_agent_task` messages off RabbitMQ
(using `rabbitmq/amqp091-go`, a plain AMQP consumer/publisher rather than
any particular task-queue framework's wire protocol), runs the exact same
orchestrator as Talon Core, and publishes a `RUN_COMPLETED` event to
`agent.pentest.output` when it's done.

The one behavioral difference from Talon Core: Talon Relay has no human
attached, so if the orchestrator pauses on the `nmap_scan` HITL gate, Relay
auto-approves it and continues. This is a deliberate, documented
simplification — see [Known Limitations](#known-limitations--upgrade-paths).

### Talon Arsenal

![Talon Arsenal](assets/talon-arsenal-badge.webp)

**Binary:** `cmd/talon-arsenal` · **Package:** `internal/arsenal`

An MCP stdio server exposing **151 tools** — everything from `nmap_scan` to
`kube_hunter_scan` to `ai_generate_attack_suite` — by proxying to the Talon
Arsenal Engine's HTTP API, a self-built backing service (see
[Deploying with Docker Compose](#deploying-with-docker-compose) for how
it's built and run).

140 of the 151 tools are **mechanically generated** from a declarative Go
table (`internal/arsenal/tools_generated.go`): name, parameters, HTTP
method, endpoint, and body-field mapping, all data-driven — because
hand-transcribing 140 near-identical proxy functions is how you introduce
140 chances for a typo in an endpoint path. The remaining 11 tools (four
fully local playbook-generators, seven with small parameter transforms
like comma-split lists) are hand-written in `tools_manual.go`.

See the [full catalog](#talon-arsenal--full-tool-catalog-151-tools) below.

### Talon Strike

![Talon Strike](assets/talon-strike-badge.webp)

**Binary:** `cmd/talon-strike` · **Package:** `internal/strike`

An MCP stdio server exposing **12 tools** wrapping a Go implementation of
the Metasploit RPC wire protocol (msgpack-RPC over HTTPS to `msfrpcd`,
`POST /api/`, `Content-Type: binary/message-pack`).

The RPC client (`internal/strike/client.go`) implements the real
`msfrpcd` handshake directly: `auth.login` → mint a fresh UUIDv4 token →
`auth.token_add` → switch to the new token for every subsequent call, the
behavior a long-lived `msfrpcd` client is expected to follow. Module
execution goes exclusively through the `module.execute` RPC job path — a
slower console-based fallback (raw byte-stream prompt-regex matching
against `msf6 >`) was deliberately left out.

See the [full catalog](#talon-strike--metasploit-tool-catalog-12-tools) below.

### Talon Forge

![Talon Forge](assets/talon-forge-badge.webp)

**Package:** `internal/forge` (no standalone binary — linked into `talon-core` and `talon-relay`)

The `custom_exploit` tool the codegen subagent reaches for when nothing in
Talon Strike's module database works. Runs a bounded generator → executor →
feedback retry loop (2 generation attempts × 2 execution-fix attempts) that:

1. Asks an LLM to write Python code for the pseudocode instruction given.
2. Extracts the fenced ` ```python ` code block from the response.
3. Ensures a persistent Docker container (`talon_forge_sandbox`,
   `python:3.13-slim`, `--network host --memory 512m --cpus 1`) exists and
   is running.
4. Copies the code in, executes it with a hard timeout via the `docker` CLI
   (no Docker SDK dependency — the CLI is already a hard requirement, so
   shelling out is the honest choice, not a shortcut).
5. If the output contains a `ModuleNotFoundError`, auto-`pip install`s the
   missing package inside the container and retries the same code.
6. Otherwise asks the LLM a strict `true`/`false` question: did this
   accomplish the task? If not, and attempts remain, regenerates the code
   from scratch with the previous attempt's output as context.

A plain bounded Go loop drives this retry cycle — a 2×2-bounded retry
doesn't need a graph execution engine.

![Talon Forge codegen sandbox](assets/codegen-sandbox.webp)

### Shared Infrastructure Packages

Four packages are deliberately **not** part of the Talon product-naming
scheme, because they're generic plumbing every component shares, not a
capability with its own identity:

| Package | Purpose |
|---|---|
| `internal/config` | Env-var-driven config loaders for Talon Arsenal Engine, Metasploit, and AMQP connections. No hardcoded credential fallbacks — see [Configuration Reference](#configuration-reference). |
| `internal/llm` | `ChatModel`, the model-agnostic tool-calling interface every agent is built against, plus a concrete AWS Bedrock Converse API implementation. |
| `internal/mcpclient` | Wraps `mark3labs/mcp-go`'s stdio client to spawn and multiplex Talon Arsenal + Talon Strike as a single flat tool namespace. |

---

## The Request Lifecycle, Step by Step

![Request lifecycle](assets/sequence-diagram-base.webp)

Walking through a single `POST /input/start` call end to end:

```
 1. Client POSTs TargetRequest{ip, cve_id, service_name, description, lhost, lport}
    to talon-core.

 2. talon-core generates a run_id (uuid.NewString()), stores a Session in
    "initializing" state, launches a goroutine, and immediately returns
    {"run_id": "...", "message": "Agent execution started"}.

 3. [goroutine] Session -> "running". Orchestrator.Run(ctx, input) begins:

    a. Seed conversation: one HumanMessage describing the target + attacker
       context (LHOST/LPORT), plus the orchestrator system prompt.

    b. Loop: call the main model with five synthetic tools --
       delegate_recon, delegate_exploit, delegate_post_exploit,
       delegate_codegen, delegate_report -- each taking one free-form
       "instructions" string. The orchestrator LLM decides which to call
       and in what order (the system prompt strongly steers it toward
       recon -> exploit -> post_exploit -> [codegen if exploit fails] ->
       report, but nothing mechanically enforces that sequence beyond the
       prompt).

    c. When delegate_recon is called: spin up a NESTED tool-calling loop
       scoped to the recon subagent's system prompt and its three real MCP
       tools (nmap_scan, smbmap_scan, nuclei_scan). This is where the HITL
       gate lives -- see the next section.

    d. Each subagent's final text answer becomes the ToolResult fed back to
       the orchestrator's delegate_* call. The subagent's own tool calls
       never reach the orchestrator directly -- it only sees the summary.

    e. This repeats for delegate_exploit (11 tools: list_exploits,
       list_payloads, generate_payload, run_exploit, run_auxiliary_module,
       run_post_module, sqlmap_scan, arp_scan_discovery, hydra_attack,
       rustscan_fast_scan, responder_credential_harvest),
       delegate_post_exploit (3 tools: list_active_sessions,
       terminate_session, send_session_command), delegate_codegen (1 tool:
       Talon Forge's custom_exploit), and finally delegate_report (0 tools
       -- it just summarizes).

    f. Two safety nets run throughout: a shared 30-call budget across every
       tool call in the run (orchestrator delegates AND the real tool calls
       they trigger), and a context-length trim that drops all but the last
       3 tool-result messages once the running transcript exceeds 100,000
       characters.

 4. Orchestrator returns final text. talon-core calls the judge model with
    a fixed prompt asking for a strict true/false verdict on whether
    exploitation actually succeeded, based on the final transcript.

 5. Session -> "completed". Output, JudgeVerdict, and the full ToolLog are
    stored.

 6. Client polls GET /output/status/{run_id} until status is "completed"
    (or "awaiting_approval" -- see next section).
```

![Judge verdict](assets/judge-verdict.webp)

---

## Human-in-the-Loop: The `nmap_scan` Gate

![Human-in-the-loop gate](assets/hitl-gate.webp)

Exactly one tool call in the entire workflow requires human sign-off by
default: `nmap_scan`, inside the recon subagent's loop.

When the recon subagent's model decides to call `nmap_scan`:

```
 orchestrator.Run() ──► recon subagent loop ──► model requests nmap_scan
                                                        │
                                                        ▼
                                          STOP. Do not execute the tool.
                                          Return RunResult{
                                            Interrupted: true,
                                            Interrupt: &PendingInterrupt{
                                              ToolName: "nmap_scan",
                                              Args: {"target": "...", ...},
                                            },
                                          }
                                                        │
                                                        ▼
                                talon-core: Session -> "awaiting_approval"
                                            PendingInterrupt stored
                                                        │
                    ┌───────────────────────────────────┴────────────────────────────────┐
                    ▼                                                                     ▼
   GET /output/status/{run_id} returns                     POST /output/resume/{run_id}
   {"status": "awaiting_approval",                          {"decision": "approve"}
    "interrupt": {"tool_name": "nmap_scan",                 {"decision": "reject"}
                  "args": {...}}}                           {"decision": "edit",
                                                              "edited_args": {...}}
                                                                        │
                                                                        ▼
                                                    orchestrator.Resume(ctx, input, decision)
                                                    picks up EXACTLY where it paused:
                                                    - approve: runs nmap_scan with the
                                                      original args, continues the loop
                                                    - reject: feeds back an error
                                                      ToolResult, continues the loop
                                                    - edit: runs nmap_scan with
                                                      decision.EditedArgs instead
```

**Why keyed by `RunInput`, not a session token:** there is no external
checkpointer here. The orchestrator itself parks the full paused
transcript (both the outer orchestrator conversation and the inner recon
subagent conversation, at the exact point of interruption) in an in-memory
map keyed by the caller's `RunInput` struct — which is comparable
(no slices/maps as fields) specifically so it can be a map key. The caller
is required to pass the *identical* `RunInput` back to `Resume()`. See
[Known Limitations](#known-limitations--upgrade-paths) for the one sharp
edge this creates.

**Talon Relay has no human to ask**, so its orchestrator construction
auto-approves any `nmap_scan` interrupt it hits and continues immediately —
appropriate for a background queue worker with no UI attached, and
explicitly called out as a simplification rather than something silently
different from Talon Core's behavior.

---

## Installation & Prerequisites

![Terminal](assets/terminal-mockup.webp)

### Runtime dependencies

| Dependency | Why | Notes |
|---|---|---|
| Go 1.25+ | Building the four binaries | `go build ./...` from repo root |
| Docker + a running daemon | Talon Forge's sandbox; also how you'll deploy everything | The `docker` CLI must be on `PATH` (or the socket mounted into the container — see `docker-compose.yml`) |
| AWS account with Bedrock access (or a local Ollama server) | Every LLM call (orchestrator, judge, and code-generation models) | Standard AWS credential chain: env vars, `~/.aws/credentials`, or an IAM role -- or set `LLM_PROVIDER=ollama` for zero AWS dependency |
| A running `msfrpcd` (Metasploit RPC daemon) | Talon Strike's entire reason to exist | `kali-msf/Dockerfile` builds one from `kalilinux/kali-rolling` |
| A running Talon Arsenal Engine | Talon Arsenal proxies to it | `arsenal-engine/Dockerfile` builds one from `kalilinux/kali-rolling` -- self-built, part of this repo |
| RabbitMQ | Only if you're running Talon Relay | Any AMQP 0-9-1 broker works; bundled in `docker-compose.yml` |

### Building from source

```bash
git clone git@github.com:anubhavg-icpl/pentester2.git
cd pentester2
go build ./...
```

This produces nothing by itself — Go builds packages, not artifacts, unless
you build a specific `cmd/` target. Build all four binaries:

```bash
go build -o bin/talon-core    ./cmd/talon-core
go build -o bin/talon-relay   ./cmd/talon-relay
go build -o bin/talon-arsenal ./cmd/talon-arsenal
go build -o bin/talon-strike  ./cmd/talon-strike
```

Or let the multi-stage `Dockerfile` do it (this is what `docker-compose.yml`
actually uses):

```bash
docker build -t talon:latest .
```

### Verifying your build

```bash
go vet ./...          # should print nothing
gofmt -l .             # should print nothing (no files need formatting)
go build ./...         # should print nothing and exit 0
```

Every one of `talon-core`, `talon-relay`, and `talon-strike` is designed to
**fail fast with a clear message** if required configuration is missing —
this is deliberate and testable:

```bash
$ ./talon-relay
2026/07/06 04:29:06 talon-relay: AMQP_URL is not set

$ ./talon-strike
2026/07/06 04:29:06 talon-strike: failed to connect to Metasploit RPC: msfrpc: MSF_PASSWORD is required
```

If you see a binary silently defaulting to a guessable credential instead
of refusing to start, that's a regression — file an issue.

---

## Configuration Reference

Every setting is an environment variable. There are no config files, and
(deliberately) no hardcoded fallback for anything that is a credential.

### Shared across `talon-core` and `talon-relay`

| Variable | Default | Required | Description |
|---|---|---|---|
| `HEXSTRIKE_SERVER_URL` | `http://localhost:8888` | No | Base URL of the Arsenal Engine that Talon Arsenal proxies to. |
| `HEXSTRIKE_TIMEOUT` | `300` | No | Request timeout (seconds) for calls to the Arsenal Engine. |
| `MSF_PASSWORD` | *(none)* | **Yes** | Metasploit RPC password. No fallback — `talon-strike` refuses to start without it. |
| `MSF_SERVER` | `msf_rpc` | No | Hostname of the `msfrpcd` instance. |
| `MSF_PORT` | `5554` | No | Port of the `msfrpcd` instance. |
| `MSF_SSL` | `true` | No | Whether to use HTTPS for the Metasploit RPC connection. |
| `PAYLOAD_SAVE_DIR` | `$HOME/payloads` | No | Where generated Metasploit payloads are written on disk. |
| AWS credential chain | *(standard AWS SDK resolution)* | **Yes** | `AWS_ACCESS_KEY_ID`/`AWS_SECRET_ACCESS_KEY`/`AWS_SESSION_TOKEN`, `~/.aws/credentials`, or an IAM role — resolved automatically by `aws-sdk-go-v2`'s default config loader. |
| `AWS_REGION` | *(SDK default, typically unset)* | No | Passed through in `docker-compose.yml`; `talon-core` currently hardcodes `us-east-1` as a Go constant rather than reading this — see note below. |

### `talon-core`-specific

| Variable | Default | Description |
|---|---|---|
| `HEXSTRIKE_MCP_PATH` | sibling binary named `talon-arsenal`, resolved relative to `talon-core`'s own executable path | Override the Talon Arsenal binary location. |
| `METASPLOIT_MCP_PATH` | sibling binary named `talon-strike` | Override the Talon Strike binary location. |

`talon-core`'s Bedrock region and model IDs are currently compiled-in
constants (`bedrockRegion = "us-east-1"`, main model
`qwen.qwen3-vl-235b-a22b`, judge model `openai.gpt-oss-120b-1:0`, code model
`us.meta.llama4-maverick-17b-instruct-v1:0`) rather than environment
variables. This is an intentional asymmetry with `talon-relay` (below) —
`talon-core`'s model selection was ported verbatim from `final.py`'s
hardcoded model construction, whereas `talon-relay` was built with env-var
overrides from day one. If you need `talon-core`'s models configurable at
runtime, that's a small, well-scoped change (promote the four `const`s in
`cmd/talon-core/main.go` to `getenv()` calls matching `talon-relay`'s
pattern) — flagged here rather than silently "fixed" without being asked.

### `talon-relay`-specific

| Variable | Default | Description |
|---|---|---|
| `AMQP_URL` | *(none)* | **Required.** Full AMQP connection string, e.g. `amqp://user:pass@rabbitmq:5672/`. No `guest:guest@localhost` fallback — the worker refuses to start without it. |
| `BEDROCK_REGION` | `us-east-1` | AWS region for all three Bedrock model clients. |
| `AGENT_MODEL_ID` | `qwen.qwen3-vl-235b-a22b` | Main orchestrator model. |
| `JUDGE_MODEL_ID` | `openai.gpt-oss-120b-1:0` | Judge model for the final true/false verdict. |
| `CODE_MODEL_ID` | `us.meta.llama4-maverick-17b-instruct-v1:0` | Talon Forge's code-generation model. |
| `HEXSTRIKE_MCP_BIN` | `talon-arsenal` (resolved via `PATH`) | Override the Talon Arsenal binary name/path. |
| `METASPLOIT_MCP_BIN` | `talon-strike` (resolved via `PATH`) | Override the Talon Strike binary name/path. |

### `kali-msf` (the Metasploit RPC container)

| Variable | Default | Description |
|---|---|---|
| `MSF_PORT` | `5554` | Port `msfrpcd` listens on. |
| `MSF_PASSWORD` | *(none)* | **Required.** The container's entrypoint checks for this and refuses to start `msfrpcd` if it's unset — no `network_msf` hardcoded default like the original image had. |

---

## Deploying with Docker Compose

`docker-compose.yml` at the repo root defines six services:

```yaml
services:
  metasploit:      # kali-msf/Dockerfile -- real Metasploit Framework + msfrpcd
  arsenal-engine:  # arsenal-engine/Dockerfile -- self-built Talon Arsenal Engine
  rabbitmq:        # bundled broker for talon-relay
  ollama:          # optional local LLM runtime (LLM_PROVIDER=ollama)
  talon-core:      # this repo's Dockerfile, command: ["/app/talon-core"]
  talon-relay:     # this repo's Dockerfile, command: ["/app/talon-relay"]
```

All services run with `network_mode: host` (Metasploit's RPC daemon and
the Arsenal Engine are expected to be reachable at `localhost:<port>` from
inside the Talon containers — this matches how the tooling underneath
actually expects to bind, and simplifies the container graph at the cost of
host-network isolation. If you need real network segmentation, replace this
with a user-defined bridge network and adjust `MSF_SERVER`/
`HEXSTRIKE_SERVER_URL` to the bridge's service DNS names).

### Minimal walkthrough

```bash
# 1. Required secrets -- docker compose will refuse to start without these
export MSF_PASSWORD='choose-a-real-password-not-this-one'
export AMQP_URL='amqp://talon:s3cr3t@your-rabbitmq-host:5672/'

# 2. Optional overrides
export MSF_PORT=5554
export HEXSTRIKE_PORT=8888
export AWS_REGION=us-east-1

# 3. AWS credentials -- either export the standard triple, or mount
#    ~/.aws into the containers, or run this on an EC2 instance with an
#    IAM role attached. docker-compose.yml does not currently forward
#    AWS_ACCESS_KEY_ID/AWS_SECRET_ACCESS_KEY -- add them to the
#    `environment:` blocks for talon-core/talon-relay if you're not using
#    an instance role or a mounted credentials file.

# 4. Build and start
docker compose build
docker compose up -d

# 5. Watch it come up
docker compose logs -f talon-core talon-relay
```

### Verifying the compose file without starting anything

```bash
docker compose config --quiet
```

This will fail loudly (`required variable MSF_PASSWORD is missing a value`)
if you haven't set the required secrets — that's the compose file's
`:?` interpolation guard doing its job, not a bug.

### The Docker build context

`Dockerfile` builds all four binaries in a single multi-stage image
(`golang:1.25-alpine` → `alpine:3.20` with `docker-cli` installed for Talon
Forge's sandbox calls). `.dockerignore` excludes `.git`, `.claude*`,
`kali-msf/`, `msf-data/`, `msf_payloads/`, `database.yml`, and
`docker-compose.yml` from the build context — none of those are needed to
build the Go binaries, and there's no reason to let a `docker build .` at
the repo root accidentally slurp local tool state into an image layer.

### The `/var/run/docker.sock` mount

Both `talon-core` and `talon-relay` mount the host's Docker socket:

```yaml
volumes:
  - /var/run/docker.sock:/var/run/docker.sock
```

This is required for Talon Forge — the codegen sandbox tool shells out to
the `docker` CLI *inside* the container, which needs a real Docker daemon
to talk to. This is sibling-container access, not Docker-in-Docker, and it
means anything running inside `talon-core`/`talon-relay` has, transitively,
root-equivalent control over the host's Docker daemon. Understand that
blast radius before deploying this on shared infrastructure — see
[Security Model](#security-model--responsible-use).

---

## Bare-Metal / systemd Deployment

If you'd rather not containerize the Go binaries (Metasploit and the
Arsenal Engine still need their own environments regardless):

```ini
# /etc/systemd/system/talon-core.service
[Unit]
Description=Talon Core - pentest orchestration control plane
After=network-online.target
Wants=network-online.target

[Service]
Type=simple
User=talon
Group=talon
EnvironmentFile=/etc/talon/core.env
WorkingDirectory=/opt/talon
ExecStart=/opt/talon/bin/talon-core
Restart=on-failure
RestartSec=5
# Talon Core spawns talon-arsenal/talon-strike as children -- make sure
# they're findable via HEXSTRIKE_MCP_PATH/METASPLOIT_MCP_PATH in the
# EnvironmentFile, or sitting next to this binary in WorkingDirectory.

[Install]
WantedBy=multi-user.target
```

```ini
# /etc/talon/core.env
HEXSTRIKE_SERVER_URL=http://127.0.0.1:8888
MSF_SERVER=127.0.0.1
MSF_PORT=5554
MSF_PASSWORD=change-me
MSF_SSL=false
AWS_REGION=us-east-1
HEXSTRIKE_MCP_PATH=/opt/talon/bin/talon-arsenal
METASPLOIT_MCP_PATH=/opt/talon/bin/talon-strike
```

Mirror the same pattern for `talon-relay.service` with its own env file
(add `AMQP_URL`, drop the MCP path overrides in favor of
`HEXSTRIKE_MCP_BIN`/`METASPLOIT_MCP_BIN` if the binaries aren't on `PATH`).

The `talon-relay` user needs to be in the `docker` group (or have socket
access some other way) for Talon Forge to work — same blast-radius caveat
as the container deployment applies here, arguably more directly since
there's no container boundary at all in this mode.

---

## HTTP API Reference (Talon Core)

![Dashboard mockup](assets/dashboard-mockup.webp)

Base URL: `http://<host>:8000` (hardcoded port, `talon-core` listens on
`0.0.0.0:8000`).

### `POST /input/start`

Start a new validation run.

**Request body:**

```json
{
  "ip": "192.168.122.250",
  "cve_id": "CVE-2011-3556",
  "service_name": "Java RMI",
  "description": "RMI registry misconfiguration allowing remote class loading",
  "lhost": "192.168.122.176",
  "lport": 4444
}
```

`lhost` defaults to `192.168.122.176` and `lport` defaults to `4444` if
omitted or zero — matching `final.py`'s `Context` defaults.

**Response (`200`):**

```json
{
  "run_id": "3fa85f64-5717-4562-b3fc-2c963f66afa6",
  "message": "Agent execution started"
}
```

**Response (`400`):** invalid JSON body.

```bash
curl -X POST http://localhost:8000/input/start \
  -H 'Content-Type: application/json' \
  -d '{"ip":"10.0.0.5","cve_id":"CVE-2021-44228","service_name":"log4j","description":"JNDI lookup RCE","lhost":"10.0.0.1","lport":4444}'
```

### `GET /output/status/{run_id}`

Poll a run's current state.

```bash
curl http://localhost:8000/output/status/3fa85f64-5717-4562-b3fc-2c963f66afa6
```

**Response, running:**

```json
{"status": "running", "output": "", "interrupt": null}
```

**Response, awaiting a human decision:**

```json
{
  "status": "awaiting_approval",
  "output": "",
  "interrupt": {
    "tool_name": "nmap_scan",
    "args": {"target": "10.0.0.5", "scan_type": "-sV", "ports": "", "additional_args": ""}
  }
}
```

**Response, completed:**

```json
{
  "status": "completed",
  "output": "Exploitation confirmed. Session 1 established via exploit/multi/http/log4j_rce...",
  "interrupt": null
}
```

**Response, unknown run:**

```json
{"status": "not_found"}
```

### `POST /output/resume/{run_id}`

Resolve a pending interrupt.

```bash
# Approve exactly as requested
curl -X POST http://localhost:8000/output/resume/3fa85f64-5717-4562-b3fc-2c963f66afa6 \
  -H 'Content-Type: application/json' \
  -d '{"decision":"approve"}'

# Reject -- the subagent gets told the human declined and has to route around it
curl -X POST http://localhost:8000/output/resume/3fa85f64-5717-4562-b3fc-2c963f66afa6 \
  -H 'Content-Type: application/json' \
  -d '{"decision":"reject"}'

# Edit -- run the SAME tool with DIFFERENT arguments than the model requested
curl -X POST http://localhost:8000/output/resume/3fa85f64-5717-4562-b3fc-2c963f66afa6 \
  -H 'Content-Type: application/json' \
  -d '{"decision":"edit","edited_args":{"target":"10.0.0.5","scan_type":"-sT","ports":"80,443","additional_args":"-Pn"}}'
```

`decision` must be exactly one of `approve`, `reject`, `edit`
(case-insensitive). Anything else is a `400`. Calling this with no pending
interrupt on the run is also a `400` (`"No pending interrupt"`).

**Response (`200`):**

```json
{"message": "Decision received, resuming orchestrator..."}
```

### `GET /monitor/traces/{run_id}`

Full message history for a run (each entry is a stringified
`llm.Message`).

```bash
curl http://localhost:8000/monitor/traces/3fa85f64-5717-4562-b3fc-2c963f66afa6
```

```json
{"history": ["...", "...", "..."]}
```

### `GET /monitor/tools?run_id={run_id}`

The tool-call log for a specific run.

> **Design note:** `ToolCallRecord`s are scoped to the `RunResult` they
> belong to, not a global tracker — so this endpoint requires a `run_id`
> query parameter and returns that run's log specifically, not a
> cross-run firehose.

```bash
curl 'http://localhost:8000/monitor/tools?run_id=3fa85f64-5717-4562-b3fc-2c963f66afa6'
```

```json
{
  "tool_log": [
    {"index": 0, "tool_name": "nmap_scan", "args": {"target": "10.0.0.5", "scan_type": "-sV"}, "output": "{\"success\":true,...}"},
    {"index": 1, "tool_name": "run_exploit", "args": {"module_name": "exploit/multi/http/log4j_rce", ...}, "output": "{\"job_id\":3,...}"}
  ]
}
```

Missing/unknown `run_id` on either monitoring endpoint returns `404` with
`{"detail": "Run not found"}`.

---

## AMQP Worker Contract (Talon Relay)

Talon Relay is a plain AMQP 0-9-1 consumer/publisher, not tied to any
particular task-queue framework's wire protocol — the message envelope is
a flat JSON object Talon Relay controls entirely.

### Consumes: `execute_agent_task`

```json
{
  "run_id": "run-abc123",
  "project_id": "proj-xyz789",
  "agent_name": "pentest-validator",
  "agent_inputs": {
    "target_ip": "10.0.0.5",
    "cve_id": "CVE-2021-44228",
    "lhost": "10.0.0.1",
    "lport": 4444,
    "description": "JNDI lookup RCE"
  }
}
```

Talon Relay runs the orchestrator to completion (auto-approving the
`nmap_scan` gate, see [above](#human-in-the-loop-the-nmap_scan-gate)),
then publishes:

### Publishes: `agent.pentest.output`

```json
{
  "event_type": "RUN_COMPLETED",
  "run_id": "run-abc123",
  "project_id": "proj-xyz789",
  "agent_name": "pentest-validator",
  "overall_status": "completed",
  "result": {
    "summary": "Exploitation confirmed. Session 1 established...",
    "raw_output": "<full stringified RunResult>"
  },
  "timestamp": "2026-07-06T04:45:00Z"
}
```

Note: `agent.pentest.output` is spelled with `agent.`, not `core.` or
`talon.` — this queue name is an external integration contract with
whatever platform is consuming completion events, and intentionally was
**not** renamed as part of the Talon rebrand. (It briefly *was*
accidentally renamed by an overly broad find-and-replace during that
rebrand and had to be reverted — a reminder that identifier renames and
string-literal protocol constants need very different handling.)

---

## Talon Arsenal — Full Tool Catalog (151 Tools)

![Recon workflow](assets/recon-workflow.webp)

Every tool below is registered on the Talon Arsenal MCP server, organized
by category. 140 are mechanically generated proxy calls; 11 (marked below)
involve either a small client-side parameter transform or fully local
playbook generation.

#### Core Network Scanning Tools (3 tools)

| Tool | Description |
|---|---|
| `nmap_scan` | Execute an enhanced Nmap scan against a target with real-time logging. |
| `gobuster_scan` | Execute Gobuster to find directories, DNS subdomains, or virtual hosts with enhanced logging. |
| `nuclei_scan` | Execute Nuclei vulnerability scanner with enhanced logging and real-time progress. |

#### Cloud Security Tools (17 tools)

| Tool | Description |
|---|---|
| `prowler_scan` | Execute Prowler for comprehensive cloud security assessment. |
| `trivy_scan` | Execute Trivy for container and filesystem vulnerability scanning. |
| `scout_suite_assessment` | Execute Scout Suite for multi-cloud security assessment. |
| `cloudmapper_analysis` | Execute CloudMapper for AWS network visualization and security analysis. |
| `pacu_exploitation` | Execute Pacu for AWS exploitation framework. |
| `kube_hunter_scan` | Execute kube-hunter for Kubernetes penetration testing. |
| `kube_bench_cis` | Execute kube-bench for CIS Kubernetes benchmark checks. |
| `docker_bench_security_scan` | Execute Docker Bench for Security for Docker security assessment. |
| `clair_vulnerability_scan` | Execute Clair for container vulnerability analysis. |
| `falco_runtime_monitoring` | Execute Falco for runtime security monitoring. |
| `checkov_iac_scan` | Execute Checkov for infrastructure as code security scanning. |
| `terrascan_iac_scan` | Execute Terrascan for infrastructure as code security scanning. |
| `create_file` | Create a file with specified content on the Arsenal Engine. |
| `modify_file` | Modify an existing file on the Arsenal Engine. |
| `delete_file` | Delete a file or directory on the Arsenal Engine. |
| `list_files` | List files in a directory on the Arsenal Engine. |
| `generate_payload` | Generate large payloads for testing and exploitation. |

#### Python Environment Management (2 tools)

| Tool | Description |
|---|---|
| `install_python_package` | Install a Python package in a virtual environment on the Arsenal Engine. |
| `execute_python_script` | Execute a Python script in a virtual environment on the Arsenal Engine. |

#### Additional Security Tools From Original Implementation (74 tools)

| Tool | Description |
|---|---|
| `dirb_scan` | Execute Dirb for directory brute forcing with enhanced logging. |
| `nikto_scan` | Execute Nikto web vulnerability scanner with enhanced logging. |
| `sqlmap_scan` | Execute SQLMap for SQL injection testing with enhanced logging. |
| `metasploit_run` | Execute a Metasploit module with enhanced logging (via the Arsenal Engine, distinct from Talon Strike's dedicated Metasploit RPC tools). |
| `hydra_attack` | Execute Hydra for password brute forcing with enhanced logging. |
| `john_crack` | Execute John the Ripper for password cracking with enhanced logging. |
| `wpscan_analyze` | Execute WPScan for WordPress vulnerability scanning with enhanced logging. |
| `enum4linux_scan` | Execute Enum4linux for SMB enumeration with enhanced logging. |
| `ffuf_scan` | Execute FFuf for web fuzzing with enhanced logging. |
| `netexec_scan` | Execute NetExec (formerly CrackMapExec) for network enumeration with enhanced logging. |
| `amass_scan` | Execute Amass for subdomain enumeration with enhanced logging. |
| `hashcat_crack` | Execute Hashcat for advanced password cracking with enhanced logging. |
| `subfinder_scan` | Execute Subfinder for passive subdomain enumeration with enhanced logging. |
| `smbmap_scan` | Execute SMBMap for SMB share enumeration with enhanced logging. |
| `rustscan_fast_scan` | Execute Rustscan for ultra-fast port scanning with enhanced logging. |
| `masscan_high_speed` | Execute Masscan for high-speed Internet-scale port scanning with intelligent rate limiting. |
| `nmap_advanced_scan` | Execute advanced Nmap scans with custom NSE scripts and optimized timing. |
| `autorecon_comprehensive` | Execute AutoRecon for comprehensive automated reconnaissance. |
| `enum4linux_ng_advanced` | Execute Enum4linux-ng for advanced SMB enumeration with enhanced logging. |
| `rpcclient_enumeration` | Execute rpcclient for RPC enumeration with enhanced logging. |
| `nbtscan_netbios` | Execute nbtscan for NetBIOS name scanning with enhanced logging. |
| `arp_scan_discovery` | Execute arp-scan for network discovery with enhanced logging. |
| `responder_credential_harvest` | Execute Responder for credential harvesting with enhanced logging. |
| `volatility_analyze` | Execute Volatility for memory forensics analysis with enhanced logging. |
| `msfvenom_generate` | Execute MSFVenom for payload generation with enhanced logging. |
| `gdb_analyze` | Execute GDB for binary analysis and debugging with enhanced logging. |
| `radare2_analyze` | Execute Radare2 for binary analysis and reverse engineering with enhanced logging. |
| `binwalk_analyze` | Execute Binwalk for firmware and file analysis with enhanced logging. |
| `ropgadget_search` | Search for ROP gadgets in a binary using ROPgadget with enhanced logging. |
| `checksec_analyze` | Check security features of a binary with enhanced logging. |
| `xxd_hexdump` | Create a hex dump of a file using xxd with enhanced logging. |
| `strings_extract` | Extract strings from a binary file with enhanced logging. |
| `objdump_analyze` | Analyze a binary using objdump with enhanced logging. |
| `ghidra_analysis` | Execute Ghidra for advanced binary analysis and reverse engineering. |
| `pwntools_exploit` | Execute Pwntools for exploit development and automation. |
| `one_gadget_search` | Execute one_gadget to find one-shot RCE gadgets in libc. |
| `libc_database_lookup` | Execute libc-database for libc identification and offset lookup. |
| `gdb_peda_debug` | Execute GDB with PEDA for enhanced debugging and exploitation. |
| `angr_symbolic_execution` | Execute angr for symbolic execution and binary analysis. |
| `ropper_gadget_search` | Execute ropper for advanced ROP/JOP gadget searching. |
| `pwninit_setup` | Execute pwninit for CTF binary exploitation setup. |
| `feroxbuster_scan` | Execute Feroxbuster for recursive content discovery with enhanced logging. |
| `dotdotpwn_scan` | Execute DotDotPwn for directory traversal testing with enhanced logging. |
| `xsser_scan` | Execute XSSer for XSS vulnerability testing with enhanced logging. |
| `wfuzz_scan` | Execute Wfuzz for web application fuzzing with enhanced logging. |
| `dirsearch_scan` | Execute Dirsearch for advanced directory and file discovery with enhanced logging. |
| `katana_crawl` | Execute Katana for next-generation crawling and spidering with enhanced logging. |
| `gau_discovery` | Execute Gau (Get All URLs) for URL discovery from multiple sources with enhanced logging. |
| `waybackurls_discovery` | Execute Waybackurls for historical URL discovery with enhanced logging. |
| `arjun_parameter_discovery` | Execute Arjun for HTTP parameter discovery with enhanced logging. |
| `paramspider_mining` | Execute ParamSpider for parameter mining from web archives with enhanced logging. |
| `x8_parameter_discovery` | Execute x8 for hidden parameter discovery with enhanced logging. |
| `jaeles_vulnerability_scan` | Execute Jaeles for advanced vulnerability scanning with custom signatures. |
| `dalfox_xss_scan` | Execute Dalfox for advanced XSS vulnerability scanning with enhanced logging. |
| `httpx_probe` | Execute HTTPx for HTTP probing with enhanced logging. |
| `anew_data_processing` | Execute anew for appending new lines to files (useful for data processing). |
| `qsreplace_parameter_replacement` | Execute qsreplace for query string parameter replacement. |
| `uro_url_filtering` | Execute uro for filtering out similar URLs. |
| `ai_generate_payload` | Generate AI-powered contextual payloads for security testing. |
| `ai_test_payload` | Test a generated payload against the target with AI analysis. |
| `ai_generate_attack_suite`* | Generate a comprehensive attack suite spanning multiple payload types, aggregating `ai_generate_payload` across each requested attack type. |
| `api_fuzzer`* | Advanced API endpoint fuzzing with intelligent parameter discovery. |
| `graphql_scanner` | Advanced GraphQL security scanning and introspection. |
| `jwt_analyzer` | Advanced JWT token analysis and vulnerability testing. |
| `api_schema_analyzer` | Analyze API schemas and identify potential security issues. |
| `comprehensive_api_audit`* | Comprehensive API security audit combining endpoint fuzzing, schema analysis, JWT analysis, and GraphQL scanning into a single report. |
| `volatility3_analyze` | Execute Volatility3 for advanced memory forensics with enhanced logging. |
| `foremost_carving` | Execute Foremost for file carving with enhanced logging. |
| `steghide_analysis` | Execute Steghide for steganography analysis with enhanced logging. |
| `exiftool_extract` | Execute ExifTool for metadata extraction with enhanced logging. |
| `hashpump_attack` | Execute HashPump for hash length extension attacks with enhanced logging. |
| `hakrawler_crawl` | Execute Hakrawler for web endpoint discovery with enhanced logging. |
| `httpx_probe` (2nd registration) | Execute HTTPx for HTTP probing with enhanced logging. |
| `paramspider_discovery` | Execute ParamSpider for parameter discovery with enhanced logging. |

<sub>* Hand-ported in `tools_manual.go`, not table-generated — see [Talon Arsenal](#talon-arsenal).</sub>

#### Advanced Web Security Tools Continued (26 tools)

| Tool | Description |
|---|---|
| `burpsuite_scan` | Execute Burp Suite with enhanced logging. |
| `zap_scan` | Execute OWASP ZAP with enhanced logging. |
| `arjun_scan` | Execute Arjun for parameter discovery with enhanced logging. |
| `wafw00f_scan` | Execute wafw00f to identify and fingerprint WAF products with enhanced logging. |
| `fierce_scan` | Execute fierce for DNS reconnaissance with enhanced logging. |
| `dnsenum_scan` | Execute dnsenum for DNS enumeration with enhanced logging. |
| `autorecon_scan` | Execute AutoRecon for comprehensive target enumeration with full parameter support. |
| `server_health` | Check the health status of the Arsenal Engine. |
| `get_cache_stats` | Get cache statistics from the Arsenal Engine. |
| `clear_cache` | Clear the cache on the Arsenal Engine. |
| `get_telemetry` | Get system telemetry from the Arsenal Engine. |
| `list_active_processes` | List all active processes on the Arsenal Engine. |
| `get_process_status` | Get the status of a specific process by PID. |
| `terminate_process` | Terminate a specific running process by PID. |
| `pause_process` | Pause a specific running process by PID. |
| `resume_process` | Resume a paused process by PID. |
| `get_process_dashboard` | Get an enhanced process dashboard with visual status indicators. |
| `execute_command` | Execute an arbitrary command on the Arsenal Engine with enhanced logging. |
| `monitor_cve_feeds` | Monitor CVE databases for new vulnerabilities with AI analysis. |
| `generate_exploit_from_cve` | Generate working exploits from CVE information using AI-powered analysis. |
| `discover_attack_chains`* | Discover multi-stage attack chains for target software with vulnerability correlation (attack depth clamped 1-5). |
| `research_zero_day_opportunities` | Automated zero-day vulnerability research using AI analysis and pattern recognition. |
| `correlate_threat_intelligence`* | Correlate threat intelligence (IOCs, CVEs, domains) across multiple sources. |
| `advanced_payload_generation` | Generate advanced payloads with AI-powered evasion techniques and contextual adaptation. |
| `vulnerability_intelligence_dashboard` | Get a comprehensive vulnerability intelligence dashboard with latest threats and trends. |
| `threat_hunting_assistant`* | AI-powered threat hunting playbook generator (detection queries, threat scenarios, investigation steps), with optional threat-intel correlation when indicators are supplied. |

<sub>* Hand-ported in `tools_manual.go`.</sub>

#### Enhanced Visual Output Tools (5 tools)

| Tool | Description |
|---|---|
| `get_live_dashboard` | Get a live dashboard showing all active processes with enhanced visual formatting. |
| `create_vulnerability_report` | Create a vulnerability report with severity-based styling and visual indicators. |
| `format_tool_output_visual` | Format tool output with visual styling, syntax highlighting, and structure. |
| `create_scan_summary` | Create a comprehensive scan summary report with visual formatting. |
| `display_system_metrics` | Display current system metrics and performance indicators with visual formatting. |

#### Intelligent Decision Engine Tools (8 tools)

| Tool | Description |
|---|---|
| `analyze_target_intelligence` | Analyze a target using AI-powered intelligence to create a comprehensive profile. |
| `select_optimal_tools_ai` | Use AI to select optimal security tools based on target analysis and testing objective. |
| `optimize_tool_parameters_ai`* | Use AI to optimize tool parameters based on target profile and free-form JSON context. |
| `create_attack_chain_ai` | Create an intelligent attack chain using AI-driven tool sequencing and optimization. |
| `intelligent_smart_scan` | Execute an intelligent scan using AI-driven tool selection and parameter optimization. |
| `detect_technologies_ai` | Use AI to detect technologies and provide technology-specific testing recommendations. |
| `ai_reconnaissance_workflow` | Execute an AI-driven reconnaissance workflow with intelligent tool chaining. |
| `ai_vulnerability_assessment` | Perform an AI-driven vulnerability assessment with intelligent prioritization. |

<sub>* Hand-ported in `tools_manual.go`.</sub>

#### Bug Bounty Hunting Specialized Workflows (16 tools)

| Tool | Description |
|---|---|
| `bugbounty_reconnaissance_workflow`* | Create a comprehensive reconnaissance workflow for bug bounty hunting, scoped by in-/out-of-scope domain lists. |
| `bugbounty_vulnerability_hunting`* | Create a vulnerability hunting workflow prioritized by impact and bounty potential. |
| `bugbounty_business_logic_testing` | Create a business logic testing workflow for advanced bug bounty hunting. |
| `bugbounty_osint_gathering` | Create an OSINT gathering workflow for bug bounty reconnaissance. |
| `bugbounty_file_upload_testing` | Create a file upload vulnerability testing workflow with bypass techniques. |
| `bugbounty_comprehensive_assessment`* | Create a comprehensive bug bounty assessment combining all specialized workflows. |
| `bugbounty_authentication_bypass_testing`* | Generate an authentication bypass testing workflow (form/JWT/OAuth/SAML-specific technique lists) for bug bounty hunting. |
| `http_framework_test` | Enhanced HTTP testing framework (Burp Suite alternative) for comprehensive web security testing. |
| `browser_agent_inspect` | AI-powered browser agent for comprehensive web application inspection and security analysis. |
| `http_set_rules` | Set match/replace rules used to rewrite parts of URL/query/headers/body before sending. |
| `http_set_scope` | Define an in-scope host (and optionally subdomains) so out-of-scope requests are skipped. |
| `http_repeater` | Send a crafted request (Burp Repeater equivalent) with a `request_spec` (`url`, `method`, `headers`, `cookies`, `data`). |
| `http_intruder` | Simple Intruder (sniper) fuzzing — iterates payloads over each parameter individually. |
| `burpsuite_alternative_scan` | Comprehensive Burp Suite alternative combining the HTTP framework and browser agent for complete web security testing. |
| `error_handling_statistics` | Get intelligent error-handling system statistics and recent error patterns. |
| `test_error_recovery` | Test the intelligent error recovery system with simulated failures. |

<sub>* Hand-ported in `tools_manual.go`.</sub>

---

## Talon Strike — Metasploit Tool Catalog (12 Tools)

![Talon Strike RPC connection](assets/msf-rpc-connection.webp)

| Tool | Description |
|---|---|
| `list_exploits` | List available Metasploit exploit modules, optionally filtered by a search term. |
| `list_payloads` | List available Metasploit payload modules, optionally filtered by platform and/or architecture. |
| `generate_payload` | Generate a Metasploit payload (msfvenom-style, via `module.execute` with `modtype=payload`) with the given options and save it to `PAYLOAD_SAVE_DIR`. |
| `run_exploit` | Run an exploit module against a target via `module.execute`, then poll `session.list` for up to 60 seconds (2-second interval) for a resulting session. |
| `run_post_module` | Run a post-exploitation module against an active session. |
| `run_auxiliary_module` | Run an auxiliary module (scanners, fuzzers, protocol-specific checks). |
| `list_active_sessions` | List all active Metasploit sessions (both shell and Meterpreter), including type and target info. |
| `send_session_command` | Send a command to an active session (dispatches to `session.shell_write`/`session.shell_read` or `session.meterpreter_write`/`session.meterpreter_read` depending on session type) and return the output. |
| `list_listeners` | List running handler jobs (`exploit/multi/handler`) currently waiting for inbound connections. |
| `start_listener` | Start a new Metasploit handler as a background job, configured with the given payload/LHOST/LPORT. |
| `stop_job` | Stop a running Metasploit job (a handler or any other backgrounded module). |
| `terminate_session` | Forcefully terminate an active Metasploit session via `session.stop`. |

**Option parsing** on every tool that accepts a Metasploit module's options
tolerates three input shapes: a real JSON object, a
`"key=value,key2=value2"` string (the common mistake format people
actually send), or nothing. Values are type-guessed afterward — an
all-digit string becomes an integer, `"true"`/`"false"` (case-insensitive)
becomes a boolean.

**What's deliberately not implemented:** a console-based execution
fallback that would read raw bytes off an interactive `msf6 >` console and
match a compiled regex against the prompt to detect command completion.
Every module execution in Talon Strike goes through the `module.execute`
RPC job path instead — simpler, and the only path the orchestrator's
subagents ever actually need.

---

## Talon Forge — The Codegen Sandbox

Talon Forge is not exposed as a standalone binary or MCP server — it's
linked directly into `talon-core` and `talon-relay` as the one tool
available to the `codegen` subagent, `custom_exploit`.

**Container lifecycle:** a single persistent container (`python:3.13-slim`,
`--network host --memory 512m --cpus 1`) is created once and reused across
calls within a process's lifetime — not recreated per execution, which
would be needlessly slow for an LLM that might need three or four attempts
per task.

**Retry bounds:** 2 code-generation attempts × 2 execution-fix attempts.
This is intentionally small. Talon Forge is a *fallback* for when Talon
Strike's module database comes up empty, not a general-purpose "keep
trying until it works" loop that could run indefinitely against a target.

**Auto-remediation:** the only failure Talon Forge attempts to fix
automatically is a missing Python package (`ModuleNotFoundError`), which it
resolves with a `pip install` inside the sandbox and a same-code retry.
Every other failure mode (syntax errors, logic errors, target refused
connection) burns a generation attempt and asks the LLM to write different
code with the previous failure as context.

**Security boundary:** the sandbox is a container with host networking and
capped resources, **not** a hardened isolation boundary. `--network host`
means generated code can reach anything the host can reach. This is
appropriate for a tool whose entire purpose is executing an LLM-written
exploit against an authorized target on the same network as the operator —
it is not appropriate to expose this to untrusted input. See
[Security Model](#security-model--responsible-use).

---

## Subagent System Prompts

These are the actual strings compiled into `internal/core/prompts.go` —
reproduced here so you don't have to go spelunking in source to understand
what's steering each subagent's behavior.

**Orchestrator:**
> You are a senior penetration testing orchestrator. You will receive
> target infrastructure details and attacker context (LHOST, LPORT). Your
> workflow MUST follow these strict sequential steps. DO NOT skip steps or
> stop early. STEP 1 (Recon): Task the 'recon' subagent to verify if the
> specified service or vulnerability exists on the target. STEP 2
> (Exploit): Once 'recon' confirms the target, immediately task the
> 'exploit' subagent... STEP 3 (Post-Exploit)... STEP 4 (Fallback to
> CodeGen)... STEP 5 (Report)...

**Recon:**
> You are a recon specialist. Your primary job is to verify target
> services and vulnerabilities. Use your scanning tools and report back
> ONLY factual findings based on tool outputs. Be concise and accurate.

**Exploit:**
> You are an exploit specialist utilizing pre-built modules (like
> Metasploit). 1. SEARCH: Find relevant modules for the target service or
> CVE. 2. CONFIGURE: When generating payloads or setting up listeners, you
> MUST use the LHOST and LPORT provided by the orchestrator. 3. EXECUTE:
> You must execute the chosen module immediately. Do not stop after
> searching. 4. VERIFY: Read stdout/stderr. If it explicitly states 'No
> session created' or 'Exploit failed', it is a FAILURE. Move to the next
> module.

**Post-Exploit:**
> You are a post-exploitation specialist. Your job is to interact with
> established sessions to retrieve proof of compromise. 1. Identify the
> active session ID. 2. Execute commands to identify the system and user
> (e.g., 'sysinfo', 'hostname', or 'whoami'). 3. Return the raw tool output
> containing the proof back to the orchestrator.

**Codegen:**
> You are a senior exploit developer. You are invoked when standard tools
> fail. You will be given recon data, target details, LHOST, and LPORT.
> Your job is to use the 'custom_exploit' tool to generate and execute a
> custom Python script (e.g., reverse shells, RCE exploits) against the
> target. Ensure the generated code properly utilizes the provided LHOST
> and LPORT for any reverse connections.

**Report:**
> You are a report writer. Generate a final validation report. Only
> generate this if an exploit actually succeeded. Summarize the IP, the
> CVE tested, the module used, and the proof of success.

---

## Security Model & Responsible Use

**Talon is a weapon.** It automates the exact chain of actions —
reconnaissance, exploit selection, exploitation, session interaction,
custom exploit generation — that separates "running a vulnerability
scanner" from "actually compromising a system." Everything below is not
boilerplate; it describes real capabilities this codebase has.

- **Only test systems you are explicitly authorized to test.** Written
  authorization, defined scope, defined rules of engagement. If you don't
  have a signed statement of work or a bug bounty program's published scope
  covering the target, you don't have authorization.
- **The `nmap_scan` human-in-the-loop gate is a safety-relevant control, not
  a UX nicety.** It exists so a human confirms the target before the first
  network packet leaves the box. Talon Relay's auto-approve-on-interrupt
  behavior removes that control for anything run through the queue —
  understand that before wiring Talon Relay into an automated pipeline that
  accepts targets from anywhere other than a pre-vetted, pre-scoped
  intake process.
- **`/var/run/docker.sock` access is root-equivalent host control.** Any
  code that can reach the Docker socket can mount the host filesystem into
  a new container and read/write anything. Talon Forge needs this to run
  generated exploits sandboxed *from the target*, not sandboxed *from your
  infrastructure*. Treat the host running `talon-core`/`talon-relay` as
  equivalent in trust level to the Metasploit and Arsenal Engine boxes it
  talks to — not as an isolated control-plane node.
- **The LLM writes and runs its own exploit code with `--network host`.**
  There is no allowlist on what Talon Forge's generated Python can connect
  to besides the target you gave it in the prompt. A misconfigured or
  compromised orchestrator model could generate code that reaches anywhere
  the host can reach. Run this on infrastructure scoped to the engagement,
  not on a shared jump box with access to unrelated internal systems.
- **Credentials are environment-variable-only by design**, specifically so
  they don't end up committed, logged by a shell history, or baked into an
  image layer. `MSF_PASSWORD` and `AMQP_URL` have no fallback and the
  binaries refuse to start without them — don't "fix" that by adding one
  back.
- **This is not a managed SaaS product with abuse detection.** There is no
  rate limiting, no scope validation against a published bug bounty
  program, no built-in check that the target IP is one you're allowed to
  touch. All of that responsibility sits with whoever operates this
  software.

If you found this repository without the context of an authorized
engagement already in hand: this is not a "try it on a random IP and see
what happens" tool. Set up a deliberately vulnerable target in a lab you
control (Metasploitable, a HackTheBox/TryHackMe VPN lab, a local VM you
built specifically to be broken) before pointing this at anything else.

---

## Known Limitations & Upgrade Paths

Every simplification below is a deliberate, documented choice — not an
oversight. Each is marked in the source with a `// ponytail:` comment
naming the ceiling and the upgrade condition.

| Limitation | Where | Upgrade when |
|---|---|---|
| Interrupted-run sessions are keyed by `RunInput` (a comparable struct of target IP/CVE/service/description/context), not an explicit session ID. Two concurrent runs against an *identical* target collide on the same map key. | `internal/core/orchestrator.go` | If concurrent identical-target runs become a real scenario (e.g. Talon Core serving overlapping requests for the same target), add an explicit session/thread ID field to `RunInput`. |
| Metasploit's console-based execution path (raw byte-stream reads, regex prompt matching against `msf6 >`) isn't implemented, in favor of the `module.execute` RPC job path. | `internal/strike/tools.go` | If a future module type genuinely needs interactive console semantics that `module.execute` can't express. |
| `send_session_command` doesn't support a nested "shell-within-Meterpreter" state machine. It's a flat write-then-poll-read against the session's native read/write RPC methods. | `internal/strike/tools.go` | If a workflow genuinely needs to drop a Meterpreter session into an OS shell mid-session and run commands there. |
| Talon Relay auto-approves the `nmap_scan` HITL gate rather than blocking. | `cmd/talon-relay/main.go` (orchestrator construction) | If queue-driven runs need the same human checkpoint as the HTTP API — would require a second AMQP round-trip (publish "awaiting approval," consume a decision message) rather than a synchronous HTTP resume. |
| `talon-core`'s Bedrock region and three Bedrock model IDs are compiled-in constants (the `LLM_PROVIDER`/`OLLAMA_*` switch is env-configurable, but per-model Bedrock overrides aren't, unlike `talon-relay`'s `AGENT_MODEL_ID`/`JUDGE_MODEL_ID`/`CODE_MODEL_ID`/`BEDROCK_REGION`). | `cmd/talon-core/main.go` | If you need to change Bedrock models without a rebuild — promote the `const`s to `getenv()` calls matching `talon-relay`'s existing pattern. |
| `GET /monitor/tools` requires a `run_id` query parameter and serves per-run logs rather than a single cross-run global tracker. | `internal/control/server.go` | Documented here as an API-shape detail, not a defect. |
| The `agent.pentest.output` AMQP queue name is a hardcoded string constant, not derived from any config. | `internal/relay/worker.go` | If the platform this integrates with ever needs a configurable output queue name. |

---

## Development Guide

### Project layout

```
cmd/
  talon-core/       main() for the HTTP control plane
  talon-relay/      main() for the AMQP queue worker
  talon-arsenal/    main() for the Talon Arsenal Engine proxy MCP server
  talon-strike/     main() for the Metasploit RPC MCP server
internal/
  core/             the orchestrator: subagent loop, HITL gate, judge, tool tracker
  control/          HTTP handlers + in-memory session store for talon-core
  relay/            AMQP consumer/publisher for talon-relay
  arsenal/          151-tool Talon Arsenal Engine proxy (tools_generated.go + tools_manual.go)
  strike/           Metasploit msgpack-RPC client + 12 MCP tools
  forge/            Docker-sandboxed Python codegen loop
  llm/              ChatModel interface + Bedrock Converse implementation
  mcpclient/        stdio MCP client wrapper (spawns/multiplexes arsenal + strike)
  config/           env-var config loaders, no hardcoded secrets
```

### Adding a new Talon Arsenal tool

If the underlying Arsenal Engine gains a new endpoint that's a pure
`params -> JSON body -> POST` proxy (the overwhelming majority of Arsenal's
tools are exactly this shape), add an entry to the `generatedTools` table
in `internal/arsenal/tools_generated.go`:

```go
{
    Name:        "my_new_tool",
    Description: "One-line description shown to the orchestrator LLM.",
    Params: []paramSpec{
        {Name: "target", Kind: "string", Required: true, Default: nil},
        {Name: "timeout", Kind: "integer", Required: false, Default: 30},
    },
    HTTPMethod:   "POST",
    Endpoint:     "api/tools/my_new_tool",
    EndpointVars: []string{},
    BodyFields:   []bodyField{{ParamName: "target", JSONKey: "target"}, {ParamName: "timeout", JSONKey: "timeout"}},
    StaticFields: map[string]any{},
},
```

`Register()` in `internal/arsenal/register.go` picks it up automatically —
no new handler code needed. If the tool needs a parameter transform (a
comma-split list, an integer clamp, a JSON-string parse) or does something
fully local rather than proxying, add it to `registerManualTools()` in
`tools_manual.go` instead, following the pattern of the existing 11.

### Adding a new Talon Strike tool

Metasploit tools are hand-written in `internal/strike/tools.go` since each
one has genuinely different RPC call shapes (`module.execute` vs.
`session.list` vs. `job.stop`, different option-parsing needs). Follow the
existing 12 as a template; the msgpack-RPC primitives (`Call`, `ParseOptionsGracefully`)
live in `client.go` and `modules.go`.

### Running the test suite

There isn't a large one yet — this codebase is young. `go vet ./...` and
`gofmt -l .` are both wired into the verification you should run before any
commit (see [Installation](#installation--prerequisites)). If you're adding
non-trivial logic (a new retry loop, a new parsing branch), add a focused
`_test.go` next to it — table-driven tests for the `paramSpec`/`bodyField`
mapping logic in `internal/arsenal/register.go` and the option-parsing
logic in `internal/strike/modules.go` would be the highest-value additions
if you're looking for a place to start.

### Adding new tool table entries

There's no maintained code-generation pipeline for `tools_generated.go` in
this repository -- new tool additions go directly into
`tools_generated.go`/`tools_manual.go` by hand, following the existing
entries as a template (see [Adding a new Talon Arsenal
tool](#adding-a-new-talon-arsenal-tool) above).

---

## Package Reference

A godoc-style listing of every exported symbol, package by package, pulled
directly from source. If you're integrating against Talon programmatically
(embedding `internal/core` in your own binary, writing a new HTTP front
end against the same orchestrator, etc.) this is the actual contract
surface — everything not listed here is unexported and considered private
implementation detail that can change without notice.

### `internal/config`

Env-var-driven configuration, no hardcoded credential fallbacks.

```go
type Context struct { LHOST string; LPORT int }
func DefaultContext() Context

type HexstrikeConfig struct { ServerURL string; Timeout int }
func LoadHexstrikeConfig() HexstrikeConfig

type MSFConfig struct {
    Password, Server, Port string
    SSL bool
    PayloadSaveDir string
}
func LoadMSFConfig() MSFConfig

type AMQPConfig struct { URL string }
func LoadAMQPConfig() AMQPConfig
```

### `internal/llm`

The model-agnostic tool-calling contract, plus the one concrete
implementation (Bedrock).

```go
type Role string
const (RoleSystem Role = "system"; RoleUser Role = "user"; RoleAssistant Role = "assistant"; RoleTool Role = "tool")

type ToolCall struct { ID, Name string; Args map[string]any }
type ToolResult struct { ToolCallID, Name, Content string; IsError bool }
type Message struct {
    Role Role
    Text string
    ToolCalls []ToolCall
    ToolResults []ToolResult
}
func UserMessage(text string) Message
func AssistantText(text string) Message
func ToolResultMessage(r ToolResult) Message

type ToolSpec struct { Name, Description string; InputSchema map[string]any }

type ChatModel interface {
    Converse(ctx context.Context, systemPrompt string, messages []Message, tools []ToolSpec) (Message, error)
}

// Bedrock implements ChatModel via the AWS Bedrock Converse API.
type Bedrock struct { /* unexported */ }
func NewBedrock(ctx context.Context, modelID, region string, temperature float32, maxTokens int32) (*Bedrock, error)
func (b *Bedrock) Converse(ctx context.Context, systemPrompt string, messages []Message, tools []ToolSpec) (Message, error)
```

Any type satisfying `ChatModel` can be dropped into `core.New(...)` in place
of a real `*llm.Bedrock` — this is the seam you'd use to write a test double
or plug in a different model provider.

### `internal/mcpclient`

Wraps `mark3labs/mcp-go`'s stdio client to spawn and flatten multiple MCP
servers into one tool namespace.

```go
type ServerSpec struct { Name, Command string; Args []string }

type Multi struct { /* unexported */ }
func NewMulti(ctx context.Context, specs []ServerSpec) (*Multi, error)
func (m *Multi) Tools() []llm.ToolSpec
func (m *Multi) Subset(names ...string) []llm.ToolSpec
func (m *Multi) Call(ctx context.Context, name string, args map[string]any) (string, error)
func (m *Multi) Close() error
```

### `internal/core`

The orchestrator. This is the biggest and most important contract in the
codebase.

```go
type RunInput struct {
    TargetIP, CVEID, ServiceName, Description string
    Context config.Context
}

type ToolCallRecord struct {
    Index int
    ToolName string
    Args map[string]any
    Output string
}

type RunResult struct {
    FinalMessage string
    ToolLog []ToolCallRecord
    JudgeVerdict bool
    Interrupted bool
    Interrupt *PendingInterrupt
}

type PendingInterrupt struct { ToolName string; Args map[string]any }
type Decision struct { Type string; EditedArgs map[string]any }

type CodegenTool interface {
    Name() string
    Description() string
    Call(ctx context.Context, query string) (string, error)
}

type Orchestrator struct { /* unexported, includes interrupt/resume session state */ }
func New(model llm.ChatModel, judge llm.ChatModel, tools *mcpclient.Multi, codegen CodegenTool) *Orchestrator
func (o *Orchestrator) Run(ctx context.Context, input RunInput) (RunResult, error)
func (o *Orchestrator) Resume(ctx context.Context, input RunInput, decision Decision) (RunResult, error)
```

Internally (not exported, but worth knowing about if you're reading the
source): `maxToolCalls = 30` is the shared tool-call budget;
`contextTrimTrigger = 100_000` / `contextTrimKeep = 3` govern the
length-based context trim; `tracker` is the per-run counter+log threaded
through every subagent call; `subagentSpec` is the (model, prompt, tools,
HITL gate, executor) tuple for each of the five delegate targets, returned
by the unexported `subagentConfig(delegateName string)` switch.

### `internal/control`

HTTP layer for `talon-core`.

```go
type Session struct {
    Status string
    Output string
    PendingInterrupt *core.PendingInterrupt
    RunInput core.RunInput
    History []string
    ToolLog []core.ToolCallRecord
}

type Store struct { /* unexported, sync.RWMutex-protected */ }
func NewStore() *Store
func (s *Store) Create(runID string, input core.RunInput)
func (s *Store) Get(runID string) (Session, bool)
func (s *Store) SetStatus(runID, status string)
func (s *Store) SetResult(runID string, result core.RunResult)
func (s *Store) SetError(runID string, err error)
func (s *Store) ClearInterrupt(runID string)
func (s *Store) ToolLog(runID string) ([]core.ToolCallRecord, bool)

type Server struct { /* unexported */ }
func NewServer(orch *core.Orchestrator, store *Store) *Server
func (s *Server) Mux() *http.ServeMux
```

### `internal/relay`

AMQP consumer/publisher for `talon-relay`.

```go
type AgentInputs struct {
    TargetIP, CVEID, LHOST, Description string
    LPORT int
}
type AgentTask struct {
    RunID, ProjectID, AgentName string
    AgentInputs AgentInputs
}
type CompletionResult struct { Summary, RawOutput string }
type CompletionPayload struct {
    EventType, RunID, ProjectID, AgentName, OverallStatus string
    Result CompletionResult
    Timestamp string
}

type Worker struct { /* unexported */ }
func NewWorker(url string) (*Worker, error)
func (w *Worker) Consume(ctx context.Context, orchestrator *core.Orchestrator) error
func (w *Worker) Close() error

func PublishCompletion(ctx context.Context, ch *amqp.Channel, payload CompletionPayload) error
```

### `internal/arsenal`

The 151-tool Talon Arsenal Engine proxy.

```go
type Client struct { /* unexported */ }
func NewClient(baseURL string, timeout time.Duration) *Client
func (c *Client) Get(endpoint string, params map[string]any) map[string]any
func (c *Client) Post(endpoint string, data map[string]any) map[string]any
func (c *Client) ExecuteCommand(command string, useCache bool) map[string]any
func (c *Client) CheckHealth() map[string]any

func Register(srv *server.MCPServer, client *Client)
```

`Register` is the only entry point most callers need — it wires all 140
generated tools plus the 11 manual ones onto an `mcp-go` server in one
call. The `toolSpec`/`paramSpec`/`bodyField` types backing the generated
table are unexported by design — they're a code-generation artifact, not a
stable public shape.

### `internal/strike`

The Metasploit msgpack-RPC client and its 12 MCP tools.

```go
type Client struct { /* unexported */ }
func NewClient(ctx context.Context, cfg config.MSFConfig) (*Client, error)
func (c *Client) Call(ctx context.Context, method string, params ...any) (map[string]any, error)
func (c *Client) CoreVersion(ctx context.Context) (map[string]any, error)

func (c *Client) ListExploits(ctx context.Context) ([]string, error)
func (c *Client) ListPayloads(ctx context.Context) ([]string, error)
func (c *Client) ListAuxiliary(ctx context.Context) ([]string, error)
func (c *Client) ListPost(ctx context.Context) ([]string, error)
func (c *Client) Execute(ctx context.Context, modtype, modname string, options map[string]any) (map[string]any, error)

func (c *Client) ListSessions(ctx context.Context) (map[string]any, error)
func (c *Client) ReadSession(ctx context.Context, id, sessionType string) (map[string]any, error)
func (c *Client) WriteSession(ctx context.Context, id, sessionType, data string) (map[string]any, error)
func (c *Client) StopSession(ctx context.Context, id string) (map[string]any, error)

func (c *Client) ListJobs(ctx context.Context) (map[string]any, error)
func (c *Client) JobInfo(ctx context.Context, id string) (map[string]any, error)
func (c *Client) StopJob(ctx context.Context, id string) (map[string]any, error)

func ParseOptionsGracefully(v any) (map[string]any, error)

func Register(srv *server.MCPServer, c *Client)
```

### `internal/forge`

The codegen sandbox and its `CodegenTool` implementation.

```go
// Coder runs the full generate/execute/feedback retry loop and returns the
// final observed output text.
func Coder(ctx context.Context, model llm.ChatModel, pseudocode string) (string, error)

type CustomExploitTool struct { Model llm.ChatModel }
func NewCustomExploitTool(model llm.ChatModel) *CustomExploitTool
func (t *CustomExploitTool) Name() string
func (t *CustomExploitTool) Description() string
func (t *CustomExploitTool) Call(ctx context.Context, query string) (string, error)

type ExecResult struct { /* stdout, stderr, exit code */ }
func EnsureContainer(ctx context.Context) error
func ExecPython(ctx context.Context, code string, timeout time.Duration) (stdout, stderr string, exitCode int, err error)
func PipInstall(ctx context.Context, pkg string) (ExecResult, error)
```

`*CustomExploitTool` is the concrete type satisfying `core.CodegenTool`
structurally — `internal/forge` deliberately does not import
`internal/core` (that would create an import cycle, since `core` is what
invokes `CodegenTool`), so the interface satisfaction is duck-typed and
wired up explicitly in `cmd/talon-core/main.go` and
`cmd/talon-relay/main.go`.

---

## Troubleshooting / FAQ

**`talon-strike` exits immediately with "MSF_PASSWORD is required."**
That's correct behavior, not a bug — set `MSF_PASSWORD` to match whatever
password `msfrpcd` was started with. There is no default.

**`talon-core` logs `mcpclient: initialize metasploit: transport error: transport closed` and then exits.**
`talon-core` spawns `talon-strike` as a child process; if `talon-strike`
itself fails to start (most commonly: `MSF_PASSWORD` unset, or `msfrpcd`
unreachable at `MSF_SERVER:MSF_PORT`), its stdio pipe closes immediately,
which surfaces in `talon-core`'s MCP client as a transport error. Check
`talon-strike`'s own startup requirements first.

**How do I know if the orchestrator is stuck on the `nmap_scan` gate versus actually still running recon?**
Poll `GET /output/status/{run_id}`. `"status": "awaiting_approval"` with a
non-null `"interrupt"` means it's parked waiting on
`POST /output/resume/{run_id}`. `"status": "running"` with no interrupt
means it's actively working.

**The judge model said the exploit succeeded but I don't see a session.**
The judge only reads the final text transcript — it has no independent
verification of session state. If the exploit subagent's summary claimed
success without a real session existing, that's an orchestrator/subagent
prompt-following failure, not the judge lying; check
`GET /monitor/tools?run_id=...` for the actual `run_exploit`/
`list_active_sessions` tool outputs to see what really happened.

**Can I run Talon Arsenal or Talon Strike standalone, without Talon Core?**
Yes — that's the point of building them as independent MCP stdio servers.
Point any MCP client at the built binary (`./talon-arsenal`,
`./talon-strike --transport stdio`) and drive the tools directly.

**Why does `talon-relay` need Docker socket access if it's "just" a queue worker?**
Because the orchestrator it runs is the exact same one `talon-core` runs,
including the codegen subagent's Talon Forge tool. If your deployment never
exercises the codegen fallback in practice, you could omit the socket mount
and accept that any run reaching that path will fail — but that's a
per-deployment risk tradeoff, not something the code enforces for you.

---

## Roadmap

Not commitments, just the honest list of what's visibly unfinished, pulled
directly from the [Known Limitations](#known-limitations--upgrade-paths)
table above plus a few structural items:

- [ ] Explicit session IDs for interrupt/resume state, replacing the
      `RunInput`-as-map-key approach.
- [ ] `talon-core` env-var-configurable Bedrock region/model IDs, matching
      `talon-relay`.
- [ ] An AMQP-native human-in-the-loop path for `talon-relay`, so
      queue-driven runs don't have to auto-approve the `nmap_scan` gate.
- [ ] A real automated test suite (`internal/arsenal` param-mapping,
      `internal/strike` option-parsing, `internal/core` orchestrator-loop
      behavior with a mocked `ChatModel`).
- [ ] A persistent session store (currently pure in-memory — a
      `talon-core` restart loses every in-flight run's state).
- [ ] Structured logging (currently `log.Printf` with string prefixes) with
      a real level/field-based logger, for anyone running this at scale.

---

## Contributing

![Go native](assets/go-native-milestone.webp)

This is a young, fast-moving codebase — expect some rough edges (see
[Known Limitations](#known-limitations--upgrade-paths) for the ones
already known about). Before opening a PR:

1. `go build ./...`, `go vet ./...`, and `gofmt -l .` must all be clean.
2. If you're touching `internal/strike`, cite the specific Metasploit RPC
   method(s) involved and how you verified their behavior — the wire
   protocol here is reverse-engineered against real `msfrpcd` behavior,
   not guessed, and that bar should hold for any changes.
3. If you're touching `internal/arsenal`'s generated table, follow the
   existing entries' shape closely — it's hand-maintained (see
   [Adding a new Talon Arsenal tool](#adding-a-new-talon-arsenal-tool)).
4. Fix bugs you find rather than working around them, and say so in the
   commit message.

## License

No license file has been added to this repository yet. Until one is,
standard copyright applies — all rights reserved by the repository owner.
If you intend to use, distribute, or build on this code, reach out to the
repository owner first, or wait for an explicit license to be added.

---

*Talon is a Go-native pentest orchestration platform. Every architectural
claim in this document is grounded in the actual source in this
repository — component boundaries, tool counts, environment variables,
and API shapes were all read directly from `internal/`, `cmd/`, and
`docker-compose.yml`, not written from memory or aspiration.*

