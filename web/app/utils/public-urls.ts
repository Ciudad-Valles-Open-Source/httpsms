/**
 * Absolute URLs of this self-hosted install (dashboard, API, GitHub project),
 * built from the runtime config so that pages never link to or load content
 * from the upstream hosted service.
 */
export interface PublicUrlOptions {
  appUrl: string
  apiBaseUrl: string
  githubUrl: string
  downloadUrl: string
  /** Used to resolve relative `appUrl` / `apiBaseUrl` values (e.g. `/api`). */
  origin: string
}

export interface PublicUrls {
  appUrl: string
  apiBaseUrl: string
  settingsUrl: string
  docsUrl: string
  sendMessageUrl: string
  githubUrl: string
  githubIssuesUrl: string
  /** `owner/repo` part of the GitHub URL (empty if it is not a GitHub URL). */
  githubRepo: string
  downloadUrl: string
}

function trimTrailingSlash(url: string): string {
  return url.replace(/\/+$/, '')
}

function toAbsolute(url: string, origin: string): string {
  try {
    return new URL(url || '/', origin).toString()
  } catch {
    return url
  }
}

export function buildPublicUrls(options: PublicUrlOptions): PublicUrls {
  const appUrl = trimTrailingSlash(toAbsolute(options.appUrl, options.origin))
  const apiBaseUrl = trimTrailingSlash(
    toAbsolute(options.apiBaseUrl, options.origin),
  )
  const githubUrl = trimTrailingSlash(options.githubUrl)
  const repoMatch = githubUrl.match(/^https?:\/\/github\.com\/([^/]+\/[^/]+)/)

  return {
    appUrl,
    apiBaseUrl,
    settingsUrl: `${appUrl}/settings`,
    docsUrl: `${appUrl}/docs`,
    sendMessageUrl: `${apiBaseUrl}/v1/messages/send`,
    githubUrl,
    githubIssuesUrl: `${githubUrl}/issues`,
    githubRepo: repoMatch?.[1] ?? '',
    downloadUrl: options.downloadUrl,
  }
}

/** Absolute URL of a file in `public/`, e.g. for `og:image`. */
export function publicAssetUrl(appUrl: string, path: string): string {
  return `${trimTrailingSlash(appUrl)}${path.startsWith('/') ? path : `/${path}`}`
}
