'use client'

/**
 * Rotating wireframe globe — ported from agentic-os-personal content-os/public/hud/globe.js
 * (Three.js wireframe + points + equator ring + pulse waves).
 *
 * Events (window):
 *   emits  "globe:click"
 *   listens "assistant:state" { state: idle|listening|thinking|speaking|running|idle }
 *   listens "assistant:level" { level: 0..1 }  // optional mic / activity pulse
 */

import { useEffect, useRef } from 'react'
import * as THREE from 'three'

export type GlobeState = 'idle' | 'listening' | 'thinking' | 'speaking' | 'running'

type TalonGlobeProps = {
  className?: string
  /** Drive pulse intensity 0..1 from live ops (optional). */
  activityLevel?: number
  /** Semantic state for spin/opacity (optional; also listens on window). */
  state?: GlobeState
  onClick?: () => void
}

const TalonGlobe = ({ className, activityLevel = 0, state = 'idle', onClick }: TalonGlobeProps) => {
  const wrapRef = useRef<HTMLDivElement>(null)
  const canvasRef = useRef<HTMLCanvasElement>(null)
  const stateRef = useRef(state)
  const levelRef = useRef(activityLevel)

  useEffect(() => {
    stateRef.current = state
  }, [state])

  useEffect(() => {
    levelRef.current = activityLevel
  }, [activityLevel])

  useEffect(() => {
    const canvas = canvasRef.current
    const wrap = wrapRef.current
    if (!canvas || !wrap) return

    const scene = new THREE.Scene()
    const camera = new THREE.PerspectiveCamera(45, 1, 0.1, 100)
    camera.position.z = 3.0

    const renderer = new THREE.WebGLRenderer({ canvas, antialias: true, alpha: true })
    renderer.setPixelRatio(Math.min(window.devicePixelRatio || 1, 2))

    const globe = new THREE.Group()
    scene.add(globe)

    // Electric cyan operator wire — never brand-red chrome
    const CYAN = new THREE.Color(0x22d3ee)
    const MINT = new THREE.Color(0xa5f3fc)
    const WHITE = new THREE.Color(0xe0f2fe)

    const core = new THREE.Mesh(
      new THREE.SphereGeometry(0.98, 48, 48),
      new THREE.MeshBasicMaterial({ color: 0x060a0e })
    )
    globe.add(core)

    const wireMat = new THREE.LineBasicMaterial({ color: CYAN, transparent: true, opacity: 0.48 })
    const wire = new THREE.LineSegments(new THREE.WireframeGeometry(new THREE.SphereGeometry(1, 28, 20)), wireMat)
    globe.add(wire)

    const pointsMat = new THREE.PointsMaterial({ color: WHITE, size: 0.012, transparent: true, opacity: 0.75 })
    const points = new THREE.Points(new THREE.SphereGeometry(1.01, 40, 30), pointsMat)
    globe.add(points)

    const ring = new THREE.Mesh(
      new THREE.RingGeometry(1.18, 1.2, 96),
      new THREE.MeshBasicMaterial({ color: CYAN, transparent: true, opacity: 0.55, side: THREE.DoubleSide })
    )
    ring.rotation.x = Math.PI / 2.2
    globe.add(ring)

    globe.rotation.z = 0.36

    const waves: THREE.Mesh[] = []
    for (let i = 0; i < 3; i++) {
      const w = new THREE.Mesh(
        new THREE.RingGeometry(1.25, 1.27, 96),
        new THREE.MeshBasicMaterial({ color: MINT, transparent: true, opacity: 0, side: THREE.DoubleSide })
      )
      w.rotation.x = Math.PI / 2.2
      w.userData.phase = i / 3
      globe.add(w)
      waves.push(w)
    }

    const baseSpeed = 0.0016
    let speed = baseSpeed
    let targetOpacity = 0.42
    let micLevel = 0
    let micSmooth = 0

    const applyState = (s: string) => {
      if (s === 'listening') {
        speed = 0.012
        targetOpacity = 0.85
        wireMat.color.copy(WHITE)
      } else if (s === 'thinking' || s === 'running') {
        speed = 0.03
        targetOpacity = 0.75
        wireMat.color.copy(MINT)
      } else if (s === 'speaking') {
        speed = 0.008
        targetOpacity = 0.65
        wireMat.color.copy(CYAN)
      } else {
        targetOpacity = 0.48
        wireMat.color.copy(CYAN)
        if (s === 'idle') micLevel = Math.max(micLevel, 0)
      }
    }

    const resize = () => {
      const w = wrap.clientWidth || 360
      const h = wrap.clientHeight || 360
      renderer.setSize(w, h, false)
      camera.aspect = w / h
      camera.updateProjectionMatrix()
    }
    resize()
    window.addEventListener('resize', resize)
    const ro = typeof ResizeObserver !== 'undefined' ? new ResizeObserver(resize) : null
    ro?.observe(wrap)

    let raf = 0
    const animate = () => {
      raf = requestAnimationFrame(animate)

      // Prefer live props each frame
      applyState(stateRef.current)
      micLevel = Math.max(micLevel * 0.92, levelRef.current)

      globe.rotation.y += speed
      speed += (baseSpeed - speed) * 0.05
      wireMat.opacity += (targetOpacity - wireMat.opacity) * 0.08
      ring.rotation.z += 0.002

      micSmooth += (micLevel - micSmooth) * 0.25
      const pulse = 1 + micSmooth * 0.12
      globe.scale.setScalar(pulse)
      pointsMat.opacity = 0.55 + micSmooth * 0.45

      const t = performance.now() / 1000
      for (const w of waves) {
        const p = (t * 0.5 + (w.userData.phase as number)) % 1
        const s = 1 + p * (0.5 + micSmooth * 1.6)
        w.scale.setScalar(s)
        const mat = w.material as THREE.MeshBasicMaterial
        mat.opacity = Math.max(0, 1 - p) * 0.5 * micSmooth
      }

      renderer.render(scene, camera)
    }
    animate()

    const onClickWrap = () => {
      window.dispatchEvent(new CustomEvent('globe:click'))
      onClick?.()
    }
    wrap.addEventListener('click', onClickWrap)

    const onAssistantState = (e: Event) => {
      const s = (e as CustomEvent).detail?.state as string | undefined
      if (s) applyState(s)
    }
    const onAssistantLevel = (e: Event) => {
      const lvl = (e as CustomEvent).detail?.level
      if (typeof lvl === 'number') micLevel = Math.max(0, Math.min(1, lvl))
    }
    window.addEventListener('assistant:state', onAssistantState)
    window.addEventListener('assistant:level', onAssistantLevel)

    return () => {
      cancelAnimationFrame(raf)
      window.removeEventListener('resize', resize)
      ro?.disconnect()
      wrap.removeEventListener('click', onClickWrap)
      window.removeEventListener('assistant:state', onAssistantState)
      window.removeEventListener('assistant:level', onAssistantLevel)
      renderer.dispose()
      wire.geometry.dispose()
      wireMat.dispose()
      core.geometry.dispose()
      ;(core.material as THREE.Material).dispose()
      points.geometry.dispose()
      pointsMat.dispose()
      ring.geometry.dispose()
      ;(ring.material as THREE.Material).dispose()
      for (const w of waves) {
        w.geometry.dispose()
        ;(w.material as THREE.Material).dispose()
      }
    }
  }, [onClick])

  return (
    <div
      ref={wrapRef}
      className={className}
      role='img'
      aria-label='Talon operations globe'
      style={{ cursor: onClick ? 'pointer' : 'default' }}
    >
      <canvas ref={canvasRef} className='h-full w-full' />
    </div>
  )
}

export default TalonGlobe
