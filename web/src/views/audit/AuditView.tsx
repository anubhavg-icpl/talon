'use client'

import { useEffect, useState } from 'react'

import { toast } from 'sonner'

import type { AuditEntry, AuditSeverity, AuditStats } from '@/lib/api'
import { auditStats, exportAudit, listAudit } from '@/lib/api'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { Card, CardContent } from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import { Skeleton } from '@/components/ui/skeleton'

const sevVariant = (
  sev: AuditSeverity | string
): 'destructive' | 'default' | 'secondary' | 'outline' => {
  switch (sev) {
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

const sevColor = (sev: string): string => {
  switch (sev) {
    case 'critical':
      return 'bg-red-500'
    case 'high':
      return 'bg-orange-500'
    case 'medium':
      return 'bg-yellow-500'
    case 'low':
      return 'bg-blue-500'
    default:
      return 'bg-zinc-500'
  }
}

const STAT_KEYS: Array<{ key: keyof Omit<AuditStats, 'total'>; label: string }> = [
  { key: 'critical', label: 'CRITICAL' },
  { key: 'high', label: 'HIGH' },
  { key: 'medium', label: 'MEDIUM' },
  { key: 'low', label: 'LOW' },
  { key: 'info', label: 'INFO' }
]

const AuditView = () => {
  const [runId, setRunId] = useState('')
  const [committedRunId, setCommittedRunId] = useState('')
  const [entries, setEntries] = useState<AuditEntry[] | null>(null)
  const [stats, setStats] = useState<AuditStats | null>(null)
  const [error, setError] = useState<string | null>(null)
  const [exporting, setExporting] = useState(false)

  useEffect(() => {
    if (!committedRunId) return
    let stale = false

    const load = async () => {
      try {
        const [ent, st] = await Promise.all([
          listAudit(committedRunId),
          auditStats(committedRunId).catch(() => null)
        ])
        if (stale) return
        setEntries(ent ?? [])
        if (st) setStats(st)
        setError(null)
      } catch (err) {
        if (stale) return
        setError(err instanceof Error ? err.message : String(err))
        setEntries([])
      }
    }

    void load()
    return () => {
      stale = true
    }
  }, [committedRunId])

  const submitRunId = (e: React.FormEvent) => {
    e.preventDefault()
    const rid = runId.trim()
    if (!rid) return
    setEntries(null)
    setStats(null)
    setCommittedRunId(rid)
  }

  const onExport = async () => {
    if (!committedRunId) return
    setExporting(true)
    try {
      const data = await exportAudit(committedRunId)
      const blob = new Blob([JSON.stringify(data, null, 2)], { type: 'application/json' })
      const url = URL.createObjectURL(blob)
      const a = document.createElement('a')
      a.href = url
      a.download = `audit-${committedRunId.slice(0, 8)}.json`
      a.click()
      URL.revokeObjectURL(url)
      toast.success(`Exported ${data.count} entries`)
    } catch (err) {
      toast.error(err instanceof Error ? err.message : 'export failed')
    } finally {
      setExporting(false)
    }
  }

  return (
    <div className='flex flex-col gap-6'>
      <div className='flex flex-wrap items-end justify-between gap-4'>
        <div>
          <h1 className='font-mono text-xl font-semibold tracking-widest'>AUDIT TRAIL</h1>
          <p className='micro-label mt-1'>TAMPER-EVIDENT COMPLIANCE LOG — SEVERITY-COLORED TIMELINE</p>
        </div>
        {committedRunId && (
          <Button
            size='sm'
            variant='outline'
            disabled={exporting}
            onClick={onExport}
            className='font-mono text-[10px] tracking-widest uppercase'
          >
            {exporting ? 'Exporting…' : 'Export JSON'}
          </Button>
        )}
      </div>

      <Card>
        <CardContent>
          <form onSubmit={submitRunId} className='flex flex-wrap gap-2'>
            <Input
              placeholder='run_id'
              className='font-mono text-xs'
              value={runId}
              onChange={e => setRunId(e.target.value)}
            />
            <Button type='submit' size='sm' className='font-mono text-[10px] tracking-widest uppercase'>
              Load trail
            </Button>
          </form>
        </CardContent>
      </Card>

      {error && (
        <div className='border-destructive/40 bg-destructive/10 text-destructive rounded-md border px-4 py-3 font-mono text-xs'>
          {error}
        </div>
      )}

      {stats && (
        <div className='grid grid-cols-2 gap-3 md:grid-cols-6'>
          {STAT_KEYS.map(({ key, label }) => (
            <Card key={String(key)} className='py-3'>
              <CardContent className='flex items-center gap-2 px-3 py-0'>
                <span className={`size-2.5 shrink-0 rounded-full ${sevColor(key)}`} />
                <div>
                  <div className='micro-label'>{label}</div>
                  <div className='font-mono text-2xl'>{stats[key]}</div>
                </div>
              </CardContent>
            </Card>
          ))}
          <Card className='py-3'>
            <CardContent className='px-3 py-0'>
              <div className='micro-label'>TOTAL</div>
              <div className='font-mono text-2xl'>{stats.total}</div>
            </CardContent>
          </Card>
        </div>
      )}

      {!committedRunId ? (
        <p className='micro-label py-12 text-center'>SELECT A RUN TO VIEW ITS AUDIT TRAIL</p>
      ) : !entries ? (
        <Skeleton className='h-40 w-full' />
      ) : entries.length === 0 ? (
        <p className='micro-label py-12 text-center'>NO AUDIT ENTRIES</p>
      ) : (
        <div className='relative flex flex-col gap-0 border-l border-border/40 pl-6'>
          {entries.map((entry, i) => (
            <div key={entry.id} className='relative pb-4'>
              <span
                className={`absolute top-1 -left-[27px] size-3 rounded-full ring-2 ring-background ${sevColor(entry.severity)}`}
              />
              <div className='flex flex-wrap items-center gap-2'>
                <Badge variant={sevVariant(entry.severity)} className='font-mono text-[9px] uppercase'>
                  {entry.severity}
                </Badge>
                <Badge variant='outline' className='font-mono text-[9px] uppercase'>
                  {entry.actor}
                </Badge>
                <span className='font-mono text-xs font-semibold'>{entry.action}</span>
                {i === 0 && (
                  <span className='text-muted-foreground ml-auto font-mono text-[9px]'>#{entries.length - i}</span>
                )}
              </div>
              <div className='mt-1 font-mono text-[11px]'>
                <span className='text-muted-foreground'>{entry.resource_type}</span>
                {entry.resource_id && (
                  <span className='text-muted-foreground'> · {entry.resource_id}</span>
                )}
                <span className='text-muted-foreground'> · {new Date(entry.timestamp).toLocaleString()}</span>
                {entry.ip_address && (
                  <span className='text-muted-foreground'> · {entry.ip_address}</span>
                )}
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  )
}

export default AuditView
