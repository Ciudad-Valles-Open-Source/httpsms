import { buildPublicUrls, publicAssetUrl } from '~/utils/public-urls'

/** URLs of this self-hosted install, from the runtime config. */
export function usePublicUrls() {
  const config = useRuntimeConfig().public
  const urls = buildPublicUrls({
    appUrl: config.appUrl,
    apiBaseUrl: config.apiBaseUrl,
    githubUrl: config.appGithubUrl,
    downloadUrl: config.appDownloadUrl,
    origin: import.meta.client ? window.location.origin : 'http://localhost',
  })

  return {
    ...urls,
    /** Absolute URL of a file in `public/` (for `og:image`). */
    assetUrl: (path: string) => publicAssetUrl(urls.appUrl, path),
  }
}
