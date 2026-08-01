'use client'

import { useMemo, useState } from 'react'
import { CheckCircle2, Circle, FileText, Shield, Target, XCircle } from 'lucide-react'

import { Badge } from '@/components/ui/badge'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Collapsible, CollapsibleContent, CollapsibleTrigger } from '@/components/ui/collapsible'
import type { Finding, KillChainAnalysis, MethodologyState, StructuredReport, ToolCallRecord } from '@/lib/api'
import { cn } from '@/lib/utils'
import { AssistantMarkdown } from '@/views/assist/ToolResultViz'

const sevTone = (sev: string) => {
  const s = sev.toLowerCase()
  if (s === 'critical') return 'bg-red-500/15 text-red-400 border-red-500/40'
  if (s === 'high') return 'bg-orange-500/15 text-orange-400 border-orange-500/40'
  if (s === 'medium') return 'bg-amber-500/15 text-amber-400 border-amber-500/40'
  if (s === 'low') return 'bg-red-300/20 text-red-300 border-red-300/40'
  return 'bg-muted text-muted-foreground border-border'
}

const sevVariant = (sev: string): 'destructive' | 'default' | 'secondary' | 'outline' => {
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

const PIPELINE = ['recon', 'exploit', 'post_exploit', 'codegen', 'report'] as const

type Props = {
  report: StructuredReport | null
  findings?: Finding[] | null
  tools?: ToolCallRecord[] | null
  methodology?: MethodologyState | null
  killChain?: KillChainAnalysis | null
  agentMode?: string
  fallbackMarkdown?: string
}

const ReportFindingCard = ({ f }: { f: Finding }) => {
  const [open, setOpen] = useState(false)
  const gatePass = f.evidence?.passed

  return (
    <Collapsible open={open} onOpenChange={setOpen}>
      <Card className='border-border/60 overflow-hidden'>
        <CollapsibleTrigger className='hover:bg-muted/20 w-full text-left transition-colors'>
          <CardHeader className='pb-2'>
            <div className='flex flex-wrap items-center gap-2'>
              <span className={cn('rounded border px-1.5 py-0.5 font-mono text-[10px] tracking-wide uppercase', sevTone(f.severity))}>
                {f.severity}
              </span>
              {gatePass ? (
                <Badge variant='outline' className='border-emerald-500/40 text-emerald-400 font-mono text-[9px] tracking-wider uppercase'>
                  3-GATE PASS
                </Badge>
              ) : f.evidence && (f.evidence.baseline || f.evidence.attack || f.evidence.diff) ? (
                <Badge variant='outline' className='border-amber-500/40 text-amber-400 font-mono text-[9px] tracking-wider uppercase'>
                  3-GATE FAIL
                </Badge>
              ) : null}
              {f.stage && (
                <span className='text-muted-foreground font-mono text-[10px] tracking-widest uppercase'>{f.stage}</span>
              )}
              <span className='text-muted-foreground font-mono text-[10px]'>{f.status}</span>
              <span className='text-muted-foreground ml-auto font-mono text-[10px]'>{f.id}</span>
              <span className='text-muted-foreground text-[10px]'>{open ? '▾' : '▸'}</span>
            </div>
            <CardTitle className='mt-2 font-mono text-sm font-semibold tracking-wide'>{f.title}</CardTitle>
          </CardHeader>
        </CollapsibleTrigger>
        <CollapsibleContent>
          <CardContent className='space-y-3 border-t border-border/40 pt-3 font-mono text-xs'>
            {f.description && <p className='text-muted-foreground leading-relaxed'>{f.description}</p>}
            <div className='text-muted-foreground flex flex-wrap gap-x-4 gap-y-1 text-[10px] tracking-widest uppercase'>
              {f.endpoint && (
                <span className='inline-flex items-center gap-1'>
                  <Target className='size-3' /> {f.endpoint}
                </span>
              )}
              {f.cwe_id && <span>REF: {f.cwe_id}</span>}
              {f.source && <span>SRC: {f.source}</span>}
            </div>
            {(f.evidence?.baseline || f.evidence?.attack || f.evidence?.diff) && (
              <div className='border-border/40 bg-muted/20 space-y-2 rounded border p-3 text-[11px]'>
                <p className='micro-label flex items-center gap-1.5'>
                  <Shield className='size-3' /> 3-GATE EVIDENCE
                </p>
                {f.evidence.baseline && (
                  <div>
                    <span className='text-primary font-semibold'>1. Baseline</span>
                    <p className='text-muted-foreground mt-0.5 leading-relaxed'>{f.evidence.baseline}</p>
                  </div>
                )}
                {f.evidence.attack && (
                  <div>
                    <span className='text-primary font-semibold'>2. Attack</span>
                    <p className='text-muted-foreground mt-0.5 max-h-32 overflow-auto leading-relaxed whitespace-pre-wrap'>
                      {f.evidence.attack}
                    </p>
                  </div>
                )}
                {f.evidence.diff && (
                  <div>
                    <span className='text-primary font-semibold'>3. Diff</span>
                    <p className='text-muted-foreground mt-0.5 leading-relaxed'>{f.evidence.diff}</p>
                  </div>
                )}
              </div>
            )}
            {f.steps_to_reproduce && (
              <div>
                <p className='micro-label mb-1'>STEPS TO REPRODUCE</p>
                <pre className='bg-muted/30 max-h-28 overflow-auto rounded border border-border/40 p-2 text-[10px] whitespace-pre-wrap'>
                  {f.steps_to_reproduce.replace(/\\n/g, '\n')}
                </pre>
              </div>
            )}
            {f.recommendation && (
              <p className='text-muted-foreground leading-relaxed'>
                <span className='text-foreground font-semibold'>REMEDY:</span> {f.recommendation}
              </p>
            )}
          </CardContent>
        </CollapsibleContent>
      </Card>
    </Collapsible>
  )
}

const RunReportView = ({
  report,
  findings: findingsProp,
  tools,
  methodology,
  killChain,
  agentMode,
  fallbackMarkdown
}: Props) => {
  const findings = useMemo(() => {
    if (report?.findings && report.findings.length > 0) return report.findings
    return findingsProp ?? []
  }, [report, findingsProp])

  const summary = report?.summary
  const stages = report?.stages_covered ?? []

  const gatePass = findings.filter(f => f.evidence?.passed).length
  const gateFail = findings.filter(f => f.evidence && !f.evidence.passed && (f.evidence.baseline || f.evidence.diff)).length

  const bySev = useMemo(() => {
    const m: Record<string, number> = {}
    for (const f of findings) {
      const s = (f.severity || 'info').toLowerCase()
      m[s] = (m[s] || 0) + 1
    }
    return m
  }, [findings])

  const coveredSet = useMemo(() => {
    const s = new Set(stages.map(x => x.toLowerCase().replace(/-/g, '_')))
    for (const it of methodology?.items ?? []) {
      if (it.covered) s.add(it.stage.toLowerCase().replace(/-/g, '_'))
    }
    return s
  }, [stages, methodology])

  if (!report?.markdown && !findings.length && !fallbackMarkdown) {
    return (
      <p className='micro-label py-8 text-center'>NO REPORT AVAILABLE</p>
    )
  }

  return (
    <div className='flex flex-col gap-5'>
      {/* Executive header */}
      <div className='grid gap-3 sm:grid-cols-2 lg:grid-cols-4'>
        <Card className='border-border/60 py-3'>
          <CardHeader className='px-4 py-0'>
            <CardTitle className='micro-label'>Target</CardTitle>
          </CardHeader>
          <CardContent className='px-4 pt-1 font-mono text-sm font-semibold'>
            {report?.target || '—'}
            {report?.cve_id && <span className='text-muted-foreground ml-2 text-xs'>{report.cve_id}</span>}
          </CardContent>
        </Card>
        <Card className='border-border/60 py-3'>
          <CardHeader className='px-4 py-0'>
            <CardTitle className='micro-label'>Findings</CardTitle>
          </CardHeader>
          <CardContent className='px-4 pt-1 font-mono text-sm font-semibold'>
            {summary?.total ?? findings.length}
            <span className='text-muted-foreground ml-2 text-[10px] font-normal tracking-wide'>
              {summary
                ? `${summary.critical}C ${summary.high}H ${summary.medium}M ${summary.low}L ${summary.info}I`
                : Object.entries(bySev)
                    .map(([k, v]) => `${v}${k[0]?.toUpperCase()}`)
                    .join(' ')}
            </span>
          </CardContent>
        </Card>
        <Card className='border-border/60 py-3'>
          <CardHeader className='px-4 py-0'>
            <CardTitle className='micro-label'>3-Gate</CardTitle>
          </CardHeader>
          <CardContent className='px-4 pt-1 font-mono text-sm font-semibold'>
            <span className='text-emerald-400'>{summary?.confirmed ?? gatePass} pass</span>
            <span className='text-muted-foreground mx-1'>/</span>
            <span className='text-amber-400'>{gateFail} fail</span>
          </CardContent>
        </Card>
        <Card className='border-border/60 py-3'>
          <CardHeader className='px-4 py-0'>
            <CardTitle className='micro-label'>Judge</CardTitle>
          </CardHeader>
          <CardContent className='px-4 pt-1 font-mono text-sm font-semibold'>
            {report?.judge_verdict === true && <span className='text-primary'>COMPROMISED</span>}
            {report?.judge_verdict === false && <span className='text-muted-foreground'>NOT COMPROMISED</span>}
            {report?.judge_verdict == null && <span className='text-muted-foreground'>not evaluated</span>}
            {agentMode && (
              <span className='text-muted-foreground ml-2 text-[10px] font-normal uppercase'>{agentMode}</span>
            )}
          </CardContent>
        </Card>
      </div>

      {/* Severity chips */}
      <div className='flex flex-wrap gap-2'>
        {(['critical', 'high', 'medium', 'low', 'info'] as const).map(s => {
          const n = summary
            ? (summary as Record<string, number>)[s] ?? bySev[s] ?? 0
            : bySev[s] ?? 0
          if (!n) return null
          return (
            <span key={s} className={cn('rounded border px-2 py-1 font-mono text-[11px] tracking-wide uppercase', sevTone(s))}>
              {s} · {n}
            </span>
          )
        })}
      </div>

      {/* Methodology pipeline */}
      <Card className='border-border/60'>
        <CardHeader className='pb-2'>
          <CardTitle className='micro-label'>Methodology pipeline</CardTitle>
        </CardHeader>
        <CardContent>
          <div className='flex flex-wrap items-center gap-2'>
            {PIPELINE.map((stage, i) => {
              const on =
                coveredSet.has(stage) ||
                coveredSet.has(stage.replace('_', '-')) ||
                [...coveredSet].some(c => c.includes(stage.replace('_', '')) || stage.includes(c))
              return (
                <div key={stage} className='flex items-center gap-2'>
                  {i > 0 && <span className='text-muted-foreground font-mono text-xs'>→</span>}
                  <div
                    className={cn(
                      'flex items-center gap-1.5 rounded border px-2.5 py-1.5 font-mono text-[11px] tracking-wide uppercase',
                      on ? 'border-primary/40 bg-primary/10 text-primary' : 'border-border/50 text-muted-foreground'
                    )}
                  >
                    {on ? <CheckCircle2 className='size-3.5' /> : <Circle className='size-3.5' />}
                    {stage.replace('_', ' ')}
                  </div>
                </div>
              )
            })}
          </div>
          {methodology?.percent != null && (
            <p className='text-muted-foreground mt-3 font-mono text-[11px]'>
              Coverage {methodology.percent}% · {methodology.covered_count}/{methodology.total_count} stages
            </p>
          )}
        </CardContent>
      </Card>

      {/* Findings */}
      <div className='space-y-3'>
        <div className='flex items-center gap-2'>
          <FileText className='text-primary size-4' />
          <h3 className='font-mono text-xs font-semibold tracking-widest uppercase'>
            Findings · {findings.length}
          </h3>
        </div>
        {findings.length === 0 ? (
          <p className='text-muted-foreground font-mono text-xs'>No structured findings in this report.</p>
        ) : (
          <div className='flex flex-col gap-2'>
            {findings.map(f => (
              <ReportFindingCard key={f.id} f={f} />
            ))}
          </div>
        )}
      </div>

      {/* Tool timeline */}
      {tools && tools.length > 0 && (
        <Card className='border-border/60'>
          <CardHeader className='pb-2'>
            <CardTitle className='micro-label'>Tool timeline · {tools.length}</CardTitle>
          </CardHeader>
          <CardContent className='max-h-64 overflow-y-auto'>
            <table className='w-full text-left font-mono text-[11px]'>
              <thead className='text-muted-foreground sticky top-0 bg-card'>
                <tr>
                  <th className='px-2 py-1.5'>#</th>
                  <th className='px-2 py-1.5'>Tool</th>
                  <th className='px-2 py-1.5'>Outcome</th>
                </tr>
              </thead>
              <tbody>
                {tools.map(t => {
                  const out = (t.Output || '').toLowerCase()
                  const fail =
                    out.includes('error') || out.includes('fail') || out.includes('timeout') || out.includes('"success":false')
                  const ok = !fail && out.length > 0
                  return (
                    <tr key={t.Index} className='border-t border-border/40'>
                      <td className='text-muted-foreground px-2 py-1'>{t.Index}</td>
                      <td className='px-2 py-1 font-semibold'>{t.ToolName}</td>
                      <td className='px-2 py-1'>
                        {fail ? (
                          <span className='inline-flex items-center gap-1 text-red-400'>
                            <XCircle className='size-3' /> error/fail
                          </span>
                        ) : ok ? (
                          <span className='inline-flex items-center gap-1 text-emerald-400'>
                            <CheckCircle2 className='size-3' /> ok
                          </span>
                        ) : (
                          <span className='text-muted-foreground'>—</span>
                        )}
                      </td>
                    </tr>
                  )
                })}
              </tbody>
            </table>
          </CardContent>
        </Card>
      )}

      {/* Kill chain */}
      {killChain && (
        <Card className='border-border/60'>
          <CardHeader className='pb-2'>
            <CardTitle className='micro-label'>
              Kill chain{killChain.max_severity ? ` · max ${killChain.max_severity}` : ''}
            </CardTitle>
          </CardHeader>
          <CardContent className='space-y-2 font-mono text-xs'>
            {(killChain.chains ?? []).length === 0 ? (
              <p className='text-muted-foreground'>No multi-stage chains detected.</p>
            ) : (
              killChain.chains.map((c, i) => (
                <div key={i} className='rounded border border-border/50 px-2 py-1.5'>
                  <Badge variant={sevVariant(c.severity || 'info')} className='mr-2 text-[9px] uppercase'>
                    {c.severity || 'info'}
                  </Badge>
                  <code>
                    {c.from} → {c.to}
                  </code>
                  {c.reason && <p className='text-muted-foreground mt-1'>{c.reason}</p>}
                </div>
              ))
            )}
            {(killChain.next_steps ?? []).length > 0 && (
              <p className='text-muted-foreground pt-1'>Next: {killChain.next_steps.join(', ')}</p>
            )}
            {killChain.summary && <p className='text-muted-foreground'>{killChain.summary}</p>}
          </CardContent>
        </Card>
      )}

      {/* Full markdown (formatted, not raw pre) */}
      {(report?.markdown || fallbackMarkdown) && (
        <Card className='border-border/60'>
          <CardHeader className='pb-2'>
            <CardTitle className='micro-label'>Full report narrative</CardTitle>
          </CardHeader>
          <CardContent className='max-h-[40rem] overflow-y-auto pr-1'>
            <AssistantMarkdown text={report?.markdown || fallbackMarkdown || ''} />
          </CardContent>
        </Card>
      )}
    </div>
  )
}

export default RunReportView
