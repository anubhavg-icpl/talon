// Component Imports
import LiveDot from '@/components/shared/LiveDot'
import { Badge } from '@/components/ui/badge'

// Util Imports
import { cn } from '@/lib/utils'

// Type Imports
import type { RunStatus } from '@/lib/api'

const statusConfig: Record<
  string,
  { label: string; tone: 'green' | 'red' | 'amber' | 'muted'; pulse: boolean; className: string }
> = {
  running: {
    label: 'RUNNING',
    tone: 'green',
    pulse: true,
    className: 'border-primary/40 bg-primary/10 text-primary'
  },
  awaiting_approval: {
    label: 'AWAITING APPROVAL',
    tone: 'amber',
    pulse: true,
    className: 'border-warning/40 bg-warning/10 text-warning'
  },
  completed: {
    label: 'COMPLETED',
    tone: 'green',
    pulse: false,
    className: 'border-primary/40 bg-primary/5 text-primary'
  },
  error: {
    label: 'ERROR',
    tone: 'red',
    pulse: false,
    className: 'border-destructive/40 bg-destructive/10 text-destructive'
  },
  initializing: {
    label: 'INITIALIZING',
    tone: 'green',
    pulse: true,
    className: 'border-primary/40 bg-primary/10 text-primary'
  },
  not_found: {
    label: 'NOT FOUND',
    tone: 'muted',
    pulse: false,
    className: 'border-border bg-muted text-muted-foreground'
  }
}

const fallback = {
  label: 'UNKNOWN',
  tone: 'muted' as const,
  pulse: false,
  className: 'border-border bg-muted text-muted-foreground'
}

const StatusBadge = ({ status, className }: { status: RunStatus | string; className?: string }) => {
  const cfg = statusConfig[status] ?? { ...fallback, label: status.toUpperCase().replace(/_/g, ' ') }

  return (
    <Badge
      variant='outline'
      className={cn('gap-1.5 rounded-sm font-mono text-[10px] tracking-widest uppercase', cfg.className, className)}
    >
      <LiveDot tone={cfg.tone} pulse={cfg.pulse} />
      {cfg.label}
    </Badge>
  )
}

export default StatusBadge
