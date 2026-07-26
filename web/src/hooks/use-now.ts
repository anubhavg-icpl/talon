'use client'

// React Imports
import { useSyncExternalStore } from 'react'

/**
 * Ticking clock that is safe against hydration mismatch:
 * the server snapshot is null, so clocks render a placeholder
 * during SSR/hydration and only tick after mount.
 */
const useNow = (intervalMs = 1000): number | null => {
  return useSyncExternalStore(
    onStoreChange => {
      const id = setInterval(onStoreChange, intervalMs)

      return () => clearInterval(id)
    },
    () => Date.now(),
    () => null
  )
}

export default useNow
