<script setup lang="ts">
import { mdiMessageTextOutline } from '@mdi/js'

definePageMeta({
  middleware: ['auth'],
})

useHead({
  title: 'Threads - httpSMS',
})

const { lgAndUp } = useVDisplay()
const authStore = useAuthStore()
const phonesStore = usePhonesStore()
const threadsStore = useThreadsStore()

onMounted(async () => {
  await authStore.loadUser()
  await phonesStore.loadPhones()
  await threadsStore.loadThreads()
})
</script>

<template>
  <div v-if="lgAndUp" class="threads-placeholder">
    <div class="empty-state">
      <span class="icon-tile icon-tile-lg icon-tile-gradient">
        <v-icon size="34" :icon="mdiMessageTextOutline" />
      </span>
      <h2 class="empty-state-title text-headline-small">
        {{ $t('threads.select_message') }}
      </h2>
      <p class="empty-state-text">
        {{ $t('threads.select_message_desc') }}
      </p>
    </div>
  </div>
  <div v-else class="threads-mobile">
    <MessageThreadHeader />
    <MessageThread />
  </div>
</template>

<style scoped>
.threads-placeholder {
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 100dvh;
  background:
    radial-gradient(
      circle at 50% 40%,
      rgb(var(--v-theme-primary), 0.1),
      transparent 45%
    ),
    rgb(var(--v-theme-background));
}

.threads-mobile {
  min-height: 100dvh;
  background: rgb(var(--v-theme-surface));
}
</style>
