'use client'

import { useCallback, useEffect, useMemo, useState } from 'react'

import Link from 'next/link'

import Lazy3D from '@/components/shared/Lazy3D'
import PageHeader from '@/components/shared/PageHeader'
import { ExamplesStage } from '@/components/shared/three'
import { Badge } from '@/components/ui/badge'
import { Button, buttonVariants } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { cn } from '@/lib/utils'

type ReelItem = {
  id: string
  title: string
  blurb: string
  poster: string
  video?: string
  /** Live interactive Three.js stage instead of still/video */
  live3d?: boolean
  tags: string[]
}

const REEL: ReelItem[] = [
  {
    id: 'live-globe',
    title: 'Operator globe — live Three.js',
    blurb: 'Drag to orbit · starfield · atmosphere · beacons · great-circle arcs · multi-axis rings. Pure WebGL.',
    poster: '/showcase/operator-globe-hud.webp',
    live3d: true,
    tags: ['three.js', 'webgl', 'live']
  },
  {
    id: 'hero',
    title: 'Talon AI — Offensive orchestration',
    blurb: 'Go-native recon → exploit → post-exploit → judge pipeline with HITL gates.',
    poster: '/showcase/talon-hero-raptor.webp',
    video: '/showcase/talon-hero.mp4',
    tags: ['brand', 'hero']
  },
  {
    id: 'dashboard',
    title: 'Operator console',
    blurb: 'Live runs, findings, skills catalog, and multi-agent control plane.',
    poster: '/showcase/talon-dashboard-product.webp',
    video: '/showcase/talon-dashboard.mp4',
    tags: ['ui', 'product']
  },
  {
    id: 'pipeline',
    title: 'Multi-agent pipeline',
    blurb: 'Orchestrator delegates recon, exploit, post-exploit, codegen, and report specialists.',
    poster: '/showcase/talon-pipeline-agents.webp',
    video: '/showcase/talon-pipeline.mp4',
    tags: ['agents', 'a2a']
  },
  {
    id: 'skills',
    title: 'CyberStrike skills pack',
    blurb: '~7.6k methodology skills — browse in UI, load via skill_search / skill_get at runtime.',
    poster: '/showcase/talon-knowledge-panel.webp',
    video: '/showcase/talon-skills.mp4',
    tags: ['skills', 'cyberstrike']
  },
  {
    id: 'globe-still',
    title: 'C2 HUD still',
    blurb: 'Marketing still of the operator globe — same red wire aesthetic as the live scene.',
    poster: '/showcase/operator-globe-hud.webp',
    tags: ['three.js', 'hud']
  },
  {
    id: 'filmstrip',
    title: 'Agent filmstrip',
    blurb: 'Specialist modes: COMMANDER, GHOST, STRIKER, PHANTOM, CIPHER.',
    poster: '/showcase/talon-agent-filmstrip.webp',
    tags: ['agents']
  },
  {
    id: 'mark',
    title: 'Brand mark',
    blurb: 'Talon raptor identity for the operator shell.',
    poster: '/showcase/talon-brand-mark.webp',
    tags: ['brand']
  }
]

const ShowcaseView = () => {
  const [active, setActive] = useState(0)
  const [playing, setPlaying] = useState(true)
  const [videoOk, setVideoOk] = useState<Record<string, boolean>>({})

  const item = REEL[active]

  useEffect(() => {
    if (!playing) return
    if (item.live3d) return // keep live 3d on stage while selected; manual next
    const hasVideo = videoOk[item.id]
    if (hasVideo) return
    const t = setTimeout(() => setActive(i => (i + 1) % REEL.length), 4500)
    return () => clearTimeout(t)
  }, [active, playing, item.id, item.live3d, videoOk])

  const onVideoMeta = useCallback((id: string, ok: boolean) => {
    setVideoOk(prev => (prev[id] === ok ? prev : { ...prev, [id]: ok }))
  }, [])

  const tags = useMemo(() => [...new Set(REEL.flatMap(r => r.tags))], [])

  return (
    <div className='flex flex-col gap-6'>
      <PageHeader
        title='SHOWCASE'
        description='Live operator globe · SkeletonUtils + GLTF · product reel — the Three.js surface, end to end'
        actions={
          <>
            <Link href='/runs/new' className={cn(buttonVariants(), 'font-mono text-xs tracking-widest uppercase')}>
              Launch operation
            </Link>
            <Link
              href='/skills'
              className={cn(buttonVariants({ variant: 'outline' }), 'font-mono text-xs tracking-widest uppercase')}
            >
              Browse skills
            </Link>
          </>
        }
      />

      {/* Single primary WebGL stage — gated so the page always loads first */}
      <Lazy3D className='min-h-[360px] w-full' label='SHOWCASE THREE.JS STAGE' poster='/showcase/operator-globe-hud.webp'>
        <ExamplesStage className='min-h-[360px] w-full' />
      </Lazy3D>

      {/* Product still/video reel — no extra WebGL */}
      <div className='grid gap-4 lg:grid-cols-[1fr_220px]'>
        <Card className='hud-corners relative overflow-hidden'>
          <CardContent className='p-0'>
            <div className='relative aspect-video w-full bg-black'>
              {videoOk[item.id] && item.video && !item.live3d ? (
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
                  <img
                    key={item.poster}
                    src={item.poster}
                    alt={item.title}
                    className='showcase-kenburns h-full w-full object-cover'
                  />
                  {item.video && !item.live3d && videoOk[item.id] === undefined && (
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
                  {item.live3d ? (
                    <Badge className='font-mono text-[10px]'>SEE STAGE ABOVE</Badge>
                  ) : videoOk[item.id] ? (
                    <Badge className='font-mono text-[10px]'>VIDEO</Badge>
                  ) : (
                    <Badge variant='secondary' className='font-mono text-[10px]'>
                      STILL
                    </Badge>
                  )}
                </div>
                <h2 className='mt-2 font-mono text-lg font-semibold tracking-wide'>{item.title}</h2>
                <p className='text-muted-foreground mt-1 max-w-2xl font-mono text-xs'>{item.blurb}</p>
              </div>
            </div>
          </CardContent>
        </Card>

        <div className='flex flex-col gap-3'>
          <Card className='hud-corners flex-1'>
            <CardHeader className='pb-2'>
              <CardTitle className='micro-label'>REEL CONTROLS</CardTitle>
              <CardDescription className='font-mono text-[11px]'>Stills & optional MP4s</CardDescription>
            </CardHeader>
            <CardContent className='flex flex-col gap-2'>
              <Button
                variant='outline'
                size='sm'
                className='font-mono text-[10px] tracking-widest uppercase'
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
                Next slide
              </Button>
              <Link
                href='/runs/new'
                className={cn(buttonVariants({ size: 'sm' }), 'font-mono text-[10px] tracking-widest uppercase')}
              >
                Engage
              </Link>
            </CardContent>
          </Card>
        </div>
      </div>

      <div className='grid grid-cols-2 gap-3 sm:grid-cols-3 lg:grid-cols-4 xl:grid-cols-8'>
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

      <div className='grid gap-4 md:grid-cols-3'>
        <Card>
          <CardHeader>
            <CardTitle className='micro-label'>THREE.JS STACK</CardTitle>
          </CardHeader>
          <CardContent className='text-muted-foreground font-mono text-xs leading-relaxed'>
            WebGLRenderer · MeshStandardMaterial · Points · Line / LineLoop · OrbitControls · fog · lights ·
            Fibonacci beacons · great-circle arcs · starfield
          </CardContent>
        </Card>
        <Card>
          <CardHeader>
            <CardTitle className='micro-label'>WHERE IT RUNS</CardTitle>
          </CardHeader>
          <CardContent className='text-muted-foreground font-mono text-xs leading-relaxed'>
            Login (background globe) · all pages (ambient starfield) · Overview · New run · Engagements · Findings ·
            Skills · Runs — every hub via the shared &lt;GlobePanel&gt;
          </CardContent>
        </Card>
        <Card>
          <CardHeader>
            <CardTitle className='micro-label'>OPERATOR SURFACES</CardTitle>
          </CardHeader>
          <CardContent className='text-muted-foreground font-mono text-xs leading-relaxed'>
            Dashboard · CLI · live SSE/WS · agent modes · 3-gate findings · CyberStrike skills
          </CardContent>
        </Card>
      </div>
    </div>
  )
}

export default ShowcaseView
