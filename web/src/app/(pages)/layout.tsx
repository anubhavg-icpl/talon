'use client'

// React Imports
import { Suspense } from 'react'
import type { ReactNode } from 'react'

// Component Imports
import Footer from '@/components/layout/Footer'
import Header from '@/components/layout/Header'
import Sidebar from '@/components/layout/Sidebar'
import { SidebarInset } from '@/components/ui/sidebar'
import { Toaster } from '@/components/ui/sonner'

const PagesLayout = ({ children }: Readonly<{ children: ReactNode }>) => {
  return (
    <div className='relative flex h-full w-full min-w-0'>
      {/* No ambient WebGL starfield — avoids multi-canvas pressure site-wide */}
      <div className='pointer-events-none fixed inset-0 z-0 bg-[radial-gradient(ellipse_at_top,oklch(0.25_0.05_25/0.35),transparent_55%)]' />
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
