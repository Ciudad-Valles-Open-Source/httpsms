import {
  getAccountBlockedCode,
  type AccountBlockedCode,
} from '../utils/account-status'

let authToken: string | null = null
let apiKey: string | null = null
let accountBlockedHandler: ((code: AccountBlockedCode) => void) | null = null

export function setAuthHeader(token: string | null) {
  authToken = token
}

export function setApiKey(key: string | null) {
  apiKey = key
}

/**
 * Registers the callback run when the API rejects a request because the
 * signed-in account is pending approval or suspended (HTTP 403 with
 * `data.code` = `account_pending` / `account_suspended`).
 */
export function setAccountBlockedHandler(
  handler: ((code: AccountBlockedCode) => void) | null,
) {
  accountBlockedHandler = handler
}

function createApiFetch() {
  const config = useRuntimeConfig()
  const publicConfig = config.public as Record<string, string>
  const baseURL = publicConfig.apiBaseUrl

  return $fetch.create({
    baseURL,
    headers: {
      'X-Client-Version': publicConfig.clientVersion || 'dev',
    },
    onRequest({ options }) {
      const headers = new Headers(options.headers)
      if (authToken) {
        headers.set('Authorization', `Bearer ${authToken}`)
      }
      if (apiKey) {
        headers.set('x-api-key', apiKey)
      }
      options.headers = headers
    },
    onResponseError({ response }) {
      if (response.status !== 403 || !accountBlockedHandler) return
      const code = getAccountBlockedCode(response._data)
      if (code) accountBlockedHandler(code)
    },
  })
}

export function useApi() {
  return { apiFetch: createApiFetch(), setAuthHeader, setApiKey }
}

export function useApiComposable() {
  return {
    useApi: createApiFetch,
  }
}
