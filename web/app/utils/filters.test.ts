import assert from 'node:assert/strict'
import test from 'node:test'
import { humanizeTime, humanizeTimeShort } from './filters.ts'

const minutesAgo = (minutes: number) =>
  new Date(Date.now() - minutes * 60_000).toISOString()

test('relative times follow the app language', () => {
  assert.equal(humanizeTimeShort(minutesAgo(2), 'es'), 'hace 2 minutos')
  assert.equal(humanizeTimeShort(minutesAgo(2), 'en'), '2 minutes ago')
  assert.equal(humanizeTimeShort(minutesAgo(180), 'es'), 'hace 3 horas')
  assert.equal(humanizeTimeShort(minutesAgo(0), 'es'), 'ahora')
})

test('durations follow the app language', () => {
  assert.match(humanizeTime(minutesAgo(5), 'es'), /^5 minutos/)
  assert.match(humanizeTime(minutesAgo(5), 'en'), /^5 minutes/)
})
