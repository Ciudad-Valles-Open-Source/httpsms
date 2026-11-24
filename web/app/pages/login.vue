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
  <div class="login-page">
    <div class="login-page-top">
      <AppLogo />
      <div class="d-flex align-center">
        <ThemeToggle />
        <LanguageSwitcher />
      </div>
    </div>
    <main class="login-page-main">
      <div class="login-card card-soft">
        <span class="hero-badge mb-5">
          <span class="dot" />
          {{ $t('pages.login.badge') }}
        </span>
        <h1 class="login-card-title">{{ $t('pages.login.welcome') }}</h1>
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

<style scoped>
.login-page {
  display: flex;
  flex-direction: column;
  min-height: 100dvh;
  background:
    radial-gradient(
      circle at 85% 0%,
      rgb(var(--v-theme-primary), 0.22),
      transparent 40%
    ),
    radial-gradient(
      circle at 0% 100%,
      rgb(var(--v-theme-secondary), 0.16),
      transparent 40%
    ),
    rgb(var(--v-theme-background));
}

.login-page-top {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 1rem 1.25rem;
}

.login-page-main {
  display: flex;
  flex: 1;
  flex-direction: column;
  justify-content: center;
  width: 100%;
  max-width: 460px;
  margin: 0 auto;
  padding: 1rem 1rem 4rem;
}

.login-card {
  padding: 2.25rem 2rem;
  text-align: center;
}

.login-card-title {
  margin: 0;
  font-size: 2rem;
  font-weight: 800;
  letter-spacing: -0.03em;
}

@media (width <= 600px) {
  .login-card {
    padding: 1.75rem 1.25rem;
  }
}
</style>
