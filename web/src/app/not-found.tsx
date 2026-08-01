import Link from 'next/link'

import Logo from '@/components/shared/Logo'
import { Button } from '@/components/ui/button'

const NotFound = () => {
  return (
    <div className='relative flex h-screen w-screen flex-col items-center justify-center gap-8 overflow-hidden p-6'>
      <div className='grid-bg pointer-events-none absolute inset-0 opacity-40' />
      <div className='scanlines pointer-events-none absolute inset-0' />
      <div className='relative z-10 flex flex-col items-center gap-8'>
        <Logo />
        <div className='hud-corners bg-card/80 flex flex-col items-center gap-3 border border-primary/20 px-10 py-8 text-center backdrop-blur-sm'>
          <p className='text-primary text-glow font-mono text-6xl font-bold tracking-widest'>404</p>
          <p className='micro-label'>SIGNAL LOST — ROUTE NOT FOUND</p>
          <p className='text-muted-foreground max-w-sm font-mono text-[11px]'>
            The path you requested is outside scope. Return to the operator overview.
          </p>
        </div>
        <Button
          className='font-mono text-xs tracking-widest uppercase'
          render={<Link href='/overview' />}
          nativeButton={false}
        >
          [ ← Back to Overview ]
        </Button>
      </div>
    </div>
  )
}

export default NotFound
