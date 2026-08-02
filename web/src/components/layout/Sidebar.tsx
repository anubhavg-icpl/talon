'use client'

/**
 * Talon operator sidebar — shadcn Sidebar composition (collapsible icon mode).
 * Structure: SidebarProvider (Providers) → Sidebar + Rail | SidebarInset (layout)
 */

import { type ComponentType } from 'react'

import Link from 'next/link'
import { usePathname, useSearchParams } from 'next/navigation'

import * as Icon from 'lucide-react'
import { ChevronRightIcon, CrosshairIcon, SettingsIcon, SquareArrowOutUpRightIcon } from 'lucide-react'

import type { MenuGroupSubItem, MenuItem, MenuSubItem } from '@/configs/navConfig'
import { Collapsible, CollapsibleContent, CollapsibleTrigger } from '@/components/ui/collapsible'
import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarGroup,
  SidebarGroupContent,
  SidebarGroupLabel,
  SidebarHeader,
  SidebarMenu,
  SidebarMenuBadge,
  SidebarMenuButton,
  SidebarMenuItem,
  SidebarMenuSub,
  SidebarMenuSubButton,
  SidebarMenuSubItem,
  SidebarRail,
  useSidebar
} from '@/components/ui/sidebar'
import GlobeWallpaper from '@/components/shared/GlobeWallpaper'
import { navItems } from '@/configs/navConfig'
import themeConfig from '@/configs/themeConfig'
import { cn } from '@/lib/utils'

const isSubGroup = (item: MenuSubItem): item is MenuGroupSubItem => 'childItems' in item

const isExternalLink = (href: string) => href.startsWith('http://') || href.startsWith('https://')

/** Dedicated routes under /runs that have their own nav item (must not light up Runs). */
const RUNS_SIBLING_PREFIXES = ['/runs/new']

function isRunsListPath(pathname: string): boolean {
  if (pathname === '/runs') return true
  // /runs/:id detail only — not /runs/new
  if (!pathname.startsWith('/runs/')) return false
  if (RUNS_SIBLING_PREFIXES.some(p => pathname === p || pathname.startsWith(p + '/'))) return false
  return true
}

function isLinkActive(
  href: string,
  activePath: string | undefined,
  pathname: string,
  searchParams: Pick<URLSearchParams, 'get'>
): boolean {
  // Exact match always wins
  if (pathname === href) return true

  if (activePath) {
    // Prefer exact activePath; for prefixes avoid sibling collisions (/runs vs /runs/new)
    if (pathname === activePath) return true
    if (activePath === '/runs' || href === '/runs') return isRunsListPath(pathname)
    return pathname.startsWith(activePath + '/')
  }

  if (href.includes('?')) {
    const [hrefPath, hrefQuery] = href.split('?')
    if (pathname !== hrefPath) return false
    const hrefParams = new URLSearchParams(hrefQuery)
    for (const [key, value] of hrefParams.entries()) {
      if (searchParams.get(key) !== value) return false
    }
    return true
  }

  // Nested routes: /runs/:id highlights Runs; /runs/new does not
  if (href === '/runs') return isRunsListPath(pathname)
  if (href !== '/' && pathname.startsWith(href + '/')) return true
  return false
}

const SidebarGroupedMenuItems = ({
  data,
  groupLabel,
  pathname,
  searchParams
}: {
  data: MenuItem[]
  groupLabel?: string
  pathname: string
  searchParams: Pick<URLSearchParams, 'get'>
}) => {
  return (
    <SidebarGroup>
      {groupLabel && (
        <SidebarGroupLabel className='text-sidebar-foreground/45 font-mono text-[10px] tracking-[0.18em] uppercase group-data-[collapsible=icon]:hidden'>
          {groupLabel}
        </SidebarGroupLabel>
      )}
      <SidebarGroupContent>
        <SidebarMenu>
          {data.map(item => {
            const Tag = item.icon ? (Icon[item.icon] as ComponentType<{ className?: string }>) : null

            const isChildActive =
              item.childItems?.some(subItem =>
                isSubGroup(subItem)
                  ? subItem.childItems.some(leaf => isLinkActive(leaf.href, leaf.activePath, pathname, searchParams))
                  : isLinkActive(subItem.href, subItem.activePath, pathname, searchParams)
              ) ?? false

            if (item.childItems) {
              return (
                <Collapsible key={item.label} className='group/collapsible' defaultOpen={isChildActive}>
                  <SidebarMenuItem>
                    <CollapsibleTrigger
                      render={
                        <SidebarMenuButton
                          tooltip={item.label}
                          isActive={isChildActive}
                          className='data-active:bg-primary/10 data-active:text-primary'
                        />
                      }
                    >
                      {Tag && <Tag className='size-4 shrink-0' />}
                      <span className={cn('min-w-0 flex-1 truncate', item.badge && 'pr-12')}>{item.label}</span>
                      {item.badge && (
                        <SidebarMenuBadge
                          className={cn(
                            'bg-primary/15 text-primary max-w-20 truncate rounded-sm px-1.5 font-mono text-[9px] tracking-wide',
                            item.badgeClassName
                          )}
                        >
                          {item.badge}
                        </SidebarMenuBadge>
                      )}
                      <ChevronRightIcon className='ml-auto size-4 shrink-0 transition-transform duration-200 group-data-open/collapsible:rotate-90 group-data-[collapsible=icon]:hidden' />
                    </CollapsibleTrigger>
                    <CollapsibleContent className='h-(--collapsible-panel-height) overflow-hidden transition-all duration-200 data-ending-style:h-0 data-starting-style:h-0 group-data-[collapsible=icon]:hidden'>
                      <SidebarMenuSub>
                        {item.childItems.map(subItem =>
                          isSubGroup(subItem) ? (
                            <Collapsible key={subItem.label} className='group/subcollapsible' defaultOpen>
                              <SidebarMenuSubItem>
                                <CollapsibleTrigger
                                  nativeButton={false}
                                  render={
                                    <SidebarMenuSubButton
                                      className='data-active:bg-primary/10 data-active:text-primary justify-between'
                                      isActive={subItem.childItems.some(leaf =>
                                        isLinkActive(leaf.href, leaf.activePath, pathname, searchParams)
                                      )}
                                    />
                                  }
                                >
                                  {subItem.label}
                                  <ChevronRightIcon className='ml-auto size-3.5 shrink-0 transition-transform duration-200 group-data-open/subcollapsible:rotate-90' />
                                </CollapsibleTrigger>
                                <CollapsibleContent className='h-(--collapsible-panel-height) overflow-hidden transition-all duration-200 data-ending-style:h-0 data-starting-style:h-0'>
                                  <SidebarMenuSub className='mx-0'>
                                    {subItem.childItems.map(leaf => (
                                      <SidebarMenuSubItem key={leaf.label}>
                                        <SidebarMenuSubButton
                                          className='data-active:bg-primary/10 data-active:text-primary justify-between'
                                          render={<Link href={leaf.href} target={leaf.target} />}
                                          isActive={isLinkActive(leaf.href, leaf.activePath, pathname, searchParams)}
                                        >
                                          <span className='min-w-0 flex-1 truncate'>{leaf.label}</span>
                                          {leaf.badge && (
                                            <SidebarMenuBadge className='bg-primary/15 text-primary rounded-sm px-1.5 font-mono text-[9px]'>
                                              {leaf.badge}
                                            </SidebarMenuBadge>
                                          )}
                                          {isExternalLink(leaf.href) && (
                                            <SquareArrowOutUpRightIcon className='ml-auto size-3.5 shrink-0 opacity-50' />
                                          )}
                                        </SidebarMenuSubButton>
                                      </SidebarMenuSubItem>
                                    ))}
                                  </SidebarMenuSub>
                                </CollapsibleContent>
                              </SidebarMenuSubItem>
                            </Collapsible>
                          ) : (
                            <SidebarMenuSubItem key={subItem.label}>
                              <SidebarMenuSubButton
                                className='data-active:bg-primary/10 data-active:text-primary justify-between'
                                render={<Link href={subItem.href} target={subItem.target} />}
                                isActive={isLinkActive(subItem.href, subItem.activePath, pathname, searchParams)}
                              >
                                <span className='min-w-0 flex-1 truncate'>{subItem.label}</span>
                                {subItem.badge && (
                                  <SidebarMenuBadge
                                    className={cn(
                                      'bg-primary/15 text-primary max-w-20 truncate rounded-sm px-1.5 font-mono text-[9px]',
                                      subItem.badgeClassName
                                    )}
                                  >
                                    {subItem.badge}
                                  </SidebarMenuBadge>
                                )}
                                {isExternalLink(subItem.href) && (
                                  <SquareArrowOutUpRightIcon className='ml-auto size-3.5 shrink-0 opacity-50' />
                                )}
                              </SidebarMenuSubButton>
                            </SidebarMenuSubItem>
                          )
                        )}
                      </SidebarMenuSub>
                    </CollapsibleContent>
                  </SidebarMenuItem>
                </Collapsible>
              )
            }

            return (
              <SidebarMenuItem key={item.label}>
                <SidebarMenuButton
                  tooltip={item.label}
                  render={<Link href={item.href} target={item.target} />}
                  isActive={isLinkActive(item.href, item.activePath, pathname, searchParams)}
                  className='data-active:bg-primary/10 data-active:text-primary'
                >
                  {Tag && <Tag className='size-4 shrink-0' />}
                  <span
                    className={cn(
                      'min-w-0 flex-1 truncate',
                      item.badge && isExternalLink(item.href) && 'pr-8',
                      item.badge && !isExternalLink(item.href) && 'pr-12',
                      !item.badge && isExternalLink(item.href) && 'pr-6'
                    )}
                  >
                    {item.label}
                  </span>
                  {item.badge && (
                    <SidebarMenuBadge
                      className={cn(
                        'bg-primary/15 text-primary max-w-20 truncate rounded-sm px-1.5 font-mono text-[9px] tracking-wide',
                        isExternalLink(item.href) && 'right-6',
                        item.badgeClassName
                      )}
                    >
                      {item.badge}
                    </SidebarMenuBadge>
                  )}
                  {isExternalLink(item.href) && (
                    <SquareArrowOutUpRightIcon className='ml-auto size-3.5 shrink-0 opacity-50' />
                  )}
                </SidebarMenuButton>
              </SidebarMenuItem>
            )
          })}
        </SidebarMenu>
      </SidebarGroupContent>
    </SidebarGroup>
  )
}

const SidebarBrand = () => {
  return (
    <SidebarMenu>
      <SidebarMenuItem>
        <SidebarMenuButton
          size='lg'
          tooltip='Talon Operator'
          className='hover:bg-sidebar-accent/80 data-[slot=sidebar-menu-button]:p-2!'
          render={<Link href={themeConfig.homePageUrl} />}
        >
          <div className='bg-primary/15 text-primary ring-primary/25 flex aspect-square size-8 shrink-0 items-center justify-center rounded-md ring-1'>
            <CrosshairIcon className='size-4 drop-shadow-[0_0_8px_var(--op-glow)]' />
          </div>
          <div className='grid min-w-0 flex-1 text-left text-sm leading-tight'>
            <span className='truncate font-mono text-[12px] font-semibold tracking-[0.16em] uppercase'>
              Talon <span className='text-primary'>//</span> Op
            </span>
            <span className='text-muted-foreground truncate font-mono text-[10px] tracking-widest uppercase'>
              Red operator shell
            </span>
          </div>
        </SidebarMenuButton>
      </SidebarMenuItem>
    </SidebarMenu>
  )
}

const SidebarLayout = () => {
  const pathname = usePathname()
  const searchParams = useSearchParams()
  const { state } = useSidebar()

  return (
    <Sidebar collapsible='icon' variant='sidebar' side='left' className='border-sidebar-border'>
      <SidebarHeader className='border-sidebar-border/80 border-b'>
        <SidebarBrand />
      </SidebarHeader>

      <SidebarContent className='gap-0 overflow-x-hidden'>
        {navItems.map((navItem, index) => (
          <SidebarGroupedMenuItems
            key={navItem.groupLabel || index}
            data={navItem.items}
            groupLabel={navItem.groupLabel}
            pathname={pathname}
            searchParams={searchParams}
          />
        ))}
      </SidebarContent>

      <SidebarFooter className='border-sidebar-border/80 gap-2 border-t'>
        {/* Same dark-planet WebM as Overview Target HUD */}
        <div className='group-data-[collapsible=icon]:hidden px-2 pt-1'>
          <GlobeWallpaper
            compact
            label='TALON · GLOBE'
            allow3d={false}
            className='max-h-40 w-full rounded-md border-primary/20 shadow-[0_0_24px_-8px_var(--op-glow)]'
          />
        </div>
        {/* Icon-collapsed: tiny looping planet */}
        <div className='hidden items-center justify-center px-1 py-1 group-data-[collapsible=icon]:flex'>
          <GlobeWallpaper
            compact
            minimal
            allow3d={false}
            className='size-8 max-h-8 rounded-md border-primary/30'
          />
        </div>
        <SidebarMenu>
          <SidebarMenuItem>
            <SidebarMenuButton
              tooltip='Settings'
              isActive={pathname.startsWith('/settings')}
              render={<Link href='/settings' />}
              className='data-active:bg-primary/10 data-active:text-primary'
            >
              <SettingsIcon className='size-4 shrink-0' />
              <span>Settings</span>
            </SidebarMenuButton>
          </SidebarMenuItem>
        </SidebarMenu>
        <div className='text-muted-foreground group-data-[collapsible=icon]:hidden px-2 pb-1 font-mono text-[9px] tracking-widest uppercase'>
          {state === 'expanded' ? '⌘B collapse' : '⌘B expand'} · authorized only
        </div>
      </SidebarFooter>

      <SidebarRail className='hover:after:bg-primary/40' />
    </Sidebar>
  )
}

export default SidebarLayout
