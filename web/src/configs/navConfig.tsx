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
        icon: 'SquareTerminal',
        label: 'Arsenal Shell',
        href: '/terminal'
      },
      {
        icon: 'Bug',
        label: 'Findings',
        href: '/findings'
      },
      {
        icon: 'GitCompare',
        label: 'Compare',
        href: '/compare'
      },
      {
        icon: 'Shield',
        label: 'Engagements',
        href: '/ops'
      },
      {
        icon: 'ShieldCheck',
        label: 'Approvals',
        href: '/approvals',
        activePath: '/approvals'
      },
      {
        icon: 'Share2',
        label: 'Sharing',
        href: '/engagements',
        activePath: '/engagements'
      }
    ]
  },
  {
    groupLabel: 'Intelligence',
    items: [
      {
        icon: 'Sparkles',
        label: 'SLM Assist',
        href: '/assist',
        badge: 'LOCAL'
      },
      {
        icon: 'Bot',
        label: 'Agents',
        href: '/agents'
      },
      {
        icon: 'BookOpen',
        label: 'Skills',
        href: '/skills'
      },
      {
        icon: 'KeyRound',
        label: 'Crypto Toolkit',
        href: '/crypto'
      },
      {
        icon: 'ShieldAlert',
        label: 'Alert Triage',
        href: '/detection'
      },
      {
        icon: 'ScrollText',
        label: 'Playbooks',
        href: '/playbooks'
      },
      {
        icon: 'BookCopy',
        label: 'Blueprints',
        href: '/blueprints',
        activePath: '/blueprints'
      },
      {
        icon: 'Radio',
        label: 'Intel Feed',
        href: '/intel'
      }
    ]
  },
  {
    groupLabel: 'System',
    items: [
      {
        icon: 'History',
        label: 'Audit Trail',
        href: '/audit',
        activePath: '/audit'
      },
      {
        icon: 'KeySquare',
        label: 'Gatekeepers',
        href: '/gatekeepers',
        activePath: '/gatekeepers'
      },
      {
        icon: 'Settings',
        label: 'Settings',
        href: '/settings'
      }
    ]
  }
]
