// Detects tables whose columns don't fit the available width, so they can
// switch to a list of cards (see `stack-table` in main.css, the v-fit-table
// directive and useTableFit).

/** True when the table's content is wider than the space it has. */
export function tableOverflows(root: HTMLElement): boolean {
  const scroller = root.querySelector<HTMLElement>('.v-table__wrapper') ?? root
  return scroller.scrollWidth > scroller.clientWidth + 1
}

/**
 * Calls `onChange` (at most once per frame, before paint) when the element's
 * width changes or its content changes (rows added, text updated). Height
 * changes alone, like the table switching to cards, are ignored.
 */
export function observeTable(
  element: HTMLElement,
  onChange: () => void,
): () => void {
  let frame = 0
  const schedule = () => {
    if (frame) return
    frame = requestAnimationFrame(() => {
      frame = 0
      onChange()
    })
  }

  let lastWidth = -1
  const resize = new ResizeObserver(() => {
    const width = element.clientWidth
    if (width === lastWidth) return
    lastWidth = width
    schedule()
  })
  resize.observe(element)

  // Class changes are not observed: the stacked state itself is a class.
  const mutation = new MutationObserver(schedule)
  mutation.observe(element, {
    childList: true,
    subtree: true,
    characterData: true,
  })

  return () => {
    cancelAnimationFrame(frame)
    resize.disconnect()
    mutation.disconnect()
  }
}
