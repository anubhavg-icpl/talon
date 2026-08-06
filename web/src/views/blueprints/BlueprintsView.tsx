'use client'

import { useEffect, useState } from 'react'

import type { Blueprint, BlueprintCategory } from '@/lib/api'
import { listBlueprints } from '@/lib/api'
import { Badge } from '@/components/ui/badge'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Skeleton } from '@/components/ui/skeleton'

const CATEGORIES: Array<BlueprintCategory | ''> = ['', 'recon', 'exploit', 'post-exploit', 'reporting']

const categoryColor = (cat: string): string => {
  switch (cat) {
    case 'recon':
      return 'border-primary/40 text-primary'
    case 'exploit':
      return 'border-destructive/40 text-destructive'
    case 'post-exploit':
      return 'border-orange-500/40 text-orange-500'
    case 'reporting':
      return 'border-emerald-500/40 text-emerald-500'
    default:
      return 'border-border text-muted-foreground'
  }
}

const BlueprintsView = () => {
  const [filter, setFilter] = useState<BlueprintCategory | ''>('')
  const [items, setItems] = useState<Blueprint[] | null>(null)
  const [error, setError] = useState<string | null>(null)
  const [expanded, setExpanded] = useState<Set<string>>(new Set())

  useEffect(() => {
    let stale = false
    setItems(null)
    listBlueprints(filter || undefined)
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
  }, [filter])

  const toggleExpand = (id: string) => {
    setExpanded(prev => {
      const next = new Set(prev)
      if (next.has(id)) next.delete(id)
      else next.add(id)
      return next
    })
  }

  return (
    <div className='flex flex-col gap-6'>
      <div>
        <h1 className='font-mono text-xl font-semibold tracking-widest'>BLUEPRINTS</h1>
        <p className='micro-label mt-1'>REUSABLE PENTEST PLAYBOOK TEMPLATES — STRUCTURED STEP-BY-STEP METHODOLOGY</p>
      </div>

      <div className='flex flex-wrap gap-2'>
        {CATEGORIES.map(c => (
          <button
            key={c || 'all'}
            type='button'
            onClick={() => setFilter(c)}
            className={`rounded border px-3 py-1 font-mono text-[10px] tracking-widest uppercase ${
              filter === c ? 'border-primary bg-primary/10 text-primary' : 'border-border text-muted-foreground'
            }`}
          >
            {c || 'all'}
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
        <p className='micro-label py-12 text-center'>NO BLUEPRINTS FOUND</p>
      ) : (
        <div className='grid gap-4 md:grid-cols-2 xl:grid-cols-3'>
          {items.map(bp => {
            const isOpen = expanded.has(bp.id)
            return (
              <Card key={bp.id} className='hud-corners'>
                <CardHeader className='pb-2'>
                  <div className='flex flex-wrap items-center gap-2'>
                    <span className={`rounded border px-2 py-0.5 font-mono text-[9px] tracking-widest uppercase ${categoryColor(bp.category)}`}>
                      {bp.category}
                    </span>
                    <Badge variant='outline' className='font-mono text-[9px]'>{bp.phase}</Badge>
                    <Badge variant='secondary' className='font-mono text-[9px]'>v{bp.version}</Badge>
                  </div>
                  <button type='button' onClick={() => toggleExpand(bp.id)} className='mt-1 flex items-center gap-2 text-left'>
                    <CardTitle className='font-mono text-sm tracking-widest'>{bp.name}</CardTitle>
                    <span className='text-muted-foreground font-mono text-[10px]'>{isOpen ? '▲' : '▼'}</span>
                  </button>
                </CardHeader>
                <CardContent className='space-y-3'>
                  <p className='text-muted-foreground font-mono text-[11px] leading-relaxed'>{bp.description}</p>

                  {isOpen && (
                    <div className='space-y-2'>
                      <div className='micro-label'>STEPS ({bp.steps?.length ?? 0})</div>
                      {bp.steps?.map((step, i) => (
                        <div key={i} className='rounded border border-border/50 bg-black/30 p-2'>
                          <div className='flex items-center gap-2'>
                            <span className='text-primary font-mono text-[10px] font-bold'>
                              {step.order || i + 1}
                            </span>
                            <span className='font-mono text-[11px] font-semibold'>{step.tool}</span>
                          </div>
                          <p className='text-muted-foreground mt-1 font-mono text-[10px]'>{step.description}</p>
                          {step.expected_result && (
                            <p className='mt-1 font-mono text-[10px] text-emerald-500/80'>
                              ✓ {step.expected_result}
                            </p>
                          )}
                          {step.on_failure && (
                            <p className='font-mono text-[10px] text-orange-500/80'>⚠ {step.on_failure}</p>
                          )}
                        </div>
                      ))}
                      <div className='flex flex-wrap gap-1 pt-1'>
                        {bp.tags?.map(t => (
                          <Badge key={t} variant='secondary' className='font-mono text-[9px]'>
                            {t}
                          </Badge>
                        ))}
                      </div>
                      <div className='micro-label pt-1'>AUTHOR</div>
                      <div className='text-muted-foreground font-mono text-[10px]'>{bp.author}</div>
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

export default BlueprintsView
