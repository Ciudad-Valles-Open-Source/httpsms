import { getAuth, signOut } from 'firebase/auth'

/** Signs the user out and clears every store holding account data. */
export function useLogout() {
  const { t } = useI18n()
  const router = useRouter()
  const authStore = useAuthStore()
  const phonesStore = usePhonesStore()
  const threadsStore = useThreadsStore()
  const notificationsStore = useNotificationsStore()
  const redirectPreferenceStore = useRedirectPreferenceStore()
  const adminStore = useAdminStore()

  return async function logout() {
    const auth = getAuth()
    await signOut(auth)
    authStore.resetState()
    phonesStore.resetState()
    threadsStore.resetState()
    adminStore.resetState()
    redirectPreferenceStore.resetState()
    notificationsStore.addNotification({
      type: 'info',
      message: t('components.messageThreadHeader.logoutSuccess'),
    })
    router.push({ name: 'index' })
  }
}
