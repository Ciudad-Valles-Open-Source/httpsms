import assert from 'node:assert/strict'
import { test } from 'node:test'
import {
  actionNeedsConfirmation,
  adminUserActions,
  formatCountBadge,
  isProtectedAdminUser,
  matchesAdminUserFilter,
} from './admin-users'

test('offers the actions that fit each account status', () => {
  assert.deepEqual(adminUserActions({ status: 'pending' }), [
    'approve',
    'suspend',
    'delete',
  ])
  assert.deepEqual(adminUserActions({ status: 'active' }), [
    'suspend',
    'delete',
  ])
  assert.deepEqual(adminUserActions({ status: 'suspended' }), [
    'reactivate',
    'delete',
  ])
})

test('only destructive actions ask for confirmation', () => {
  assert.equal(actionNeedsConfirmation('approve'), false)
  assert.equal(actionNeedsConfirmation('reactivate'), false)
  assert.equal(actionNeedsConfirmation('suspend'), true)
  assert.equal(actionNeedsConfirmation('delete'), true)
})

test('administrators and the signed-in account are protected', () => {
  assert.equal(isProtectedAdminUser({ id: 'a', is_admin: true }, 'me'), true)
  assert.equal(isProtectedAdminUser({ id: 'me', is_admin: false }, 'me'), true)
  assert.equal(isProtectedAdminUser({ id: 'a', is_admin: false }, 'me'), false)
  assert.equal(isProtectedAdminUser({ id: 'a', is_admin: false }, null), false)
})

test('matches accounts against the selected filter', () => {
  assert.equal(matchesAdminUserFilter({ status: 'active' }, 'all'), true)
  assert.equal(matchesAdminUserFilter({ status: 'active' }, 'active'), true)
  assert.equal(matchesAdminUserFilter({ status: 'active' }, 'pending'), false)
})

test('formats count badges with a cap at the page limit', () => {
  assert.equal(formatCountBadge(null, 100), null)
  assert.equal(formatCountBadge(0, 100), null)
  assert.equal(formatCountBadge(7, 100), '7')
  assert.equal(formatCountBadge(99, 100), '99')
  assert.equal(formatCountBadge(100, 100), '99+')
})
