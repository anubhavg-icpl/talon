/**
 * Typed client for the talon-core API.
 *
 * The browser NEVER talks to talon-core directly — everything goes through
 * the Next.js proxy at /api/talon/* (see src/app/api/talon/[...path]/route.ts).
 */

export type RunStatus = 'initializing' | 'running' | 'awaiting_approval' | 'completed' | 'error' | 'not_found'

export type RunSummary = {
  run_id: string
  target: string
  cve_id?: string
  service_name?: string
  status: RunStatus
  judge_verdict?: boolean
  tool_calls: number
  findings_count?: number
  agent_mode?: string
  started_at: string // RFC3339
  ended_at?: string // RFC3339, present once the run reached a terminal state
}

export type ListRunsResponse = {
  runs: RunSummary[] | null
  total: number
  limit: number
  offset: number
}

export type StartRunRequest = {
  ip: string
  cve_id?: string
  service_name?: string
  description?: string
  lhost?: string
  lport?: number
  agent_mode?: string
  playbook_id?: string
}

export type StartRunResponse = {
  run_id: string
  message: string
}

/** Pending HITL interrupt — JSON keys are capitalized Go field names. */
export type Interrupt = {
  ToolName: string
  Args: Record<string, unknown>
}

export type FindingsSummary = {
  total: number
  critical: number
  high: number
  medium: number
  low: number
  info: number
  confirmed: number
}

export type StatusResponse = {
  status: RunStatus
  output: string
  interrupt: Interrupt | null
  judge_verdict?: boolean
  findings_summary?: FindingsSummary
  findings_count?: number
  has_report?: boolean
  agent_mode?: string
  methodology_percent?: number
}

export type GateEvidence = {
  baseline?: string
  attack?: string
  diff?: string
  passed: boolean
}

export type Finding = {
  id: string
  severity: 'critical' | 'high' | 'medium' | 'low' | 'info' | string
  title: string
  description: string
  cwe_id?: string
  endpoint?: string
  attack_vector?: string
  steps_to_reproduce?: string
  business_impact?: string
  recommendation?: string
  poc?: string
  evidence: GateEvidence
  status: string
  source: string
  stage?: string
  created_at?: string
}

export type FindingsResponse = {
  run_id: string
  findings: Finding[]
  summary: FindingsSummary
}

export type ResumeDecision = 'approve' | 'reject' | 'edit'

export type ResumeRequest = {
  decision: ResumeDecision
  edited_args?: Record<string, unknown>
}

/** Tool call record — JSON keys are capitalized Go field names. */
export type ToolCallRecord = {
  Index: number
  ToolName: string
  Args: Record<string, unknown> | string
  Output: string
}

export type ToolsResponse = {
  tool_log: ToolCallRecord[] | null
}

export type TracesResponse = {
  history: string[] | null
}

export type HealthResponse = {
  status: string
  service: string
}

const BASE = '/api/talon'

async function request<T>(path: string, init?: RequestInit): Promise<T> {
  const res = await fetch(`${BASE}${path}`, {
    ...init,
    headers: init?.body ? { 'content-type': 'application/json', ...init?.headers } : init?.headers,
    cache: 'no-store'
  })

  if (!res.ok) {
    // Session expired/missing — back to the gate (client-side only).
    if (res.status === 401 && typeof window !== 'undefined' && !window.location.pathname.startsWith('/login')) {
      window.location.href = '/login'
    }

    const text = await res.text().catch(() => '')

    throw new Error(`talon-core ${res.status}: ${text || res.statusText}`)
  }

  return (await res.json()) as T
}

export const health = () => request<HealthResponse>('/health')

export type ServiceHealth = {
  name: string
  endpoint: string
  status: 'online' | 'offline' | 'unconfigured'
  detail: string
  latency_ms: number
}

/** Live per-service probes run server-side by talon-core (GET /health/services). */
export const serviceHealth = () => request<{ services: ServiceHealth[] }>('/health/services')

export type MeResponse = {
  username: string
  auth?: string
}

export const me = () => request<MeResponse>('/auth/me')

export const logout = () => request<unknown>('/auth/logout', { method: 'POST' })

export const listRuns = (limit?: number, offset?: number) => {
  const params = new URLSearchParams()

  if (limit !== undefined) params.set('limit', String(limit))
  if (offset !== undefined) params.set('offset', String(offset))

  const qs = params.toString()

  return request<ListRunsResponse>(`/runs${qs ? `?${qs}` : ''}`)
}

/** Aggregate run counters computed server-side (GET /runs/summary). */
export type RunsSummaryResponse = {
  total: number
  active: number
  compromised: number
  awaiting_approval: number
  completed: number
  errored: number
}

export const runsSummary = () => request<RunsSummaryResponse>('/runs/summary')

/** Live config entry exposed by talon-core (GET /config). Secret values arrive masked. */
export type ConfigEntry = {
  key: string
  label: string
  value: string
  set: boolean
  secret: boolean
  hot: boolean
  source: 'database' | 'env' | 'default'
  writable: boolean
}

export const getConfig = () => request<{ config: ConfigEntry[] }>('/config')

/** Empty string or the masked value on a secret key means "leave unchanged". */
export const putConfig = (values: Record<string, unknown>) =>
  request<{ updated: string[]; note: string }>('/config', { method: 'PUT', body: JSON.stringify(values) })

export type MCPServerInfo = {
  name: string
  tools: string[]
}

export type MCPServersResponse = {
  servers: MCPServerInfo[]
  skill_stats?: Record<string, number>
  agent_to_agent?: {
    model?: string
    notes?: string[]
  }
}

export const getMCPServers = () => request<MCPServersResponse>('/mcp/servers')

export const startRun = (body: StartRunRequest) =>
  request<StartRunResponse>('/input/start', { method: 'POST', body: JSON.stringify(body) })

export const getStatus = (runId: string) => request<StatusResponse>(`/output/status/${runId}`)

export const resumeRun = (runId: string, body: ResumeRequest) =>
  request<unknown>(`/output/resume/${runId}`, { method: 'POST', body: JSON.stringify(body) })

export const getTools = (runId: string) => request<ToolsResponse>(`/monitor/tools?run_id=${encodeURIComponent(runId)}`)

export const getTraces = (runId: string) => request<TracesResponse>(`/monitor/traces/${runId}`)

export type StructuredReport = {
  markdown: string
  generated_at?: string
  sections?: string[]
  findings?: Finding[]
  summary?: FindingsSummary
  judge_verdict?: boolean
  target?: string
  cve_id?: string
  stages_covered?: string[]
  message?: string
}

export const getFindings = (runId: string) => request<FindingsResponse>(`/runs/${runId}/findings`)

export const getReport = (runId: string) => request<StructuredReport>(`/runs/${runId}/report`)

export type Skill = {
  id: string
  name: string
  stage: string
  category?: string
  body?: string
  source?: string
  path?: string
}

export type CategoryCount = {
  name: string
  count: number
}

export type SkillsListResponse = {
  skills: Skill[]
  count: number
  total: number
  offset: number
  limit: number
  stats?: Record<string, number>
  categories?: CategoryCount[]
}

export type SkillsQuery = {
  brief?: boolean
  full?: boolean
  stage?: string
  category?: string
  q?: string
  limit?: number
  offset?: number
}

export const getSkills = (opts?: SkillsQuery) => {
  const params = new URLSearchParams()
  if (opts?.brief) params.set('brief', '1')
  if (opts?.full) params.set('full', '1')
  if (opts?.stage) params.set('stage', opts.stage)
  if (opts?.category) params.set('category', opts.category)
  if (opts?.q) params.set('q', opts.q)
  if (opts?.limit !== undefined) params.set('limit', String(opts.limit))
  if (opts?.offset !== undefined) params.set('offset', String(opts.offset))
  const qs = params.toString()
  return request<SkillsListResponse>(`/skills${qs ? `?${qs}` : ''}`)
}

export const getSkill = (id: string) => request<Skill>(`/skills/${encodeURIComponent(id)}`)

export type AgentInfo = {
  id: string
  name: string
  codename: string
  focus: string
  description: string
  delegates: string[]
}

export const getAgents = () => request<{ agents: AgentInfo[]; count: number }>('/agents')

export type Playbook = {
  id: string
  name: string
  codename: string
  description: string
  agent_mode: string
  prompt: string
  tags: string[]
}

export const getPlaybooks = () => request<{ playbooks: Playbook[]; count: number }>('/playbooks')

export type IntelEvent = {
  at: string
  run_id: string
  target: string
  kind: string
  label: string
  detail?: string
  severity?: string
}

export const getIntel = (limit?: number) =>
  request<{ events: IntelEvent[] }>(`/intel${limit ? `?limit=${limit}` : ''}`)

export type TimelineEvent = {
  index: number
  kind: string
  label: string
  stage?: string
  detail?: string
  severity?: string
}

export const getTimeline = (runId: string) =>
  request<{ run_id: string; timeline: TimelineEvent[] }>(`/runs/${runId}/timeline`)

export type OperatorNote = {
  id: string
  author?: string
  body: string
  created_at: string
}

export const getNotes = (runId: string) =>
  request<{ run_id: string; notes: OperatorNote[] }>(`/runs/${runId}/notes`)

export const addNote = (runId: string, body: string, author?: string) =>
  request<OperatorNote>(`/runs/${runId}/notes`, {
    method: 'POST',
    body: JSON.stringify({ body, author: author || 'operator' })
  })

export const compareRuns = (a: string, b: string) =>
  request<Record<string, unknown>>(`/runs/compare?a=${encodeURIComponent(a)}&b=${encodeURIComponent(b)}`)

export const exportRun = (runId: string) => request<Record<string, unknown>>(`/runs/${runId}/export`)

export const batchStart = (body: {
  ips: string[]
  cve_id?: string
  service_name?: string
  description?: string
  lhost?: string
  lport?: number
  agent_mode?: string
  playbook_id?: string
}) => request<{ started: { run_id: string; ip: string }[]; count: number }>('/input/batch', {
  method: 'POST',
  body: JSON.stringify(body)
})

export type ScopePolicy = {
  enabled: boolean
  allowed_cidrs: string[]
  denied_cidrs: string[]
  denied_ports?: number[]
  max_concurrent: number
  require_auth_label: boolean
  auto_approve_nmap_private: boolean
  updated_at?: string
}

export const getScope = () => request<ScopePolicy>('/scope')
export const putScope = (body: ScopePolicy) =>
  request<ScopePolicy>('/scope', { method: 'PUT', body: JSON.stringify(body) })

export type Target = {
  id: string
  address: string
  url?: string
  label?: string
  tags?: string[]
  notes?: string
  last_run_id?: string
  last_status?: string
  created_at?: string
  updated_at?: string
}

export const listTargets = () => request<{ targets: Target[] }>('/targets')
export const upsertTarget = (t: Partial<Target> & { address?: string; url?: string }) =>
  request<Target>('/targets', { method: 'POST', body: JSON.stringify(t) })
export const deleteTarget = (id: string) =>
  request<{ deleted: string }>(`/targets/${id}`, { method: 'DELETE' })

export type Schedule = {
  id: string
  name: string
  interval: string
  target: string
  playbook_id?: string
  agent_mode?: string
  enabled: boolean
  last_run_at?: string
  next_run_at?: string
  created_at?: string
}

export const listSchedules = () => request<{ schedules: Schedule[] }>('/schedules')
export const upsertSchedule = (s: Partial<Schedule> & { name: string; target: string }) =>
  request<Schedule>('/schedules', { method: 'POST', body: JSON.stringify(s) })
export const deleteSchedule = (id: string) =>
  request<{ deleted: string }>(`/schedules/${id}`, { method: 'DELETE' })

export type NotifyConfig = {
  webhook_url: string
  on_complete: boolean
  on_hitl: boolean
  on_critical_finding: boolean
  on_error: boolean
}

export const getNotify = () => request<NotifyConfig>('/notify')
export const putNotify = (n: NotifyConfig) =>
  request<NotifyConfig>('/notify', { method: 'PUT', body: JSON.stringify(n) })

export type Credential = {
  id: string
  name: string
  kind: string
  username?: string
  has_secret: boolean
  scope?: string
  created_at?: string
}

export const listCredentials = () => request<{ credentials: Credential[] }>('/credentials')
export const addCredential = (body: {
  name: string
  kind?: string
  username?: string
  secret: string
  scope?: string
}) => request<Credential>('/credentials', { method: 'POST', body: JSON.stringify(body) })
export const deleteCredential = (id: string) =>
  request<{ deleted: string }>(`/credentials/${id}`, { method: 'DELETE' })

export type EvidenceItem = {
  id: string
  run_id: string
  finding_id?: string
  kind: string
  title: string
  body: string
  created_at: string
}

export const listEvidence = (runId?: string) =>
  request<{ evidence: EvidenceItem[] }>(`/evidence${runId ? `?run_id=${encodeURIComponent(runId)}` : ''}`)
export const addEvidence = (e: Partial<EvidenceItem> & { run_id: string; title: string }) =>
  request<EvidenceItem>('/evidence', { method: 'POST', body: JSON.stringify(e) })

export type BudgetStats = {
  llm_calls: number
  tool_calls: number
  runs_started: number
  runs_completed: number
  critical_findings: number
}

export const getBudget = () => request<BudgetStats>('/budget')

export const retestRun = (runId: string, findingId?: string) =>
  request<{ run_id: string; message: string }>(`/runs/${runId}/retest`, {
    method: 'POST',
    body: JSON.stringify({ finding_id: findingId || '' })
  })

export const reportHTMLUrl = (runId: string) => `/api/talon/runs/${runId}/report.html`

export type KillChainLink = {
  from: string
  to: string
  severity: string
  reason: string
}

export type KillChainAnalysis = {
  chains: KillChainLink[]
  next_steps: string[]
  summary: string
  max_severity: string
}

export const getKillChain = (runId: string) => request<KillChainAnalysis>(`/runs/${runId}/killchain`)

export type CoverageItem = {
  stage: string
  label: string
  covered: boolean
  tools?: string[]
  notes?: string
}

export type MethodologyState = {
  items: CoverageItem[]
  covered_count: number
  total_count: number
  percent: number
  agent_mode?: string
}

export const getMethodology = (runId: string) => request<MethodologyState>(`/runs/${runId}/methodology`)

export const triageFinding = (runId: string, findingId: string, status: string, duplicateOf?: string) =>
  request<Finding>(`/runs/${runId}/findings/${findingId}/triage`, {
    method: 'POST',
    body: JSON.stringify({ status, ...(duplicateOf ? { duplicate_of: duplicateOf } : {}) })
  })

export type GlobalFinding = {
  run_id: string
  target: string
  finding: Finding
}

export const getGlobalFindings = (opts?: { severity?: string; limit?: number }) => {
  const params = new URLSearchParams()
  if (opts?.severity) params.set('severity', opts.severity)
  if (opts?.limit !== undefined) params.set('limit', String(opts.limit))
  const qs = params.toString()
  return request<{ findings: GlobalFinding[]; count: number }>(`/findings${qs ? `?${qs}` : ''}`)
}

export type AnalyzeResponse = {
  run_id: string
  analysis: string
}

/** One-shot LLM analysis of a run (report + tool log) — server-side analyst model. */
export const analyzeRun = (runId: string) =>
  request<AnalyzeResponse>(`/analyze/${runId}`, { method: 'POST' })

export type FindingsStreamPayload = {
  findings_count: number
  findings_summary?: FindingsSummary
  findings?: Finding[]
}

export type StreamRunHandlers = {
  onTool?: (tool: ToolCallRecord) => void
  onStatus?: (status: StatusResponse) => void
  onFindings?: (payload: FindingsStreamPayload) => void
  onError?: (error: Error) => void
}

const TERMINAL_STATUSES: RunStatus[] = ['completed', 'error', 'not_found']

/**
 * Live-follow a run with a resilience chain:
 *   WebSocket (direct to talon-core /monitor/ws) → SSE (proxied /monitor/stream) → 3s polling.
 */
export const streamRun = (runId: string, handlers: StreamRunHandlers): (() => void) => {
  let stopped = false
  let pollTimer: ReturnType<typeof setTimeout> | null = null
  let es: EventSource | null = null
  let ws: WebSocket | null = null
  let lastToolIndex = -1

  const stopPolling = () => {
    if (pollTimer !== null) {
      clearTimeout(pollTimer)
      pollTimer = null
    }
  }

  const closeSockets = () => {
    if (ws) {
      ws.onclose = null
      ws.onerror = null
      ws.close()
      ws = null
    }

    if (es) {
      es.close()
      es = null
    }
  }

  const handleTool = (tool: ToolCallRecord) => {
    if (tool.Index > lastToolIndex) lastToolIndex = tool.Index
    handlers.onTool?.(tool)
  }

  /** Returns true when the status is terminal (stream is done). */
  const handleStatus = (status: StatusResponse): boolean => {
    handlers.onStatus?.(status)

    if (TERMINAL_STATUSES.includes(status.status)) {
      stopped = true
      closeSockets()
      stopPolling()

      return true
    }

    return false
  }

  const startPolling = () => {
    if (stopped || pollTimer !== null) return

    const tick = async () => {
      if (stopped) return

      try {
        const [status, tools, findingsRes] = await Promise.all([
          getStatus(runId),
          getTools(runId),
          getFindings(runId).catch(() => null)
        ])

        if (stopped) return

        for (const tool of tools.tool_log ?? []) {
          if (tool.Index > lastToolIndex) handleTool(tool)
        }

        if (findingsRes) {
          handlers.onFindings?.({
            findings_count: findingsRes.findings?.length ?? 0,
            findings_summary: findingsRes.summary,
            findings: findingsRes.findings ?? []
          })
        }

        if (handleStatus(status)) return
      } catch (err) {
        handlers.onError?.(err instanceof Error ? err : new Error(String(err)))
      }

      if (!stopped) pollTimer = setTimeout(tick, 3000)
    }

    pollTimer = setTimeout(tick, 0)
  }

  const startSSE = () => {
    if (stopped) return

    try {
      es = new EventSource(`${BASE}/monitor/stream/${runId}`)

      es.addEventListener('tool', ev => {
        try {
          handleTool(JSON.parse((ev as MessageEvent).data) as ToolCallRecord)
        } catch {
          // ignore malformed events
        }
      })

      es.addEventListener('status', ev => {
        try {
          handleStatus(JSON.parse((ev as MessageEvent).data) as StatusResponse)
        } catch {
          // ignore malformed events
        }
      })

      es.addEventListener('findings', ev => {
        try {
          handlers.onFindings?.(JSON.parse((ev as MessageEvent).data) as FindingsStreamPayload)
        } catch {
          // ignore
        }
      })

      es.onerror = () => {
        es?.close()
        es = null

        if (!stopped) startPolling()
      }
    } catch {
      startPolling()
    }
  }

  const startWS = () => {
    if (stopped) return

    try {
      const base = process.env.NEXT_PUBLIC_TALON_WS_URL ?? `ws://${window.location.hostname}:8000`

      ws = new WebSocket(`${base}/monitor/ws/${runId}`)

      ws.onmessage = ev => {
        try {
          const msg = JSON.parse(ev.data as string) as
            | { type: 'tool'; data: ToolCallRecord }
            | { type: 'status'; data: StatusResponse }
            | { type: 'findings'; data: FindingsStreamPayload }

          if (msg.type === 'tool') handleTool(msg.data)
          else if (msg.type === 'status') handleStatus(msg.data)
          else if (msg.type === 'findings') handlers.onFindings?.(msg.data)
        } catch {
          // ignore malformed messages
        }
      }

      ws.onerror = () => {
        // Defer to onclose for the fallback decision
      }

      ws.onclose = () => {
        ws = null
        if (!stopped) startSSE()
      }
    } catch {
      startSSE()
    }
  }

  startWS()

  return () => {
    stopped = true
    closeSockets()
    stopPolling()
  }
}

// ---------------------------------------------------------------------------
// LLM token streaming (POST /llm/stream) — SmolLM/ONNX/OpenAI via Go SSE
// ---------------------------------------------------------------------------

export type LLMStreamMessage = { role: 'system' | 'user' | 'assistant'; content: string }

export type LLMStreamHandlers = {
  onMeta?: (meta: { provider: string; model: string; role?: string }) => void
  onToken?: (token: string) => void
  onDone?: (result: { text: string; ms: number }) => void
  onError?: (error: Error) => void
}

/**
 * Stream chat tokens from talon-core POST /llm/stream (SSE).
 * Prefer LLM_PROVIDER=onnx for millisecond local SmolLM tokens.
 * Returns an abort function.
 */
export const streamLLM = (
  body: {
    messages: LLMStreamMessage[]
    system?: string
    role?: string
    model?: string
    max_tokens?: number
    temperature?: number
  },
  handlers: LLMStreamHandlers
): (() => void) => {
  const ac = new AbortController()

  ;(async () => {
    try {
      const res = await fetch(`${BASE}/llm/stream`, {
        method: 'POST',
        headers: { 'content-type': 'application/json', accept: 'text/event-stream' },
        body: JSON.stringify(body),
        signal: ac.signal,
        cache: 'no-store'
      })

      if (!res.ok) {
        if (res.status === 401 && typeof window !== 'undefined' && !window.location.pathname.startsWith('/login')) {
          window.location.href = '/login'
        }

        const text = await res.text().catch(() => '')

        throw new Error(`llm/stream ${res.status}: ${text || res.statusText}`)
      }

      if (!res.body) throw new Error('llm/stream: empty body')

      const reader = res.body.getReader()
      const decoder = new TextDecoder()
      let buffer = ''
      let event = 'message'

      const dispatch = (ev: string, data: string) => {
        if (!data || data === '[DONE]') return

        try {
          const parsed = JSON.parse(data) as Record<string, unknown>

          if (ev === 'meta') {
            handlers.onMeta?.(parsed as { provider: string; model: string; role?: string })
          } else if (ev === 'token') {
            const c = typeof parsed.content === 'string' ? parsed.content : ''

            if (c) handlers.onToken?.(c)
          } else if (ev === 'done') {
            handlers.onDone?.({
              text: String(parsed.text ?? ''),
              ms: Number(parsed.ms ?? 0)
            })
          } else if (ev === 'error') {
            handlers.onError?.(new Error(String(parsed.error ?? 'stream error')))
          }
        } catch {
          // ignore malformed frames
        }
      }

      while (true) {
        const { done, value } = await reader.read()

        if (done) break
        buffer += decoder.decode(value, { stream: true })
        const parts = buffer.split('\n')

        buffer = parts.pop() ?? ''
        for (const line of parts) {
          if (line.startsWith('event:')) {
            event = line.slice(6).trim()
          } else if (line.startsWith('data:')) {
            dispatch(event, line.slice(5).trim())
            event = 'message'
          } else if (line === '') {
            event = 'message'
          }
        }
      }
    } catch (err) {
      if (ac.signal.aborted) return
      handlers.onError?.(err instanceof Error ? err : new Error(String(err)))
    }
  })()

  return () => ac.abort()
}

export const llmInfo = (role?: string) =>
  request<{
    provider: string
    model: string
    role: string
    onnx_base_url: string
    ollama_url: string
    openai_base: string
    stream_path: string
    assist_path?: string
    tools_path?: string
    tool_count?: number
  }>(`/llm/info${role ? `?role=${encodeURIComponent(role)}` : ''}`)

export type SLMToolDef = {
  name: string
  description: string
  parameters: Record<string, unknown>
  safe: boolean
}

export const listLLMTools = () =>
  request<{
    tools: SLMToolDef[]
    count: number
    protocol: string
    assist_path: string
    stream_path: string
    note: string
  }>('/llm/tools')

export type LLMAssistHandlers = {
  onMeta?: (meta: Record<string, unknown>) => void
  onToken?: (token: string) => void
  onToolStart?: (tool: { name: string; arguments?: Record<string, unknown> }) => void
  onToolResult?: (tool: { name: string; result: string; ms: number }) => void
  onRound?: (round: { n: number; max: number }) => void
  onStatus?: (status: { phase: string; message: string }) => void
  onDone?: (result: { text: string; ms: number; tool_calls?: number; rounds?: number }) => void
  onError?: (error: Error) => void
}

/**
 * End-to-end SLM assist: POST /llm/assist with curated codebase tools.
 * SSE events: meta | status | token | tool_start | tool_result | round | done | error.
 */
export const streamLLMAssist = (
  body: {
    messages: LLMStreamMessage[]
    system?: string
    role?: string
    model?: string
    max_tokens?: number
    max_rounds?: number
    disable_tools?: boolean
  },
  handlers: LLMAssistHandlers
): (() => void) => {
  const ac = new AbortController()

  ;(async () => {
    let finished = false
    const finishDone = (r: { text: string; ms: number; tool_calls?: number; rounds?: number }) => {
      if (finished) return
      finished = true
      handlers.onDone?.(r)
    }
    const finishErr = (err: Error) => {
      if (finished) return
      finished = true
      handlers.onError?.(err)
    }

    try {
      const res = await fetch(`${BASE}/llm/assist`, {
        method: 'POST',
        headers: { 'content-type': 'application/json', accept: 'text/event-stream' },
        body: JSON.stringify(body),
        signal: ac.signal,
        cache: 'no-store'
      })

      if (!res.ok) {
        if (res.status === 401 && typeof window !== 'undefined' && !window.location.pathname.startsWith('/login')) {
          window.location.href = '/login'
        }

        const text = await res.text().catch(() => '')

        throw new Error(`llm/assist ${res.status}: ${text || res.statusText}`)
      }

      if (!res.body) throw new Error('llm/assist: empty body')

      const reader = res.body.getReader()
      const decoder = new TextDecoder()
      let buffer = ''
      let event = 'message'
      let sawDone = false
      let lastText = ''

      const dispatch = (ev: string, data: string) => {
        if (!data || data === '[DONE]') return

        try {
          const parsed = JSON.parse(data) as Record<string, unknown>

          switch (ev) {
            case 'meta':
              handlers.onMeta?.(parsed)
              break
            case 'status':
              handlers.onStatus?.({
                phase: String(parsed.phase ?? ''),
                message: String(parsed.message ?? '')
              })
              break
            case 'token':
              if (typeof parsed.content === 'string' && parsed.content) {
                lastText += parsed.content
                handlers.onToken?.(parsed.content)
              }
              break
            case 'tool_start':
              handlers.onToolStart?.({
                name: String(parsed.name ?? ''),
                arguments: (parsed.arguments as Record<string, unknown>) || undefined
              })
              break
            case 'tool_result':
              handlers.onToolResult?.({
                name: String(parsed.name ?? ''),
                result: String(parsed.result ?? ''),
                ms: Number(parsed.ms ?? 0)
              })
              break
            case 'round':
              handlers.onRound?.({ n: Number(parsed.n ?? 0), max: Number(parsed.max ?? 0) })
              break
            case 'done':
              sawDone = true
              finishDone({
                text: String(parsed.text ?? lastText),
                ms: Number(parsed.ms ?? 0),
                tool_calls: Number(parsed.tool_calls ?? 0),
                rounds: Number(parsed.rounds ?? 0)
              })
              break
            case 'error':
              finishErr(new Error(String(parsed.error ?? 'assist error')))
              break
          }
        } catch {
          // ignore malformed frames / keepalive comments
        }
      }

      const consume = (chunk: string) => {
        buffer += chunk
        const parts = buffer.split('\n')

        buffer = parts.pop() ?? ''
        for (const line of parts) {
          if (line.startsWith('event:')) {
            event = line.slice(6).trim()
          } else if (line.startsWith('data:')) {
            dispatch(event, line.slice(5).trim())
            event = 'message'
          } else if (line === '' || line.startsWith(':')) {
            // blank frame or SSE comment keepalive
            event = 'message'
          }
        }
      }

      while (true) {
        const { done, value } = await reader.read()

        if (done) break
        consume(decoder.decode(value, { stream: true }))
      }
      // Flush decoder + trailing buffer (final data: line without trailing \n).
      consume(decoder.decode())
      if (buffer.trim()) {
        if (buffer.startsWith('data:')) dispatch(event, buffer.slice(5).trim())
      }

      if (!sawDone && !finished) {
        finishDone({ text: lastText, ms: 0 })
      }
    } catch (err) {
      if (ac.signal.aborted) {
        if (!finished) {
          finished = true
        }

        return
      }
      finishErr(err instanceof Error ? err : new Error(String(err)))
    }
  })()

  return () => ac.abort()
}


// ─── Pentest agent: evidence, target state, traffic, recap, crypto ───

export type EvidenceRecord = {
  index: number
  tool: string
  summary: string
  size: number
}

export type RunEvidenceResponse = {
  run_id: string
  total: number
  items: EvidenceRecord[]
}

export type TrafficRecord = {
  seq: number
  tool: string
  output_snippet: string
}

export type RunTrafficResponse = {
  run_id: string
  total: number
  items: TrafficRecord[]
}

export type RecapStep = {
  step: number
  action: string
  result: string
  timestamp: string
}

export type RecapEvidence = {
  id: string
  tool: string
  summary: string
  relevant: boolean
}

export type RecapRepro = {
  label: string
  command: string
}

export type RunRecap = {
  target: string
  run_id: string
  start_time: string
  end_time: string
  duration: string
  solve_path: RecapStep[]
  key_evidence: RecapEvidence[]
  reproduction: RecapRepro[]
  finding_count: number
  verified_count: number
}

export type TargetFinding = {
  id: string
  title: string
  severity: string
  status: string
  evidence_ids?: string[]
  created_at: string
}

export type ReconDimension = {
  name: string
  status: string
  summary: string
  updated_at: string
}

export type TargetState = {
  target: string
  slug: string
  findings: TargetFinding[]
  recon_dimensions: ReconDimension[]
  failed_vectors: Array<{ vector: string; target: string; reason: string; tried_at: string }>
  attack_path: Array<{ step: number; action: string; result: string; timestamp: string }>
  runtime: { os?: string; services?: string; credentials?: string; notes?: string }
  schema_version: number
  updated_at: string
  created_at: string
}

export type PlanStep = {
  priority: number
  category: string
  action: string
  dimension?: string
  confidence: number
}

export type ResumePlanResponse = {
  target: string
  steps: PlanStep[]
  total: number
  summary: string
}

export type CryptoResult = {
  success: boolean
  operation?: string
  result?: string
  error?: string
}

export const getRunEvidence = (runId: string) =>
  request<RunEvidenceResponse>(`/runs/${runId}/evidence`)

export const getRunTraffic = (runId: string) =>
  request<RunTrafficResponse>(`/runs/${runId}/traffic`)

export const getRunRecap = (runId: string, format?: 'json' | 'markdown') =>
  request<RunRecap>(`/runs/${runId}/recap${format === 'markdown' ? '?format=markdown' : ''}`)

export const getRunRecapMarkdown = (runId: string) =>
  request<string>(`/runs/${runId}/recap?format=markdown`)

export const getTargetState = (addr: string) =>
  request<TargetState>(`/targets/${encodeURIComponent(addr)}/state`)

export const getTargetResumePlan = (addr: string) =>
  request<ResumePlanResponse>(`/targets/${encodeURIComponent(addr)}/resume-plan`)

export const getCryptoOperations = () =>
  request<{ operations: string[] }>('/crypto/operations')

export const executeCryptoDecode = (body: {
  operation: string
  input: string
  key?: string
  iv?: string
  shift?: number
}) =>
  request<CryptoResult>('/crypto/decode', { method: 'POST', body: JSON.stringify(body) })

// ─── SOC analysis: triage, investigation, tuning pipeline ───

export type DetectionSkillType = 'triage' | 'investigation' | 'tuning'

export type DetectionSkill = {
  id: string
  name: string
  category: string
  stage: string
  path?: string
}

export type DetectionVerdict =
  | 'escalate' | 'dismiss'
  | 'malicious' | 'suspicious' | 'inconclusive' | 'benign'

export type TriageCheck = {
  name: string
  outcome: 'RISK' | 'CLEAR'
  label?: string
  detail?: string
}

export type InvestigationSignal = {
  name: string
  fired: boolean
  detail?: string
}

export type DetectionCase = {
  id: string
  alert_type: string
  title: string
  entity: string
  entity_type?: string
  severity: string
  source_data?: Record<string, unknown>
  triage_state?: {
    verdict: DetectionVerdict
    risk_count: number
    checks: TriageCheck[]
    evidence?: string
    reasoning?: string
    skill_id?: string
  }
  investigation_state?: {
    verdict: DetectionVerdict
    signals: InvestigationSignal[]
    evidence?: string
    reasoning?: string
    recommended_actions?: string[]
    skill_id?: string
  }
  tuning_state?: {
    action: string
    target?: string
    value?: string
    rationale?: string
    skill_id?: string
  }
}

export const getDetectionSkills = (opts?: { q?: string; type?: DetectionSkillType }) => {
  const params = new URLSearchParams()
  if (opts?.q) params.set('q', opts.q)
  if (opts?.type) params.set('category', opts.type)
  const qs = params.toString()
  return request<{ total: number; skills: DetectionSkill[] }>(`/detection/skills${qs ? '?' + qs : ''}`)
}

export const getDetectionSkillsByType = (type: DetectionSkillType) =>
  request<{ type: string; total: number; skills: DetectionSkill[] }>(`/detection/skills/${type}`)

export const getDetectionCases = () =>
  request<{ cases: DetectionCase[] }>(`/detection/cases`)

// ===========================================================================
// CF-derived features: VFS, Approvals, Gatekeepers, Blueprints, Audit,
// MCP Gateway, and Sharing/Engagements.
// Routes registered in internal/control/cf_handlers.go (RegisterCFRoutes).
// ===========================================================================

// ─── VFS — virtual filesystem (internal/vfs) ────────────────────────────────

/** Entry.type values mirror vfs.NodeType. */
export type VFSNodeType = 'file' | 'dir' | 'symlink'

/** Matches vfs.Entry (Go json tags). */
export type VFSEntry = {
  inode: number
  name: string
  type: VFSNodeType
  mode: number
  size: number
  mtime: string // RFC3339
  rev: number
  target?: string // symlinks
  content?: string // small files on read
}

/** GET /vfs?dir=PATH — list directory entries. */
export const listVFS = (dir?: string) =>
  request<VFSEntry[] | null>(`/vfs${dir ? `?dir=${encodeURIComponent(dir)}` : ''}`)

/** POST /vfs/file {path, content} — write (or overwrite) a file. */
export const writeVFSFile = (path: string, content: string) =>
  request<{ status: string; path: string }>('/vfs/file', {
    method: 'POST',
    body: JSON.stringify({ path, content })
  })

/**
 * GET /vfs/file?path=PATH — read a file's bytes. Unlike the other VFS routes,
 * the backend answers with raw bytes (Content-Type: application/octet-stream),
 * so this bypasses the JSON `request` wrapper and resolves to text.
 */
export const readVFSFile = async (path: string): Promise<string> => {
  const res = await fetch(`${BASE}/vfs/file?path=${encodeURIComponent(path)}`, { cache: 'no-store' })

  if (!res.ok) {
    if (res.status === 401 && typeof window !== 'undefined' && !window.location.pathname.startsWith('/login')) {
      window.location.href = '/login'
    }
    const text = await res.text().catch(() => '')
    throw new Error(`talon-core ${res.status}: ${text || res.statusText}`)
  }

  return res.text()
}

/** DELETE /vfs/file?path=PATH — remove a file or directory recursively. */
export const deleteVFSFile = (path: string) =>
  request<{ status: string }>(`/vfs/file?path=${encodeURIComponent(path)}`, { method: 'DELETE' })

/** POST /vfs/mkdir {path} — create a directory (idempotent, recursive). */
export const mkdirVFS = (path: string) =>
  request<{ status: string }>('/vfs/mkdir', { method: 'POST', body: JSON.stringify({ path }) })

/** GET /vfs/stat?path=PATH — stat a single entry. */
export const statVFS = (path: string) =>
  request<VFSEntry>(`/vfs/stat?path=${encodeURIComponent(path)}`)

// ─── Approvals — human-in-the-loop action gating (internal/approval) ─────────

/** approval.ActionState lifecycle states. */
export type ApprovalActionState = 'pending' | 'applying' | 'applied' | 'rejected' | 'failed' | 'unknown'

/** approval.RiskLevel. */
export type ApprovalRiskLevel = 'low' | 'medium' | 'high' | 'critical'

/** Matches approval.Action (Go json tags). */
export type ApprovalAction = {
  id: string
  run_id: string
  tool_name: string
  args: Record<string, unknown> | string // json.RawMessage
  state: ApprovalActionState
  risk_level: ApprovalRiskLevel
  summary: string
  result?: unknown // json.RawMessage, omitempty
  created_at: string // RFC3339
  resolved_at?: string // *time.Time, omitempty
  claimed_at?: string // *time.Time, omitempty
}

/** POST /approvals — create a new approval-gated action. */
export const createApproval = (action: Partial<ApprovalAction> & { run_id: string; tool_name: string }) =>
  request<ApprovalAction>('/approvals', { method: 'POST', body: JSON.stringify(action) })

/** GET /approvals?run_id=X — list all actions for a run. */
export const listApprovals = (runId: string) =>
  request<ApprovalAction[] | null>(`/approvals?run_id=${encodeURIComponent(runId)}`)

/** GET /approvals/pending?run_id=X — list actions still awaiting resolution. */
export const listPendingApprovals = (runId: string) =>
  request<ApprovalAction[] | null>(`/approvals/pending?run_id=${encodeURIComponent(runId)}`)

/** GET /approvals/{id} — fetch a single action. */
export const getApproval = (id: string) => request<ApprovalAction>(`/approvals/${encodeURIComponent(id)}`)

/** POST /approvals/{id}/approve {result} — approve and record the result. */
export const approveApproval = (id: string, result?: unknown) =>
  request<{ status: string }>(`/approvals/${encodeURIComponent(id)}/approve`, {
    method: 'POST',
    body: JSON.stringify({ result })
  })

/** POST /approvals/{id}/reject {reason} — reject with a reason. */
export const rejectApproval = (id: string, reason: string) =>
  request<{ status: string }>(`/approvals/${encodeURIComponent(id)}/reject`, {
    method: 'POST',
    body: JSON.stringify({ reason })
  })

/** GET /approvals/check/{tool} — is this tool considered dangerous? */
export const checkApprovalDangerous = (tool: string) =>
  request<{ tool: string; dangerous: boolean }>(`/approvals/check/${encodeURIComponent(tool)}`)

// ─── Gatekeepers — capability-based access control (internal/gatekeeper) ──────

/** gatekeeper.AuthType credential schemes. */
export type GatekeeperAuthType = 'oauth' | 'apikey' | 'basic'

/** gatekeeper.SessionStatus lifecycle states. */
export type GatekeeperSessionStatus = 'active' | 'revoked' | 'expired'

/** Matches gatekeeper.Capability (Go json tags). */
export type Capability = {
  tool: string[]
  scope: string[]
  read_only: boolean
  expires_at?: string // *time.Time, omitempty
}

/** Matches gatekeeper.GatekeeperConfig (Go json tags). Credentials are never
 *  serialized server-side (json:"-"). */
export type GatekeeperConfig = {
  name: string
  type: string
  auth_type: GatekeeperAuthType
  allowed_tools?: string[]
  scopes?: string[]
  require_approval: boolean
}

/** Matches gatekeeper.Session (Go json tags). */
export type GatekeeperSession = {
  id: string
  gatekeeper_name: string
  capabilities: Capability[]
  created_at: string // RFC3339
  expires_at?: string // *time.Time, omitempty
  status: GatekeeperSessionStatus
}

/** Matches gatekeeper.ActionLog (Go json tags). */
export type GatekeeperActionLog = {
  id: string
  session_id: string
  action: string
  resource: string
  result: string
  approved: boolean
  timestamp: string // RFC3339
  approval_id?: string // *string, omitempty
}

/** GET /gatekeepers — list all registered gatekeeper configs. */
export const listGatekeepers = () => request<GatekeeperConfig[]>('/gatekeepers')

/** POST /gatekeepers {GatekeeperConfig} — register (or replace) a gatekeeper. */
export const registerGatekeeper = (config: GatekeeperConfig) =>
  request<{ status: string; name: string }>('/gatekeepers', {
    method: 'POST',
    body: JSON.stringify(config)
  })

/** DELETE /gatekeepers/{name} — unregister a gatekeeper. */
export const removeGatekeeper = (name: string) =>
  request<{ status: string }>(`/gatekeepers/${encodeURIComponent(name)}`, { method: 'DELETE' })

/** POST /gatekeepers/{name}/access {Capability} — request an access session. */
export const requestGatekeeperAccess = (name: string, capability: Capability) =>
  request<GatekeeperSession>(`/gatekeepers/${encodeURIComponent(name)}/access`, {
    method: 'POST',
    body: JSON.stringify(capability)
  })

/** GET /gatekeepers/{name}/actions?session=X — audit trail for a session. */
export const getGatekeeperActions = (name: string, session?: string) =>
  request<GatekeeperActionLog[] | null>(
    `/gatekeepers/${encodeURIComponent(name)}/actions${session ? `?session=${encodeURIComponent(session)}` : ''}`
  )

/** POST /gatekeepers/{name}/sessions/{sid}/revoke — revoke a session. */
export const revokeGatekeeperSession = (name: string, sid: string) =>
  request<{ status: string }>(`/gatekeepers/${encodeURIComponent(name)}/sessions/${encodeURIComponent(sid)}/revoke`, {
    method: 'POST'
  })

// ─── Blueprints — reusable pentest playbooks (internal/blueprint) ────────────

/** blueprint.Category values (DB CHECK constraint). */
export type BlueprintCategory = 'recon' | 'exploit' | 'post-exploit' | 'reporting'

/** Matches blueprint.BlueprintStep (Go json tags). */
export type BlueprintStep = {
  order: number
  tool: string
  description: string
  args?: Record<string, string>
  expected_result?: string
  on_failure?: string
}

/** Matches blueprint.Blueprint (Go json tags). */
export type Blueprint = {
  id: string
  name: string
  description: string
  category: BlueprintCategory
  phase: string
  steps: BlueprintStep[]
  tags?: string[]
  version: string
  author: string
  created_at: string // RFC3339
}

/** GET /blueprints?category=X — list blueprints, optionally filtered. */
export const listBlueprints = (category?: BlueprintCategory) =>
  request<Blueprint[] | null>(`/blueprints${category ? `?category=${encodeURIComponent(category)}` : ''}`)

/** POST /blueprints {Blueprint} — create a blueprint. */
export const createBlueprint = (bp: Partial<Blueprint> & { name: string; category: BlueprintCategory }) =>
  request<Blueprint>('/blueprints', { method: 'POST', body: JSON.stringify(bp) })

/** GET /blueprints/{id} — fetch a single blueprint. */
export const getBlueprint = (id: string) => request<Blueprint>(`/blueprints/${encodeURIComponent(id)}`)

/** PUT /blueprints/{id} {Blueprint} — replace a blueprint's mutable fields. */
export const updateBlueprint = (id: string, bp: Partial<Blueprint> & { name: string; category: BlueprintCategory }) =>
  request<Blueprint>(`/blueprints/${encodeURIComponent(id)}`, { method: 'PUT', body: JSON.stringify(bp) })

/** DELETE /blueprints/{id} — remove a blueprint. */
export const deleteBlueprint = (id: string) =>
  request<{ status: string }>(`/blueprints/${encodeURIComponent(id)}`, { method: 'DELETE' })

// ─── Audit — tamper-evident compliance trail (internal/audit) ────────────────

/** audit.Actor values. */
export type AuditActor = 'user' | 'agent' | 'system'

/** audit.Severity values. */
export type AuditSeverity = 'info' | 'low' | 'medium' | 'high' | 'critical'

/** Matches audit.AuditEntry (Go json tags). */
export type AuditEntry = {
  id: string
  run_id: string
  actor: AuditActor
  action: string
  resource_type: string
  resource_id?: string
  details?: unknown // json.RawMessage, omitempty
  ip_address?: string
  timestamp: string // RFC3339
  severity: AuditSeverity
}

/** Shape of the audit export envelope (GET /audit/{run_id}/export). */
export type AuditExport = {
  run_id: string
  count: number
  entries: AuditEntry[] | null
}

/** Severity roll-up returned by GET /audit/{run_id}/stats. */
export type AuditStats = {
  info: number
  low: number
  medium: number
  high: number
  critical: number
  total: number
}

/** GET /audit/{run_id} — list audit entries for a run (oldest first). */
export const listAudit = (runId: string) =>
  request<AuditEntry[] | null>(`/audit/${encodeURIComponent(runId)}`)

/** POST /audit {AuditEntry} — append a compliance entry. */
export const logAudit = (entry: Partial<AuditEntry> & { run_id: string; actor: AuditActor; action: string }) =>
  request<AuditEntry>('/audit', { method: 'POST', body: JSON.stringify(entry) })

/** GET /audit/{run_id}/export — full audit trail as a JSON export envelope. */
export const exportAudit = (runId: string) =>
  request<AuditExport>(`/audit/${encodeURIComponent(runId)}/export`)

/** GET /audit/{run_id}/stats — per-severity counts plus a total. */
export const auditStats = (runId: string) =>
  request<AuditStats>(`/audit/${encodeURIComponent(runId)}/stats`)

// ─── MCP Gateway — tool classification & approval routing (internal/mcpgw) ───

/** mcpgw.TrustTier. */
export type MCPTrustTier = 'byo' | 'vetted'

/** mcpgw.ToolClassification. */
export type MCPToolClassification = 'observation' | 'action'

/** Matches mcpgw.ToolDescriptor (Go json tags). */
export type MCPToolDescriptor = {
  name: string
  description: string
  endpoint: string
  tier: MCPTrustTier
  class: MCPToolClassification
  hints?: Record<string, boolean> // readOnlyHint, destructiveHint, …
  schema?: Record<string, unknown>
  vetted: boolean
}

/** Matches mcpgw.CallRequest (Go json tags). */
export type MCPCallRequest = {
  tool_name: string
  args: Record<string, unknown>
  run_id: string
  caller: string // "agent" | "user"
}

/** Matches mcpgw.CallResult (Go json tags). */
export type MCPCallResult = {
  success: boolean
  output?: unknown
  error?: string
  approval_id?: string
  auto_approved?: boolean
  simulated?: boolean // true while awaiting approval
  timestamp: string // RFC3339
}

/** Tool breakdown returned by GET /mcp/stats. */
export type MCPStats = {
  total: number
  observations: number
  actions: number
  vetted: number
  byo: number
}

/** GET /mcp/tools — list all registered tool descriptors. */
export const listMCPTools = () => request<MCPToolDescriptor[]>('/mcp/tools')

/** POST /mcp/tools {ToolDescriptor} — register a tool. */
export const registerMCPTool = (tool: MCPToolDescriptor) =>
  request<MCPToolDescriptor>('/mcp/tools', { method: 'POST', body: JSON.stringify(tool) })

/** POST /mcp/call {CallRequest} — invoke a tool through the gateway. */
export const callMCPTool = (req: MCPCallRequest) =>
  request<MCPCallResult>('/mcp/call', { method: 'POST', body: JSON.stringify(req) })

/** POST /mcp/approve/{id} — execute a previously-queued (simulated) action. */
export const approveMCPAction = (id: string) =>
  request<MCPCallResult>(`/mcp/approve/${encodeURIComponent(id)}`, { method: 'POST' })

/** POST /mcp/reject/{id} {reason} — reject a queued action. */
export const rejectMCPAction = (id: string, reason: string) =>
  request<{ status: string }>(`/mcp/reject/${encodeURIComponent(id)}`, {
    method: 'POST',
    body: JSON.stringify({ reason })
  })

/** POST /mcp/vet/{tool} — promote a tool from byo to vetted tier. */
export const vetMCPTool = (tool: string) =>
  request<{ status: string; tool: string }>(`/mcp/vet/${encodeURIComponent(tool)}`, { method: 'POST' })

/** GET /mcp/stats — summary counts of registered tools. */
export const mcpStats = () => request<MCPStats>('/mcp/stats')

// ─── Sharing / Engagements — collaborative pentest scopes (internal/sharing) ─

/** sharing.Role values (hierarchy: owner > build > use). */
export type EngagementRole = 'owner' | 'build' | 'use'

/** Matches sharing.ShareLink (Go json tags). */
export type ShareLink = {
  id: string
  engagement_id: string
  role: EngagementRole
  token: string
  created_by: string
  created_at: string // RFC3339
  expires_at?: string // *time.Time, omitempty
  revoked: boolean
  revoked_at?: string // *time.Time, omitempty
  label?: string
}

/** Matches sharing.Collaborator (Go json tags). */
export type Collaborator = {
  user_id: string
  username: string
  role: EngagementRole
  granted_at: string // RFC3339
  granted_by: string
  share_link_id?: string
}

/** Matches sharing.Engagement (Go json tags). */
export type Engagement = {
  id: string
  name: string
  owner_id: string
  run_ids: string[]
  created_at: string // RFC3339
  metadata?: Record<string, string>
}

/** GET /engagements — list engagements visible to the current user. */
export const listEngagements = () => request<Engagement[] | null>('/engagements')

/** POST /engagements {Engagement} — create an engagement (owner_id set server-side). */
export const createEngagement = (eng: Partial<Engagement> & { name: string }) =>
  request<Engagement>('/engagements', { method: 'POST', body: JSON.stringify(eng) })

/** GET /engagements/{id} — fetch a single engagement. */
export const getEngagement = (id: string) => request<Engagement>(`/engagements/${encodeURIComponent(id)}`)

/** POST /engagements/{id}/shares {role, label} — mint a shareable link + token. */
export const createEngagementShare = (id: string, role: EngagementRole, label?: string) =>
  request<ShareLink>(`/engagements/${encodeURIComponent(id)}/shares`, {
    method: 'POST',
    body: JSON.stringify({ role, label })
  })

/** GET /engagements/{id}/shares — list share links for an engagement. */
export const listEngagementShares = (id: string) =>
  request<ShareLink[] | null>(`/engagements/${encodeURIComponent(id)}/shares`)

/** POST /engagements/{id}/shares/{linkID}/revoke — revoke a share link. */
export const revokeEngagementShare = (id: string, linkId: string) =>
  request<{ status: string; engagement: string }>(
    `/engagements/${encodeURIComponent(id)}/shares/${encodeURIComponent(linkId)}/revoke`,
    { method: 'POST' }
  )

/** POST /share/accept {token, username} — accept an invitation, become a collaborator. */
export const acceptShare = (token: string, username: string) =>
  request<Collaborator>('/share/accept', { method: 'POST', body: JSON.stringify({ token, username }) })

/** GET /engagements/{id}/collaborators — list collaborators on an engagement. */
export const listCollaborators = (id: string) =>
  request<Collaborator[] | null>(`/engagements/${encodeURIComponent(id)}/collaborators`)

/** DELETE /engagements/{id}/collaborators/{userID} — remove a collaborator. */
export const removeCollaborator = (id: string, userId: string) =>
  request<{ status: string }>(`/engagements/${encodeURIComponent(id)}/collaborators/${encodeURIComponent(userId)}`, {
    method: 'DELETE'
  })

/** POST /engagements/{id}/runs/{runID} — associate a run with an engagement. */
export const addRunToEngagement = (id: string, runId: string) =>
  request<{ status: string }>(`/engagements/${encodeURIComponent(id)}/runs/${encodeURIComponent(runId)}`, {
    method: 'POST'
  })
