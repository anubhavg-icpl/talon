'use client'

/**
 * Production entry for Three.js surfaces — client-only dynamic imports
 * so the main bundle stays lean and SSR never touches WebGL.
 */

import dynamic from 'next/dynamic'

import WebGLFallback from '@/components/shared/WebGLFallback'

export const TalonGlobe = dynamic(() => import('@/components/shared/TalonGlobe'), {
  ssr: false,
  loading: () => <WebGLFallback label='LOADING GLOBE…' />
})

export const Starfield = dynamic(() => import('@/components/shared/three/Starfield'), {
  ssr: false,
  loading: () => null
})

export const ExamplesStage = dynamic(() => import('@/components/shared/three/ExamplesStage'), {
  ssr: false,
  loading: () => <WebGLFallback className='min-h-[320px]' label='LOADING THREE.JS STAGE…' />
})

export const SkeletonDemo = dynamic(() => import('@/components/shared/three/SkeletonDemo'), {
  ssr: false,
  loading: () => <WebGLFallback className='min-h-[280px]' label='LOADING SKELETON SCENE…' />
})
