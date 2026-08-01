'use client'

import { useEffect } from 'react'
import Link from 'next/link'

import { Button } from '@/components/ui/button'

/** Route error boundary — stops a single view crash from killing the shell. */
export default function PagesError({
  error,
  reset
}: {
  error: Error & { digest?: string }
  reset: () => void
}) {
  useEffect(() => {
    console.error('[talon pages error]', error)
  }, [error])

  return (
    <div className='flex min-h-[50vh] flex-col items-center justify-center gap-6 p-8 text-center'>
      <p className='text-primary font-mono text-4xl font-bold tracking-widest'>FAULT</p>
      <p className='micro-label text-muted-foreground max-w-md'>
        This view failed to render (often WebGL / memory). The rest of the console is still up.
      </p>
      <p className='text-muted-foreground max-w-lg font-mono text-[11px] break-all'>{error.message}</p>
      <div className='flex flex-wrap justify-center gap-2'>
        <Button type='button' onClick={reset} className='font-mono text-xs tracking-widest uppercase'>
          Retry view
        </Button>
        <Link
          href='/overview'
          className='border-input bg-background hover:bg-accent inline-flex h-9 items-center rounded-md border px-3 font-mono text-xs tracking-widest uppercase'
        >
          Overview
        </Link>
        <Link
          href='/assist'
          className='border-input bg-background hover:bg-accent inline-flex h-9 items-center rounded-md border px-3 font-mono text-xs tracking-widest uppercase'
        >
          SLM Assist
        </Link>
        <Link
          href='/login'
          className='border-input bg-background hover:bg-accent inline-flex h-9 items-center rounded-md border px-3 font-mono text-xs tracking-widest uppercase'
        >
          Login
        </Link>
      </div>
    </div>
  )
}
