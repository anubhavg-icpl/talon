'use client'

import { useState } from 'react'

import { useRouter } from 'next/navigation'

import Logo from '@/components/shared/Logo'
import MatrixRain from '@/components/shared/MatrixRain'
import TalonGlobe from '@/components/shared/TalonGlobe'
import Starfield from '@/components/shared/three/Starfield'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'

const Login = () => {
  const router = useRouter()
  const [username, setUsername] = useState('')
  const [password, setPassword] = useState('')
  const [busy, setBusy] = useState(false)
  const [error, setError] = useState<string | null>(null)

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

      router.push('/overview')
      router.refresh()
    } catch {
      setError('CORE UNREACHABLE — IS TALON-CORE UP?')
      setBusy(false)
    }
  }

  return (
    <div className='relative flex min-h-svh items-center justify-center overflow-hidden p-4'>
      {/* Full Three.js C2 background */}
      <div className='pointer-events-none absolute inset-0'>
        <TalonGlobe className='h-full w-full opacity-80' variant='background' state='running' activityLevel={0.4} />
      </div>
      <Starfield className='opacity-50' count={500} opacity={0.35} />
      <MatrixRain className='absolute inset-0 opacity-[0.08]' />
      <div className='scanlines pointer-events-none absolute inset-0' />
      <div className='pointer-events-none absolute inset-0 bg-gradient-to-b from-background/40 via-background/55 to-background/80' />

      <div className='hud-corners bg-card/85 relative z-10 w-full max-w-md rounded-sm border border-primary/20 p-8 backdrop-blur-md'>
        <div className='mb-8 flex flex-col items-center gap-3'>
          <Logo />
          <p className='micro-label text-primary/80 text-center'>C2 GATE — THREE.JS OPERATOR SHELL</p>
        </div>

        <form onSubmit={submit} className='flex flex-col gap-5'>
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
            <p className='border-destructive/40 bg-destructive/10 text-destructive rounded-sm border px-3 py-2 font-mono text-xs tracking-widest uppercase'>
              {error}
            </p>
          )}

          <Button
            type='submit'
            disabled={busy}
            className='glow-cyan mt-2 font-mono text-xs font-semibold tracking-widest uppercase'
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
