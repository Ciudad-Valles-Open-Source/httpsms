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
  <div class="auth-page">
    <div class="auth-page-top">
      <AppLogo />
      <div class="d-flex align-center">
        <ThemeToggle />
        <LanguageSwitcher />
      </div>
    </div>
    <main class="auth-page-main">
      <div class="auth-card card-soft">
        <span
          class="icon-tile icon-tile-lg"
          :class="suspended ? 'status-tile-error' : 'status-tile-warning'"
        >
          <VIcon
            size="36"
            :icon="suspended ? mdiAccountCancel : mdiAccountClock"
          />
        </span>
        <h1 class="auth-card-title status-title mt-5">
          {{
            suspended
              ? $t('pages.accountPending.suspendedTitle')
              : $t('pages.accountPending.pendingTitle')
          }}
        </h1>
        <p class="text-body-medium text-medium-emphasis mt-3 mb-0">
          {{
            suspended
              ? $t('pages.accountPending.suspendedBody')
              : $t('pages.accountPending.pendingBody')
          }}
        </p>
        <div v-if="email" class="signed-in-as mt-5">
          <span class="text-body-small text-medium-emphasis">
            {{ $t('pages.accountPending.signedInAs') }}
          </span>
          <strong class="text-body-medium text-break">{{ email }}</strong>
        </div>
        <div class="d-flex flex-column ga-3 mt-6">
          <VBtn
            block
            size="large"
            class="btn-gradient"
            :prepend-icon="mdiRefresh"
            :loading="refreshing"
            @click="refreshStatus"
          >
            {{ $t('pages.accountPending.refresh') }}
          </VBtn>
          <VBtn
            block
            size="large"
            variant="outlined"
            :prepend-icon="mdiLogout"
            :loading="loggingOut"
            @click="logout"
          >
            {{ $t('pages.accountPending.logout') }}
          </VBtn>
        </div>
        <VDivider class="my-5" />
        <VBtn
          color="error"
          variant="text"
          size="small"
          :prepend-icon="mdiDelete"
          :loading="deletingAccount"
          @click="showDeleteAccountDialog = true"
        >
          {{ $t('pages.accountPending.deleteAccount') }}
        </VBtn>
      </div>
    </main>
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
  </div>
</template>

<style scoped>
.status-title {
  font-size: 1.6rem;
}

.status-tile-warning {
  background: rgb(var(--v-theme-warning), 0.14);
  color: rgb(var(--v-theme-warning));
}

.status-tile-error {
  background: rgb(var(--v-theme-error), 0.14);
  color: rgb(var(--v-theme-error));
}

.signed-in-as {
  display: flex;
  flex-direction: column;
  gap: 0.15rem;
  padding: 0.75rem 1rem;
  border: 1px solid rgb(var(--v-border-color), var(--v-border-opacity));
  border-radius: 0.9rem;
  background: rgb(var(--v-theme-on-surface), 0.03);
}
</style>
