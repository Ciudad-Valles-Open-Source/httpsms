/**
 * Pure helpers for the administrator's user management panel.
 */

export type AdminUserStatus = 'active' | 'pending' | 'suspended'

export type AdminUserFilter = AdminUserStatus | 'all'

export type AdminUserAction = 'approve' | 'suspend' | 'reactivate' | 'delete'

/** An account as returned by `GET /v1/admin/users`. */
export interface AdminUser {
  id: string
  email: string
  status: AdminUserStatus
  is_admin: boolean
  created_at: string
  updated_at: string
  approved_at?: string | null
  approved_by?: string | null
  phones_count: number
}

export const ADMIN_USER_FILTERS: AdminUserFilter[] = [
  'pending',
  'active',
  'suspended',
  'all',
]

/** The status each action sets through `PUT /v1/admin/users/{id}/status`. */
export const ACTION_TARGET_STATUS: Record<
  Exclude<AdminUserAction, 'delete'>,
  AdminUserStatus
> = {
  approve: 'active',
  suspend: 'suspended',
  reactivate: 'active',
}

/** Actions offered for an account, in display order. */
export function adminUserActions(
  user: Pick<AdminUser, 'status'>,
): AdminUserAction[] {
  switch (user.status) {
    case 'pending':
      return ['approve', 'suspend', 'delete']
    case 'active':
      return ['suspend', 'delete']
    case 'suspended':
      return ['reactivate', 'delete']
    default:
      return ['delete']
  }
}

/** Destructive actions ask for confirmation first. */
export function actionNeedsConfirmation(action: AdminUserAction): boolean {
  return action === 'suspend' || action === 'delete'
}

/**
 * Administrators and the signed-in account cannot be changed from the panel
 * (the API answers 422 as well).
 */
export function isProtectedAdminUser(
  user: Pick<AdminUser, 'id' | 'is_admin'>,
  currentUserId: string | null | undefined,
): boolean {
  return user.is_admin || (!!currentUserId && user.id === currentUserId)
}

/** Whether an account belongs in the list shown for a filter. */
export function matchesAdminUserFilter(
  user: Pick<AdminUser, 'status'>,
  filter: AdminUserFilter,
): boolean {
  return filter === 'all' || user.status === filter
}

/**
 * Text for a count badge when the count was fetched with a page `limit`:
 * hidden when empty and capped (e.g. "99+") when the page was full.
 */
export function formatCountBadge(
  count: number | null | undefined,
  limit: number,
): string | null {
  if (!count || count <= 0) return null
  if (count >= limit) return `${limit - 1}+`
  return String(count)
}
