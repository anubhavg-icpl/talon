'use client'

import { useEffect, useState } from 'react'

import Link from 'next/link'

import type { AgentInfo } from '@/lib/api'
import HudStill from '@/components/shared/HudStill'
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar'
import { Badge } from '@/components/ui/badge'
import { buttonVariants } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Skeleton } from '@/components/ui/skeleton'
import { agentAvatarSrc } from '@/lib/agent-avatars'
import { getAgents } from '@/lib/api'
import { cn } from '@/lib/utils'

const AgentsView = () => {
  const [agents, setAgents] = useState<AgentInfo[] | null>(null)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    getAgents()
      .then(res => setAgents(res.agents ?? []))
      .catch(err => setError(err instanceof Error ? err.message : String(err)))
  }, [])

  return (
    <div className='flex flex-col gap-6'>
      <div className='flex flex-wrap items-end justify-between gap-3'>
        <div>
          <h1 className='font-mono text-xl font-semibold tracking-widest'>AGENTS</h1>
          <p className='micro-label mt-1'>
            SPECIALIST MODES · A2A VIA ORCHESTRATOR DELEGATES · CYBERSTRIKE skill_search / skill_get
          </p>
        </div>
        <Link href='/runs/new' className={cn(buttonVariants(), 'font-mono text-xs tracking-widest uppercase')}>
          New Operation
        </Link>
      </div>

      <HudStill
        src='/showcase/talon-pipeline-agents.webp'
        alt='Multi-agent pipeline'
        variant='banner'
        className='max-h-28'
      />

      {error && (
        <div className='border-destructive/40 bg-destructive/10 text-destructive rounded-md border px-4 py-3 font-mono text-xs'>
          {error}
        </div>
      )}

      {!agents ? (
        <Skeleton className='h-40 w-full' />
      ) : (
        <div className='grid gap-4 md:grid-cols-2 xl:grid-cols-3'>
          {agents.map(a => {
            const src = agentAvatarSrc(a.id)
            const initials = a.codename.slice(0, 2)
            return (
              <Card key={a.id} className='hud-corners overflow-hidden pt-0'>
                <div className='relative h-36 w-full border-b border-primary/15 bg-black'>
                  {/* eslint-disable-next-line @next/next/no-img-element */}
                  <img
                    src={src}
                    alt={a.codename}
                    className='absolute inset-0 size-full object-cover object-[center_20%]'
                  />
                  <div className='pointer-events-none absolute inset-0 bg-gradient-to-t from-card via-card/20 to-transparent' />
                  <div className='absolute right-3 bottom-3 left-3 flex items-end gap-3'>
                    <Avatar className='size-14 ring-2 ring-primary/40' size='lg'>
                      <AvatarImage src={src} alt={a.codename} />
                      <AvatarFallback className='font-mono text-xs'>{initials}</AvatarFallback>
                    </Avatar>
                    <div className='min-w-0 pb-0.5'>
                      <p className='font-mono text-sm font-semibold tracking-widest text-primary'>{a.codename}</p>
                      <p className='micro-label truncate text-zinc-300'>{a.focus}</p>
                    </div>
                  </div>
                </div>
                <CardHeader className='pt-4'>
                  <div className='flex flex-wrap gap-2'>
                    <Badge className='font-mono text-[10px] tracking-widest'>{a.codename}</Badge>
                    <Badge variant='outline' className='font-mono text-[10px] tracking-widest uppercase'>
                      {a.focus}
                    </Badge>
                  </div>
                  <CardTitle className='mt-2 font-mono text-sm tracking-widest'>{a.name}</CardTitle>
                  <CardDescription className='font-mono text-xs'>{a.description}</CardDescription>
                </CardHeader>
                <CardContent>
                  <p className='micro-label mb-2'>DELEGATES</p>
                  <div className='flex flex-wrap gap-1'>
                    {a.delegates.map(d => (
                      <Badge key={d} variant='secondary' className='font-mono text-[10px]'>
                        {d.replace('delegate_', '')}
                      </Badge>
                    ))}
                  </div>
                  <Link
                    href={`/runs/new?mode=${a.id}`}
                    className={cn(
                      buttonVariants({ variant: 'outline', size: 'sm' }),
                      'mt-4 font-mono text-[10px] tracking-widest uppercase'
                    )}
                  >
                    Launch as {a.codename}
                  </Link>
                </CardContent>
              </Card>
            )
          })}
        </div>
      )}

      <Card className='border-primary/30'>
        <CardHeader>
          <CardTitle className='micro-label'>HOW AGENTS COMMUNICATE</CardTitle>
        </CardHeader>
        <CardContent className='text-muted-foreground space-y-2 font-mono text-xs'>
          <p>
            ▸ <span className='text-primary'>Orchestrator</span> is the only agent that calls others — via{' '}
            <code>delegate_recon</code>, <code>delegate_exploit</code>, <code>delegate_post_exploit</code>,{' '}
            <code>delegate_codegen</code>, <code>delegate_report</code>.
          </p>
          <p>▸ Subagents return text to the orchestrator; they do not peer-call each other.</p>
          <p>
            ▸ <span className='text-primary'>MCP</span>: arsenal (nmap/nuclei/…) + strike (Metasploit) over stdio.
          </p>
          <p>
            ▸ <span className='text-primary'>CyberStrike skills</span>: every subagent has{' '}
            <code>skill_search</code> / <code>skill_get</code> against the ~7.6k pack (see Skills page).
          </p>
          <p>
            ▸ Findings: <code>report_finding</code> / <code>triage_finding</code> with 3-gate evidence.
          </p>
        </CardContent>
      </Card>
    </div>
  )
}

export default AgentsView
