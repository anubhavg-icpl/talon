import { Suspense } from 'react'

// Component Imports
import NewRun from '@/views/runs/NewRun'

const NewRunPage = () => {
  return (
    <Suspense fallback={<div className='micro-label p-8'>LOADING…</div>}>
      <NewRun />
    </Suspense>
  )
}

export default NewRunPage
