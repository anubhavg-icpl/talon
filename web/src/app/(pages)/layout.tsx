'use client'

// React Imports
import { Suspense } from 'react'
import type { ReactNode } from 'react'

// Component Imports
import Footer from '@/components/layout/Footer'
import Header from '@/components/layout/Header'
import Sidebar from '@/components/layout/Sidebar'
import { Starfield } from '@/components/shared/three'
import { SidebarInset } from '@/components/ui/sidebar'
import { Toaster } from '@/components/ui/sonner'

const PagesLayout = ({ children }: Readonly<{ children: ReactNode }>) => {
  return (
    <div className='relative flex h-full w-full min-w-0'>
      {/* Light ambient stars — fewer particles to avoid GPU/tab OOM on tight hosts */}
      <Starfield className='fixed inset-0 z-0 opacity-25' count={160} opacity={0.28} speed={0.00018} />
      <div className='relative z-10 flex h-full w-full min-w-0'>
        <Suspense fallback={null}>
          <Sidebar />
        </Suspense>
        <SidebarInset className='bg-background/75 flex flex-1 flex-col backdrop-blur-[2px]'>
          <Header />
          <main className='operator-stage mx-auto size-full max-w-360 flex-1 px-4 py-6 sm:px-6'>{children}</main>
          <Toaster position='bottom-right' richColors closeButton />
          <Footer />
        </SidebarInset>
      </div>
    </div>
  )
}

export default PagesLayout
