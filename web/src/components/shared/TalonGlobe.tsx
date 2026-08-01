'use client'

/**
 * Operator command globe — pure Three.js (no external glTF).
 * Variants: hero (interactive orbit) · compact · background (fullscreen C2).
 */

import { useEffect, useRef } from 'react'
import * as THREE from 'three'
import { OrbitControls } from 'three/examples/jsm/controls/OrbitControls.js'

import { cn } from '@/lib/utils'

export type GlobeState = 'idle' | 'listening' | 'thinking' | 'speaking' | 'running'
export type GlobeVariant = 'hero' | 'compact' | 'background'

type TalonGlobeProps = {
  className?: string
  activityLevel?: number
  state?: GlobeState
  variant?: GlobeVariant
  /** Enable orbit drag (default true for hero). */
  interactive?: boolean
  onClick?: () => void
}

const CYAN = 0x22d3ee
const MINT = 0xa5f3fc
const DEEP = 0x060a0e

function fibonacciSphere(n: number, radius: number): THREE.Vector3[] {
  const out: THREE.Vector3[] = []
  const golden = Math.PI * (3 - Math.sqrt(5))
  for (let i = 0; i < n; i++) {
    const y = 1 - (i / Math.max(1, n - 1)) * 2
    const r = Math.sqrt(Math.max(0, 1 - y * y))
    const theta = golden * i
    out.push(new THREE.Vector3(Math.cos(theta) * r * radius, y * radius, Math.sin(theta) * r * radius))
  }
  return out
}

function greatCircle(a: THREE.Vector3, b: THREE.Vector3, segments = 48): THREE.BufferGeometry {
  const pts: THREE.Vector3[] = []
  const len = a.length()
  for (let i = 0; i <= segments; i++) {
    const t = i / segments
    const v = new THREE.Vector3().copy(a).lerp(b, t).normalize()
    const lift = Math.sin(t * Math.PI) * 0.08
    pts.push(v.multiplyScalar(len * 1.02 * (1 + lift / len)))
  }
  return new THREE.BufferGeometry().setFromPoints(pts)
}

function starField(count: number): THREE.BufferGeometry {
  const pos = new Float32Array(count * 3)
  for (let i = 0; i < count; i++) {
    const r = 8 + Math.random() * 22
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

const TalonGlobe = ({
  className,
  activityLevel = 0,
  state = 'idle',
  variant = 'compact',
  interactive,
  onClick
}: TalonGlobeProps) => {
  const wrapRef = useRef<HTMLDivElement>(null)
  const canvasRef = useRef<HTMLCanvasElement>(null)
  const stateRef = useRef(state)
  const levelRef = useRef(activityLevel)
  const canOrbit = interactive ?? variant === 'hero'

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
    const isBg = variant === 'background'
    const isHero = variant === 'hero'
    const starCount = isBg ? 1400 : isHero ? 1100 : 700
    const beaconN = isBg ? 20 : isHero ? 16 : 12
    const arcN = isBg ? 12 : isHero ? 10 : 6

    const scene = new THREE.Scene()
    if (isBg) scene.fog = new THREE.FogExp2(0x060a0e, 0.035)

    const camera = new THREE.PerspectiveCamera(isBg ? 48 : 42, 1, 0.1, 100)
    camera.position.set(0, isBg ? 0.25 : 0.15, isBg ? 3.6 : isHero ? 3.0 : 3.15)

    const renderer = new THREE.WebGLRenderer({ canvas, antialias: true, alpha: true })
    renderer.setPixelRatio(Math.min(window.devicePixelRatio || 1, isBg ? 1.75 : 2))
    renderer.setClearColor(0x000000, 0)

    const key = new THREE.DirectionalLight(0x7dd3fc, 1.15)
    key.position.set(2.5, 1.8, 3)
    scene.add(key)
    const rim = new THREE.DirectionalLight(0x22d3ee, 0.55)
    rim.position.set(-2, -0.5, -1.5)
    scene.add(rim)
    scene.add(new THREE.AmbientLight(0x0e7490, 0.4))
    const point = new THREE.PointLight(0x22d3ee, 0.85, 14, 2)
    point.position.set(0, 0.5, 2.5)
    scene.add(point)

    const starsGeo = starField(starCount)
    const starsMat = new THREE.PointsMaterial({
      color: 0xbae6fd,
      size: isBg ? 0.035 : 0.028,
      transparent: true,
      opacity: isBg ? 0.65 : 0.55,
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

    const coreMat = new THREE.MeshStandardMaterial({
      color: DEEP,
      roughness: 0.82,
      metalness: 0.28,
      emissive: 0x022c36,
      emissiveIntensity: 0.5
    })
    const core = new THREE.Mesh(new THREE.SphereGeometry(0.96, isHero || isBg ? 72 : 48, isHero || isBg ? 72 : 48), coreMat)
    globe.add(core)

    const atmoMat = new THREE.MeshBasicMaterial({
      color: CYAN,
      transparent: true,
      opacity: isBg ? 0.1 : 0.07,
      side: THREE.BackSide,
      depthWrite: false
    })
    const atmo = new THREE.Mesh(new THREE.SphereGeometry(1.14, 48, 48), atmoMat)
    globe.add(atmo)

    // Inner energy core
    const energy = new THREE.Mesh(
      new THREE.SphereGeometry(0.22, 24, 24),
      new THREE.MeshBasicMaterial({ color: MINT, transparent: true, opacity: 0.55 })
    )
    globe.add(energy)

    const wireMat = new THREE.LineBasicMaterial({ color: CYAN, transparent: true, opacity: 0.42 })
    const wire = new THREE.LineSegments(new THREE.WireframeGeometry(new THREE.SphereGeometry(1.0, 36, 28)), wireMat)
    globe.add(wire)

    const surfacePts = new THREE.Points(
      new THREE.SphereGeometry(1.005, 52, 40),
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

    const ringGroup = new THREE.Group()
    root.add(ringGroup)
    const ringMats: THREE.MeshBasicMaterial[] = []
    const ringConfigs = [
      { r: 1.22, tilt: Math.PI / 2.15, color: CYAN, opacity: 0.5 },
      { r: 1.34, tilt: Math.PI / 2.6, color: MINT, opacity: 0.28 },
      { r: 1.48, tilt: Math.PI / 1.75, color: CYAN, opacity: 0.18 },
      ...(isHero || isBg
        ? [{ r: 1.62, tilt: Math.PI / 3.1, color: MINT, opacity: 0.12 }]
        : [])
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

    const beacons = fibonacciSphere(beaconN, 1.02)
    const beaconGroup = new THREE.Group()
    globe.add(beaconGroup)
    const beaconMat = new THREE.MeshBasicMaterial({ color: MINT, transparent: true, opacity: 0.9 })
    const beaconGeo = new THREE.SphereGeometry(0.018, 10, 10)
    for (const p of beacons) {
      const m = new THREE.Mesh(beaconGeo, beaconMat)
      m.position.copy(p)
      beaconGroup.add(m)
    }

    const arcGroup = new THREE.Group()
    globe.add(arcGroup)
    const arcMats: THREE.LineBasicMaterial[] = []
    const disposables: THREE.BufferGeometry[] = []
    for (let i = 0; i < arcN; i++) {
      const a = beacons[i % beacons.length]
      const b = beacons[(i * 3 + 5) % beacons.length]
      const geo = greatCircle(a, b, 56)
      disposables.push(geo)
      const mat = new THREE.LineBasicMaterial({
        color: i % 2 === 0 ? CYAN : MINT,
        transparent: true,
        opacity: 0.35
      })
      arcMats.push(mat)
      arcGroup.add(new THREE.Line(geo, mat))
    }

    // Latitude bands
    for (const lat of [-0.55, 0, 0.55]) {
      const pts: THREE.Vector3[] = []
      for (let i = 0; i <= 64; i++) {
        const a = (i / 64) * Math.PI * 2
        const r = Math.cos(lat)
        pts.push(new THREE.Vector3(Math.cos(a) * r * 1.01, Math.sin(lat) * 1.01, Math.sin(a) * r * 1.01))
      }
      const g = new THREE.BufferGeometry().setFromPoints(pts)
      disposables.push(g)
      globe.add(
        new THREE.LineLoop(
          g,
          new THREE.LineBasicMaterial({ color: CYAN, transparent: true, opacity: 0.22 })
        )
      )
    }

    const reticle = new THREE.Group()
    if (!isBg) root.add(reticle)
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

    let controls: OrbitControls | null = null
    if (canOrbit && !prefersReduced) {
      controls = new OrbitControls(camera, canvas)
      controls.enableDamping = true
      controls.dampingFactor = 0.06
      controls.enablePan = false
      controls.minDistance = 2.2
      controls.maxDistance = 5.5
      controls.autoRotate = true
      controls.autoRotateSpeed = isHero ? 0.6 : 0.35
      controls.rotateSpeed = 0.55
    }

    const pointer = { x: 0, y: 0 }
    const onPointer = (e: PointerEvent) => {
      if (canOrbit) return
      const rect = wrap.getBoundingClientRect()
      pointer.x = ((e.clientX - rect.left) / rect.width) * 2 - 1
      pointer.y = -(((e.clientY - rect.top) / rect.height) * 2 - 1)
    }
    wrap.addEventListener('pointermove', onPointer)

    const baseSpeed = prefersReduced ? 0 : isBg ? 0.0018 : 0.0022
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
      if (controls) {
        controls.autoRotateSpeed = s === 'running' || s === 'thinking' ? 1.4 : isHero ? 0.6 : 0.35
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

      if (!prefersReduced && !canOrbit) {
        globe.rotation.y += speed
        ringGroup.rotation.y -= speed * 0.35
        ringGroup.rotation.z += speed * 0.15
        stars.rotation.y += 0.00025
        reticle.rotation.z -= 0.0015
      } else if (!prefersReduced && canOrbit) {
        globe.rotation.y += speed * 0.35
        ringGroup.rotation.z += speed * 0.08
        stars.rotation.y += 0.0002
      }

      speed += (baseSpeed - speed) * 0.05
      wireMat.opacity += (targetOpacity - wireMat.opacity) * 0.08
      energy.scale.setScalar(1 + Math.sin(performance.now() / 500) * 0.08 + micSmooth * 0.15)

      micSmooth += (micLevel - micSmooth) * 0.22
      const pulse = 1 + micSmooth * 0.1
      if (!canOrbit) globe.scale.setScalar(pulse)
      surfacePtsMat.opacity = 0.5 + micSmooth * 0.45
      atmoMat.opacity = (isBg ? 0.08 : 0.05) + micSmooth * 0.08
      beaconMat.opacity = 0.55 + micSmooth * 0.45

      if (!canOrbit) {
        root.rotation.y += (pointer.x * 0.18 - root.rotation.y) * 0.06
        root.rotation.x += (pointer.y * 0.12 - root.rotation.x) * 0.06
      }

      const t = performance.now() / 1000
      for (const w of waves) {
        const p = (t * 0.45 + (w.userData.phase as number)) % 1
        const s = 1 + p * (0.55 + micSmooth * 1.4)
        w.scale.setScalar(s)
        ;(w.material as THREE.MeshBasicMaterial).opacity =
          Math.max(0, 1 - p) * 0.45 * (0.25 + micSmooth)
      }
      for (let i = 0; i < arcMats.length; i++) {
        arcMats[i].opacity = 0.22 + 0.2 * Math.sin(t * 1.4 + i) * (0.4 + micSmooth)
      }

      controls?.update()
      renderer.render(scene, camera)
    }
    animate()

    const onClickWrap = () => {
      window.dispatchEvent(new CustomEvent('globe:click'))
      onClick?.()
    }
    if (onClick) wrap.addEventListener('click', onClickWrap)

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
      if (onClick) wrap.removeEventListener('click', onClickWrap)
      wrap.removeEventListener('pointermove', onPointer)
      window.removeEventListener('assistant:state', onAssistantState)
      window.removeEventListener('assistant:level', onAssistantLevel)
      controls?.dispose()
      renderer.dispose()
      starsGeo.dispose()
      starsMat.dispose()
      core.geometry.dispose()
      coreMat.dispose()
      atmo.geometry.dispose()
      atmoMat.dispose()
      energy.geometry.dispose()
      ;(energy.material as THREE.Material).dispose()
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
  }, [onClick, variant, canOrbit])

  return (
    <div
      ref={wrapRef}
      className={cn(className)}
      role='img'
      aria-label='Talon operator Three.js globe'
      style={{ cursor: canOrbit ? 'grab' : onClick ? 'pointer' : 'default' }}
    >
      <canvas ref={canvasRef} className='h-full w-full' />
    </div>
  )
}

export default TalonGlobe
