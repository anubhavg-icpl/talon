'use client'

import { useCallback, useEffect, useMemo, useState } from 'react'

import Link from 'next/link'

import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import TalonGlobe from '@/components/shared/TalonGlobe'

type ReelItem = {
  id: string
  title: string
  blurb: string
  /** Poster / still (always works) */
  poster: string
  /** Optional product video path under /showcase/*.mp4 */
  video?: string
  tags: string[]
}

const REEL: ReelItem[] = [
  {
    id: 'hero',
    title: 'Talon AI — Offensive orchestration',
    blurb: 'Go-native recon → exploit → post-exploit → judge pipeline with HITL gates.',
    poster: '/showcase/talon-hero-raptor.jpg',
    video: '/showcase/talon-hero.mp4',
    tags: ['brand', 'hero']
  },
  {
    id: 'dashboard',
    title: 'Ops console',
    blurb: 'Live runs, findings, skills catalog, and multi-agent control plane.',
    poster: '/showcase/talon-dashboard-product.jpg',
    video: '/showcase/talon-dashboard.mp4',
    tags: ['ui', 'product']
  },
  {
    id: 'pipeline',
    title: 'Multi-agent pipeline',
    blurb: 'Orchestrator delegates recon, exploit, post-exploit, codegen, and report specialists.',
    poster: '/showcase/talon-pipeline-agents.jpg',
    video: '/showcase/talon-pipeline.mp4',
    tags: ['agents', 'a2a']
  },
  {
    id: 'skills',
    title: 'CyberStrike skills pack',
    blurb: '~7.6k methodology skills — browse in UI, load via skill_search / skill_get at runtime.',
    poster: '/showcase/talon-knowledge-panel.jpg',
    video: '/showcase/talon-skills.mp4',
    tags: ['skills', 'cyberstrike']
  },
  {
    id: 'filmstrip',
    title: 'Agent filmstrip',
    blurb: 'Specialist modes: COMMANDER, GHOST, STRIKER, PHANTOM, CIPHER.',
    poster: '/showcase/talon-agent-filmstrip.jpg',
    tags: ['agents']
  },
  {
    id: 'mark',
    title: 'Brand mark',
    blurb: 'Talon raptor identity for the red/black ops console.',
    poster: '/showcase/talon-brand-mark.jpg',
    tags: ['brand']
  }
]

const ShowcaseView = () => {
  const [active, setActive] = useState(0)
  const [playing, setPlaying] = useState(true)
  const [videoOk, setVideoOk] = useState<Record<string, boolean>>({})

  const item = REEL[active]

  // Auto-advance reel when no video is playing through
  useEffect(() => {
    if (!playing) return
    const hasVideo = videoOk[item.id]
    if (hasVideo) return // let video onEnded advance
    const t = setTimeout(() => setActive(i => (i + 1) % REEL.length), 4500)
    return () => clearTimeout(t)
  }, [active, playing, item.id, videoOk])

  const onVideoMeta = useCallback((id: string, ok: boolean) => {
    setVideoOk(prev => (prev[id] === ok ? prev : { ...prev, [id]: ok }))
  }, [])

  const tags = useMemo(() => [...new Set(REEL.flatMap(r => r.tags))], [])

  return (
    <div className='flex flex-col gap-6'>
      <div className='flex flex-wrap items-end justify-between gap-3'>
        <div>
          <h1 className='font-mono text-xl font-semibold tracking-widest'>SHOWCASE</h1>
          <p className='micro-label mt-1'>PRODUCT REEL · GLOBE · PIPELINE · SKILLS — ALIGNED E2E DEMO SURFACE</p>
        </div>
        <div className='flex flex-wrap gap-2'>
          <Button asChild className='font-mono text-xs tracking-widest uppercase'>
            <Link href='/runs/new'>Launch operation</Link>
          </Button>
          <Button asChild variant='outline' className='font-mono text-xs tracking-widest uppercase'>
            <Link href='/skills'>Browse skills</Link>
          </Button>
        </div>
      </div>

      {/* Main stage */}
      <div className='grid gap-4 lg:grid-cols-[1fr_280px]'>
        <Card className='hud-corners scanlines relative overflow-hidden'>
          <CardContent className='p-0'>
            <div className='relative aspect-video w-full bg-black'>
              {videoOk[item.id] ? (
                <video
                  key={item.id}
                  className='h-full w-full object-cover'
                  poster={item.poster}
                  src={item.video}
                  autoPlay
                  muted
                  playsInline
                  onEnded={() => playing && setActive(i => (i + 1) % REEL.length)}
                  onError={() => onVideoMeta(item.id, false)}
                  onLoadedData={() => onVideoMeta(item.id, true)}
                />
              ) : (
                <>
                  {/* Cinematic still with slow ken-burns via CSS */}
                  <img
                    key={item.poster}
                    src={item.poster}
                    alt={item.title}
                    className='showcase-kenburns h-full w-full object-cover'
                  />
                  {/* Probe optional mp4 once */}
                  {item.video && videoOk[item.id] === undefined && (
                    <video
                      className='hidden'
                      src={item.video}
                      preload='metadata'
                      onLoadedData={() => onVideoMeta(item.id, true)}
                      onError={() => onVideoMeta(item.id, false)}
                    />
                  )}
                </>
              )}
              <div className='absolute inset-x-0 bottom-0 bg-gradient-to-t from-black/90 via-black/50 to-transparent p-6'>
                <div className='flex flex-wrap gap-2'>
                  {item.tags.map(t => (
                    <Badge key={t} variant='outline' className='font-mono text-[10px] uppercase'>
                      {t}
                    </Badge>
                  ))}
                  {videoOk[item.id] ? (
                    <Badge className='font-mono text-[10px]'>VIDEO</Badge>
                  ) : (
                    <Badge variant='secondary' className='font-mono text-[10px]'>
                      STILL REEL
                    </Badge>
                  )}
                </div>
                <h2 className='mt-2 font-mono text-lg font-semibold tracking-wide'>{item.title}</h2>
                <p className='text-muted-foreground mt-1 max-w-2xl font-mono text-xs'>{item.blurb}</p>
              </div>
            </div>
          </CardContent>
        </Card>

        <div className='flex flex-col gap-4'>
          <Card className='hud-corners flex flex-1 flex-col overflow-hidden'>
            <CardHeader className='pb-2'>
              <CardTitle className='micro-label'>OPS GLOBE</CardTitle>
              <CardDescription className='font-mono text-[11px]'>
                From agentic-os HUD · click to start a run
              </CardDescription>
            </CardHeader>
            <CardContent className='flex flex-1 items-center justify-center p-2'>
              <TalonGlobe
                className='aspect-square w-full max-w-[240px]'
                state='running'
                activityLevel={0.45}
                onClick={() => {
                  window.location.href = '/runs/new'
                }}
              />
            </CardContent>
          </Card>
          <div className='flex gap-2'>
            <Button
              variant='outline'
              size='sm'
              className='flex-1 font-mono text-[10px] tracking-widest uppercase'
              onClick={() => setPlaying(p => !p)}
            >
              {playing ? 'Pause reel' : 'Play reel'}
            </Button>
            <Button
              variant='outline'
              size='sm'
              className='font-mono text-[10px] tracking-widest uppercase'
              onClick={() => setActive(i => (i + 1) % REEL.length)}
            >
              Next
            </Button>
          </div>
        </div>
      </div>

      {/* Thumbnails */}
      <div className='grid grid-cols-2 gap-3 sm:grid-cols-3 lg:grid-cols-6'>
        {REEL.map((r, i) => (
          <button
            key={r.id}
            type='button'
            onClick={() => setActive(i)}
            className={`overflow-hidden rounded-md border text-left transition-colors ${
              i === active ? 'border-primary ring-primary/40 ring-1' : 'border-border/60 hover:border-primary/40'
            }`}
          >
            <img src={r.poster} alt={r.title} className='aspect-video w-full object-cover' />
            <p className='truncate px-2 py-1.5 font-mono text-[10px] tracking-wide'>{r.title}</p>
          </button>
        ))}
      </div>

      {/* Capability grid */}
      <div className='grid gap-4 md:grid-cols-3'>
        <Card>
          <CardHeader>
            <CardTitle className='micro-label'>PIPELINE</CardTitle>
          </CardHeader>
          <CardContent className='text-muted-foreground font-mono text-xs leading-relaxed'>
            recon → exploit → post-exploit → codegen → judge · HITL on nmap · MCP arsenal + Metasploit strike
          </CardContent>
        </Card>
        <Card>
          <CardHeader>
            <CardTitle className='micro-label'>INTELLIGENCE</CardTitle>
          </CardHeader>
          <CardContent className='text-muted-foreground font-mono text-xs leading-relaxed'>
            CyberStrike skill pack · skill_search / skill_get · 3-gate findings · kill chain · methodology coverage
          </CardContent>
        </Card>
        <Card>
          <CardHeader>
            <CardTitle className='micro-label'>OPERATOR SURFACES</CardTitle>
          </CardHeader>
          <CardContent className='text-muted-foreground font-mono text-xs leading-relaxed'>
            Dashboard · CLI · live SSE/WS · globe ops hero · agent modes (full / web / network / post)
          </CardContent>
        </Card>
      </div>

      <Card className='border-primary/20'>
        <CardHeader>
          <CardTitle className='micro-label'>ADD REAL MP4 SHOWCASE VIDEOS</CardTitle>
          <CardDescription className='font-mono text-xs'>
            Drop files into <code className='text-primary'>web/public/showcase/</code> — auto-detected:
          </CardDescription>
        </CardHeader>
        <CardContent>
          <ul className='text-muted-foreground space-y-1 font-mono text-[11px]'>
            <li>▸ talon-hero.mp4</li>
            <li>▸ talon-dashboard.mp4</li>
            <li>▸ talon-pipeline.mp4</li>
            <li>▸ talon-skills.mp4</li>
          </ul>
          <p className='micro-label mt-3'>Until then the reel uses cinematic product stills with Ken Burns motion.</p>
          <div className='mt-3 flex flex-wrap gap-1'>
            {tags.map(t => (
              <Badge key={t} variant='outline' className='font-mono text-[9px] uppercase'>
                {t}
              </Badge>
            ))}
          </div>
        </CardContent>
      </Card>
    </div>
  )
}

export default ShowcaseView
