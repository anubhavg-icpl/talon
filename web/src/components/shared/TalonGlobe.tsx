'use client'

/**
 * TalonGlobe — production operator globe (vasturiano/three-globe + three.js).
 *
 * Pattern matches the official three-globe basic example:
 *   https://github.com/vasturiano/three-globe/blob/master/example/basic/index.html
 *
 * Critical: AmbientLight + DirectionalLight (r155+ intensity scale),
 * correct camera Z, resize after layout, CDN earth textures with local fallback.
 */

import { useEffect, useRef } from 'react'
import * as THREE from 'three'
import { OrbitControls } from 'three/examples/jsm/controls/OrbitControls.js'
import ThreeGlobe from 'three-globe'

import { listRuns, type RunSummary } from '@/lib/api'

export type GlobeState = 'idle' | 'listening' | 'thinking' | 'speaking' | 'running'
export type GlobeVariant = 'hero' | 'compact' | 'background'

type TalonGlobeProps = {
  className?: string
  activityLevel?: number
  state?: GlobeState
  variant?: GlobeVariant
  interactive?: boolean
  onClick?: () => void
}

const C = {
  compromised: '#ef4444',
  active: '#f59e0b',
  error: '#7f1d1d',
  done: '#f87171',
  muted: '#71717a',
  ops: '#ffffff',
  atmosphere: '#ef4444'
}

// CDN textures (same as three-globe official examples) — always available.
const TEX = {
  earth: 'https://cdn.jsdelivr.net/npm/three-globe/example/img/earth-dark.jpg',
  bump: 'https://cdn.jsdelivr.net/npm/three-globe/example/img/earth-topology.png',
  // Local overrides when present (public/globe/*)
  earthLocal: '/globe/earth-dark.jpg',
  bumpLocal: '/globe/earth-topology.png'
}

const OPS: [number, number] = [38.9, -77.0]

function hashStr(s: string): number {
  let h = 2166136261
  for (let i = 0; i < s.length; i++) {
    h ^= s.charCodeAt(i)
    h = Math.imul(h, 16777619)
  }
  return h >>> 0
}

function ipToCoord(ip: string): [number, number] {
  const h = hashStr(ip || '0.0.0.0')
  const lat = (((h % 130000) / 130000) * 120 - 60)
  const lng = ((((h >> 16) % 200000) / 200000) * 360 - 180)
  return [lat, lng]
}

type RunStatus = RunSummary['status']

function statusColor(st: RunStatus | undefined, findings: number | undefined): string {
  if (!st) return C.muted
  if (st === 'error') return C.error
  if (st === 'running' || st === 'awaiting_approval' || st === 'initializing') return C.active
  if (st === 'completed') return findings && findings > 0 ? C.compromised : C.done
  return C.muted
}

function isActive(st?: RunStatus) {
  return st === 'running' || st === 'awaiting_approval' || st === 'initializing'
}

const STATUS_PRIORITY: RunStatus[] = ['running', 'awaiting_approval', 'initializing', 'error', 'completed']

function pickStatus(a: RunStatus | undefined, b: RunStatus): RunStatus {
  if (!a) return b
  return STATUS_PRIORITY.indexOf(a) <= STATUS_PRIORITY.indexOf(b) ? a : b
}

type Point = { lat: number; lng: number; color: string; size: number; label: string }
type Arc = { startLat: number; startLng: number; endLat: number; endLng: number; color: string }
type Ring = { lat: number; lng: number; color: string; maxR: number }

function demoData() {
  // Visible default markers so compact HUD is never an empty black sphere.
  const points: Point[] = [
    { lat: OPS[0], lng: OPS[1], color: C.ops, size: 0.55, label: 'OPS' },
    { lat: 51.5, lng: -0.1, color: C.active, size: 0.35, label: 'EU' },
    { lat: 35.6, lng: 139.7, color: C.done, size: 0.3, label: 'APAC' },
    { lat: -33.8, lng: 151.2, color: C.muted, size: 0.28, label: 'AU' },
    { lat: 1.3, lng: 103.8, color: C.compromised, size: 0.32, label: 'SG' }
  ]
  const arcs: Arc[] = points.slice(1).map(p => ({
    startLat: OPS[0],
    startLng: OPS[1],
    endLat: p.lat,
    endLng: p.lng,
    color: p.color
  }))
  const rings: Ring[] = [{ lat: OPS[0], lng: OPS[1], color: C.atmosphere, maxR: 4 }]
  return { points, arcs, rings }
}

function runsToGlobe(runs: RunSummary[] | null | undefined) {
  const demo = demoData()
  if (!runs?.length) return demo

  const points: Point[] = [{ lat: OPS[0], lng: OPS[1], color: C.ops, size: 0.55, label: 'OPS CENTER' }]
  const arcs: Arc[] = []
  const rings: Ring[] = [{ lat: OPS[0], lng: OPS[1], color: C.atmosphere, maxR: 3.5 }]

  const byTarget = new Map<string, { st: RunStatus | undefined; fc: number; n: number; label: string }>()
  for (const r of runs) {
    const t = r.target || '0.0.0.0'
    const cur = byTarget.get(t) || {
      st: undefined as RunStatus | undefined,
      fc: 0,
      n: 0,
      label: r.service_name || r.cve_id || t
    }
    cur.fc = Math.max(cur.fc, r.findings_count ?? 0)
    cur.st = pickStatus(cur.st, r.status)
    cur.n++
    byTarget.set(t, cur)
  }

  for (const [ip, info] of byTarget) {
    const [lat, lng] = ipToCoord(ip)
    const color = statusColor(info.st, info.fc)
    points.push({
      lat,
      lng,
      color,
      size: 0.28 + Math.min(0.4, info.n * 0.06),
      label: `${ip} · ${info.st}`
    })
    arcs.push({ startLat: OPS[0], startLng: OPS[1], endLat: lat, endLng: lng, color })
    if (isActive(info.st) || info.fc > 0) {
      rings.push({ lat, lng, color, maxR: info.fc > 0 ? 5 : 3 })
    }
  }
  return { points, arcs, rings }
}

function camZ(variant: GlobeVariant) {
  if (variant === 'background') return 420
  if (variant === 'hero') return 280
  return 320 // compact — framed HUD
}

const TalonGlobe = ({
  className,
  activityLevel = 0,
  state = 'idle',
  variant = 'compact',
  interactive,
  onClick
}: TalonGlobeProps) => {
  const mountRef = useRef<HTMLDivElement | null>(null)
  const live = useRef({ activityLevel, state, onClick })
  live.current = { activityLevel, state, onClick }

  useEffect(() => {
    const mount = mountRef.current
    if (!mount) return

    let cancelled = false
    let raf = 0
    let renderer: THREE.WebGLRenderer | null = null
    let controls: OrbitControls | null = null
    let ro: ResizeObserver | null = null

    const fail = (msg: string) => {
      if (!mount) return
      mount.innerHTML = `<div class="flex h-full min-h-[100px] w-full items-center justify-center bg-black/60 font-mono text-[9px] tracking-widest text-zinc-500">${msg}</div>`
    }

    const boot = () => {
      if (cancelled || !mount) return

      // Wait for layout — Lazy3D / flex can report 0×0 on first paint.
      let width = mount.clientWidth
      let height = mount.clientHeight
      if (width < 40 || height < 40) {
        width = Math.max(width, 200)
        height = Math.max(height, 200)
      }

      try {
        const probe = document.createElement('canvas')
        if (!(probe.getContext('webgl') || probe.getContext('experimental-webgl'))) {
          fail('WEBGL UNAVAILABLE')
          return
        }
      } catch {
        fail('WEBGL UNAVAILABLE')
        return
      }

      try {
        renderer = new THREE.WebGLRenderer({
          antialias: true,
          alpha: true,
          powerPreference: 'default',
          failIfMajorPerformanceCaveat: false
        })
      } catch {
        fail('WEBGL INIT FAILED')
        return
      }

      renderer.setPixelRatio(Math.min(window.devicePixelRatio || 1, 2))
      renderer.setSize(width, height, false)
      renderer.setClearColor(0x000000, 0)
      renderer.domElement.style.display = 'block'
      renderer.domElement.style.width = '100%'
      renderer.domElement.style.height = '100%'
      mount.replaceChildren(renderer.domElement)

      const scene = new THREE.Scene()

      // Lights are required for three-globe earth materials (official example).
      // three r155+ uses physical intensity — Math.PI scale matches docs.
      scene.add(new THREE.AmbientLight(0xcccccc, Math.PI))
      const dir = new THREE.DirectionalLight(0xffffff, 0.6 * Math.PI)
      dir.position.set(1, 1, 1)
      scene.add(dir)

      const camera = new THREE.PerspectiveCamera(50, width / height, 0.1, 2000)
      camera.position.z = camZ(variant)

      // Prefer local textures; fall back to CDN (always works in prod).
      const globe = new ThreeGlobe({ animateIn: variant !== 'compact', waitForGlobeReady: true })
        .globeImageUrl(TEX.earthLocal)
        .bumpImageUrl(TEX.bumpLocal)
        .showAtmosphere(true)
        .atmosphereColor(C.atmosphere)
        .atmosphereAltitude(0.2)
        .showGraticules(false)

      // If local assets 404, swap to CDN after a short delay when image fails.
      // three-globe loads async — also set CDN as primary if env wants reliability.
      try {
        const img = new Image()
        img.onerror = () => {
          if (cancelled) return
          globe.globeImageUrl(TEX.earth).bumpImageUrl(TEX.bump)
        }
        img.src = TEX.earthLocal
      } catch {
        globe.globeImageUrl(TEX.earth).bumpImageUrl(TEX.bump)
      }

      // Seed with demo markers so HUD is never empty black.
      const seed = demoData()
      globe
        .pointsData(seed.points)
        .pointLat('lat')
        .pointLng('lng')
        .pointColor('color')
        .pointAltitude(0.012)
        .pointRadius('size')
        .pointsMerge(false)
        .arcsData(seed.arcs)
        .arcColor('color')
        .arcDashLength(0.35)
        .arcDashGap(1.2)
        .arcDashAnimateTime(2000)
        .arcStroke(0.45)
        .arcAltitudeAutoScale(0.35)
        .ringsData(seed.rings)
        .ringColor(() => (t: number) => `rgba(239,68,68,${1 - t})`)
        .ringMaxRadius('maxR')
        .ringPropagationSpeed(2.2)
        .ringRepeatPeriod(1100)

      scene.add(globe)

      const canOrbit = interactive ?? (variant === 'hero' || variant === 'background')
      controls = new OrbitControls(camera, renderer.domElement)
      controls.enableDamping = true
      controls.dampingFactor = 0.08
      controls.rotateSpeed = 0.55
      controls.enableZoom = canOrbit
      controls.enablePan = false
      controls.minDistance = 140
      controls.maxDistance = 700
      controls.autoRotate = true
      controls.autoRotateSpeed = 0.5

      renderer.domElement.style.cursor = onClick ? 'pointer' : canOrbit ? 'grab' : 'default'
      const onClickDom = () => live.current.onClick?.()
      renderer.domElement.addEventListener('click', onClickDom)

      let alive = true
      listRuns(200)
        .then(res => {
          if (!alive || cancelled || !res?.runs?.length) return
          const d = runsToGlobe(res.runs)
          globe.pointsData(d.points).arcsData(d.arcs).ringsData(d.rings)
        })
        .catch(() => {
          /* keep demo data */
        })

      const animate = () => {
        if (cancelled || !renderer) return
        const { activityLevel: al, state: st } = live.current
        if (controls) {
          controls.autoRotateSpeed =
            0.35 + al * 1.0 + (st === 'running' || st === 'thinking' ? 0.65 : 0)
          controls.update()
        }
        renderer.render(scene, camera)
        raf = requestAnimationFrame(animate)
      }
      raf = requestAnimationFrame(animate)

      const onResize = () => {
        if (!mount || !renderer) return
        const w = mount.clientWidth || width
        const h = mount.clientHeight || height
        if (w < 2 || h < 2) return
        width = w
        height = h
        camera.aspect = w / h
        camera.updateProjectionMatrix()
        renderer.setSize(w, h, false)
      }
      ro = new ResizeObserver(onResize)
      ro.observe(mount)
      // Second pass after CSS layout settles (Lazy3D expand, flex).
      requestAnimationFrame(onResize)
      setTimeout(onResize, 120)

      const onState = (e: Event) => {
        const s = (e as CustomEvent).detail?.state as GlobeState | undefined
        if (s) live.current = { ...live.current, state: s }
      }
      window.addEventListener('assistant:state', onState as EventListener)

      // Store cleanup on mount element
      ;(mount as any).__talonGlobeCleanup = () => {
        alive = false
        cancelAnimationFrame(raf)
        window.removeEventListener('assistant:state', onState as EventListener)
        renderer?.domElement.removeEventListener('click', onClickDom)
        ro?.disconnect()
        controls?.dispose()
        renderer?.dispose()
        scene.clear()
        if (renderer?.domElement.parentNode === mount) {
          mount.removeChild(renderer.domElement)
        }
      }
    }

    // Defer one frame so Lazy3D / flex has real dimensions.
    const t = window.setTimeout(boot, 0)

    return () => {
      cancelled = true
      clearTimeout(t)
      const cleanup = (mount as any).__talonGlobeCleanup as (() => void) | undefined
      cleanup?.()
      ;(mount as any).__talonGlobeCleanup = undefined
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [variant, interactive])

  return (
    <div
      ref={mountRef}
      className={className}
      style={{ width: '100%', height: '100%', minHeight: 120, minWidth: 120, position: 'relative' }}
    />
  )
}

export default TalonGlobe
