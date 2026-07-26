// React Imports
import Link from 'next/link'

// Component Imports
import Logo from '@/components/shared/Logo'
import { Button } from '@/components/ui/button'

const NotFound = () => {
  return (
    <div className='grid-bg flex h-screen w-screen flex-col items-center justify-center gap-8 p-6'>
      <Logo />
      <div className='flex flex-col items-center gap-3 text-center'>
        <p className='text-primary font-mono text-6xl font-bold tracking-widest'>404</p>
        <p className='micro-label'>SIGNAL LOST — PAGE NOT FOUND</p>
      </div>
      <Button className='font-mono text-xs tracking-widest uppercase' render={<Link href='/overview' />} nativeButton={false}>
        [ ← Back to Overview ]
      </Button>
    </div>
  )
}

export default NotFound
