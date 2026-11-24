import {
  mdiAccountCog,
  mdiAccountMultiple,
  mdiCellphoneKey,
  mdiChartBar,
  mdiCommentTextMultipleOutline,
  mdiMagnify,
  mdiMessageTextOutline,
  mdiPlus,
  mdiShieldAccount,
} from '@mdi/js'
import { formatCountBadge } from '~/utils/admin-users'
import { PENDING_COUNT_LIMIT } from '~/stores/admin'

export interface PanelNavItem {
  /** Route name, also used to decide which item is active. */
  name: string
  icon: string
  /** i18n key of the label. */
  label: string
  /** Only shown once the user has picked one of their phones. */
  requiresOwner?: boolean
  /** Only shown to administrators. */
  requiresAdmin?: boolean
  /** Shows the number of accounts waiting for approval. */
  showsPendingUsers?: boolean
  /** Other route names that belong to this section. */
  matches?: string[]
}

export const PANEL_NAV_ITEMS: PanelNavItem[] = [
  {
    name: 'threads',
    icon: mdiMessageTextOutline,
    label: 'components.panelNav.threads',
    matches: ['threads', 'threads-id', 'heartbeats-id'],
  },
  {
    name: 'messages',
    icon: mdiPlus,
    label: 'components.messageThreadHeader.newMessage',
    requiresOwner: true,
  },
  {
    name: 'bulk-messages',
    icon: mdiCommentTextMultipleOutline,
    label: 'components.messageThreadHeader.bulkMessages',
    requiresOwner: true,
  },
  {
    name: 'search-messages',
    icon: mdiMagnify,
    label: 'components.messageThreadHeader.searchMessages',
    requiresOwner: true,
  },
  {
    name: 'contacts',
    icon: mdiAccountMultiple,
    label: 'components.messageThreadHeader.contacts',
  },
  {
    name: 'statistics',
    icon: mdiChartBar,
    label: 'components.messageThreadHeader.statistics',
  },
  {
    name: 'phone-api-keys',
    icon: mdiCellphoneKey,
    label: 'components.messageThreadHeader.phoneApiKeys',
  },
  {
    name: 'settings',
    icon: mdiAccountCog,
    label: 'components.messageThreadHeader.settings',
  },
  {
    name: 'admin-users',
    icon: mdiShieldAccount,
    label: 'components.messageThreadHeader.administration',
    requiresAdmin: true,
    showsPendingUsers: true,
  },
]

/** Navigation entries of the dashboard, shared by the rail and the menus. */
export function usePanelNav() {
  const route = useRoute()
  const phonesStore = usePhonesStore()
  const authStore = useAuthStore()
  const adminStore = useAdminStore()

  const items = computed(() =>
    PANEL_NAV_ITEMS.filter(
      (item) =>
        (!item.requiresOwner || phonesStore.owner) &&
        (!item.requiresAdmin || authStore.isAdmin),
    ),
  )

  // Administrators see how many accounts are waiting for approval.
  const pendingUsersBadge = computed(() =>
    formatCountBadge(adminStore.pendingCount, PENDING_COUNT_LIMIT),
  )

  function badgeFor(item: PanelNavItem): string | null {
    return item.showsPendingUsers ? pendingUsersBadge.value : null
  }

  function isActive(item: PanelNavItem): boolean {
    const current = (route.name as string) ?? ''
    return (item.matches ?? [item.name]).includes(current)
  }

  return { items, isActive, badgeFor }
}
