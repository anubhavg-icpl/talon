'use client'

// Hook Imports
import useNow from '@/hooks/use-now'

// Util Imports
import { elapsed } from '@/lib/format'

/**
 * HH:MM:SS elapsed counter from an RFC3339 start time. Ticks live while the run
 * is active; pass `until` (the run's end time) to freeze it at the final
 * duration instead of counting forever after completion.
 */
const Elapsed = ({ since, until, className }: { since: string; until?: string | null; className?: string }) => {
  const now = useNow()
  const end = until ? Date.parse(until) : NaN

  // Frozen: run finished — show the fixed duration, no ticking.
  if (!Number.isNaN(end)) {
    return <span className={className}>{elapsed(since, end)}</span>
  }

  return <span className={className}>{now === null ? '--:--:--' : elapsed(since, now)}</span>
}

export default Elapsed
