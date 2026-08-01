'use client'

import { useEffect, useState } from 'react'

import Link from 'next/link'

import type { GlobalFinding } from '@/lib/api'
import { Badge } from '@/components/ui/badge'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Skeleton } from '@/components/ui/skeleton'
import { getGlobalFindings } from '@/lib/api'

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

const FindingsView = () => {
  const [items, setItems] = useState<GlobalFinding[] | null>(null)
  const [filter, setFilter] = useState('')
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    getGlobalFindings({ severity: filter || undefined, limit: 100 })
      .then(res => {
        setItems(res.findings ?? [])
        setError(null)
      })
      .catch(err => setError(err instanceof Error ? err.message : String(err)))
  }, [filter])

  return (
    <div className='flex flex-col gap-6'>
      <div>
        <h1 className='font-mono text-xl font-semibold tracking-widest'>FINDINGS</h1>
        <p className='micro-label mt-1'>GLOBAL REGISTRY — 3-GATE STRUCTURED FINDINGS ACROSS RUNS</p>
      </div>

      <div className='flex flex-wrap gap-2'>
        {['', 'critical', 'high', 'medium', 'low', 'info'].map(s => (
          <button
            key={s || 'all'}
            type='button'
            onClick={() => setFilter(s)}
            className={`rounded border px-3 py-1 font-mono text-[10px] tracking-widest uppercase ${
              filter === s ? 'border-primary bg-primary/10 text-primary' : 'border-border text-muted-foreground'
            }`}
          >
            {s || 'all'}
          </button>
        ))}
      </div>

      {error && (
        <div className='border-destructive/40 bg-destructive/10 text-destructive rounded-md border px-4 py-3 font-mono text-xs'>
          {error}
        </div>
      )}

      {!items ? (
        <Skeleton className='h-40 w-full' />
      ) : items.length === 0 ? (
        <p className='micro-label py-12 text-center'>NO FINDINGS YET — COMPLETE A RUN TO POPULATE</p>
      ) : (
        <div className='flex flex-col gap-3'>
          {items.map((it, i) => (
            <Card key={`${it.run_id}-${it.finding.id}-${i}`}>
              <CardHeader className='pb-2'>
                <div className='flex flex-wrap items-center gap-2'>
                  <Badge variant={sevVariant(it.finding.severity)} className='font-mono text-[10px] uppercase'>
                    {it.finding.severity}
                  </Badge>
                  {it.finding.evidence?.passed && (
                    <Badge variant='outline' className='border-primary/40 text-primary font-mono text-[10px]'>
                      3-GATE
                    </Badge>
                  )}
                  <span className='text-muted-foreground font-mono text-[10px]'>{it.finding.id}</span>
                  <Link href={`/runs/${it.run_id}`} className='text-primary ml-auto font-mono text-[10px] tracking-widest underline'>
                    {it.target} · {it.run_id.slice(0, 8)}
                  </Link>
                </div>
                <CardTitle className='mt-1 font-mono text-sm'>{it.finding.title}</CardTitle>
              </CardHeader>
              <CardContent className='text-muted-foreground font-mono text-xs'>{it.finding.description}</CardContent>
            </Card>
          ))}
        </div>
      )}
    </div>
  )
}

export default FindingsView
