import { readStoredMode } from '~/composables/useThemePreference'

// Applies the remembered light/dark mode before Vuetify renders, so the page
// does not flash the default theme first. The dashboard layout then adds the
// user's palette (see useThemeScope).
export default defineNuxtPlugin({
  name: 'theme-preference',
  enforce: 'pre',
  setup(nuxtApp) {
    nuxtApp.hook('vuetify:before-create', ({ vuetifyOptions }) => {
      const prefersLight =
        window.matchMedia?.('(prefers-color-scheme: light)').matches ?? false
      const mode = readStoredMode() ?? (prefersLight ? 'light' : 'dark')
      vuetifyOptions.theme = {
        ...(typeof vuetifyOptions.theme === 'object'
          ? vuetifyOptions.theme
          : {}),
        defaultTheme: mode,
      }
    })
  },
})
