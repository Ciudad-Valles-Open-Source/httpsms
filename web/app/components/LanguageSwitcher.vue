<script setup lang="ts">
import { mdiTranslate } from '@mdi/js'

const { locale, setLocale } = useI18n()

function toggleLanguage() {
  setLocale(locale.value === 'en' ? 'es' : 'en')
}
</script>

<template>
  <v-btn
    variant="text"
    size="small"
    class="language-switcher"
    :aria-label="$t('components.languageSwitcher.label')"
    @click="toggleLanguage"
  >
    <v-icon :icon="mdiTranslate" size="16" class="mr-1" />
    {{ locale === 'en' ? 'ES' : 'EN' }}
  </v-btn>
</template>

<style scoped>
.language-switcher {
  min-width: 0;
  padding-inline: 0.6rem;
}
</style>
