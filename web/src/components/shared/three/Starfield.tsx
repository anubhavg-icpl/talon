'use client'

/**
 * Lightweight ambient Three.js starfield for app chrome (login / pages layout).
 * One canvas, low GPU cost, respects prefers-reduced-motion.
 */

import { useEffect, useRef } from 'react'
import * as THREE from 'three'

import { cn } from '@/lib/utils'

type StarfieldProps = {
  className?: string
  /** Particle count (default 700). */
  count?: number
  /** Opacity of points 0..1. */
  opacity?: number
  /** Slow drift speed. */
  speed?: number
}

function buildStars(count: number): THREE.BufferGeometry {
  const pos = new Float32Array(count * 3)
  for (let i = 0; i < count; i++) {
    const r = 12 + Math.random() * 40
    const theta = Math.random() * Math.PI * 2
    const phi = Math.acos(2 * Math.random() - 1)
    pos[i * 3] = r * Math.sin(phi) * Math.cos(theta)
    pos[i * 3 + 1] = r * Math.sin(phi) * Math.sin(theta)
    pos[i * 3 + 2] = r * Math.cos(phi)
  }
  const g = new THREE.BufferGeometry()
  g.setAttribute('position', new THREE.BufferAttribute(pos, 3))
  return g
}

const Starfield = ({ className, count = 700, opacity = 0.45, speed = 0.00035 }: StarfieldProps) => {
  const wrapRef = useRef<HTMLDivElement>(null)
  const canvasRef = useRef<HTMLCanvasElement>(null)

  useEffect(() => {
    const wrap = wrapRef.current
    const canvas = canvasRef.current
    if (!wrap || !canvas) return

    const reduced = window.matchMedia('(prefers-reduced-motion: reduce)').matches
    const scene = new THREE.Scene()
    const camera = new THREE.PerspectiveCamera(55, 1, 0.1, 120)
    camera.position.z = 1

    const renderer = new THREE.WebGLRenderer({ canvas, antialias: false, alpha: true })
    renderer.setPixelRatio(Math.min(window.devicePixelRatio || 1, 1.75))
    renderer.setClearColor(0x000000, 0)

    const geo = buildStars(count)
    const mat = new THREE.PointsMaterial({
      color: 0x7dd3fc,
      size: 0.04,
      transparent: true,
      opacity,
      sizeAttenuation: true,
      depthWrite: false
    })
    const points = new THREE.Points(geo, mat)
    scene.add(points)

    // Subtle cyan nebula fog plane (additive-ish via low opacity sphere)
    const haze = new THREE.Mesh(
      new THREE.SphereGeometry(18, 16, 16),
      new THREE.MeshBasicMaterial({
        color: 0x0e7490,
        transparent: true,
        opacity: 0.035,
        side: THREE.BackSide,
        depthWrite: false
      })
    )
    scene.add(haze)

    const resize = () => {
      const w = wrap.clientWidth || window.innerWidth
      const h = wrap.clientHeight || window.innerHeight
      renderer.setSize(w, h, false)
      camera.aspect = w / h
      camera.updateProjectionMatrix()
    }
    resize()
    const ro = new ResizeObserver(resize)
    ro.observe(wrap)

    let raf = 0
    const tick = () => {
      raf = requestAnimationFrame(tick)
      if (!reduced) {
        points.rotation.y += speed
        points.rotation.x += speed * 0.35
        haze.rotation.y -= speed * 0.2
      }
      renderer.render(scene, camera)
    }
    tick()

    return () => {
      cancelAnimationFrame(raf)
      ro.disconnect()
      renderer.dispose()
      geo.dispose()
      mat.dispose()
      haze.geometry.dispose()
      ;(haze.material as THREE.Material).dispose()
    }
  }, [count, opacity, speed])

  return (
    <div ref={wrapRef} className={cn('pointer-events-none absolute inset-0 overflow-hidden', className)} aria-hidden>
      <canvas ref={canvasRef} className='h-full w-full' />
    </div>
  )
}

export default Starfield
