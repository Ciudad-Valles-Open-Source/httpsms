<script setup lang="ts">
definePageMeta({
  layout: 'website',
})

const { t } = useI18n()

useSeoMeta({
  title: computed(() => t('legal.terms.metaTitle')),
  description: computed(() => t('legal.terms.metaDescription')),
  ogTitle: computed(() => t('legal.terms.metaTitle')),
  ogDescription: computed(() => t('legal.terms.metaDescription')),
})

const REPO_URL = 'https://github.com/Ciudad-Valles-Open-Source/httpsms'
const ISSUES_URL = `${REPO_URL}/issues`
const LICENSE_URL = `${REPO_URL}/blob/main/LICENSE`

interface LegalSection {
  id: string
  /** Number of paragraphs (p1..pN) shown before the list. */
  paragraphs: number
  /** Number of list items (items.i1..iN). */
  items?: number
}

const sections: LegalSection[] = [
  { id: 'about', paragraphs: 2 },
  { id: 'license', paragraphs: 2 },
  { id: 'operators', paragraphs: 2 },
  { id: 'responsibilities', paragraphs: 1, items: 8 },
  { id: 'prohibited', paragraphs: 1, items: 8 },
  { id: 'suspension', paragraphs: 2 },
  { id: 'warranty', paragraphs: 2 },
  { id: 'liability', paragraphs: 3 },
  { id: 'indemnity', paragraphs: 1 },
  { id: 'thirdParty', paragraphs: 1 },
  { id: 'changes', paragraphs: 1 },
  { id: 'contact', paragraphs: 2 },
]

const sectionKey = (id: string) => `legal.terms.sections.${id}`
</script>

<template>
  <VContainer>
    <VRow>
      <VCol cols="12" md="8" offset-md="2">
        <h1 class="text-display-medium mb-2">
          {{ $t('legal.terms.title') }}
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
        <p class="text-body-large mb-4">{{ $t('legal.terms.intro.p1') }}</p>
        <p class="text-body-large mb-4">{{ $t('legal.terms.intro.p2') }}</p>

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
            <template #privacy>
              <NuxtLink to="/privacy-policy">{{
                $t('legal.links.privacy')
              }}</NuxtLink>
            </template>
            <template #repo>
              <a :href="REPO_URL" target="_blank" rel="noopener">{{
                $t('legal.links.repo')
              }}</a>
            </template>
            <template #license>
              <a :href="LICENSE_URL" target="_blank" rel="noopener">{{
                $t('legal.links.license')
              }}</a>
            </template>
            <template #issues>
              <a :href="ISSUES_URL" target="_blank" rel="noopener">{{
                $t('legal.links.issues')
              }}</a>
            </template>
          </i18n-t>

          <ul v-if="section.items" class="text-body-large pl-6 mb-6">
            <li v-for="n in section.items" :key="n" class="mb-2">
              {{ $t(`${sectionKey(section.id)}.items.i${n}`) }}
            </li>
          </ul>
        </section>

        <VDivider class="mx-16 mt-10" />
        <div class="text-center mt-8 mb-4">
          <BackButton />
        </div>
      </VCol>
    </VRow>
  </VContainer>
</template>
