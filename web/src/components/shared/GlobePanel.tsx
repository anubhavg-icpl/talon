'use client'

// Framed operator globe. WebGL is gated behind Lazy3D so hubs don't crash tabs.

import Lazy3D from '@/components/shared/Lazy3D'
import { TalonGlobe } from '@/components/shared/three'
import { cn } from '@/lib/utils'

type GlobePanelProps = {
  className?: string
  variant?: 'hero' | 'compact'
  state?: 'idle' | 'listening' | 'thinking' | 'speaking' | 'running'
  activityLevel?: number
  interactive?: boolean
  label?: string
  /** Use compact "3D" chip gate (thumbnails) */
  gateCompact?: boolean
}

const GlobePanel = ({
  className,
  variant = 'compact',
  state = 'idle',
  activityLevel = 0,
  interactive,
  label,
  gateCompact
}: GlobePanelProps) => {
  const compactGate = gateCompact ?? variant === 'compact'

  return (
    <div
      className={cn(
        'relative aspect-square overflow-hidden rounded-sm border border-primary/20 bg-black/40',
        className
      )}
    >
      <Lazy3D compact={compactGate} className='h-full w-full' label={label || 'OPERATOR GLOBE'}>
        <TalonGlobe
          className='h-full w-full'
          variant={variant}
          state={state}
          activityLevel={activityLevel}
          interactive={interactive}
        />
      </Lazy3D>
      {label ? <span className='micro-label absolute bottom-1.5 left-2 z-10 text-primary/70'>{label}</span> : null}
    </div>
  )
}

export default GlobePanel
