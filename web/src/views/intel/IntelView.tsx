'use client'

import { useEffect, useState } from 'react'

import Link from 'next/link'

import type { IntelEvent } from '@/lib/api'
import { Badge } from '@/components/ui/badge'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Skeleton } from '@/components/ui/skeleton'
import { getIntel } from '@/lib/api'

const IntelView = () => {
  const [events, setEvents] = useState<IntelEvent[] | null>(null)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    let mounted = true
    const load = () =>
      getIntel(80)
        .then(r => {
          if (!mounted) return
          setEvents(r.events ?? [])
          setError(null)
        })
        .catch(e => mounted && setError(e instanceof Error ? e.message : String(e)))
    load()
    const id = setInterval(load, 8000)
    return () => {
      mounted = false
      clearInterval(id)
    }
  }, [])

  return (
    <div className='flex flex-col gap-6'>
      <div>
        <h1 className='font-mono text-xl font-semibold tracking-widest'>INTEL FEED</h1>
        <p className='micro-label mt-1'>CROSS-RUN FINDINGS + HISTORY · LIVE REFRESH</p>
      </div>
      {error && (
        <div className='border-destructive/40 bg-destructive/10 text-destructive rounded-md border px-4 py-3 font-mono text-xs'>
          {error}
        </div>
      )}
      <Card>
        <CardHeader>
          <CardTitle className='micro-label'>LIVE STREAM</CardTitle>
        </CardHeader>
        <CardContent className='space-y-2'>
          {!events ? (
            <Skeleton className='h-40 w-full' />
          ) : events.length === 0 ? (
            <p className='micro-label py-8 text-center'>NO INTEL YET — RUN AN OPERATION</p>
          ) : (
            events.map((ev, i) => (
              <div key={`${ev.run_id}-${i}`} className='border-border/50 rounded border px-3 py-2 font-mono text-xs'>
                <div className='flex flex-wrap items-center gap-2'>
                  <Badge variant={ev.kind === 'finding' ? 'destructive' : 'outline'} className='text-[9px] uppercase'>
                    {ev.kind}
                  </Badge>
                  {ev.severity && (
                    <Badge variant='secondary' className='text-[9px] uppercase'>
                      {ev.severity}
                    </Badge>
                  )}
                  <Link href={`/runs/${ev.run_id}`} className='text-primary ml-auto text-[10px] underline'>
                    {ev.target} · {ev.run_id.slice(0, 8)}
                  </Link>
                </div>
                <p className='mt-1 font-medium'>{ev.label}</p>
                {ev.detail && <p className='text-muted-foreground mt-0.5 line-clamp-2 text-[11px]'>{ev.detail}</p>}
                <p className='text-muted-foreground mt-1 text-[10px]'>{ev.at}</p>
              </div>
            ))
          )}
        </CardContent>
      </Card>
    </div>
  )
}

export default IntelView
