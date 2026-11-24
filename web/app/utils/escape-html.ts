/**
 * Escapes text so it can be interpolated safely into a string rendered with
 * v-html (e.g. a translated message that contains markup).
 */
export function escapeHtml(text: string | null | undefined): string {
  return String(text ?? '')
    .replace(/&/g, '&amp;')
    .replace(/</g, '&lt;')
    .replace(/>/g, '&gt;')
    .replace(/"/g, '&quot;')
    .replace(/'/g, '&#39;')
}
