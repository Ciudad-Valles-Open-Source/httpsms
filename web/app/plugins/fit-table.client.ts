import type { Directive } from 'vue'
import { observeTable, tableOverflows } from '~/utils/table-fit'

const STACKED = 'is-stacked'
const observers = new WeakMap<HTMLElement, () => void>()

// Measures the table laid out as a table and stacks it if it overflows. Both
// steps happen in the same frame, so the switch never flickers.
function fit(el: HTMLElement) {
  el.classList.remove(STACKED)
  if (tableOverflows(el)) el.classList.add(STACKED)
}

/**
 * `v-fit-table` on a `stack-table` VTable: turns its rows into cards whenever
 * the columns don't fit, and back into a table when there is room again.
 */
const fitTable: Directive<HTMLElement> = {
  mounted(el) {
    observers.set(
      el,
      observeTable(el, () => fit(el)),
    )
    fit(el)
  },
  // Vue may have re-rendered the class list.
  updated(el) {
    fit(el)
  },
  beforeUnmount(el) {
    observers.get(el)?.()
    observers.delete(el)
  },
}

export default defineNuxtPlugin((nuxtApp) => {
  nuxtApp.vueApp.directive('fit-table', fitTable)
})
