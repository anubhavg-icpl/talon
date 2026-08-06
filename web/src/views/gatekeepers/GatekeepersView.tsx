'use client'

import { useEffect, useState } from 'react'

import { toast } from 'sonner'

import type { GatekeeperActionLog, GatekeeperConfig } from '@/lib/api'
import { getGatekeeperActions, listGatekeepers } from '@/lib/api'
import { Badge } from '@/components/ui/badge'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Skeleton } from '@/components/ui/skeleton'

const GatekeepersView = () => {
  const [items, setItems] = useState<GatekeeperConfig[] | null>(null)
  const [error, setError] = useState<string | null>(null)
  const [expanded, setExpanded] = useState<Set<string>>(new Set())
  const [actions, setActions] = useState<Record<string, GatekeeperActionLog[]>>({})

  useEffect(() => {
    let stale = false
    listGatekeepers()
      .then(res => {
        if (stale) return
        setItems(res ?? [])
        setError(null)
      })
      .catch(err => {
        if (stale) return
        setError(err instanceof Error ? err.message : String(err))
        setItems([])
      })
    return () => {
      stale = true
    }
  }, [])

  const toggleExpand = async (name: string) => {
    const next = new Set(expanded)
    if (next.has(name)) {
      next.delete(name)
      setExpanded(next)
      return
    }
    next.add(name)
    setExpanded(next)
    try {
      const res = await getGatekeeperActions(name)
      setActions(prev => ({ ...prev, [name]: res ?? [] }))
    } catch (err) {
      toast.error(err instanceof Error ? err.message : 'failed to load actions')
      setActions(prev => ({ ...prev, [name]: [] }))
    }
  }

  return (
    <div className='flex flex-col gap-6'>
      <div>
        <h1 className='font-mono text-xl font-semibold tracking-widest'>GATEKEEPERS</h1>
        <p className='micro-label mt-1'>CAPABILITY-BASED ACCESS CONTROL — REGISTERED TOOL GATEKEEPERS & ACTION LOGS</p>
      </div>

      {error && (
        <div className='border-destructive/40 bg-destructive/10 text-destructive rounded-md border px-4 py-3 font-mono text-xs'>
          {error}
        </div>
      )}

      {!items ? (
        <Skeleton className='h-40 w-full' />
      ) : items.length === 0 ? (
        <p className='micro-label py-12 text-center'>NO GATEKEEPERS REGISTERED</p>
      ) : (
        <div className='flex flex-col gap-3'>
          {items.map(gk => {
            const isOpen = expanded.has(gk.name)
            const gkActions = actions[gk.name]
            return (
              <Card key={gk.name}>
                <CardHeader className='pb-2'>
                  <div className='flex flex-wrap items-center gap-2'>
                    <Badge variant='outline' className='font-mono text-[9px] uppercase'>{gk.type}</Badge>
                    <Badge
                      variant={gk.require_approval ? 'destructive' : 'secondary'}
                      className='font-mono text-[9px] uppercase'
                    >
                      {gk.require_approval ? 'REQUIRES APPROVAL' : 'AUTO'}
                    </Badge>
                    <span className='text-muted-foreground font-mono text-[10px] uppercase'>{gk.auth_type}</span>
                  </div>
                  <button type='button' onClick={() => toggleExpand(gk.name)} className='mt-1 flex items-center gap-2 text-left'>
                    <CardTitle className='font-mono text-sm tracking-widest'>{gk.name}</CardTitle>
                    <span className='text-muted-foreground font-mono text-[10px]'>{isOpen ? '▲' : '▼'}</span>
                  </button>
                </CardHeader>
                <CardContent className='space-y-3'>
                  <div className='grid grid-cols-2 gap-4 font-mono text-[11px]'>
                    <div>
                      <div className='micro-label mb-1'>ALLOWED TOOLS</div>
                      <div className='flex flex-wrap gap-1'>
                        {(gk.allowed_tools ?? []).length === 0 ? (
                          <span className='text-muted-foreground'>—</span>
                        ) : (
                          (gk.allowed_tools ?? []).map(t => (
                            <Badge key={t} variant='secondary' className='font-mono text-[9px]'>{t}</Badge>
                          ))
                        )}
                      </div>
                    </div>
                    <div>
                      <div className='micro-label mb-1'>SCOPES</div>
                      <div className='flex flex-wrap gap-1'>
                        {(gk.scopes ?? []).length === 0 ? (
                          <span className='text-muted-foreground'>—</span>
                        ) : (
                          (gk.scopes ?? []).map(s => (
                            <Badge key={s} variant='outline' className='font-mono text-[9px]'>{s}</Badge>
                          ))
                        )}
                      </div>
                    </div>
                  </div>

                  {isOpen && (
                    <div className='space-y-2'>
                      <div className='micro-label'>ACTION LOG</div>
                      {!gkActions ? (
                        <p className='text-muted-foreground font-mono text-[11px]'>Loading…</p>
                      ) : gkActions.length === 0 ? (
                        <p className='text-muted-foreground font-mono text-[11px]'>No actions recorded</p>
                      ) : (
                        <div className='max-h-64 space-y-1 overflow-auto'>
                          {gkActions.map(a => (
                            <div
                              key={a.id}
                              className='flex items-center gap-2 border-b border-border/40 py-1 font-mono text-[11px]'
                            >
                              <Badge
                                variant={a.approved ? 'secondary' : 'destructive'}
                                className='font-mono text-[9px]'
                              >
                                {a.approved ? 'APPROVED' : 'DENIED'}
                              </Badge>
                              <span className='font-semibold'>{a.action}</span>
                              <span className='text-muted-foreground truncate'>{a.resource}</span>
                              <span className='text-muted-foreground ml-auto shrink-0'>
                                {new Date(a.timestamp).toLocaleTimeString()}
                              </span>
                            </div>
                          ))}
                        </div>
                      )}
                    </div>
                  )}
                </CardContent>
              </Card>
            )
          })}
        </div>
      )}
    </div>
  )
}

export default GatekeepersView
