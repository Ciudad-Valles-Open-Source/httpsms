# Environment Variables Reference

This document lists every environment variable accepted by the httpSMS API and web frontend, with descriptions, expected values, and whether each variable is required.

Variables marked **Required** must be set before the application will start correctly. Variables marked **Optional** fall back to safe defaults when omitted.

---

## API Environment Variables

These variables are read by the Go API server at runtime. They must be supplied via the container environment — the production image is built with `--dotenv=false`, which disables `.env` file loading.

### Core

| Variable | Required | Default | Description |
|---|---|---|---|
| `ENV` | No | `local` | Runtime environment identifier. Set to `production` in production deployments. Controls logging verbosity and other environment-specific behavior. |
| `APP_PORT` | No | `8000` | TCP port on which the API server listens. |
| `APP_URL` | No | `http://localhost:3000` | Public URL of the web frontend. Used to generate links in outbound emails. |
| `APP_NAME` | No | `httpSMS` | Display name of the application. Used in emails and logs. |
| `SWAGGER_HOST` | No | `localhost:8000` | Host shown in the Swagger UI. Set to your API's public hostname in production. |
| `USE_HTTP_LOGGER` | No | `true` | When `true`, logs incoming HTTP requests. |

### Database

| Variable | Required | Default | Description |
|---|---|---|---|
| `DATABASE_URL` | Yes | — | PostgreSQL connection string used for the primary connection pool. Format: `postgresql://user:password@host:port/database` |
| `DATABASE_URL_DEDICATED` | Yes | — | PostgreSQL connection string used for the dedicated connection pool (migrations and long-running queries). May be the same as `DATABASE_URL`. |
| `DATABASE_MIGRATION_SKIP` | No | — | When set to any non-empty value, skips running GORM AutoMigrate on startup. Useful when migrations are managed externally. |

### Redis

| Variable | Required | Default | Description |
|---|---|---|---|
| `REDIS_URL` | Yes | — | Redis connection string. Format: `redis://[:password@]host:port[/database]` |

### Firebase

| Variable | Required | Default | Description |
|---|---|---|---|
| `GCP_PROJECT_ID` | Yes | — | Firebase project ID. Also used as the Google Cloud project identifier. |
| `FIREBASE_CREDENTIALS` | Yes | — | Full content of the Firebase service account JSON file, provided as a single-line string. |

### Event Queue (System User)

These variables configure the internal system user that drives the asynchronous event queue.

| Variable | Required | Default | Description |
|---|---|---|---|
| `EVENTS_QUEUE_TYPE` | No | `emulator` | Queue backend. `emulator` uses an in-process HTTP emulator suitable for Docker deployments. Set to `cloud_tasks` for Google Cloud Tasks. |
| `EVENTS_QUEUE_NAME` | No | `events-local` | Name of the queue used by the emulator or Cloud Tasks. |
| `EVENTS_QUEUE_ENDPOINT` | No | `http://localhost:8000/v1/events` | HTTP endpoint where the queue delivers event payloads. Update to the API's internal hostname in multi-container deployments. |
| `EVENTS_QUEUE_USER_ID` | No | Auto-generated | UUID of the system user. Auto-generated on first boot if empty; printed to stdout. Persist this value across restarts. |
| `EVENTS_QUEUE_USER_API_KEY` | No | Auto-generated | API key of the system user. Auto-generated on first boot if empty; printed to stdout. Persist this value across restarts. |
| `EVENTS_QUEUE_USER_EMAIL` | No | `system@httpsms.local` | Email address of the system user record in the database. |

> **Note:** On the first startup against an empty database, the API automatically creates the system user using the values above. If `EVENTS_QUEUE_USER_ID` or `EVENTS_QUEUE_USER_API_KEY` are not provided, cryptographically secure values are generated and printed to the container logs. Capture and persist these values before subsequent restarts. See [self-hosting.md — Step 7](./self-hosting.md#step-7--system-user-initialization) for details.

### SMTP

| Variable | Required | Default | Description |
|---|---|---|---|
| `SMTP_HOST` | Yes | — | Hostname of the SMTP server. |
| `SMTP_PORT` | No | `587` | Port of the SMTP server. |
| `SMTP_USERNAME` | Yes | — | SMTP authentication username. |
| `SMTP_PASSWORD` | Yes | — | SMTP authentication password. |
| `SMTP_FROM_NAME` | No | `httpSMS` | Display name shown in the From field of outbound emails. |
| `SMTP_FROM_EMAIL` | No | `noreply@httpsms.local` | Email address used in the From field of outbound emails. Must match the domain authorized by your SMTP provider. |

### Cloudflare Turnstile

| Variable | Required | Default | Description |
|---|---|---|---|
| `CLOUDFLARE_TURNSTILE_SECRET_KEY` | No | — | Turnstile secret key used by the API to validate CAPTCHA tokens on the `/v1/messages/search` route. The route is unprotected if this value is empty. |

### Entitlements

| Variable | Required | Default | Description |
|---|---|---|---|
| `ENTITLEMENT_ENABLED` | No | `false` | When `true`, enforces subscription-based message limits per user. Set to `false` for self-hosted deployments where all users should have unrestricted access. |

### Optional Integrations

| Variable | Required | Default | Description |
|---|---|---|---|
| `PUSHER_APP_ID` | No | — | Pusher application ID. Enables real-time updates in the web UI via WebSockets. |
| `PUSHER_KEY` | No | — | Pusher public key. |
| `PUSHER_SECRET` | No | — | Pusher secret key. |
| `PUSHER_CLUSTER` | No | — | Pusher cluster identifier (e.g., `mt1`, `eu`). |
| `AXIOM_TOKEN` | No | — | Axiom API token for centralized log and trace ingestion. |
| `AXIOM_DATASET_EVENTS` | No | — | Axiom dataset name for logs and traces. |
| `AXIOM_DATASET_METRICS` | No | — | Axiom dataset name for metrics. |
| `GCS_BUCKET_NAME` | No | — | Google Cloud Storage bucket name for MMS attachment storage. When empty, attachments are stored in memory and lost on restart. |

---

## Web Frontend Environment Variables

The web frontend is a statically generated Nuxt SPA. All environment variables are **injected at image build time** via Docker build arguments. They are not read from the container environment at runtime.

When using Docker Compose, the `docker-compose.yml` passes these values from the host environment to the build process automatically. When building the image manually or in a CI pipeline, pass them as `--build-arg` flags.

| Variable | Required | Default | Description |
|---|---|---|---|
| `API_BASE_URL` | Yes | `http://localhost:8000` | Public URL of the API. The browser makes requests to this URL directly. Must be reachable from the end user's browser, not just from within the Docker network. |
| `APP_URL` | No | `http://localhost:3000` | Public URL of the web application. Used for SEO metadata and email links. |
| `APP_NAME` | No | `httpSMS` | Display name of the application. |
| `APP_ENV` | No | `production` | Runtime environment label. |
| `APP_GITHUB_URL` | No | `https://github.com/NdoleStudio/httpsms` | Link to the project source repository shown in the UI. |
| `APP_DOCUMENTATION_URL` | No | `https://docs.httpsms.com` | Link to external documentation shown in the UI. |
| `APP_DOWNLOAD_URL` | No | (GitHub releases URL) | Direct download URL for the Android APK. |
| `FIREBASE_API_KEY` | Yes | — | Firebase Web SDK API key. |
| `FIREBASE_AUTH_DOMAIN` | Yes | — | Firebase authentication domain. |
| `FIREBASE_PROJECT_ID` | Yes | — | Firebase project ID. |
| `FIREBASE_STORAGE_BUCKET` | Yes | — | Firebase Storage bucket name. |
| `FIREBASE_MESSAGING_SENDER_ID` | Yes | — | Firebase Cloud Messaging sender ID. |
| `FIREBASE_APP_ID` | Yes | — | Firebase App ID. |
| `FIREBASE_MEASUREMENT_ID` | No | — | Firebase Analytics measurement ID. Optional. |
| `CLOUDFLARE_TURNSTILE_SITE_KEY` | No | — | Turnstile site key embedded in the search page form. |
| `PUSHER_KEY` | No | — | Pusher public key for real-time updates. |
| `PUSHER_CLUSTER` | No | — | Pusher cluster for real-time updates. |

---

## PostgreSQL Container Variables

These variables configure the `postgres` service in the Docker Compose stack.

| Variable | Required | Default | Description |
|---|---|---|---|
| `POSTGRES_DB` | No | `httpsms` | Name of the database to create. |
| `POSTGRES_USER` | No | `httpsms` | PostgreSQL username. |
| `POSTGRES_PASSWORD` | Yes | `changeme` | PostgreSQL password. Change this before any production deployment. |

---

## Notes on Variable Precedence

- Variables defined in `.env` at the repository root are automatically loaded by Docker Compose.
- When deploying through an orchestrator such as Coolify, Portainer, or Docker Swarm, variables are injected by the platform and take precedence over any `.env` file.
- The API binary is started with `--dotenv=false`. It never reads `.env` files from disk at runtime; all configuration must arrive as process environment variables.
- The web frontend's configuration is frozen at the time the Docker image is built. Changing a web environment variable requires rebuilding the image.
