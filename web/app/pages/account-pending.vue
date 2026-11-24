<script setup lang="ts">
import {
  mdiAccountClock,
  mdiAccountCancel,
  mdiRefresh,
  mdiLogout,
  mdiDelete,
} from '@mdi/js'
import { getAuth, signOut } from 'firebase/auth'
import { getApiErrorMessage } from '~/utils/api-error'
import { DASHBOARD_PATH } from '~/utils/account-status'

definePageMeta({
  middleware: ['auth'],
  layout: 'blank',
})

const { t } = useI18n()

useHead({
  title: () => t('pages.accountPending.metaTitle'),
})

const router = useRouter()
const { smAndUp } = useVDisplay()
const authStore = useAuthStore()
const phonesStore = usePhonesStore()
const threadsStore = useThreadsStore()
const adminStore = useAdminStore()
const notificationsStore = useNotificationsStore()
const redirectPreferenceStore = useRedirectPreferenceStore()

const suspended = computed(() => authStore.user?.status === 'suspended')
const email = computed(
  () => authStore.user?.email || authStore.authUser?.email || '',
)

const refreshing = ref(false)
const loggingOut = ref(false)
const deletingAccount = ref(false)
const showDeleteAccountDialog = ref(false)

// Leave as soon as the account becomes usable (e.g. after a refresh).
watch(
  () => authStore.user && !authStore.accountBlocked,
  (usable) => {
    if (usable) router.replace(DASHBOARD_PATH)
  },
)

async function refreshStatus() {
  refreshing.value = true
  try {
    await authStore.loadUser()
    if (authStore.accountBlocked) {
      notificationsStore.addNotification({
        type: 'info',
        message: suspended.value
          ? t('pages.accountPending.stillSuspended')
          : t('pages.accountPending.stillPending'),
      })
    }
  } catch (error: unknown) {
    notificationsStore.addNotification({
      type: 'error',
      message: getApiErrorMessage(
        error,
        t('pages.accountPending.refreshError'),
      ),
    })
  } finally {
    refreshing.value = false
  }
}

function resetStores() {
  authStore.resetState()
  phonesStore.resetState()
  threadsStore.resetState()
  adminStore.resetState()
  redirectPreferenceStore.resetState()
}

async function logout() {
  loggingOut.value = true
  try {
    await signOut(getAuth())
    resetStores()
    notificationsStore.addNotification({
      type: 'info',
      message: t('settings.logged_out'),
    })
    await router.push({ name: 'index' })
  } finally {
    loggingOut.value = false
  }
}

// Same flow as "Delete your Account" in Settings.
async function deleteUserAccount() {
  deletingAccount.value = true
  try {
    await authStore.deleteUserAccount()
    notificationsStore.addNotification({
      message: t('settings.account_deleted'),
      type: 'success',
    })
    await signOut(getAuth())
    resetStores()
    await router.push({ name: 'index' })
  } catch (error: unknown) {
    notificationsStore.addNotification({
      message: getApiErrorMessage(error, t('settings.failed_delete_account')),
      type: 'error',
    })
  } finally {
    deletingAccount.value = false
    showDeleteAccountDialog.value = false
  }
}
</script>

<template>
  <VContainer class="fill-height d-flex align-center justify-center">
    <VRow align="center" justify="center">
      <VCol cols="12" sm="10" md="7" lg="5" xl="4">
        <div class="text-center mb-6">
          <VAvatar color="#121212" size="45" rounded="0">
            <VImg :src="'/img/logo.svg'" />
          </VAvatar>
        </div>
        <VCard class="pa-2 pa-sm-4">
          <VCardText class="text-center">
            <VIcon
              :icon="suspended ? mdiAccountCancel : mdiAccountClock"
              :color="suspended ? 'error' : 'warning'"
              size="64"
              class="mb-4"
            />
            <h1 class="text-headline-medium mb-3">
              {{
                suspended
                  ? $t('pages.accountPending.suspendedTitle')
                  : $t('pages.accountPending.pendingTitle')
              }}
            </h1>
            <p class="text-body-large text-medium-emphasis mb-4">
              {{
                suspended
                  ? $t('pages.accountPending.suspendedBody')
                  : $t('pages.accountPending.pendingBody')
              }}
            </p>
            <p v-if="email" class="text-body-medium mb-0">
              {{ $t('pages.accountPending.signedInAs') }}
              <strong class="text-break">{{ email }}</strong>
            </p>
          </VCardText>
          <VCardActions class="flex-column flex-sm-row ga-2 px-4 pb-4">
            <VBtn
              color="primary"
              variant="flat"
              :block="!smAndUp"
              :loading="refreshing"
              @click="refreshStatus"
            >
              <VIcon start :icon="mdiRefresh" />
              {{ $t('pages.accountPending.refresh') }}
            </VBtn>
            <VBtn
              variant="tonal"
              :block="!smAndUp"
              :loading="loggingOut"
              @click="logout"
            >
              <VIcon start :icon="mdiLogout" />
              {{ $t('pages.accountPending.logout') }}
            </VBtn>
            <VSpacer v-if="smAndUp" />
            <VBtn
              color="error"
              variant="text"
              :block="!smAndUp"
              :loading="deletingAccount"
              @click="showDeleteAccountDialog = true"
            >
              <VIcon start :icon="mdiDelete" />
              {{ $t('pages.accountPending.deleteAccount') }}
            </VBtn>
          </VCardActions>
        </VCard>
      </VCol>
    </VRow>
    <VDialog v-model="showDeleteAccountDialog" max-width="600px">
      <VCard>
        <VCardTitle class="text-center text-wrap">{{
          $t('settings.delete_account_title')
        }}</VCardTitle>
        <VCardText class="mt-2 text-center text-medium-emphasis">
          {{ $t('settings.delete_account_confirm') }}
        </VCardText>
        <VCardActions>
          <VBtn
            color="error"
            variant="text"
            :loading="deletingAccount"
            @click="deleteUserAccount"
          >
            <VIcon v-if="smAndUp" start :icon="mdiDelete" />
            {{ $t('settings.delete_my_account') }}
          </VBtn>
          <VSpacer />
          <VBtn
            color="primary"
            variant="flat"
            @click="showDeleteAccountDialog = false"
          >
            {{ $t('settings.keep_my_account') }}
          </VBtn>
        </VCardActions>
      </VCard>
    </VDialog>
  </VContainer>
</template>
