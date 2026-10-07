<script setup lang="ts">
import { mdiCallMade, mdiCallReceived } from '@mdi/js'
import { useBillingStore } from '~/stores/billing'

const billingStore = useBillingStore()
const { lgAndUp } = useVDisplay()
const { formatDecimal } = useFilters()

defineProps<{
  hideTitle?: boolean
}>()
</script>

<template>
  <div>
    <!-- Overview -->
    <h4 v-if="!hideTitle" class="text-headline-large mb-3 mt-8">{{ $t('pages.billing.overview') }}</h4>
    <p class="text-medium-emphasis">
      {{ $t('pages.billing.overviewDesc') }}
      <v-code v-if="billingStore.billingUsage" class="font-weight-bold">
        <BillingDateOrdinal
          :value="billingStore.billingUsage.start_timestamp"
        />
      </v-code>
      {{ $t('pages.billing.overviewDescTo') }}
      <v-code v-if="billingStore.billingUsage" class="font-weight-bold">
        <BillingDateOrdinal
          :value="billingStore.billingUsage.end_timestamp"
        />
      </v-code>.
    </p>
    <VRow v-if="billingStore.billingUsage">
      <VCol cols="12" md="6">
        <VAlert
          type="info"
          variant="tonal"
          :icon="mdiCallMade"
          prominent
        >
          <h2 class="text-headline-large my-0">
            {{ formatDecimal(billingStore.billingUsage.sent_messages) }}
          </h2>
          <p class="text-medium-emphasis mt-n1">{{ $t('pages.billing.messagesSent') }}</p>
        </VAlert>
      </VCol>
      <VCol cols="12" md="6">
        <VAlert
          type="warning"
          variant="tonal"
          :icon="mdiCallReceived"
          prominent
        >
          <h2 class="text-headline-large font-weight-bold my-0">
            {{ formatDecimal(billingStore.billingUsage.received_messages) }}
          </h2>
          <p class="text-medium-emphasis mt-n1">{{ $t('pages.billing.messagesReceived') }}</p>
        </VAlert>
      </VCol>
    </VRow>

    <!-- Usage History -->
    <h4 class="text-headline-large mb-3 mt-8">{{ $t('pages.billing.usageHistory') }}</h4>
    <p class="text-medium-emphasis">
      {{ $t('pages.billing.usageHistoryDesc') }}
    </p>
    <VTable density="comfortable">
      <thead>
        <tr class="text-uppercase text-medium-emphasis">
          <th class="text-left">{{ $t('pages.billing.thStartDate') }}</th>
          <th class="text-left">{{ $t('pages.billing.thEndDate') }}</th>
          <th class="text-left">
            {{ $t('pages.billing.thSent') }}
            <span v-if="lgAndUp">{{ $t('pages.billing.thMessages') }}</span>
          </th>
          <th class="text-left">
            {{ $t('pages.billing.thReceived') }}
            <span v-if="lgAndUp">{{ $t('pages.billing.thMessages') }}</span>
          </th>
        </tr>
      </thead>
      <tbody>
        <tr
          v-for="billingUsage in billingStore.billingUsageHistory"
          :key="billingUsage.id"
        >
          <td>
            <BillingDateOrdinal :value="billingUsage.start_timestamp" />
          </td>
          <td>
            <BillingDateOrdinal :value="billingUsage.end_timestamp" />
          </td>
          <td>{{ formatDecimal(billingUsage.sent_messages) }}</td>
          <td>{{ formatDecimal(billingUsage.received_messages) }}</td>
        </tr>
      </tbody>
    </VTable>
  </div>
</template>
