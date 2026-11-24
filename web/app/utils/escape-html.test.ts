import { test } from 'node:test'
import assert from 'node:assert/strict'
import { escapeHtml } from './escape-html'

test('escapes HTML special characters', () => {
  assert.equal(
    escapeHtml('<img src=x onerror="alert(1)">'),
    '&lt;img src=x onerror=&quot;alert(1)&quot;&gt;',
  )
  assert.equal(escapeHtml(`Tom & Jerry's`), 'Tom &amp; Jerry&#39;s')
})

test('handles empty values', () => {
  assert.equal(escapeHtml(undefined), '')
  assert.equal(escapeHtml(null), '')
  assert.equal(escapeHtml(''), '')
})
