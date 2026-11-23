// Navigation structure of the built-in documentation (/docs).
//
// The Markdown for every page lives in `web/docs-content/{es,en}/<slug>.md`.
// Titles shown in the sidebar come from the `docs.*` i18n keys so they follow
// the current locale. This file is plain data (no Nuxt auto-imports) so that
// nuxt.config.ts can import it to prerender the routes.

export interface DocPage {
  slug: string
  titleKey: string
}

export interface DocGroup {
  titleKey: string
  pages: DocPage[]
}

export const DOCS_BASE_PATH = '/docs'

export const DOCS_DEFAULT_SLUG = 'getting-started'

export const DOCS_GROUPS: DocGroup[] = [
  {
    titleKey: 'docs.groups.introduction',
    pages: [{ slug: 'getting-started', titleKey: 'docs.pages.gettingStarted' }],
  },
  {
    titleKey: 'docs.groups.webhooks',
    pages: [
      {
        slug: 'webhooks-introduction',
        titleKey: 'docs.pages.webhooksIntroduction',
      },
      { slug: 'webhooks-events', titleKey: 'docs.pages.webhooksEvents' },
    ],
  },
  {
    titleKey: 'docs.groups.features',
    pages: [
      { slug: 'phone-api-keys', titleKey: 'docs.pages.phoneApiKeys' },
      {
        slug: 'control-sms-send-rate',
        titleKey: 'docs.pages.controlSmsSendRate',
      },
      {
        slug: 'scheduling-sms-messages',
        titleKey: 'docs.pages.schedulingSmsMessages',
      },
      {
        slug: 'outgoing-message-queue',
        titleKey: 'docs.pages.outgoingMessageQueue',
      },
    ],
  },
]

export const DOCS_PAGES: DocPage[] = DOCS_GROUPS.flatMap((group) => group.pages)

export function docPath(slug: string): string {
  return `${DOCS_BASE_PATH}/${slug}`
}

export function findDocPage(slug: string): DocPage | undefined {
  return DOCS_PAGES.find((page) => page.slug === slug)
}

export function adjacentDocPages(slug: string): {
  previous: DocPage | null
  next: DocPage | null
} {
  const index = DOCS_PAGES.findIndex((page) => page.slug === slug)
  if (index === -1) return { previous: null, next: null }
  return {
    previous: DOCS_PAGES[index - 1] ?? null,
    next: DOCS_PAGES[index + 1] ?? null,
  }
}

/** Every docs URL, used to prerender the routes with `nuxt generate`. */
export function docRoutes(): string[] {
  return [DOCS_BASE_PATH, ...DOCS_PAGES.map((page) => docPath(page.slug))]
}
