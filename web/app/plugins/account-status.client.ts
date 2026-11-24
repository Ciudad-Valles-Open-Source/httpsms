import { setAccountBlockedHandler } from '../composables/useApi'
import { ACCOUNT_PENDING_PATH } from '../utils/account-status'

/**
 * Sends the user to the account status page whenever an API call is rejected
 * because their account is pending approval or suspended.
 */
export default defineNuxtPlugin(() => {
  const router = useRouter()
  const authStore = useAuthStore()

  setAccountBlockedHandler((code) => {
    authStore.markAccountBlocked(code)
    if (router.currentRoute.value.path !== ACCOUNT_PENDING_PATH) {
      router.replace(ACCOUNT_PENDING_PATH).catch(console.error)
    }
  })
})
