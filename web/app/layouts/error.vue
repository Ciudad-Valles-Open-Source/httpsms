<script setup lang="ts">
import { mdiArrowLeft } from '@mdi/js'

useThemeScope('public')

defineProps<{
  error: {
    statusCode: number
    message: string
  }
}>()
</script>

<template>
  <v-app>
    <v-main>
      <div class="error-page">
        <AppLogo class="mb-10" />
        <p class="error-page-code text-gradient">{{ error.statusCode }}</p>
        <p class="text-title-large text-medium-emphasis mt-2 mb-8">
          {{ error.message }}
        </p>
        <v-btn class="btn-gradient" size="large" to="/">
          <v-icon start :icon="mdiArrowLeft" />
          {{ $t('error.goHome') }}
        </v-btn>
      </div>
    </v-main>
  </v-app>
</template>

<style scoped>
.error-page {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  min-height: 100dvh;
  padding: 2rem 1rem;
  text-align: center;
}

.error-page-code {
  margin: 0;
  font-size: clamp(4.5rem, 14vw, 8rem);
  font-weight: 800;
  line-height: 1;
  letter-spacing: -0.04em;
}
</style>
