'use client'

import { useCallback, useEffect, useRef, useState } from 'react'
import { Activity, Loader2, Send, Sparkles, Square, Wrench } from 'lucide-react'
import { toast } from 'sonner'

import PageHeader from '@/components/shared/PageHeader'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Collapsible, CollapsibleContent, CollapsibleTrigger } from '@/components/ui/collapsible'
import { ScrollArea } from '@/components/ui/scroll-area'
import { Textarea } from '@/components/ui/textarea'
import type { LLMStreamMessage, SLMToolDef } from '@/lib/api'
import { listLLMTools, llmInfo, streamLLMAssist } from '@/lib/api'
import { cn } from '@/lib/utils'

import { AssistantMarkdown, ToolResultViz } from './ToolResultViz'

type ChatRole = 'user' | 'assistant' | 'tool'

type ChatItem = {
  id: string
  role: ChatRole
  content: string
  toolName?: string
  toolArgs?: string
  ms?: number
  streaming?: boolean
  /** tool lifecycle: running | done | error */
  toolState?: 'running' | 'done' | 'error'
}

const uid = () => `${Date.now()}-${Math.random().toString(36).slice(2, 8)}`

const SUGGESTIONS = [
  'Is the stack healthy (ollama, onnx-slm, msf)?',
  'What runs are active right now?',
  'Search skills for SSRF methodology',
  'List agents and playbooks',
  'Summarize recent findings from intel feed'
]

/** Insert item immediately before the streaming assistant bubble (chronological tools → answer). */
const insertBeforeAssistant = (prev: ChatItem[], asstId: string, item: ChatItem): ChatItem[] => {
  const idx = prev.findIndex(m => m.id === asstId)
  if (idx < 0) return [...prev, item]
  const next = [...prev]
  next.splice(idx, 0, item)
  return next
}

const ToolCard = ({ m }: { m: ChatItem }) => {
  const [open, setOpen] = useState(true)
  const viz =
    m.toolState === 'done' && m.content ? (
      <ToolResultViz toolName={m.toolName} content={m.content} />
    ) : null

  return (
    <Collapsible open={open} onOpenChange={setOpen} className='mr-4'>
      <div
        className={cn(
          'rounded-md border border-dashed px-3 py-2 font-mono text-xs',
          m.toolState === 'running' && 'border-primary/40 bg-primary/5',
          m.toolState === 'done' && 'bg-muted/30 border-border/70',
          m.toolState === 'error' && 'border-red-400/40 bg-red-500/5'
        )}
      >
        <CollapsibleTrigger className='flex w-full items-center gap-2 text-left'>
          {m.toolState === 'running' ? (
            <Loader2 className='text-primary size-3.5 shrink-0 animate-spin' />
          ) : (
            <Wrench className='text-muted-foreground size-3.5 shrink-0' />
          )}
          <span className='font-semibold tracking-wide'>{m.toolName || 'tool'}</span>
          {m.toolState === 'running' && (
            <Badge variant='outline' className='h-5 font-mono text-[9px] tracking-wider uppercase'>
              running
            </Badge>
          )}
          {m.toolState === 'done' && m.ms != null && (
            <span className='text-muted-foreground text-[10px]'>{m.ms}ms</span>
          )}
          {m.toolArgs && m.toolArgs !== '{}' && (
            <span className='text-muted-foreground truncate text-[10px]'>{m.toolArgs}</span>
          )}
          <span className='text-muted-foreground ml-auto text-[10px]'>{open ? '▾' : '▸'}</span>
        </CollapsibleTrigger>
        <CollapsibleContent className='mt-2'>
          {viz || (
            <pre className='bg-background/60 max-h-48 overflow-auto rounded border border-border/50 p-2 text-[10px] leading-relaxed whitespace-pre-wrap'>
              {m.content || (m.toolState === 'running' ? 'awaiting result…' : '')}
            </pre>
          )}
        </CollapsibleContent>
      </div>
    </Collapsible>
  )
}

const AssistView = () => {
  const [items, setItems] = useState<ChatItem[]>([])
  const [draft, setDraft] = useState('')
  const [busy, setBusy] = useState(false)
  const [tools, setTools] = useState<SLMToolDef[]>([])
  const [meta, setMeta] = useState<Record<string, unknown> | null>(null)
  const [phase, setPhase] = useState<string | null>(null)
  const [info, setInfo] = useState<{ provider?: string; model?: string; tool_count?: number } | null>(null)
  const stopRef = useRef<(() => void) | null>(null)
  const bottomRef = useRef<HTMLDivElement>(null)
  /** Maps tool name → chat item id for the active running call (update in place). */
  const toolRunRef = useRef<Map<string, string>>(new Map())

  useEffect(() => {
    listLLMTools()
      .then(r => setTools(r.tools ?? []))
      .catch(() => setTools([]))
    llmInfo()
      .then(r => setInfo({ provider: r.provider, model: r.model, tool_count: r.tool_count }))
      .catch(() => setInfo(null))
  }, [])

  useEffect(() => {
    bottomRef.current?.scrollIntoView({ behavior: 'smooth' })
  }, [items, busy, phase])

  const stop = useCallback(() => {
    stopRef.current?.()
    stopRef.current = null
    setBusy(false)
    setPhase(null)
    setItems(prev =>
      prev.map(m =>
        m.streaming || m.toolState === 'running'
          ? { ...m, streaming: false, toolState: m.toolState === 'running' ? 'error' : m.toolState }
          : m
      )
    )
  }, [])

  const send = useCallback(
    (text: string) => {
      const content = text.trim()
      if (!content || busy) return

      const userMsg: ChatItem = { id: uid(), role: 'user', content }
      const asstId = uid()
      const asstMsg: ChatItem = { id: asstId, role: 'assistant', content: '', streaming: true }

      toolRunRef.current = new Map()
      // Order: user → (tools inserted before asst) → assistant answer
      setItems(prev => [...prev, userMsg, asstMsg])
      setDraft('')
      setBusy(true)
      setMeta(null)
      setPhase('connecting…')

      const history: LLMStreamMessage[] = []
      for (const m of [...items, userMsg]) {
        if (m.role === 'user' || m.role === 'assistant') {
          if (m.content) history.push({ role: m.role, content: m.content })
        }
      }

      stopRef.current = streamLLMAssist(
        { messages: history, max_rounds: 5 },
        {
          onMeta: m => setMeta(m),
          onStatus: s => setPhase(s.message || s.phase || null),
          onToken: tok => {
            setPhase(null)
            setItems(prev => prev.map(m => (m.id === asstId ? { ...m, content: m.content + tok } : m)))
          },
          onToolStart: t => {
            setPhase(`tool · ${t.name}`)
            const tid = uid()
            toolRunRef.current.set(t.name, tid)
            const args = t.arguments ? JSON.stringify(t.arguments) : '{}'
            setItems(prev =>
              insertBeforeAssistant(prev, asstId, {
                id: tid,
                role: 'tool',
                toolName: t.name,
                toolArgs: args,
                toolState: 'running',
                content: ''
              })
            )
          },
          onToolResult: t => {
            setPhase(`tool done · ${t.name} (${t.ms}ms)`)
            const tid = toolRunRef.current.get(t.name)
            setItems(prev => {
              if (tid && prev.some(m => m.id === tid)) {
                return prev.map(m =>
                  m.id === tid
                    ? {
                        ...m,
                        toolState: 'done',
                        ms: t.ms,
                        content: t.result.slice(0, 8000)
                      }
                    : m
                )
              }
              // Fallback: insert completed tool before assistant
              return insertBeforeAssistant(prev, asstId, {
                id: uid(),
                role: 'tool',
                toolName: t.name,
                toolState: 'done',
                ms: t.ms,
                content: t.result.slice(0, 8000)
              })
            })
            toolRunRef.current.delete(t.name)
          },
          onDone: r => {
            setPhase(null)
            setItems(prev =>
              prev.map(m =>
                m.id === asstId
                  ? {
                      ...m,
                      content: (r.text && r.text.trim()) || m.content || '(no reply)',
                      streaming: false,
                      ms: r.ms
                    }
                  : m.toolState === 'running'
                    ? { ...m, toolState: 'done' as const }
                    : m
              )
            )
            setBusy(false)
            stopRef.current = null
            if (r.tool_calls) {
              toast.success(`Assist · ${r.tool_calls} tool(s) · ${r.ms}ms`)
            }
          },
          onError: err => {
            setPhase(null)
            setItems(prev =>
              prev.map(m =>
                m.id === asstId
                  ? { ...m, content: m.content || err.message, streaming: false }
                  : m.toolState === 'running'
                    ? { ...m, toolState: 'error' as const }
                    : m
              )
            )
            setBusy(false)
            stopRef.current = null
            toast.error(err.message)
          }
        }
      )
    },
    [busy, items]
  )

  return (
    <div className='flex h-full min-h-[70vh] flex-col gap-4'>
      <PageHeader
        title={
          <span className='inline-flex items-center gap-3'>
            {/* eslint-disable-next-line @next/next/no-img-element */}
            <img
              src='/assist/slm-assist-mark.webp'
              alt=''
              width={40}
              height={40}
              className='size-10 rounded-md border border-primary/30 object-cover shadow-[0_0_20px_oklch(0.62_0.22_25/0.25)]'
            />
            SLM Assist
          </span>
        }
        description='Local / hosted copilot with curated codebase tools — runs, skills, health. Read-only; full exploits stay on agent runs.'
      />

      {/* Pipeline hero */}
      <div className='hud-corners relative min-h-[100px] overflow-hidden rounded-sm border border-primary/20 sm:min-h-[120px]'>
        {/* eslint-disable-next-line @next/next/no-img-element */}
        <img
          src='/assist/slm-assist-banner.webp'
          alt=''
          className='absolute inset-0 size-full object-cover object-center'
        />
        <div className='absolute inset-0 bg-gradient-to-r from-background via-background/75 to-background/30' />
        <div className='absolute inset-0 bg-gradient-to-t from-background/90 via-transparent to-background/20' />
        <div className='relative flex h-full min-h-[100px] flex-wrap items-end gap-3 p-4 sm:min-h-[120px] sm:p-5'>
          <div className='flex flex-1 flex-col gap-1.5'>
            <p className='micro-label text-primary flex items-center gap-1.5'>
              <Sparkles className='size-3' />
              LOCAL COPILOT PIPELINE
            </p>
            <p className='font-mono text-[11px] tracking-wide text-zinc-200 sm:text-xs'>
              UI → <code className='text-primary'>/api/talon/llm/assist</code> → Go tools → provider · SSE tokens
            </p>
          </div>
          {info?.provider && (
            <Badge
              variant='outline'
              className='border-primary/40 bg-black/50 font-mono text-[10px] tracking-wider text-zinc-100 uppercase backdrop-blur-sm'
            >
              {info.provider} · {info.model}
            </Badge>
          )}
        </div>
      </div>

      <div className='grid flex-1 gap-4 lg:grid-cols-[1fr_300px]'>
        <Card className='hud-corners flex min-h-[560px] flex-col overflow-hidden border-border/80 pt-0'>
          <CardHeader className='flex flex-row items-center justify-between space-y-0 border-b border-primary/10 py-3'>
            <div className='flex min-w-0 items-center gap-2'>
              {/* eslint-disable-next-line @next/next/no-img-element */}
              <img src='/assist/slm-assist-mark.webp' alt='' className='size-6 rounded-sm object-cover' />
              <CardTitle className='font-mono text-sm tracking-widest uppercase'>Console chat</CardTitle>
              {busy && <Loader2 className='text-primary size-3.5 shrink-0 animate-spin' />}
              {busy && phase && (
                <span className='text-muted-foreground max-w-[200px] truncate font-mono text-[10px] tracking-wide sm:max-w-[280px]'>
                  {phase}
                </span>
              )}
            </div>
            <div className='flex shrink-0 items-center gap-2'>
              {meta?.tools_enabled != null && (
                <Badge variant='secondary' className='font-mono text-[10px]'>
                  tools {meta.tools_enabled ? 'on' : 'off'}
                </Badge>
              )}
            </div>
          </CardHeader>
          <CardContent className='flex flex-1 flex-col gap-3 p-0'>
            <ScrollArea className='h-[440px] flex-1 px-4 py-3'>
              {items.length === 0 && (
                <div className='text-muted-foreground space-y-5 py-6 text-center text-sm'>
                  {/* eslint-disable-next-line @next/next/no-img-element */}
                  <img
                    src='/assist/slm-assist-empty.webp'
                    alt=''
                    className='mx-auto h-36 w-auto rounded-sm border border-primary/15 object-cover shadow-[0_0_40px_oklch(0.62_0.22_25/0.12)] sm:h-44'
                  />
                  <div>
                    <p className='font-mono text-xs tracking-widest text-foreground uppercase'>
                      Ask about runs, skills, stack health
                    </p>
                    <p className='micro-label mt-1'>Curated tools · no exploit MCP here</p>
                  </div>
                  <div className='flex flex-wrap justify-center gap-2'>
                    {SUGGESTIONS.map(s => (
                      <Button
                        key={s}
                        type='button'
                        size='sm'
                        variant='outline'
                        className='h-auto max-w-xs whitespace-normal border-primary/25 py-1.5 text-left font-mono text-[11px]'
                        onClick={() => send(s)}
                        disabled={busy}
                      >
                        {s}
                      </Button>
                    ))}
                  </div>
                </div>
              )}
              <div className='space-y-3'>
                {items.map(m => {
                  if (m.role === 'tool') {
                    return <ToolCard key={m.id} m={m} />
                  }
                  return (
                    <div
                      key={m.id}
                      className={cn(
                        'rounded-md border px-3 py-2 text-sm',
                        m.role === 'user' && 'bg-primary/5 border-primary/20 ml-6 sm:ml-10',
                        m.role === 'assistant' && 'bg-card border-border mr-4'
                      )}
                    >
                      <div className='text-muted-foreground mb-1 flex items-center gap-1.5 font-mono text-[10px] tracking-widest uppercase'>
                        {m.role === 'assistant' ? (
                          <>
                            {/* eslint-disable-next-line @next/next/no-img-element */}
                            <img src='/assist/slm-assist-mark.webp' alt='' className='size-3.5 rounded-sm object-cover' />
                            assistant
                          </>
                        ) : (
                          'user'
                        )}
                        {m.streaming && <span className='text-primary animate-pulse'>streaming</span>}
                        {!m.streaming && m.ms != null && m.role === 'assistant' && (
                          <span className='opacity-70'>{m.ms}ms</span>
                        )}
                      </div>
                      {m.role === 'assistant' && !m.content && m.streaming ? (
                        <p className='text-muted-foreground font-mono text-xs'>
                          {phase ? `⏳ ${phase}` : '⏳ waiting on tools + model…'}
                        </p>
                      ) : m.role === 'assistant' ? (
                        <AssistantMarkdown text={m.content} />
                      ) : (
                        <div className='text-sm leading-relaxed whitespace-pre-wrap'>{m.content}</div>
                      )}
                    </div>
                  )
                })}
                <div ref={bottomRef} />
              </div>
            </ScrollArea>

            <div className='border-t border-primary/10 p-3'>
              <div className='flex gap-2'>
                <Textarea
                  value={draft}
                  onChange={e => setDraft(e.target.value)}
                  placeholder='Ask Talon Assist… (uses codebase tools)'
                  className='min-h-[72px] flex-1 font-mono text-xs'
                  disabled={busy}
                  onKeyDown={e => {
                    if (e.key === 'Enter' && !e.shiftKey) {
                      e.preventDefault()
                      send(draft)
                    }
                  }}
                />
                <div className='flex flex-col gap-2'>
                  {busy ? (
                    <Button type='button' variant='destructive' size='icon' onClick={stop} title='Stop'>
                      <Square className='size-4' />
                    </Button>
                  ) : (
                    <Button type='button' size='icon' onClick={() => send(draft)} disabled={!draft.trim()} title='Send'>
                      <Send className='size-4' />
                    </Button>
                  )}
                </div>
              </div>
              <p className='text-muted-foreground mt-2 flex items-center gap-1.5 font-mono text-[10px] tracking-wide'>
                <Activity className='size-3' />
                POST /llm/assist · SSE · tools before answer · Enter to send
              </p>
            </div>
          </CardContent>
        </Card>

        <Card className='hud-corners h-fit overflow-hidden border-border/80 pt-0'>
          <div className='relative h-24 border-b border-primary/15'>
            {/* eslint-disable-next-line @next/next/no-img-element */}
            <img
              src='/assist/slm-assist-tools.webp'
              alt=''
              className='absolute inset-0 size-full object-cover object-top'
            />
            <div className='absolute inset-0 bg-gradient-to-t from-card via-card/50 to-transparent' />
            <div className='absolute right-3 bottom-2 left-3'>
              <CardTitle className='flex items-center gap-2 font-mono text-xs tracking-widest text-zinc-100 uppercase'>
                <Wrench className='size-3.5 text-primary' />
                Tools ({tools.length || info?.tool_count || 0})
              </CardTitle>
            </div>
          </div>
          <CardContent className='space-y-2 px-3 pt-3 pb-4'>
            <p className='text-muted-foreground text-xs'>
              Curated read-only catalog. Full MCP exploits stay on agent runs.
            </p>
            <ScrollArea className='h-[380px] pr-2'>
              <ul className='space-y-2'>
                {tools.map(t => (
                  <li
                    key={t.name}
                    className='rounded-sm border border-border/60 border-l-primary/40 border-l-2 p-2 transition-colors hover:border-primary/30 hover:bg-primary/5'
                  >
                    <div className='font-mono text-[11px] font-semibold tracking-wide'>{t.name}</div>
                    <p className='text-muted-foreground mt-0.5 text-[11px] leading-snug'>{t.description}</p>
                  </li>
                ))}
                {tools.length === 0 && (
                  <li className='text-muted-foreground text-xs'>Loading GET /llm/tools…</li>
                )}
              </ul>
            </ScrollArea>
          </CardContent>
        </Card>
      </div>
    </div>
  )
}

export default AssistView
