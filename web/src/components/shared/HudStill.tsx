'use client'

/**
 * Lightweight product still for page chrome — not a showcase reel.
 * Use where media already belongs (agents strip, skills banner, etc.).
 */

import { cn } from '@/lib/utils'

type HudStillProps = {
  src: string
  alt?: string
  className?: string
  /** Cover height for banner strips */
  variant?: 'banner' | 'panel' | 'filmstrip'
}

const HudStill = ({ src, alt = '', className, variant = 'banner' }: HudStillProps) => {
  return (
    <div
      className={cn(
        'relative overflow-hidden rounded-sm border border-primary/15 bg-black',
        variant === 'banner' && 'aspect-[21/6] w-full max-h-36',
        variant === 'panel' && 'aspect-video w-full',
        variant === 'filmstrip' && 'aspect-[32/9] w-full max-h-28',
        className
      )}
    >
      {/* eslint-disable-next-line @next/next/no-img-element */}
      <img src={src} alt={alt} className='absolute inset-0 size-full object-cover object-center opacity-90' />
      <div className='pointer-events-none absolute inset-0 bg-gradient-to-r from-background/70 via-transparent to-background/40' />
      <div className='pointer-events-none absolute inset-0 ring-1 ring-inset ring-primary/10' />
    </div>
  )
}

export default HudStill
