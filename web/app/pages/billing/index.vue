<script setup lang="ts">
import {
  mdiAlert,
  mdiArrowLeft,
  mdiCallMade,
  mdiCallReceived,
  mdiCheck,
  mdiDownloadOutline,
  mdiInvoice,
} from '@mdi/js'
import type {
  RequestsUserPaymentInvoice,
  ResponsesUserSubscriptionPaymentsResponse,
} from '~~/shared/types/api'

import { countries, getStateOptions } from '~/utils/countries'
import { getApiErrorMessage } from '~/utils/api-error'

type SubscriptionPayment =
  ResponsesUserSubscriptionPaymentsResponse['data'][number]

definePageMeta({
  middleware: ['auth'],
})

useHead({
  title: computed(() => $t('pages.billing.title')),
})

const config = useRuntimeConfig()
const { lgAndUp } = useVDisplay()
const authStore = useAuthStore()
const billingStore = useBillingStore()
const notificationsStore = useNotificationsStore()
const { formatDecimal, formatTimestamp } = useFilters()

const loading = ref(true)
const loadingSubscriptionPayments = ref(false)
const dialog = ref(false)
const subscriptionInvoiceDialog = ref(false)
const payments = ref<ResponsesUserSubscriptionPaymentsResponse | null>(null)
const selectedPayment = ref<SubscriptionPayment | null>(null)
const invoiceFormName = ref('')
const invoiceFormAddress = ref('')
const invoiceFormCity = ref('')
const invoiceFormState = ref('')
const invoiceFormZipCode = ref('')
const invoiceFormCountry = ref('')
const invoiceFormNotes = ref('')
const errorMessages = ref(new Map<string, string>())

type PaymentPlan = {
  name: string
  id: string
  price: number
  messagesPerMonth: number
}

const plans: PaymentPlan[] = [
  { name: 'Free', id: 'free', messagesPerMonth: 200, price: 0 },
  {
    name: 'PRO - Monthly',
    id: 'pro-monthly',
    messagesPerMonth: 5000,
    price: 10,
  },
  {
    name: 'PRO - Yearly',
    id: 'pro-yearly',
    messagesPerMonth: 5000,
    price: 100,
  },
  {
    name: 'Ultra - Monthly',
    id: 'ultra-monthly',
    messagesPerMonth: 10000,
    price: 20,
  },
  {
    name: 'Ultra - Yearly',
    id: 'ultra-yearly',
    messagesPerMonth: 10000,
    price: 200,
  },
  {
    name: '20k - Monthly',
    id: '20k-monthly',
    messagesPerMonth: 20000,
    price: 35,
  },
  {
    name: '20k - Yearly',
    id: '20k-yearly',
    messagesPerMonth: 20000,
    price: 350,
  },
  {
    name: '50k - Monthly',
    id: '50k-monthly',
    messagesPerMonth: 50000,
    price: 89,
  },
  {
    name: '100k - Monthly',
    id: '100k-monthly',
    messagesPerMonth: 100000,
    price: 175,
  },
  {
    name: '200k - Monthly',
    id: '200k-monthly',
    messagesPerMonth: 200000,
    price: 350,
  },
  {
    name: 'PRO - Lifetime',
    id: 'pro-lifetime',
    messagesPerMonth: 10000,
    price: 1000,
  },
]

const plan = computed<PaymentPlan>(() => {
  return (plans.find(
    (x) => x.id === (authStore.user?.subscription_name || 'free'),
  ) ?? plans[0])!
})

const isOnFreePlan = computed(() => plan.value.id === 'free')
const isOnLifetimePlan = computed(() => plan.value.id === 'pro-lifetime')
const subscriptionIsCancelled = computed(
  () => authStore.user?.subscription_status === 'cancelled',
)

const invoiceStateOptions = computed(() =>
  getStateOptions(invoiceFormCountry.value),
)

const totalMessages = computed(() => {
  if (!billingStore.billingUsage) return 0
  return (
    billingStore.billingUsage.sent_messages +
    billingStore.billingUsage.received_messages
  )
})

const checkoutURL = computed(() => {
  const url = new URL(config.public.checkoutUrl as string)
  const user = authStore.authUser
  if (user) {
    url.searchParams.append('checkout[custom][user_id]', user.id)
    if (user.email) {
      url.searchParams.append('checkout[email]', user.email)
    }
    if (user.displayName) {
      url.searchParams.append('checkout[name]', user.displayName)
    }
  }
  return url.toString()
})

const enterpriseCheckoutURL = computed(() => {
  const url = new URL(config.public.enterpriseCheckoutUrl as string)
  const user = authStore.authUser
  if (user) {
    url.searchParams.append('checkout[custom][user_id]', user.id)
    if (user.email) {
      url.searchParams.append('checkout[email]', user.email)
    }
    if (user.displayName) {
      url.searchParams.append('checkout[name]', user.displayName)
    }
  }
  return url.toString()
})

async function loadData() {
  await Promise.all([
    authStore.loadUser(),
    billingStore.loadBillingUsage(),
    billingStore.loadBillingUsageHistory(),
  ])
  loading.value = false
  loadSubscriptionInvoices()
}

async function loadSubscriptionInvoices() {
  if (!authStore.user?.subscription_id) return
  loadingSubscriptionPayments.value = true
  try {
    payments.value = await billingStore.indexSubscriptionPayments()
  } finally {
    loadingSubscriptionPayments.value = false
  }
}

async function updateDetails() {
  loading.value = true
  try {
    window.location.href = await billingStore.getSubscriptionUpdateLink()
  } catch (error: unknown) {
    notificationsStore.addNotification({
      message: getApiErrorMessage(
        error,
        'We could not redirect you to the subscription update portal.',
      ),
      type: 'error',
    })
    loading.value = false
  }
}

async function cancelPlan() {
  loading.value = true
  try {
    await billingStore.cancelSubscription()
    notificationsStore.addNotification({
      message: 'Subscription cancelled successfully',
      type: 'success',
    })
    navigateTo('/')
  } catch (error: unknown) {
    notificationsStore.addNotification({
      message: getApiErrorMessage(
        error,
        'We could not cancel your subscription.',
      ),
      type: 'error',
    })
    loading.value = false
  }
}

async function generateInvoice() {
  errorMessages.value = new Map()
  loading.value = true
  try {
    await billingStore.generateSubscriptionPaymentInvoice(
      selectedPayment.value?.id || '',
      {
        name: invoiceFormName.value,
        address: invoiceFormAddress.value,
        city: invoiceFormCity.value,
        state: invoiceFormState.value,
        zip_code: invoiceFormZipCode.value,
        country: invoiceFormCountry.value,
        notes: invoiceFormNotes.value,
      } as RequestsUserPaymentInvoice,
    )
    subscriptionInvoiceDialog.value = false
  } catch (error: unknown) {
    if (error instanceof Map) {
      errorMessages.value = error
    }
  } finally {
    loading.value = false
  }
}

function showInvoiceDialog(payment: SubscriptionPayment) {
  selectedPayment.value = payment
  subscriptionInvoiceDialog.value = true
}

onMounted(async () => {
  await loadData()
})
</script>

<template>
  <VContainer fluid class="px-0 pt-0" :class="{ 'fill-height': lgAndUp }">
    <div class="w-100 h-100">
      <VAppBar>
        <VBtn icon to="/threads">
          <VIcon :icon="mdiArrowLeft" />
        </VBtn>
        <VToolbarTitle>{{ $t('pages.billing.accountUsage') }}</VToolbarTitle>
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
            <!-- Current Plan -->
            <template v-if="config.public.enableBilling">
            <h4 class="text-headline-large mb-3 mt-0">{{ $t('pages.billing.currentPlan') }}</h4>
            <VRow v-if="authStore.user">
              <VCol md="6">
                <VAlert type="info" :icon="false" variant="tonal" prominent>
                  <div>
                    <h1
                      class="text-title-large mt-0 mb-0 font-weight-bold text-uppercase"
                    >
                      <span v-if="isOnFreePlan">{{ plan.name }}</span>
                      <span v-else-if="subscriptionIsCancelled" v-html="$t('pages.billing.downgradedFree').replace('{date}', new Date(authStore.user.subscription_ends_at!).toLocaleDateString())">
                      </span>
                      <span v-else>{{ plan.name }}</span>
                    </h1>
                    <p
                      v-if="
                        !isOnFreePlan &&
                        !isOnLifetimePlan &&
                        !subscriptionIsCancelled
                      "
                      class="text-medium-emphasis mt-1"
                      v-html="$t('pages.billing.nextBill').replace('{price}', String(plan.price)).replace('{date}', new Date(authStore.user.subscription_renews_at!).toLocaleDateString())"
                    >
                    </p>
                    <p v-if="isOnLifetimePlan" class="text-medium-emphasis" v-html="$t('pages.billing.lifetimePlanDesc').replace('{price}', String(plan.price))">
                    </p>
                    <p
                      v-else-if="subscriptionIsCancelled"
                      class="text-medium-emphasis"
                      v-html="$t('pages.billing.downgradedFree').replace('{date}', new Date(authStore.user.subscription_ends_at!).toLocaleDateString())"
                    >
                    </p>
                    <p v-if="!config.public.enableBilling" class="text-medium-emphasis mt-1">
                      {{ formatDecimal(totalMessages) }} {{ $t('pages.billing.messagesSent') }}
                    </p>
                    <p v-else class="text-medium-emphasis mt-1">
                      {{ $t('pages.billing.messagesRatio').replace('{used}', formatDecimal(totalMessages)).replace('{total}', formatDecimal(plan.messagesPerMonth)) }}
                    </p>
                  </div>
                  <div class="d-flex mb-1 mt-1">
                    <VBtn
                      v-if="
                        !subscriptionIsCancelled &&
                        !isOnFreePlan &&
                        !isOnLifetimePlan
                      "
                      color="primary"
                      :loading="loading"
                      @click="updateDetails"
                    >
                      {{ $t('pages.billing.updatePlan') }}
                    </VBtn>
                    <VBtn
                      v-else-if="!isOnLifetimePlan"
                      color="primary"
                      :href="checkoutURL"
                    >
                      {{ $t('pages.billing.upgradePlan') }}
                    </VBtn>
                    <VSpacer />
                    <VDialog
                      v-if="
                        !subscriptionIsCancelled &&
                        !isOnFreePlan &&
                        !isOnLifetimePlan
                      "
                      v-model="dialog"
                      max-width="590"
                      opacity="0.9"
                    >
                      <template #activator="{ props: activatorProps }">
                        <VBtn
                          v-bind="activatorProps"
                          color="error"
                          variant="text"
                        >
                          {{ $t('pages.billing.cancelPlan') }}
                        </VBtn>
                      </template>
                      <VCard>
                        <VCardText class="pt-4">
                          <h2 class="text-headline-medium mt-0 mb-2">
                            {{ $t('pages.billing.cancelConfirmTitle') }}
                          </h2>
                          <p class="text-medium-emphasis" v-html="$t('pages.billing.cancelConfirmDesc').replace('{date}', new Date(authStore.user.subscription_renews_at!).toLocaleDateString())">
                          </p>
                        </VCardText>
                        <VCardActions class="mt-n6 px-6 pb-6">
                          <VBtn
                            color="primary"
                            variant="flat"
                            @click="dialog = false"
                          >
                            {{ $t('pages.billing.keepSubscription') }}
                          </VBtn>
                          <VSpacer />
                          <VBtn
                            v-if="!isOnFreePlan"
                            variant="text"
                            :loading="loading"
                            color="error"
                            @click="cancelPlan"
                          >
                            {{ $t('pages.billing.cancelPlan') }}
                          </VBtn>
                        </VCardActions>
                      </VCard>
                    </VDialog>
                  </div>
                </VAlert>
              </VCol>
            </VRow>

            <!-- Upgrade Plan (only for free users) -->
            <template v-if="isOnFreePlan">
              <h2 class="text-headline-large mt-4 mb-2">{{ $t('pages.billing.upgradePlanTitle') }}</h2>
              <VRow>
                <VCol cols="12" md="6">
                  <VCard :href="checkoutURL" link>
                    <VCardText>
                      <VRow align="center">
                        <VCol class="flex-grow-1 flex-shrink-1">
                          <h1
                            class="text-title-large font-weight-bold text-uppercase mt-3"
                          >
                            {{ $t('pages.billing.proPlanTitle') }}
                          </h1>
                          <p class="text-medium-emphasis">
                            {{ $t('pages.billing.proPlanFeatures') }}
                          </p>
                        </VCol>
                        <VCol class="flex-grow-0 flex-shrink-0 text-center">
                          <span class="text-headline-medium">$10</span>{{ $t('pages.billing.perMonth') }}
                        </VCol>
                      </VRow>
                    </VCardText>
                  </VCard>
                </VCol>
                <VCol cols="12" md="6">
                  <VCard :href="enterpriseCheckoutURL" link>
                    <VCardText>
                      <VRow align="center">
                        <VCol class="flex-grow-1 flex-shrink-1">
                          <h1
                            class="text-title-large font-weight-bold text-uppercase mt-3"
                          >
                            {{ $t('pages.billing.enterprisePlanTitle') }}
                          </h1>
                          <p class="text-medium-emphasis">
                            {{ $t('pages.billing.enterprisePlanFeatures') }}
                          </p>
                        </VCol>
                        <VCol class="flex-grow-0 flex-shrink-0 text-center">
                          <span class="text-headline-medium">$89</span>{{ $t('pages.billing.perMonth') }}
                        </VCol>
                      </VRow>
                    </VCardText>
                  </VCard>
                </VCol>
              </VRow>
            </template>
            </template>

            <!-- Overview & Usage Statistics -->
            <UsageStatistics />

            <!-- Subscription Payments -->
            <template v-if="authStore.user?.subscription_id != null">
              <h4 class="text-headline-large mb-3 mt-8">
                {{ $t('pages.billing.subscriptionPayments') }}
              </h4>
              <p class="text-medium-emphasis" v-html="$t('pages.billing.subscriptionPaymentsDesc')">
              </p>
              <VProgressCircular
                v-if="payments == null && loadingSubscriptionPayments"
                :size="20"
                :width="2"
                color="primary"
                indeterminate
              />
              <VTable v-if="payments">
                <thead>
                  <tr class="text-uppercase">
                    <th v-if="lgAndUp" class="text-left">{{ $t('pages.billing.thId') }}</th>
                    <th class="text-left">{{ $t('pages.billing.thTimestamp') }}</th>
                    <th class="text-left">{{ $t('pages.billing.thStatus') }}</th>
                    <th v-if="lgAndUp" class="text-left">{{ $t('pages.billing.thTax') }}</th>
                    <th class="text-left">{{ $t('pages.billing.thTotal') }}</th>
                    <th></th>
                  </tr>
                </thead>
                <tbody>
                  <tr v-for="payment in payments.data" :key="payment.id">
                    <td v-if="lgAndUp">{{ payment.id }}</td>
                    <td>
                      {{ formatTimestamp(payment.attributes.created_at) }}
                    </td>
                    <td>
                      <VChip
                        v-if="payment.attributes.status === 'paid'"
                        color="success"
                      >
                        <template #prepend>
                          <VIcon size="small" :icon="mdiCheck" />
                        </template>
                        {{ payment.attributes.status_formatted }}
                      </VChip>
                      <VChip v-else color="error">
                        <template #prepend>
                          <VIcon size="small" :icon="mdiAlert" />
                        </template>
                        {{ payment.attributes.status_formatted }}
                      </VChip>
                    </td>
                    <td v-if="lgAndUp">
                      {{ payment.attributes.tax_formatted }}
                    </td>
                    <td class="font-weight-bold">
                      {{ payment.attributes.total_formatted }}
                    </td>
                    <td class="text-right">
                      <VBtn
                        color="primary"
                        size="small"
                        @click="showInvoiceDialog(payment)"
                      >
                        <VIcon start :icon="mdiInvoice" />
                        {{ $t('pages.billing.btnInvoice') }}
                      </VBtn>
                    </td>
                  </tr>
                </tbody>
              </VTable>
            </template>

          </VCol>
        </VRow>
      </VContainer>
    </div>

    <!-- Invoice Dialog -->
    <VDialog
      v-model="subscriptionInvoiceDialog"
      persistent
      max-width="600"
      opacity="0.9"
    >
      <VCard>
        <VCardTitle class="text-headline-large">{{ $t('pages.billing.invoiceDialogTitle') }}</VCardTitle>
        <VCardSubtitle class="mt-n1" v-html="$t('pages.billing.invoiceDialogDesc').replace('{amount}', selectedPayment?.attributes.total_formatted || '').replace('{date}', formatTimestamp(selectedPayment?.attributes.created_at ?? ''))">
        </VCardSubtitle>
        <VCardText class="pb-0">
          <VContainer>
            <VRow>
              <VCol cols="12">
                <VTextField
                  v-model="invoiceFormName"
                  density="compact"
                  color="primary"
                  :disabled="loading"
                  :error="errorMessages.has('name')"
                  :error-messages="errorMessages.get('name')"
                  :label="$t('pages.billing.formNameLabel')"
                  :placeholder="$t('pages.billing.formNamePlaceholder')"
                  persistent-placeholder
                  variant="outlined"
                />
              </VCol>
              <VCol cols="12">
                <VTextField
                  v-model="invoiceFormAddress"
                  density="compact"
                  color="primary"
                  :disabled="loading"
                  :error="errorMessages.has('address')"
                  :error-messages="errorMessages.get('address')"
                  :label="$t('pages.billing.formAddressLabel')"
                  :placeholder="$t('pages.billing.formAddressPlaceholder')"
                  persistent-placeholder
                  variant="outlined"
                />
              </VCol>
            </VRow>
            <VRow>
              <VCol cols="6">
                <VTextField
                  v-model="invoiceFormCity"
                  density="compact"
                  :disabled="loading"
                  :error="errorMessages.has('city')"
                  :error-messages="errorMessages.get('city')"
                  :label="$t('pages.billing.formCityLabel')"
                  :placeholder="$t('pages.billing.formCityPlaceholder')"
                  persistent-placeholder
                  variant="outlined"
                />
              </VCol>
              <VCol cols="6">
                <VTextField
                  v-if="invoiceStateOptions.length === 0"
                  v-model="invoiceFormState"
                  density="compact"
                  color="primary"
                  :disabled="loading"
                  :error="errorMessages.has('state')"
                  :error-messages="errorMessages.get('state')"
                  :label="$t('pages.billing.formStateLabel')"
                  :placeholder="$t('pages.billing.formStatePlaceholder')"
                  persistent-placeholder
                  variant="outlined"
                />
                <VAutocomplete
                  v-else
                  v-model="invoiceFormState"
                  density="compact"
                  color="primary"
                  :disabled="loading"
                  :error="errorMessages.has('state')"
                  :error-messages="errorMessages.get('state')"
                  :items="invoiceStateOptions"
                  :label="$t('pages.billing.formStateLabel')"
                  :placeholder="$t('pages.billing.formStatePlaceholder')"
                  persistent-placeholder
                  variant="outlined"
                />
              </VCol>
            </VRow>
            <VRow>
              <VCol cols="6">
                <VTextField
                  v-model="invoiceFormZipCode"
                  density="compact"
                  color="primary"
                  :disabled="loading"
                  :error="errorMessages.has('zip_code')"
                  :error-messages="errorMessages.get('zip_code')"
                  :label="$t('pages.billing.formZipCodeLabel')"
                  :placeholder="$t('pages.billing.formZipCodePlaceholder')"
                  persistent-placeholder
                  variant="outlined"
                />
              </VCol>
              <VCol cols="6">
                <VAutocomplete
                  v-model="invoiceFormCountry"
                  density="compact"
                  color="primary"
                  :disabled="loading"
                  :error="errorMessages.has('country')"
                  :error-messages="errorMessages.get('country')"
                  :items="countries"
                  :label="$t('pages.billing.formCountryLabel')"
                  :placeholder="$t('pages.billing.formCountryPlaceholder')"
                  persistent-placeholder
                  variant="outlined"
                />
              </VCol>
            </VRow>
            <VRow>
              <VCol cols="12">
                <VTextarea
                  v-model="invoiceFormNotes"
                  density="compact"
                  color="primary"
                  :disabled="loading"
                  :error="errorMessages.has('notes')"
                  :error-messages="errorMessages.get('notes')"
                  rows="3"
                  :label="$t('pages.billing.formNotesLabel')"
                  :placeholder="$t('pages.billing.formNotesPlaceholder')"
                  persistent-placeholder
                  variant="outlined"
                />
              </VCol>
            </VRow>
          </VContainer>
        </VCardText>
        <VCardActions class="pb-4 mt-n4">
          <VBtn
            variant="flat"
            :loading="loading"
            color="primary"
            @click="generateInvoice"
          >
            <VIcon start :icon="mdiDownloadOutline" />
            {{ $t('pages.billing.btnDownloadInvoice') }}
          </VBtn>
          <VSpacer />
          <VBtn
            color="warning"
            variant="text"
            @click="subscriptionInvoiceDialog = false"
          >
            {{ $t('pages.billing.btnClose') }}
          </VBtn>
        </VCardActions>
      </VCard>
    </VDialog>
  </VContainer>
</template>
