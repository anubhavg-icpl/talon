// React Imports
import type { ReactNode } from 'react'

// Util Imports
import { cn } from '@/lib/utils'

/**
 * Standard page title block — one h1 + optional micro-label subtitle, with an
 * optional right-aligned action slot. Keeps every top-level view visually
 * aligned (typography, spacing, baseline) without repeating the markup.
 */
const PageHeader = ({
  title,
  subtitle,
  action,
  className
}: {
  title: ReactNode
  subtitle?: ReactNode
  action?: ReactNode
  className?: string
}) => {
  return (
    <div className={cn('flex flex-wrap items-end justify-between gap-3', className)}>
      <div className='min-w-0'>
        <h1 className='font-mono text-xl font-semibold tracking-widest'>{title}</h1>
        {subtitle && <p className='micro-label mt-1'>{subtitle}</p>}
      </div>
      {action}
    </div>
  )
}

export default PageHeader
