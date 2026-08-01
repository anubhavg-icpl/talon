'use client'

import { useSyncExternalStore } from 'react'

import { SquareTerminal } from 'lucide-react'

import LiveDot from '@/components/shared/LiveDot'
import PageHeader from '@/components/shared/PageHeader'
import { Badge } from '@/components/ui/badge'

// talon-core reverse-proxies ttyd at /shell behind the console session (SSO).
const CORE_URL = process.env.NEXT_PUBLIC_TALON_CORE_URL

const useMounted = () =>
  useSyncExternalStore(
    () => () => {},
    () => true,
    () => false
  )

/**
 * Integrated Kali shell — embeds ttyd (interactive PTY) from arsenal_engine,
 * reverse-proxied by talon-core at /shell. Theme is set server-side (carbon + red).
 */
const Terminal = () => {
  const mounted = useMounted()

  const src = mounted
    ? `${CORE_URL ?? `${window.location.protocol}//${window.location.hostname}:8000`}/shell/`
    : null

  return (
    <div className='flex flex-col gap-4'>
      <PageHeader
        title={
          <span className='inline-flex items-center gap-2.5'>
            <SquareTerminal className='text-primary size-6' />
            KALI SHELL
          </span>
        }
        subtitle='Live PTY into arsenal_engine (Kali) · session SSO'
        action={
          <Badge
            variant='outline'
            className='border-primary/40 bg-primary/10 font-mono text-[10px] tracking-widest text-primary uppercase'
          >
            <LiveDot tone='green' className='mr-1.5' />
            ARSENAL · LIVE
          </Badge>
        }
      />

      <div className='hud-corners relative overflow-hidden rounded-sm border border-primary/30 bg-black shadow-[0_0_40px_oklch(0.62_0.22_25/0.12)]'>
        {/* Top chrome bar */}
        <div className='flex items-center justify-between gap-3 border-b border-primary/20 bg-[oklch(0.08_0.02_25)] px-3 py-2'>
          <div className='flex items-center gap-2'>
            <span className='size-2 rounded-full bg-primary shadow-[0_0_8px_oklch(0.62_0.22_25)]' />
            <span className='font-mono text-[10px] tracking-[0.2em] text-primary/90 uppercase'>
              root@arsenal · bash
            </span>
          </div>
          <span className='text-muted-foreground font-mono text-[9px] tracking-widest uppercase'>
            xterm-256 · carbon / red
          </span>
        </div>

        {src ? (
          <iframe
            src={src}
            title='Kali terminal'
            className='h-[min(72vh,820px)] w-full border-0 bg-[#0a0a0a]'
            allow='clipboard-read; clipboard-write'
          />
        ) : (
          <div className='flex h-[min(72vh,820px)] items-center justify-center bg-[#0a0a0a]'>
            <p className='font-mono text-xs tracking-widest text-primary/80 uppercase'>
              Initializing terminal
              <span className='terminal-cursor ml-1' />
            </p>
          </div>
        )}

        {/* Bottom status */}
        <div className='flex flex-wrap items-center justify-between gap-2 border-t border-primary/15 bg-[oklch(0.07_0.015_25)] px-3 py-2'>
          <p className='font-mono text-[9px] tracking-widest text-zinc-400 uppercase'>
            bash inside arsenal_engine · auth via console session (SSO)
          </p>
          <p className='font-mono text-[9px] tracking-widest text-primary/70 uppercase'>
            Authorized use only
          </p>
        </div>
      </div>
    </div>
  )
}

export default Terminal
