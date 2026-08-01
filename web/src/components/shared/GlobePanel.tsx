'use client'

// Distributes the operator Three.js globe across views with one consistent
// framed treatment so every hub feels like the same C2 console. All props
// optional — defaults to a static compact globe (no live-run state coupling).

import { cn } from '@/lib/utils'
import { TalonGlobe } from '@/components/shared/three'

type GlobePanelProps = {
  className?: string
  variant?: 'hero' | 'compact'
  state?: 'idle' | 'listening' | 'thinking' | 'speaking' | 'running'
  activityLevel?: number
  interactive?: boolean
  label?: string
}

const GlobePanel = ({
  className,
  variant = 'compact',
  state = 'idle',
  activityLevel = 0,
  interactive,
  label
}: GlobePanelProps) => {
  return (
    <div
      className={cn(
        'relative aspect-square overflow-hidden rounded-sm border border-primary/20 bg-black/40',
        className
      )}
    >
      <TalonGlobe
        className='h-full w-full'
        variant={variant}
        state={state}
        activityLevel={activityLevel}
        interactive={interactive}
      />
      {label ? <span className='micro-label absolute bottom-1.5 left-2 text-primary/70'>{label}</span> : null}
    </div>
  )
}

export default GlobePanel
