import assert from 'node:assert/strict'
import { test } from 'node:test'

import {
  buildDocPlaceholders,
  renderDocMarkdown,
  slugifyDocHeading,
  substituteDocPlaceholders,
} from './docs-markdown'
import { adjacentDocPages, docRoutes, DOCS_PAGES } from './docs-manifest'

test('buildDocPlaceholders', async (t) => {
  const base = {
    apiBaseUrl: 'https://api.example.com/',
    appUrl: 'https://app.example.com/',
    appDownloadUrl: 'https://app.example.com/app.apk',
  }

  await t.test(
    'strips trailing slashes and points to the English Swagger UI',
    () => {
      assert.deepEqual(buildDocPlaceholders({ ...base, locale: 'en' }), {
        API_BASE_URL: 'https://api.example.com',
        APP_URL: 'https://app.example.com',
        APP_DOWNLOAD_URL: 'https://app.example.com/app.apk',
        API_DOCS_URL: 'https://api.example.com/index.html',
      })
    },
  )

  await t.test('points to the Spanish Swagger UI for es', () => {
    assert.equal(
      buildDocPlaceholders({ ...base, locale: 'es' }).API_DOCS_URL,
      'https://api.example.com/es/index.html',
    )
  })
})

test('substituteDocPlaceholders replaces known names and keeps unknown ones', () => {
  assert.equal(
    substituteDocPlaceholders(
      'curl {{API_BASE_URL}}/v1 {{ APP_URL }} {{OTHER}}',
      {
        API_BASE_URL: 'https://api.test',
        APP_URL: 'https://app.test',
      },
    ),
    'curl https://api.test/v1 https://app.test {{OTHER}}',
  )
})

test('slugifyDocHeading strips accents and punctuation', () => {
  assert.equal(
    slugifyDocHeading('Introducción a los Webhooks'),
    'introduccion-a-los-webhooks',
  )
  assert.equal(
    slugifyDocHeading('message.phone.received'),
    'messagephonereceived',
  )
  assert.equal(slugifyDocHeading('  ¿Qué es?  '), 'que-es')
})

test('renderDocMarkdown', async (t) => {
  const doc = renderDocMarkdown(
    [
      '# Primeros pasos',
      '',
      'Envía tu **primer** SMS',
      'con la API.',
      '',
      '## Paso 1',
      '## Paso 1',
      '',
      '```bash',
      'curl -X POST "https://api.test"',
      '```',
      '',
      '> **Nota:** importante.',
      '',
      '> **Advertencia:** cuidado.',
      '',
      '| a | b |',
      '|---|---|',
      '| 1 | 2 |',
      '',
      '[externo](https://example.com) [interno](/docs/webhooks-events)',
      '',
      '<script>alert(1)</script>',
    ].join('\n'),
  )

  await t.test('extracts the title and description', () => {
    assert.equal(doc.title, 'Primeros pasos')
    assert.equal(doc.description, 'Envía tu primer SMS con la API.')
  })

  await t.test('adds unique heading ids and anchors', () => {
    assert.match(doc.html, /<h1 id="primeros-pasos">Primeros pasos<\/h1>/)
    assert.match(
      doc.html,
      /<h2 id="paso-1">Paso 1 <a class="docs-heading-anchor" href="#paso-1"/,
    )
    assert.match(doc.html, /<h2 id="paso-1-1">/)
  })

  await t.test('highlights fenced code', () => {
    assert.match(
      doc.html,
      /<pre class="docs-code code-window-panel" data-code-panel="0" data-lang="bash"><code class="hljs language-bash">/,
    )
    assert.match(doc.html, /hljs-string/)
  })

  await t.test('groups consecutive fences into one code window', () => {
    const grouped = renderDocMarkdown(
      '# T\n\n```bash\ncurl x\n```\n\n```python\nprint(1)\n```\n\nText\n\n```json\n{}\n```\n',
    )
    assert.equal(grouped.html.match(/<div class="code-window">/g)?.length, 2)
    assert.match(
      grouped.html,
      /<button type="button" role="tab" class="code-window-tab is-active" data-code-tab="0" aria-selected="true">.*cURL<\/button><button[^>]*data-code-tab="1" aria-selected="false">.*Python<\/button>/,
    )
    assert.match(grouped.html, /data-code-panel="1" data-lang="python" hidden>/)
    // A lone fence gets a static label, not a clickable tab.
    assert.match(
      grouped.html,
      /<span class="code-window-tab is-active">.*JSON<\/span>/,
    )
  })

  await t.test('marks callouts, wraps tables and handles links', () => {
    assert.match(
      doc.html,
      /<blockquote class="docs-callout docs-callout-info">/,
    )
    assert.match(
      doc.html,
      /<blockquote class="docs-callout docs-callout-warning">/,
    )
    assert.match(doc.html, /<div class="docs-table"><table>/)
    assert.match(
      doc.html,
      /<a href="https:\/\/example.com" target="_blank" rel="noopener noreferrer">/,
    )
    assert.match(doc.html, /<a href="\/docs\/webhooks-events">interno<\/a>/)
  })

  await t.test('does not render raw HTML', () => {
    assert.doesNotMatch(doc.html, /<script>/)
  })
})

test('docs manifest', async (t) => {
  await t.test('lists every route including the index', () => {
    assert.equal(docRoutes().length, DOCS_PAGES.length + 1)
    assert.ok(docRoutes().includes('/docs'))
    assert.ok(docRoutes().includes('/docs/webhooks-events'))
  })

  await t.test('computes previous and next pages', () => {
    assert.deepEqual(
      [
        adjacentDocPages('getting-started').previous,
        adjacentDocPages('getting-started').next?.slug,
      ],
      [null, 'webhooks-introduction'],
    )
    assert.equal(adjacentDocPages('outgoing-message-queue').next, null)
    assert.deepEqual(adjacentDocPages('missing'), {
      previous: null,
      next: null,
    })
  })
})
