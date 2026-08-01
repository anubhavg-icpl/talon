'use client'

import { useCallback, useEffect, useRef, useState } from 'react'
import { Bot, Loader2, Send, Square, Wrench } from 'lucide-react'
import { toast } from 'sonner'

import PageHeader from '@/components/shared/PageHeader'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { ScrollArea } from '@/components/ui/scroll-area'
import { Textarea } from '@/components/ui/textarea'
import type { LLMStreamMessage, SLMToolDef } from '@/lib/api'
import { listLLMTools, llmInfo, streamLLMAssist } from '@/lib/api'
import { cn } from '@/lib/utils'

type ChatRole = 'user' | 'assistant' | 'system' | 'tool'

type ChatItem = {
  id: string
  role: ChatRole
  content: string
  toolName?: string
  ms?: number
  streaming?: boolean
}

const uid = () => `${Date.now()}-${Math.random().toString(36).slice(2, 8)}`

const SUGGESTIONS = [
  'What runs are active right now?',
  'Is the stack healthy (ollama, onnx-slm, msf)?',
  'Search skills for SSRF methodology',
  'List agents and playbooks',
  'Summarize recent findings from intel feed'
]

const AssistView = () => {
  const [items, setItems] = useState<ChatItem[]>([])
  const [draft, setDraft] = useState('')
  const [busy, setBusy] = useState(false)
  const [tools, setTools] = useState<SLMToolDef[]>([])
  const [meta, setMeta] = useState<Record<string, unknown> | null>(null)
  const [info, setInfo] = useState<{ provider?: string; model?: string; tool_count?: number } | null>(null)
  const stopRef = useRef<(() => void) | null>(null)
  const bottomRef = useRef<HTMLDivElement>(null)

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
  }, [items, busy])

  const stop = useCallback(() => {
    stopRef.current?.()
    stopRef.current = null
    setBusy(false)
    setItems(prev => prev.map(m => (m.streaming ? { ...m, streaming: false } : m)))
  }, [])

  const send = useCallback(
    (text: string) => {
      const content = text.trim()
      if (!content || busy) return

      const userMsg: ChatItem = { id: uid(), role: 'user', content }
      const asstId = uid()
      const asstMsg: ChatItem = { id: asstId, role: 'assistant', content: '', streaming: true }

      setItems(prev => [...prev, userMsg, asstMsg])
      setDraft('')
      setBusy(true)
      setMeta(null)

      // Build history for the API (user/assistant only).
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
          onToken: tok => {
            setItems(prev =>
              prev.map(m => (m.id === asstId ? { ...m, content: m.content + tok } : m))
            )
          },
          onToolStart: t => {
            setItems(prev => [
              ...prev,
              {
                id: uid(),
                role: 'tool',
                toolName: t.name,
                content: `→ calling ${t.name}(${JSON.stringify(t.arguments ?? {})})`
              }
            ])
          },
          onToolResult: t => {
            setItems(prev => [
              ...prev,
              {
                id: uid(),
                role: 'tool',
                toolName: t.name,
                ms: t.ms,
                content: t.result.slice(0, 4000)
              }
            ])
          },
          onDone: r => {
            setItems(prev => {
              const withoutEmptyAsst = prev.filter(m => !(m.id === asstId && !m.content && !r.text))
              const hasAsst = withoutEmptyAsst.some(m => m.id === asstId)
              if (hasAsst) {
                return withoutEmptyAsst.map(m =>
                  m.id === asstId
                    ? { ...m, content: r.text || m.content, streaming: false, ms: r.ms }
                    : m
                )
              }
              return [
                ...withoutEmptyAsst,
                { id: asstId, role: 'assistant', content: r.text || '(no reply)', streaming: false, ms: r.ms }
              ]
            })
            setBusy(false)
            stopRef.current = null
            if (r.tool_calls) {
              toast.success(`Assist done · ${r.tool_calls} tool call(s) · ${r.ms}ms`)
            }
          },
          onError: err => {
            setItems(prev =>
              prev.map(m =>
                m.id === asstId
                  ? { ...m, content: m.content || err.message, streaming: false }
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
        title='SLM Assist'
        description='Local SmolLM / ONNX copilot with curated codebase tools — runs, skills, health, agents. Read-only; full exploits stay on agent runs.'
      />

      <div className='grid flex-1 gap-4 lg:grid-cols-[1fr_280px]'>
        <Card className='flex min-h-[560px] flex-col border-border/80'>
          <CardHeader className='flex flex-row items-center justify-between space-y-0 border-b py-3'>
            <div className='flex items-center gap-2'>
              <Bot className='text-primary size-4' />
              <CardTitle className='font-mono text-sm tracking-widest uppercase'>Console chat</CardTitle>
              {busy && <Loader2 className='text-muted-foreground size-3.5 animate-spin' />}
            </div>
            <div className='flex items-center gap-2'>
              {info?.provider && (
                <Badge variant='outline' className='font-mono text-[10px] tracking-wider uppercase'>
                  {String(info.provider)} · {String(info.model ?? '')}
                </Badge>
              )}
              {meta?.tools_enabled != null && (
                <Badge variant='secondary' className='font-mono text-[10px]'>
                  tools {meta.tools_enabled ? 'on' : 'off'}
                </Badge>
              )}
            </div>
          </CardHeader>
          <CardContent className='flex flex-1 flex-col gap-3 p-0'>
            <ScrollArea className='h-[420px] flex-1 px-4 py-3'>
              {items.length === 0 && (
                <div className='text-muted-foreground space-y-3 py-8 text-center text-sm'>
                  <p className='font-mono text-xs tracking-widest uppercase'>Ask about runs, skills, stack health</p>
                  <div className='flex flex-wrap justify-center gap-2'>
                    {SUGGESTIONS.map(s => (
                      <Button
                        key={s}
                        type='button'
                        size='sm'
                        variant='outline'
                        className='h-auto max-w-xs whitespace-normal py-1.5 text-left font-mono text-[11px]'
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
                {items.map(m => (
                  <div
                    key={m.id}
                    className={cn(
                      'rounded-md border px-3 py-2 text-sm',
                      m.role === 'user' && 'bg-primary/5 border-primary/20 ml-8',
                      m.role === 'assistant' && 'bg-card border-border mr-8',
                      m.role === 'tool' && 'bg-muted/40 border-dashed font-mono text-xs'
                    )}
                  >
                    <div className='text-muted-foreground mb-1 flex items-center gap-1.5 font-mono text-[10px] tracking-widest uppercase'>
                      {m.role === 'tool' ? (
                        <>
                          <Wrench className='size-3' />
                          {m.toolName || 'tool'}
                          {m.ms != null && <span className='opacity-70'>{m.ms}ms</span>}
                        </>
                      ) : (
                        m.role
                      )}
                      {m.streaming && <span className='text-primary animate-pulse'>streaming</span>}
                    </div>
                    <pre className='font-sans wrap-break-word whitespace-pre-wrap'>{m.content || (m.streaming ? '…' : '')}</pre>
                  </div>
                ))}
                <div ref={bottomRef} />
              </div>
            </ScrollArea>

            <div className='border-t p-3'>
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
              <p className='text-muted-foreground mt-2 font-mono text-[10px] tracking-wide'>
                POST /llm/assist · SSE tools · read-only catalog · Enter to send
              </p>
            </div>
          </CardContent>
        </Card>

        <Card className='border-border/80 h-fit'>
          <CardHeader className='py-3'>
            <CardTitle className='flex items-center gap-2 font-mono text-xs tracking-widest uppercase'>
              <Wrench className='size-3.5' />
              Tools ({tools.length || info?.tool_count || 0})
            </CardTitle>
          </CardHeader>
          <CardContent className='space-y-2 px-3 pb-4'>
            <p className='text-muted-foreground text-xs'>
              SmolLM calls these via <code className='text-primary'>TOOL_CALL</code> against the live store / skills / health — same surface the UI uses.
            </p>
            <ScrollArea className='h-[420px] pr-2'>
              <ul className='space-y-2'>
                {tools.map(t => (
                  <li key={t.name} className='rounded border border-border/60 p-2'>
                    <div className='font-mono text-[11px] font-semibold tracking-wide'>{t.name}</div>
                    <p className='text-muted-foreground mt-0.5 text-[11px] leading-snug'>{t.description}</p>
                  </li>
                ))}
                {tools.length === 0 && (
                  <li className='text-muted-foreground text-xs'>Loading tool catalog from GET /llm/tools…</li>
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
