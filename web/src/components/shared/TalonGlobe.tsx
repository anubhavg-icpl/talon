'use client'

/**
 * TalonGlobe — interactive 3D operator globe on vasturiano/three-globe.
 *
 * Renders the real earth (dark texture + topology bump + red atmosphere) and
 * overlays LIVE engagement data pulled from the control-plane API:
 *   • points — every distinct target IP (coord derived deterministically; lab
 *              IPs carry no geo), sized by run count, colored by status/findings
 *   • arcs   — ops-center → each target (great-circle), colored by severity
 *   • rings  — pulsing markers on active / compromised targets
 *
 * Drop-in for the previous procedural globe: identical prop surface
 * (variant / state / activityLevel / interactive / onClick), so GlobePanel,
 * Overview, OpsHub, NewRun and Login keep working unchanged. Red-only palette.
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
  /** Enable orbit drag + zoom (default true for hero / background). */
  interactive?: boolean
  onClick?: () => void
}

// Red-only palette (status → color).
const C = {
  compromised: '#ef4444', // red-500 — findings confirmed
  active: '#f59e0b', // amber-500 — running / awaiting
  error: '#7f1d1d', // red-900 — errored
  done: '#f87171', // red-400 — completed, no findings
  muted: '#52525b', // zinc-600 — no signal
  ops: '#ffffff', // ops center
  atmosphere: '#ef4444'
}

// Notional ops center (operator HQ). Private-range lab target IPs have no real
// geo, so each is mapped to a stable deterministic coord via a string hash.
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
  const lat = (((h % 130000) / 130000) * 120 - 60) // -60..60 (avoid poles)
  const lng = ((((h >> 16) % 200000) / 200000) * 360 - 180) // -180..180
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

// Lower index = higher priority when a target appears across multiple runs.
const STATUS_PRIORITY: RunStatus[] = ['running', 'awaiting_approval', 'initializing', 'error', 'completed']
function pickStatus(a: RunStatus | undefined, b: RunStatus): RunStatus {
  if (!a) return b
  const ia = STATUS_PRIORITY.indexOf(a)
  const ib = STATUS_PRIORITY.indexOf(b)
  return ia <= ib ? a : b
}

type Point = { lat: number; lng: number; color: string; label: string; size: number }
type Arc = { startLat: number; startLng: number; endLat: number; endLng: number; color: string }
type Ring = { lat: number; lng: number; color: string; maxR: number }

function runsToGlobe(runs: RunSummary[] | null | undefined) {
  const points: Point[] = [{ lat: OPS[0], lng: OPS[1], color: C.ops, label: 'OPS CENTER', size: 0.45 }]
  const arcs: Arc[] = []
  const rings: Ring[] = []
  if (!runs || !runs.length) return { points, arcs, rings }

  const byTarget = new Map<string, { st: RunStatus | undefined; fc: number; n: number; label: string }>()
  for (const r of runs) {
    const t = r.target || '0.0.0.0'
    const cur = byTarget.get(t) || { st: undefined as RunStatus | undefined, fc: 0, n: 0, label: r.service_name || r.cve_id || t }
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
      label: `${ip} · ${info.label} · ${info.n} run(s) · ${info.st}`,
      size: 0.3 + Math.min(0.35, info.n * 0.05)
    })
    arcs.push({ startLat: OPS[0], startLng: OPS[1], endLat: lat, endLng: lng, color })
    if (isActive(info.st) || info.fc > 0) rings.push({ lat, lng, color, maxR: info.fc > 0 ? 5 : 3 })
  }
  return { points, arcs, rings }
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
  // Live values read inside the rAF loop without forcing a globe rebuild.
  const live = useRef({ activityLevel, state, onClick })
  live.current = { activityLevel, state, onClick }

  useEffect(() => {
    const mount = mountRef.current
    if (!mount) return

    let width = mount.clientWidth || 300
    let height = mount.clientHeight || 300
    const isBg = variant === 'background'
    const canOrbit = interactive ?? (variant === 'hero' || isBg)

    const scene = new THREE.Scene()
    const camera = new THREE.PerspectiveCamera(45, width / height, 0.1, 4000)
    camera.position.set(0, 0, isBg ? 380 : 240)

    const renderer = new THREE.WebGLRenderer({ antialias: true, alpha: !isBg })
    renderer.setPixelRatio(Math.min(window.devicePixelRatio, 2))
    renderer.setSize(width, height)
    mount.appendChild(renderer.domElement)

    // Night-sky backdrop only for the fullscreen variant; otherwise transparent
    // so the page-level starfield shows through.
    let sky: THREE.Mesh | null = null
    if (isBg) {
      const skyTex = new THREE.TextureLoader().load('/globe/night-sky.png')
      sky = new THREE.Mesh(
        new THREE.SphereGeometry(900, 64, 64),
        new THREE.MeshBasicMaterial({ map: skyTex, side: THREE.BackSide })
      )
      scene.add(sky)
    }

    const globe = new ThreeGlobe({ animateIn: true })
      .globeImageUrl('/globe/earth-dark.jpg')
      .bumpImageUrl('/globe/earth-topology.png')
      .showAtmosphere(true)
      .atmosphereColor(C.atmosphere)
      .atmosphereAltitude(0.18)
    scene.add(globe)

    // Empty layers first; refreshed when live runs resolve. three-globe's .d.ts
    // is incomplete for several accessors (pointLabel, string-arg color/size),
    // so the data-layer surface is cast loose; the texture/atmosphere methods
    // above stay typed.
    const g = globe as any
    g.pointsData([])
      .pointColor('color')
      .pointAltitude(0.01)
      .pointRadius('size')
      .pointLabel('label')
    g.arcsData([])
      .arcColor('color')
      .arcDashLength(0.4)
      .arcDashGap(2)
      .arcDashInitialGap(() => Math.random() * 2)
      .arcDashAnimateTime(2200)
      .arcStroke(0.18)
    g.ringsData([])
      .ringColor('color')
      .ringMaxRadius('maxR')
      .ringPropagationSpeed(3)
      .ringRepeatPeriod(900)

    const controls = new OrbitControls(camera, renderer.domElement)
    controls.enableDamping = true
    controls.dampingFactor = 0.12
    controls.rotateSpeed = 0.6
    controls.enableZoom = canOrbit
    controls.enablePan = false
    controls.minDistance = 150
    controls.maxDistance = 600
    controls.autoRotate = true

    renderer.domElement.style.cursor = onClick ? 'pointer' : 'grab'
    const clickHandler = () => live.current.onClick?.()
    renderer.domElement.addEventListener('click', clickHandler)

    // Pull live engagement data from the control plane (best-effort).
    let alive = true
    listRuns(200)
      .then(res => {
        if (!alive || !res?.runs) return
        const d = runsToGlobe(res.runs)
        g.pointsData(d.points).arcsData(d.arcs).ringsData(d.rings)
      })
      .catch(() => {
        /* globe still renders the earth without overlays */
      })

    const frame = () => {
      const { activityLevel: al, state: st } = live.current
      controls.autoRotateSpeed = 0.35 + al * 1.1 + (st === 'running' || st === 'thinking' ? 0.7 : 0)
      controls.update()
      renderer.render(scene, camera)
      raf = requestAnimationFrame(frame)
    }
    let raf = requestAnimationFrame(frame)

    const onResize = () => {
      width = mount.clientWidth
      height = mount.clientHeight
      if (!width || !height) return
      camera.aspect = width / height
      camera.updateProjectionMatrix()
      renderer.setSize(width, height)
    }
    const ro = new ResizeObserver(onResize)
    ro.observe(mount)

    // Overview dispatches assistant:state to nudge rotation.
    const onState = (e: Event) => {
      const s = (e as CustomEvent).detail?.state as GlobeState | undefined
      if (s) live.current = { ...live.current, state: s }
    }
    window.addEventListener('assistant:state', onState as EventListener)

    return () => {
      alive = false
      cancelAnimationFrame(raf)
      ro.disconnect()
      window.removeEventListener('assistant:state', onState as EventListener)
      renderer.domElement.removeEventListener('click', clickHandler)
      controls.dispose()
      renderer.dispose()
      if (renderer.domElement.parentNode === mount) mount.removeChild(renderer.domElement)
    }
    // Rebuild only on structural changes; state/activityLevel flow through `live`.
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [variant, interactive])

  return <div ref={mountRef} className={className} />
}

export default TalonGlobe
