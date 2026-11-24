import assert from 'node:assert/strict'
import { test } from 'node:test'
import { createPinia, setActivePinia } from 'pinia'
import type { AdminUser } from '../utils/admin-users'
import { ADMIN_USERS_PAGE_SIZE, useAdminStore } from './admin'

interface Call {
  path: string
  method?: string
  params?: Record<string, unknown>
  body?: Record<string, unknown>
}

function account(id: string, status: AdminUser['status']): AdminUser {
  return {
    id,
    email: `${id}@example.test`,
    status,
    is_admin: false,
    created_at: '2026-11-20T10:00:00Z',
    updated_at: '2026-11-20T10:00:00Z',
    approved_at: null,
    approved_by: null,
    phones_count: 0,
  }
}

function setupStore(accounts: AdminUser[]) {
  const calls: Call[] = []
  let data = [...accounts]
  const deferred: Promise<unknown>[] = []
  Object.assign(globalThis, {
    useApi: () => ({
      apiFetch: async (
        path: string,
        options: Omit<Call, 'path'> = {},
      ): Promise<unknown> => {
        calls.push({ path, ...options })
        const pending = deferred.shift()
        if (pending) return pending
        if (options.method === 'PUT') {
          const id = path.split('/').at(-2)!
          const status = options.body!.status as AdminUser['status']
          data = data.map((item) =>
            item.id === id ? { ...item, status } : item,
          )
          return { data: data.find((item) => item.id === id) }
        }
        if (options.method === 'DELETE') {
          data = data.filter((item) => item.id !== path.split('/').at(-1))
          return { message: 'deleted' }
        }
        const params = options.params ?? {}
        const filtered = data.filter(
          (item) =>
            (!params.status || item.status === params.status) &&
            (!params.query || item.email.includes(String(params.query))),
        )
        const skip = Number(params.skip ?? 0)
        return {
          data: filtered.slice(skip, skip + Number(params.limit ?? 20)),
        }
      },
    }),
  })
  setActivePinia(createPinia())
  return {
    store: useAdminStore(),
    calls,
    defer: (response: Promise<unknown>) => deferred.push(response),
  }
}

test('loads pending accounts first and pages through them', async () => {
  const accounts = Array.from({ length: ADMIN_USERS_PAGE_SIZE + 5 }, (_, i) =>
    account(`user-${i}`, 'pending'),
  )
  const { store, calls } = setupStore(accounts)

  await store.loadUsers()
  assert.equal(calls[0]?.path, '/v1/admin/users')
  assert.deepEqual(calls[0]?.params, {
    skip: 0,
    limit: ADMIN_USERS_PAGE_SIZE,
    status: 'pending',
  })
  assert.equal(store.users.length, ADMIN_USERS_PAGE_SIZE)
  assert.equal(store.hasMore, true)

  await store.loadMore()
  assert.equal(calls[1]?.params?.skip, ADMIN_USERS_PAGE_SIZE)
  assert.equal(store.users.length, ADMIN_USERS_PAGE_SIZE + 5)
  assert.equal(store.hasMore, false)
})

test('the all filter omits the status and search sends a trimmed query', async () => {
  const { store, calls } = setupStore([
    account('ana', 'active'),
    account('bob', 'pending'),
  ])
  store.filter = 'all'
  await store.setQuery('  ana ')
  assert.deepEqual(calls.at(-1)?.params, {
    skip: 0,
    limit: ADMIN_USERS_PAGE_SIZE,
    query: 'ana',
  })
  assert.deepEqual(
    store.users.map((user) => user.id),
    ['ana'],
  )
})

test('a stale response does not overwrite a newer filter', async () => {
  const { store, defer } = setupStore([])
  let resolveStale!: (value: unknown) => void
  defer(new Promise((resolve) => (resolveStale = resolve)))
  const stale = store.loadUsers()
  const fresh = store.setFilter('active')
  await fresh
  resolveStale({ data: [account('stale', 'pending')] })
  await stale
  assert.deepEqual(store.users, [])
  assert.equal(store.loading, false)
})

test('approving removes the account from the pending list', async () => {
  const { store, calls } = setupStore([
    account('ana', 'pending'),
    account('bob', 'pending'),
  ])
  await store.loadUsers()
  const updated = await store.changeStatus(store.users[0]!, 'approve')

  const put = calls.find((call) => call.method === 'PUT')
  assert.equal(put?.path, '/v1/admin/users/ana/status')
  assert.deepEqual(put?.body, { status: 'active' })
  assert.equal(updated.status, 'active')
  assert.deepEqual(
    store.users.map((user) => user.id),
    ['bob'],
  )
})

test('changing status under the all filter updates the row in place', async () => {
  const { store } = setupStore([account('ana', 'active')])
  await store.setFilter('all')
  await store.changeStatus(store.users[0]!, 'suspend')
  assert.equal(store.users[0]?.status, 'suspended')
})

test('deleting removes the account and counts pending accounts', async () => {
  const { store, calls } = setupStore([
    account('ana', 'pending'),
    account('bob', 'pending'),
  ])
  await store.loadUsers()
  await store.deleteUser(store.users[0]!)
  assert.equal(
    calls.find((c) => c.method === 'DELETE')?.path,
    '/v1/admin/users/ana',
  )
  assert.deepEqual(
    store.users.map((user) => user.id),
    ['bob'],
  )

  await store.loadPendingCount()
  assert.equal(store.pendingCount, 1)
})

test('a failed load keeps the error for the page', async () => {
  const { store, defer } = setupStore([])
  defer(Promise.reject(new Error('forbidden')))
  await store.loadUsers()
  assert.ok(store.error instanceof Error)
  assert.equal(store.loading, false)
  assert.deepEqual(store.users, [])
})
