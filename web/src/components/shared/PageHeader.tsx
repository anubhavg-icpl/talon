import type { ReactNode } from 'react'

import { cn } from '@/lib/utils'

// ponytail: accept BOTH prop conventions so merged callers compile without
// rewriting either side -- origin views pass description/actions, local views
// (Terminal, Runs, NewRun, Settings) pass subtitle/action.
type PageHeaderProps = {
  title: ReactNode
  description?: ReactNode
  subtitle?: ReactNode
  actions?: ReactNode
  action?: ReactNode
  className?: string
}

/** Consistent operator page chrome. */
const PageHeader = ({ title, description, subtitle, actions, action, className }: PageHeaderProps) => {
  const desc = description ?? subtitle
  const acts = actions ?? action
  return (
    <div className={cn('flex flex-wrap items-end justify-between gap-4', className)}>
      <div className='min-w-0 space-y-1.5'>
        <h1 className='font-mono text-xl font-semibold tracking-[0.18em] text-foreground sm:text-2xl'>{title}</h1>
        {desc ? <p className='micro-label text-muted-foreground max-w-2xl'>{desc}</p> : null}
      </div>
      {acts ? <div className='flex flex-wrap items-center gap-2'>{acts}</div> : null}
    </div>
  )
}

export default PageHeader
