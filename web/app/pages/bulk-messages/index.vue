<script setup lang="ts">
import { mdiArrowLeft, mdiMicrosoftExcel, mdiSendCheck } from '@mdi/js'
import { ErrorMessages } from '~/utils/errors'
import { toApiError } from '~/utils/api-error'
import type { EntitiesBulkMessage } from '~~/shared/types/api'
import capitalize from '~/utils/capitalize'

definePageMeta({
  middleware: ['auth'],
})

useHead({
  title: computed(() => $t('pages.bulkMessages.title')),
})

const router = useRouter()
const { mdAndUp } = useVDisplay()
const authStore = useAuthStore()
const notificationsStore = useNotificationsStore()
const { formatTimestamp } = useFilters()
const { t } = useI18n()
const { useApi } = useApiComposable()

const formFile = ref<File | null>(null)
const loading = ref(true)
const loadingHistory = ref(true)
const errorTitle = ref('')
const errorMessages = ref(new ErrorMessages())
const bulkOrders = ref<EntitiesBulkMessage[]>([])

function parseErrors(error: unknown): ErrorMessages {
  const bag = new ErrorMessages()
  const data = toApiError(error).data?.data
  if (data && typeof data === 'object') {
    Object.keys(data).forEach((key) => bag.addMany(key, data[key]))
  }
  return bag
}

function cleanName(requestId: string): string {
  if (requestId.startsWith('bulk-csv-'))
    return requestId.replace(/^bulk-csv-/, '') + '.csv'
  if (requestId.startsWith('bulk-xls-'))
    return requestId.replace(/^bulk-xls-/, '') + '.xlsx'
  const newFormatMatch = requestId.match(/^bulk-[0-9A-Za-z]+-(.+)$/)
  if (newFormatMatch) return newFormatMatch[1]
  return requestId.replace(/^bulk-/, '')
}

async function fetchBulkOrders() {
  loadingHistory.value = true
  try {
    const api = useApi()
    const response = await api<{ data: EntitiesBulkMessage[] }>(
      '/v1/bulk-messages',
      {
        method: 'GET',
      },
    )
    bulkOrders.value = response.data ?? []
  } catch {
    notificationsStore.addNotification({
      message: t('pages.bulkMessages.errorMessageHistory'),
      type: 'error',
    })
  } finally {
    loadingHistory.value = false
  }
}

async function sendBulkMessages() {
  loading.value = true
  errorMessages.value = new ErrorMessages()
  errorTitle.value = ''

  try {
    const api = useApi()
    const formData = new FormData()
    if (formFile.value) formData.append('document', formFile.value)
    const response = await api<{ message?: string }>('/v1/bulk-messages', {
      method: 'POST',
      body: formData,
    })
    notificationsStore.addNotification({
      message: response?.message ?? t('pages.bulkMessages.successMessageSent'),
      type: 'success',
    })
    loading.value = false
    formFile.value = null
    fetchBulkOrders()
  } catch (error: unknown) {
    errorTitle.value = capitalize(
      toApiError(error).data?.message ?? t('pages.bulkMessages.defaultErrorSending'),
    )
    errorMessages.value = parseErrors(error)
    notificationsStore.addNotification({
      message:
        toApiError(error).data?.message ?? t('pages.bulkMessages.defaultErrorsSending'),
      type: 'error',
    })
    loading.value = false
  }
}

onMounted(async () => {
  await authStore.loadUser()
  loading.value = false
  fetchBulkOrders()
})
</script>

<template>
  <VContainer fluid class="px-0 pt-0" :class="{ 'fill-height': true }">
    <div class="w-100 h-100">
      <VAppBar>
        <VBtn icon to="/threads">
          <VIcon :icon="mdiArrowLeft" />
        </VBtn>
        <VToolbarTitle>
          <div class="py-16">{{ $t('pages.bulkMessages.bulkMessages') }}</div>
        </VToolbarTitle>
        <VProgressLinear
          :active="loading"
          color="primary"
          :indeterminate="loading"
          location="bottom"
          absolute
        />
      </VAppBar>
      <VContainer>
        <VRow>
          <VCol cols="12" md="10" offset-md="1" xxl="8" offset-xxl="2">
            <h5 class="text-headline-large mb-3 mt-3">{{ $t('pages.bulkMessages.bulkMessages') }}</h5>
            <p v-html="$t('pages.bulkMessages.description')">
            </p>
            <VAlert v-if="errorTitle" variant="tonal" type="warning" prominent>
              <h6 class="text-title-large font-weight-bold">
                {{ errorTitle }}
              </h6>
              <ul class="text-body-medium">
                <li
                  v-for="message in errorMessages.get('document')"
                  :key="message"
                >
                  {{ message }}
                </li>
              </ul>
            </VAlert>
            <form>
              <VFileInput
                v-model="formFile"
                :label="$t('pages.bulkMessages.fileLabel')"
                color="primary"
                accept=".csv,application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
                :error-messages="errorMessages.get('document')"
                :append-inner-icon="mdiMicrosoftExcel"
                variant="outlined"
              />
              <div class="d-flex">
                <loading-button
                  color="primary"
                  type="submit"
                  size="large"
                  :loading="loading"
                  :disabled="loading"
                  :icon="mdiSendCheck"
                  @click="sendBulkMessages"
                >
                  {{ $t('pages.bulkMessages.btnSendBulk') }}
                </loading-button>
                <VSpacer />
                <VBtn
                  v-if="mdAndUp"
                  variant="plain"
                  color="info"
                  href="mailto:arnold@httpsms.com?subject=I'm having trouble with the bulk messages"
                >
                  {{ $t('pages.bulkMessages.btnNeedHelp') }}
                </VBtn>
              </div>
            </form>
          </VCol>
        </VRow>
        <VRow class="mt-8">
          <VCol cols="12" md="10" offset-md="1" xxl="8" offset-xxl="2">
            <h4 class="text-headline-large mb-3">{{ $t('pages.bulkMessages.historyTitle') }}</h4>
            <p class="text-medium-emphasis">
              {{ $t('pages.bulkMessages.historyDesc') }}
            </p>
            <VProgressLinear
              v-if="loadingHistory"
              color="primary"
              indeterminate
              class="mb-4"
            />
            <VTable v-else density="comfortable">
              <thead>
                <tr class="text-uppercase text-medium-emphasis">
                  <th class="text-left">{{ $t('pages.bulkMessages.thName') }}</th>
                  <th class="text-center">{{ $t('pages.bulkMessages.thCreatedAt') }}</th>
                  <th class="text-center">{{ $t('pages.bulkMessages.thTotal') }}</th>
                  <th class="text-center">{{ $t('pages.bulkMessages.thPending') }}</th>
                  <th class="text-center">{{ $t('pages.bulkMessages.thScheduled') }}</th>
                  <th class="text-center">{{ $t('pages.bulkMessages.thSent') }}</th>
                  <th class="text-center">{{ $t('pages.bulkMessages.thDelivered') }}</th>
                  <th class="text-center">{{ $t('pages.bulkMessages.thFailed') }}</th>
                  <th class="text-center">{{ $t('pages.bulkMessages.thExpired') }}</th>
                </tr>
              </thead>
              <tbody>
                <tr
                  v-for="order in bulkOrders"
                  :key="order.request_id"
                  class="clickable-row"
                  @click="
                    router.push(`/search-messages?query=${order.request_id}`)
                  "
                >
                  <td class="text-left">{{ cleanName(order.request_id) }}</td>
                  <td class="text-center">
                    {{ formatTimestamp(order.created_at) }}
                  </td>
                  <td class="text-center">{{ order.total }}</td>
                  <td class="text-center">{{ order.pending_count }}</td>
                  <td class="text-center">{{ order.scheduled_count }}</td>
                  <td class="text-center">{{ order.sent_count }}</td>
                  <td class="text-center">{{ order.delivered_count }}</td>
                  <td class="text-center">{{ order.failed_count }}</td>
                  <td class="text-center">{{ order.expired_count }}</td>
                </tr>
              </tbody>
            </VTable>
          </VCol>
        </VRow>
      </VContainer>
    </div>
  </VContainer>
</template>

<style scoped>
.clickable-row {
  cursor: pointer;
}

.clickable-row:hover {
  background-color: rgb(0 0 0 / 4%);
}
</style>
