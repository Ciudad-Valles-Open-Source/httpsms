import { defineStore } from 'pinia'
import { ref } from 'vue'
import {
  ACTION_TARGET_STATUS,
  matchesAdminUserFilter,
  type AdminUser,
  type AdminUserAction,
  type AdminUserFilter,
} from '../utils/admin-users'

export const ADMIN_USERS_PAGE_SIZE = 20

/** Pending accounts are counted with one request of at most this size. */
export const PENDING_COUNT_LIMIT = 100

export const useAdminStore = defineStore('admin', () => {
  const { apiFetch } = useApi()

  const users = ref<AdminUser[]>([])
  const filter = ref<AdminUserFilter>('pending')
  const query = ref('')
  const loading = ref(false)
  const loadingMore = ref(false)
  const loaded = ref(false)
  const hasMore = ref(false)
  const error = ref<unknown>(null)
  const pendingCount = ref<number | null>(null)
  let listVersion = 0

  function listParams(skip: number) {
    const params: Record<string, string | number> = {
      skip,
      limit: ADMIN_USERS_PAGE_SIZE,
    }
    if (filter.value !== 'all') params.status = filter.value
    const search = query.value.trim()
    if (search) params.query = search
    return params
  }

  async function loadUsers() {
    const version = ++listVersion
    loading.value = true
    loadingMore.value = false
    error.value = null
    try {
      const response = await apiFetch<{ data: AdminUser[] }>(
        '/v1/admin/users',
        { params: listParams(0) },
      )
      if (version !== listVersion) return
      users.value = response.data ?? []
      hasMore.value = users.value.length === ADMIN_USERS_PAGE_SIZE
      loaded.value = true
    } catch (e: unknown) {
      if (version !== listVersion) return
      users.value = []
      hasMore.value = false
      error.value = e
    } finally {
      if (version === listVersion) loading.value = false
    }
  }

  async function loadMore() {
    if (loading.value || loadingMore.value || !hasMore.value) return
    const version = listVersion
    loadingMore.value = true
    try {
      const response = await apiFetch<{ data: AdminUser[] }>(
        '/v1/admin/users',
        { params: listParams(users.value.length) },
      )
      if (version !== listVersion) return
      const page = response.data ?? []
      const known = new Set(users.value.map((user) => user.id))
      users.value = [
        ...users.value,
        ...page.filter((user) => !known.has(user.id)),
      ]
      hasMore.value = page.length === ADMIN_USERS_PAGE_SIZE
    } finally {
      if (version === listVersion) loadingMore.value = false
    }
  }

  async function setFilter(value: AdminUserFilter) {
    filter.value = value
    await loadUsers()
  }

  async function setQuery(value: string) {
    query.value = value
    await loadUsers()
  }

  async function loadPendingCount() {
    const response = await apiFetch<{ data: AdminUser[] }>('/v1/admin/users', {
      params: { status: 'pending', skip: 0, limit: PENDING_COUNT_LIMIT },
    })
    pendingCount.value = (response.data ?? []).length
  }

  function refreshPendingCount() {
    loadPendingCount().catch(console.error)
  }

  async function changeStatus(
    user: AdminUser,
    action: Exclude<AdminUserAction, 'delete'>,
  ): Promise<AdminUser> {
    const response = await apiFetch<{ data: AdminUser }>(
      `/v1/admin/users/${encodeURIComponent(user.id)}/status`,
      { method: 'PUT', body: { status: ACTION_TARGET_STATUS[action] } },
    )
    const updated = response.data
    users.value = matchesAdminUserFilter(updated, filter.value)
      ? users.value.map((item) => (item.id === updated.id ? updated : item))
      : users.value.filter((item) => item.id !== updated.id)
    refreshPendingCount()
    return updated
  }

  async function deleteUser(user: AdminUser) {
    await apiFetch(`/v1/admin/users/${encodeURIComponent(user.id)}`, {
      method: 'DELETE',
    })
    users.value = users.value.filter((item) => item.id !== user.id)
    refreshPendingCount()
  }

  function resetState() {
    listVersion++
    users.value = []
    filter.value = 'pending'
    query.value = ''
    loading.value = false
    loadingMore.value = false
    loaded.value = false
    hasMore.value = false
    error.value = null
    pendingCount.value = null
  }

  return {
    users,
    filter,
    query,
    loading,
    loadingMore,
    loaded,
    hasMore,
    error,
    pendingCount,
    loadUsers,
    loadMore,
    setFilter,
    setQuery,
    loadPendingCount,
    changeStatus,
    deleteUser,
    resetState,
  }
})
