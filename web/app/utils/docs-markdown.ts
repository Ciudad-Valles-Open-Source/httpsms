// Markdown rendering for the built-in documentation (/docs).
//
// The content is written by us (web/docs-content), so raw HTML is still
// disabled (`html: false`) but no further sanitizing is done.

import markdownIt, { type MarkdownIt, type Token } from 'markdown-it'
import hljs from 'highlight.js/lib/core'
import bash from 'highlight.js/lib/languages/bash'
import csharp from 'highlight.js/lib/languages/csharp'
import go from 'highlight.js/lib/languages/go'
import http from 'highlight.js/lib/languages/http'
import java from 'highlight.js/lib/languages/java'
import javascript from 'highlight.js/lib/languages/javascript'
import json from 'highlight.js/lib/languages/json'
import php from 'highlight.js/lib/languages/php'
import plaintext from 'highlight.js/lib/languages/plaintext'
import python from 'highlight.js/lib/languages/python'
import shell from 'highlight.js/lib/languages/shell'
import typescript from 'highlight.js/lib/languages/typescript'
import xml from 'highlight.js/lib/languages/xml'
import yaml from 'highlight.js/lib/languages/yaml'

const HIGHLIGHT_LANGUAGES = {
  bash,
  csharp,
  go,
  http,
  java,
  javascript,
  json,
  php,
  plaintext,
  python,
  shell,
  typescript,
  xml,
  yaml,
}

for (const [name, language] of Object.entries(HIGHLIGHT_LANGUAGES)) {
  if (!hljs.getLanguage(name)) hljs.registerLanguage(name, language)
}

export interface DocPlaceholderOptions {
  apiBaseUrl: string
  appUrl: string
  appDownloadUrl: string
  locale: string
}

export interface RenderedDoc {
  html: string
  title: string
  description: string
}

const CALLOUT_PATTERNS: { variant: string; pattern: RegExp }[] = [
  {
    variant: 'warning',
    pattern:
      /^\*\*(advertencia|warning|importante|important|cuidado|caution|atenci[oó]n)\b/i,
  },
  {
    variant: 'tip',
    pattern: /^\*\*(consejo|tip|sugerencia|hint)\b/i,
  },
  {
    variant: 'info',
    pattern: /^\*\*(nota|note|info|información|information)\b/i,
  },
]

function trimTrailingSlash(url: string): string {
  return url.replace(/\/+$/, '')
}

/** Values for the `{{PLACEHOLDER}}` tokens that may appear in the Markdown. */
export function buildDocPlaceholders(
  options: DocPlaceholderOptions,
): Record<string, string> {
  const apiBaseUrl = trimTrailingSlash(options.apiBaseUrl)
  return {
    API_BASE_URL: apiBaseUrl,
    APP_URL: trimTrailingSlash(options.appUrl),
    APP_DOWNLOAD_URL: options.appDownloadUrl,
    API_DOCS_URL:
      options.locale === 'es'
        ? `${apiBaseUrl}/es/index.html`
        : `${apiBaseUrl}/index.html`,
  }
}

/** Replaces `{{NAME}}` with `values[NAME]`; unknown placeholders are kept. */
export function substituteDocPlaceholders(
  markdown: string,
  values: Record<string, string>,
): string {
  return markdown.replace(/\{\{\s*([A-Z0-9_]+)\s*\}\}/g, (match, name) =>
    Object.hasOwn(values, name) ? (values[name] as string) : match,
  )
}

/** GitHub-like heading slug which also strips accents (`Introducción` → `introduccion`). */
export function slugifyDocHeading(text: string): string {
  return text
    .normalize('NFD')
    .replace(/[̀-ͯ]/g, '')
    .toLowerCase()
    .trim()
    .replace(/[^a-z0-9\s-]/g, '')
    .replace(/[\s-]+/g, '-')
    .replace(/^-+|-+$/g, '')
}

function inlineText(token: Token | undefined): string {
  if (!token) return ''
  if (!token.children) return token.content
  return token.children
    .map((child) => {
      if (child.type === 'text' || child.type === 'code_inline') {
        return child.content
      }
      return child.type === 'softbreak' || child.type === 'hardbreak' ? ' ' : ''
    })
    .join('')
}

function truncate(text: string, max: number): string {
  if (text.length <= max) return text
  const cut = text.slice(0, max - 1)
  const lastSpace = cut.lastIndexOf(' ')
  return `${(lastSpace > max / 2 ? cut.slice(0, lastSpace) : cut).trimEnd()}…`
}

function createRenderer(): MarkdownIt {
  const md: MarkdownIt = markdownIt({
    html: false,
    linkify: true,
    typographer: false,
    highlight(code: string, lang: string): string {
      const language = lang && hljs.getLanguage(lang) ? lang : null
      const body = language
        ? hljs.highlight(code, { language, ignoreIllegals: true }).value
        : md.utils.escapeHtml(code)
      const label = lang ? ` data-lang="${md.utils.escapeHtml(lang)}"` : ''
      const className = language ? ` language-${language}` : ''
      return `<pre class="docs-code"${label}><code class="hljs${className}">${body}</code></pre>`
    },
  })

  // Heading anchors (h2+ get a "#" link) and callout variants for blockquotes.
  md.core.ruler.push('docs_decorations', (state) => {
    const used = new Map<string, number>()
    const tokens = state.tokens
    tokens.forEach((token, index) => {
      if (token.type === 'heading_open') {
        const base =
          slugifyDocHeading(inlineText(tokens[index + 1])) || 'section'
        const count = used.get(base) ?? 0
        used.set(base, count + 1)
        const id = count === 0 ? base : `${base}-${count}`
        token.attrSet('id', id)
        if (token.tag !== 'h1') {
          const close = tokens
            .slice(index)
            .find((t) => t.type === 'heading_close')
          if (close) close.meta = { anchorId: id }
        }
      }

      if (token.type === 'blockquote_open') {
        const firstInline = tokens
          .slice(index + 1)
          .find((t) => t.type === 'inline')
        const content = firstInline?.content.trim() ?? ''
        const callout = CALLOUT_PATTERNS.find(({ pattern }) =>
          pattern.test(content),
        )
        token.attrJoin(
          'class',
          `docs-callout docs-callout-${callout?.variant ?? 'default'}`,
        )
      }
    })
  })

  md.renderer.rules.heading_close = (tokens, index, options, env, self) => {
    const token = tokens[index]!
    const anchorId = token.meta?.anchorId as string | undefined
    const anchor = anchorId
      ? ` <a class="docs-heading-anchor" href="#${anchorId}" aria-hidden="true" tabindex="-1">#</a>`
      : ''
    return anchor + self.renderToken(tokens, index, options)
  }

  // External links open in a new tab.
  md.renderer.rules.link_open = (tokens, index, options, env, self) => {
    const token = tokens[index]!
    const href = token.attrGet('href') ?? ''
    if (/^https?:\/\//i.test(String(href))) {
      token.attrSet('target', '_blank')
      token.attrSet('rel', 'noopener noreferrer')
    }
    return self.renderToken(tokens, index, options)
  }

  // Wide tables scroll horizontally instead of breaking the layout.
  md.renderer.rules.table_open = (tokens, index, options, env, self) =>
    `<div class="docs-table">${self.renderToken(tokens, index, options)}`
  md.renderer.rules.table_close = (tokens, index, options, env, self) =>
    `${self.renderToken(tokens, index, options)}</div>`

  return md
}

let renderer: MarkdownIt | null = null

/** Renders a docs page and extracts its title (first `# h1`) and description (first paragraph). */
export function renderDocMarkdown(markdown: string): RenderedDoc {
  renderer ??= createRenderer()
  const env = {}
  const tokens = renderer.parse(markdown, env)

  let title = ''
  let description = ''
  tokens.forEach((token, index) => {
    if (!title && token.type === 'heading_open' && token.tag === 'h1') {
      title = inlineText(tokens[index + 1]).trim()
    }
    if (!description && token.type === 'paragraph_open' && token.level === 0) {
      description = truncate(
        inlineText(tokens[index + 1])
          .replace(/\s+/g, ' ')
          .trim(),
        160,
      )
    }
  })

  return {
    html: renderer.renderer.render(tokens, renderer.options, env),
    title,
    description,
  }
}
