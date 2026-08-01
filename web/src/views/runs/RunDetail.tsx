'use client'

// React Imports
import { useCallback, useEffect, useMemo, useRef, useState } from 'react'

// Third-party Imports
import { toast } from 'sonner'

// Type Imports
import type { Interrupt, RunSummary, StatusResponse, ToolCallRecord } from '@/lib/api'

// Component Imports
import Elapsed from '@/components/shared/Elapsed'
import StatusBadge from '@/components/shared/StatusBadge'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Collapsible, CollapsibleContent, CollapsibleTrigger } from '@/components/ui/collapsible'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle
} from '@/components/ui/dialog'
import { Skeleton } from '@/components/ui/skeleton'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { Textarea } from '@/components/ui/textarea'

// Util Imports
import { analyzeRun, getStatus, getTools, getTraces, listRuns, resumeRun, streamRun } from '@/lib/api'
import { shortId } from '@/lib/format'

const isActive = (status?: string) => status === 'running' || status === 'awaiting_approval' || status === 'initializing'

const prettyArgs = (args: ToolCallRecord['Args']) =>
  typeof args === 'string' ? args : JSON.stringify(args)

// Strip ANSI escape sequences (tool output is full of color codes).
const stripAnsi = (s: string) => s.replace(/\x1b\[[0-9;]*m/g, '')

/**
 * Render a tool's raw Output for the terminal. Arsenal tools return a JSON
 * envelope ({stdout, stderr, success, message}); strike tools return
 * {message, output, status}. Unwrap those into readable text, strip ANSI,
 * pretty-print other JSON, and pass plain text through untouched.
 */
const formatToolOutput = (raw: string): string => {
  if (!raw) return ''

  try {
    const parsed: unknown = JSON.parse(raw)

    if (parsed && typeof parsed === 'object' && !Array.isArray(parsed)) {
      const obj = parsed as Record<string, unknown>

      // Strike-style: prefer the captured command output.
      if (typeof obj.output === 'string' && obj.output.trim()) return stripAnsi(obj.output)

      // Arsenal-style envelope.
      const stdout = typeof obj.stdout === 'string' ? stripAnsi(obj.stdout) : ''
      const stderr = typeof obj.stderr === 'string' ? stripAnsi(obj.stderr) : ''

      if (stdout.trim()) {
        // Surface stderr too when the tool failed — that's where errors live.
        return obj.success === false && stderr.trim() ? `${stdout}\n[stderr]\n${stderr}` : stdout
      }

      if (stderr.trim()) return stderr
      if (typeof obj.message === 'string' && obj.message.trim()) return obj.message

      return JSON.stringify(parsed, null, 2)
    }
  } catch {
    // Not JSON — fall through to plain text.
  }

  return stripAnsi(raw)
}

// Trigger a client-side file download of text content.
const downloadText = (filename: string, content: string, mime = 'text/plain') => {
  const url = URL.createObjectURL(new Blob([content], { type: mime }))
  const a = document.createElement('a')

  a.href = url
  a.download = filename
  a.click()
  URL.revokeObjectURL(url)
}

/** Copy + Save actions for a text artifact (report / analysis / raw JSON). */
const ExportButtons = ({ filename, content, mime = 'text/plain' }: { filename: string; content: string; mime?: string }) => (
  <div className='flex shrink-0 gap-2'>
    <Button
      variant='outline'
      size='sm'
      className='font-mono text-[10px] tracking-widest uppercase'
      onClick={() =>
        navigator.clipboard
          .writeText(content)
          .then(() => toast.success('Copied to clipboard'))
          .catch(() => toast.error('Copy failed'))
      }
    >
      [ COPY ]
    </Button>
    <Button
      variant='outline'
      size='sm'
      className='font-mono text-[10px] tracking-widest uppercase'
      onClick={() => downloadText(filename, content, mime)}
    >
      [ ↓ SAVE ]
    </Button>
  </div>
)

const HitlGate = ({
  interrupt,
  onDecision,
  busy
}: {
  interrupt: Interrupt
  onDecision: (decision: 'approve' | 'reject' | 'edit', editedArgs?: Record<string, unknown>) => void
  busy: boolean
}) => {
  const [editOpen, setEditOpen] = useState(false)
  const [editText, setEditText] = useState('')
  const [editError, setEditError] = useState<string | null>(null)

  const openEdit = () => {
    setEditText(JSON.stringify(interrupt.Args, null, 2))
    setEditError(null)
    setEditOpen(true)
  }

  const submitEdit = () => {
    try {
      const parsed = JSON.parse(editText)

      if (typeof parsed !== 'object' || parsed === null || Array.isArray(parsed)) {
        setEditError('Edited args must be a JSON object')

        return
      }

      onDecision('edit', parsed)
      setEditOpen(false)
    } catch (err) {
      setEditError(err instanceof Error ? err.message : 'Invalid JSON')
    }
  }

  return (
    <div className='hud-corners hud-corners-amber border-warning/40 bg-warning/5 relative rounded-md border p-4 sm:p-5'>
      <div className='flex flex-col gap-4'>
        <div className='flex items-center gap-3'>
          <span className='text-warning text-xl'>⚠</span>
          <div>
            <p className='text-warning font-mono text-sm font-semibold tracking-widest uppercase'>
              HUMAN AUTHORIZATION REQUIRED
            </p>
            <p className='micro-label mt-0.5'>AGENT PAUSED — TOOL CALL PENDING OPERATOR DECISION</p>
          </div>
        </div>

        <div className='rounded-sm border border-white/5 bg-black/60 p-3'>
          <p className='text-primary font-mono text-xs tracking-widest'>$ {interrupt.ToolName}</p>
          <pre className='text-muted-foreground mt-2 overflow-x-auto font-mono text-xs whitespace-pre-wrap'>
            {JSON.stringify(interrupt.Args, null, 2)}
          </pre>
        </div>

        <div className='flex flex-wrap gap-2'>
          <Button
            disabled={busy}
            onClick={() => onDecision('approve')}
            className='font-mono text-xs font-semibold tracking-widest uppercase'
          >
            [ ✓ APPROVE ]
          </Button>
          <Button
            disabled={busy}
            variant='destructive'
            onClick={() => onDecision('reject')}
            className='font-mono text-xs font-semibold tracking-widest uppercase'
          >
            [ ✗ REJECT ]
          </Button>
          <Button
            disabled={busy}
            variant='outline'
            onClick={openEdit}
            className='font-mono text-xs font-semibold tracking-widest uppercase'
          >
            [ ✎ EDIT ARGS ]
          </Button>
        </div>
      </div>

      <Dialog open={editOpen} onOpenChange={setEditOpen}>
        <DialogContent className='sm:max-w-xl'>
          <DialogHeader>
            <DialogTitle className='font-mono text-sm tracking-widest uppercase'>EDIT TOOL ARGS</DialogTitle>
            <DialogDescription className='font-mono text-xs'>
              {interrupt.ToolName} — must be a valid JSON object.
            </DialogDescription>
          </DialogHeader>
          <Textarea
            value={editText}
            onChange={e => setEditText(e.target.value)}
            className='min-h-64 font-mono text-xs'
            spellCheck={false}
          />
          {editError && <p className='text-destructive font-mono text-xs'>{editError}</p>}
          <DialogFooter>
            <Button variant='outline' onClick={() => setEditOpen(false)} className='font-mono text-xs tracking-widest uppercase'>
              CANCEL
            </Button>
            <Button onClick={submitEdit} className='font-mono text-xs tracking-widest uppercase'>
              [ SEND EDITED ARGS ]
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  )
}

const Terminal = ({ tools, active }: { tools: ToolCallRecord[]; active: boolean }) => {
  const scrollRef = useRef<HTMLDivElement>(null)

  useEffect(() => {
    const el = scrollRef.current

    if (el && active) el.scrollTop = el.scrollHeight
  }, [tools, active])

  return (
    <div className='scanlines relative overflow-hidden rounded-md border bg-black/80'>
      <div ref={scrollRef} className='h-[32rem] overflow-y-auto p-3 font-mono text-[11px] leading-relaxed sm:p-4 sm:text-xs'>
        {tools.length === 0 ? (
          <p className='text-muted-foreground'>
            {active ? '// awaiting tool activity…' : '// no tool activity recorded'}
          </p>
        ) : (
          tools.map(tool => (
            <div key={tool.Index} className='mb-3'>
              <p className='text-primary text-glow'>
                $ {tool.ToolName} <span className='text-primary/70'>{prettyArgs(tool.Args)}</span>
              </p>
              {tool.Output && formatToolOutput(tool.Output) && (
                <pre className='text-muted-foreground mt-1 overflow-x-auto pl-4 whitespace-pre-wrap'>
                  {formatToolOutput(tool.Output)}
                </pre>
              )}
            </div>
          ))
        )}
        {active && (
          <p className='text-primary'>
            $ <span className='terminal-cursor' />
          </p>
        )}
      </div>
    </div>
  )
}

const RunDetail = ({ runId }: { runId: string }) => {
  const [summary, setSummary] = useState<RunSummary | null>(null)
  const [status, setStatus] = useState<StatusResponse | null>(null)
  const [tools, setTools] = useState<ToolCallRecord[]>([])
  const [traces, setTraces] = useState<string[] | null>(null)
  const [loaded, setLoaded] = useState(false)
  const [loadError, setLoadError] = useState<string | null>(null)
  const [resumeBusy, setResumeBusy] = useState(false)
  const [analysis, setAnalysis] = useState<string | null>(null)
  const [analysisBusy, setAnalysisBusy] = useState(false)
  const [analysisError, setAnalysisError] = useState<string | null>(null)

  const seenToolIndexes = useRef<Set<number>>(new Set())

  const ingestTool = useCallback((tool: ToolCallRecord) => {
    if (seenToolIndexes.current.has(tool.Index)) return
    seenToolIndexes.current.add(tool.Index)
    setTools(prev => [...prev, tool].sort((a, b) => a.Index - b.Index))
  }, [])

  // Initial load: status, tool log, run summary, traces
  useEffect(() => {
    let mounted = true

    Promise.allSettled([getStatus(runId), getTools(runId), listRuns(), getTraces(runId)]).then(
      ([statusRes, toolsRes, runsRes, tracesRes]) => {
        if (!mounted) return

        if (statusRes.status === 'fulfilled') setStatus(statusRes.value)
        else setLoadError(statusRes.reason instanceof Error ? statusRes.reason.message : String(statusRes.reason))

        if (toolsRes.status === 'fulfilled') {
          for (const tool of toolsRes.value.tool_log ?? []) {
            if (!seenToolIndexes.current.has(tool.Index)) {
              seenToolIndexes.current.add(tool.Index)
            }
          }

          setTools([...(toolsRes.value.tool_log ?? [])].sort((a, b) => a.Index - b.Index))
        }

        if (runsRes.status === 'fulfilled') {
          setSummary((runsRes.value.runs ?? []).find(r => r.run_id === runId) ?? null)
        }

        if (tracesRes.status === 'fulfilled') setTraces(tracesRes.value.history ?? [])

        setLoaded(true)
      }
    )

    return () => {
      mounted = false
    }
  }, [runId])

  // Live stream (with polling fallback handled by streamRun)
  const currentStatus = status?.status

  useEffect(() => {
    if (!loaded) return
    if (currentStatus === 'not_found') return
    if (currentStatus && !isActive(currentStatus)) return

    const stop = streamRun(runId, {
      onTool: ingestTool,
      onStatus: s => {
        setStatus(prev => ({ ...prev, ...s }))

        if (!isActive(s.status)) {
          getTraces(runId)
            .then(res => setTraces(res.history ?? []))
            .catch(() => {})
        }
      },
      onError: err => setLoadError(err.message)
    })

    return stop
  }, [loaded, runId, currentStatus, ingestTool])

  const handleDecision = useCallback(
    async (decision: 'approve' | 'reject' | 'edit', editedArgs?: Record<string, unknown>) => {
      setResumeBusy(true)

      try {
        await resumeRun(runId, { decision, ...(editedArgs ? { edited_args: editedArgs } : {}) })
        toast.success(`Decision sent: ${decision.toUpperCase()}`)
        setStatus(prev => (prev ? { ...prev, status: 'running', interrupt: null } : prev))
      } catch (err) {
        toast.error(err instanceof Error ? err.message : 'Failed to send decision')
      } finally {
        setResumeBusy(false)
      }
    },
    [runId]
  )

  const rawJson = useMemo(() => JSON.stringify({ status, tools }, null, 2), [status, tools])

  const handleAnalyze = useCallback(async () => {
    setAnalysisBusy(true)
    setAnalysisError(null)

    try {
      const res = await analyzeRun(runId)

      setAnalysis(res.analysis)
    } catch (err) {
      setAnalysisError(err instanceof Error ? err.message : 'Analysis failed')
    } finally {
      setAnalysisBusy(false)
    }
  }, [runId])

  if (loaded && status?.status === 'not_found') {
    return (
      <div className='grid-bg flex h-96 flex-col items-center justify-center gap-4 rounded-md border'>
        <p className='text-destructive font-mono text-2xl font-semibold tracking-widest'>RUN NOT FOUND</p>
        <p className='micro-label'>NO OPERATION WITH ID {shortId(runId, 16)}</p>
        <Button
          variant='outline'
          className='font-mono text-xs tracking-widest uppercase'
          render={<a href='/runs' />}
          nativeButton={false}
        >
          [ ← BACK TO OPERATIONS ]
        </Button>
      </div>
    )
  }

  return (
    <div className='flex flex-col gap-6'>
      {/* Header */}
      <div className='flex flex-wrap items-center justify-between gap-4'>
        <div className='flex flex-col gap-1'>
          <div className='flex flex-wrap items-center gap-3'>
            <h1 className='font-mono text-xl font-semibold tracking-widest'>
              {summary?.target ?? <Skeleton className='inline-block h-6 w-36 align-middle' />}
            </h1>
            {status && <StatusBadge status={status.status} />}
          </div>
          <p className='micro-label'>
            RUN {shortId(runId, 16)}
            {summary?.cve_id && ` — ${summary.cve_id}`}
            {summary?.service_name && ` — ${summary.service_name}`}
          </p>
        </div>
        {summary && (
          <div className='text-right'>
            <p className='micro-label'>ELAPSED</p>
            <Elapsed
              since={summary.started_at}
              until={summary.ended_at}
              className='text-primary font-mono text-lg tracking-widest'
            />
          </div>
        )}
      </div>

      {loadError && !status && (
        <div className='border-destructive/40 bg-destructive/10 text-destructive rounded-md border px-4 py-3 font-mono text-xs tracking-widest uppercase'>
          CORE UNREACHABLE — {loadError}
        </div>
      )}

      {/* HITL gate */}
      {status?.status === 'awaiting_approval' && status.interrupt && (
        <HitlGate interrupt={status.interrupt} onDecision={handleDecision} busy={resumeBusy} />
      )}

      {/* Body tabs */}
      <Tabs defaultValue='terminal'>
        <TabsList
          variant='line'
          className='w-full justify-start overflow-x-auto border-b font-mono text-xs tracking-widest uppercase'
        >
          <TabsTrigger value='terminal'>Terminal</TabsTrigger>
          <TabsTrigger value='report'>Report</TabsTrigger>
          <TabsTrigger value='analysis'>AI Analysis</TabsTrigger>
          <TabsTrigger value='traces'>Traces</TabsTrigger>
          <TabsTrigger value='raw'>Raw</TabsTrigger>
        </TabsList>

        <TabsContent value='terminal' className='pt-4'>
          {!loaded ? <Skeleton className='h-[32rem] w-full' /> : <Terminal tools={tools} active={isActive(status?.status)} />}
        </TabsContent>

        <TabsContent value='report' className='flex flex-col gap-4 pt-4'>
          {status?.judge_verdict === true && (
            <div className='hud-corners border-primary/40 bg-primary/10 text-primary rounded-md border px-4 py-3 font-mono text-sm font-semibold tracking-widest uppercase'>
              ✓ TARGET COMPROMISED — PoC CONFIRMED
            </div>
          )}
          {status?.judge_verdict === false && status.status === 'completed' && (
            <div className='border-destructive/40 bg-destructive/10 text-destructive rounded-md border px-4 py-3 font-mono text-sm font-semibold tracking-widest uppercase'>
              ✗ NOT COMPROMISED
            </div>
          )}
          <Card>
            <CardHeader className='flex flex-row items-center justify-between gap-3'>
              <CardTitle className='micro-label'>FINAL REPORT</CardTitle>
              {status?.output && (
                <ExportButtons
                  filename={`talon-${shortId(runId, 8)}-report.md`}
                  content={status.output}
                  mime='text/markdown'
                />
              )}
            </CardHeader>
            <CardContent>
              {!loaded ? (
                <Skeleton className='h-48 w-full' />
              ) : status?.output ? (
                <pre className='max-h-[32rem] overflow-auto font-mono text-xs leading-relaxed whitespace-pre-wrap'>
                  {status.output}
                </pre>
              ) : (
                <p className='micro-label py-8 text-center'>
                  {isActive(status?.status) ? 'REPORT PENDING — OPERATION IN PROGRESS' : 'NO REPORT AVAILABLE'}
                </p>
              )}
            </CardContent>
          </Card>
        </TabsContent>

        <TabsContent value='analysis' className='pt-4'>
          <Card className='hud-corners'>
            <CardHeader>
              <div className='flex flex-wrap items-center justify-between gap-3'>
                <div>
                  <CardTitle className='micro-label'>AI ANALYST BRIEFING</CardTitle>
                  <p className='micro-label mt-1'>ONE-SHOT LLM ANALYSIS — REPORT + TOOL LOG, SERVER-SIDE MODEL — NOT PERSISTED</p>
                </div>
                <div className='flex shrink-0 items-center gap-2'>
                  {analysis && (
                    <ExportButtons filename={`talon-${shortId(runId, 8)}-analysis.txt`} content={analysis} />
                  )}
                  <Button
                    onClick={handleAnalyze}
                    disabled={analysisBusy || !loaded}
                    className='font-mono text-xs font-semibold tracking-widest uppercase'
                  >
                    {analysisBusy ? '[ ⟳ ANALYZING… ]' : analysis ? '[ ✦ RE-ANALYZE ]' : '[ ✦ ANALYZE RUN ]'}
                  </Button>
                </div>
              </div>
            </CardHeader>
            <CardContent>
              {analysisError && (
                <div className='border-destructive/40 bg-destructive/10 text-destructive mb-4 rounded-md border px-4 py-3 font-mono text-xs tracking-widest uppercase'>
                  ANALYSIS FAILED — {analysisError}
                </div>
              )}
              {analysisBusy ? (
                <div className='flex flex-col gap-2'>
                  <Skeleton className='h-4 w-3/4' />
                  <Skeleton className='h-4 w-full' />
                  <Skeleton className='h-4 w-5/6' />
                  <Skeleton className='h-4 w-2/3' />
                  <p className='micro-label mt-2'>QUERYING ANALYST MODEL — THIS CAN TAKE UP TO 2 MINUTES</p>
                </div>
              ) : analysis ? (
                <pre className='max-h-[32rem] overflow-auto font-mono text-xs leading-relaxed whitespace-pre-wrap'>
                  {analysis}
                </pre>
              ) : (
                <p className='micro-label py-8 text-center'>
                  NO ANALYSIS YET — PRESS ANALYZE TO GENERATE AN OPERATOR BRIEFING
                </p>
              )}
            </CardContent>
          </Card>
        </TabsContent>

        <TabsContent value='traces' className='pt-4'>
          <Card>
            <CardHeader>
              <CardTitle className='micro-label'>AGENT TRACE LOG</CardTitle>
            </CardHeader>
            <CardContent>
              {traces === null ? (
                <Skeleton className='h-48 w-full' />
              ) : traces.length === 0 ? (
                <p className='micro-label py-8 text-center'>NO TRACES RECORDED</p>
              ) : (
                <div className='flex max-h-[32rem] flex-col gap-1 overflow-y-auto font-mono text-xs'>
                  {traces.map((trace, i) => (
                    <p key={i} className='text-muted-foreground border-l border-white/10 py-1 pl-3 break-words'>
                      {trace}
                    </p>
                  ))}
                </div>
              )}
            </CardContent>
          </Card>
        </TabsContent>

        <TabsContent value='raw' className='pt-4'>
          <Collapsible defaultOpen>
            <Card>
              <CardHeader className='flex flex-row items-center justify-between gap-3'>
                <CollapsibleTrigger className='flex-1 text-left'>
                  <CardTitle className='micro-label hover:text-foreground'>STATUS + TOOL LOG (JSON) ▾</CardTitle>
                </CollapsibleTrigger>
                {loaded && (
                  <ExportButtons
                    filename={`talon-${shortId(runId, 8)}-raw.json`}
                    content={rawJson}
                    mime='application/json'
                  />
                )}
              </CardHeader>
              <CollapsibleContent>
                <CardContent>
                  <pre className='text-muted-foreground max-h-[32rem] overflow-auto rounded-sm bg-black/60 p-4 font-mono text-xs'>
                    {loaded ? rawJson : '// loading…'}
                  </pre>
                </CardContent>
              </CollapsibleContent>
            </Card>
          </Collapsible>
        </TabsContent>
      </Tabs>
    </div>
  )
}

export default RunDetail
