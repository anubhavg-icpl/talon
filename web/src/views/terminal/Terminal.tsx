'use client'

// React Imports
import { useSyncExternalStore } from 'react'

// Component Imports
import PageHeader from '@/components/shared/PageHeader'

// talon-core reverse-proxies ttyd at /shell behind the console session (SSO).
// Same host as the run WebSocket (:8000) so the talon_session cookie rides along.
const CORE_URL = process.env.NEXT_PUBLIC_TALON_CORE_URL

// Client-only flag without state-in-effect: false during SSR, true after mount.
const useMounted = () =>
  useSyncExternalStore(
    () => () => {},
    () => true,
    () => false
  )

/**
 * Integrated Kali shell — embeds the ttyd web terminal (full interactive PTY)
 * from arsenal_engine, reverse-proxied by talon-core at /shell. No second login:
 * the console session cookie authenticates the proxy (SSO).
 */
const Terminal = () => {
  const mounted = useMounted()

  const src = mounted
    ? `${CORE_URL ?? `${window.location.protocol}//${window.location.hostname}:8000`}/shell/`
    : null

  return (
    <div className='flex flex-col gap-6'>
      <PageHeader title='KALI SHELL' subtitle='LIVE PTY INTO THE ARSENAL (KALI) CONTAINER' />

      <div className='hud-corners relative overflow-hidden rounded-md border bg-black/80'>
        {src ? (
          <iframe src={src} title='Kali terminal' className='h-[70vh] w-full border-0' />
        ) : (
          <p className='micro-label p-6'>
            INITIALIZING TERMINAL<span className='terminal-cursor ml-1' />
          </p>
        )}
      </div>

      <p className='micro-label'>
        bash RUNS INSIDE arsenal_engine — AUTHENTICATED VIA YOUR CONSOLE SESSION (SSO) · AUTHORIZED USE ONLY
      </p>
    </div>
  )
}

export default Terminal
