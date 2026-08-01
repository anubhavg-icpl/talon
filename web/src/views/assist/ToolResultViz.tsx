'use client'

import type { ReactNode } from 'react'
import Link from 'next/link'
import {
  Activity,
  AlertTriangle,
  Bug,
  CheckCircle2,
  CircleDashed,
  History,
  Radio,
  ShieldAlert,
  Target,
  XCircle
} from 'lucide-react'

import { Badge } from '@/components/ui/badge'
import { cn } from '@/lib/utils'

type HealthRow = {
  name: string
  status: string
  detail: string
  latency_ms?: number
  endpoint?: string
}

type IntelEv = {
  at?: string
  run_id?: string
  target?: string
  kind?: string
  label?: string
  detail?: string
  severity?: string
}

type RunRow = {
  run_id?: string
  target?: string
  status?: string
  tool_calls?: number
  findings_count?: number
  agent_mode?: string
  started_at?: string
}

type FindingRow = {
  id?: string
  title?: string
  severity?: string
  status?: string
  gate_passed?: boolean
}

const sevTone = (sev?: string) => {
  const s = (sev || '').toLowerCase()
  if (s === 'critical') return 'bg-red-500/15 text-red-400 border-red-500/40'
  if (s === 'high') return 'bg-orange-500/15 text-orange-400 border-orange-500/40'
  if (s === 'medium') return 'bg-amber-500/15 text-amber-400 border-amber-500/40'
  if (s === 'low') return 'bg-red-300/20 text-red-300 border-red-300/40'
  return 'bg-muted text-muted-foreground border-border'
}

const statusIcon = (status: string) => {
  const s = status.toLowerCase()
  if (s === 'online') return <CheckCircle2 className='size-3.5 shrink-0 text-emerald-500' />
  if (s === 'offline') return <XCircle className='size-3.5 shrink-0 text-red-400' />
  return <CircleDashed className='text-muted-foreground size-3.5 shrink-0' />
}

const parseJSON = (raw: string): unknown | null => {
  try {
    return JSON.parse(raw)
  } catch {
    return null
  }
}

const HealthTable = ({ rows }: { rows: HealthRow[] }) => (
  <div className='overflow-hidden rounded border border-border/70'>
    <table className='w-full text-left font-mono text-[11px]'>
      <thead className='bg-muted/50 text-muted-foreground'>
        <tr>
          <th className='px-2 py-1.5 font-medium tracking-wider uppercase'>Service</th>
          <th className='px-2 py-1.5 font-medium tracking-wider uppercase'>Status</th>
          <th className='hidden px-2 py-1.5 font-medium tracking-wider uppercase sm:table-cell'>ms</th>
          <th className='px-2 py-1.5 font-medium tracking-wider uppercase'>Detail</th>
        </tr>
      </thead>
      <tbody>
        {rows.map(r => (
          <tr key={r.name} className='border-t border-border/50'>
            <td className='px-2 py-1.5 font-semibold'>{r.name}</td>
            <td className='px-2 py-1.5'>
              <span className='inline-flex items-center gap-1'>
                {statusIcon(r.status)}
                <span
                  className={cn(
                    'uppercase',
                    r.status === 'online' && 'text-emerald-500',
                    r.status === 'offline' && 'text-red-400'
                  )}
                >
                  {r.status}
                </span>
              </span>
            </td>
            <td className='text-muted-foreground hidden px-2 py-1.5 sm:table-cell'>{r.latency_ms ?? '—'}</td>
            <td className='text-muted-foreground max-w-[200px] truncate px-2 py-1.5' title={r.detail}>
              {r.detail}
            </td>
          </tr>
        ))}
      </tbody>
    </table>
  </div>
)

const IntelFeedViz = ({ events, count }: { events: IntelEv[]; count?: number }) => {
  const findings = events.filter(e => e.kind === 'finding')
  const bySev: Record<string, number> = {}
  for (const f of findings) {
    const s = (f.severity || 'info').toLowerCase()
    bySev[s] = (bySev[s] || 0) + 1
  }
  const order = ['critical', 'high', 'medium', 'low', 'info']

  return (
    <div className='space-y-3'>
      <div className='flex flex-wrap items-center gap-2'>
        <Radio className='text-primary size-3.5' />
        <span className='font-mono text-[10px] tracking-widest uppercase'>
          Intel · {count ?? events.length} events · {findings.length} findings
        </span>
        <div className='flex flex-wrap gap-1'>
          {order.map(s =>
            bySev[s] ? (
              <span
                key={s}
                className={cn('rounded border px-1.5 py-0.5 font-mono text-[9px] tracking-wide uppercase', sevTone(s))}
              >
                {s} {bySev[s]}
              </span>
            ) : null
          )}
        </div>
      </div>

      <div className='relative max-h-72 space-y-0 overflow-y-auto pl-3'>
        <div className='bg-border absolute top-1 bottom-1 left-[7px] w-px' />
        {events.slice(0, 24).map((ev, i) => {
          const isFinding = ev.kind === 'finding'
          return (
            <div key={`${ev.run_id}-${i}-${ev.label}`} className='relative pb-3 pl-5'>
              <div
                className={cn(
                  'absolute top-1.5 left-0 size-3.5 rounded-full border-2 border-background',
                  isFinding
                    ? (ev.severity || '').toLowerCase() === 'critical'
                      ? 'bg-red-500'
                      : (ev.severity || '').toLowerCase() === 'high'
                        ? 'bg-orange-500'
                        : (ev.severity || '').toLowerCase() === 'medium'
                          ? 'bg-amber-500'
                          : 'bg-primary'
                    : 'bg-muted-foreground/50'
                )}
              />
              <div className='bg-card/80 rounded-md border border-border/60 px-2.5 py-2'>
                <div className='flex flex-wrap items-center gap-1.5'>
                  <Badge
                    variant={isFinding ? 'destructive' : 'outline'}
                    className='h-5 font-mono text-[9px] tracking-wider uppercase'
                  >
                    {isFinding ? <Bug className='mr-0.5 size-2.5' /> : <History className='mr-0.5 size-2.5' />}
                    {ev.kind || 'event'}
                  </Badge>
                  {ev.severity && (
                    <span
                      className={cn(
                        'rounded border px-1.5 py-0.5 font-mono text-[9px] tracking-wide uppercase',
                        sevTone(ev.severity)
                      )}
                    >
                      {ev.severity}
                    </span>
                  )}
                  {ev.target && (
                    <span className='text-muted-foreground inline-flex items-center gap-0.5 font-mono text-[10px]'>
                      <Target className='size-2.5' />
                      {ev.target}
                    </span>
                  )}
                  {ev.run_id && (
                    <Link
                      href={`/runs/${ev.run_id}`}
                      className='text-primary ml-auto font-mono text-[10px] underline-offset-2 hover:underline'
                    >
                      {ev.run_id.slice(0, 8)}…
                    </Link>
                  )}
                </div>
                <p className='mt-1 text-[12px] leading-snug font-medium'>{ev.label}</p>
                {ev.detail && (
                  <p className='text-muted-foreground mt-0.5 line-clamp-2 text-[11px] leading-snug'>{ev.detail}</p>
                )}
                {ev.at && <p className='text-muted-foreground mt-1 font-mono text-[10px]'>{ev.at}</p>}
              </div>
            </div>
          )
        })}
        {events.length > 24 && (
          <p className='text-muted-foreground pl-5 font-mono text-[10px]'>+{events.length - 24} more…</p>
        )}
      </div>
    </div>
  )
}

const RunsListViz = ({ runs, total }: { runs: RunRow[]; total?: number }) => (
  <div className='space-y-2'>
    <div className='flex items-center gap-2 font-mono text-[10px] tracking-widest uppercase'>
      <Activity className='text-primary size-3.5' />
      Runs · {total ?? runs.length}
    </div>
    <div className='max-h-64 overflow-y-auto rounded border border-border/70'>
      <table className='w-full text-left font-mono text-[11px]'>
        <thead className='bg-muted/50 text-muted-foreground sticky top-0'>
          <tr>
            <th className='px-2 py-1.5'>Target</th>
            <th className='px-2 py-1.5'>Status</th>
            <th className='px-2 py-1.5'>Find</th>
            <th className='px-2 py-1.5'>Run</th>
          </tr>
        </thead>
        <tbody>
          {runs.map(r => (
            <tr key={r.run_id} className='border-t border-border/50'>
              <td className='px-2 py-1.5 font-semibold'>{r.target || '—'}</td>
              <td className='px-2 py-1.5 uppercase'>{r.status}</td>
              <td className='px-2 py-1.5'>{r.findings_count ?? 0}</td>
              <td className='px-2 py-1.5'>
                {r.run_id ? (
                  <Link href={`/runs/${r.run_id}`} className='text-primary underline-offset-2 hover:underline'>
                    {r.run_id.slice(0, 8)}
                  </Link>
                ) : (
                  '—'
                )}
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  </div>
)

const FindingsViz = ({ findings, runId }: { findings: FindingRow[]; runId?: string }) => (
  <div className='space-y-2'>
    <div className='flex items-center gap-2 font-mono text-[10px] tracking-widest uppercase'>
      <ShieldAlert className='text-primary size-3.5' />
      Findings · {findings.length}
      {runId && (
        <Link href={`/runs/${runId}`} className='text-primary ml-auto normal-case tracking-normal underline'>
          open run
        </Link>
      )}
    </div>
    <div className='max-h-64 space-y-1.5 overflow-y-auto'>
      {findings.map((f, i) => (
        <div key={f.id || i} className='rounded border border-border/60 px-2 py-1.5'>
          <div className='flex flex-wrap items-center gap-1.5'>
            <span className={cn('rounded border px-1.5 py-0.5 font-mono text-[9px] uppercase', sevTone(f.severity))}>
              {f.severity || 'info'}
            </span>
            {f.status && (
              <Badge variant='outline' className='h-5 text-[9px] uppercase'>
                {f.status}
              </Badge>
            )}
            {f.gate_passed != null && (
              <span className='text-muted-foreground font-mono text-[9px]'>
                gate {f.gate_passed ? 'pass' : 'fail'}
              </span>
            )}
          </div>
          <p className='mt-1 text-[12px] leading-snug'>{f.title}</p>
        </div>
      ))}
    </div>
  </div>
)

const SummaryCards = ({ data }: { data: Record<string, unknown> }) => {
  const keys = [
    ['total', 'Total'],
    ['active', 'Active'],
    ['completed', 'Done'],
    ['errored', 'Error'],
    ['awaiting_approval', 'HITL'],
    ['compromised', 'Compromised']
  ] as const
  return (
    <div className='grid grid-cols-3 gap-1.5 sm:grid-cols-6'>
      {keys.map(([k, label]) =>
        data[k] != null ? (
          <div key={k} className='rounded border border-border/60 bg-muted/20 px-2 py-1.5 text-center'>
            <div className='text-foreground font-mono text-sm font-semibold'>{String(data[k])}</div>
            <div className='text-muted-foreground font-mono text-[9px] tracking-wider uppercase'>{label}</div>
          </div>
        ) : null
      )}
    </div>
  )
}

const SkillHits = ({
  hits,
  total
}: {
  hits: { ID?: string; id?: string; Name?: string; name?: string; Category?: string; category?: string; Stage?: string; stage?: string }[]
  total?: number
}) => (
  <div className='max-h-56 space-y-1 overflow-y-auto'>
    <p className='text-muted-foreground mb-1 font-mono text-[10px] tracking-wider uppercase'>
      Skills · {total ?? hits.length} hits
    </p>
    {hits.map((h, i) => {
      const id = h.id || h.ID || String(i)
      const name = h.name || h.Name || id
      const cat = h.category || h.Category
      const stage = h.stage || h.Stage
      return (
        <div key={id} className='rounded border border-border/50 px-2 py-1.5'>
          <p className='text-[12px] font-medium'>{name}</p>
          <p className='text-muted-foreground font-mono text-[10px]'>
            {id}
            {cat ? ` · ${cat}` : ''}
            {stage ? ` · ${stage}` : ''}
          </p>
        </div>
      )
    })}
  </div>
)

/**
 * Renders structured tool JSON as operator-friendly visuals.
 * Falls back to null so the parent can show raw pre.
 */
export function ToolResultViz({ toolName, content }: { toolName?: string; content: string }) {
  if (!content?.trim()) return null
  const data = parseJSON(content)
  if (!data || typeof data !== 'object') return null
  const obj = data as Record<string, unknown>
  const name = toolName || ''

  // service_health
  if (name === 'service_health' || Array.isArray(obj.services)) {
    const services = obj.services as HealthRow[] | undefined
    if (Array.isArray(services) && services.length && services[0]?.name && services[0]?.status != null) {
      return <HealthTable rows={services} />
    }
  }

  // intel_feed
  if (name === 'intel_feed' || Array.isArray(obj.events)) {
    const events = obj.events as IntelEv[] | undefined
    if (Array.isArray(events) && events.length && (events[0]?.label != null || events[0]?.kind != null)) {
      return <IntelFeedViz events={events} count={typeof obj.count === 'number' ? obj.count : events.length} />
    }
  }

  // list_runs
  if (name === 'list_runs' || Array.isArray(obj.runs)) {
    const runs = obj.runs as RunRow[] | undefined
    if (Array.isArray(runs) && runs.length && (runs[0]?.run_id || runs[0]?.target)) {
      return <RunsListViz runs={runs} total={typeof obj.total === 'number' ? obj.total : undefined} />
    }
  }

  // get_findings
  if (name === 'get_findings' || Array.isArray(obj.findings)) {
    const findings = obj.findings as FindingRow[] | undefined
    if (Array.isArray(findings) && findings.length && findings[0]?.title) {
      return <FindingsViz findings={findings} runId={typeof obj.run_id === 'string' ? obj.run_id : undefined} />
    }
  }

  // runs_summary
  if (name === 'runs_summary' || (obj.total != null && obj.completed != null)) {
    return <SummaryCards data={obj} />
  }

  // search_skills
  if (name === 'search_skills' || Array.isArray(obj.hits)) {
    const hits = obj.hits as Parameters<typeof SkillHits>[0]['hits']
    if (Array.isArray(hits) && hits.length) {
      return <SkillHits hits={hits} total={typeof obj.total === 'number' ? obj.total : undefined} />
    }
  }

  // list_agents
  if (name === 'list_agents' && Array.isArray(obj.agents)) {
    const agents = obj.agents as { id?: string; name?: string; codename?: string; focus?: string }[]
    return (
      <div className='max-h-56 space-y-1 overflow-y-auto'>
        {agents.map(a => (
          <div key={a.id || a.codename} className='rounded border border-border/50 px-2 py-1.5'>
            <p className='text-[12px] font-medium'>
              {a.name || a.id} <span className='text-muted-foreground font-mono text-[10px]'>{a.codename}</span>
            </p>
            {a.focus && <p className='text-muted-foreground text-[11px]'>{a.focus}</p>}
          </div>
        ))}
      </div>
    )
  }

  // list_playbooks
  if (name === 'list_playbooks' && Array.isArray(obj.playbooks)) {
    const pbs = obj.playbooks as { id?: string; name?: string; agent_mode?: string; description?: string }[]
    return (
      <div className='max-h-56 space-y-1 overflow-y-auto'>
        {pbs.map(p => (
          <div key={p.id} className='rounded border border-border/50 px-2 py-1.5'>
            <p className='text-[12px] font-medium'>{p.name || p.id}</p>
            <p className='text-muted-foreground font-mono text-[10px]'>{p.agent_mode}</p>
            {p.description && <p className='text-muted-foreground line-clamp-2 text-[11px]'>{p.description}</p>}
          </div>
        ))}
      </div>
    )
  }

  // list_targets
  if (name === 'list_targets' && Array.isArray(obj.targets)) {
    const targets = obj.targets as { id?: string; address?: string; url?: string; status?: string; name?: string }[]
    if (targets.length === 0) {
      return <p className='text-muted-foreground font-mono text-[11px]'>No targets registered</p>
    }
    return (
      <div className='max-h-48 space-y-1 overflow-y-auto'>
        {targets.map((t, i) => (
          <div key={t.id || i} className='flex items-center gap-2 rounded border border-border/50 px-2 py-1.5 font-mono text-[11px]'>
            <Target className='text-primary size-3' />
            <span>{t.address || t.url || t.name || t.id}</span>
            {t.status && <Badge variant='outline' className='ml-auto h-5 text-[9px]'>{t.status}</Badge>}
          </div>
        ))}
      </div>
    )
  }

  // get_run_status compact
  if (name === 'get_run_status' && obj.run_id) {
    return (
      <div className='space-y-1.5 rounded border border-border/60 px-2.5 py-2 font-mono text-[11px]'>
        <div className='flex flex-wrap gap-2'>
          <Badge variant='outline' className='uppercase'>{String(obj.status)}</Badge>
          {obj.target != null && <span className='text-muted-foreground'>{String(obj.target)}</span>}
          {obj.agent_mode != null && <span className='text-muted-foreground'>mode={String(obj.agent_mode)}</span>}
        </div>
        <Link href={`/runs/${obj.run_id}`} className='text-primary text-[10px] underline'>
          {String(obj.run_id)}
        </Link>
        {obj.findings_count != null && (
          <p className='text-muted-foreground'>findings: {String(obj.findings_count)} · tools: {String(obj.tool_calls ?? 0)}</p>
        )}
        {typeof obj.output_preview === 'string' && obj.output_preview && (
          <pre className='bg-muted/30 max-h-24 overflow-auto rounded p-1.5 text-[10px] whitespace-pre-wrap'>
            {obj.output_preview}
          </pre>
        )}
      </div>
    )
  }

  // Unknown structured object — compact key chips, not a wall of JSON
  if (name) {
    return (
      <div className='flex flex-wrap gap-1.5'>
        <AlertTriangle className='text-muted-foreground size-3.5' />
        {Object.entries(obj)
          .slice(0, 12)
          .map(([k, v]) => (
            <span key={k} className='bg-muted/40 rounded border border-border/50 px-1.5 py-0.5 font-mono text-[10px]'>
              <span className='text-muted-foreground'>{k}:</span>{' '}
              {typeof v === 'object' ? (Array.isArray(v) ? `[${v.length}]` : '{…}') : String(v).slice(0, 48)}
            </span>
          ))}
      </div>
    )
  }

  return null
}

/** Lightweight markdown-ish render for assistant replies (tables + bold + code). */
export function AssistantMarkdown({ text }: { text: string }) {
  if (!text) return null

  const blocks = text.split(/\n\n+/)
  return (
    <div className='space-y-2.5 text-sm leading-relaxed'>
      {blocks.map((block, bi) => {
        const lines = block.split('\n')
        // Markdown table
        if (lines.length >= 2 && lines[0].includes('|') && lines[1]?.match(/^\|?[\s:-]+\|/)) {
          const rows = lines
            .filter(l => l.includes('|') && !l.match(/^\|?[\s:-]+\|/))
            .map(l =>
              l
                .trim()
                .replace(/^\|/, '')
                .replace(/\|$/, '')
                .split('|')
                .map(c => c.trim())
            )
          if (rows.length === 0) return null
          return (
            <div key={bi} className='overflow-x-auto rounded border border-border/60'>
              <table className='w-full text-left font-mono text-[11px]'>
                <thead className='bg-muted/40'>
                  <tr>
                    {rows[0].map((c, i) => (
                      <th key={i} className='px-2 py-1.5 font-medium'>
                        {inlineFmt(c)}
                      </th>
                    ))}
                  </tr>
                </thead>
                <tbody>
                  {rows.slice(1).map((row, ri) => (
                    <tr key={ri} className='border-t border-border/40'>
                      {row.map((c, i) => (
                        <td key={i} className='px-2 py-1.5 align-top'>
                          {inlineFmt(c)}
                        </td>
                      ))}
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          )
        }

        // Heading
        if (lines[0].startsWith('### ')) {
          return (
            <h4 key={bi} className='font-mono text-xs font-semibold tracking-wide'>
              {inlineFmt(lines[0].slice(4))}
            </h4>
          )
        }
        if (lines[0].startsWith('## ')) {
          return (
            <h3 key={bi} className='font-mono text-sm font-semibold tracking-wide'>
              {inlineFmt(lines[0].slice(3))}
            </h3>
          )
        }

        // List
        if (lines.every(l => !l.trim() || /^[-*•]|\d+\./.test(l.trim()))) {
          return (
            <ul key={bi} className='list-inside list-disc space-y-0.5 pl-0.5 text-[13px]'>
              {lines
                .filter(l => l.trim())
                .map((l, i) => (
                  <li key={i}>{inlineFmt(l.replace(/^[-*•]\s*|^\d+\.\s*/, ''))}</li>
                ))}
            </ul>
          )
        }

        return (
          <p key={bi} className='text-[13px] whitespace-pre-wrap'>
            {lines.map((l, i) => (
              <span key={i}>
                {i > 0 && <br />}
                {inlineFmt(l)}
              </span>
            ))}
          </p>
        )
      })}
    </div>
  )
}

function inlineFmt(s: string): ReactNode {
  // **bold**, `code`
  const parts: React.ReactNode[] = []
  const re = /(\*\*[^*]+\*\*|`[^`]+`)/g
  let last = 0
  let m: RegExpExecArray | null
  let k = 0
  while ((m = re.exec(s))) {
    if (m.index > last) parts.push(s.slice(last, m.index))
    const tok = m[0]
    if (tok.startsWith('**')) {
      parts.push(
        <strong key={k++} className='font-semibold'>
          {tok.slice(2, -2)}
        </strong>
      )
    } else {
      parts.push(
        <code key={k++} className='bg-muted/60 rounded px-1 py-0.5 font-mono text-[11px]'>
          {tok.slice(1, -1)}
        </code>
      )
    }
    last = m.index + tok.length
  }
  if (last < s.length) parts.push(s.slice(last))
  return parts.length === 1 ? parts[0] : <>{parts}</>
}
