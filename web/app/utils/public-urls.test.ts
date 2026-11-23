import assert from 'node:assert/strict'
import { test } from 'node:test'

import { buildPublicUrls, publicAssetUrl } from './public-urls'

test('buildPublicUrls', async (t) => {
  const base = {
    appUrl: 'https://sms.example.com/',
    apiBaseUrl: 'https://api.example.com/',
    githubUrl: 'https://github.com/Ciudad-Valles-Open-Source/httpsms/',
    downloadUrl: 'https://sms.example.com/HttpSms.apk',
    origin: 'https://sms.example.com',
  }

  await t.test('strips trailing slashes and builds derived URLs', () => {
    const urls = buildPublicUrls(base)
    assert.equal(urls.appUrl, 'https://sms.example.com')
    assert.equal(urls.apiBaseUrl, 'https://api.example.com')
    assert.equal(urls.settingsUrl, 'https://sms.example.com/settings')
    assert.equal(urls.docsUrl, 'https://sms.example.com/docs')
    assert.equal(
      urls.sendMessageUrl,
      'https://api.example.com/v1/messages/send',
    )
    assert.equal(
      urls.githubIssuesUrl,
      'https://github.com/Ciudad-Valles-Open-Source/httpsms/issues',
    )
    assert.equal(urls.githubRepo, 'Ciudad-Valles-Open-Source/httpsms')
    assert.equal(urls.downloadUrl, base.downloadUrl)
  })

  await t.test('resolves a relative API URL against the origin', () => {
    const urls = buildPublicUrls({ ...base, apiBaseUrl: '/api' })
    assert.equal(urls.apiBaseUrl, 'https://sms.example.com/api')
  })

  await t.test('leaves githubRepo empty for non-GitHub URLs', () => {
    const urls = buildPublicUrls({
      ...base,
      githubUrl: 'https://gitlab.com/x/y',
    })
    assert.equal(urls.githubRepo, '')
  })
})

test('publicAssetUrl', () => {
  assert.equal(
    publicAssetUrl('https://sms.example.com/', '/header.png'),
    'https://sms.example.com/header.png',
  )
  assert.equal(
    publicAssetUrl('https://sms.example.com', 'img/a.png'),
    'https://sms.example.com/img/a.png',
  )
})
