<script setup lang="ts">
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
const { t } = useI18n()

onMounted(async () => {
  await authStore.loadUser()
  await phonesStore.loadPhones()
  await threadsStore.loadThreads()
})
</script>

<template>
  <VContainer class="pa-0" fluid :class="{ 'fill-height': lgAndUp }">
    <VRow
      v-if="lgAndUp"
      align="center"
      :class="{ 'fill-height': lgAndUp }"
      justify="center"
    >
      <div :class="{ 'mt-n16': lgAndUp }">
        <VImg
          class="mx-auto mb-4"
          :class="{ 'mt-n16': lgAndUp }"
          max-height="400"
          max-width="90%"
          :src="'/img/person-texting.svg'"
        />
        <div class="text-center">
          <h3 class="text-headline-medium mt-4 mb-0">{{ $t('threads.select_message') }}</h3>
          <p class="text-medium-emphasis mt-0" v-html="$t('threads.discord_help_desc')"></p>
        </div>
      </div>
    </VRow>
    <VRow v-else justify="end">
      <VCol class="px-0 py-0">
        <MessageThreadHeader />
        <MessageThread />
      </VCol>
    </VRow>
  </VContainer>
</template>
