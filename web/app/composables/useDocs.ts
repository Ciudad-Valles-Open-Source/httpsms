// Loads the documentation Markdown (web/docs-content/{locale}/<slug>.md) at
// build time and renders it for the current i18n locale.

import type { DocPage } from '~/utils/docs-manifest'
import type { RenderedDoc } from '~/utils/docs-markdown'

const DOC_FILES = import.meta.glob<string>('../../docs-content/**/*.md', {
  query: '?raw',
  import: 'default',
  eager: true,
})

const DOC_FALLBACK_LOCALE = 'en'

// "<locale>/<slug>" → raw Markdown
const DOC_SOURCES = new Map<string, string>(
  Object.entries(DOC_FILES).flatMap(([path, source]) => {
    const match = path.match(/docs-content\/([^/]+)\/([^/]+)\.md$/)
    return match ? [[`${match[1]}/${match[2]}`, source] as const] : []
  }),
)

export interface DocView extends RenderedDoc {
  page: DocPage
  /** True when the current locale has no file and English is shown instead. */
  isFallback: boolean
}

/** `{{PLACEHOLDER}}` values for the current locale and runtime config. */
export function useDocPlaceholders() {
  const { locale } = useI18n()
  const config = useRuntimeConfig()
  return computed(() =>
    buildDocPlaceholders({
      apiBaseUrl: config.public.apiBaseUrl,
      appUrl: config.public.appUrl,
      appDownloadUrl: config.public.appDownloadUrl,
      locale: locale.value,
    }),
  )
}

export function useDocPage(slug: MaybeRefOrGetter<string>) {
  const { locale } = useI18n()
  const placeholders = useDocPlaceholders()

  return computed<DocView | null>(() => {
    const page = findDocPage(toValue(slug))
    if (!page) return null

    const localized = DOC_SOURCES.get(`${locale.value}/${page.slug}`)
    const source =
      localized ?? DOC_SOURCES.get(`${DOC_FALLBACK_LOCALE}/${page.slug}`)
    if (source === undefined) return null

    return {
      page,
      isFallback: localized === undefined,
      ...renderDocMarkdown(
        substituteDocPlaceholders(source, placeholders.value),
      ),
    }
  })
}
