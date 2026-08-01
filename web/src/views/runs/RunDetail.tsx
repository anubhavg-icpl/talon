'use client'

// React Imports
import { useCallback, useEffect, useMemo, useRef, useState } from 'react'

// Third-party Imports
import { toast } from 'sonner'

// Type Imports
import type {
  Finding,
  Interrupt,
  KillChainAnalysis,
  MethodologyState,
  RunSummary,
  StatusResponse,
  StructuredReport,
  ToolCallRecord
} from '@/lib/api'

// Component Imports
import Elapsed from '@/components/shared/Elapsed'
import StatusBadge from '@/components/shared/StatusBadge'
import { Badge } from '@/components/ui/badge'
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
import {
  addNote,
  analyzeRun,
  exportRun,
  getFindings,
  getKillChain,
  getMethodology,
  getNotes,
  getReport,
  getStatus,
  getTimeline,
  getTools,
  getTraces,
  listRuns,
  resumeRun,
  streamRun,
  triageFinding
} from '@/lib/api'
import type { OperatorNote, TimelineEvent } from '@/lib/api'
import { shortId } from '@/lib/format'

const severityVariant = (sev: string): 'destructive' | 'default' | 'secondary' | 'outline' => {
  switch (sev.toLowerCase()) {
    case 'critical':
    case 'high':
      return 'destructive'
    case 'medium':
      return 'default'
    case 'low':
      return 'secondary'
    default:
      return 'outline'
  }
}

const FindingCard = ({ finding }: { finding: Finding }) => (
  <Card className='border-border/60'>
    <CardHeader className='pb-2'>
      <div className='flex flex-wrap items-center gap-2'>
        <Badge variant={severityVariant(finding.severity)} className='font-mono text-[10px] tracking-widest uppercase'>
          {finding.severity}
        </Badge>
        {finding.evidence?.passed && (
          <Badge variant='outline' className='border-primary/50 text-primary font-mono text-[10px] tracking-widest uppercase'>
            3-GATE PASS
          </Badge>
        )}
        {finding.stage && (
          <span className='text-muted-foreground font-mono text-[10px] tracking-widest uppercase'>{finding.stage}</span>
        )}
        <span className='text-muted-foreground ml-auto font-mono text-[10px]'>{finding.id}</span>
      </div>
      <CardTitle className='mt-2 font-mono text-sm font-semibold tracking-wide'>{finding.title}</CardTitle>
    </CardHeader>
    <CardContent className='space-y-2 font-mono text-xs'>
      {finding.description && <p className='text-muted-foreground leading-relaxed'>{finding.description}</p>}
      <div className='text-muted-foreground flex flex-wrap gap-x-4 gap-y-1 text-[10px] tracking-widest uppercase'>
        {finding.endpoint && <span>EP: {finding.endpoint}</span>}
        {finding.cwe_id && <span>REF: {finding.cwe_id}</span>}
        {finding.source && <span>SRC: {finding.source}</span>}
      </div>
      {(finding.evidence?.baseline || finding.evidence?.attack || finding.evidence?.diff) && (
        <div className='border-border/40 bg-muted/20 mt-2 space-y-1 rounded border p-2 text-[11px]'>
          <p className='micro-label'>3-GATE EVIDENCE</p>
          {finding.evidence.baseline && (
            <p>
              <span className='text-primary'>1. BASELINE:</span> {finding.evidence.baseline}
            </p>
          )}
          {finding.evidence.attack && (
            <p>
              <span className='text-primary'>2. ATTACK:</span> {finding.evidence.attack}
            </p>
          )}
          {finding.evidence.diff && (
            <p>
              <span className='text-primary'>3. DIFF:</span> {finding.evidence.diff}
            </p>
          )}
        </div>
      )}
      {finding.recommendation && (
        <p className='text-muted-foreground'>
          <span className='text-foreground'>REMEDY:</span> {finding.recommendation}
        </p>
      )}
    </CardContent>
  </Card>
)

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
          <p className='text-muted-foreground'>{'// awaiting tool activity…'}</p>
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
  const [findings, setFindings] = useState<Finding[]>([])
  const [report, setReport] = useState<StructuredReport | null>(null)
  const [killchain, setKillchain] = useState<KillChainAnalysis | null>(null)
  const [methodology, setMethodology] = useState<MethodologyState | null>(null)
  const [timeline, setTimeline] = useState<TimelineEvent[]>([])
  const [notes, setNotes] = useState<OperatorNote[]>([])
  const [noteDraft, setNoteDraft] = useState('')

  const seenToolIndexes = useRef<Set<number>>(new Set())

  const loadFindingsReport = useCallback(async () => {
    const [fRes, rRes, kRes, mRes, tRes, nRes] = await Promise.allSettled([
      getFindings(runId),
      getReport(runId),
      getKillChain(runId),
      getMethodology(runId),
      getTimeline(runId),
      getNotes(runId)
    ])
    if (fRes.status === 'fulfilled') setFindings(fRes.value.findings ?? [])
    if (rRes.status === 'fulfilled') setReport(rRes.value)
    if (kRes.status === 'fulfilled') setKillchain(kRes.value)
    if (mRes.status === 'fulfilled') setMethodology(mRes.value)
    if (tRes.status === 'fulfilled') setTimeline(tRes.value.timeline ?? [])
    if (nRes.status === 'fulfilled') setNotes(nRes.value.notes ?? [])
  }, [runId])

  const exportBundle = async () => {
    try {
      const bundle = await exportRun(runId)
      const blob = new Blob([JSON.stringify(bundle, null, 2)], { type: 'application/json' })
      const url = URL.createObjectURL(blob)
      const a = document.createElement('a')
      a.href = url
      a.download = `talon-export-${runId.slice(0, 8)}.json`
      a.click()
      URL.revokeObjectURL(url)
      toast.success('Export bundle downloaded')
    } catch (err) {
      toast.error(err instanceof Error ? err.message : 'Export failed')
    }
  }

  const submitNote = async () => {
    if (!noteDraft.trim()) return
    try {
      await addNote(runId, noteDraft.trim())
      setNoteDraft('')
      const n = await getNotes(runId)
      setNotes(n.notes ?? [])
      toast.success('Note saved')
    } catch (err) {
      toast.error(err instanceof Error ? err.message : 'Note failed')
    }
  }

  const exportReport = () => {
    const md = report?.markdown || status?.output || ''
    if (!md) {
      toast.error('No report to export')
      return
    }
    const blob = new Blob([md], { type: 'text/markdown' })
    const url = URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = `talon-report-${runId.slice(0, 8)}.md`
    a.click()
    URL.revokeObjectURL(url)
    toast.success('Report exported')
  }

  const exportFindingsJSON = () => {
    if (!findings.length) {
      toast.error('No findings to export')
      return
    }
    const blob = new Blob([JSON.stringify({ run_id: runId, findings }, null, 2)], {
      type: 'application/json'
    })
    const url = URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = `talon-findings-${runId.slice(0, 8)}.json`
    a.click()
    URL.revokeObjectURL(url)
    toast.success('Findings JSON exported')
  }

  const handleTriage = async (id: string, st: string) => {
    try {
      await triageFinding(runId, id, st)
      toast.success(`Triaged ${id} → ${st}`)
      await loadFindingsReport()
    } catch (err) {
      toast.error(err instanceof Error ? err.message : 'Triage failed')
    }
  }

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

  // Structured findings + report once the run is terminal (or has a report flag).
  useEffect(() => {
    if (!loaded) return
    const st = status?.status
    if (st === 'completed' || st === 'error' || status?.has_report || (status?.findings_count ?? 0) > 0) {
      void loadFindingsReport()
    }
  }, [loaded, status?.status, status?.has_report, status?.findings_count, loadFindingsReport])

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
        if (s.status === 'completed' || s.status === 'error') {
          void loadFindingsReport()
        }

        if (!isActive(s.status)) {
          getTraces(runId)
            .then(res => setTraces(res.history ?? []))
            .catch(() => {})
        }
      },
      onFindings: payload => {
        if (payload.findings) setFindings(payload.findings)
        setStatus(prev =>
          prev
            ? {
                ...prev,
                findings_count: payload.findings_count,
                findings_summary: payload.findings_summary
              }
            : prev
        )
      },
      onError: err => setLoadError(err.message)
    })

    return stop
  }, [loaded, runId, currentStatus, ingestTool, loadFindingsReport])

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
            <Elapsed since={summary.started_at} className='text-primary font-mono text-lg tracking-widest' />
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
          <TabsTrigger value='findings'>
            Findings{findings.length > 0 ? ` (${findings.length})` : status?.findings_count ? ` (${status.findings_count})` : ''}
          </TabsTrigger>
          <TabsTrigger value='killchain'>Kill Chain</TabsTrigger>
          <TabsTrigger value='timeline'>Timeline</TabsTrigger>
          <TabsTrigger value='notes'>Notes</TabsTrigger>
          <TabsTrigger value='report'>Report</TabsTrigger>
          <TabsTrigger value='analysis'>AI Analysis</TabsTrigger>
          <TabsTrigger value='traces'>Traces</TabsTrigger>
          <TabsTrigger value='raw'>Raw</TabsTrigger>
        </TabsList>

        <TabsContent value='terminal' className='pt-4'>
          {!loaded ? <Skeleton className='h-[32rem] w-full' /> : <Terminal tools={tools} active={isActive(status?.status)} />}
        </TabsContent>

        <TabsContent value='findings' className='flex flex-col gap-4 pt-4'>
          <div className='flex flex-wrap items-center justify-between gap-2'>
            {status?.findings_summary ? (
              <div className='flex flex-wrap gap-2 font-mono text-[10px] tracking-widest uppercase'>
                <Badge variant='outline'>TOTAL {status.findings_summary.total}</Badge>
                {status.findings_summary.critical > 0 && (
                  <Badge variant='destructive'>CRIT {status.findings_summary.critical}</Badge>
                )}
                {status.findings_summary.high > 0 && (
                  <Badge variant='destructive'>HIGH {status.findings_summary.high}</Badge>
                )}
                {status.findings_summary.medium > 0 && (
                  <Badge variant='default'>MED {status.findings_summary.medium}</Badge>
                )}
                {status.findings_summary.low > 0 && (
                  <Badge variant='secondary'>LOW {status.findings_summary.low}</Badge>
                )}
                {status.findings_summary.info > 0 && (
                  <Badge variant='outline'>INFO {status.findings_summary.info}</Badge>
                )}
                <Badge variant='outline' className='border-primary/40 text-primary'>
                  3-GATE {status.findings_summary.confirmed}
                </Badge>
              </div>
            ) : (
              <span />
            )}
            <Button
              variant='outline'
              size='sm'
              onClick={exportFindingsJSON}
              className='font-mono text-[10px] tracking-widest uppercase'
            >
              Export JSON
            </Button>
          </div>
          {!loaded ? (
            <Skeleton className='h-48 w-full' />
          ) : findings.length === 0 ? (
            <p className='micro-label py-8 text-center'>
              {isActive(status?.status)
                ? 'FINDINGS PENDING — OPERATION IN PROGRESS'
                : 'NO STRUCTURED FINDINGS — COMPLETE A RUN TO EXTRACT EVIDENCE'}
            </p>
          ) : (
            <div className='flex flex-col gap-3'>
              {findings.map(f => (
                <div key={f.id} className='flex flex-col gap-2'>
                  <FindingCard finding={f} />
                  {f.status === 'new' && (
                    <div className='flex gap-2'>
                      <Button
                        size='sm'
                        variant='outline'
                        className='font-mono text-[10px] tracking-widest uppercase'
                        onClick={() => handleTriage(f.id, 'approved')}
                      >
                        Approve
                      </Button>
                      <Button
                        size='sm'
                        variant='outline'
                        className='font-mono text-[10px] tracking-widest uppercase'
                        onClick={() => handleTriage(f.id, 'duplicate')}
                      >
                        Mark Duplicate
                      </Button>
                    </div>
                  )}
                </div>
              ))}
            </div>
          )}
        </TabsContent>

        <TabsContent value='killchain' className='flex flex-col gap-4 pt-4'>
          {methodology && (
            <Card>
              <CardHeader>
                <CardTitle className='micro-label'>
                  METHODOLOGY COVERAGE — {methodology.percent}% ({methodology.covered_count}/{methodology.total_count})
                </CardTitle>
              </CardHeader>
              <CardContent className='space-y-2 font-mono text-xs'>
                {methodology.items?.map(it => (
                  <div key={it.stage} className='flex flex-wrap items-center gap-2'>
                    <span>{it.covered ? '☑' : '☐'}</span>
                    <span>{it.label}</span>
                    {it.tools && it.tools.length > 0 && (
                      <span className='text-muted-foreground'>({it.tools.join(', ')})</span>
                    )}
                    {it.notes && <span className='text-muted-foreground text-[10px]'>{it.notes}</span>}
                  </div>
                ))}
              </CardContent>
            </Card>
          )}
          <Card>
            <CardHeader>
              <CardTitle className='micro-label'>
                KILL CHAIN{killchain?.max_severity ? ` — MAX ${killchain.max_severity.toUpperCase()}` : ''}
              </CardTitle>
            </CardHeader>
            <CardContent>
              {!killchain ? (
                <p className='micro-label py-6 text-center'>KILL CHAIN PENDING</p>
              ) : (
                <div className='space-y-3 font-mono text-xs'>
                  {(killchain.chains ?? []).length === 0 ? (
                    <p className='text-muted-foreground'>No multi-stage chains detected yet.</p>
                  ) : (
                    killchain.chains.map((c, i) => (
                      <div key={i} className='border-border/50 rounded border p-2'>
                        <span className='text-primary uppercase'>{c.severity}</span>:{' '}
                        <code>
                          {c.from} → {c.to}
                        </code>
                        <p className='text-muted-foreground mt-1'>{c.reason}</p>
                      </div>
                    ))
                  )}
                  {killchain.next_steps && killchain.next_steps.length > 0 && (
                    <p>
                      <span className='micro-label'>NEXT STEPS: </span>
                      {killchain.next_steps.join(', ')}
                    </p>
                  )}
                  {killchain.summary && (
                    <pre className='text-muted-foreground max-h-48 overflow-auto whitespace-pre-wrap'>{killchain.summary}</pre>
                  )}
                </div>
              )}
            </CardContent>
          </Card>
        </TabsContent>

        <TabsContent value='timeline' className='flex flex-col gap-3 pt-4'>
          {timeline.length === 0 ? (
            <p className='micro-label py-8 text-center'>NO TIMELINE EVENTS YET</p>
          ) : (
            timeline.map(ev => (
              <div key={`${ev.kind}-${ev.index}`} className='border-border/50 rounded border px-3 py-2 font-mono text-xs'>
                <div className='flex flex-wrap gap-2'>
                  <Badge variant={ev.kind === 'finding' ? 'destructive' : 'outline'} className='text-[9px] uppercase'>
                    {ev.kind}
                  </Badge>
                  {ev.stage && (
                    <Badge variant='secondary' className='text-[9px]'>
                      {ev.stage}
                    </Badge>
                  )}
                  {ev.severity && (
                    <Badge className='text-[9px] uppercase'>{ev.severity}</Badge>
                  )}
                  <span className='text-muted-foreground ml-auto text-[10px]'>#{ev.index}</span>
                </div>
                <p className='mt-1 font-medium'>{ev.label}</p>
                {ev.detail && <p className='text-muted-foreground mt-0.5 line-clamp-3 text-[11px]'>{ev.detail}</p>}
              </div>
            ))
          )}
        </TabsContent>

        <TabsContent value='notes' className='flex flex-col gap-4 pt-4'>
          <Card>
            <CardHeader>
              <CardTitle className='micro-label'>OPERATOR NOTES</CardTitle>
            </CardHeader>
            <CardContent className='space-y-3'>
              <Textarea
                value={noteDraft}
                onChange={e => setNoteDraft(e.target.value)}
                placeholder='Engagement context, HITL rationale, customer notes…'
                className='min-h-24 font-mono text-xs'
              />
              <Button
                size='sm'
                onClick={submitNote}
                className='font-mono text-[10px] tracking-widest uppercase'
              >
                Save note
              </Button>
              <div className='space-y-2 pt-2'>
                {notes.length === 0 ? (
                  <p className='micro-label py-4 text-center'>NO NOTES YET</p>
                ) : (
                  notes.map(n => (
                    <div key={n.id} className='border-border/50 rounded border px-3 py-2 font-mono text-xs'>
                      <div className='text-muted-foreground flex justify-between text-[10px]'>
                        <span>{n.id} · {n.author || 'operator'}</span>
                        <span>{n.created_at}</span>
                      </div>
                      <p className='mt-1 whitespace-pre-wrap'>{n.body}</p>
                    </div>
                  ))
                )}
              </div>
            </CardContent>
          </Card>
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
            <CardHeader>
              <div className='flex flex-wrap items-center justify-between gap-2'>
                <div>
                  <CardTitle className='micro-label'>STRUCTURED VALIDATION REPORT</CardTitle>
                  {report?.stages_covered && report.stages_covered.length > 0 && (
                    <p className='micro-label mt-1'>STAGES: {report.stages_covered.join(' → ')}</p>
                  )}
                  {status?.agent_mode && (
                    <p className='micro-label mt-1'>AGENT MODE: {status.agent_mode.toUpperCase()}</p>
                  )}
                </div>
                <div className='flex gap-2'>
                  <Button
                    variant='outline'
                    size='sm'
                    onClick={exportReport}
                    className='font-mono text-[10px] tracking-widest uppercase'
                  >
                    Export .md
                  </Button>
                  <Button
                    variant='outline'
                    size='sm'
                    onClick={exportBundle}
                    className='font-mono text-[10px] tracking-widest uppercase'
                  >
                    Export JSON
                  </Button>
                </div>
              </div>
            </CardHeader>
            <CardContent>
              {!loaded ? (
                <Skeleton className='h-48 w-full' />
              ) : report?.markdown ? (
                <pre className='max-h-[32rem] overflow-auto font-mono text-xs leading-relaxed whitespace-pre-wrap'>
                  {report.markdown}
                </pre>
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
                  <p className='micro-label mt-1'>ONE-SHOT LLM ANALYSIS — REPORT + TOOL LOG, SERVER-SIDE MODEL</p>
                </div>
                <Button
                  onClick={handleAnalyze}
                  disabled={analysisBusy || !loaded}
                  className='font-mono text-xs font-semibold tracking-widest uppercase'
                >
                  {analysisBusy ? '[ ⟳ ANALYZING… ]' : analysis ? '[ ✦ RE-ANALYZE ]' : '[ ✦ ANALYZE RUN ]'}
                </Button>
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
              <CardHeader>
                <CollapsibleTrigger className='w-full text-left'>
                  <CardTitle className='micro-label hover:text-foreground'>STATUS + TOOL LOG (JSON) ▾</CardTitle>
                </CollapsibleTrigger>
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
