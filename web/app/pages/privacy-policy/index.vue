<script setup lang="ts">
definePageMeta({
  layout: 'website',
})

const { t } = useI18n()

useSeoMeta({
  title: computed(() => t('legal.privacy.metaTitle')),
  description: computed(() => t('legal.privacy.metaDescription')),
  ogTitle: computed(() => t('legal.privacy.metaTitle')),
  ogDescription: computed(() => t('legal.privacy.metaDescription')),
})

const ISSUES_URL = 'https://github.com/Ciudad-Valles-Open-Source/httpsms/issues'

interface LegalSection {
  id: string
  /** Number of paragraphs (p1..pN) shown before the list. */
  paragraphs: number
  /** Keys of a list rendered as "label: text" (items.<key>.label/text). */
  labeledItems?: string[]
  /** Number of plain list items (items.i1..iN). */
  items?: number
  /** Third-party providers rendered as links (items.<key>). */
  providers?: { key: string; name?: string; nameKey?: string; url?: string }[]
  /** Whether a "closing" paragraph follows the list. */
  closing?: boolean
}

const sections: LegalSection[] = [
  { id: 'controller', paragraphs: 2 },
  {
    id: 'data',
    paragraphs: 1,
    labeledItems: [
      'account',
      'phones',
      'messages',
      'contacts',
      'heartbeats',
      'apiKeys',
      'integrations',
      'logs',
    ],
  },
  { id: 'encryption', paragraphs: 2 },
  { id: 'uses', paragraphs: 1, items: 5, closing: true },
  {
    id: 'providers',
    paragraphs: 1,
    closing: true,
    providers: [
      {
        key: 'firebaseAuth',
        name: 'Google Firebase Authentication',
        url: 'https://firebase.google.com/support/privacy',
      },
      {
        key: 'fcm',
        name: 'Firebase Cloud Messaging (FCM)',
        url: 'https://firebase.google.com/support/privacy',
      },
      {
        key: 'googleCloud',
        name: 'Google Cloud',
        url: 'https://cloud.google.com/terms/cloud-privacy-notice',
      },
      { key: 'smtp', nameKey: 'smtpName' },
      {
        key: 'turnstile',
        name: 'Cloudflare Turnstile',
        url: 'https://www.cloudflare.com/privacypolicy/',
      },
      { key: 'axiom', name: 'Axiom', url: 'https://axiom.co/privacy' },
      { key: 'soketi', name: 'Soketi / WebSockets' },
      { key: 'discord', name: 'Discord', url: 'https://discord.com/privacy' },
      { key: 'androidApp', nameKey: 'androidAppName' },
    ],
  },
  {
    id: 'cookies',
    paragraphs: 1,
    labeledItems: ['language', 'session', 'preferences', 'turnstile'],
    closing: true,
  },
  { id: 'retention', paragraphs: 3 },
  { id: 'security', paragraphs: 2 },
  { id: 'rights', paragraphs: 2 },
  { id: 'children', paragraphs: 1 },
  { id: 'changes', paragraphs: 1 },
  { id: 'contact', paragraphs: 2 },
]

const sectionKey = (id: string) => `legal.privacy.sections.${id}`
</script>

<template>
  <VContainer>
    <VRow>
      <VCol cols="12" md="8" offset-md="2">
        <h1 class="text-display-medium mb-2">
          {{ $t('legal.privacy.title') }}
        </h1>
        <p class="text-body-medium text-medium-emphasis mb-6">
          {{ $t('legal.lastUpdated') }}
        </p>
        <VAlert
          type="info"
          variant="tonal"
          class="mb-6"
          :title="$t('legal.noticeTitle')"
          :text="$t('legal.notice')"
        />
        <p class="text-body-large mb-4">{{ $t('legal.privacy.intro.p1') }}</p>
        <p class="text-body-large mb-4">{{ $t('legal.privacy.intro.p2') }}</p>

        <section
          v-for="(section, index) in sections"
          :id="section.id"
          :key="section.id"
        >
          <h2 class="text-headline-medium mt-8 mb-4">
            {{ index + 1 }}. {{ $t(`${sectionKey(section.id)}.title`) }}
          </h2>
          <i18n-t
            v-for="n in section.paragraphs"
            :key="`p${n}`"
            :keypath="`${sectionKey(section.id)}.p${n}`"
            tag="p"
            scope="global"
            class="text-body-large mb-4"
          >
            <template #settings>
              <NuxtLink to="/settings">{{
                $t('legal.links.settings')
              }}</NuxtLink>
            </template>
            <template #issues>
              <a :href="ISSUES_URL" target="_blank" rel="noopener">{{
                $t('legal.links.issues')
              }}</a>
            </template>
          </i18n-t>

          <ul v-if="section.labeledItems" class="text-body-large pl-6 mb-6">
            <li v-for="item in section.labeledItems" :key="item" class="mb-2">
              <strong>{{
                $t(`${sectionKey(section.id)}.items.${item}.label`)
              }}</strong>
              {{ $t(`${sectionKey(section.id)}.items.${item}.text`) }}
            </li>
          </ul>

          <ul v-if="section.items" class="text-body-large pl-6 mb-6">
            <li v-for="n in section.items" :key="n" class="mb-2">
              {{ $t(`${sectionKey(section.id)}.items.i${n}`) }}
            </li>
          </ul>

          <ul v-if="section.providers" class="text-body-large pl-6 mb-6">
            <li
              v-for="provider in section.providers"
              :key="provider.key"
              class="mb-2"
            >
              <strong>
                <a
                  v-if="provider.url"
                  :href="provider.url"
                  target="_blank"
                  rel="noopener"
                  >{{ provider.name }}</a
                >
                <template v-else>{{
                  provider.nameKey
                    ? $t(`${sectionKey(section.id)}.${provider.nameKey}`)
                    : provider.name
                }}</template
                >:
              </strong>
              {{ $t(`${sectionKey(section.id)}.items.${provider.key}`) }}
            </li>
          </ul>

          <p v-if="section.closing" class="text-body-large mb-4">
            {{ $t(`${sectionKey(section.id)}.closing`) }}
          </p>
        </section>

        <VDivider class="mx-16 mt-10" />
        <div class="text-center mt-8 mb-4">
          <BackButton />
        </div>
      </VCol>
    </VRow>
  </VContainer>
</template>
