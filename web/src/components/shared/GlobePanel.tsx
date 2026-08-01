'use client'

// Framed operator globe — wallpaper video by default, optional Live 3D.

import GlobeWallpaper from '@/components/shared/GlobeWallpaper'
import { TalonGlobe } from '@/components/shared/three'
import { cn } from '@/lib/utils'

type GlobePanelProps = {
  className?: string
  variant?: 'hero' | 'compact'
  state?: 'idle' | 'listening' | 'thinking' | 'speaking' | 'running'
  activityLevel?: number
  interactive?: boolean
  label?: string
  gateCompact?: boolean
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
    <div className={cn('relative aspect-square w-full', className)}>
      <GlobeWallpaper
        compact
        label={label || 'GLOBE'}
        className='size-full'
        live3d={
          <TalonGlobe
            className='h-full w-full'
            variant={variant}
            state={state}
            activityLevel={activityLevel}
            interactive={interactive}
          />
        }
      />
    </div>
  )
}

export default GlobePanel
