<script setup lang="ts">
import { mdiArrowLeft } from '@mdi/js'

definePageMeta({
  middleware: ['auth'],
})

useHead({
  title: computed(() => $t('pages.statistics.title')),
})

const { lgAndUp } = useVDisplay()
const billingStore = useBillingStore()

const loading = ref(true)

async function loadData() {
  await Promise.all([
    billingStore.loadBillingUsage(),
    billingStore.loadBillingUsageHistory(),
  ])
  loading.value = false
}

onMounted(() => {
  loadData()
})
</script>

<template>
  <VContainer fluid class="px-0 pt-0" :class="{ 'fill-height': lgAndUp }">
    <div class="w-100 h-100">
      <VAppBar>
        <VBtn icon to="/threads">
          <VIcon :icon="mdiArrowLeft" />
        </VBtn>
        <VToolbarTitle>{{ $t('pages.statistics.title') }}</VToolbarTitle>
        <VProgressLinear
          color="primary"
          :active="loading"
          :indeterminate="loading"
          absolute
          location="bottom"
        />
      </VAppBar>
      <VContainer>
        <VRow>
          <VCol cols="12" md="9" offset-md="1" xl="8" offset-xl="2">
            <h4 class="text-headline-large mb-3 mt-0">{{ $t('pages.statistics.title') }}</h4>
            <UsageStatistics :hide-title="true" />
          </VCol>
        </VRow>
      </VContainer>
    </div>
  </VContainer>
</template>
