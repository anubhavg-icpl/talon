import { cn } from '@/lib/utils'

/** Placeholder while Three.js chunks load or if WebGL is unavailable. */
const WebGLFallback = ({ className, label = 'INITIALIZING WEBGL…' }: { className?: string; label?: string }) => (
  <div
    className={cn(
      'bg-card/40 border-primary/15 flex h-full min-h-[160px] w-full items-center justify-center rounded-sm border',
      className
    )}
    role='status'
    aria-live='polite'
  >
    <div className='flex flex-col items-center gap-2 px-4 py-8'>
      <div className='border-primary/40 border-t-primary size-6 animate-spin rounded-full border-2' />
      <p className='micro-label text-primary/80'>{label}</p>
    </div>
  </div>
)

export default WebGLFallback
