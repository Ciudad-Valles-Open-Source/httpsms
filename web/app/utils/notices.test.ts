import assert from 'node:assert/strict'
import test from 'node:test'
import { findNoticeTranslation, translateNotice } from './notices.ts'

test('translates known API and store messages', () => {
  assert.deepEqual(findNoticeTranslation('phone updated successfully'), {
    key: 'notices.phoneUpdated',
  })
  assert.deepEqual(findNoticeTranslation('Contact created successfully'), {
    key: 'notices.contactCreated',
  })
  assert.deepEqual(
    findNoticeTranslation('duplicate send event received.'),
    null,
  )
})

test('translates messages with values', () => {
  assert.deepEqual(
    findNoticeTranslation('Added 3 out of 5 messages to the queue'),
    { key: 'notices.bulkQueued', params: { added: '3', total: '5' } },
  )
  assert.equal(
    findNoticeTranslation('validation errors while sending message')?.key,
    'notices.validationErrors',
  )
  assert.equal(
    findNoticeTranslation('cannot find contact with ID [abc]')?.key,
    'notices.notFound',
  )
})

test('leaves unknown messages untouched', () => {
  const t = (key: string) => `T:${key}`
  assert.equal(translateNotice('Something custom', t), 'Something custom')
  assert.equal(
    translateNotice('webhook deleted successfully', t),
    'T:notices.webhookDeleted',
  )
})
