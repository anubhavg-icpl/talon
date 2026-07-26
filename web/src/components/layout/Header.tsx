'use client'

// React Imports
import { Fragment, useEffect, useState } from 'react'

import { usePathname } from 'next/navigation'

// Component Imports
import LiveDot from '@/components/shared/LiveDot'
import {
  Breadcrumb,
  BreadcrumbItem,
  BreadcrumbLink,
  BreadcrumbList,
  BreadcrumbPage,
  BreadcrumbSeparator
} from '@/components/ui/breadcrumb'
import { Separator } from '@/components/ui/separator'
import { SidebarTrigger } from '@/components/ui/sidebar'

// Util Imports
import { health, logout, me } from '@/lib/api'

const Header = () => {
  const pathname = usePathname()

  const segments = pathname.split('/').filter(Boolean)

  const [coreUp, setCoreUp] = useState<boolean | null>(null)
  const [username, setUsername] = useState<string | null>(null)

  useEffect(() => {
    let mounted = true

    const check = () =>
      health()
        .then(() => mounted && setCoreUp(true))
        .catch(() => mounted && setCoreUp(false))

    check()
    const id = setInterval(check, 10000)

    me()
      .then(res => mounted && setUsername(res.auth === 'disabled' ? null : res.username))
      .catch(() => {})

    return () => {
      mounted = false
      clearInterval(id)
    }
  }, [])

  const handleLogout = async () => {
    try {
      await logout()
    } catch {
      // session already gone — redirect anyway
    }

    window.location.href = '/login'
  }

  return (
    <header className='bg-card sticky top-0 z-50 border-b'>
      <div className='mx-auto flex max-w-360 items-center justify-between gap-6 px-4 py-2 sm:px-6'>
        <div className='flex items-center gap-4'>
          <SidebarTrigger className='[&_svg]:size-5!' />
          <Separator orientation='vertical' className='hidden h-4! data-vertical:self-center sm:block' />
          <Breadcrumb className='hidden sm:block'>
            <BreadcrumbList>
              {segments.map((segment, index) => {
                const isLast = index === segments.length - 1
                const label = segment.replace(/-/g, ' ').replace(/\b\w/g, c => c.toUpperCase())
                const href = '/' + segments.slice(0, index + 1).join('/')

                return (
                  <Fragment key={href}>
                    <BreadcrumbItem>
                      {isLast ? <BreadcrumbPage>{label}</BreadcrumbPage> : <BreadcrumbLink>{label}</BreadcrumbLink>}
                    </BreadcrumbItem>
                    {!isLast && <BreadcrumbSeparator />}
                  </Fragment>
                )
              })}
            </BreadcrumbList>
          </Breadcrumb>
        </div>
        <div className='flex items-center gap-4 font-mono text-[10px] tracking-widest uppercase'>
          {coreUp === null ? (
            <span className='text-muted-foreground'>CORE …</span>
          ) : coreUp ? (
            <span className='text-primary flex items-center gap-2'>
              <LiveDot tone='green' /> CORE ONLINE
            </span>
          ) : (
            <span className='text-destructive flex items-center gap-2'>
              <LiveDot tone='red' /> CORE OFFLINE
            </span>
          )}
          {username && (
            <span className='text-muted-foreground hidden items-center gap-3 sm:flex'>
              <span className='text-foreground'>@{username}</span>
              <button onClick={handleLogout} className='hover:text-destructive cursor-pointer transition-colors'>
                [ LOGOUT ]
              </button>
            </span>
          )}
        </div>
      </div>
    </header>
  )
}

export default Header
