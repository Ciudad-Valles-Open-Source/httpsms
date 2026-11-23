# Handling secrets

A **secret** is any value that grants access to something: passwords, API keys, tokens, private keys, Firebase/Google credentials, Android keystores. **They must never be in the repository**, even a private one: git history is permanent and a single push or fork makes them public.

The rule in this project:

> **Git only contains templates with empty or example values. Real values live in `.env` (ignored by git), in Coolify's environment variables or in GitHub Actions secrets.**

---

## Where each thing goes

| Environment | Where secrets are stored | In git? |
| --- | --- | --- |
| Your computer (development, `docker compose`) | `.env` at the repo root (copy of `.env.example`) | No, ignored |
| Production (Coolify) | The resource's *Environment Variables*, marked as secret | No |
| Building the Android app locally | Environment variables or `~/.gradle/gradle.properties` | No |
| Building the Android app in GitHub Actions | *Settings → Secrets and variables → Actions* | No |
| Templates (`.env.example`, `web/.env.docker`, `web/.env.production`, `api/.env.*`, `tests/.env.test`) | Variable names with empty or example values only | Yes |

### Server secrets

Put them in `.env` (local) and Coolify (production), never in `.env.example` or the `.env.docker` / `.env.production` templates: `POSTGRES_PASSWORD`, `FIREBASE_CREDENTIALS` (service account JSON with a private key), `EVENTS_QUEUE_USER_API_KEY`, `SOKETI_SECRET` (keep `SOKETI_KEY` with it), `SMTP_PASSWORD`, `CLOUDFLARE_TURNSTILE_SECRET_KEY`, `AXIOM_TOKEN`.

> **About `FIREBASE_API_KEY` (web):** the Firebase web key is not strictly secret (it ships inside the browser JavaScript), but **don't commit it** anyway, and restrict it to your domains in Google Cloud Console → *APIs & Services → Credentials*.

### Android secrets

- `google-services.json` (your Firebase project): `android/app/google-services.json`, ignored by git.
- Signing keystore (`*.jks`) and passwords: keep the file outside the repo and reference it from `android/key.properties` (ignored).
- Optional `HTTPSMS_SENTRY_DSN`, `HTTPSMS_LOGS_URL`, `HTTPSMS_LOGS_TOKEN`: environment variables or Gradle `-P` properties.

---

## Storing them as secrets

**Coolify:** resource → **Environment Variables** → add each variable and tick **Is Secret**. Variables the web app needs at build time (`FIREBASE_*`, `SOKETI_KEY`, `API_BASE_URL`…) must also be **Build Variables**. Redeploy.

**GitHub Actions:** **Settings → Secrets and variables → Actions → New repository secret**. Store files base64-encoded (`base64 -w0 android/app/google-services.json`) and restore them in the workflow:

```yaml
- name: Restore google-services.json
  run: echo "${{ secrets.GOOGLE_SERVICES_JSON_B64 }}" | base64 -d > android/app/google-services.json
- name: Build APK
  env:
    HTTPSMS_DEFAULT_SERVER_URL: ${{ vars.HTTPSMS_DEFAULT_SERVER_URL }}
    HTTPSMS_SENTRY_DSN: ${{ secrets.HTTPSMS_SENTRY_DSN }}
  run: cd android && ./gradlew assembleRelease
```

---

## Built-in protections

1. **`.gitignore`** ignores `.env`, `google-services.json`, keystores, private keys, `key.properties` and Firebase/GCP credential files.
2. **Pre-commit secret scanner** (`scripts/check-secrets.sh`): scans only what you are committing and blocks sensitive files (even if forced with `git add -f`), well-known key formats (Google `AIza…`, `PRIVATE KEY` blocks, GitHub/Slack/Stripe/AWS tokens, Sentry DSNs, Axiom tokens, JWTs) and non-placeholder values of `…_SECRET`, `…_PASSWORD`, `…_TOKEN`, `…_API_KEY`, `…_CREDENTIALS`, `…_DSN` variables in `.env*` files. Run it manually with `bash scripts/check-secrets.sh`; mark false positives with a `secret-scan: ignore` comment. The hook is installed by `pnpm install` in `web/`.
3. Also enable **Secret scanning** and **Push protection** in GitHub (**Settings → Code security**) so pushes containing known secrets are rejected server-side.

---

## If a secret leaks

1. **Revoke and rotate it first.** Deleting it from the repo does not help: it may already be copied.
2. Update the new value in `.env`, Coolify and GitHub Secrets, then redeploy.
3. Optionally rewrite history with [`git filter-repo`](https://github.com/newren/git-filter-repo) and force-push.

> The fork's history includes upstream commits with configuration for **upstream's** services (its `google-services.json` and Firebase web config). They are not yours and give no access to your installation; they were removed from the current tree. None of your real `.env` values appear in the history.
