'use client'

import { useEffect, useState } from 'react'

import Link from 'next/link'

import type { AgentInfo } from '@/lib/api'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Skeleton } from '@/components/ui/skeleton'
import { getAgents } from '@/lib/api'

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
        <Button asChild className='font-mono text-xs tracking-widest uppercase'>
          <Link href='/runs/new'>New Operation</Link>
        </Button>
      </div>

      {error && (
        <div className='border-destructive/40 bg-destructive/10 text-destructive rounded-md border px-4 py-3 font-mono text-xs'>
          {error}
        </div>
      )}

      {!agents ? (
        <Skeleton className='h-40 w-full' />
      ) : (
        <div className='grid gap-4 md:grid-cols-2 xl:grid-cols-3'>
          {agents.map(a => (
            <Card key={a.id} className='hud-corners'>
              <CardHeader>
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
                <Button asChild variant='outline' size='sm' className='mt-4 font-mono text-[10px] tracking-widest uppercase'>
                  <Link href={`/runs/new?mode=${a.id}`}>Launch as {a.codename}</Link>
                </Button>
              </CardContent>
            </Card>
          ))}
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
