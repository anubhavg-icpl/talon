'use client'

import { useEffect, useRef, useState } from 'react'

import { useRouter } from 'next/navigation'

import Logo from '@/components/shared/Logo'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'

const VIDEO_WEBM = '/globe/dark-planet.webm'
const VIDEO_MP4 = '/globe/dark-planet.mp4'
const POSTER = '/globe/operator-globe-hud.webp'

const Login = () => {
  const router = useRouter()
  const videoRef = useRef<HTMLVideoElement | null>(null)
  const [username, setUsername] = useState('')
  const [password, setPassword] = useState('')
  const [busy, setBusy] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [videoOk, setVideoOk] = useState(true)

  useEffect(() => {
    const v = videoRef.current
    if (!v || !videoOk) return
    v.play().catch(() => {
      /* autoplay blocked — poster + CSS still animate */
    })
  }, [videoOk])

  const submit = async (e: React.FormEvent) => {
    e.preventDefault()
    setBusy(true)
    setError(null)

    try {
      const res = await fetch('/api/talon/auth/login', {
        method: 'POST',
        headers: { 'content-type': 'application/json' },
        body: JSON.stringify({ username, password })
      })

      if (!res.ok) {
        const body = await res.json().catch(() => null)

        setError(res.status === 401 ? 'ACCESS DENIED — INVALID CREDENTIALS' : (body?.detail ?? `core error ${res.status}`))
        setBusy(false)

        return
      }

      const body = (await res.json().catch(() => null)) as { token?: string; username?: string } | null
      // Persist token for Arsenal Shell iframe (core often on :8000; cookie is origin-scoped).
      if (body?.token && typeof window !== 'undefined') {
        try {
          sessionStorage.setItem('talon_token', body.token)
          sessionStorage.setItem('talon_user', body.username || username)
        } catch {
          /* private mode */
        }
        // Dual-login to talon-core origin so HttpOnly talon_session is also set there
        // (needed for /shell WebSocket SSO when dashboard ≠ core port).
        const coreBase =
          process.env.NEXT_PUBLIC_TALON_CORE_URL ||
          `${window.location.protocol}//${window.location.hostname}:8000`
        fetch(`${coreBase}/auth/login`, {
          method: 'POST',
          credentials: 'include',
          headers: { 'content-type': 'application/json' },
          body: JSON.stringify({ username, password })
        }).catch(() => {
          /* core may be same-origin already or blocked — sessionStorage token is fallback */
        })
      }

      router.push('/overview')
      router.refresh()
    } catch {
      setError('CORE UNREACHABLE — IS TALON-CORE UP?')
      setBusy(false)
    }
  }

  return (
    <div className='login-stage relative flex min-h-svh items-center justify-center overflow-hidden p-4'>
      {/* Layer 0 — dark base */}
      <div className='absolute inset-0 bg-black' aria-hidden />

      {/* Layer 1 — dark-planet wallpaper (same asset as operator globe HUD) */}
      <div className='login-planet pointer-events-none absolute inset-0' aria-hidden>
        {videoOk ? (
          <video
            ref={videoRef}
            className='login-planet-video absolute inset-0 size-full object-cover'
            poster={POSTER}
            autoPlay
            muted
            loop
            playsInline
            preload='metadata'
            onError={() => setVideoOk(false)}
          >
            <source src={VIDEO_WEBM} type='video/webm' />
            <source src={VIDEO_MP4} type='video/mp4' />
          </video>
        ) : (
          // eslint-disable-next-line @next/next/no-img-element
          <img src={POSTER} alt='' className='absolute inset-0 size-full object-cover' />
        )}
      </div>

      {/* Layer 2 — animated ambient orbs + grid (CSS only, no WebGL) */}
      <div className='login-orbs pointer-events-none absolute inset-0' aria-hidden />
      <div className='login-grid pointer-events-none absolute inset-0' aria-hidden />
      <div className='login-sweep pointer-events-none absolute inset-0' aria-hidden />

      {/* Layer 3 — vignette + readability wash */}
      <div
        className='pointer-events-none absolute inset-0 bg-[radial-gradient(ellipse_at_center,transparent_0%,oklch(0.12_0.02_25/0.35)_45%,oklch(0.08_0.02_25/0.92)_100%)]'
        aria-hidden
      />
      <div className='pointer-events-none absolute inset-0 bg-gradient-to-b from-black/40 via-transparent to-black/80' aria-hidden />
      <div className='scanlines pointer-events-none absolute inset-0 opacity-40' aria-hidden />

      {/* Auth card */}
      <div className='login-card hud-corners bg-card/85 relative z-10 w-full max-w-md rounded-sm border border-primary/25 p-8 shadow-[0_0_64px_oklch(0.62_0.22_25/0.14),0_24px_48px_oklch(0_0_0/0.55)] backdrop-blur-xl'>
        <div className='mb-8 flex flex-col items-center gap-3'>
          {/* eslint-disable-next-line @next/next/no-img-element */}
          <img
            src='/showcase/talon-brand-mark.webp'
            alt=''
            className='login-mark mb-1 size-14 object-contain opacity-95'
            width={56}
            height={56}
          />
          <Logo />
          <p className='micro-label text-primary/85 text-center'>AUTHORIZED OPERATORS ONLY</p>
          <div className='flex items-center gap-2'>
            <span className='size-1.5 animate-pulse rounded-full bg-primary shadow-[0_0_8px_var(--primary)]' />
            <span className='font-mono text-[9px] tracking-[0.22em] text-zinc-400 uppercase'>
              C2 link · dark planet
            </span>
          </div>
        </div>

        <form onSubmit={submit} className='flex flex-col gap-5' noValidate>
          <div className='flex flex-col gap-2'>
            <Label htmlFor='username' className='micro-label'>
              OPERATOR ID
            </Label>
            <Input
              id='username'
              value={username}
              onChange={e => setUsername(e.target.value)}
              autoComplete='username'
              autoFocus
              required
              className='font-mono'
              placeholder='admin'
            />
          </div>

          <div className='flex flex-col gap-2'>
            <Label htmlFor='password' className='micro-label'>
              PASSPHRASE
            </Label>
            <Input
              id='password'
              type='password'
              value={password}
              onChange={e => setPassword(e.target.value)}
              autoComplete='current-password'
              required
              className='font-mono'
              placeholder='••••••••••••'
            />
          </div>

          {error && (
            <p
              role='alert'
              className='border-destructive/40 bg-destructive/10 text-destructive rounded-sm border px-3 py-2 font-mono text-xs tracking-widest uppercase'
            >
              {error}
            </p>
          )}

          <Button
            type='submit'
            disabled={busy}
            className='glow-red mt-2 font-mono text-xs font-semibold tracking-widest uppercase'
          >
            {busy ? '[ ⟳ LINKING… ]' : '[ ENTER C2 ]'}
          </Button>
        </form>

        <p className='micro-label mt-6 text-center'>$ talon auth login — same credentials as the CLI</p>
      </div>
    </div>
  )
}

export default Login
