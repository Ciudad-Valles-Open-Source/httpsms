# Continuous integration (GitHub Actions)

Two workflows live in `.github/workflows/`. Neither deploys your server: **Coolify deploys the API and web app straight from the repository**. GitHub Actions only checks the code and publishes the Android APK.

| Workflow | Runs on | Does |
| --- | --- | --- |
| `ci.yaml` | Every push to `main` and every pull request | Secret scan, API tests, integration tests, web and Android builds, Docker image builds |
| `release.yaml` | Pushing a `v*` tag (or manually) | Builds the signed APK and publishes it as `HttpSms.apk` on a GitHub Release |

## `ci.yaml` jobs

| Job | Checks | Secrets |
| --- | --- | --- |
| **Secret Scan** | No commit in the push/PR contains keys, passwords, keystores or credentials (`scripts/check-secrets.sh --range`) | None |
| **API (Go)** | `gofmt`, `go vet`, build, unit tests (including the Swagger translation coverage test) | None |
| **Integration Tests** | API + emulators in Docker, end-to-end tests; Firebase credentials are generated on the fly | None |
| **Web (Nuxt)** | Frozen-lockfile install, lint, tests, `pnpm run generate` | None (optional variables) |
| **Android App** | Unit tests and debug APK build | Optional `GOOGLE_SERVICES_JSON_B64` |
| **Docker Images** | API and web Dockerfiles build (nothing is pushed) | None |

Pull requests from forks get no secrets, so the Android job falls back to a compile-only placeholder `google-services.json`.

## Releasing the APK

Configure the secrets below once, then:

```bash
git tag v1.0.0
git push origin v1.0.0
```

The workflow builds and tests the signed APK and creates the `v1.0.0` release with `HttpSms.apk` and `HttpSms.apk.sha256`. Because the file name never changes, `https://github.com/<org>/<repo>/releases/latest/download/HttpSms.apk` (the dashboard's default `APP_DOWNLOAD_URL`) always serves the newest version. You can also run it from **Actions → Release Android App → Run workflow** with an existing tag.

## Secrets and variables

**Settings → Secrets and variables → Actions.**

Secrets: `GOOGLE_SERVICES_JSON_B64` (`base64 -w0 android/app/google-services.json`), `ANDROID_KEYSTORE_B64` (`base64 -w0 your-release.jks`), `ANDROID_KEYSTORE_PASSWORD`, `ANDROID_KEY_ALIAS`, `ANDROID_KEY_PASSWORD` — all required for releases; optional `HTTPSMS_SENTRY_DSN`, `HTTPSMS_LOGS_URL`, `HTTPSMS_LOGS_TOKEN`.

Variables (public): `HTTPSMS_DEFAULT_SERVER_URL` (API URL pre-filled in the app), optional `API_BASE_URL` / `APP_URL` for the CI web build.

Create a keystore if you don't have one, store it **outside the repo** and back it up — losing it means users can't update the app in place:

```bash
keytool -genkeypair -v -keystore httpsms-release.jks -alias httpsms \
  -keyalg RSA -keysize 4096 -validity 10000
```

## Recommended GitHub settings

Enable *Secret scanning* + *Push protection* (Settings → Code security), require the `Secret Scan`, `API (Go)`, `Web (Nuxt)` and `Android App` checks on `main` (Settings → Branches), and add yourself to `.github/CODEOWNERS`.

See also [secrets.md](./secrets.md).
