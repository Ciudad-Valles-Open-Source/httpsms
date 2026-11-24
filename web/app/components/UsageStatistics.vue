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
    <h4 v-if="!hideTitle" class="text-headline-large mb-3 mt-8">
      {{ $t('pages.billing.overview') }}
    </h4>
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
        /> </v-code
      >.
    </p>
    <VRow v-if="billingStore.billingUsage" class="mt-2">
      <VCol cols="12" sm="6">
        <div class="card-soft stat-card">
          <span class="icon-tile">
            <VIcon :icon="mdiCallMade" />
          </span>
          <div>
            <p class="stat-card-num">
              {{ formatDecimal(billingStore.billingUsage.sent_messages) }}
            </p>
            <p class="stat-card-label">
              {{ $t('pages.billing.messagesSent') }}
            </p>
          </div>
        </div>
      </VCol>
      <VCol cols="12" sm="6">
        <div class="card-soft stat-card">
          <span class="icon-tile icon-tile-secondary">
            <VIcon :icon="mdiCallReceived" />
          </span>
          <div>
            <p class="stat-card-num">
              {{ formatDecimal(billingStore.billingUsage.received_messages) }}
            </p>
            <p class="stat-card-label">
              {{ $t('pages.billing.messagesReceived') }}
            </p>
          </div>
        </div>
      </VCol>
    </VRow>
    <VRow v-else>
      <VCol v-for="n in 2" :key="n" cols="12" sm="6">
        <VSkeletonLoader type="list-item-avatar-two-line" class="card-soft" />
      </VCol>
    </VRow>

    <!-- Usage History -->
    <h4 class="panel-section-title mt-10">
      {{ $t('pages.billing.usageHistory') }}
    </h4>
    <p class="text-medium-emphasis">
      {{ $t('pages.billing.usageHistoryDesc') }}
    </p>
    <VTable class="stack-table" density="comfortable">
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
          v-if="
            billingStore.billingUsage &&
            billingStore.billingUsageHistory.length === 0
          "
        >
          <td colspan="4" class="text-center text-medium-emphasis py-6">
            {{ $t('pages.billing.noHistory') }}
          </td>
        </tr>
        <tr
          v-for="billingUsage in billingStore.billingUsageHistory"
          :key="billingUsage.id"
        >
          <td :data-label="$t('pages.billing.thStartDate')">
            <BillingDateOrdinal :value="billingUsage.start_timestamp" />
          </td>
          <td :data-label="$t('pages.billing.thEndDate')">
            <BillingDateOrdinal :value="billingUsage.end_timestamp" />
          </td>
          <td :data-label="$t('pages.billing.thSent')">
            {{ formatDecimal(billingUsage.sent_messages) }}
          </td>
          <td :data-label="$t('pages.billing.thReceived')">
            {{ formatDecimal(billingUsage.received_messages) }}
          </td>
        </tr>
      </tbody>
    </VTable>
  </div>
</template>

<style scoped>
.stat-card {
  display: flex;
  align-items: center;
  gap: 1rem;
  padding: 1.25rem 1.4rem;
}

.stat-card-num {
  margin: 0;
  font-size: 2rem;
  font-weight: 800;
  line-height: 1.1;
  letter-spacing: -0.02em;
}

.stat-card-label {
  margin: 0.25rem 0 0;
  color: rgb(var(--v-theme-on-surface), var(--v-medium-emphasis-opacity));
  font-size: 0.875rem;
}

.icon-tile-secondary {
  background: rgb(var(--v-theme-secondary), 0.14);
  color: rgb(var(--v-theme-secondary));
}
</style>
