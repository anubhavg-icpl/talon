'use client'

/**
 * Gate heavy WebGL until the operator opts in.
 * Prevents Chrome "This page couldn't load" tab crashes on low-RAM hosts
 * when Overview / Showcase / multi-globe hubs all spin up at once.
 */

import { useCallback, useEffect, useState, type ReactNode } from 'react'
import { Box, Sparkles } from 'lucide-react'

import { Button } from '@/components/ui/button'
import { cn } from '@/lib/utils'

const STORAGE_KEY = 'talon-load-3d'

type Lazy3DProps = {
  children: ReactNode
  className?: string
  /** Poster / still while 3D is off */
  poster?: string
  label?: string
  /** Auto-start if user previously enabled this session */
  remember?: boolean
  /** Compact chip-style gate (sidebar thumbs) */
  compact?: boolean
}

export function usePrefer3D(remember = true): [boolean, (v: boolean) => void] {
  const [on, setOn] = useState(false)

  useEffect(() => {
    if (!remember) return
    try {
      if (sessionStorage.getItem(STORAGE_KEY) === '1') setOn(true)
    } catch {
      /* ignore */
    }
  }, [remember])

  const set = useCallback(
    (v: boolean) => {
      setOn(v)
      if (remember) {
        try {
          if (v) sessionStorage.setItem(STORAGE_KEY, '1')
          else sessionStorage.removeItem(STORAGE_KEY)
        } catch {
          /* ignore */
        }
      }
    },
    [remember]
  )

  return [on, set]
}

const Lazy3D = ({
  children,
  className,
  poster = '/showcase/operator-globe-hud.webp',
  label = 'OPERATOR GLOBE',
  remember = true,
  compact = false
}: Lazy3DProps) => {
  const [on, setOn] = usePrefer3D(remember)

  if (on) {
    return <div className={cn('relative size-full min-h-[80px]', className)}>{children}</div>
  }

  if (compact) {
    return (
      <button
        type='button'
        onClick={() => setOn(true)}
        className={cn(
          'group relative flex size-full min-h-[48px] flex-col items-center justify-center gap-1 overflow-hidden rounded-sm border border-primary/20 bg-black/50',
          className
        )}
        title='Load 3D'
      >
        {poster && (
          // eslint-disable-next-line @next/next/no-img-element
          <img src={poster} alt='' className='absolute inset-0 size-full object-cover opacity-40' />
        )}
        <Box className='text-primary relative size-4' />
        <span className='relative font-mono text-[8px] tracking-widest text-primary/90 uppercase'>3D</span>
      </button>
    )
  }

  return (
    <div
      className={cn(
        'relative flex size-full min-h-[160px] flex-col items-center justify-center gap-3 overflow-hidden rounded-sm border border-primary/20 bg-black/50',
        className
      )}
    >
      {poster && (
        // eslint-disable-next-line @next/next/no-img-element
        <img src={poster} alt='' className='absolute inset-0 size-full object-cover opacity-35' />
      )}
      <div className='pointer-events-none absolute inset-0 bg-gradient-to-t from-background/90 via-background/40 to-transparent' />
      <div className='relative z-10 flex flex-col items-center gap-2 px-4 text-center'>
        <Sparkles className='text-primary size-5' />
        <p className='micro-label text-primary/90'>{label}</p>
        <p className='text-muted-foreground max-w-[220px] font-mono text-[10px] leading-relaxed'>
          3D is off by default to keep this host stable. Load WebGL only when you want it.
        </p>
        <Button
          type='button'
          size='sm'
          className='mt-1 font-mono text-[10px] tracking-widest uppercase'
          onClick={() => setOn(true)}
        >
          Load 3D
        </Button>
      </div>
    </div>
  )
}

export default Lazy3D
