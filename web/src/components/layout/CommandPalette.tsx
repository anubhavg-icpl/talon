'use client'

// React Imports
import { type ComponentType, useCallback, useEffect, useMemo, useState } from 'react'

// Next Imports
import { useRouter } from 'next/navigation'

// Third-party Imports
import * as Icon from 'lucide-react'
import { CornerDownLeftIcon, SearchIcon } from 'lucide-react'

// Component Imports
import { Dialog, DialogContent, DialogDescription, DialogTitle } from '@/components/ui/dialog'

// Config Imports
import { navItems } from '@/configs/navConfig'

// Util Imports
import { cn } from '@/lib/utils'

type Command = { label: string; href: string; group: string; icon: keyof typeof Icon }

// Flatten the nav config into a single searchable command list.
const COMMANDS: Command[] = navItems.flatMap(group =>
  group.items.flatMap(item =>
    'href' in item && item.href
      ? [{ label: item.label, href: item.href, group: (group.groupLabel ?? 'MENU').toUpperCase(), icon: item.icon }]
      : []
  )
)

/**
 * ⌘K / Ctrl+K command palette — keyboard-first jump-to navigation over the nav
 * config. Built on the existing base-ui Dialog (no extra dependency). Arrow keys
 * move the selection, Enter navigates, Esc closes.
 */
const CommandPalette = () => {
  const router = useRouter()
  const [open, setOpen] = useState(false)
  const [query, setQuery] = useState('')
  const [active, setActive] = useState(0)

  const results = useMemo(() => {
    const q = query.trim().toLowerCase()

    if (!q) return COMMANDS

    return COMMANDS.filter(c => c.label.toLowerCase().includes(q) || c.group.toLowerCase().includes(q))
  }, [query])

  // Event-driven open/close that resets search state on open (no state-in-effect).
  const setOpenReset = useCallback((next: boolean) => {
    if (next) {
      setQuery('')
      setActive(0)
    }

    setOpen(next)
  }, [])

  // Global ⌘K / Ctrl+K toggle.
  useEffect(() => {
    const onKey = (e: KeyboardEvent) => {
      if ((e.metaKey || e.ctrlKey) && e.key.toLowerCase() === 'k') {
        e.preventDefault()
        setOpenReset(!open)
      }
    }

    window.addEventListener('keydown', onKey)

    return () => window.removeEventListener('keydown', onKey)
  }, [open, setOpenReset])

  const go = useCallback(
    (href: string) => {
      setOpen(false)
      router.push(href)
    },
    [router]
  )

  const onKeyDown = (e: React.KeyboardEvent) => {
    if (e.key === 'ArrowDown') {
      e.preventDefault()
      setActive(a => Math.min(a + 1, results.length - 1))
    } else if (e.key === 'ArrowUp') {
      e.preventDefault()
      setActive(a => Math.max(a - 1, 0))
    } else if (e.key === 'Enter') {
      e.preventDefault()
      const sel = results[active]

      if (sel) go(sel.href)
    }
  }

  return (
    <>
      <button
        onClick={() => setOpen(true)}
        className='border-border bg-muted/40 text-muted-foreground hover:border-primary/40 hover:text-foreground flex items-center gap-2 rounded-sm border px-2 py-1 font-mono text-xs tracking-widest uppercase transition-colors'
        aria-label='Open command palette'
      >
        <SearchIcon className='size-3.5' />
        <span className='hidden sm:inline'>SEARCH</span>
        <kbd className='border-border ml-1 hidden rounded-sm border px-1 py-px text-[10px] sm:inline'>⌘K</kbd>
      </button>

      <Dialog open={open} onOpenChange={setOpen}>
        <DialogContent
          showCloseButton={false}
          className='top-[12%] w-full max-w-lg translate-y-0 gap-0 overflow-hidden p-0 sm:max-w-lg'
        >
          <DialogTitle className='sr-only'>Command palette</DialogTitle>
          <DialogDescription className='sr-only'>Jump to a page</DialogDescription>

          <div className='flex items-center gap-2 border-b px-3'>
            <SearchIcon className='text-muted-foreground size-4 shrink-0' />
            <input
              autoFocus
              value={query}
              onChange={e => {
                setQuery(e.target.value)
                setActive(0)
              }}
              onKeyDown={onKeyDown}
              placeholder='jump to…'
              spellCheck={false}
              className='placeholder:text-muted-foreground/60 h-11 w-full bg-transparent font-mono text-sm outline-none'
            />
            <kbd className='micro-label border-border rounded-sm border px-1.5 py-0.5'>ESC</kbd>
          </div>

          <div className='max-h-80 overflow-y-auto p-1'>
            {results.length === 0 ? (
              <p className='micro-label px-3 py-8 text-center'>NO MATCHES</p>
            ) : (
              results.map((c, i) => {
                const Tag = Icon[c.icon] as ComponentType<{ className?: string }> | undefined

                return (
                  <button
                    key={c.href}
                    onMouseEnter={() => setActive(i)}
                    onClick={() => go(c.href)}
                    className={cn(
                      'flex w-full items-center gap-3 rounded-sm px-3 py-2 text-left font-mono text-sm transition-colors',
                      i === active ? 'bg-primary/10 text-primary' : 'text-foreground'
                    )}
                  >
                    {Tag && <Tag className='size-4 shrink-0' />}
                    <span className='flex-1 truncate'>{c.label}</span>
                    <span className='micro-label shrink-0'>{c.group}</span>
                    {i === active && <CornerDownLeftIcon className='size-3.5 shrink-0 opacity-60' />}
                  </button>
                )
              })
            )}
          </div>
        </DialogContent>
      </Dialog>
    </>
  )
}

export default CommandPalette
