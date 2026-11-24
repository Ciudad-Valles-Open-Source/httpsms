/**
 * Helpers for the account approval flow: new accounts start "pending" until an
 * administrator approves them, and administrators can suspend accounts.
 */

export type AccountStatus = 'active' | 'pending' | 'suspended'

export type AccountBlockedCode = 'account_pending' | 'account_suspended'

export const ACCOUNT_PENDING_PATH = '/account-pending'

/** Where the dashboard sends signed-in users once their account is usable. */
export const DASHBOARD_PATH = '/threads'

export interface AccountStatusUser {
  status?: string | null
  is_admin?: boolean | null
}

/**
 * Whether the account cannot use the dashboard yet. Administrators are never
 * blocked, and a missing status (older API) is treated as active so nobody is
 * locked out by a version mismatch.
 */
export function isAccountBlocked(
  user: AccountStatusUser | null | undefined,
): boolean {
  if (!user || user.is_admin) return false
  return !!user.status && user.status !== 'active'
}

/**
 * Decides where an authenticated navigation to `path` must go instead, or
 * returns null to let it through. Blocked accounts may only see the account
 * status page, and usable accounts are sent from that page to the dashboard.
 */
export function accountStatusRedirect(
  user: AccountStatusUser | null | undefined,
  path: string,
): string | null {
  if (!user) return null
  const onStatusPage = path.replace(/\/+$/, '') === ACCOUNT_PENDING_PATH
  if (isAccountBlocked(user)) {
    return onStatusPage ? null : ACCOUNT_PENDING_PATH
  }
  return onStatusPage ? DASHBOARD_PATH : null
}

function isBlockedCode(value: unknown): value is AccountBlockedCode {
  return value === 'account_pending' || value === 'account_suspended'
}

/**
 * Extracts the account status code from a 403 API response body such as
 * `{"status":"error","data":{"code":"account_pending"}}`. Accepts either the
 * response body or an ofetch error (whose `data` is the body).
 */
export function getAccountBlockedCode(
  body: unknown,
): AccountBlockedCode | null {
  if (body === null || typeof body !== 'object') return null
  const record = body as { data?: unknown; code?: unknown }
  if (isBlockedCode(record.code)) return record.code
  if (record.data !== null && typeof record.data === 'object') {
    const nested = record.data as { code?: unknown; data?: unknown }
    if (isBlockedCode(nested.code)) return nested.code
    if (nested.data !== null && typeof nested.data === 'object') {
      const code = (nested.data as { code?: unknown }).code
      if (isBlockedCode(code)) return code
    }
  }
  return null
}

/** The status implied by an account status error code. */
export function statusFromBlockedCode(code: AccountBlockedCode): AccountStatus {
  return code === 'account_suspended' ? 'suspended' : 'pending'
}
