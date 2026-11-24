<script setup lang="ts">
import { mdiArrowLeft, mdiMagnify, mdiRefresh, mdiAccountSearch } from '@mdi/js'
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

const { mdAndUp } = useVDisplay()
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
      <VContainer>
        <VRow>
          <VCol cols="12" lg="10" offset-lg="1" xl="8" offset-xl="2">
            <h4 class="text-headline-large mb-2 mt-0">
              {{ $t('pages.adminUsers.heading') }}
            </h4>
            <p class="text-medium-emphasis mb-4">
              {{ $t('pages.adminUsers.description') }}
            </p>

            <div class="d-flex flex-column flex-md-row ga-3 mb-4">
              <VChipGroup
                :model-value="adminStore.filter"
                mandatory
                selected-class="text-primary"
                class="flex-grow-1"
                @update:model-value="selectFilter"
              >
                <VChip
                  v-for="item in ADMIN_USER_FILTERS"
                  :key="item"
                  :value="item"
                  variant="tonal"
                  filter
                >
                  {{ $t(`pages.adminUsers.filters.${item}`) }}
                  <VBadge
                    v-if="item === 'pending' && pendingBadge"
                    :content="pendingBadge"
                    color="warning"
                    inline
                  />
                </VChip>
              </VChipGroup>
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
              class="mb-4"
            >
              {{ loadErrorMessage }}
              <template #append>
                <VBtn variant="text" color="error" @click="reload">
                  {{ $t('pages.adminUsers.retry') }}
                </VBtn>
              </template>
            </VAlert>

            <template
              v-else-if="adminStore.loading && !adminStore.users.length"
            >
              <VSkeletonLoader
                v-for="n in 3"
                :key="n"
                type="list-item-two-line"
                class="mb-2"
              />
            </template>

            <div
              v-else-if="adminStore.loaded && !adminStore.users.length"
              class="text-center py-12"
            >
              <VIcon
                :icon="mdiAccountSearch"
                size="56"
                class="text-medium-emphasis mb-3"
              />
              <p class="text-body-large text-medium-emphasis">
                {{ emptyMessage }}
              </p>
            </div>

            <template v-else-if="adminStore.users.length">
              <VTable v-if="mdAndUp" density="comfortable">
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
                    <td class="text-break">
                      {{ user.email }}
                      <VChip
                        v-if="user.is_admin"
                        size="x-small"
                        color="primary"
                        label
                        class="ml-1"
                      >
                        {{ $t('pages.adminUsers.adminBadge') }}
                      </VChip>
                      <VChip
                        v-else-if="user.id === authStore.user?.id"
                        size="x-small"
                        label
                        class="ml-1"
                      >
                        {{ $t('pages.adminUsers.youBadge') }}
                      </VChip>
                    </td>
                    <td>
                      <VChip
                        size="small"
                        label
                        :color="STATUS_COLOR[user.status]"
                      >
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
                    <td class="text-center">{{ user.phones_count ?? 0 }}</td>
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
                        :user="user"
                        :disabled="!!busy && busy.id !== user.id"
                        :busy-action="busyAction(user)"
                        @action="onAction(user, $event)"
                      />
                    </td>
                  </tr>
                </tbody>
              </VTable>

              <template v-else>
                <VCard
                  v-for="user in adminStore.users"
                  :key="user.id"
                  variant="tonal"
                  class="mb-3"
                >
                  <VCardText>
                    <div class="d-flex align-start ga-2 mb-2">
                      <span class="text-body-large text-break flex-grow-1">
                        {{ user.email }}
                      </span>
                      <VChip
                        size="small"
                        label
                        :color="STATUS_COLOR[user.status]"
                      >
                        {{ $t(`pages.adminUsers.status.${user.status}`) }}
                      </VChip>
                    </div>
                    <div class="text-body-small text-medium-emphasis">
                      <div>
                        {{
                          $t('pages.adminUsers.createdOn', {
                            date: formatDate(user.created_at),
                          })
                        }}
                      </div>
                      <div v-if="user.approved_at">
                        {{
                          $t('pages.adminUsers.approvedOn', {
                            date: formatDate(user.approved_at),
                          })
                        }}
                        <span v-if="user.approved_by">
                          ·
                          {{
                            $t('pages.adminUsers.approvedBy', {
                              name: user.approved_by,
                            })
                          }}
                        </span>
                      </div>
                      <div>
                        {{
                          $t('pages.adminUsers.phonesCount', {
                            count: user.phones_count ?? 0,
                          })
                        }}
                      </div>
                    </div>
                    <div class="mt-3">
                      <VChip
                        v-if="user.is_admin"
                        size="small"
                        color="primary"
                        label
                      >
                        {{ $t('pages.adminUsers.adminBadge') }}
                      </VChip>
                      <VChip
                        v-else-if="user.id === authStore.user?.id"
                        size="small"
                        label
                      >
                        {{ $t('pages.adminUsers.youBadge') }}
                      </VChip>
                      <AdminUserActions
                        v-else
                        :user="user"
                        :disabled="!!busy && busy.id !== user.id"
                        :busy-action="busyAction(user)"
                        @action="onAction(user, $event)"
                      />
                    </div>
                  </VCardText>
                </VCard>
              </template>

              <div v-if="adminStore.hasMore" class="text-center mt-4">
                <VBtn
                  variant="tonal"
                  :loading="adminStore.loadingMore"
                  @click="loadMore"
                >
                  {{ $t('pages.adminUsers.loadMore') }}
                </VBtn>
              </div>
            </template>
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
.admin-users-search {
  min-width: 260px;
  max-width: 100%;
}
</style>
