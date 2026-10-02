# Self-Hosting httpSMS

This document covers everything required to run httpSMS on your own infrastructure. The application is distributed as a set of Docker images and can be deployed on any container platform that supports Docker Compose or individual Dockerfiles.

## Table of Contents

- [Architecture Overview](#architecture-overview)
- [Prerequisites](#prerequisites)
- [Step 1 — Configure Firebase](#step-1--configure-firebase)
- [Step 2 — Configure SMTP](#step-2--configure-smtp)
- [Step 3 — Configure Cloudflare Turnstile](#step-3--configure-cloudflare-turnstile)
- [Step 4 — Clone the Repository](#step-4--clone-the-repository)
- [Step 5 — Set Environment Variables](#step-5--set-environment-variables)
- [Step 6 — Build and Start](#step-6--build-and-start)
- [Step 7 — System User Initialization](#step-7--system-user-initialization)
- [Step 8 — Build the Android App](#step-8--build-the-android-app)
- [Post-Deployment Verification](#post-deployment-verification)

---

## Architecture Overview

A complete httpSMS deployment consists of four services:

| Service | Image | Default Port | Purpose |
|---|---|---|---|
| `api` | Built from `api/Dockerfile` | 8000 | Go/Fiber REST API and event processor |
| `web` | Built from `web/Dockerfile` | 3000 | Nuxt SPA (served by nginx) |
| `postgres` | `postgres:16-alpine` | 5432 | Primary relational database |
| `redis` | `redis:7-alpine` | 6379 | Queue and cache backend |

The API runs database migrations automatically on startup using GORM AutoMigrate. No manual schema management is required.

---

## Prerequisites

- Docker Engine 24+ and Docker Compose v2
- A [Firebase](https://console.firebase.google.com/) project (free tier is sufficient)
- An SMTP provider (any standard SMTP server or service such as Mailgun, Postmark, Mailtrap, or a self-hosted solution)
- A [Cloudflare Turnstile](https://developers.cloudflare.com/turnstile/get-started/) widget (required for the message search feature)
- An Android phone with the httpSMS Android app installed

---

## Step 1 — Configure Firebase

httpSMS uses Firebase for two purposes:

- **Firebase Cloud Messaging (FCM):** Sends push notifications to the Android app to trigger outbound SMS delivery.
- **Firebase Authentication:** Authenticates users on the web interface.

### Create a Firebase project

1. Open the [Firebase Console](https://console.firebase.google.com/) and create a new project.
2. Register a **Web App** within the project. Follow the [Web App setup guide](https://firebase.google.com/docs/web/setup#register-app) and record the SDK configuration values:

```js
const firebaseConfig = {
  apiKey: "AIzaSy...",
  authDomain: "your-project.firebaseapp.com",
  projectId: "your-project",
  storageBucket: "your-project.appspot.com",
  messagingSenderId: "123456789",
  appId: "1:123456789:web:abc123",
  measurementId: "G-XXXXX",
};
```

These values are used by the web frontend at build time.

### Enable Email/Password sign-in

1. In the Firebase Console, navigate to **Authentication > Sign-in method**.
2. Enable **Email/Password** and click **Save**.

> **Known Issue:** The Firebase email/password sign-in provider has a bug with email enumeration protection that can prevent login. If users cannot sign in, [disable email enumeration protection](https://cloud.google.com/identity-platform/docs/admin/email-enumeration-protection#disable) in the Identity Platform settings.

### Generate a service account key

1. Follow the [Firebase Admin SDK setup guide](https://firebase.google.com/docs/admin/setup#initialize_the_sdk_in_non-google_environments).
2. Download the service account JSON file.
3. Convert it to a single-line string for use in the `FIREBASE_CREDENTIALS` environment variable:

```bash
cat firebase-credentials.json | tr -d '\n'
```

### Generate google-services.json for the Android app

Follow the [Android setup instructions](https://support.google.com/firebase/answer/7015592?hl=en#android) to download `google-services.json`. This file is needed when building the Android app in Step 8.

---

## Step 2 — Configure SMTP

httpSMS sends transactional emails for events such as when an Android phone has been offline for an extended period.

Any SMTP server is supported. For local development, [Mailtrap](https://mailtrap.io/) provides a free sandbox. For production, services such as Mailgun, Postmark, or SendGrid are suitable.

Record the following values for use in environment variables:

- `SMTP_HOST`
- `SMTP_PORT` (typically `587` for TLS or `465` for SSL)
- `SMTP_USERNAME`
- `SMTP_PASSWORD`
- `SMTP_FROM_EMAIL` (the sender address)

---

## Step 3 — Configure Cloudflare Turnstile

The `/v1/messages/search` route is protected by a [Cloudflare Turnstile](https://developers.cloudflare.com/turnstile/get-started/) CAPTCHA to prevent abuse.

1. Open the [Cloudflare dashboard](https://dash.cloudflare.com/) and navigate to **Turnstile**.
2. Create a new site widget. For local development, use `localhost` as the domain.
3. Record the **Site Key** (`CLOUDFLARE_TURNSTILE_SITE_KEY`, used by the web frontend) and the **Secret Key** (`CLOUDFLARE_TURNSTILE_SECRET_KEY`, used by the API).

---

## Step 4 — Clone the Repository

```bash
git clone https://github.com/NdoleStudio/httpsms.git
cd httpsms
```

---

## Step 5 — Set Environment Variables

Copy the provided template and fill in all required values:

```bash
cp .env.example .env
```

Open `.env` and configure the following required variables:

```dotenv
# Database
POSTGRES_PASSWORD=<a-strong-random-password>

# Firebase (API — server-side)
GCP_PROJECT_ID=<your-firebase-project-id>
FIREBASE_CREDENTIALS=<single-line-json-from-step-1>

# Firebase (Web — injected at build time)
FIREBASE_API_KEY=
FIREBASE_AUTH_DOMAIN=
FIREBASE_PROJECT_ID=
FIREBASE_STORAGE_BUCKET=
FIREBASE_MESSAGING_SENDER_ID=
FIREBASE_APP_ID=
FIREBASE_MEASUREMENT_ID=

# SMTP
SMTP_HOST=
SMTP_PORT=587
SMTP_USERNAME=
SMTP_PASSWORD=
SMTP_FROM_EMAIL=

# Cloudflare Turnstile
CLOUDFLARE_TURNSTILE_SECRET_KEY=
CLOUDFLARE_TURNSTILE_SITE_KEY=

# Application URLs (update to your domain for production)
APP_URL=http://localhost:3000
API_BASE_URL=http://localhost:8000
```

For a complete reference of all available variables, see [environment-variables.md](./environment-variables.md).

---

## Step 6 — Build and Start

```bash
docker compose up --build
```

The first build downloads base images and compiles the Go binary; this may take several minutes. When complete:

- Web UI: [http://localhost:3000](http://localhost:3000)
- API: [http://localhost:8000](http://localhost:8000)
- Swagger UI: [http://localhost:8000/index.html](http://localhost:8000/index.html)

---

## Step 7 — System User Initialization

httpSMS uses a dedicated internal system user to process asynchronous events through the event queue. This user does not correspond to any human operator.

**This step is fully automatic.** When the API starts against an empty database, it detects the absence of users and creates the system user automatically from the values of `EVENTS_QUEUE_USER_ID` and `EVENTS_QUEUE_USER_API_KEY`. If those variables are not set, secure values are generated and printed to the API container logs:

```
╔══════════════════════════════════════════════════════════════╗
║              SYSTEM USER CREATED SUCCESSFULLY              ║
╠══════════════════════════════════════════════════════════════╣
║  User ID : a1b2c3d4-e5f6-7890-abcd-ef1234567890            ║
║  API Key : uk_Abc123...                                     ║
║  Email   : system@httpsms.local                             ║
╠══════════════════════════════════════════════════════════════╣
║  Save these values in your environment variables:          ║
║    EVENTS_QUEUE_USER_ID                                    ║
║    EVENTS_QUEUE_USER_API_KEY                               ║
║                                                            ║
║  If auto-generated, set them before the next restart       ║
║  to avoid creating a new system user.                      ║
╚══════════════════════════════════════════════════════════════╝
```

**After the first boot**, retrieve these values from the logs and persist them in your environment configuration:

```dotenv
EVENTS_QUEUE_USER_ID=<value-from-logs>
EVENTS_QUEUE_USER_API_KEY=<value-from-logs>
```

Restart the API service to pick up the changes:

```bash
docker compose restart api
```

If you prefer to pre-define these values before the first boot, set them in your `.env` file beforehand. The seed process will use the provided values instead of auto-generating them.

---

## Step 8 — Build the Android App

Before building the Android app, replace `android/app/google-services.json` with the file downloaded in Step 1. This enables FCM push notifications from your own Firebase project.

Open the project in [Android Studio](https://developer.android.com/studio) and build a release APK. Install it on the Android phone that will function as your SMS gateway.

---

## Post-Deployment Verification

1. Open the web UI and create a user account via the registration form.
2. Navigate to **Settings** and verify your API key is displayed.
3. Open the Android app, log in with the same credentials, and confirm the phone appears in the dashboard.
4. Send a test SMS message via the API or the web UI.

For platform-specific deployment instructions (Coolify, Portainer, CapRover, Docker Swarm), see [deployment-platforms.md](./deployment-platforms.md).
