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
  <VContainer fluid class="fill-height d-flex align-center justify-center">
    <VRow align="center" justify="center">
      <VCol cols="12" md="4" xl="3" class="mt-n16">
        <div class="text-center mb-5">
          <VAvatar color="#121212" size="45" rounded="0" class="mt-n8 mr-4">
            <VImg :src="'/img/logo.svg'" />
          </VAvatar>
          <span class="text-display-small">{{ $t('pages.login.welcome') }}</span>
        </div>
        <p class="text-body-large text-center text-medium-emphasis mt-1 mb-4" v-html="$t('pages.login.joinUsers')">
        </p>
        <VCard max-width="500" class="mx-auto pa-4">
          <VCardText class="px-4 pt-2 pb-0">
            <ClientOnly>
              <FirebaseAuth :to="to" />
            </ClientOnly>
          </VCardText>
        </VCard>
        <div class="text-center mt-4">
          <BackButton />
        </div>
      </VCol>
    </VRow>
  </VContainer>
</template>
