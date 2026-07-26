'use client'

// React Imports
import { type ComponentType } from 'react'

// Next Imports
import Link from 'next/link'
import { usePathname, useSearchParams } from 'next/navigation'

// Third-party Imports
import * as Icon from 'lucide-react'
import { ChevronRightIcon, SquareArrowOutUpRightIcon } from 'lucide-react'

// Type Imports
import type { MenuGroupSubItem, MenuItem, MenuSubItem } from '@/configs/navConfig'

// Component Imports
import Logo from '@/components/shared/Logo'
import { Collapsible, CollapsibleContent, CollapsibleTrigger } from '@/components/ui/collapsible'
import {
  Sidebar,
  SidebarContent,
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
  SidebarMenuSubItem
} from '@/components/ui/sidebar'

// Config Imports
import { navItems } from '@/configs/navConfig'
import themeConfig from '@/configs/themeConfig'

// Util Imports
import { cn } from '@/lib/utils'

const isSubGroup = (item: MenuSubItem): item is MenuGroupSubItem => 'childItems' in item

const isExternalLink = (href: string) => href.startsWith('http://') || href.startsWith('https://')

function isLinkActive(
  href: string,
  activePath: string | undefined,
  pathname: string,
  searchParams: Pick<URLSearchParams, 'get'>
): boolean {
  if (activePath) {
    return pathname.startsWith(activePath)
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

  return pathname === href
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
        <SidebarGroupLabel className='text-sidebar-foreground/50 tracking-wider uppercase'>
          {groupLabel}
        </SidebarGroupLabel>
      )}
      <SidebarGroupContent>
        <SidebarMenu>
          {data.map(item => {
            const Tag = item.icon ? (Icon[item.icon] as ComponentType) : null

            const isChildActive =
              item.childItems?.some(subItem =>
                isSubGroup(subItem)
                  ? subItem.childItems.some(leaf => isLinkActive(leaf.href, leaf.activePath, pathname, searchParams))
                  : isLinkActive(subItem.href, subItem.activePath, pathname, searchParams)
              ) ?? false

            return item.childItems ? (
              <Collapsible className='group/collapsible' key={item.label}>
                <SidebarMenuItem>
                  <CollapsibleTrigger
                    render={
                      <SidebarMenuButton
                        tooltip={item.label}
                        isActive={isChildActive}
                        className='data-active:bg-primary/5!'
                      />
                    }
                  >
                    {Tag && <Tag />}
                    <span className={cn('min-w-0 flex-1 truncate', item.badge && 'pr-14')}>{item.label}</span>
                    {item.badge && (
                      <SidebarMenuBadge
                        className={cn(
                          'bg-primary/10 max-w-24 truncate rounded-full px-1.5 font-normal',
                          item.badgeClassName
                        )}
                      >
                        {item.badge}
                      </SidebarMenuBadge>
                    )}
                    <ChevronRightIcon className='ml-auto transition-transform duration-200 group-data-open/collapsible:rotate-90' />
                  </CollapsibleTrigger>
                  <CollapsibleContent className='h-(--collapsible-panel-height) overflow-hidden transition-all duration-200 data-ending-style:h-0 data-starting-style:h-0'>
                    <SidebarMenuSub>
                      {item.childItems.map(subItem =>
                        isSubGroup(subItem) ? (
                          <Collapsible className='group/subcollapsible' key={subItem.label}>
                            <SidebarMenuSubItem>
                              <CollapsibleTrigger
                                nativeButton={false}
                                render={
                                  <SidebarMenuSubButton
                                    className='data-active:bg-primary/10! justify-between'
                                    isActive={subItem.childItems.some(leaf =>
                                      isLinkActive(leaf.href, leaf.activePath, pathname, searchParams)
                                    )}
                                  />
                                }
                              >
                                {subItem.label}
                                <ChevronRightIcon className='ml-auto shrink-0 transition-transform duration-200 group-data-open/subcollapsible:rotate-90' />
                              </CollapsibleTrigger>
                              <CollapsibleContent className='h-(--collapsible-panel-height) overflow-hidden transition-all duration-200 data-ending-style:h-0 data-starting-style:h-0'>
                                <SidebarMenuSub className='mx-0'>
                                  {subItem.childItems.map(leaf => (
                                    <SidebarMenuSubItem key={leaf.label}>
                                      <SidebarMenuSubButton
                                        className='data-active:bg-primary/10! justify-between'
                                        render={<Link href={leaf.href} target={leaf.target} />}
                                        isActive={isLinkActive(leaf.href, leaf.activePath, pathname, searchParams)}
                                      >
                                        <span
                                          className={cn(
                                            'min-w-0 flex-1 truncate',
                                            leaf.badge && isExternalLink(leaf.href) && 'pr-8',
                                            leaf.badge && !isExternalLink(leaf.href) && 'pr-14',
                                            !leaf.badge && isExternalLink(leaf.href) && 'pr-6'
                                          )}
                                        >
                                          {leaf.label}
                                        </span>
                                        {leaf.badge && (
                                          <SidebarMenuBadge
                                            className={cn(
                                              'bg-primary/10 max-w-24 truncate rounded-full px-1.5 font-normal',
                                              isExternalLink(leaf.href) && 'right-6',
                                              leaf.badgeClassName
                                            )}
                                          >
                                            {leaf.badge}
                                          </SidebarMenuBadge>
                                        )}
                                        {isExternalLink(leaf.href) && (
                                          <SquareArrowOutUpRightIcon className='ml-auto size-3.5! shrink-0 opacity-50' />
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
                              className='data-active:bg-primary/10! justify-between'
                              render={<Link href={subItem.href} target={subItem.target} />}
                              isActive={isLinkActive(subItem.href, subItem.activePath, pathname, searchParams)}
                            >
                              <span
                                className={cn(
                                  'min-w-0 flex-1 truncate',
                                  subItem.badge && isExternalLink(subItem.href) && 'pr-8',
                                  subItem.badge && !isExternalLink(subItem.href) && 'pr-14',
                                  !subItem.badge && isExternalLink(subItem.href) && 'pr-6'
                                )}
                              >
                                {subItem.label}
                              </span>
                              {subItem.badge && (
                                <SidebarMenuBadge
                                  className={cn(
                                    'bg-primary/10 max-w-24 truncate rounded-full px-1.5 font-normal',
                                    isExternalLink(subItem.href) && 'right-6',
                                    subItem.badgeClassName
                                  )}
                                >
                                  {subItem.badge}
                                </SidebarMenuBadge>
                              )}
                              {isExternalLink(subItem.href) && (
                                <SquareArrowOutUpRightIcon className='ml-auto size-3.5! shrink-0 opacity-50' />
                              )}
                            </SidebarMenuSubButton>
                          </SidebarMenuSubItem>
                        )
                      )}
                    </SidebarMenuSub>
                  </CollapsibleContent>
                </SidebarMenuItem>
              </Collapsible>
            ) : (
              <SidebarMenuItem key={item.label}>
                <SidebarMenuButton
                  tooltip={item.label}
                  render={<Link href={item.href} target={item.target} />}
                  isActive={isLinkActive(item.href, item.activePath, pathname, searchParams)}
                  className='data-active:bg-primary/10!'
                >
                  {Tag && <Tag />}
                  <span
                    className={cn(
                      'min-w-0 flex-1 truncate',
                      item.badge && isExternalLink(item.href) && 'pr-8',
                      item.badge && !isExternalLink(item.href) && 'pr-14',
                      !item.badge && isExternalLink(item.href) && 'pr-6'
                    )}
                  >
                    {item.label}
                  </span>
                  {item.badge && (
                    <SidebarMenuBadge
                      className={cn(
                        'bg-primary/10 max-w-24 truncate rounded-full px-1.5 font-normal',
                        isExternalLink(item.href) && 'right-6',
                        item.badgeClassName
                      )}
                    >
                      {item.badge}
                    </SidebarMenuBadge>
                  )}
                  {isExternalLink(item.href) && (
                    <SquareArrowOutUpRightIcon className='ml-auto size-3.5! shrink-0 opacity-50' />
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

const SidebarLayout = () => {
  const pathname = usePathname()
  const searchParams = useSearchParams()

  return (
    <Sidebar collapsible='icon' variant='sidebar'>
      <SidebarHeader>
        <SidebarMenu>
          <SidebarMenuItem>
            <SidebarMenuButton
              size='lg'
              className='gap-2.5 bg-transparent! hover:bg-transparent!'
              render={<Link href={`${themeConfig.homePageUrl}`} />}
            >
              <Logo compact />
            </SidebarMenuButton>
          </SidebarMenuItem>
        </SidebarMenu>
      </SidebarHeader>
      <SidebarContent className='group-data-[collapsible=icon]:overflow-y-auto'>
        {navItems.map((navItem, index) => {
          return (
            <SidebarGroupedMenuItems
              key={navItem.groupLabel || index}
              data={navItem.items}
              groupLabel={navItem.groupLabel}
              pathname={pathname}
              searchParams={searchParams}
            />
          )
        })}
      </SidebarContent>
    </Sidebar>
  )
}

export default SidebarLayout
