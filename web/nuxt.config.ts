// https://nuxt.com/docs/api/configuration/nuxt-config
import { docRoutes } from './app/utils/docs-manifest'

// Public URL of this self-hosted dashboard (used for SEO / og:image).
const appUrl = (process.env.APP_URL || 'http://localhost:3000').replace(
  /\/+$/,
  '',
)

export default defineNuxtConfig({
  compatibilityDate: '2025-01-01',

  ssr: false,

  modules: [
    '@nuxt/eslint',
    'vuetify-nuxt-module',
    '@pinia/nuxt',
    '@nuxtjs/google-fonts',
    '@nuxtjs/seo',
    '@nuxtjs/i18n',
  ],

  googleFonts: {
    families: {
      'Plus Jakarta Sans': [400, 500, 600, 700, 800],
    },
    display: 'swap',
    download: true,
  },

  css: ['vuetify/styles', 'v-phone-input/styles', '~/assets/css/main.css'],

  i18n: {
    compilation: {
      strictMessage: false,
    },
    locales: [
      { code: 'en', name: 'English', files: ['en.json', 'vuetify-en.ts'] },
      { code: 'es', name: 'Español', files: ['es.json', 'vuetify-es.ts'] },
    ],
    defaultLocale: 'en',
    lazy: true,
    langDir: 'locales/',
    strategy: 'no_prefix',
    detectBrowserLanguage: {
      useCookie: true,
      cookieKey: 'i18n_redirected',
      redirectOn: 'root',
    },
  },

  site: {
    url: appUrl,
    name: process.env.APP_NAME || 'httpSMS',
    description:
      'Turn your Android phone into an SMS gateway. Send and receive text messages worldwide with a simple HTTP API — no SMS provider or short code required.',
    defaultLocale: 'en',
  },

  // Authenticated app routes that should never appear in search engines or the
  // sitemap. The marketing/blog/legal pages remain public and indexable.
  robots: {
    disallow: [
      '/messages',
      '/threads',
      '/settings',
      '/billing',
      '/bulk-messages',
      '/heartbeats',
      '/phone-api-keys',
      '/search-messages',
      '/admin',
      '/account-pending',
    ],
  },

  sitemap: {
    urls: docRoutes(),
    exclude: [
      '/messages',
      '/threads',
      '/threads/**',
      '/settings',
      '/billing',
      '/bulk-messages',
      '/heartbeats/**',
      '/phone-api-keys',
      '/search-messages',
      '/login',
      '/admin/**',
      '/account-pending',
    ],
  },

  // The app ships as a client-rendered SPA (ssr: false) served statically, so
  // runtime Satori OG-image generation is not available. Curated static OG
  // images are set per page instead.
  ogImage: {
    enabled: false,
  },

  // Static download assets in /public are valid at runtime but are not Nuxt
  // routes, so exclude them from link validation to avoid false positives.
  linkChecker: {
    excludeLinks: ['/templates/**'],
  },

  build: {
    transpile: ['vuetify', 'chart.js', 'vue-chartjs', 'v-phone-input'],
  },

  vite: {
    define: {
      'process.env.DEBUG': false,
    },
    optimizeDeps: {
      include: [
        '@mdi/js',
        'chartjs-adapter-moment',
        'date-fns',
        'firebase/app',
        'firebase/auth',
        'highlight.js/lib/core',
        'libphonenumber-js',
        'pusher-js',
        'qrcode',
      ],
    },
  },

  vuetify: {
    moduleOptions: {
      prefixComposables: true,
    },
    vuetifyOptions: {
      // Visual language shared with the author's other products: violet → blue
      // brand gradient, Plus Jakarta Sans, soft bordered cards and a
      // dark-first palette with a light alternative (see ThemeToggle.vue).
      theme: {
        defaultTheme: 'dark',
        themes: {
          // Default palette: blues of the author's Arquos products (Galatea /
          // Glauqus). The violet palette below stays available in the
          // appearance picker (see useThemePreference.ts).
          dark: {
            dark: true,
            colors: {
              background: '#0f111a',
              'on-background': '#e3e8f4',
              surface: '#161a27',
              'on-surface': '#e3e8f4',
              'surface-bright': '#1f2434',
              'surface-light': '#252b3d',
              'surface-variant': '#2a3146',
              'on-surface-variant': '#eef2fb',
              primary: '#5b8def',
              'on-primary': '#06112b',
              secondary: '#4fd1ff',
              'on-secondary': '#04202b',
              error: '#ff6b81',
              'on-error': '#2a0710',
              info: '#4fd1ff',
              'on-info': '#04202b',
              success: '#34d399',
              'on-success': '#05261a',
              warning: '#fbbf24',
              'on-warning': '#2a1d02',
            },
            variables: {
              'border-color': '#ffffff',
              'border-opacity': 0.1,
              'medium-emphasis-opacity': 0.7,
              'theme-code': '#1f2434',
              'theme-on-code': '#e3e8f4',
            },
          },
          light: {
            dark: false,
            colors: {
              background: '#f2f7ff',
              'on-background': '#1f2d52',
              surface: '#ffffff',
              'on-surface': '#1f2d52',
              'surface-bright': '#ffffff',
              'surface-light': '#e8effc',
              'surface-variant': '#25396f',
              'on-surface-variant': '#f2f7ff',
              primary: '#435ebe',
              'on-primary': '#ffffff',
              secondary: '#0050c8',
              'on-secondary': '#ffffff',
              error: '#d92d4b',
              info: '#0050c8',
              success: '#13865a',
              warning: '#b86e00',
            },
            variables: {
              'border-color': '#25396f',
              'border-opacity': 0.12,
              'medium-emphasis-opacity': 0.68,
              'theme-code': '#e8effc',
              'theme-on-code': '#1f2d52',
            },
          },
          violetDark: {
            dark: true,
            colors: {
              background: '#0f0c16',
              'on-background': '#ebe6f4',
              surface: '#16121f',
              'on-surface': '#ebe6f4',
              'surface-bright': '#221c2e',
              'surface-light': '#2a2338',
              'surface-variant': '#2d2640',
              'on-surface-variant': '#efeaf7',
              primary: '#b671ff',
              'on-primary': '#14081f',
              secondary: '#6e9bff',
              'on-secondary': '#0a1230',
              error: '#ff6b81',
              'on-error': '#2a0710',
              info: '#6e9bff',
              'on-info': '#0a1230',
              success: '#34d399',
              'on-success': '#05261a',
              warning: '#fbbf24',
              'on-warning': '#2a1d02',
            },
            variables: {
              'border-color': '#ffffff',
              'border-opacity': 0.1,
              'medium-emphasis-opacity': 0.68,
              'theme-code': '#221c2e',
              'theme-on-code': '#ebe6f4',
            },
          },
          violetLight: {
            dark: false,
            colors: {
              background: '#f6f3fb',
              'on-background': '#1d1828',
              surface: '#ffffff',
              'on-surface': '#1d1828',
              'surface-bright': '#ffffff',
              'surface-light': '#efeaf7',
              'surface-variant': '#2d2640',
              'on-surface-variant': '#f4f0fa',
              primary: '#7a1fc2',
              'on-primary': '#ffffff',
              secondary: '#2f6bf0',
              'on-secondary': '#ffffff',
              error: '#d92d4b',
              info: '#2f6bf0',
              success: '#13865a',
              warning: '#b86e00',
            },
            variables: {
              'border-color': '#1e143c',
              'border-opacity': 0.12,
              'medium-emphasis-opacity': 0.66,
              'theme-code': '#efeaf7',
              'theme-on-code': '#1d1828',
            },
          },
        },
      },
      defaults: {
        VBtn: {
          rounded: 'lg',
          class: 'text-none font-weight-bold',
          style: 'letter-spacing: 0',
        },
        VCard: { rounded: 'xl', elevation: 0, border: true },
        VAlert: { rounded: 'lg' },
        VTextField: { rounded: 'lg', color: 'primary' },
        VTextarea: { rounded: 'lg', color: 'primary' },
        VSelect: { rounded: 'lg', color: 'primary' },
        VAutocomplete: { rounded: 'lg', color: 'primary' },
        VCombobox: { rounded: 'lg', color: 'primary' },
        VFileInput: { rounded: 'lg', color: 'primary' },
        VSnackbar: { rounded: 'lg' },
        VAppBar: { flat: true, color: 'background', border: 'b' },
        VTabs: { color: 'primary' },
        VExpansionPanels: { flat: true },
        // Rows turn into label/value cards on phones.
        VDataTable: { mobileBreakpoint: 'sm' },
        VDataTableServer: { mobileBreakpoint: 'sm' },
      },
      icons: {
        defaultSet: 'mdi-svg',
      },
    },
  },

  runtimeConfig: {
    public: {
      apiBaseUrl: process.env.API_BASE_URL || 'http://localhost:8000',
      clientVersion: process.env.GITHUB_RUN_ID || 'dev',
      appUrl,
      appName: process.env.APP_NAME || 'HTTP SMS',
      appGithubUrl:
        process.env.APP_GITHUB_URL ||
        'https://github.com/Ciudad-Valles-Open-Source/httpsms',
      appDownloadUrl:
        process.env.APP_DOWNLOAD_URL ||
        'https://github.com/Ciudad-Valles-Open-Source/httpsms/releases/latest/download/HttpSms.apk',
      appEnv: process.env.APP_ENV || 'production',
      checkoutUrl: process.env.CHECKOUT_URL || '',
      enterpriseCheckoutUrl: process.env.ENTERPRISE_CHECKOUT_URL || '',
      cloudflareTurnstileSiteKey:
        process.env.CLOUDFLARE_TURNSTILE_SITE_KEY || '',
      soketiKey: process.env.SOKETI_KEY || '',
      soketiCluster: process.env.SOKETI_CLUSTER || '',
      firebaseApiKey: process.env.FIREBASE_API_KEY || '',
      firebaseAuthDomain: process.env.FIREBASE_AUTH_DOMAIN || '',
      firebaseProjectId: process.env.FIREBASE_PROJECT_ID || '',
      firebaseStorageBucket: process.env.FIREBASE_STORAGE_BUCKET || '',
      firebaseMessagingSenderId: process.env.FIREBASE_MESSAGING_SENDER_ID || '',
      firebaseAppId: process.env.FIREBASE_APP_ID || '',
      firebaseMeasurementId: process.env.FIREBASE_MEASUREMENT_ID || '',
      useSoketi: process.env.USE_SOKETI === 'true',
    },
  },

  nitro: {
    prerender: {
      // The /docs/[[slug]] pages are dynamic routes, so list them explicitly.
      routes: docRoutes(),
      failOnError: false,
    },
  },

  routeRules: {
    '/messages': { robots: false },
    '/threads': { robots: false },
    '/threads/**': { robots: false },
    '/settings': { robots: false },
    '/billing': { robots: false },
    '/bulk-messages': { robots: false },
    '/heartbeats/**': { robots: false },
    '/phone-api-keys': { robots: false },
    '/search-messages': { robots: false },
    '/admin/**': { robots: false },
    '/account-pending': { robots: false },
  },

  app: {
    head: {
      titleTemplate: '%s',
      title: 'Convert your android phone into an SMS gateway - httpSMS',
      htmlAttrs: { lang: 'en' },
      script: [
        {
          src: 'https://challenges.cloudflare.com/turnstile/v0/api.js?render=explicit',
        },
      ],
      meta: [
        { charset: 'utf-8' },
        { name: 'viewport', content: 'width=device-width, initial-scale=1' },
        {
          name: 'description',
          content:
            'Use your android phone to send and receive SMS messages using a simple HTTP API.',
        },
        { name: 'format-detection', content: 'telephone=no' },
        { name: 'twitter:card', content: 'summary_large_image' },
        {
          property: 'og:title',
          content: 'Convert your android phone into an SMS gateway - httpSMS',
        },
        {
          property: 'og:description',
          content:
            'Use your android phone to send and receive SMS messages using a simple HTTP API.',
        },
        {
          property: 'og:image',
          content: `${appUrl}/header.png`,
        },
      ],
      link: [{ rel: 'icon', type: 'image/x-icon', href: '/favicon.ico' }],
    },
  },
})
