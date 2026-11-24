<script setup lang="ts">
withDefaults(
  defineProps<{
    to?: string
    size?: number
    showName?: boolean
  }>(),
  { to: '/', size: 36, showName: true },
)
</script>

<template>
  <NuxtLink
    :to="to"
    class="app-logo"
    :aria-label="$t('components.appLogo.home')"
  >
    <span
      class="app-logo-mark"
      :style="{ width: `${size}px`, height: `${size}px` }"
    >
      <img src="/img/logo.svg" alt="" />
    </span>
    <span v-if="showName" class="app-logo-name">
      http<span class="text-gradient">SMS</span>
    </span>
  </NuxtLink>
</template>

<style scoped>
.app-logo {
  display: inline-flex;
  align-items: center;
  gap: 0.6rem;
  color: rgb(var(--v-theme-on-background)) !important;
  text-decoration: none;
}

.app-logo-mark {
  display: inline-flex;
  flex: 0 0 auto;
  align-items: center;
  justify-content: center;
  padding: 0.3rem;
  border: 1px solid rgb(var(--v-border-color), var(--v-border-opacity));
  border-radius: 0.7rem;
  background: rgb(var(--v-theme-surface));
}

.app-logo-mark img {
  width: 100%;
  height: 100%;
  object-fit: contain;
}

.app-logo-name {
  font-size: 1.2rem;
  font-weight: 800;
  letter-spacing: -0.02em;
}
</style>
