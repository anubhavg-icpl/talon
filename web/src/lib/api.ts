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

export type StatusResponse = {
  status: RunStatus
  output: string
  interrupt: Interrupt | null
  judge_verdict?: boolean
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

export const getMCPServers = () => request<{ servers: MCPServerInfo[] }>('/mcp/servers')

export const startRun = (body: StartRunRequest) =>
  request<StartRunResponse>('/input/start', { method: 'POST', body: JSON.stringify(body) })

export const getStatus = (runId: string) => request<StatusResponse>(`/output/status/${runId}`)

export const resumeRun = (runId: string, body: ResumeRequest) =>
  request<unknown>(`/output/resume/${runId}`, { method: 'POST', body: JSON.stringify(body) })

export const getTools = (runId: string) => request<ToolsResponse>(`/monitor/tools?run_id=${encodeURIComponent(runId)}`)

export const getTraces = (runId: string) => request<TracesResponse>(`/monitor/traces/${runId}`)

export type AnalyzeResponse = {
  run_id: string
  analysis: string
}

/** One-shot LLM analysis of a run (report + tool log) — server-side analyst model. */
export const analyzeRun = (runId: string) =>
  request<AnalyzeResponse>(`/analyze/${runId}`, { method: 'POST' })

export type StreamRunHandlers = {
  onTool?: (tool: ToolCallRecord) => void
  onStatus?: (status: StatusResponse) => void
  onError?: (error: Error) => void
}

const TERMINAL_STATUSES: RunStatus[] = ['completed', 'error', 'not_found']

/**
 * Live-follow a run with a resilience chain:
 *   WebSocket (direct to talon-core /monitor/ws) → SSE (proxied /monitor/stream) → 3s polling.
 * Each tier falls through to the next on error before a terminal status, so a
 * broken WS or SSE path never breaks the view. Returns an unsubscribe function.
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
        const [status, tools] = await Promise.all([getStatus(runId), getTools(runId)])

        if (stopped) return

        for (const tool of tools.tool_log ?? []) {
          if (tool.Index > lastToolIndex) handleTool(tool)
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

      es.onerror = () => {
        // SSE endpoint unreachable or dropped — switch to polling
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

          if (msg.type === 'tool') handleTool(msg.data)
          else if (msg.type === 'status') handleStatus(msg.data)
        } catch {
          // ignore malformed messages
        }
      }

      ws.onerror = () => {
        // Defer to onclose for the fallback decision
      }

      ws.onclose = () => {
        ws = null

        // Clean close after terminal status sets stopped=true; any other
        // close means WS is unavailable — degrade to SSE.
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
