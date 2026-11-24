import assert from 'node:assert/strict'
import { test } from 'node:test'
import {
  ACCOUNT_PENDING_PATH,
  DASHBOARD_PATH,
  accountStatusRedirect,
  getAccountBlockedCode,
  isAccountBlocked,
  statusFromBlockedCode,
} from './account-status'

test('pending and suspended accounts are blocked unless they are admins', () => {
  assert.equal(isAccountBlocked({ status: 'pending', is_admin: false }), true)
  assert.equal(isAccountBlocked({ status: 'suspended' }), true)
  assert.equal(isAccountBlocked({ status: 'active', is_admin: false }), false)
  assert.equal(isAccountBlocked({ status: 'pending', is_admin: true }), false)
})

test('a missing user or status never blocks the account', () => {
  assert.equal(isAccountBlocked(null), false)
  assert.equal(isAccountBlocked(undefined), false)
  assert.equal(isAccountBlocked({}), false)
})

test('blocked accounts are sent to the account status page', () => {
  const user = { status: 'pending', is_admin: false }
  assert.equal(accountStatusRedirect(user, '/threads'), ACCOUNT_PENDING_PATH)
  assert.equal(accountStatusRedirect(user, '/settings'), ACCOUNT_PENDING_PATH)
  assert.equal(
    accountStatusRedirect(user, '/admin/users'),
    ACCOUNT_PENDING_PATH,
  )
  assert.equal(accountStatusRedirect(user, ACCOUNT_PENDING_PATH), null)
  assert.equal(accountStatusRedirect(user, `${ACCOUNT_PENDING_PATH}/`), null)
})

test('usable accounts leave the account status page for the dashboard', () => {
  const user = { status: 'active', is_admin: false }
  assert.equal(accountStatusRedirect(user, '/threads'), null)
  assert.equal(
    accountStatusRedirect(user, ACCOUNT_PENDING_PATH),
    DASHBOARD_PATH,
  )
  assert.equal(
    accountStatusRedirect({ status: 'pending', is_admin: true }, '/settings'),
    null,
  )
})

test('an unknown user lets the navigation through', () => {
  assert.equal(accountStatusRedirect(null, '/threads'), null)
  assert.equal(accountStatusRedirect(null, ACCOUNT_PENDING_PATH), null)
})

test('reads the account status code from API response bodies and errors', () => {
  const body = {
    status: 'error',
    message: 'pending',
    data: { code: 'account_pending' },
  }
  assert.equal(getAccountBlockedCode(body), 'account_pending')
  assert.equal(
    getAccountBlockedCode({ data: { data: { code: 'account_suspended' } } }),
    'account_suspended',
  )
  assert.equal(
    getAccountBlockedCode({ code: 'account_suspended' }),
    'account_suspended',
  )
})

test('ignores unrelated or malformed error bodies', () => {
  assert.equal(getAccountBlockedCode(null), null)
  assert.equal(getAccountBlockedCode('account_pending'), null)
  assert.equal(getAccountBlockedCode({ data: { code: 'forbidden' } }), null)
  assert.equal(
    getAccountBlockedCode({ data: { email: ['is required'] } }),
    null,
  )
})

test('maps error codes to account statuses', () => {
  assert.equal(statusFromBlockedCode('account_pending'), 'pending')
  assert.equal(statusFromBlockedCode('account_suspended'), 'suspended')
})
