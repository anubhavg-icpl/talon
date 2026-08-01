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
