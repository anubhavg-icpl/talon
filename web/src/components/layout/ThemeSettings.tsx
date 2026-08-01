'use client'

// React Imports
import { useState } from 'react'

// Third-party Imports
import { CheckIcon, PaletteIcon } from 'lucide-react'

// Component Imports
import { Button } from '@/components/ui/button'
import { Popover, PopoverContent, PopoverTrigger } from '@/components/ui/popover'

// Util Imports
import { THEME_PRESETS, getStoredPreset, setStoredPreset, type PresetId } from '@/lib/theme-presets'
import { cn } from '@/lib/utils'

/**
 * Accent theme switcher — the console's take on the shadcn-admin preferences
 * popover. Swaps the --talon accent hue (primary/ring/charts/glow/HUD) across
 * a set of on-brand presets; Talon green is the default. Persisted in
 * localStorage and applied pre-paint by the boot script in the root layout.
 */
const ThemeSettings = () => {
  // Lazy init reads localStorage on the client; the popup only renders on open
  // (portalled), so there is no SSR/hydration mismatch on the checkmark.
  const [preset, setPreset] = useState<PresetId>(() => getStoredPreset())

  const choose = (id: PresetId) => {
    setStoredPreset(id)
    setPreset(id)
  }

  return (
    <Popover>
      <PopoverTrigger render={<Button variant='ghost' size='icon' aria-label='Accent theme' />}>
        <PaletteIcon className='size-4' />
      </PopoverTrigger>
      <PopoverContent align='end' className='w-56 gap-2'>
        <p className='micro-label'>ACCENT THEME</p>
        <div className='flex flex-col gap-1'>
          {THEME_PRESETS.map(p => (
            <button
              key={p.id}
              onClick={() => choose(p.id)}
              className={cn(
                'flex items-center gap-2.5 rounded-sm px-2 py-1.5 font-mono text-xs tracking-widest uppercase transition-colors',
                p.id === preset ? 'bg-primary/10 text-primary' : 'hover:bg-muted/50'
              )}
            >
              <span
                className='size-3.5 shrink-0 rounded-full border border-white/10'
                style={{ background: p.swatch }}
              />
              <span className='flex-1 text-left'>{p.label}</span>
              {p.id === preset && <CheckIcon className='size-3.5 shrink-0' />}
            </button>
          ))}
        </div>
      </PopoverContent>
    </Popover>
  )
}

export default ThemeSettings
