import type { ReactNode } from 'react'

import { cn } from '@/lib/utils'

type PageHeaderProps = {
  title: string
  description?: string
  actions?: ReactNode
  className?: string
}

/** Consistent operator page chrome. */
const PageHeader = ({ title, description, actions, className }: PageHeaderProps) => {
  return (
    <div className={cn('flex flex-wrap items-end justify-between gap-4', className)}>
      <div className='min-w-0 space-y-1.5'>
        <h1 className='font-mono text-xl font-semibold tracking-[0.18em] text-foreground sm:text-2xl'>{title}</h1>
        {description ? <p className='micro-label text-muted-foreground max-w-2xl'>{description}</p> : null}
      </div>
      {actions ? <div className='flex flex-wrap items-center gap-2'>{actions}</div> : null}
    </div>
  )
}

export default PageHeader
