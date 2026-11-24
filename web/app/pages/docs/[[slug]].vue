<script setup lang="ts">
import { mdiArrowLeft, mdiArrowRight, mdiTranslate } from '@mdi/js'

definePageMeta({ layout: 'docs' })

const route = useRoute()
const router = useRouter()
const { t } = useI18n()

const slug = computed(
  () => (route.params.slug as string | undefined) || DOCS_DEFAULT_SLUG,
)
const doc = useDocPage(slug)
const adjacent = computed(() => {
  const { previous, next } = adjacentDocPages(slug.value)
  return {
    previous: previous && { ...previous, to: docPath(previous.slug) },
    next: next && { ...next, to: docPath(next.slug) },
  }
})
const docsHome = DOCS_BASE_PATH

const metaTitle = computed(() =>
  doc.value
    ? `${doc.value.title} - ${t('docs.metaTitleSuffix')}`
    : `${t('docs.notFoundTitle')} - ${t('docs.metaTitleSuffix')}`,
)
const metaDescription = computed(
  () => doc.value?.description || t('docs.metaDescription'),
)

const { assetUrl } = usePublicUrls()

useSeoMeta({
  title: metaTitle,
  description: metaDescription,
  ogTitle: metaTitle,
  ogDescription: metaDescription,
  ogImage: assetUrl('/header.png'),
  twitterCard: 'summary_large_image',
})

// Code windows with several languages: show the clicked tab's panel.
function selectCodeTab(event: MouseEvent): boolean {
  const tab = (event.target as HTMLElement | null)?.closest<HTMLElement>(
    '[data-code-tab]',
  )
  const codeWindow = tab?.closest('.code-window')
  if (!tab || !codeWindow) return false
  const selected = tab.dataset.codeTab
  codeWindow.querySelectorAll<HTMLElement>('[data-code-tab]').forEach((el) => {
    const active = el.dataset.codeTab === selected
    el.classList.toggle('is-active', active)
    el.setAttribute('aria-selected', String(active))
  })
  codeWindow
    .querySelectorAll<HTMLElement>('[data-code-panel]')
    .forEach((el) => {
      el.hidden = el.dataset.codePanel !== selected
    })
  return true
}

// Links inside the rendered Markdown are plain <a> tags: route internal ones
// through the router so they don't reload the SPA.
function onArticleClick(event: MouseEvent) {
  if (selectCodeTab(event)) return
  if (
    event.defaultPrevented ||
    event.button !== 0 ||
    event.metaKey ||
    event.ctrlKey ||
    event.shiftKey ||
    event.altKey
  ) {
    return
  }
  const anchor = (event.target as HTMLElement | null)?.closest('a')
  const href = anchor?.getAttribute('href')
  if (!anchor || !href || anchor.target === '_blank') return

  if (href.startsWith('#')) {
    event.preventDefault()
    router.push({ hash: href })
  } else if (href.startsWith('/') && !href.startsWith('//')) {
    event.preventDefault()
    router.push(href)
  }
}

// The article renders synchronously, so scroll to a deep-linked heading once
// it exists.
onMounted(() => {
  if (!route.hash) return
  nextTick(() => {
    document
      .getElementById(decodeURIComponent(route.hash.slice(1)))
      ?.scrollIntoView()
  })
})
</script>

<template>
  <VContainer class="docs-page py-8 px-4 px-md-8">
    <template v-if="doc">
      <VAlert
        v-if="doc.isFallback"
        type="info"
        variant="tonal"
        density="compact"
        class="mb-6"
        :icon="mdiTranslate"
      >
        {{ $t('docs.fallbackNotice') }}
      </VAlert>

      <!-- eslint-disable-next-line vue/no-v-html -- trusted, build-time Markdown from web/docs-content -->
      <article class="docs-article" @click="onArticleClick" v-html="doc.html" />

      <VDivider class="mt-12 mb-6" />
      <VRow>
        <VCol cols="12" sm="6">
          <VCard
            v-if="adjacent.previous"
            :to="adjacent.previous.to"
            class="pa-5 h-100"
          >
            <div class="text-body-small text-medium-emphasis">
              <VIcon size="small" :icon="mdiArrowLeft" />
              {{ $t('docs.previous') }}
            </div>
            <div class="text-title-medium font-weight-bold text-primary mt-1">
              {{ $t(adjacent.previous.titleKey) }}
            </div>
          </VCard>
        </VCol>
        <VCol cols="12" sm="6">
          <VCard
            v-if="adjacent.next"
            :to="adjacent.next.to"
            class="pa-5 h-100 text-right"
          >
            <div class="text-body-small text-medium-emphasis">
              {{ $t('docs.next') }}
              <VIcon size="small" :icon="mdiArrowRight" />
            </div>
            <div class="text-title-medium font-weight-bold text-primary mt-1">
              {{ $t(adjacent.next.titleKey) }}
            </div>
          </VCard>
        </VCol>
      </VRow>
    </template>

    <div v-else class="text-center py-16">
      <h1 class="text-display-small mb-4">{{ $t('docs.notFoundTitle') }}</h1>
      <p class="text-body-large text-medium-emphasis mb-8">
        {{ $t('docs.notFoundText') }}
      </p>
      <VBtn color="primary" :to="docsHome">{{ $t('docs.backToDocs') }}</VBtn>
    </div>
  </VContainer>
</template>

<style scoped>
.docs-page {
  max-width: 920px;
}

.docs-page :deep(.v-card) {
  transition: border-color 0.15s;
}

.docs-page :deep(.v-card:hover) {
  border-color: rgb(var(--v-theme-primary), 0.5);
}

.docs-article {
  font-size: 1rem;
  line-height: 1.7;
  overflow-wrap: break-word;
}

.docs-article :deep(h1),
.docs-article :deep(h2),
.docs-article :deep(h3),
.docs-article :deep(h4) {
  line-height: 1.25;
  font-weight: 700;
  letter-spacing: -0.015em;
  scroll-margin-top: 80px;
}

.docs-article :deep(h1) {
  font-size: 2.6rem;
  font-weight: 800;
  letter-spacing: -0.03em;
  margin-bottom: 1rem;
}

.docs-article :deep(h2) {
  font-size: 1.75rem;
  margin: 2.5rem 0 1rem;
  padding-bottom: 0.4rem;
  border-bottom: 1px solid rgb(var(--v-border-color), var(--v-border-opacity));
}

.docs-article :deep(h3) {
  font-size: 1.35rem;
  margin: 2rem 0 0.75rem;
}

.docs-article :deep(h4) {
  font-size: 1.1rem;
  margin: 1.5rem 0 0.5rem;
}

.docs-article :deep(.docs-heading-anchor) {
  color: rgb(var(--v-theme-primary));
  text-decoration: none;
  opacity: 0;
  margin-left: 0.25rem;
  transition: opacity 0.15s;
}

.docs-article :deep(h2:hover .docs-heading-anchor),
.docs-article :deep(h3:hover .docs-heading-anchor),
.docs-article :deep(h4:hover .docs-heading-anchor),
.docs-article :deep(.docs-heading-anchor:focus) {
  opacity: 1;
}

.docs-article :deep(p),
.docs-article :deep(ul),
.docs-article :deep(ol),
.docs-article :deep(pre),
.docs-article :deep(blockquote),
.docs-article :deep(.docs-table) {
  margin-bottom: 1rem;
}

.docs-article :deep(ul),
.docs-article :deep(ol) {
  padding-left: 1.5rem;
}

.docs-article :deep(li + li),
.docs-article :deep(li > ul),
.docs-article :deep(li > ol) {
  margin-top: 0.25rem;
}

.docs-article :deep(li > ul),
.docs-article :deep(li > ol) {
  margin-bottom: 0;
}

.docs-article :deep(a) {
  color: rgb(var(--v-theme-primary));
  text-decoration: none;
}

.docs-article :deep(a:hover) {
  text-decoration: underline;
}

.docs-article :deep(hr) {
  border: 0;
  border-top: 1px solid rgb(var(--v-border-color), var(--v-border-opacity));
  margin: 2rem 0;
}

.docs-article :deep(img) {
  max-width: 100%;
  border-radius: 4px;
}

.docs-article :deep(:not(pre) > code) {
  font-size: 0.875em;
  padding: 0.15em 0.4em;
  border-radius: 4px;
  background: rgb(var(--v-theme-on-surface), 0.08);
  color: rgb(var(--v-theme-on-surface));
}

.docs-article :deep(.code-window) {
  margin: 1.25rem 0;
  overflow: hidden;
  border: 1px solid rgb(255 255 255 / 10%);
  border-radius: var(--card-radius, 1rem);
  background: #0d1117;
  box-shadow: var(--card-shadow);
}

.docs-article :deep(.code-window-bar) {
  display: flex;
  gap: 0.4rem;
  padding: 0.8rem 1rem 0.4rem;
}

.docs-article :deep(.code-window-bar span) {
  width: 0.7rem;
  height: 0.7rem;
  border-radius: 50%;
  background: rgb(255 255 255 / 18%);
}

.docs-article :deep(.code-window-tabs) {
  display: flex;
  overflow-x: auto;
  border-bottom: 1px solid rgb(255 255 255 / 10%);
  scrollbar-width: none;
}

.docs-article :deep(.code-window-tab) {
  display: inline-flex;
  flex: none;
  align-items: center;
  gap: 0.45rem;
  padding: 0.6rem 1.1rem;
  border: 0;
  border-bottom: 2px solid transparent;
  background: transparent;
  color: #e6edf3;
  font: inherit;
  font-size: 0.875rem;
  font-weight: 700;
}

.docs-article :deep(button.code-window-tab) {
  cursor: pointer;
}

.docs-article :deep(button.code-window-tab:hover) {
  background: rgb(255 255 255 / 4%);
}

.docs-article :deep(.code-window-tab.is-active) {
  border-bottom-color: rgb(var(--v-theme-primary));
  color: rgb(var(--v-theme-primary));
}

.docs-article :deep(.code-window-tab svg) {
  width: 1.1rem;
  height: 1.1rem;
}

.docs-article :deep(pre.docs-code) {
  margin: 0;
  overflow-x: auto;
  padding: 1.25rem 1.5rem;
  font-size: 0.875rem;
  line-height: 1.6;
  color: #e6edf3;
  background: transparent;
}

.docs-article :deep(pre.docs-code[hidden]) {
  display: none;
}

.docs-article :deep(pre.docs-code code) {
  background: transparent;
  padding: 0;
  font-size: inherit;
}

.docs-article :deep(.docs-table) {
  overflow-x: auto;
  border: 1px solid rgb(var(--v-border-color), var(--v-border-opacity));
  border-radius: 0.9rem;
}

.docs-article :deep(table) {
  width: 100%;
  border-collapse: collapse;
  font-size: 0.925rem;
}

.docs-article :deep(th),
.docs-article :deep(td) {
  padding: 0.6rem 0.9rem;
  text-align: left;
  vertical-align: top;
  border-bottom: 1px solid rgb(var(--v-border-color), var(--v-border-opacity));
}

.docs-article :deep(th) {
  font-weight: 600;
  background: rgb(var(--v-theme-on-surface), 0.05);
}

.docs-article :deep(tr:last-child td) {
  border-bottom: 0;
}

.docs-article :deep(blockquote.docs-callout) {
  padding: 0.75rem 1rem;
  border-left: 4px solid rgb(var(--callout-color));
  border-radius: 0 0.75rem 0.75rem 0;
  background: rgb(var(--callout-color), 0.08);

  --callout-color: var(--v-theme-on-surface);
}

.docs-article :deep(blockquote.docs-callout-info) {
  --callout-color: var(--v-theme-info);
}

.docs-article :deep(blockquote.docs-callout-tip) {
  --callout-color: var(--v-theme-success);
}

.docs-article :deep(blockquote.docs-callout-warning) {
  --callout-color: var(--v-theme-warning);
}

.docs-article :deep(blockquote > :last-child) {
  margin-bottom: 0;
}

@media (width <= 600px) {
  .docs-article :deep(h1) {
    font-size: 2rem;
  }

  .docs-article :deep(h2) {
    font-size: 1.5rem;
  }
}
</style>
