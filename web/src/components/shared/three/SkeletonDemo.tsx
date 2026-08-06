'use client'

/**
 * Port of the Three.js SkeletonUtils example the user provided:
 * - GLTF load (Soldier.glb from three.js examples)
 * - SkeletonUtils.clone for independent skinned instances
 * - Per-clone AnimationMixer (idle / walk / run)
 * - Optional shared-skeleton mode (DetachedBindMode)
 * - OrbitControls + red operator lighting
 *
 * Model: web/public/showcase/models/Soldier.glb
 * (https://threejs.org/examples/models/gltf/Soldier.glb)
 */

import { useEffect, useRef, useState } from 'react'
import * as THREE from 'three'
import { OrbitControls } from 'three/examples/jsm/controls/OrbitControls.js'
import { GLTFLoader } from 'three/examples/jsm/loaders/GLTFLoader.js'
import { clone as skeletonClone } from 'three/examples/jsm/utils/SkeletonUtils.js'

import { cn } from '@/lib/utils'

export type SkeletonMode = 'independent' | 'shared'

type SkeletonDemoProps = {
  className?: string
  mode?: SkeletonMode
  modelUrl?: string
}

const MODEL_URL = '/showcase/models/Soldier.glb'

const SkeletonDemo = ({
  className,
  mode = 'independent',
  modelUrl = MODEL_URL
}: SkeletonDemoProps) => {
  const wrapRef = useRef<HTMLDivElement>(null)
  const canvasRef = useRef<HTMLCanvasElement>(null)
  const modeRef = useRef(mode)
  const [status, setStatus] = useState<'loading' | 'ready' | 'error'>('loading')
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    modeRef.current = mode
  }, [mode])

  useEffect(() => {
    const wrap = wrapRef.current
    const canvas = canvasRef.current
    if (!wrap || !canvas) return

    let disposed = false
    const prefersReduced = window.matchMedia('(prefers-reduced-motion: reduce)').matches

    const scene = new THREE.Scene()
    scene.background = new THREE.Color(0x060a0e)
    scene.fog = new THREE.Fog(0x060a0e, 8, 28)

    const camera = new THREE.PerspectiveCamera(45, 1, 0.1, 100)
    camera.position.set(2.5, 2.2, 5.5)

    const renderer = new THREE.WebGLRenderer({ canvas, antialias: true, alpha: false })
    renderer.setPixelRatio(Math.min(window.devicePixelRatio || 1, 2))
    renderer.outputColorSpace = THREE.SRGBColorSpace
    renderer.shadowMap.enabled = true

    // Operator red key + fill (red-only theme)
    const hemi = new THREE.HemisphereLight(0xfca5a5, 0x0c1220, 1.1)
    scene.add(hemi)
    const dir = new THREE.DirectionalLight(0xef4444, 2.2)
    dir.position.set(4, 8, 3)
    dir.castShadow = true
    dir.shadow.mapSize.set(1024, 1024)
    scene.add(dir)
    const fill = new THREE.DirectionalLight(0xf87171, 0.6)
    fill.position.set(-3, 2, -2)
    scene.add(fill)

    // Ground grid — HUD floor
    const ground = new THREE.Mesh(
      new THREE.CircleGeometry(12, 64),
      new THREE.MeshStandardMaterial({
        color: 0x140505,
        roughness: 0.9,
        metalness: 0.2,
        transparent: true,
        opacity: 0.95
      })
    )
    ground.rotation.x = -Math.PI / 2
    ground.receiveShadow = true
    scene.add(ground)

    const grid = new THREE.GridHelper(16, 32, 0xef4444, 0x4c1414)
    grid.position.y = 0.01
    ;(grid.material as THREE.Material).transparent = true
    ;(grid.material as THREE.Material).opacity = 0.35
    scene.add(grid)

    // Ring reticle under actors
    const ring = new THREE.Mesh(
      new THREE.RingGeometry(3.2, 3.28, 96),
      new THREE.MeshBasicMaterial({
        color: 0xef4444,
        transparent: true,
        opacity: 0.45,
        side: THREE.DoubleSide
      })
    )
    ring.rotation.x = -Math.PI / 2
    ring.position.y = 0.02
    scene.add(ring)

    const controls = new OrbitControls(camera, canvas)
    controls.enableDamping = true
    controls.target.set(0, 1, 0)
    controls.minDistance = 2
    controls.maxDistance = 14
    controls.maxPolarAngle = Math.PI * 0.48
    controls.update()

    const clock = new THREE.Clock()
    const objects: THREE.Object3D[] = []
    const mixers: THREE.AnimationMixer[] = []

    const clearActors = () => {
      for (const o of objects) {
        scene.remove(o)
        o.traverse(child => {
          if (child instanceof THREE.Mesh) {
            child.geometry?.dispose()
            const mats = Array.isArray(child.material) ? child.material : [child.material]
            for (const m of mats) m?.dispose?.()
          }
        })
      }
      objects.length = 0
      mixers.length = 0
    }

    const setupIndependent = (model: THREE.Object3D, animations: THREE.AnimationClip[]) => {
      // three clones — each has its own skeleton + mixer (user example setupDefaultScene)
      const clips = {
        idle: animations.find(c => /idle/i.test(c.name)) ?? animations[0],
        walk: animations.find(c => /walk/i.test(c.name)) ?? animations[Math.min(1, animations.length - 1)],
        run: animations.find(c => /run/i.test(c.name)) ?? animations[Math.min(2, animations.length - 1)]
      }

      const placements = [
        { x: -2, clip: clips.idle },
        { x: 0, clip: clips.walk },
        { x: 2, clip: clips.run }
      ]

      for (const p of placements) {
        const clone = skeletonClone(model)
        clone.position.set(p.x, 0, 0)
        clone.traverse(c => {
          if (c instanceof THREE.Mesh) {
            c.castShadow = true
            c.receiveShadow = true
          }
        })
        scene.add(clone)
        objects.push(clone)
        const mixer = new THREE.AnimationMixer(clone)
        if (p.clip) mixer.clipAction(p.clip).play()
        mixers.push(mixer)
      }
    }

    const setupShared = (model: THREE.Object3D, animations: THREE.AnimationClip[]) => {
      // shared skeleton — all meshes drive from one bone hierarchy (user example setupSharedSkeletonScene)
      const sharedModel = skeletonClone(model)
      const skinnedMeshes: THREE.SkinnedMesh[] = []
      let hips: THREE.Object3D | null = null

      sharedModel.traverse(c => {
        if ((c as THREE.SkinnedMesh).isSkinnedMesh) skinnedMeshes.push(c as THREE.SkinnedMesh)
        if (c.name.toLowerCase().includes('hip') || c.name === 'mixamorigHips') hips = c
      })

      const skinned = skinnedMeshes[0]
      if (!skinned) {
        setupIndependent(model, animations)
        return
      }

      const skeleton = skinned.skeleton
      const parentBone = hips ?? skeleton.bones[0]
      if (parentBone) scene.add(parentBone)

      const identity = new THREE.Matrix4()
      const xs = [-2, 0, 2]
      for (const x of xs) {
        const mesh = skinned.clone()
        mesh.bindMode = THREE.DetachedBindMode
        mesh.bind(skeleton, identity)
        mesh.position.x = x
        mesh.castShadow = true
        scene.add(mesh)
        objects.push(mesh)
      }
      if (parentBone) objects.push(parentBone)

      const mixer = new THREE.AnimationMixer(parentBone ?? sharedModel)
      const run = animations.find(c => /run/i.test(c.name)) ?? animations[0]
      if (run) mixer.clipAction(run).play()
      mixers.push(mixer)
    }

    let sourceModel: THREE.Object3D | null = null
    let sourceClips: THREE.AnimationClip[] = []

    const rebuild = () => {
      if (!sourceModel) return
      clearActors()
      if (modeRef.current === 'shared') setupShared(sourceModel, sourceClips)
      else setupIndependent(sourceModel, sourceClips)
    }

    const loader = new GLTFLoader()
    loader.load(
      modelUrl,
      gltf => {
        if (disposed) return
        sourceModel = gltf.scene
        sourceClips = gltf.animations ?? []
        // Normalize scale — Soldier is ~1 unit tall
        sourceModel.traverse(c => {
          if (c instanceof THREE.Mesh) {
            c.castShadow = true
            c.receiveShadow = true
          }
        })
        rebuild()
        setStatus('ready')
        setError(null)
      },
      undefined,
      err => {
        if (disposed) return
        console.error(err)
        setStatus('error')
        setError(err instanceof Error ? err.message : 'Failed to load Soldier.glb')
      }
    )

    // Re-build when mode prop changes via custom event from effect cleanup pattern
    const onModeHint = () => rebuild()
    wrap.addEventListener('skeleton-mode', onModeHint)

    const resize = () => {
      const w = wrap.clientWidth || 640
      const h = wrap.clientHeight || 360
      renderer.setSize(w, h, false)
      camera.aspect = w / h
      camera.updateProjectionMatrix()
    }
    resize()
    const ro = new ResizeObserver(resize)
    ro.observe(wrap)

    let raf = 0
    let pageVisible = document.visibilityState === 'visible'
    const onVis = () => {
      pageVisible = document.visibilityState === 'visible'
    }
    document.addEventListener('visibilitychange', onVis)

    const animate = () => {
      raf = requestAnimationFrame(animate)
      if (!pageVisible) return
      const delta = clock.getDelta()
      if (!prefersReduced) {
        for (const m of mixers) m.update(delta)
        ring.rotation.z += 0.002
      }
      controls.update()
      renderer.render(scene, camera)
    }
    animate()

    // Expose rebuild for mode switches without remounting whole scene
    ;(wrap as HTMLDivElement & { __rebuild?: () => void }).__rebuild = rebuild

    return () => {
      disposed = true
      cancelAnimationFrame(raf)
      document.removeEventListener('visibilitychange', onVis)
      ro.disconnect()
      wrap.removeEventListener('skeleton-mode', onModeHint)
      controls.dispose()
      clearActors()
      ground.geometry.dispose()
      ;(ground.material as THREE.Material).dispose()
      ring.geometry.dispose()
      ;(ring.material as THREE.Material).dispose()
      grid.geometry.dispose()
      ;(grid.material as THREE.Material).dispose()
      scene.remove(grid)
      renderer.dispose()
    }
  }, [modelUrl])

  // Rebuild actors when mode changes
  useEffect(() => {
    const wrap = wrapRef.current as (HTMLDivElement & { __rebuild?: () => void }) | null
    wrap?.__rebuild?.()
  }, [mode])

  return (
    <div ref={wrapRef} className={cn('relative h-full w-full', className)}>
      <canvas ref={canvasRef} className='h-full w-full' />
      {status === 'loading' && (
        <div className='micro-label text-primary absolute inset-0 flex items-center justify-center bg-black/50'>
          LOADING GLTF · SkeletonUtils…
        </div>
      )}
      {status === 'error' && (
        <div className='absolute inset-0 flex items-center justify-center bg-black/70 p-4 text-center font-mono text-xs text-destructive'>
          {error ?? 'MODEL LOAD FAILED'}
          <br />
          place Soldier.glb in public/showcase/models/
        </div>
      )}
    </div>
  )
}

export default SkeletonDemo
