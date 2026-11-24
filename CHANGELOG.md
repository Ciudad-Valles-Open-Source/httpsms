# CHANGELOG

*🌎 [Leer en Español](CHANGELOG_ES.md)*

This repository is a 100% free, self-hosted fork of
[NdoleStudio/httpsms](https://github.com/NdoleStudio/httpsms). Versioning
continues after the last release of the original project (`v1.2.1`).

## [2.1.0] - 2026-11-24

### Added

- **Account approval and admin panel.** With the `ADMIN_EMAILS` variable, new
  accounts stay pending until an administrator approves them in `/admin/users`
  (approve, suspend, reactivate, or delete). Pending or suspended accounts
  cannot use the dashboard, the API, or their phone keys. Administrators must
  sign in with a verified email, and API keys never grant admin permissions.
  Optionally notifies administrators and the approved user by email. See
  [docs/admin.md](docs/admin.md).
- `/account-pending` page for pending or suspended accounts.
- `GET /v1/users/me` includes `status` and `is_admin`; new `/v1/admin/users`
  endpoints documented in Swagger (English and Spanish).
- **Complete UI redesign**, aligned with the other Arquos products: Plus Jakarta
  Sans typography, cards, badges, blue palette, and light/dark themes. The
  dashboard features a desktop sidebar navigation, redesigned conversation list
  and thread view (fixing outgoing messages getting clipped on the right),
  loading and empty states, and sectioned Settings. The public site has a new
  header with a mobile menu, landing page, login, docs, and error pages. Tables
  turn into cards on mobile.
- Appearance selector in the dashboard (light or dark mode, and palette).
- Swagger theme matches the dashboard palette.

### Changed

- Sign-in offers only **Google** and **email**. Removed GitHub, which was never
  enabled in Firebase on this installation, so no existing accounts are
  affected.
- If the browser blocks the Google sign-in popup, authentication falls back to a
  redirect flow.

### Fixed

- Administrators lost their permissions because the API key identity overrode
  the Firebase session when the dashboard sent both.
- Sign-in errors were not displayed: clear messages are now shown for blocked
  popups, unauthorized domains, and disabled sign-in methods in Firebase.
- The release workflow failed when looking for `release` unit tests.

### Security

- The API no longer logs full tokens or API keys (only `****` and the last 4
  characters).

## [2.0.0] - 2026-11-24

First standalone release of the fork, based on the original `v1.2.1` with its
security patches and dependency updates.

### Added

- Self-hosted documentation at `/docs` (Spanish and English) and Swagger at
  `/index.html` and `/es/index.html` with a dark theme and full translation.
- **All / Unread** conversation filter and messages up to 2048 characters.
- Working real-time updates with Soketi (single shared WebSocket connection).
- Android app: configurable server URL at build time and from the app, login QR
  code including server and key, and `http://` support with a warning.
- Fork-specific Privacy Policy and Terms of Service.
- Independent GitHub Actions CI and signed APK publication as a release.
- Secret scanner in pre-commit and CI, plus a secrets management guide.
- `JWT_ISSUER` and `OTEL_SERVICE_NAME` environment variables.
- Includes previous fork work: Spanish translation, Docker and Coolify
  self-hosting guide, and server configuration in the app.

### Changed

- **Breaking:** removed paid plans, limits, subscription endpoints, and Lemon
  Squeezy; everything is free with no limits.
- **Breaking:** the issuer (`iss`) of webhook and notification JWTs is
  `JWT_ISSUER` (defaults to normalized `APP_NAME`) instead of `api.httpsms.com`.
- **Breaking:** the login QR code uses the `httpsms://login?server=…&key=…`
  format (the app still accepts previous QR codes).
- `SOKETI_KEY` must be alphanumeric.
- PostgreSQL and Redis are no longer exposed on the host (Docker internal
  network only).
- No dependencies on the hosted service: removed links and assets from
  httpsms.com, as well as Microsoft Clarity, Sentry, and Axiom from the original
  project.

### Fixed

- `PUT /v1/phones` no longer clears the FCM token, sending schedule, or
  auto-reply when omitted from the request.
- The heartbeat monitor restarts when history expires.
- XSS vulnerabilities in text containing contact names, schedules, and phone
  numbers.
- Flaky tests for contacts and conversation archiving.

[Unreleased]: https://github.com/Ciudad-Valles-Open-Source/httpsms/compare/v2.1.0...HEAD
[2.1.0]: https://github.com/Ciudad-Valles-Open-Source/httpsms/compare/v2.0.0...v2.1.0
[2.0.0]: https://github.com/Ciudad-Valles-Open-Source/httpsms/releases/tag/v2.0.0
