import { useTheme } from 'vuetify'

export type ThemeMode = 'dark' | 'light'
export type ThemePalette = 'blue' | 'violet'

export const THEME_STORAGE_KEY = 'httpsms-theme'
export const PALETTE_STORAGE_KEY = 'httpsms-palette'

export const THEME_PALETTES: {
  value: ThemePalette
  label: string
  swatch: [string, string]
}[] = [
  {
    value: 'blue',
    label: 'components.theme.paletteBlue',
    swatch: ['#435ebe', '#4fd1ff'],
  },
  {
    value: 'violet',
    label: 'components.theme.paletteViolet',
    swatch: ['#7a1fc2', '#6e9bff'],
  },
]

/** Vuetify theme name for a mode/palette pair (see nuxt.config.ts). */
export function themeName(mode: ThemeMode, palette: ThemePalette): string {
  if (palette === 'violet')
    return mode === 'dark' ? 'violetDark' : 'violetLight'
  return mode
}

function readStorage(key: string): string | null {
  if (typeof window === 'undefined') return null
  try {
    return window.localStorage.getItem(key)
  } catch {
    return null
  }
}

function writeStorage(key: string, value: string) {
  try {
    window.localStorage.setItem(key, value)
  } catch {
    // Storage can be unavailable (private mode); the change still applies.
  }
}

export function readStoredMode(): ThemeMode | null {
  const value = readStorage(THEME_STORAGE_KEY)
  return value === 'light' || value === 'dark' ? value : null
}

export function readStoredPalette(): ThemePalette {
  return readStorage(PALETTE_STORAGE_KEY) === 'violet' ? 'violet' : 'blue'
}

export type ThemeScope = 'public' | 'dashboard'

/**
 * Appearance preference: light or dark mode plus one of the preset color
 * palettes (blue by default, violet as an alternative). Both choices are
 * remembered in localStorage; without a saved mode the app follows the
 * operating system and falls back to dark.
 *
 * The palette only applies to the dashboard. Public pages (landing, docs,
 * login) always use the default blue palette and only follow the mode. Each
 * layout declares its scope with `useThemeScope()`.
 */
export function useThemePreference() {
  const theme = useTheme()
  const scope = useState<ThemeScope>('theme-scope', () => 'public')
  const palette = useState<ThemePalette>('theme-palette', readStoredPalette)

  const isDark = computed(() => theme.global.current.value.dark)
  const mode = computed<ThemeMode>(() => (isDark.value ? 'dark' : 'light'))
  const activePalette = computed<ThemePalette>(() =>
    scope.value === 'dashboard' ? palette.value : 'blue',
  )
  /** Dark variant of the active palette, for always-dark areas (code). */
  const darkThemeName = computed(() => themeName('dark', activePalette.value))

  function apply(nextMode: ThemeMode = mode.value) {
    theme.change(themeName(nextMode, activePalette.value))
  }

  function setMode(value: ThemeMode) {
    writeStorage(THEME_STORAGE_KEY, value)
    apply(value)
  }

  function setPalette(value: ThemePalette) {
    palette.value = value
    writeStorage(PALETTE_STORAGE_KEY, value)
    apply()
  }

  function setScope(value: ThemeScope) {
    scope.value = value
    apply()
  }

  function toggleTheme() {
    setMode(isDark.value ? 'light' : 'dark')
  }

  return {
    isDark,
    mode,
    palette,
    darkThemeName,
    setMode,
    setPalette,
    setScope,
    toggleTheme,
  }
}

/** Called by each layout: the dashboard uses the user's palette. */
export function useThemeScope(scope: ThemeScope) {
  const { setScope } = useThemePreference()
  setScope(scope)
}
