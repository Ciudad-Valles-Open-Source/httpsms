<script setup lang="ts">
definePageMeta({
  layout: 'website',
})

const { t } = useI18n()

useSeoMeta({
  title: computed(() => t('terms.meta_title')),
  description: computed(() => t('terms.meta_description')),
})
</script>

<template>
  <VContainer>
    <VRow>
      <VCol cols="12" md="8" offset-md="2">
        <h1 class="text-display-medium mb-2">{{ $t('terms.title') }}</h1>
        <p class="text-body-large mb-4">{{ $t('terms.p1') }}</p>
        <p class="text-body-large mb-4">{{ $t('terms.p2') }}</p>
        <p class="text-body-large mb-4">{{ $t('terms.p3') }}</p>
        <p class="text-body-large mb-4">{{ $t('terms.p4') }}</p>
        <ul class="text-body-large pl-6 mb-6">
          <li>
            <a
              href="https://policies.google.com/terms"
              target="_blank"
              rel="noopener"
              >{{ $t('terms.google_play') }}</a
            >
          </li>
          <li>
            <a
              href="https://firebase.google.com/terms/analytics"
              target="_blank"
              rel="noopener"
              >{{ $t('terms.google_analytics') }}</a
            >
          </li>
          <li>
            <a href="https://sentry.io/terms/" target="_blank" rel="noopener"
              >{{ $t('terms.sentry') }}</a
            >
          </li>
          <li>
            <a
              href="https://segment.com/legal/terms/"
              target="_blank"
              rel="noopener"
              >{{ $t('terms.segment') }}</a
            >
          </li>
        </ul>
        <p class="text-body-large mb-4">{{ $t('terms.p5') }}</p>
        <p class="text-body-large mb-4">{{ $t('terms.p6') }}</p>
        <p class="text-body-large mb-4">{{ $t('terms.p7') }}</p>
        <p class="text-body-large mb-4">{{ $t('terms.p8') }}</p>
        <p class="text-body-large mb-4">{{ $t('terms.p9') }}</p>

        <h2 class="text-headline-medium mt-8 mb-4">{{ $t('terms.changes_title') }}</h2>
        <p class="text-body-large mb-4">{{ $t('terms.changes_desc') }}</p>
        <p class="text-body-large mb-4">{{ $t('terms.effective_date') }}</p>

        <h2 class="text-headline-medium mt-8 mb-4">{{ $t('terms.contact_us') }}</h2>
        <p class="text-body-large mb-10" v-html="$t('terms.contact_desc')"></p>
        <VDivider class="mx-16" />
        <div class="text-center mt-8 mb-4">
          <BackButton />
        </div>
      </VCol>
    </VRow>
  </VContainer>
</template>
