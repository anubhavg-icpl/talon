'use client'

import { useState } from 'react'

import { useRouter } from 'next/navigation'

import Logo from '@/components/shared/Logo'
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
      {/* Static backdrop only — no WebGL on login (prevents tab crash before auth) */}
      <div
        className='pointer-events-none absolute inset-0 opacity-80'
        aria-hidden
        style={{
          backgroundImage: 'url(/globe/operator-globe-hud.webp)',
          backgroundSize: 'cover',
          backgroundPosition: 'center'
        }}
      />
      <div className='scanlines pointer-events-none absolute inset-0' />
      <div className='pointer-events-none absolute inset-0 bg-gradient-to-b from-background/50 via-background/65 to-background/90' />

      <div className='hud-corners bg-card/90 relative z-10 w-full max-w-md rounded-sm border border-primary/20 p-8 shadow-[0_0_48px_oklch(0.84_0.14_25/0.08)] backdrop-blur-md'>
        <div className='mb-8 flex flex-col items-center gap-3'>
          {/* eslint-disable-next-line @next/next/no-img-element */}
          <img
            src='/showcase/talon-brand-mark.webp'
            alt=''
            className='mb-1 size-14 object-contain opacity-90'
            width={56}
            height={56}
          />
          <Logo />
          <p className='micro-label text-primary/80 text-center'>AUTHORIZED OPERATORS ONLY</p>
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
