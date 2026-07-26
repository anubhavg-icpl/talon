// Util Imports
import { cn } from '@/lib/utils'

const CrosshairGlyph = ({ className }: { className?: string }) => (
  <svg viewBox='0 0 24 24' fill='none' stroke='currentColor' strokeWidth='1.5' className={className} aria-hidden>
    <circle cx='12' cy='12' r='7' />
    <circle cx='12' cy='12' r='1.5' fill='currentColor' stroke='none' />
    <path d='M12 1v5M12 18v5M1 12h5M18 12h5' />
  </svg>
)

const Logo = ({ className, compact }: { className?: string; compact?: boolean }) => {
  return (
    <span className={cn('inline-flex items-center gap-2 font-mono', className)}>
      <CrosshairGlyph className='text-primary size-5 shrink-0' />
      <span className='text-sm font-semibold tracking-widest whitespace-nowrap'>
        TALON <span className='text-primary'>{'//'}</span> {compact ? 'OPS' : 'OPS CONSOLE'}
      </span>
    </span>
  )
}

export default Logo
