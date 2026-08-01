// Third-party Imports
import type * as Icon from 'lucide-react'

type IconName = keyof typeof Icon

export type MenuLeafSubItem = {
  label: string
  href: string
  activePath?: string
  badge?: string
  badgeClassName?: string
  target?: '_blank' | '_self' | '_parent' | '_top'
}

export type MenuGroupSubItem = {
  label: string
  childItems: MenuLeafSubItem[]
}

export type MenuSubItem = MenuLeafSubItem | MenuGroupSubItem

export type MenuItem = {
  icon: IconName
  label: string
} & (
  | {
      href: string
      activePath?: string
      badge?: string
      badgeClassName?: string
      childItems?: never
      target?: '_blank' | '_self' | '_parent' | '_top'
    }
  | {
      href?: never
      badge?: string
      badgeClassName?: string
      childItems: MenuSubItem[]
    }
)

export type NavItem = {
  groupLabel?: string
  items: MenuItem[]
}

export const navItems: NavItem[] = [
  {
    groupLabel: 'Console',
    items: [
      {
        icon: 'Radar',
        label: 'Overview',
        href: '/overview'
      },
      {
        icon: 'Clapperboard',
        label: 'Showcase',
        href: '/showcase'
      }
    ]
  },
  {
    groupLabel: 'Operations',
    items: [
      {
        icon: 'Terminal',
        label: 'Runs',
        href: '/runs',
        activePath: '/runs'
      },
      {
        icon: 'Crosshair',
        label: 'New Operation',
        href: '/runs/new'
      },
      {
        icon: 'Bug',
        label: 'Findings',
        href: '/findings'
      }
    ]
  },
  {
    groupLabel: 'Intelligence',
    items: [
      {
        icon: 'Bot',
        label: 'Agents',
        href: '/agents'
      },
      {
        icon: 'BookOpen',
        label: 'Skills',
        href: '/skills'
      }
    ]
  },
  {
    groupLabel: 'System',
    items: [
      {
        icon: 'Settings',
        label: 'Settings',
        href: '/settings'
      }
    ]
  }
]
