<script setup lang="ts">
const { locale, setLocale } = useI18n()

function toggleLanguage() {
  setLocale(locale.value === 'en' ? 'es' : 'en')
}
</script>

<template>
  <v-btn variant="text" @click="toggleLanguage" class="mx-1" size="small">
    {{ locale === 'en' ? 'ES' : 'EN' }}
  </v-btn>
</template>
