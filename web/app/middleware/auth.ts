import { useAuthStore } from '../stores/auth'
import {
  ACCOUNT_PENDING_PATH,
  accountStatusRedirect,
} from '../utils/account-status'

export default defineNuxtRouteMiddleware(async (to: { path: string }) => {
  const authStore = useAuthStore()

  if (!authStore.authStateChanged) {
    await new Promise<void>((resolve) => {
      const stop = watch(
        () => authStore.authStateChanged,
        (changed) => {
          if (changed) {
            stop()
            resolve()
          }
        },
        { immediate: true },
      )
    })
  }

  if (authStore.authUser === null) {
    return navigateTo({ path: '/login', query: { to: to.path } })
  }

  // Pending and suspended accounts may only see the account status page. The
  // status page always reloads the user so an approval is picked up.
  if (!authStore.user || to.path === ACCOUNT_PENDING_PATH) {
    try {
      await authStore.loadUser()
    } catch (error) {
      console.error(error)
    }
  }

  const redirect = accountStatusRedirect(authStore.user, to.path)
  if (redirect) {
    return navigateTo(redirect)
  }
})
