'use client'

/**
 * Operator command globe — pure Three.js (no external glTF).
 * Patterns drawn from classic Three.js HUD / recon visuals:
 * wireframe sphere, atmosphere shell, starfield, great-circle arcs,
 * beacon nodes, multi-axis rings, activity-driven spin.
 *
 * Events (window):
 *   emits  "globe:click"
 *   listens "assistant:state" { state }
 *   listens "assistant:level" { level: 0..1 }
 */

import { useEffect, useRef } from 'react'
import * as THREE from 'three'

export type GlobeState = 'idle' | 'listening' | 'thinking' | 'speaking' | 'running'

type TalonGlobeProps = {
  className?: string
  activityLevel?: number
  state?: GlobeState
  onClick?: () => void
}

const CYAN = 0x22d3ee
const MINT = 0xa5f3fc
const DEEP = 0x060a0e

/** Fibonacci sphere samples — even beacon placement. */
function fibonacciSphere(n: number, radius: number): THREE.Vector3[] {
  const out: THREE.Vector3[] = []
  const golden = Math.PI * (3 - Math.sqrt(5))
  for (let i = 0; i < n; i++) {
    const y = 1 - (i / (n - 1)) * 2
    const r = Math.sqrt(1 - y * y)
    const theta = golden * i
    out.push(new THREE.Vector3(Math.cos(theta) * r * radius, y * radius, Math.sin(theta) * r * radius))
  }
  return out
}

/** Great-circle arc between two points on a sphere. */
function greatCircle(a: THREE.Vector3, b: THREE.Vector3, segments = 48): THREE.BufferGeometry {
  const pts: THREE.Vector3[] = []
  for (let i = 0; i <= segments; i++) {
    const t = i / segments
    const v = new THREE.Vector3().copy(a).lerp(b, t).normalize().multiplyScalar(a.length() * 1.02)
    // lift mid-arc slightly
    const lift = Math.sin(t * Math.PI) * 0.08
    pts.push(v.clone().multiplyScalar(1 + lift / a.length()))
  }
  return new THREE.BufferGeometry().setFromPoints(pts)
}

function starField(count: number): THREE.BufferGeometry {
  const pos = new Float32Array(count * 3)
  for (let i = 0; i < count; i++) {
    const r = 8 + Math.random() * 18
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

    const prefersReduced = window.matchMedia('(prefers-reduced-motion: reduce)').matches

    const scene = new THREE.Scene()
    const camera = new THREE.PerspectiveCamera(42, 1, 0.1, 80)
    camera.position.set(0, 0.15, 3.15)

    const renderer = new THREE.WebGLRenderer({ canvas, antialias: true, alpha: true })
    renderer.setPixelRatio(Math.min(window.devicePixelRatio || 1, 2))
    renderer.setClearColor(0x000000, 0)

    // Soft key + rim — operator cyan, not red
    const key = new THREE.DirectionalLight(0x7dd3fc, 1.1)
    key.position.set(2.5, 1.8, 3)
    scene.add(key)
    const rim = new THREE.DirectionalLight(0x22d3ee, 0.55)
    rim.position.set(-2, -0.5, -1.5)
    scene.add(rim)
    scene.add(new THREE.AmbientLight(0x0e7490, 0.35))

    // --- Starfield ---
    const starsGeo = starField(900)
    const starsMat = new THREE.PointsMaterial({
      color: 0xbae6fd,
      size: 0.028,
      transparent: true,
      opacity: 0.55,
      sizeAttenuation: true,
      depthWrite: false
    })
    const stars = new THREE.Points(starsGeo, starsMat)
    scene.add(stars)

    const root = new THREE.Group()
    scene.add(root)

    const globe = new THREE.Group()
    root.add(globe)
    globe.rotation.z = 0.28

    // Core + atmosphere (backside fresnel-ish via double sphere)
    const coreMat = new THREE.MeshStandardMaterial({
      color: DEEP,
      roughness: 0.85,
      metalness: 0.25,
      emissive: 0x022c36,
      emissiveIntensity: 0.45
    })
    const core = new THREE.Mesh(new THREE.SphereGeometry(0.96, 64, 64), coreMat)
    globe.add(core)

    const atmoMat = new THREE.MeshBasicMaterial({
      color: CYAN,
      transparent: true,
      opacity: 0.07,
      side: THREE.BackSide,
      depthWrite: false
    })
    const atmo = new THREE.Mesh(new THREE.SphereGeometry(1.12, 48, 48), atmoMat)
    globe.add(atmo)

    // Wireframe grid
    const wireMat = new THREE.LineBasicMaterial({ color: CYAN, transparent: true, opacity: 0.42 })
    const wire = new THREE.LineSegments(
      new THREE.WireframeGeometry(new THREE.SphereGeometry(1.0, 32, 24)),
      wireMat
    )
    globe.add(wire)

    // Surface points
    const surfacePts = new THREE.Points(
      new THREE.SphereGeometry(1.005, 48, 36),
      new THREE.PointsMaterial({
        color: 0xe0f2fe,
        size: 0.011,
        transparent: true,
        opacity: 0.65,
        sizeAttenuation: true,
        depthWrite: false
      })
    )
    globe.add(surfacePts)
    const surfacePtsMat = surfacePts.material as THREE.PointsMaterial

    // Multi-axis orbit rings
    const ringGroup = new THREE.Group()
    root.add(ringGroup)
    const ringMats: THREE.MeshBasicMaterial[] = []
    const ringConfigs = [
      { r: 1.22, tilt: Math.PI / 2.15, color: CYAN, opacity: 0.5 },
      { r: 1.34, tilt: Math.PI / 2.6, color: MINT, opacity: 0.28 },
      { r: 1.48, tilt: Math.PI / 1.75, color: CYAN, opacity: 0.18 }
    ]
    for (const cfg of ringConfigs) {
      const mat = new THREE.MeshBasicMaterial({
        color: cfg.color,
        transparent: true,
        opacity: cfg.opacity,
        side: THREE.DoubleSide,
        depthWrite: false
      })
      ringMats.push(mat)
      const mesh = new THREE.Mesh(new THREE.RingGeometry(cfg.r, cfg.r + 0.012, 128), mat)
      mesh.rotation.x = cfg.tilt
      ringGroup.add(mesh)
    }

    // Pulse waves on primary ring plane
    const waves: THREE.Mesh[] = []
    for (let i = 0; i < 3; i++) {
      const mat = new THREE.MeshBasicMaterial({
        color: MINT,
        transparent: true,
        opacity: 0,
        side: THREE.DoubleSide,
        depthWrite: false
      })
      const w = new THREE.Mesh(new THREE.RingGeometry(1.2, 1.22, 96), mat)
      w.rotation.x = Math.PI / 2.15
      w.userData.phase = i / 3
      root.add(w)
      waves.push(w)
    }

    // Beacon nodes on sphere
    const beacons = fibonacciSphere(14, 1.02)
    const beaconGroup = new THREE.Group()
    globe.add(beaconGroup)
    const beaconMat = new THREE.MeshBasicMaterial({ color: MINT, transparent: true, opacity: 0.9 })
    const beaconGeo = new THREE.SphereGeometry(0.018, 10, 10)
    for (const p of beacons) {
      const m = new THREE.Mesh(beaconGeo, beaconMat)
      m.position.copy(p)
      beaconGroup.add(m)
    }

    // Attack / recon arcs between random beacons
    const arcGroup = new THREE.Group()
    globe.add(arcGroup)
    const arcMats: THREE.LineBasicMaterial[] = []
    const disposables: THREE.BufferGeometry[] = []
    for (let i = 0; i < 8; i++) {
      const a = beacons[i % beacons.length]
      const b = beacons[(i * 3 + 5) % beacons.length]
      const geo = greatCircle(a, b, 56)
      disposables.push(geo)
      const mat = new THREE.LineBasicMaterial({
        color: i % 2 === 0 ? CYAN : MINT,
        transparent: true,
        opacity: 0.35 + (i % 3) * 0.08
      })
      arcMats.push(mat)
      arcGroup.add(new THREE.Line(geo, mat))
    }

    // Crosshair reticle (fixed, does not spin with globe)
    const reticle = new THREE.Group()
    root.add(reticle)
    const retMat = new THREE.LineBasicMaterial({ color: CYAN, transparent: true, opacity: 0.35 })
    const retRing = new THREE.LineLoop(
      new THREE.BufferGeometry().setFromPoints(
        Array.from({ length: 64 }, (_, i) => {
          const a = (i / 64) * Math.PI * 2
          return new THREE.Vector3(Math.cos(a) * 1.62, Math.sin(a) * 1.62, 0)
        })
      ),
      retMat
    )
    reticle.add(retRing)
    const tickGeo = new THREE.BufferGeometry().setFromPoints([
      new THREE.Vector3(0, 1.55, 0),
      new THREE.Vector3(0, 1.72, 0),
      new THREE.Vector3(0, -1.55, 0),
      new THREE.Vector3(0, -1.72, 0),
      new THREE.Vector3(1.55, 0, 0),
      new THREE.Vector3(1.72, 0, 0),
      new THREE.Vector3(-1.55, 0, 0),
      new THREE.Vector3(-1.72, 0, 0)
    ])
    reticle.add(new THREE.LineSegments(tickGeo, retMat))

    // Pointer parallax
    const pointer = { x: 0, y: 0 }
    const onPointer = (e: PointerEvent) => {
      const rect = wrap.getBoundingClientRect()
      pointer.x = ((e.clientX - rect.left) / rect.width) * 2 - 1
      pointer.y = -(((e.clientY - rect.top) / rect.height) * 2 - 1)
    }
    wrap.addEventListener('pointermove', onPointer)

    const baseSpeed = prefersReduced ? 0 : 0.0022
    let speed = baseSpeed
    let targetOpacity = 0.42
    let micLevel = 0
    let micSmooth = 0

    const applyState = (s: string) => {
      if (s === 'listening') {
        speed = prefersReduced ? 0 : 0.014
        targetOpacity = 0.85
        wireMat.color.setHex(0xe0f2fe)
      } else if (s === 'thinking' || s === 'running') {
        speed = prefersReduced ? 0 : 0.032
        targetOpacity = 0.78
        wireMat.color.setHex(MINT)
      } else if (s === 'speaking') {
        speed = prefersReduced ? 0 : 0.01
        targetOpacity = 0.65
        wireMat.color.setHex(CYAN)
      } else {
        speed = baseSpeed
        targetOpacity = 0.45
        wireMat.color.setHex(CYAN)
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
      applyState(stateRef.current)
      micLevel = Math.max(micLevel * 0.92, levelRef.current)

      if (!prefersReduced) {
        globe.rotation.y += speed
        ringGroup.rotation.y -= speed * 0.35
        ringGroup.rotation.z += speed * 0.15
        stars.rotation.y += 0.00025
        reticle.rotation.z -= 0.0015
      }
      speed += (baseSpeed - speed) * 0.05
      wireMat.opacity += (targetOpacity - wireMat.opacity) * 0.08

      micSmooth += (micLevel - micSmooth) * 0.22
      const pulse = 1 + micSmooth * 0.1
      globe.scale.setScalar(pulse)
      surfacePtsMat.opacity = 0.5 + micSmooth * 0.45
      atmoMat.opacity = 0.05 + micSmooth * 0.08
      beaconMat.opacity = 0.55 + micSmooth * 0.45

      // Parallax root
      root.rotation.y += (pointer.x * 0.18 - root.rotation.y) * 0.06
      root.rotation.x += (pointer.y * 0.12 - root.rotation.x) * 0.06

      const t = performance.now() / 1000
      for (const w of waves) {
        const p = (t * 0.45 + (w.userData.phase as number)) % 1
        const s = 1 + p * (0.55 + micSmooth * 1.4)
        w.scale.setScalar(s)
        const mat = w.material as THREE.MeshBasicMaterial
        mat.opacity = Math.max(0, 1 - p) * 0.45 * (0.25 + micSmooth)
      }
      // Arc breathe
      for (let i = 0; i < arcMats.length; i++) {
        arcMats[i].opacity = 0.22 + 0.2 * Math.sin(t * 1.4 + i) * (0.4 + micSmooth)
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
      wrap.removeEventListener('pointermove', onPointer)
      window.removeEventListener('assistant:state', onAssistantState)
      window.removeEventListener('assistant:level', onAssistantLevel)
      renderer.dispose()
      starsGeo.dispose()
      starsMat.dispose()
      core.geometry.dispose()
      coreMat.dispose()
      atmo.geometry.dispose()
      atmoMat.dispose()
      wire.geometry.dispose()
      wireMat.dispose()
      surfacePts.geometry.dispose()
      surfacePtsMat.dispose()
      beaconGeo.dispose()
      beaconMat.dispose()
      for (const m of ringMats) m.dispose()
      ringGroup.traverse(o => {
        if (o instanceof THREE.Mesh) o.geometry.dispose()
      })
      for (const w of waves) {
        w.geometry.dispose()
        ;(w.material as THREE.Material).dispose()
      }
      for (const g of disposables) g.dispose()
      for (const m of arcMats) m.dispose()
      retRing.geometry.dispose()
      tickGeo.dispose()
      retMat.dispose()
    }
  }, [onClick])

  return (
    <div
      ref={wrapRef}
      className={className}
      role='img'
      aria-label='Talon operator globe — click to launch engagement'
      style={{ cursor: onClick ? 'pointer' : 'default' }}
    >
      <canvas ref={canvasRef} className='h-full w-full' />
    </div>
  )
}

export default TalonGlobe
