'use client'

import { useEffect, useState } from 'react'

import Link from 'next/link'

import type { Playbook } from '@/lib/api'
import HudStill from '@/components/shared/HudStill'
import { Badge } from '@/components/ui/badge'
import { buttonVariants } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Skeleton } from '@/components/ui/skeleton'
import { getPlaybooks } from '@/lib/api'
import { cn } from '@/lib/utils'

const PlaybooksView = () => {
  const [playbooks, setPlaybooks] = useState<Playbook[] | null>(null)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    getPlaybooks()
      .then(r => setPlaybooks(r.playbooks ?? []))
      .catch(e => setError(e instanceof Error ? e.message : String(e)))
  }, [])

  return (
    <div className='flex flex-col gap-6'>
      <div>
        <h1 className='font-mono text-xl font-semibold tracking-widest'>PLAYBOOKS</h1>
        <p className='micro-label mt-1'>ENGAGEMENT TEMPLATES — LAUNCH WITH PRESET AGENT MODE + PROMPT</p>
      </div>

      <HudStill
        src='/showcase/talon-pipeline-agents.webp'
        alt='Engagement pipeline templates'
        variant='banner'
        className='max-h-32'
      />

      {error && (
        <div className='border-destructive/40 bg-destructive/10 text-destructive rounded-md border px-4 py-3 font-mono text-xs'>
          {error}
        </div>
      )}
      {!playbooks ? (
        <Skeleton className='h-40 w-full' />
      ) : (
        <div className='grid gap-4 md:grid-cols-2 xl:grid-cols-3'>
          {playbooks.map(pb => (
            <Card key={pb.id} className='hud-corners'>
              <CardHeader>
                <div className='flex flex-wrap gap-2'>
                  <Badge className='font-mono text-[10px]'>{pb.codename}</Badge>
                  <Badge variant='outline' className='font-mono text-[10px] uppercase'>
                    {pb.agent_mode}
                  </Badge>
                </div>
                <CardTitle className='mt-2 font-mono text-sm tracking-widest'>{pb.name}</CardTitle>
                <CardDescription className='font-mono text-xs'>{pb.description}</CardDescription>
              </CardHeader>
              <CardContent className='space-y-3'>
                <p className='text-muted-foreground font-mono text-[11px] leading-relaxed'>{pb.prompt}</p>
                <div className='flex flex-wrap gap-1'>
                  {pb.tags?.map(t => (
                    <Badge key={t} variant='secondary' className='font-mono text-[9px]'>
                      {t}
                    </Badge>
                  ))}
                </div>
                <Link
                  href={`/runs/new?playbook=${pb.id}&mode=${pb.agent_mode}`}
                  className={cn(buttonVariants({ size: 'sm' }), 'font-mono text-[10px] tracking-widest uppercase')}
                >
                  Launch playbook
                </Link>
              </CardContent>
            </Card>
          ))}
        </div>
      )}
    </div>
  )
}

export default PlaybooksView
