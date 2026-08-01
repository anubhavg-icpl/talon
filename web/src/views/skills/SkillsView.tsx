'use client'

import { useCallback, useEffect, useMemo, useState } from 'react'

import type { CategoryCount, Skill } from '@/lib/api'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import { ScrollArea } from '@/components/ui/scroll-area'
import { Skeleton } from '@/components/ui/skeleton'
import { getSkill, getSkills } from '@/lib/api'

const STAGES = ['', 'recon', 'exploit', 'post_exploit', 'codegen', 'report', 'all'] as const
const PAGE = 40
const RECENT_KEY = 'talon.skills.recent'

const loadRecent = (): string[] => {
  try {
    const raw = localStorage.getItem(RECENT_KEY)
    return raw ? (JSON.parse(raw) as string[]) : []
  } catch {
    return []
  }
}

const pushRecent = (id: string) => {
  try {
    const prev = loadRecent().filter(x => x !== id)
    const next = [id, ...prev].slice(0, 24)
    localStorage.setItem(RECENT_KEY, JSON.stringify(next))
    return next
  } catch {
    return [id]
  }
}

const SkillsView = () => {
  const [skills, setSkills] = useState<Skill[]>([])
  const [categories, setCategories] = useState<CategoryCount[]>([])
  const [stats, setStats] = useState<Record<string, number> | null>(null)
  const [total, setTotal] = useState(0)
  const [offset, setOffset] = useState(0)
  const [stage, setStage] = useState('')
  const [category, setCategory] = useState('')
  const [q, setQ] = useState('')
  const [qDraft, setQDraft] = useState('')
  const [error, setError] = useState<string | null>(null)
  const [loading, setLoading] = useState(true)
  const [selected, setSelected] = useState<Skill | null>(null)
  const [detailLoading, setDetailLoading] = useState(false)
  const [recent, setRecent] = useState<string[]>([])
  const [catFilter, setCatFilter] = useState('')

  useEffect(() => {
    setRecent(loadRecent())
  }, [])

  const load = useCallback(async () => {
    setLoading(true)
    try {
      const res = await getSkills({
        brief: true,
        stage: stage || undefined,
        category: category || undefined,
        q: q || undefined,
        limit: PAGE,
        offset
      })
      setSkills(res.skills ?? [])
      setTotal(res.total ?? 0)
      setStats(res.stats ?? null)
      setCategories(res.categories ?? [])
      setError(null)
    } catch (err) {
      setError(err instanceof Error ? err.message : String(err))
    } finally {
      setLoading(false)
    }
  }, [stage, category, q, offset])

  useEffect(() => {
    void load()
  }, [load])

  // Debounce search
  useEffect(() => {
    const t = setTimeout(() => {
      setOffset(0)
      setQ(qDraft.trim())
    }, 300)
    return () => clearTimeout(t)
  }, [qDraft])

  const openSkill = async (id: string) => {
    setDetailLoading(true)
    try {
      const full = await getSkill(id)
      setSelected(full)
      setRecent(pushRecent(id))
    } catch (err) {
      setError(err instanceof Error ? err.message : String(err))
    } finally {
      setDetailLoading(false)
    }
  }

  const filteredCats = useMemo(() => {
    const needle = catFilter.toLowerCase()
    if (!needle) return categories
    return categories.filter(c => c.name.toLowerCase().includes(needle))
  }, [categories, catFilter])

  const page = Math.floor(offset / PAGE) + 1
  const pages = Math.max(1, Math.ceil(total / PAGE))

  return (
    <div className='flex h-[calc(100vh-8rem)] flex-col gap-4'>
      <div className='flex flex-wrap items-end justify-between gap-3'>
        <div>
          <h1 className='font-mono text-xl font-semibold tracking-widest'>SKILLS</h1>
          <p className='micro-label mt-1'>
            CYBERSTRIKE + BUILTIN CATALOG — BROWSE · SEARCH · REVISIT · INJECTED INTO AGENTS
          </p>
        </div>
        {stats && (
          <div className='flex flex-wrap gap-2 font-mono text-[10px] tracking-widest uppercase'>
            <Badge variant='outline'>TOTAL {stats.total ?? 0}</Badge>
            <Badge variant='outline'>DISK {stats.src_disk ?? 0}</Badge>
            <Badge variant='outline'>BUILTIN {stats.src_builtin ?? 0}</Badge>
          </div>
        )}
      </div>

      {error && (
        <div className='border-destructive/40 bg-destructive/10 text-destructive rounded-md border px-4 py-3 font-mono text-xs'>
          {error}
        </div>
      )}

      <div className='grid min-h-0 flex-1 gap-4 lg:grid-cols-[240px_1fr_minmax(280px,380px)]'>
        {/* Categories sidebar */}
        <Card className='hud-corners flex min-h-0 flex-col overflow-hidden'>
          <CardHeader className='pb-2'>
            <CardTitle className='micro-label'>CATEGORIES</CardTitle>
            <Input
              value={catFilter}
              onChange={e => setCatFilter(e.target.value)}
              placeholder='Filter…'
              className='mt-2 font-mono text-[11px]'
            />
          </CardHeader>
          <CardContent className='min-h-0 flex-1 p-0'>
            <ScrollArea className='h-full max-h-[calc(100vh-16rem)] px-3 pb-3'>
              <button
                type='button'
                onClick={() => {
                  setCategory('')
                  setOffset(0)
                }}
                className={`mb-1 w-full rounded px-2 py-1.5 text-left font-mono text-[11px] ${
                  !category ? 'bg-primary/15 text-primary' : 'text-muted-foreground hover:bg-muted/40'
                }`}
              >
                All categories
              </button>
              {filteredCats.map(c => (
                <button
                  key={c.name}
                  type='button'
                  onClick={() => {
                    setCategory(c.name)
                    setOffset(0)
                  }}
                  className={`mb-0.5 flex w-full items-center justify-between rounded px-2 py-1.5 text-left font-mono text-[11px] ${
                    category === c.name ? 'bg-primary/15 text-primary' : 'text-muted-foreground hover:bg-muted/40'
                  }`}
                >
                  <span className='truncate'>{c.name}</span>
                  <span className='text-[10px] opacity-70'>{c.count}</span>
                </button>
              ))}
            </ScrollArea>
          </CardContent>
        </Card>

        {/* List */}
        <div className='flex min-h-0 flex-col gap-3'>
          <div className='flex flex-wrap items-center gap-2'>
            <Input
              value={qDraft}
              onChange={e => setQDraft(e.target.value)}
              placeholder='Search name, id, path, body…'
              className='max-w-md font-mono text-xs'
            />
            <div className='flex flex-wrap gap-1'>
              {STAGES.map(st => (
                <button
                  key={st || 'any'}
                  type='button'
                  onClick={() => {
                    setStage(st)
                    setOffset(0)
                  }}
                  className={`rounded border px-2 py-1 font-mono text-[10px] tracking-widest uppercase ${
                    stage === st ? 'border-primary bg-primary/10 text-primary' : 'border-border text-muted-foreground'
                  }`}
                >
                  {st || 'any stage'}
                </button>
              ))}
            </div>
          </div>

          {recent.length > 0 && (
            <div className='flex flex-wrap items-center gap-1'>
              <span className='micro-label mr-1'>RECENT</span>
              {recent.slice(0, 8).map(id => (
                <button
                  key={id}
                  type='button'
                  onClick={() => void openSkill(id)}
                  className='border-border text-muted-foreground hover:text-primary max-w-[140px] truncate rounded border px-1.5 py-0.5 font-mono text-[10px]'
                  title={id}
                >
                  {id.replace(/^disk-cyberstrike-/, '').slice(0, 28)}
                </button>
              ))}
            </div>
          )}

          <div className='text-muted-foreground font-mono text-[10px] tracking-widest uppercase'>
            {total} match · page {page}/{pages}
          </div>

          <ScrollArea className='min-h-0 flex-1 rounded-md border'>
            <div className='flex flex-col gap-1 p-2'>
              {loading ? (
                Array.from({ length: 8 }).map((_, i) => <Skeleton key={i} className='h-14 w-full' />)
              ) : skills.length === 0 ? (
                <p className='micro-label py-12 text-center'>NO SKILLS MATCH</p>
              ) : (
                skills.map(sk => (
                  <button
                    key={sk.id}
                    type='button'
                    onClick={() => void openSkill(sk.id)}
                    className={`hover:border-primary/50 rounded-md border px-3 py-2 text-left transition-colors ${
                      selected?.id === sk.id ? 'border-primary bg-primary/10' : 'border-border/60 bg-card/40'
                    }`}
                  >
                    <div className='flex flex-wrap items-center gap-2'>
                      <Badge variant='outline' className='font-mono text-[9px] uppercase'>
                        {sk.stage}
                      </Badge>
                      {sk.category && (
                        <Badge variant='secondary' className='font-mono text-[9px]'>
                          {sk.category}
                        </Badge>
                      )}
                      {sk.source && (
                        <span className='text-muted-foreground font-mono text-[9px]'>{sk.source}</span>
                      )}
                    </div>
                    <p className='mt-1 font-mono text-xs font-medium tracking-wide'>{sk.name}</p>
                    <p className='text-muted-foreground truncate font-mono text-[10px]'>{sk.id}</p>
                  </button>
                ))
              )}
            </div>
          </ScrollArea>

          <div className='flex items-center justify-between gap-2'>
            <Button
              variant='outline'
              size='sm'
              disabled={offset <= 0 || loading}
              onClick={() => setOffset(o => Math.max(0, o - PAGE))}
              className='font-mono text-[10px] tracking-widest uppercase'
            >
              Prev
            </Button>
            <span className='font-mono text-[10px]'>
              {offset + 1}–{Math.min(offset + PAGE, total)} of {total}
            </span>
            <Button
              variant='outline'
              size='sm'
              disabled={offset + PAGE >= total || loading}
              onClick={() => setOffset(o => o + PAGE)}
              className='font-mono text-[10px] tracking-widest uppercase'
            >
              Next
            </Button>
          </div>
        </div>

        {/* Detail pane */}
        <Card className='hud-corners flex min-h-0 flex-col overflow-hidden'>
          <CardHeader className='pb-2'>
            <CardTitle className='micro-label'>SKILL DETAIL</CardTitle>
          </CardHeader>
          <CardContent className='min-h-0 flex-1 overflow-hidden'>
            {detailLoading ? (
              <Skeleton className='h-48 w-full' />
            ) : !selected ? (
              <p className='micro-label py-16 text-center'>SELECT A SKILL TO READ FULL METHODOLOGY</p>
            ) : (
              <ScrollArea className='h-full max-h-[calc(100vh-16rem)] pr-3'>
                <div className='flex flex-wrap gap-2'>
                  <Badge variant='outline' className='font-mono text-[10px] uppercase'>
                    {selected.stage}
                  </Badge>
                  {selected.category && (
                    <Badge variant='secondary' className='font-mono text-[10px]'>
                      {selected.category}
                    </Badge>
                  )}
                  {selected.source && (
                    <Badge variant='outline' className='font-mono text-[10px]'>
                      {selected.source}
                    </Badge>
                  )}
                </div>
                <h2 className='mt-3 font-mono text-sm font-semibold tracking-wide'>{selected.name}</h2>
                <p className='text-muted-foreground mt-1 break-all font-mono text-[10px]'>{selected.id}</p>
                {selected.path && (
                  <p className='text-muted-foreground mt-0.5 font-mono text-[10px]'>path: {selected.path}</p>
                )}
                <pre className='text-muted-foreground mt-4 whitespace-pre-wrap font-mono text-[11px] leading-relaxed'>
                  {selected.body || '(empty body)'}
                </pre>
              </ScrollArea>
            )}
          </CardContent>
        </Card>
      </div>
    </div>
  )
}

export default SkillsView
