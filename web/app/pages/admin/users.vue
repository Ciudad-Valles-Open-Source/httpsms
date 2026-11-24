<script setup lang="ts">
import {
  mdiArrowLeft,
  mdiMagnify,
  mdiRefresh,
  mdiAccountSearch,
  mdiShieldAccount,
  mdiCalendarAccount,
  mdiCheckDecagram,
  mdiCellphone,
} from '@mdi/js'
import { getApiErrorMessage } from '~/utils/api-error'
import {
  ADMIN_USER_FILTERS,
  actionNeedsConfirmation,
  formatCountBadge,
  isProtectedAdminUser,
  type AdminUser,
  type AdminUserAction,
  type AdminUserFilter,
  type AdminUserStatus,
} from '~/utils/admin-users'
import { PENDING_COUNT_LIMIT } from '~/stores/admin'

definePageMeta({
  middleware: ['auth', 'admin'],
})

const { t, locale } = useI18n()

useHead({
  title: () => t('pages.adminUsers.metaTitle'),
})

// The table becomes a list of cards whenever its columns don't fit.
const { target: usersTable, stacked: usersStacked } = useTableFit()
const authStore = useAuthStore()
const adminStore = useAdminStore()
const notificationsStore = useNotificationsStore()

const SEARCH_DEBOUNCE_MS = 400

const STATUS_COLOR: Record<AdminUserStatus, string> = {
  pending: 'warning',
  active: 'success',
  suspended: 'error',
}

const search = ref(adminStore.query)
let searchTimer: ReturnType<typeof setTimeout> | null = null

watch(search, (value) => {
  if (searchTimer) clearTimeout(searchTimer)
  searchTimer = setTimeout(() => {
    searchTimer = null
    if ((value ?? '').trim() !== adminStore.query.trim()) {
      adminStore.setQuery(value ?? '')
    }
  }, SEARCH_DEBOUNCE_MS)
})

const pendingBadge = computed(() =>
  formatCountBadge(adminStore.pendingCount, PENDING_COUNT_LIMIT),
)

const loadErrorMessage = computed(() =>
  adminStore.error
    ? getApiErrorMessage(adminStore.error, t('pages.adminUsers.loadError'))
    : '',
)

const emptyMessage = computed(() =>
  adminStore.query.trim()
    ? t('pages.adminUsers.empty.search', { query: adminStore.query.trim() })
    : t(`pages.adminUsers.empty.${adminStore.filter}`),
)

const busy = ref<{ id: string; action: AdminUserAction } | null>(null)
const confirmation = ref<{ user: AdminUser; action: AdminUserAction } | null>(
  null,
)
const confirmationOpen = computed({
  get: () => confirmation.value !== null,
  set: (open: boolean) => {
    if (!open && !busy.value) confirmation.value = null
  },
})

function isProtected(user: AdminUser): boolean {
  return isProtectedAdminUser(user, authStore.user?.id)
}

function busyAction(user: AdminUser): AdminUserAction | null {
  return busy.value?.id === user.id ? busy.value.action : null
}

function formatDate(value?: string | null): string {
  if (!value) return ''
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) return value
  return date.toLocaleDateString(locale.value, {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
  })
}

function initial(email: string): string {
  return (email.trim()[0] ?? '?').toUpperCase()
}

function selectFilter(value: AdminUserFilter | null | undefined) {
  if (!value || value === adminStore.filter) return
  adminStore.setFilter(value)
}

function onAction(user: AdminUser, action: AdminUserAction) {
  if (isProtected(user) || busy.value) return
  if (actionNeedsConfirmation(action)) {
    confirmation.value = { user, action }
    return
  }
  runAction(user, action)
}

async function runAction(user: AdminUser, action: AdminUserAction) {
  busy.value = { id: user.id, action }
  try {
    if (action === 'delete') {
      await adminStore.deleteUser(user)
    } else {
      await adminStore.changeStatus(user, action)
    }
    notificationsStore.addNotification({
      type: 'success',
      message: t(`pages.adminUsers.success.${action}`, { email: user.email }),
    })
  } catch (error: unknown) {
    notificationsStore.addNotification({
      type: 'error',
      message: getApiErrorMessage(error, t('pages.adminUsers.actionError')),
    })
  } finally {
    busy.value = null
    confirmation.value = null
  }
}

async function loadMore() {
  try {
    await adminStore.loadMore()
  } catch (error: unknown) {
    notificationsStore.addNotification({
      type: 'error',
      message: getApiErrorMessage(error, t('pages.adminUsers.loadError')),
    })
  }
}

function reload() {
  adminStore.loadUsers()
  adminStore.loadPendingCount().catch(console.error)
}

onMounted(reload)

onBeforeUnmount(() => {
  if (searchTimer) clearTimeout(searchTimer)
})
</script>

<template>
  <VContainer fluid class="px-0 pt-0">
    <div class="w-100">
      <VAppBar>
        <VBtn icon to="/threads" :aria-label="$t('pages.adminUsers.back')">
          <VIcon :icon="mdiArrowLeft" />
        </VBtn>
        <VToolbarTitle>{{ $t('pages.adminUsers.title') }}</VToolbarTitle>
        <VBtn
          icon
          :disabled="adminStore.loading"
          :aria-label="$t('pages.adminUsers.retry')"
          @click="reload"
        >
          <VIcon :icon="mdiRefresh" />
        </VBtn>
        <VProgressLinear
          color="primary"
          :active="adminStore.loading"
          indeterminate
          absolute
          location="bottom"
        />
      </VAppBar>
      <VContainer class="pt-6 px-sm-4">
        <VRow>
          <VCol cols="12">
            <header class="admin-header">
              <span class="icon-tile icon-tile-gradient">
                <VIcon :icon="mdiShieldAccount" />
              </span>
              <div>
                <h1 class="admin-title">
                  {{ $t('pages.adminUsers.heading') }}
                </h1>
                <p class="text-body-medium text-medium-emphasis mb-0">
                  {{ $t('pages.adminUsers.description') }}
                </p>
              </div>
            </header>

            <section class="panel-section">
              <div class="admin-toolbar">
                <div
                  role="group"
                  :aria-label="$t('pages.adminUsers.columns.status')"
                  class="d-flex flex-wrap ga-2"
                >
                  <VBtn
                    v-for="item in ADMIN_USER_FILTERS"
                    :key="item"
                    size="small"
                    rounded="pill"
                    :color="adminStore.filter === item ? 'primary' : undefined"
                    :variant="adminStore.filter === item ? 'flat' : 'outlined'"
                    :class="{
                      'text-medium-emphasis': adminStore.filter !== item,
                    }"
                    :aria-pressed="adminStore.filter === item"
                    @click="selectFilter(item)"
                  >
                    {{ $t(`pages.adminUsers.filters.${item}`) }}
                    <VBadge
                      v-if="item === 'pending' && pendingBadge"
                      :content="pendingBadge"
                      color="warning"
                      inline
                      class="ml-1"
                    />
                  </VBtn>
                </div>
                <VTextField
                  v-model="search"
                  :label="$t('pages.adminUsers.searchLabel')"
                  :prepend-inner-icon="mdiMagnify"
                  variant="outlined"
                  density="compact"
                  color="primary"
                  type="search"
                  clearable
                  hide-details
                  class="admin-users-search"
                />
              </div>

              <VAlert
                v-if="loadErrorMessage"
                type="error"
                variant="tonal"
                class="mt-5"
              >
                {{ loadErrorMessage }}
                <template #append>
                  <VBtn variant="text" color="error" @click="reload">
                    {{ $t('pages.adminUsers.retry') }}
                  </VBtn>
                </template>
              </VAlert>

              <div
                v-else-if="adminStore.loading && !adminStore.users.length"
                class="mt-4"
              >
                <VSkeletonLoader
                  v-for="n in 3"
                  :key="n"
                  type="list-item-avatar-two-line"
                  color="transparent"
                />
              </div>

              <div
                v-else-if="adminStore.loaded && !adminStore.users.length"
                class="empty-state"
              >
                <span class="icon-tile icon-tile-lg">
                  <VIcon size="32" :icon="mdiAccountSearch" />
                </span>
                <p class="empty-state-text mt-4 mb-0">
                  {{ emptyMessage }}
                </p>
              </div>

              <div v-else-if="adminStore.users.length" ref="usersTable">
                <VTable v-if="!usersStacked" density="comfortable" class="mt-5">
                  <thead>
                    <tr>
                      <th>{{ $t('pages.adminUsers.columns.email') }}</th>
                      <th>{{ $t('pages.adminUsers.columns.status') }}</th>
                      <th>{{ $t('pages.adminUsers.columns.created') }}</th>
                      <th>{{ $t('pages.adminUsers.columns.approved') }}</th>
                      <th class="text-center">
                        {{ $t('pages.adminUsers.columns.phones') }}
                      </th>
                      <th class="text-end">
                        {{ $t('pages.adminUsers.columns.actions') }}
                      </th>
                    </tr>
                  </thead>
                  <tbody>
                    <tr v-for="user in adminStore.users" :key="user.id">
                      <td>
                        <div class="d-flex align-center ga-3 py-2">
                          <VAvatar
                            size="34"
                            color="primary"
                            variant="tonal"
                            class="user-initial"
                          >
                            {{ initial(user.email) }}
                          </VAvatar>
                          <div class="user-text">
                            <div class="font-weight-medium user-email">
                              {{ user.email }}
                            </div>
                            <VChip
                              v-if="user.is_admin"
                              size="x-small"
                              color="primary"
                              variant="tonal"
                              class="mt-1"
                            >
                              {{ $t('pages.adminUsers.adminBadge') }}
                            </VChip>
                            <VChip
                              v-else-if="user.id === authStore.user?.id"
                              size="x-small"
                              variant="tonal"
                              class="mt-1"
                            >
                              {{ $t('pages.adminUsers.youBadge') }}
                            </VChip>
                          </div>
                        </div>
                      </td>
                      <td>
                        <VChip
                          size="small"
                          variant="tonal"
                          :color="STATUS_COLOR[user.status]"
                        >
                          <span class="chip-dot" />
                          {{ $t(`pages.adminUsers.status.${user.status}`) }}
                        </VChip>
                      </td>
                      <td class="text-no-wrap">
                        {{ formatDate(user.created_at) }}
                      </td>
                      <td>
                        <template v-if="user.approved_at">
                          <div class="text-no-wrap">
                            {{ formatDate(user.approved_at) }}
                          </div>
                          <div
                            v-if="user.approved_by"
                            class="text-body-small text-medium-emphasis text-break"
                          >
                            {{
                              $t('pages.adminUsers.approvedBy', {
                                name: user.approved_by,
                              })
                            }}
                          </div>
                        </template>
                        <span v-else class="text-medium-emphasis">—</span>
                      </td>
                      <td class="text-center">
                        {{ user.phones_count ?? 0 }}
                      </td>
                      <td class="text-end py-2">
                        <span
                          v-if="isProtected(user)"
                          class="text-body-small text-medium-emphasis"
                        >
                          {{ $t('pages.adminUsers.protectedHint') }}
                        </span>
                        <AdminUserActions
                          v-else
                          compact
                          class="flex-nowrap"
                          :user="user"
                          :disabled="!!busy && busy.id !== user.id"
                          :busy-action="busyAction(user)"
                          @action="onAction(user, $event)"
                        />
                      </td>
                    </tr>
                  </tbody>
                </VTable>

                <ul v-else class="user-list">
                  <li
                    v-for="user in adminStore.users"
                    :key="user.id"
                    class="user-item"
                  >
                    <div class="d-flex align-start ga-3">
                      <VAvatar
                        size="40"
                        color="primary"
                        variant="tonal"
                        class="user-initial"
                      >
                        {{ initial(user.email) }}
                      </VAvatar>
                      <div class="flex-grow-1 user-text">
                        <div class="font-weight-bold text-break">
                          {{ user.email }}
                        </div>
                        <div class="d-flex flex-wrap ga-1 mt-1">
                          <VChip
                            size="small"
                            variant="tonal"
                            :color="STATUS_COLOR[user.status]"
                          >
                            <span class="chip-dot" />
                            {{ $t(`pages.adminUsers.status.${user.status}`) }}
                          </VChip>
                          <VChip
                            v-if="user.is_admin"
                            size="small"
                            color="primary"
                            variant="tonal"
                          >
                            {{ $t('pages.adminUsers.adminBadge') }}
                          </VChip>
                          <VChip
                            v-else-if="user.id === authStore.user?.id"
                            size="small"
                            variant="tonal"
                          >
                            {{ $t('pages.adminUsers.youBadge') }}
                          </VChip>
                        </div>
                        <ul class="user-meta">
                          <li>
                            <VIcon size="16" :icon="mdiCalendarAccount" />
                            {{
                              $t('pages.adminUsers.createdOn', {
                                date: formatDate(user.created_at),
                              })
                            }}
                          </li>
                          <li v-if="user.approved_at">
                            <VIcon size="16" :icon="mdiCheckDecagram" />
                            <span>
                              {{
                                $t('pages.adminUsers.approvedOn', {
                                  date: formatDate(user.approved_at),
                                })
                              }}
                              <template v-if="user.approved_by">
                                ·
                                {{
                                  $t('pages.adminUsers.approvedBy', {
                                    name: user.approved_by,
                                  })
                                }}
                              </template>
                            </span>
                          </li>
                          <li>
                            <VIcon size="16" :icon="mdiCellphone" />
                            {{
                              $t('pages.adminUsers.phonesCount', {
                                count: user.phones_count ?? 0,
                              })
                            }}
                          </li>
                        </ul>
                      </div>
                    </div>
                    <AdminUserActions
                      v-if="!isProtected(user)"
                      class="mt-3"
                      :user="user"
                      :disabled="!!busy && busy.id !== user.id"
                      :busy-action="busyAction(user)"
                      @action="onAction(user, $event)"
                    />
                  </li>
                </ul>

                <div v-if="adminStore.hasMore" class="text-center mt-5">
                  <VBtn
                    variant="tonal"
                    color="primary"
                    :loading="adminStore.loadingMore"
                    @click="loadMore"
                  >
                    {{ $t('pages.adminUsers.loadMore') }}
                  </VBtn>
                </div>
              </div>
            </section>
          </VCol>
        </VRow>
      </VContainer>
    </div>

    <VDialog v-model="confirmationOpen" max-width="500px">
      <VCard v-if="confirmation">
        <VCardTitle class="text-wrap">
          {{
            confirmation.action === 'delete'
              ? $t('pages.adminUsers.confirm.deleteTitle')
              : $t('pages.adminUsers.confirm.suspendTitle')
          }}
        </VCardTitle>
        <VCardText class="text-medium-emphasis">
          {{
            confirmation.action === 'delete'
              ? $t('pages.adminUsers.confirm.deleteBody', {
                  email: confirmation.user.email,
                })
              : $t('pages.adminUsers.confirm.suspendBody', {
                  email: confirmation.user.email,
                })
          }}
        </VCardText>
        <VCardActions>
          <VBtn
            :color="confirmation.action === 'delete' ? 'error' : 'warning'"
            variant="flat"
            :loading="!!busy"
            @click="runAction(confirmation.user, confirmation.action)"
          >
            {{ $t(`pages.adminUsers.actions.${confirmation.action}`) }}
          </VBtn>
          <VSpacer />
          <VBtn variant="text" :disabled="!!busy" @click="confirmation = null">
            {{ $t('pages.adminUsers.confirm.cancel') }}
          </VBtn>
        </VCardActions>
      </VCard>
    </VDialog>
  </VContainer>
</template>

<style scoped>
.admin-header {
  display: flex;
  align-items: flex-start;
  gap: 1rem;
  margin-bottom: 1.5rem;
}

.admin-title {
  margin: 0 0 0.25rem;
  font-size: 1.9rem;
  font-weight: 800;
  line-height: 1.2;
  letter-spacing: -0.03em;
}

.admin-toolbar {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  justify-content: space-between;
  gap: 0.75rem 1rem;
}

.admin-users-search {
  flex: 0 1 320px;
  min-width: 240px;
}

.user-initial {
  font-weight: 800;
}

.user-text {
  min-width: 0;
}

.user-email {
  min-width: 12rem;
  overflow-wrap: anywhere;
}

.chip-dot {
  width: 0.45rem;
  height: 0.45rem;
  margin-right: 0.4rem;
  border-radius: 50%;
  background: currentcolor;
}

.user-list {
  margin: 1rem 0 0;
  padding: 0;
  list-style: none;
}

.user-item {
  padding: 1rem 0;
  border-top: 1px solid rgb(var(--v-border-color), var(--v-border-opacity));
}

.user-item:last-child {
  padding-bottom: 0;
}

.user-meta {
  display: flex;
  flex-direction: column;
  gap: 0.2rem;
  margin: 0.6rem 0 0;
  padding: 0;
  color: rgb(var(--v-theme-on-surface), var(--v-medium-emphasis-opacity));
  font-size: 0.85rem;
  list-style: none;
}

.user-meta li {
  display: flex;
  align-items: flex-start;
  gap: 0.4rem;
}

.user-meta .v-icon {
  margin-top: 0.1rem;
}

@media (width <= 600px) {
  .admin-title {
    font-size: 1.5rem;
  }

  .admin-users-search {
    flex-basis: 100%;
    min-width: 0;
  }
}
</style>
