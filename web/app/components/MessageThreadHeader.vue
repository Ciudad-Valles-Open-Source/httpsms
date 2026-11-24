<script setup lang="ts">
import { useDisplay } from 'vuetify'
import {
  mdiLogout,
  mdiDownload,
  mdiBatteryChargingHigh,
  mdiPackageUp,
  mdiPackageDown,
  mdiDotsVertical,
  mdiCellphone,
  mdiBookOpenVariant,
} from '@mdi/js'
import type { EntitiesPhone } from '~~/shared/types/api'

const { t } = useI18n()

const router = useRouter()
const route = useRoute()
const { lgAndUp } = useDisplay()
const authStore = useAuthStore()
const phonesStore = usePhonesStore()
const threadsStore = useThreadsStore()
const appStore = useAppStore()
const notificationsStore = useNotificationsStore()
const { formatPhoneNumber, phoneCountry, humanizeTime } = useFilters()
const { items: navItems, isActive, badgeFor } = usePanelNav()
const logout = useLogout()

const selectedMenuItem = ref(-1)

interface SelectItem {
  title: string
  value: string
}

const owners = computed<SelectItem[]>(() => {
  return phonesStore.phones.map((phone: EntitiesPhone) => ({
    title: formatPhoneNumber(phone.phone_number),
    value: phone.phone_number,
  }))
})

async function onOwnerChanged(owner: string) {
  await authStore.updateUser({ owner })
  if (route.name !== 'threads') {
    threadsStore.setThreadId(null)
    await router.push({ name: 'threads' })
    return
  }
  await threadsStore.loadThreads()
}

async function toggleArchive() {
  threadsStore.toggleArchive()
  setTimeout(() => {
    selectedMenuItem.value = -1
  }, 1000)
  if (route.name !== 'threads') {
    threadsStore.setThreadId(null)
    await router.push({ name: 'threads' })
    return
  }
  await threadsStore.loadThreads()
}

async function filterThreads(unreadOnly: boolean) {
  if (threadsStore.unreadOnly === unreadOnly) return

  const previous = threadsStore.unreadOnly
  threadsStore.unreadOnly = unreadOnly
  threadsStore.loadingThreads = true
  try {
    await threadsStore.loadThreads()
  } catch (error) {
    threadsStore.unreadOnly = previous
    notificationsStore.addNotification({
      type: 'error',
      message: t('components.messageThreadHeader.filterError'),
    })
    console.error(error)
  }
}
</script>

<template>
  <div class="thread-header flex-shrink-0">
    <div class="d-flex align-center ga-1">
      <h1 class="thread-header-title text-truncate">
        {{ $t('components.panelNav.threads') }}
      </h1>
      <v-progress-circular
        v-if="appStore.polling"
        indeterminate
        :size="16"
        :width="2"
        class="ml-2"
        color="primary"
        :aria-label="$t('components.messageThreadHeader.syncing')"
      />
      <v-spacer />
      <template v-if="!lgAndUp">
        <ThemeToggle with-palette />
        <LanguageSwitcher />
        <v-menu>
          <template #activator="{ props: menuProps }">
            <v-btn
              v-bind="menuProps"
              icon
              variant="text"
              size="small"
              :aria-label="$t('components.messageThreadHeader.menu')"
            >
              <v-icon :icon="mdiDotsVertical" />
            </v-btn>
          </template>
          <v-list density="compact" prepend-gap="16" min-width="240">
            <v-list-item
              v-for="item in navItems"
              :key="item.name"
              :to="{ name: item.name }"
              :active="isActive(item)"
              color="primary"
            >
              <template #prepend><v-icon :icon="item.icon" /></template>
              <v-list-item-title>{{ $t(item.label) }}</v-list-item-title>
              <template v-if="badgeFor(item)" #append>
                <v-badge
                  inline
                  color="warning"
                  :content="badgeFor(item)"
                  :aria-label="
                    $t('components.messageThreadHeader.pendingUsers', {
                      count: badgeFor(item),
                    })
                  "
                />
              </template>
            </v-list-item>
            <v-divider class="my-1" />
            <v-list-item
              v-if="phonesStore.owner"
              :href="appStore.appData.appDownloadUrl"
            >
              <template #prepend><v-icon :icon="mdiDownload" /></template>
              <v-list-item-title>{{
                $t('components.messageThreadHeader.downloadApp')
              }}</v-list-item-title>
            </v-list-item>
            <v-list-item to="/docs">
              <template #prepend
                ><v-icon :icon="mdiBookOpenVariant"
              /></template>
              <v-list-item-title>{{
                $t('components.messageThreadHeader.documentation')
              }}</v-list-item-title>
            </v-list-item>
            <v-list-item @click.prevent="logout">
              <template #prepend><v-icon :icon="mdiLogout" /></template>
              <v-list-item-title>{{
                $t('components.messageThreadHeader.logout')
              }}</v-list-item-title>
            </v-list-item>
          </v-list>
        </v-menu>
      </template>
    </div>

    <v-select
      variant="outlined"
      density="compact"
      class="mt-4"
      hide-details
      :prepend-inner-icon="mdiCellphone"
      :disabled="owners.length === 0"
      :placeholder="$t('components.messageThreadHeader.phoneNumbers')"
      :aria-label="$t('components.messageThreadHeader.phoneNumbers')"
      :items="owners"
      :model-value="phonesStore.owner"
      @update:model-value="onOwnerChanged"
    />
    <div
      v-if="phonesStore.owner"
      class="d-flex align-center flex-wrap ga-2 mt-2 text-body-small"
    >
      <span class="text-medium-emphasis">
        {{ phoneCountry(phonesStore.owner) }}
      </span>
      <v-tooltip v-if="phonesStore.heartbeat" location="bottom">
        <template #activator="{ props: tooltipProps }">
          <v-chip
            v-bind="tooltipProps"
            size="x-small"
            color="success"
            variant="tonal"
            class="font-weight-bold"
            :to="{
              name: 'heartbeats-id',
              params: { id: phonesStore.owner },
            }"
          >
            <v-icon
              v-if="phonesStore.heartbeat.charging"
              start
              :icon="mdiBatteryChargingHigh"
            />
            <span v-else class="status-dot mr-2" />
            {{
              $t('components.messageThreadHeader.timeAgo', {
                time: humanizeTime(phonesStore.heartbeat.timestamp),
              })
            }}
          </v-chip>
        </template>
        {{ $t('components.messageThreadHeader.lastHeartbeat') }}
      </v-tooltip>
    </div>

    <div class="d-flex align-center ga-2 pt-4">
      <div
        role="group"
        :aria-label="$t('components.messageThreadHeader.filtersLabel')"
        class="d-flex ga-2"
      >
        <v-btn
          size="small"
          rounded="pill"
          :color="!threadsStore.unreadOnly ? 'primary' : undefined"
          :variant="!threadsStore.unreadOnly ? 'flat' : 'outlined'"
          :class="{ 'text-medium-emphasis': threadsStore.unreadOnly }"
          :aria-pressed="!threadsStore.unreadOnly"
          :disabled="!phonesStore.owner || threadsStore.loadingThreads"
          @click="filterThreads(false)"
        >
          {{ $t('components.messageThreadHeader.filterAll') }}
        </v-btn>
        <v-btn
          size="small"
          rounded="pill"
          :color="threadsStore.unreadOnly ? 'primary' : undefined"
          :variant="threadsStore.unreadOnly ? 'flat' : 'outlined'"
          :class="{ 'text-medium-emphasis': !threadsStore.unreadOnly }"
          :aria-pressed="threadsStore.unreadOnly"
          :disabled="!phonesStore.owner || threadsStore.loadingThreads"
          @click="filterThreads(true)"
        >
          {{ $t('components.messageThreadHeader.filterUnread') }}
        </v-btn>
      </div>
      <v-spacer />
      <v-tooltip location="bottom">
        <template #activator="{ props: tooltipProps }">
          <v-btn
            v-bind="tooltipProps"
            size="small"
            rounded="pill"
            :color="threadsStore.archivedThreads ? 'warning' : undefined"
            :variant="threadsStore.archivedThreads ? 'tonal' : 'text'"
            :class="{ 'text-medium-emphasis': !threadsStore.archivedThreads }"
            :aria-pressed="threadsStore.archivedThreads"
            :prepend-icon="
              threadsStore.archivedThreads ? mdiPackageUp : mdiPackageDown
            "
            @click.prevent="toggleArchive"
          >
            {{ $t('components.messageThreadHeader.archivedFilter') }}
          </v-btn>
        </template>
        {{
          threadsStore.archivedThreads
            ? $t('components.messageThreadHeader.showInbox')
            : $t('components.messageThreadHeader.showArchived')
        }}
      </v-tooltip>
    </div>
  </div>
</template>

<style scoped>
.thread-header {
  padding: 1.1rem 1rem 0.9rem;
  border-bottom: 1px solid rgb(var(--v-border-color), var(--v-border-opacity));
  background: rgb(var(--v-theme-surface));
}

.thread-header-title {
  margin: 0;
  font-size: 1.35rem;
  font-weight: 800;
  letter-spacing: -0.02em;
}
</style>
