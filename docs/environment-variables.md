# Environment Variables Reference

*🌎 [Leer en Español](environment-variables_ES.md)*

This document lists every environment variable accepted by the httpSMS API and web frontend, with descriptions, expected values, and whether each variable is required.

Variables marked **Required** must be set before the application will start correctly. Variables marked **Optional** fall back to safe defaults when omitted.

---

## API Environment Variables

These variables are read by the Go API server at runtime. They must be supplied via the container environment — the production image is built with `--dotenv=false`, which disables `.env` file loading.

### Core

| Variable          | Required | Default                 | Description                                                                                                                                        |
| ----------------- | -------- | ----------------------- | -------------------------------------------------------------------------------------------------------------------------------------------------- |
| `ENV`             | No       | `local`                 | Runtime environment identifier. Set to `production` in production deployments. Controls logging verbosity and other environment-specific behavior. |
| `APP_PORT`        | No       | `8000`                  | TCP port on which the API server listens.                                                                                                          |
| `APP_URL`         | No       | `http://localhost:3000` | Public URL of the web frontend. Used to generate links in outbound emails.                                                                         |
| `APP_NAME`        | No       | `httpSMS`               | Display name of the application. Used in emails and logs.                                                                                          |
| `JWT_ISSUER`      | No       | Slug of `APP_NAME` (`httpsms`) | Value of the `iss` claim in the JWTs the API signs (the `Authorization` header of webhook requests and of HTTP push notifications to URL-backed phone gateways). When empty it defaults to a slug of `APP_NAME`: lowercase, spaces become `_`, accents are stripped to ASCII (`Ñandú` → `nandu`) and every character outside `[a-z0-9._-]` is removed, e.g. `Arquos SMS!` → `arquos_sms`. Falls back to `httpsms` if the slug is empty. |
| `SWAGGER_HOST`    | No       | `localhost:8000`        | Host shown in the Swagger UI. Set to your API's public hostname in production.                                                                     |
| `USE_HTTP_LOGGER` | No       | `true`                  | When `true`, logs incoming HTTP requests.                                                                                                          |

### Database

| Variable                  | Required | Default | Description                                                                                                                                   |
| ------------------------- | -------- | ------- | --------------------------------------------------------------------------------------------------------------------------------------------- |
| `DATABASE_URL`            | Yes      | —       | PostgreSQL connection string used for the primary connection pool. Format: `postgresql://user:password@host:port/database`                    |
| `DATABASE_URL_DEDICATED`  | Yes      | —       | PostgreSQL connection string used for the dedicated connection pool (migrations and long-running queries). May be the same as `DATABASE_URL`. |
| `DATABASE_MIGRATION_SKIP` | No       | —       | When set to any non-empty value, skips running GORM AutoMigrate on startup. Useful when migrations are managed externally.                    |

### Redis

| Variable    | Required | Default | Description                                                                 |
| ----------- | -------- | ------- | --------------------------------------------------------------------------- |
| `REDIS_URL` | Yes      | —       | Redis connection string. Format: `redis://[:password@]host:port[/database]` |

### Firebase

| Variable               | Required | Default | Description                                                                               |
| ---------------------- | -------- | ------- | ----------------------------------------------------------------------------------------- |
| `GCP_PROJECT_ID`       | Yes      | —       | Firebase project ID. Also used as the Google Cloud project identifier.                    |
| `FIREBASE_CREDENTIALS` | Yes      | —       | Full content of the Firebase service account JSON file, provided as a single-line string. |

### Account Approval (Administrators)

| Variable       | Required | Default | Description |
| -------------- | -------- | ------- | ----------- |
| `ADMIN_EMAILS` | No       | —       | Comma-separated list of administrator emails, compared case-insensitively after trimming (e.g. `admin@example.com,ops@example.com`). When set, the **approval gate** is enabled: anyone can sign up with Firebase, but new accounts stay `pending` and every API request made with them (dashboard token, user API key or phone API key) is rejected with `403` and `data.code` = `account_pending` (or `account_suspended`) until an administrator approves them. `GET /v1/users/me` and `DELETE /v1/users/me` keep working so a pending user can see their status or delete their account. A request is an administrator request **only** when it is signed in on the dashboard with a Firebase ID token whose email is listed here **and verified** (`email_verified`); API keys never grant admin rights. Administrators manage accounts with `/v1/admin/users`. Accounts that existed before the variable was set stay `active`. When empty, registrations are open (every account is active immediately) and the API logs a warning at startup. Status changes apply immediately on the API instance that handled them; when running several API replicas, other replicas may take up to 5 seconds to see the change. If SMTP is configured, administrators are emailed when a new account waits for approval and users are emailed when approved. |

### Event Queue (System User)

These variables configure the internal system user that drives the asynchronous event queue.

| Variable                    | Required | Default                           | Description                                                                                                                              |
| --------------------------- | -------- | --------------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------- |
| `EVENTS_QUEUE_TYPE`         | No       | `emulator`                        | Queue backend. `emulator` uses an in-process HTTP emulator suitable for Docker deployments. Set to `cloud_tasks` for Google Cloud Tasks. |
| `EVENTS_QUEUE_NAME`         | No       | `events-local`                    | Name of the queue used by the emulator or Cloud Tasks.                                                                                   |
| `EVENTS_QUEUE_ENDPOINT`     | No       | `http://localhost:8000/v1/events` | HTTP endpoint where the queue delivers event payloads. Update to the API's internal hostname in multi-container deployments.             |
| `EVENTS_QUEUE_USER_ID`      | No       | Auto-generated                    | UUID of the system user. Auto-generated on first boot if empty; printed to stdout. Persist this value across restarts.                   |
| `EVENTS_QUEUE_USER_API_KEY` | No       | Auto-generated                    | API key of the system user. Auto-generated on first boot if empty; printed to stdout. Persist this value across restarts.                |
| `EVENTS_QUEUE_USER_EMAIL`   | No       | `system@httpsms.local`            | Email address of the system user record in the database.                                                                                 |

> **Note:** On the first startup against an empty database, the API automatically creates the system user using the values above. If `EVENTS_QUEUE_USER_ID` or `EVENTS_QUEUE_USER_API_KEY` are not provided, cryptographically secure values are generated and printed to the container logs. Capture and persist these values before subsequent restarts. See [self-hosting.md — Step 7](./self-hosting.md#step-7--system-user-initialization) for details.

### SMTP

| Variable          | Required | Default                 | Description                                                                                                      |
| ----------------- | -------- | ----------------------- | ---------------------------------------------------------------------------------------------------------------- |
| `SMTP_HOST`       | Yes      | —                       | Hostname of the SMTP server.                                                                                     |
| `SMTP_PORT`       | No       | `587`                   | Port of the SMTP server.                                                                                         |
| `SMTP_USERNAME`   | Yes      | —                       | SMTP authentication username.                                                                                    |
| `SMTP_PASSWORD`   | Yes      | —                       | SMTP authentication password.                                                                                    |
| `SMTP_FROM_NAME`  | No       | `httpSMS`               | Display name shown in the From field of outbound emails.                                                         |
| `SMTP_FROM_EMAIL` | No       | `noreply@httpsms.local` | Email address used in the From field of outbound emails. Must match the domain authorized by your SMTP provider. |

### Cloudflare Turnstile

| Variable                          | Required | Default | Description                                                                                                                                          |
| --------------------------------- | -------- | ------- | ---------------------------------------------------------------------------------------------------------------------------------------------------- |
| `CLOUDFLARE_TURNSTILE_SECRET_KEY` | No       | —       | Turnstile secret key used by the API to validate CAPTCHA tokens on the `/v1/messages/search` route. The route is unprotected if this value is empty. |

### Optional Integrations

| Variable                | Required | Default | Description                                                                                                                    |
| ----------------------- | -------- | ------- | ------------------------------------------------------------------------------------------------------------------------------ |
| `SOKETI_APP_ID`         | No       | —       | Soketi application ID. Enables real-time updates in the web UI via WebSockets. See [realtime.md](./realtime.md). |
| `SOKETI_KEY`            | No       | —       | Soketi public key. Must be URL-safe (letters and numbers only, e.g. `openssl rand -hex 32`) because it is part of the WebSocket URL. |
| `SOKETI_SECRET`         | No       | —       | Soketi secret key.                                                                                                             |
| `SOKETI_HOST`           | No       | —       | Soketi Host (optional). Used when connecting to the WebSocket server.                                         |
| `SOKETI_SECURE`         | No       | `false` | Soketi Secure flag (optional). Set to `false` when using local HTTP for the WebSocket server.                             |
| `SOKETI_CLUSTER`        | No       | —       | Soketi cluster identifier (e.g., `mt1`, `eu`).                                                                                 |
| `AXIOM_TOKEN`           | No       | —       | Axiom API token for centralized log and trace ingestion.                                                                       |
| `AXIOM_DATASET_EVENTS`  | No       | —       | Axiom dataset name for logs and traces.                                                                                        |
| `AXIOM_DATASET_METRICS` | No       | —       | Axiom dataset name for metrics.                                                                                                |
| `OTEL_SERVICE_NAME`     | No       | Slug of `APP_NAME` (`httpsms`) | OpenTelemetry `service.name` reported with traces and metrics (and the `service.name` attribute sent with Plunk contacts). Defaults to the same `APP_NAME` slug as `JWT_ISSUER`. |
| `GCS_BUCKET_NAME`       | No       | —       | Google Cloud Storage bucket name for MMS attachment storage. When empty, attachments are stored in memory and lost on restart. |

---

## Web Frontend Environment Variables

The web frontend is a statically generated Nuxt SPA. All environment variables are **injected at image build time** via Docker build arguments. They are not read from the container environment at runtime.

When using Docker Compose, the `docker-compose.yml` passes these values from the host environment to the build process automatically. When building the image manually or in a CI pipeline, pass them as `--build-arg` flags.

| Variable                        | Required | Default                                  | Description                                                                                                                                                     |
| ------------------------------- | -------- | ---------------------------------------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `API_BASE_URL`                  | Yes      | `http://localhost:8000`                  | Public URL of the API. The browser makes requests to this URL directly. Must be reachable from the end user's browser, not just from within the Docker network. |
| `APP_URL`                       | No       | `http://localhost:3000`                  | Public URL of the web application. Used for SEO metadata and email links.                                                                                       |
| `APP_NAME`                      | No       | `httpSMS`                                | Display name of the application.                                                                                                                                |
| `APP_ENV`                       | No       | `production`                             | Runtime environment label.                                                                                                                                      |
| `APP_GITHUB_URL`                | No       | `https://github.com/NdoleStudio/httpsms` | Link to the project source repository shown in the UI.                                                                                                          |
| `APP_DOWNLOAD_URL`              | No       | (GitHub releases URL)                    | Direct download URL for the Android APK.                                                                                                                        |
| `FIREBASE_API_KEY`              | Yes      | —                                        | Firebase Web SDK API key.                                                                                                                                       |
| `FIREBASE_AUTH_DOMAIN`          | Yes      | —                                        | Firebase authentication domain.                                                                                                                                 |
| `FIREBASE_PROJECT_ID`           | Yes      | —                                        | Firebase project ID.                                                                                                                                            |
| `FIREBASE_STORAGE_BUCKET`       | Yes      | —                                        | Firebase Storage bucket name.                                                                                                                                   |
| `FIREBASE_MESSAGING_SENDER_ID`  | Yes      | —                                        | Firebase Cloud Messaging sender ID.                                                                                                                             |
| `FIREBASE_APP_ID`               | Yes      | —                                        | Firebase App ID.                                                                                                                                                |
| `FIREBASE_MEASUREMENT_ID`       | No       | —                                        | Firebase Analytics measurement ID. Optional.                                                                                                                    |
| `CLOUDFLARE_TURNSTILE_SITE_KEY` | No       | —                                        | Turnstile site key embedded in the search page form.                                                                                                            |
| `SOKETI_KEY`                    | No       | —                                        | Soketi public key for real-time updates. Must be URL-safe (letters and numbers only).                                                                         |
| `SOKETI_CLUSTER`                | No       | —                                        | Soketi cluster for real-time updates.                                                                                                                           |
| `USE_SOKETI`                    | No       | `false`                                  | Whether to use a Soketi instance. If true, uses the browser's current host and `/app/` route.                                     |

---

## PostgreSQL Container Variables

These variables configure the `postgres` service in the Docker Compose stack.

| Variable            | Required | Default    | Description                                                        |
| ------------------- | -------- | ---------- | ------------------------------------------------------------------ |
| `POSTGRES_DB`       | No       | `httpsms`  | Name of the database to create.                                    |
| `POSTGRES_USER`     | No       | `httpsms`  | PostgreSQL username.                                               |
| `POSTGRES_PASSWORD` | Yes      | `changeme` | PostgreSQL password. Change this before any production deployment. |

---

## Notes on Variable Precedence

- Variables defined in `.env` at the repository root are automatically loaded by Docker Compose.
- When deploying through an orchestrator such as Coolify, Portainer, or Docker Swarm, variables are injected by the platform and take precedence over any `.env` file.
- The API binary is started with `--dotenv=false`. It never reads `.env` files from disk at runtime; all configuration must arrive as process environment variables.
- The web frontend's configuration is frozen at the time the Docker image is built. Changing a web environment variable requires rebuilding the image.
