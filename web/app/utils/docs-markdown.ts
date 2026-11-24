// Markdown rendering for the built-in documentation (/docs).
//
// The content is written by us (web/docs-content), so raw HTML is still
// disabled (`html: false`) but no further sanitizing is done.

import {
  mdiCodeJson,
  mdiConsole,
  mdiLanguageCsharp,
  mdiLanguageGo,
  mdiLanguageJava,
  mdiLanguageJavascript,
  mdiLanguagePhp,
  mdiLanguagePython,
  mdiLanguageTypescript,
  mdiPowershell,
} from '@mdi/js'
import markdownIt, { type MarkdownIt, type Token } from 'markdown-it'
import hljs from 'highlight.js/lib/core'
import bash from 'highlight.js/lib/languages/bash'
import csharp from 'highlight.js/lib/languages/csharp'
import dart from 'highlight.js/lib/languages/dart'
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
  dart,
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
  })

  // Fenced code renders as the landing page's "code window": three dots and
  // a language tab. Consecutive fences (e.g. the same request in cURL,
  // JavaScript and Python) share one window with a tab per language; the
  // docs page switches tabs through [data-code-tab] (see pages/docs).
  md.core.ruler.push('docs_code_groups', (state) => {
    const tokens = state.tokens
    let start = 0
    while (start < tokens.length) {
      if (tokens[start]!.type !== 'fence') {
        start++
        continue
      }
      let end = start
      while (tokens[end + 1]?.type === 'fence') end++
      for (let i = start; i <= end; i++) {
        tokens[i]!.meta = {
          ...tokens[i]!.meta,
          codeGroup: { index: i - start, last: i === end },
        }
      }
      start = end + 1
    }
  })

  md.renderer.rules.fence = (tokens, index) => {
    const token = tokens[index]!
    const lang = token.info.trim().split(/\s+/)[0] ?? ''
    const group = token.meta?.codeGroup as
      { index: number; last: boolean } | undefined
    const position = group?.index ?? 0
    const language = lang && hljs.getLanguage(lang) ? lang : null
    const body = language
      ? hljs.highlight(token.content, { language, ignoreIllegals: true }).value
      : md.utils.escapeHtml(token.content)
    const className = language ? ` language-${language}` : ''
    const hidden = position === 0 ? '' : ' hidden'
    const panel = `<pre class="docs-code code-window-panel" data-code-panel="${position}" data-lang="${md.utils.escapeHtml(lang)}"${hidden}><code class="hljs${className}">${body}</code></pre>`

    let open = ''
    if (position === 0) {
      const members: Token[] = []
      for (let i = index; tokens[i]?.type === 'fence'; i++) {
        members.push(tokens[i]!)
      }
      const tabs = members
        .map((member, i) => {
          const memberLang = member.info.trim().split(/\s+/)[0] ?? ''
          return codeTab(memberLang, i, members.length > 1, md)
        })
        .join('')
      open =
        '<div class="code-window">' +
        '<div class="code-window-bar"><span></span><span></span><span></span></div>' +
        `<div class="code-window-tabs" role="tablist">${tabs}</div>`
    }
    const close = !group || group.last ? '</div>' : ''
    return `${open}${panel}${close}`
  }

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

interface CodeLanguage {
  label: string
  /** Single-colour Material Design icon path, painted with `color`. */
  icon?: string
  color?: string
  /** Full-colour logo markup (24×24 viewBox) for brands without an mdi icon. */
  svg?: string
}

// Dart's multi-colour mark, simplified to flat polygons.
const DART_LOGO =
  '<svg viewBox="0 0 1024 1024" aria-hidden="true">' +
  '<path fill="#0081C6" d="M215 213 10 520l20 80 118 120 74 27z"/>' +
  '<path fill="#00D2B8" d="m215 213 310-205 75 12 125 125 25 73z"/>' +
  '<path fill="#00D2B8" d="m215 213 660 660-128 59-525-185z"/>' +
  '<path fill="#0081C6" d="m215 213 535 5 185 525-60 130z"/>' +
  '<path fill="#00A8E1" d="m750 218 105 52 165 168v390l-145 45 60-130z"/>' +
  '<path fill="#55DDCA" d="m222 747 53 103 175 172h380l45-149-128 59z"/>' +
  '</svg>'

// Same icons and colors as the landing page's code window.
const CODE_LANGUAGES: Record<string, CodeLanguage> = {
  bash: { label: 'cURL', icon: mdiPowershell, color: '#4f8ef7' },
  shell: { label: 'Shell', icon: mdiConsole, color: '#4f8ef7' },
  javascript: {
    label: 'JavaScript',
    icon: mdiLanguageJavascript,
    color: '#efd81d',
  },
  typescript: {
    label: 'TypeScript',
    icon: mdiLanguageTypescript,
    color: '#3178c6',
  },
  python: { label: 'Python', icon: mdiLanguagePython, color: '#ffffff' },
  php: { label: 'PHP', icon: mdiLanguagePhp, color: '#777bb3' },
  go: { label: 'Go', icon: mdiLanguageGo, color: '#00aed8' },
  java: { label: 'Java', icon: mdiLanguageJava, color: '#0c89c7' },
  csharp: { label: 'C#', icon: mdiLanguageCsharp, color: '#a179dc' },
  dart: { label: 'Dart', svg: DART_LOGO },
  json: { label: 'JSON', icon: mdiCodeJson, color: '#f5a623' },
}

function codeTab(
  lang: string,
  index: number,
  interactive: boolean,
  md: MarkdownIt,
): string {
  const known = CODE_LANGUAGES[lang]
  const label = md.utils.escapeHtml(
    known?.label ?? (lang && lang !== 'text' ? lang.toUpperCase() : 'Text'),
  )
  let icon = ''
  if (known?.svg) icon = known.svg
  else if (known?.icon) {
    icon = `<svg viewBox="0 0 24 24" aria-hidden="true" style="color:${known.color}"><path fill="currentColor" d="${known.icon}"/></svg>`
  }
  const active = index === 0
  const classes = `code-window-tab${active ? ' is-active' : ''}`
  if (!interactive) return `<span class="${classes}">${icon}${label}</span>`
  return `<button type="button" role="tab" class="${classes}" data-code-tab="${index}" aria-selected="${active}">${icon}${label}</button>`
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
