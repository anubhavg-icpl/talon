'use client'

// Hook Imports
import useNow from '@/hooks/use-now'

// Util Imports
import { elapsed } from '@/lib/format'

/** Ticking HH:MM:SS elapsed counter from an RFC3339 start time. Renders after mount to avoid hydration mismatch. */
const Elapsed = ({ since, className }: { since: string; className?: string }) => {
  const now = useNow()

  return <span className={className}>{now === null ? '--:--:--' : elapsed(since, now)}</span>
}

export default Elapsed
