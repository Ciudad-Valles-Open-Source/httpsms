<script setup lang="ts">
definePageMeta({
  middleware: ['guest'],
  layout: 'blank',
})

const { t } = useI18n()

useSeoMeta({
  title: () => t('pages.login.metaTitle'),
  description: () => t('pages.login.metaDescription'),
})

const route = useRoute()
const to = computed(() => (route.query.to as string) || '/threads')
</script>

<template>
  <div class="auth-page">
    <div class="auth-page-top">
      <AppLogo />
      <div class="d-flex align-center">
        <ThemeToggle />
        <LanguageSwitcher />
      </div>
    </div>
    <main class="auth-page-main">
      <div class="auth-card card-soft">
        <span class="hero-badge mb-5">
          <span class="dot" />
          {{ $t('pages.login.badge') }}
        </span>
        <h1 class="auth-card-title">{{ $t('pages.login.welcome') }}</h1>
        <p
          class="text-body-medium text-medium-emphasis mt-2 mb-6"
          v-html="$t('pages.login.joinUsers')"
        ></p>
        <div class="text-left">
          <ClientOnly>
            <FirebaseAuth :to="to" />
            <template #fallback>
              <VSkeletonLoader
                type="button, button, button"
                color="transparent"
              />
            </template>
          </ClientOnly>
        </div>
      </div>
      <div class="text-center mt-6">
        <BackButton />
      </div>
    </main>
  </div>
</template>
