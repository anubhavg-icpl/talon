'use client'

// Hook Imports
import useNow from '@/hooks/use-now'

/** Live UTC clock. Renders after mount to avoid hydration mismatch. */
const UtcClock = ({ className }: { className?: string }) => {
  const now = useNow()

  return (
    <span className={className}>{now === null ? '---------- --:--:-- UTC' : `${new Date(now).toISOString().slice(0, 19).replace('T', ' ')} UTC`}</span>
  )
}

export default UtcClock
