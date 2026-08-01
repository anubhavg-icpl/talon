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
      <CrosshairGlyph className='text-primary size-5 shrink-0 drop-shadow-[0_0_8px_var(--op-glow)]' />
      <span className='text-sm font-semibold tracking-[0.2em] whitespace-nowrap'>
        TALON <span className='text-primary text-glow'>{'//'}</span>{' '}
        <span className='text-primary/90'>{compact ? 'OP' : 'OPERATOR'}</span>
      </span>
    </span>
  )
}

export default Logo
