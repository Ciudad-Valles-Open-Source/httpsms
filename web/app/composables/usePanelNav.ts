import {
  mdiAccountCog,
  mdiAccountMultiple,
  mdiCellphoneKey,
  mdiChartBar,
  mdiCommentTextMultipleOutline,
  mdiMagnify,
  mdiMessageTextOutline,
  mdiPlus,
} from '@mdi/js'

export interface PanelNavItem {
  /** Route name, also used to decide which item is active. */
  name: string
  icon: string
  /** i18n key of the label. */
  label: string
  /** Only shown once the user has picked one of their phones. */
  requiresOwner?: boolean
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
]

/** Navigation entries of the dashboard, shared by the rail and the menus. */
export function usePanelNav() {
  const route = useRoute()
  const phonesStore = usePhonesStore()

  const items = computed(() =>
    PANEL_NAV_ITEMS.filter((item) => !item.requiresOwner || phonesStore.owner),
  )

  function isActive(item: PanelNavItem): boolean {
    const current = (route.name as string) ?? ''
    return (item.matches ?? [item.name]).includes(current)
  }

  return { items, isActive }
}
