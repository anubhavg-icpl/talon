'use client'

// React Imports
import { Suspense } from 'react'
import type { ReactNode } from 'react'

// Component Imports
import Footer from '@/components/layout/Footer'
import Header from '@/components/layout/Header'
import Sidebar from '@/components/layout/Sidebar'
import Starfield from '@/components/shared/three/Starfield'
import { SidebarInset } from '@/components/ui/sidebar'
import { Toaster } from '@/components/ui/sonner'

const PagesLayout = ({ children }: Readonly<{ children: ReactNode }>) => {
  return (
    <div className='relative flex h-full w-full min-w-0'>
      {/* App-wide Three.js ambient starfield (all authenticated pages) */}
      <Starfield className='fixed inset-0 z-0 opacity-40' count={600} opacity={0.4} speed={0.00028} />
      <div className='relative z-10 flex h-full w-full min-w-0'>
        <Suspense>
          <Sidebar />
        </Suspense>
        <SidebarInset className='bg-background/70 flex flex-1 flex-col backdrop-blur-[1px]'>
          <Header />
          <main className='operator-stage mx-auto size-full max-w-360 flex-1 px-4 py-6 sm:px-6'>{children}</main>
          <Toaster />
          <Footer />
        </SidebarInset>
      </div>
    </div>
  )
}

export default PagesLayout
