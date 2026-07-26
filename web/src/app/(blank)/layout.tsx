// React Imports
import type { ReactNode } from 'react'

const BlankLayout = ({ children }: { children: ReactNode }) => {
  return <div className='bg-background min-h-svh w-full'>{children}</div>
}

export default BlankLayout
