'use client'

import { useMemo, useState } from 'react'
import type { LucideIcon } from 'lucide-react'
import {
  AlertTriangle,
  Bug,
  CheckCircle2,
  ChevronDown,
  CircleDot,
  Crosshair,
  FileText,
  Radar,
  Shield,
  Terminal,
  Wrench,
  XCircle
} from 'lucide-react'

import { Badge } from '@/components/ui/badge'
import { Collapsible, CollapsibleContent, CollapsibleTrigger } from '@/components/ui/collapsible'
import type { Finding, TimelineEvent, ToolCallRecord } from '@/lib/api'
import { cn } from '@/lib/utils'

const STAGE_ORDER = ['recon', 'exploit', 'post_exploit', 'codegen', 'report', 'other'] as const

const stageMeta: Record<string, { label: string; icon: LucideIcon; tone: string }> = {
  recon: { label: 'Recon', icon: Radar, tone: 'text-red-400 border-red-500/40 bg-red-500/10' },
  exploit: { label: 'Exploit', icon: Crosshair, tone: 'text-orange-400 border-orange-500/40 bg-orange-500/10' },
  post_exploit: { label: 'Post-exploit', icon: Shield, tone: 'text-violet-400 border-violet-500/40 bg-violet-500/10' },
  codegen: { label: 'Codegen', icon: Terminal, tone: 'text-emerald-400 border-emerald-500/40 bg-emerald-500/10' },
  report: { label: 'Report', icon: FileText, tone: 'text-primary border-primary/40 bg-primary/10' },
  other: { label: 'Other', icon: Wrench, tone: 'text-muted-foreground border-border bg-muted/30' }
}

const sevTone = (sev?: string) => {
  const s = (sev || '').toLowerCase()
  if (s === 'critical') return 'bg-red-500/15 text-red-400 border-red-500/40'
  if (s === 'high') return 'bg-orange-500/15 text-orange-400 border-orange-500/40'
  if (s === 'medium') return 'bg-amber-500/15 text-amber-400 border-amber-500/40'
  if (s === 'low') return 'bg-red-300/20 text-red-300 border-red-300/40'
  return 'bg-muted text-muted-foreground border-border'
}

const stripAnsi = (s: string) => s.replace(/\x1b\[[0-9;]*m/g, '')

const formatToolOutput = (raw: string): string => {
  if (!raw) return ''
  try {
    const parsed: unknown = JSON.parse(raw)
    if (parsed && typeof parsed === 'object' && !Array.isArray(parsed)) {
      const obj = parsed as Record<string, unknown>
      if (typeof obj.output === 'string' && obj.output.trim()) return stripAnsi(obj.output)
      const stdout = typeof obj.stdout === 'string' ? stripAnsi(obj.stdout) : ''
      const stderr = typeof obj.stderr === 'string' ? stripAnsi(obj.stderr) : ''
      if (stdout.trim()) {
        return obj.success === false && stderr.trim() ? `${stdout}\n[stderr]\n${stderr}` : stdout
      }
      if (stderr.trim()) return stderr
      if (typeof obj.message === 'string' && obj.message.trim()) return obj.message
      if (typeof obj.error === 'string' && obj.error.trim()) return obj.error
      return JSON.stringify(parsed, null, 2)
    }
  } catch {
    /* plain */
  }
  return stripAnsi(raw)
}

const toolOutcome = (output: string): 'ok' | 'fail' | 'warn' | 'empty' => {
  if (!output?.trim()) return 'empty'
  const low = output.toLowerCase()
  if (
    low.includes('"success":false') ||
    low.includes('error') ||
    low.includes('failed') ||
    low.includes('not found') ||
    low.includes('invalid module') ||
    low.includes('return_code":127') ||
    low.includes('http 400')
  ) {
    return 'fail'
  }
  if (low.includes('warning') || low.includes('no session') || low.includes('timeout')) return 'warn'
  return 'ok'
}

const stageForTool = (name: string): string => {
  const n = name.toLowerCase()
  if (n.includes('nmap') || n.includes('nuclei') || n.includes('rustscan') || n.includes('recon') || n.includes('enum'))
    return 'recon'
  if (n.includes('exploit') || n.includes('sqlmap') || n.includes('auxiliary') || n.includes('msf')) return 'exploit'
  if (n.includes('session') || n.includes('post') || n.includes('shell')) return 'post_exploit'
  if (n.includes('codegen') || n.includes('custom_exploit')) return 'codegen'
  if (n.includes('report') || n.includes('judge')) return 'report'
  if (n.includes('skill')) return 'other'
  return 'other'
}

type TimelineItem = {
  key: string
  index: number
  kind: 'tool' | 'finding'
  stage: string
  label: string
  detail: string
  severity?: string
  outcome?: 'ok' | 'fail' | 'warn' | 'empty'
  args?: string
  finding?: Finding
}

type Props = {
  timeline?: TimelineEvent[]
  tools?: ToolCallRecord[]
  findings?: Finding[]
}

const RunTimelineView = ({ timeline = [], tools = [], findings = [] }: Props) => {
  const items = useMemo((): TimelineItem[] => {
    // Prefer full tool log + findings over truncated API timeline detail.
    if (tools.length > 0 || findings.length > 0) {
      const out: TimelineItem[] = []
      for (const t of [...tools].sort((a, b) => a.Index - b.Index)) {
        const args =
          typeof t.Args === 'string' ? t.Args : t.Args ? JSON.stringify(t.Args) : ''
        out.push({
          key: `tool-${t.Index}`,
          index: t.Index,
          kind: 'tool',
          stage: stageForTool(t.ToolName),
          label: t.ToolName,
          detail: t.Output || '',
          outcome: toolOutcome(t.Output || ''),
          args: args.length > 200 ? args.slice(0, 200) + '…' : args
        })
      }
      const base = tools.length
      findings.forEach((f, i) => {
        out.push({
          key: `finding-${f.id}`,
          index: base + i,
          kind: 'finding',
          stage: (f.stage || 'other').toLowerCase().replace(/-/g, '_'),
          label: f.title,
          detail: f.description || '',
          severity: f.severity,
          finding: f
        })
      })
      return out
    }
    return timeline.map(ev => ({
      key: `${ev.kind}-${ev.index}`,
      index: ev.index,
      kind: (ev.kind === 'finding' ? 'finding' : 'tool') as 'tool' | 'finding',
      stage: (ev.stage || 'other').toLowerCase().replace(/-/g, '_'),
      label: ev.label,
      detail: ev.detail || '',
      severity: ev.severity,
      outcome: ev.kind === 'tool' ? toolOutcome(ev.detail || '') : undefined
    }))
  }, [timeline, tools, findings])

  const stats = useMemo(() => {
    const toolsN = items.filter(i => i.kind === 'tool').length
    const findingsN = items.filter(i => i.kind === 'finding').length
    const fails = items.filter(i => i.outcome === 'fail').length
    const oks = items.filter(i => i.outcome === 'ok').length
    return { toolsN, findingsN, fails, oks, total: items.length }
  }, [items])

  const grouped = useMemo(() => {
    const map = new Map<string, TimelineItem[]>()
    for (const it of items) {
      const st = STAGE_ORDER.includes(it.stage as (typeof STAGE_ORDER)[number]) ? it.stage : 'other'
      if (!map.has(st)) map.set(st, [])
      map.get(st)!.push(it)
    }
    return STAGE_ORDER.filter(s => map.has(s)).map(s => ({ stage: s, items: map.get(s)! }))
  }, [items])

  if (items.length === 0) {
    return <p className='micro-label py-8 text-center'>NO TIMELINE EVENTS YET</p>
  }

  return (
    <div className='flex flex-col gap-4'>
      {/* Summary strip */}
      <div className='flex flex-wrap gap-2 font-mono text-[10px] tracking-widest uppercase'>
        <Badge variant='outline'>Events {stats.total}</Badge>
        <Badge variant='outline'>Tools {stats.toolsN}</Badge>
        <Badge variant='outline'>Findings {stats.findingsN}</Badge>
        {stats.oks > 0 && (
          <Badge variant='outline' className='border-emerald-500/40 text-emerald-400'>
            Ok {stats.oks}
          </Badge>
        )}
        {stats.fails > 0 && (
          <Badge variant='outline' className='border-red-500/40 text-red-400'>
            Fail {stats.fails}
          </Badge>
        )}
      </div>

      {/* Vertical timeline by stage */}
      <div className='space-y-6'>
        {grouped.map(({ stage, items: stageItems }) => {
          const meta = stageMeta[stage] || stageMeta.other
          const Icon = meta.icon
          return (
            <div key={stage}>
              <div className={cn('mb-3 inline-flex items-center gap-2 rounded border px-2.5 py-1 font-mono text-[11px] tracking-wide uppercase', meta.tone)}>
                <Icon className='size-3.5' />
                {meta.label}
                <span className='opacity-70'>· {stageItems.length}</span>
              </div>

              <div className='relative pl-5'>
                <div className='bg-border absolute top-1 bottom-1 left-[9px] w-px' />
                <div className='space-y-2'>
                  {stageItems.map(it => (
                    <TimelineRow key={it.key} item={it} />
                  ))}
                </div>
              </div>
            </div>
          )
        })}
      </div>
    </div>
  )
}

const TimelineRow = ({ item }: { item: TimelineItem }) => {
  const [open, setOpen] = useState(false)
  const isFinding = item.kind === 'finding'
  const pretty = item.kind === 'tool' ? formatToolOutput(item.detail) : item.detail

  const nodeColor =
    isFinding
      ? (item.severity || '').toLowerCase() === 'critical'
        ? 'bg-red-500'
        : (item.severity || '').toLowerCase() === 'high'
          ? 'bg-orange-500'
          : (item.severity || '').toLowerCase() === 'medium'
            ? 'bg-amber-500'
            : 'bg-primary'
      : item.outcome === 'fail'
        ? 'bg-red-500'
        : item.outcome === 'warn'
          ? 'bg-amber-500'
          : item.outcome === 'ok'
            ? 'bg-emerald-500'
            : 'bg-muted-foreground/50'

  return (
    <Collapsible open={open} onOpenChange={setOpen}>
      <div className='relative pb-1 pl-4'>
        <div className={cn('absolute top-2 left-0 size-3.5 rounded-full border-2 border-background', nodeColor)} />

        <div
          className={cn(
            'rounded-md border px-3 py-2',
            isFinding ? 'border-border/70 bg-card/80' : 'border-border/50 bg-muted/15'
          )}
        >
          <CollapsibleTrigger className='flex w-full items-start gap-2 text-left'>
            <div className='mt-0.5 shrink-0'>
              {isFinding ? (
                <Bug className='text-primary size-3.5' />
              ) : item.outcome === 'fail' ? (
                <XCircle className='size-3.5 text-red-400' />
              ) : item.outcome === 'ok' ? (
                <CheckCircle2 className='size-3.5 text-emerald-400' />
              ) : item.outcome === 'warn' ? (
                <AlertTriangle className='size-3.5 text-amber-400' />
              ) : (
                <CircleDot className='text-muted-foreground size-3.5' />
              )}
            </div>
            <div className='min-w-0 flex-1'>
              <div className='flex flex-wrap items-center gap-1.5'>
                <Badge
                  variant={isFinding ? 'destructive' : 'outline'}
                  className='h-5 font-mono text-[9px] tracking-wider uppercase'
                >
                  {item.kind}
                </Badge>
                {isFinding && item.severity && (
                  <span className={cn('rounded border px-1.5 py-0.5 font-mono text-[9px] uppercase', sevTone(item.severity))}>
                    {item.severity}
                  </span>
                )}
                {!isFinding && item.outcome && item.outcome !== 'empty' && (
                  <span
                    className={cn(
                      'font-mono text-[9px] tracking-wider uppercase',
                      item.outcome === 'ok' && 'text-emerald-400',
                      item.outcome === 'fail' && 'text-red-400',
                      item.outcome === 'warn' && 'text-amber-400'
                    )}
                  >
                    {item.outcome}
                  </span>
                )}
                <span className='text-muted-foreground ml-auto font-mono text-[10px]'>#{item.index}</span>
                <ChevronDown className={cn('text-muted-foreground size-3.5 transition-transform', open && 'rotate-180')} />
              </div>
              <p className='mt-1 font-mono text-[12px] font-semibold tracking-wide'>{item.label}</p>
              {!open && pretty && (
                <p className='text-muted-foreground mt-0.5 line-clamp-2 font-mono text-[11px]'>{pretty.slice(0, 180)}</p>
              )}
              {!open && !isFinding && item.args && (
                <p className='text-muted-foreground/80 mt-0.5 line-clamp-1 font-mono text-[10px]'>{item.args}</p>
              )}
            </div>
          </CollapsibleTrigger>

          <CollapsibleContent className='mt-2 border-t border-border/40 pt-2'>
            {isFinding && item.finding ? (
              <div className='space-y-2 font-mono text-[11px]'>
                {item.finding.endpoint && (
                  <p className='text-muted-foreground'>
                    <span className='text-foreground'>Endpoint:</span> {item.finding.endpoint}
                  </p>
                )}
                <p className='text-muted-foreground leading-relaxed'>{item.finding.description}</p>
                {item.finding.evidence && (
                  <div className='bg-muted/20 space-y-1 rounded border border-border/40 p-2'>
                    <p className='micro-label'>3-GATE</p>
                    {item.finding.evidence.baseline && <p>1. {item.finding.evidence.baseline}</p>}
                    {item.finding.evidence.attack && (
                      <p className='max-h-20 overflow-auto'>2. {item.finding.evidence.attack}</p>
                    )}
                    {item.finding.evidence.diff && <p>3. {item.finding.evidence.diff}</p>}
                    <p className={item.finding.evidence.passed ? 'text-emerald-400' : 'text-amber-400'}>
                      {item.finding.evidence.passed ? 'PASS' : 'FAIL / incomplete'}
                    </p>
                  </div>
                )}
                {item.finding.recommendation && (
                  <p className='text-muted-foreground'>
                    <span className='text-foreground'>Remedy:</span> {item.finding.recommendation}
                  </p>
                )}
              </div>
            ) : (
              <div className='space-y-2'>
                {item.args && (
                  <div>
                    <p className='micro-label mb-1'>ARGS</p>
                    <pre className='bg-muted/30 max-h-24 overflow-auto rounded border border-border/40 p-2 font-mono text-[10px] whitespace-pre-wrap'>
                      {item.args}
                    </pre>
                  </div>
                )}
                <div>
                  <p className='micro-label mb-1'>OUTPUT</p>
                  <pre className='bg-black/40 max-h-56 overflow-auto rounded border border-border/40 p-2 font-mono text-[10px] leading-relaxed whitespace-pre-wrap text-zinc-300'>
                    {pretty || '(empty)'}
                  </pre>
                </div>
              </div>
            )}
          </CollapsibleContent>
        </div>
      </div>
    </Collapsible>
  )
}

export default RunTimelineView
