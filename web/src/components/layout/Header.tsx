'use client'

// React Imports
import { Fragment, useEffect, useState } from 'react'

import Link from 'next/link'
import { usePathname } from 'next/navigation'

// Third-party Imports
import { ChevronDownIcon, LogOutIcon, SettingsIcon } from 'lucide-react'

// Component Imports
import CommandPalette from '@/components/layout/CommandPalette'
import ThemeSettings from '@/components/layout/ThemeSettings'
import LiveDot from '@/components/shared/LiveDot'
import {
  Breadcrumb,
  BreadcrumbItem,
  BreadcrumbLink,
  BreadcrumbList,
  BreadcrumbPage,
  BreadcrumbSeparator
} from '@/components/ui/breadcrumb'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger
} from '@/components/ui/dropdown-menu'
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
        <div className='flex items-center gap-3 font-mono text-[10px] tracking-widest uppercase'>
          <CommandPalette />
          <ThemeSettings />
          <Separator orientation='vertical' className='hidden h-4! data-vertical:self-center sm:block' />
          {coreUp === null ? (
            <span className='text-muted-foreground'>CORE …</span>
          ) : coreUp ? (
            <span className='text-primary flex items-center gap-2'>
              <LiveDot tone='green' /> <span className='hidden sm:inline'>CORE ONLINE</span>
            </span>
          ) : (
            <span className='text-destructive flex items-center gap-2'>
              <LiveDot tone='red' /> <span className='hidden sm:inline'>CORE OFFLINE</span>
            </span>
          )}
          {username && (
            <DropdownMenu>
              <DropdownMenuTrigger className='text-muted-foreground hover:text-foreground flex items-center gap-1.5 font-mono text-[10px] tracking-widest uppercase transition-colors'>
                <span className='text-foreground'>@{username}</span>
                <ChevronDownIcon className='size-3' />
              </DropdownMenuTrigger>
              <DropdownMenuContent align='end' className='min-w-44'>
                <DropdownMenuLabel className='micro-label'>OPERATOR</DropdownMenuLabel>
                <DropdownMenuItem
                  render={<Link href='/settings' />}
                  className='font-mono text-xs tracking-widest uppercase'
                >
                  <SettingsIcon /> SYSTEM
                </DropdownMenuItem>
                <DropdownMenuSeparator />
                <DropdownMenuItem
                  variant='destructive'
                  onClick={handleLogout}
                  className='font-mono text-xs tracking-widest uppercase'
                >
                  <LogOutIcon /> LOGOUT
                </DropdownMenuItem>
              </DropdownMenuContent>
            </DropdownMenu>
          )}
        </div>
      </div>
    </header>
  )
}

export default Header
