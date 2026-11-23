<script setup lang="ts">
import { useDisplay } from 'vuetify'
import { setAuthHeader } from '~/composables/useApi'
import { bindRealtimeEvents } from '~/composables/useRealtime'
import { getAuth } from 'firebase/auth'

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
    <v-navigation-drawer v-if="lgAndUp && hasDrawer" :width="400" permanent>
      <template #default>
        <div class="message-thread-sidebar d-flex flex-column">
          <div v-if="appStore.isLocal" class="py-2 bg-warning" />
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
      padding-left: 400px;
    }
  }
  .message-thread-sidebar {
    height: 100%;
  }
  .v-navigation-drawer__message-thread {
    flex: 1;
    min-height: 0;
    &::-webkit-scrollbar {
      width: 8px;
    }
    &::-webkit-scrollbar-track {
      background: #363636;
    }
    &::-webkit-scrollbar-thumb {
      background: #666666;
      border-radius: 8px;
    }
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
