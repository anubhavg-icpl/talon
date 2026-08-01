'use client'

/**
 * Production Target HUD: wallpaper video (dark planet) as the globe surface.
 * Always works — no WebGL. Optional "Live 3D" mounts TalonGlobe on demand.
 */

import { useEffect, useRef, useState, type ReactNode } from 'react'
import { Box, Sparkles } from 'lucide-react'

import { Button } from '@/components/ui/button'
import { cn } from '@/lib/utils'

/** Prefer WebM (VP9); MP4 kept as fallback for older engines */
const VIDEO_SRC_WEBM = '/globe/dark-planet.webm'
const VIDEO_SRC_MP4 = '/globe/dark-planet.mp4'
const POSTER_SRC = 'data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAIAAACQd1PeAAAADElEQVR4nGNgYGAAAAAEAAH2FzhVAAAAAElFTkSuQmCC'

type GlobeWallpaperProps = {
  className?: string
  label?: string
  /** Optional WebGL globe shown after user opts in */
  live3d?: ReactNode
  /** Show the Live 3D control (default true if live3d provided) */
  allow3d?: boolean
  /** Compact square HUD (Target HUD) */
  compact?: boolean
}

const GlobeWallpaper = ({
  className,
  label = 'TARGET HUD',
  live3d,
  allow3d = Boolean(live3d),
  compact = false
}: GlobeWallpaperProps) => {
  const videoRef = useRef<HTMLVideoElement | null>(null)
  const [mode, setMode] = useState<'video' | '3d'>('video')
  const [videoOk, setVideoOk] = useState(true)

  useEffect(() => {
    const v = videoRef.current
    if (!v || mode !== 'video') return
    v.play().catch(() => {
      /* autoplay may be blocked — poster still shows */
    })
  }, [mode, videoOk])

  return (
    <div
      className={cn(
        'relative overflow-hidden rounded-sm border border-primary/25 bg-black',
        compact ? 'aspect-square w-full' : 'aspect-video w-full',
        className
      )}
    >
      {mode === 'video' || !live3d ? (
        <>
          {videoOk ? (
            <video
              ref={videoRef}
              className='absolute inset-0 size-full object-cover'
              poster={POSTER_SRC}
              autoPlay
              muted
              loop
              playsInline
              preload='metadata'
              onError={() => setVideoOk(false)}
            >
              <source src={VIDEO_SRC_WEBM} type='video/webm' />
              <source src={VIDEO_SRC_MP4} type='video/mp4' />
            </video>
          ) : (
            // eslint-disable-next-line @next/next/no-img-element
            <img src={POSTER_SRC} alt='' className='absolute inset-0 size-full object-cover' />
          )}
          <div className='pointer-events-none absolute inset-0 bg-gradient-to-t from-black/70 via-transparent to-black/30' />
          <div className='pointer-events-none absolute inset-0 ring-1 ring-inset ring-primary/15' />
          {/* HUD chrome */}
          <div className='absolute top-2 left-2 z-10 flex items-center gap-1.5'>
            <span className='size-1.5 animate-pulse rounded-full bg-primary' />
            <span className='font-mono text-[9px] tracking-[0.2em] text-primary/90 uppercase'>{label}</span>
          </div>
          <div className='absolute right-2 bottom-2 left-2 z-10 flex items-end justify-between gap-2'>
            <span className='font-mono text-[8px] tracking-widest text-zinc-400 uppercase'>
              {videoOk ? 'WALLPAPER · DARK PLANET' : 'STILL · HUD'}
            </span>
            {allow3d && live3d && (
              <Button
                type='button'
                size='sm'
                variant='outline'
                className='pointer-events-auto h-7 border-primary/40 bg-black/50 font-mono text-[9px] tracking-widest uppercase backdrop-blur-sm'
                onClick={() => setMode('3d')}
              >
                <Box className='mr-1 size-3' />
                Live 3D
              </Button>
            )}
          </div>
        </>
      ) : (
        <>
          <div className='absolute inset-0'>{live3d}</div>
          <div className='absolute top-2 right-2 z-10'>
            <Button
              type='button'
              size='sm'
              variant='outline'
              className='h-7 border-primary/40 bg-black/50 font-mono text-[9px] tracking-widest uppercase backdrop-blur-sm'
              onClick={() => setMode('video')}
            >
              <Sparkles className='mr-1 size-3' />
              Wallpaper
            </Button>
          </div>
        </>
      )}
    </div>
  )
}

export default GlobeWallpaper
