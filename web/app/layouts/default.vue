<script setup lang="ts">
import { useDisplay } from 'vuetify'
import { setAuthHeader } from '~/composables/useApi'
import { bindRealtimeEvents } from '~/composables/useRealtime'
import { getAuth } from 'firebase/auth'

useThemeScope('dashboard')

const route = useRoute()
const { lgAndUp } = useDisplay()
const authStore = useAuthStore()
const phonesStore = usePhonesStore()
const threadsStore = useThreadsStore()
const appStore = useAppStore()

// Events pushed by the API (through Soketi) which change the thread list.
const REALTIME_EVENTS = [
  'phone.updated',
  'message.phone.received',
  'message.phone.sent',
  'message.send.failed',
  'message.call.missed',
]

let unbindRealtime: () => void = () => {}
let refreshTimer: ReturnType<typeof setTimeout> | null = null

const hasDrawer = computed(() => {
  return ['threads', 'threads-id'].includes((route.name as string) ?? '')
})

// Several events usually arrive together, so refresh once after they settle.
function scheduleRefresh() {
  if (refreshTimer) clearTimeout(refreshTimer)
  refreshTimer = setTimeout(() => void refresh(), 500)
}

async function refresh() {
  if (authStore.authUser == null || !phonesStore.owner) return

  appStore.setPolling(true)
  try {
    const token = await getAuth().currentUser?.getIdToken()
    if (token) setAuthHeader(token)

    await Promise.all([
      phonesStore.loadPhones(true),
      threadsStore.loadThreads(),
      phonesStore.getHeartbeat(),
    ])
  } catch (error) {
    console.error(error)
  } finally {
    setTimeout(() => appStore.setPolling(false), 1000)
  }
}

onMounted(() => {
  watch(
    () => authStore.authUser?.id,
    (userId) => {
      unbindRealtime()
      unbindRealtime = bindRealtimeEvents(
        userId,
        Object.fromEntries(REALTIME_EVENTS.map((e) => [e, scheduleRefresh])),
      )
    },
    { immediate: true },
  )
})

onBeforeUnmount(() => {
  unbindRealtime()
  if (refreshTimer) clearTimeout(refreshTimer)
})
</script>

<template>
  <v-app>
    <v-navigation-drawer
      v-if="lgAndUp && authStore.authUser"
      :width="76"
      permanent
      color="background"
      class="panel-rail-drawer"
    >
      <PanelRail />
    </v-navigation-drawer>
    <v-navigation-drawer
      v-if="lgAndUp && hasDrawer"
      :width="380"
      permanent
      color="surface"
    >
      <template #default>
        <div class="message-thread-sidebar d-flex flex-column">
          <div v-if="appStore.isLocal" class="local-env-bar bg-warning" />
          <MessageThreadHeader />
          <div class="overflow-y-auto v-navigation-drawer__message-thread">
            <MessageThread />
          </div>
        </div>
      </template>
    </v-navigation-drawer>
    <v-main :class="{ 'has-drawer': hasDrawer && lgAndUp }">
      <AppToast />
      <slot v-if="authStore.authStateChanged" />
      <LoadingDashboard v-else />
    </v-main>
  </v-app>
</template>

<style lang="scss">
.v-application {
  .w-full {
    width: 100%;
  }

  .h-full {
    height: 100%;
  }

  .has-drawer {
    .v-snackbar {
      padding-left: 456px;
    }
  }

  .message-thread-sidebar {
    height: 100%;
  }

  .local-env-bar {
    height: 4px;
  }

  .v-navigation-drawer__message-thread {
    flex: 1;
    min-height: 0;
  }

  .hover\:text-decoration-underline {
    &:hover {
      text-decoration: underline;
    }
  }

  code.hljs {
    font-size: 16px;
  }
}
</style>
