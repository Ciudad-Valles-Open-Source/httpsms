import { observeTable, tableOverflows } from '~/utils/table-fit'

/**
 * For tables rendered by Vuetify components (VDataTable's `mobile` prop) or
 * swapped for a custom list: `stacked` turns true when the table, laid out as
 * a table, doesn't fit inside `target`, and false again once `target` is as
 * wide as the table needed.
 */
export function useTableFit() {
  // A wrapper element or a component (e.g. VDataTable) via its template ref.
  const target = ref<HTMLElement | { $el: HTMLElement } | null>(null)
  const element = () =>
    target.value instanceof HTMLElement
      ? target.value
      : (target.value?.$el ?? null)
  const stacked = ref(false)
  let requiredWidth = 0
  let stopObserving: (() => void) | null = null

  async function check() {
    const root = element()
    if (!root) return
    if (stacked.value) {
      if (root.clientWidth < requiredWidth) return
      // Lay it out as a table again and measure (before the next paint).
      stacked.value = false
      await nextTick()
    }
    if (tableOverflows(root)) {
      const scroller =
        root.querySelector<HTMLElement>('.v-table__wrapper') ?? root
      requiredWidth = scroller.scrollWidth
      stacked.value = true
    }
  }

  // The element can appear later (e.g. after the data loads), so (re)attach
  // whenever the template ref changes.
  watch(
    element,
    (root) => {
      stopObserving?.()
      stopObserving = null
      if (!root) return
      stopObserving = observeTable(root, check)
      check()
    },
    { flush: 'post' },
  )
  onBeforeUnmount(() => stopObserving?.())

  return { target, stacked }
}
