'use client'

/**
 * Optional Three.js stage (beacons / arcs / rings) — used only if Live 3D is opted in:
 * 1) Operator globe (procedural HUD)
 * 2) SkeletonUtils multi-instance skinned mesh (official Soldier.glb)
 */

import { useState } from 'react'

import TalonGlobe from '@/components/shared/TalonGlobe'
import SkeletonDemo, { type SkeletonMode } from '@/components/shared/three/SkeletonDemo'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { cn } from '@/lib/utils'

export type ExampleId = 'globe' | 'skeleton'

const ExamplesStage = ({ className }: { className?: string }) => {
  const [example, setExample] = useState<ExampleId>('globe')
  const [skeletonMode, setSkeletonMode] = useState<SkeletonMode>('independent')

  return (
    <Card className={cn('hud-corners overflow-hidden border-primary/30', className)}>
      <CardHeader className='pb-2'>
        <div className='flex flex-wrap items-start justify-between gap-3'>
          <div>
            <CardTitle className='micro-label text-primary'>OPERATOR GLOBE · LIVE WEBGL STAGE</CardTitle>
            <CardDescription className='font-mono text-[11px]'>
              Procedural C2 HUD · OrbitControls · starfield · great-circle arcs · SkeletonUtils + GLTF
            </CardDescription>
          </div>
          <div className='flex flex-wrap gap-2'>
            <Button
              size='sm'
              variant={example === 'globe' ? 'default' : 'outline'}
              className='font-mono text-[10px] tracking-widest uppercase'
              onClick={() => setExample('globe')}
            >
              Operator globe
            </Button>
            <Button
              size='sm'
              variant={example === 'skeleton' ? 'default' : 'outline'}
              className='font-mono text-[10px] tracking-widest uppercase'
              onClick={() => setExample('skeleton')}
            >
              SkeletonUtils
            </Button>
          </div>
        </div>
      </CardHeader>
      <CardContent className='space-y-3 p-0 pt-0'>
        {example === 'skeleton' && (
          <div className='flex flex-wrap gap-2 px-6'>
            <Button
              size='sm'
              variant={skeletonMode === 'independent' ? 'default' : 'outline'}
              className='font-mono text-[10px] tracking-widest uppercase'
              onClick={() => setSkeletonMode('independent')}
            >
              Independent skeletons
            </Button>
            <Button
              size='sm'
              variant={skeletonMode === 'shared' ? 'default' : 'outline'}
              className='font-mono text-[10px] tracking-widest uppercase'
              onClick={() => setSkeletonMode('shared')}
            >
              Shared skeleton
            </Button>
            <Badge variant='outline' className='font-mono text-[9px]'>
              Soldier.glb · three.js examples
            </Badge>
          </div>
        )}

        <div className='relative h-[min(56vh,480px)] w-full bg-black'>
          {example === 'skeleton' ? (
            <SkeletonDemo key={skeletonMode} mode={skeletonMode} className='h-full w-full' />
          ) : (
            <TalonGlobe className='h-full w-full' variant='hero' state='running' activityLevel={0.7} interactive />
          )}
        </div>

        <div className='text-muted-foreground space-y-1 px-6 pb-4 font-mono text-[10px] leading-relaxed'>
          {example === 'skeleton' ? (
            skeletonMode === 'independent' ? (
              <p>
                <span className='text-primary'>setupDefaultScene</span> — three{' '}
                <code className='text-primary'>SkeletonUtils.clone</code> models, each with its own{' '}
                <code className='text-primary'>AnimationMixer</code> (idle · walk · run).
              </p>
            ) : (
              <p>
                <span className='text-primary'>setupSharedSkeletonScene</span> — one bone tree, three meshes in{' '}
                <code className='text-primary'>DetachedBindMode</code> bound to the shared skeleton (same animation
                state).
              </p>
            )
          ) : (
            <p>
              Procedural operator HUD: starfield, atmosphere, beacons, great-circle arcs, multi-axis rings,
              OrbitControls.
            </p>
          )}
        </div>
      </CardContent>
    </Card>
  )
}

export default ExamplesStage
