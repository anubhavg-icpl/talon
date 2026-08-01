// Accent theme presets. Each maps to a `data-theme-preset` value that overrides
// the --talon-l/c/h hue triple in globals.css. "talon" (green) is the default
// and sets no attribute.

export const THEME_PRESETS = [
  { id: 'talon', label: 'TALON', swatch: 'oklch(0.87 0.2 145)' },
  { id: 'amber', label: 'AMBER', swatch: 'oklch(0.8 0.15 85)' },
  { id: 'violet', label: 'VIOLET', swatch: 'oklch(0.72 0.18 300)' },
  { id: 'crimson', label: 'CRIMSON', swatch: 'oklch(0.7 0.2 20)' }
] as const

export type PresetId = (typeof THEME_PRESETS)[number]['id']

export const PRESET_STORAGE_KEY = 'talon-theme-preset'

const isPreset = (v: unknown): v is PresetId => THEME_PRESETS.some(p => p.id === v)

/** Set/clear the document attribute that activates a preset. */
export const applyPreset = (id: PresetId) => {
  if (id === 'talon') delete document.documentElement.dataset.themePreset
  else document.documentElement.dataset.themePreset = id
}

export const getStoredPreset = (): PresetId => {
  if (typeof window === 'undefined') return 'talon'

  const v = window.localStorage.getItem(PRESET_STORAGE_KEY)

  return isPreset(v) ? v : 'talon'
}

export const setStoredPreset = (id: PresetId) => {
  window.localStorage.setItem(PRESET_STORAGE_KEY, id)
  applyPreset(id)
}
