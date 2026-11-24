import { useAuthStore } from '../stores/auth'
import { DASHBOARD_PATH } from '../utils/account-status'

/**
 * Restricts a page to administrators. Use it after the `auth` middleware,
 * which loads the signed-in user.
 */
export default defineNuxtRouteMiddleware(() => {
  const authStore = useAuthStore()
  if (!authStore.isAdmin) {
    return navigateTo(DASHBOARD_PATH)
  }
})
