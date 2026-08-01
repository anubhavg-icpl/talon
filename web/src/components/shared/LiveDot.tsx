// Util Imports
import { cn } from '@/lib/utils'

type Tone = 'green' | 'red' | 'amber' | 'cyan' | 'muted'

const toneClasses: Record<Tone, string> = {
  green: 'bg-primary',
  red: 'bg-destructive',
  amber: 'bg-warning',
  cyan: 'bg-primary',
  muted: 'bg-muted-foreground'
}

const LiveDot = ({ tone = 'green', pulse = true, className }: { tone?: Tone; pulse?: boolean; className?: string }) => {
  return (
    <span className={cn('relative inline-flex size-2 shrink-0', className)}>
      {pulse && (
        <span className={cn('absolute inline-flex size-full animate-ping rounded-full opacity-60', toneClasses[tone])} />
      )}
      <span className={cn('relative inline-flex size-2 rounded-full', toneClasses[tone])} />
    </span>
  )
}

export default LiveDot
