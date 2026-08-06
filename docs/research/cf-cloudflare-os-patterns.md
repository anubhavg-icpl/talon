# Cloudflare OS Patterns → Talon (Go) Adaptation Notes

Study of `/tmp/cf-cloudflare-os` (a "Gadgets Workshop" platform: AI agents in a capability-sandboxed environment). This doc distills the 7 requested patterns, the key TypeScript types, and how each maps to Go for a pentest automation platform.

---

## 1. Gatekeepers (Capability-Based Tool Security)

**Source:** `packages/workshop-shared/src/gatekeeper.ts` (the kernel API), `packages/gatekeeper-github/src/github.ts`, `packages/gatekeeper-mcp/src/mcp.ts`.

### Core idea
Every external integration is a **Gatekeeper** — a separate worker (process) that owns credentials and exposes a *capability* to the sandbox, never the raw secret. The model is a 3-level capability chain, each level narrower than the one before:

```
GatekeeperVendor        (per-service: "github", "mcp")  — connectAccount(), describe(), OAuth flow
   └─ GatekeeperUser    (per-user connected account)    — getGatekeeperClassFor(url), getVerifier(), revoke()
        └─ Gatekeeper<Session>  (per-resource binding)  — startSession(queue), addObserver(), applyAction()
             └─ Session   (handed INTO the sandbox/agent) — the only thing untrusted code can touch
```

Key invariant: the **Session** is the *only* object the sandboxed agent ever sees. Everything above it is privileged and held server-side.

### Key types
```ts
// A vendor declares what resource URL patterns it can connect (e.g. "https://github.com/*")
type SupportedResource = { urlPattern: string; title: string; grantable?: boolean }

// Per-resource binding inside a workspace — created via getGatekeeperClassFor(url)
interface Gatekeeper<Session> extends DurableObject {
  describe(): Promise<ResourceDescription>          // pre-grant metadata
  startSession(queue: RpcStub<ApprovalQueue>): Promise<Session>  // the capability handed out
  addObserver(id, verifier): Promise<void>          // sharing/access verification (§6)
  applyAction(action: number): Promise<void>        // execute an approved write (§2)
  rejectAction(action: number): Promise<void>
}

// GatekeeperUser proves who you are to the gatekeeper for sharing checks (§6)
interface GatekeeperUser {
  getGatekeeperClassFor(url: string): Promise<{class, resource}>
  getVerifier(): Promise<Fetcher<GatekeeperUserVerifier>>
}
```

The **endpoint/nonce security** is explicit: `connectAccount()` URLs *must* embed a cryptographic nonce to prevent replay (documented as a SECURITY note on the interface).

### Go adaptation for Talon
This maps directly to a **pentest tool plugin interface**. Each tool (nmap, sqlmap, nuclei, a custom MCP server) is a gatekeeper:

```go
// pkg/capability/gatekeeper.go
type Vendor interface {
    Describe() VendorDescription
    ConnectAccount(cb ConnectCallback) (authURL string, err error) // OAuth/manual cred flow
    SupportedResources() []SupportedResource                       // e.g. urlPattern per scan target type
    NewToolBinding(target string) (Gatekeeper, error)              // = getGatekeeperClassFor
}

// A bound tool instance scoped to one target — the privileged side
type Gatekeeper interface {
    Describe() ResourceDescription
    StartSession(q ApprovalQueue) (Session, error)  // hands the narrow capability to the agent
    ApplyAction(actionID int64) error               // run an approved dangerous action
    RejectAction(actionID int64) error
    AddObserver(observerID string, v Verifier) error
    RemoveObserver(observerID string) error
}

// The ONLY interface the LLM agent loop ever receives — least privilege
type Session interface { /* tool-specific methods, e.g. Scan(), Enumerate() */ }
```

Each Go plugin = one `Vendor` impl + one `Gatekeeper` impl + one `Session` impl. Run tools as **separate processes** (or goroutines with explicit capability passing) so a compromised agent can't reach credentials directly — match CF's "separate Worker per gatekeeper."

---

## 2. Human-in-the-Loop Approval System

**Source:** `gatekeeper.ts` (`ApprovalQueue`, `ObservationAuthorizer`, `ActionDescription`), `mcp-shared/src/action-store.ts`, `mcp-shared/src/session.ts`.

### Core idea — the Read/Write split
**Every** operation through a gatekeeper goes through a queue. The queue has two modes:

- **Observation (read):** `authorizeObservation()` is *synchronous* — it must succeed *before* data returns to the agent. Used for recon/scanning reads. Can carry `prohibitAllSharing` (data too sensitive to share) or `excludeObservers` (per-user ACLs, §6).
- **Action (write/side-effect):** `submitAction()` is *asynchronous* — returns immediately with a pending state; the agent polls `getActionResult(id)`. The action is **never executed until a human approves** (or an auto-approval rule fires). The gatekeeper may *simulate* the pending action so the agent can keep reasoning.

### Key types
```ts
interface ObservationAuthorizer extends RpcTarget {
  authorizeObservation(description: ObservationDescription): Promise<void>  // throws to block
}
interface ApprovalQueue extends ObservationAuthorizer {
  submitAction(action: number, description: ActionDescription): Promise<void>  // async, returns fast
  bindHook(controller, callback, description): Promise<void>                   // event registration
}

type ObservationDescription = {
  title: string; description: string
  prohibitAllSharing?: boolean   // lock down workspace after this read
  excludeObservers?: string[]    // block if named observers still have access (§6)
}

type ActionDescription = {
  title: string; description: string       // Markdown shown to the approver
  implementsRevert: boolean                 // can it be undone?
  awaitDecision?: boolean                   // agent should pause until decided (if not simulated)
  autoApprovable?: boolean                  // gatekeeper says this specific call is safe to auto-run
  actionKind?: { tag: string; label: string } // stable tag for auto-approval rule matching
}
```

### The state machine (`action-store.ts`) — this is the gold pattern for Talon
SQLite table `mcp_actions` with a 5-state lifecycle and **at-most-once** claim semantics:

```
pending ──apply──> applying ──success──> applied
   │                  │
   │reject            │ failure/timeout
   ▼                  ▼
 rejected          failed (retryable: true/false)
```

Critical guarantees (directly relevant to pentest where actions hit real systems):
1. **Claim before dispatch:** row flips `pending → applying` and is persisted *before* the network call. A second concurrent `applyAction` finds it `applying` and is refused — **prevents double-execution**.
2. **Recover from crash:** on store init, any row left in `applying` (from a dead process) is force-set to `failed` with `retryable=false` and message "may or may not have taken effect" — because the write *might have landed*.
3. **Outcome classification:** `callMayHaveTakenEffect(err)` — only a 401/403 *proves* the call was refused pre-dispatch. Dropped connections, malformed replies, oversized bodies = outcome unknown = not retryable.
4. **Result written before attached:** the applied state is persisted in its own small write *before* the (potentially untrusted/huge) result payload is stored, so handling a bad payload can't lose the fact the write happened.
5. **Bounds:** MAX_PENDING=50, MAX_RETAINED=100, result capped at 128KB, args at 64KB.

### Go adaptation
```go
// pkg/approval/store.go — SQLite (or Postgres) backed
type ActionState string
const (
    StatePending ActionState = "pending"
    StateApplying ActionState = "applying"  // claimed, in-flight
    StateApplied ActionState = "applied"
    StateRejected ActionState = "rejected"
    StateFailed ActionState = "failed"
)

type StoredAction struct {
    ID          int64
    ToolName    string
    Args        json.RawMessage
    State       ActionState
    SubmittedAt time.Time
    ClaimedAt   *time.Time
    Retryable   *bool            // nil=retryable, false=outcome-unknown
    Result      *json.RawMessage
    Error       string
}

// pkg/approval/queue.go
type ApprovalQueue interface {
    AuthorizeObservation(desc Observation) error   // synchronous, before data returns
    SubmitAction(id int64, desc Action) error       // async, returns immediately
}

// The call path in every tool session (mirrors session.ts callTool):
func (s *Session) CallTool(name string, args map[string]any) (Result, error) {
    classified := classify(name) // read vs action (§4)
    if classified.IsRead {
        res, _ := s.host.Call(...)               // execute read
        s.queue.AuthorizeObservation(describeCall(...)) // gate BEFORE returning data
        return res, nil
    }
    // action: stage, submit, return pending
    staged := s.host.StageAction(name, args)
    s.queue.SubmitAction(staged.ID, describeCall(...))
    return Result{Status: "pending", ActionID: staged.ID}, nil
}
```
**For Talon:** `awaitDecision: true` is essential — a pentest agent must *pause* after submitting e.g. "run sqlmap against prod host" rather than reasoning against a world where it didn't happen. The at-most-once claim prevents an LLM retry loop from hitting a target twice.

---

## 3. Blueprints (Playbook Templates)

**Source:** `docs/blueprints.md`, `packages/workshop-backend/format-blueprints/` (`<name>.gadget` + `<name>.json`).

### Core idea
A blueprint = a **versioned, shareable template** capturing *code + binding requirements* but NOT credentials, storage, or history. Two distribution modes:
- **Published blueprints:** random 128-bit hex ID, shareable by link, owner-attributed, versioned (old versions retained to avoid race during concurrent instantiation).
- **Bundled/format blueprints:** stable readable IDs (`format.document`), shipped as committed data, no owning user — installed on first deploy, promoted once ever.

### Structure
```
format-blueprints/
  workspace-docs.gadget    # binary archive: magic(8B) + version(4B) + JSON-metadata-len(4B) +
                           #   content-len(8B) + BlueprintMetadata + raw content bytes
  workspace-docs.json      # human-curated sidecar (the .gadget's own title is inert):
                           #   { blueprintId, title, description, output:{id,noun,plural,icon}, author, revision }
```

`.gadget` binary container:
```
0xec2e2d3a2300e317 (8-byte magic) | version=1 (4B) | JSON-len (4B) | content-len (8B) | JSON metadata | content
```
Metadata capped at 64KB, content at 32MB.

A blueprint declares **binding requirements** (3 types), each optionally annotated with friendly name/description/suggested-value:
- `gatekeeper` — external resource (records vendor + urlPattern; instantiator picks their own account)
- `aiModel` — LLM binding (instantiator picks from their configured models)
- `agentSpawner` — sub-agent binding (carries spawner config)

### Storage: 3-tier, one-way propagation
`Gadget DO (authoritative, +dirty flag) → User DO (denormalized for listing) → Workers KV (public lookup)`. Code content in R2 keyed `<blueprintId>/<version>`. The `dirty` flag enables a "Retry publish" UX.

### Go adaptation for Talon → "Pentest Playbooks"
Map directly to reusable pentest templates (e.g. "External recon → nmap → service enum → web exploit chain"):

```go
// pkg/playbook/blueprint.go
type Blueprint struct {
    ID          string            // 128-bit hex, or "talon.recon" for bundled
    Title       string
    Description string
    Version     int
    Revision    int               // reinstall trigger for bundled (archive byte changes)
    Output      OutputFormat      // grouping id + noun/plural/icon for the results page
    Author      Author
    Bindings    []BindingRequirement  // each: gatekeeper/aiModel/agentSpawner + optional annotation
    Code        []byte            // the playbook DAG/workflow definition (YAML or serialized graph)
}

type BindingRequirement struct {
    Name        string            // stable key the playbook code references
    Type        BindingType       // "tool" | "model" | "subagent"
    VendorID    string            // e.g. "nmap", "openai"
    URLPattern  string            // target scope, e.g. "https://*.target.com/*"
    Suggest     string            // optional suggested value
    Description string            // helper text shown during instantiation
}
```
- Ship default playbooks as committed data (`talon.recon`, `talon.webapp`) using the `.json` sidecar pattern so title/description are reviewable in PRs without touching binary.
- Version retention: keep old playbook versions so in-flight engagements don't break on an upgrade.
- `output.id` = the grouping key on a findings/results page (e.g. `report` groups all report-producing playbooks).

---

## 4. MCP Integration

**Source:** `mcp-shared/README.md`, `mcp-shared/src/tools.ts`, `scope.ts`, `session.ts`, `facet.ts`, `action-store.ts`, `gatekeeper-mcp/src/mcp.ts`.

This is the most directly portable subsystem for Talon — it shows how to connect *any* MCP server as a capability with read/write classification and approval gating.

### Trust tiers (`tools.ts`) — the single decision point
```ts
type ServerTrust = "vetted" | "byo";
// classifyTool is the ONLY place a tool's annotations become a policy decision:
function classifyTool(tool, trust): ClassifiedTool {
  const readOnly = tool.annotations?.readOnlyHint === true;   // strict === true
  const autoApprovable = !readOnly
    && trust === "vetted"                                      // deployment, not server, decides
    && annotations.destructiveHint === false
    && annotations.idempotentHint === true;
  return { mode: readOnly ? "read" : "action", autoApprovable, classifiedBy: readOnly ? "server-annotation" : "default" };
}
```
- `byo` (user-typed URL): `readOnlyHint` classifies reads, but **nothing can auto-apply a write**.
- `vetted` (admin-asserted endpoint): may auto-approve non-destructive idempotent writes.
- Unannotated tools default to action (needs approval), never auto-apply — fail-safe.
- **Neither tier is shareable** (owner-only bindings).

### Scope grammar (`scope.ts`) — capability narrowing via URL fragment
```
<endpoint>                       → every tool, now and later
<endpoint>#server=github         → one portal upstream server
<endpoint>#tool=a&tool=b         → only these exact tools
```
`scopeAllows(scope, toolName, isPortal)` is the per-call gate. Namespaced action tags (`mcp:<encoded-endpoint>:<tool>`) prevent an always-approve decision on one server leaking to another.

### The session call path (`session.ts`) — one path every call takes
`callTool()` → classify → if read: execute + `authorizeObservation` before returning; if action: `stageAction` (SQLite insert pending) + `submitAction` + return `{status:"pending", actionId}`. Then `getActionResult(id)` polled until applied/rejected/failed. The applied-result handoff is *also* an observation (authorized before handed to agent).

### Approval prompt hardening (`describeCall`)
Server-controlled text (tool descriptions, args) rendered into the approval prompt is **defused**: Markdown fences neutralized, headings stripped, length-capped (desc 600 chars, args 4000 chars), backticks removed from code spans. Prevents a malicious MCP server from injecting its own prose into the prompt that asks a human to approve.

### Go adaptation for Talon
Talon can be both an **MCP client** (connect external MCP tool servers as gatekeepers) and expose its own tools as MCP. The `tools.ts`/`scope.ts`/`session.ts` trio ports almost line-for-line:

```go
// pkg/mcp/classify.go
type ServerTrust string
const ( TrustVetted ServerTrust = "vetted"; TrustBYO = "byo" )

type ClassifiedTool struct {
    Tool          McpTool
    Mode          string   // "read" | "action"
    AutoApprovable bool
    ClassifiedBy  string   // "server-annotation" | "default"
}
func ClassifyTool(tool McpTool, trust ServerTrust) ClassifiedTool {
    readOnly := ptrBool(tool.Annotations.ReadOnlyHint) == true
    auto := !readOnly && trust == TrustVetted &&
        ptrBool(tool.Annotations.DestructiveHint) == false &&
        ptrBool(tool.Annotations.IdempotentHint) == true
    return ClassifiedTool{Mode: pick(readOnly, "read", "action"), AutoApprovable: auto, ...}
}

// pkg/mcp/scope.go — URL fragment grammar + scopeAllows()
// pkg/mcp/store.go  — the §2 at-most-once action store
```
Pentest relevance: connect MCP tool servers (e.g. a vuln-scanner MCP) as gatekeepers; reads (recon) run free, writes (exploit) queue for approval. A `vetted` internal tool server can auto-approve idempotent "mark finding" writes.

---

## 5. Sharing & Collaboration

**Source:** `docs/sharing.md`.

### Core idea
Capability-based access at `open()`: compute the caller's **effective role** from a permission graph, hand back a *different object* by role. Two roles, totally ordered `build > use`.

### Roles
- `build` — full edit/chat/bindings access (like owner, minus: can't delete, BYOK models, own connected accounts, limited revocation).
- `use` — render + interact with deployed UI only. `UseOverseerInterface` implements the full interface but throws `Unauthorized` everywhere except an explicit allowlist. Because it `implements Overseer`, **newly-added methods fail to compile until a dev decides** whether `use` may call them — **default-deny by construction**.

### Permission graph + lazy revocation (the clever part)
Edges: `user` (sharer→target) and `share-link` (link→target). Effective role = max over valid edges of `min(edgeRole, sharerEffectiveRole)`, iterated to fixed point. The owner is implicit `build` root.
- **Lazy revocation:** severing an edge doesn't cascade-delete; dependents just become *unreachable* and are denied at next `open()`. Roles only increase during iteration → terminates.
- **Non-destructive/undoable:** re-adding an intermediary restores all downstream grants because records were never pruned.
- Live recomputation at every `open()` is the **sole source of truth** — no eager cleanup whose bug could grant access.

### Share links
Random 128-bit key, stored as **HMAC-SHA256 hash** (server can't reconstruct; DB leak exposes nothing). Redeemed atomically with `openGadget(id, shareKey)` in one RPC.

### Resource isolation
Collaborators share code/storage/chat but AI models, gatekeeper bindings resolve from *each user's own accounts* — no implicit access to another's credentials.

### Live-session eviction
Since auth is only checked at `open()`, revocation that affects someone calls `ctx.abort()` to restart the DO, forcing reconnects (with a ~100ms delay + storage sync so the revoker's own response lands first).

### Go adaptation for Talon → "Engagement sharing / team access"
```go
// pkg/sharing/role.go
type Role string
const ( RoleBuild Role = "build"; RoleUse Role = "use" )

// pkg/sharing/graph.go
type PermissionEdge struct {
    Target  string  // user ID
    Sharer  string  // user ID or "" for owner
    Kind    string  // "user" | "link"
    LinkID  string
    Role    Role
}
// EffectiveRoles: fixed-point propagation, roles monotonically increase → terminates.
// Sole source of truth at OpenEngagement(); no denormalized access state.

// pkg/sharing/capability.go — return DIFFERENT interface by role
type BuildCapability interface { Edit(); Chat(); AddBinding(); /* full set */ }
type UseCapability interface { Render(); GetMetadata() }       // narrow allowlist
// Implement UseCapability to also satisfy a "full" interface so adding a method is a compile error
// until a dev decides whether use-role may call it (default-deny).
```
For Talon: share a pentest engagement with read-only "view findings" (use) vs. "run scans" (build). Share-link = HMAC-hashed. Lazy revocation keeps the audit graph intact. Pentest-specific: a `build` operator's tool bindings use *their own* API keys/tokens, never the engagement owner's.

---

## 6. Action Logging / Observers (Information-Flow Control)

**Source:** `docs/observers.md`.

### Core idea
When you share a gadget that has *read sensitive data*, the system enforces that **collaborators can't see data they themselves lack permission to read directly**. This is the security invariant most relevant to a pentest platform (prevent leaking recon findings to someone unauthorized for that target).

### Mechanism
- **Observer:** a non-owner collaborator who can see data the gadget read. On first open (and every open), `Gatekeeper.addObserver(observerID, verifier)` runs *inside the gatekeeper's trust domain* to verify this user can directly read everything observed so far — throws to deny.
- **Verifier:** minted by the *observer's own* connected account (`GatekeeperUser.getVerifier()`), passed back to the same vendor gatekeeper, which "unwraps" it to learn the observer's vendor identity and check ACLs.
- **Forward exclusion:** for observations made *after* a user became an observer, the gatekeeper sets `ObservationDescription.excludeObservers: string[]`. The overseer must block the observation (or the named observer must have already lost access) — otherwise the read throws.
- Observer IDs are **opaque random strings** (not emails) to avoid tempting gatekeeper authors to parse identity.

### Four per-resource strategies (chosen per resource type)
- **A — Private-only:** `addObserver` always throws (replaces `prohibitAllSharing`).
- **B — ACL check (single unit):** verify against one atomic resource's ACL. Cache per-open.
- **C — Data-set tracking:** gatekeeper logs data sets observed; re-verifies all observers when a new set is touched; sets `excludeObservers` for failures.
- **D — Low-stakes:** no-op (any collaborator may observe).

### Authorization keys off the sharing graph (not live sessions)
Because observed data may be *stored and re-displayed later*, every exclusion decision keys off `computeEffectiveRoles` (still in the sharing graph?), never "is this session currently open."

### Go adaptation for Talon → "Need-to-know scoping for shared engagements"
Most directly: when sharing an engagement, verify each viewer is authorized for every finding/target already collected. Strategy **C** (data-set tracking) fits pentest: track which in-scope hosts were observed; when a new target is added, re-verify all observers; block the observation if a current viewer lacks authorization for that target.

```go
// pkg/capability/observer.go
type Verifier interface { /* opaque identity for ACL checks */ }

type ObservableGatekeeper interface {
    AddObserver(observerID string, v Verifier) error  // throws → deny open; verify can read all prior observations
    RemoveObserver(observerID string)                  // idempotent
}

type Observation struct {
    Title, Description string
    ProhibitAllSharing *bool   // lock engagement after this
    ExcludeObservers   []string // block if named viewers still authorized
}

// In the approval queue:
func (q *Queue) AuthorizeObservation(obs Observation) error {
    for _, oid := range obs.ExcludeObservers {
        if sharingGraph.IsAuthorized(oid) {
            return ErrBlockedByObserverACL // degrade to per-observation lockdown
        }
    }
    if obs.ProhibitAllSharing != nil && *obs.ProhibitAllSharing {
        q.lockdown() // engagement can only observe, no actions
    }
    return nil
}
```

---

## 7. The workshop-shared RPC API

**Source:** `packages/workshop-shared/src/api.ts` (2904 lines — the full client/server contract).

### Structure
The entire client↔server boundary is an **RPC API over a persistent WebSocket** (Cap'n Web protocol). Three tiers:
- **`PublicApi`** — unauthenticated: `getServerConfig()`, `getBlueprint(id)` (no auth — "a blueprint is just data"), `downloadBlueprint()`, `login()`/`createAccount()` (Argon2id client-side hashing with a service salt + username), `authenticate()`.
- **`AuthenticatedApi`** — per-user: model management, connected accounts (`subscribeConnectedAccounts` — push-based), `openGadget(id, shareKey?, configureObservers?)` (returns an `Overseer` stub), `newGadget()`, `listGatekeeperVendors()`, `connectAccount(vendorId)`, blueprint library, sharing RPCs.
- **`Overseer`** — per-workspace: the live editor/agent surface.

### Notable design choices
- **Promise pipelining:** `openGadget` returns an `Overseer` stub you can immediately pipeline calls on; share-key redemption + open happen in one round trip.
- **Push subscriptions** via `RpcStub<Subscriber>` callbacks (e.g. `ConnectedAccountsSubscriber.add/remove/ready`) rather than polling — ideal for OAuth-completes-in-another-tab flows.
- **Capability callback for observer config:** `openGadget(..., configureObservers?: RpcStub<ObserverConfigCallback>)` — invoked *only* when a non-owner needs to pick accounts; common case (owner) is zero extra round trips.
- **Binding name validation:** `validateBindingName()` rejects reserved words, `__proto__`, `constructor` etc. — applied at *every* chokepoint that writes a binding name (maps keyed by name).
- **Stable error codes:** `OPEN_GADGET_ERROR_CODES` with `createOpenGadgetError`/`getOpenGadgetErrorCode` for machine-readable expected failures (no structured/typed errors for control flow).
- Argon2id params: `parallelism:1, iterations:3, memorySize:64MiB, hashLength:32`, salt = `SERVICE_SALT + encode(username)`.

### Go adaptation for Talon
Talon's Go backend ↔ Next.js frontend should adopt a similar clean tiered RPC boundary. Two options:
1. **gRPC / Connect-RPC** over WebSocket or HTTP/2 — closest to the typed-interface model, gives streaming for subscriptions.
2. **tRPC-style over WebSocket** (as Talon's Next.js already implies TS) with a Go server emitting JSON-RPC.

Key API shapes to mirror:
```go
// Public (no auth): playbook fetch, login (Argon2id client-side), server config
// Authenticated: open engagement (returns a live workspace handle), connected tools,
//   push-subscribe to tool/account changes via a streaming/callback channel
type AuthenticatedAPI interface {
    OpenEngagement(id string, shareKey string, cfg ObserverConfig) (Overseer, error)
    ConnectTool(vendorID string) (authURL string, err error)
    SubscribeTools(sub Subscriber) (cancel func())   // push-based, not poll
}
// Overseer: live agent chat, code/scan edits, action approval feed
```
Push subscriptions (not polling) for "OAuth completed in another tab" and live action-approval feeds. Binding/target-name validation at every write site. Stable string error codes for expected failures.

---

## Cross-Cutting Relevance Ranking for Talon

| Pattern | Pentest relevance | Effort | Notes |
|---|---|---|---|
| **§2 Approval queue (at-most-once)** | Critical | Medium | Every dangerous action (exploit, brute force) gated; claim-before-dispatch prevents LLM-retry double-fire. |
| **§1 Gatekeepers / capabilities** | Critical | Medium | Each tool = narrow Session capability; creds never in agent reach. |
| **§4 MCP read/write classify** | High | Low-Med | Port `classifyTool`+`scope` directly; connect external MCP tools; expose Talon tools as MCP. |
| **§3 Blueprints** | High | Low | Reusable pentest playbook templates; ship defaults as committed data. |
| **§6 Observers / IFC** | High | High | Need-to-know scoping when sharing engagements; strategy C (data-set tracking) for per-target ACL. |
| **§5 Sharing/roles** | Medium-High | Medium | Team access to engagements; lazy-revocation permission graph; HMAC share links. |
| **§7 Tiered RPC API** | Medium | Low | Clean API boundary; push subscriptions; mirror AuthenticatedApi→Overseer split. |

**Strongest single takeaway:** the **§2 action store state machine** (claim-before-dispatch, crash-recovery to "outcome unknown / not retryable", result-written-before-attached) is the precise safety primitive a pentest agent needs when an LLM action touches a real production target — adopt it verbatim in Go.
