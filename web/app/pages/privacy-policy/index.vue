<script setup lang="ts">
definePageMeta({
  layout: 'website',
})

useSeoMeta({
  title: computed(() => $t('pages.privacy.title')),
  description: computed(() => $t('pages.privacy.descriptionMeta')),
})
</script>

<template>
  <VContainer>
    <VRow>
      <VCol cols="12" md="8" offset-md="2">
        <h1 class="text-display-medium mb-2">{{ $t('pages.privacy.titlePrivacy') }}</h1>
        <p class="text-body-large mb-4">
          {{ $t('pages.privacy.p1') }}
        </p>

        <h2 class="text-headline-medium mt-8 mb-4">
          {{ $t('pages.privacy.titleInfo') }}
        </h2>
        <p class="text-body-large mb-4">
          {{ $t('pages.privacy.p2') }}
        </p>
        <p class="text-body-large mb-4">
          {{ $t('pages.privacy.p3') }}
        </p>
        <ul class="text-body-large pl-6 mb-6">
          <li>
            <a
              href="https://www.google.com/policies/privacy/"
              target="_blank"
              rel="noopener"
              >Google Play Services</a
            >
          </li>
          <li>
            <a
              href="https://firebase.google.com/policies/analytics"
              target="_blank"
              rel="noopener"
              >Google Analytics for Firebase</a
            >
          </li>
          <li>
            <a href="https://sentry.io/privacy/" target="_blank" rel="noopener"
              >Sentry</a
            >
          </li>
          <li>
            <a
              href="https://segment.com/legal/privacy/"
              target="_blank"
              rel="noopener"
              >Segment</a
            >
          </li>
          <li>
            <a
              href="https://privacy.microsoft.com/en-GB/privacystatement"
              target="_blank"
              rel="noopener"
              >Microsoft Clarity</a
            >
          </li>
        </ul>

        <h2 class="text-headline-medium mt-8 mb-4">{{ $t('pages.privacy.titleLogData') }}</h2>
        <p class="text-body-large mb-4">
          {{ $t('pages.privacy.pLogData') }}
        </p>

        <h2 class="text-headline-medium mt-8 mb-4">{{ $t('pages.privacy.titleCookies') }}</h2>
        <p class="text-body-large mb-4">
          {{ $t('pages.privacy.pCookies') }}
        </p>

        <h2 class="text-headline-medium mt-8 mb-4">{{ $t('pages.privacy.titleProviders') }}</h2>
        <p class="text-body-large mb-4">
          {{ $t('pages.privacy.pProviders1') }}
        </p>
        <ul class="text-body-large pl-6 mb-6">
          <li>{{ $t('pages.privacy.liProv1') }}</li>
          <li>{{ $t('pages.privacy.liProv2') }}</li>
          <li>{{ $t('pages.privacy.liProv3') }}</li>
          <li>{{ $t('pages.privacy.liProv4') }}</li>
        </ul>
        <p class="text-body-large mb-4">
          {{ $t('pages.privacy.pProviders2') }}
        </p>

        <h2 class="text-headline-medium mt-8 mb-4">{{ $t('pages.privacy.titleSecurity') }}</h2>
        <p class="text-body-large mb-4">
          {{ $t('pages.privacy.pSecurity') }}
        </p>

        <h2 class="text-headline-medium mt-8 mb-4">{{ $t('pages.privacy.titleLinks') }}</h2>
        <p class="text-body-large mb-4">
          {{ $t('pages.privacy.pLinks') }}
        </p>

        <h2 class="text-headline-medium mt-8 mb-4">{{ $t('pages.privacy.titleChildren') }}</h2>
        <p class="text-body-large mb-4">
          {{ $t('pages.privacy.pChildren') }}
        </p>

        <h2 class="text-headline-medium mt-8 mb-4">
          {{ $t('pages.privacy.titleChanges') }}
        </h2>
        <p class="text-body-large mb-4">
          {{ $t('pages.privacy.pChanges1') }}
        </p>
        <p class="text-body-large mb-4">
          {{ $t('pages.privacy.pChanges2') }}
        </p>

        <h2 class="text-headline-medium mt-8 mb-4">{{ $t('pages.privacy.titleContact') }}</h2>
        <p class="text-body-large mb-10" v-html="$t('pages.privacy.pContact')">
        </p>
        <VDivider class="mx-16" />
        <div class="text-center mt-8 mb-4">
          <BackButton />
        </div>
      </VCol>
    </VRow>
  </VContainer>
</template>
