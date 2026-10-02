# Deployment Platforms

This document provides platform-specific instructions for deploying httpSMS on common container orchestration platforms. All platforms use the same Docker images and the same set of environment variables described in [environment-variables.md](./environment-variables.md).

## Table of Contents

- [Coolify](#coolify)
- [Portainer](#portainer)
- [CapRover](#caprover)
- [Docker Swarm](#docker-swarm)
- [Generic Docker Compose](#generic-docker-compose)
- [Important Notes on the Web Frontend](#important-notes-on-the-web-frontend)

---

## Coolify

[Coolify](https://coolify.io/) is a self-hosted Heroku/Netlify alternative. It supports Git-based deployments of Docker Compose stacks.

### Docker Compose deployment (recommended)

This approach deploys all four services (API, web, PostgreSQL, Redis) as a single stack.

1. In Coolify, navigate to **New Resource > Docker Compose**.
2. Connect your Git repository and select the branch to deploy.
3. Coolify will detect `docker-compose.yml` automatically.
4. Navigate to the **Environment Variables** panel and set all required variables. Refer to [environment-variables.md](./environment-variables.md) for the full list.

**Minimum required variables for a functional deployment:**

```dotenv
POSTGRES_PASSWORD=<strong-password>

GCP_PROJECT_ID=<firebase-project-id>
FIREBASE_CREDENTIALS=<single-line-service-account-json>

FIREBASE_API_KEY=<web-sdk-value>
FIREBASE_AUTH_DOMAIN=<web-sdk-value>
FIREBASE_PROJECT_ID=<web-sdk-value>
FIREBASE_STORAGE_BUCKET=<web-sdk-value>
FIREBASE_MESSAGING_SENDER_ID=<web-sdk-value>
FIREBASE_APP_ID=<web-sdk-value>

APP_URL=https://your-domain.com
API_BASE_URL=https://api.your-domain.com
SWAGGER_HOST=api.your-domain.com

SMTP_HOST=<smtp-host>
SMTP_USERNAME=<smtp-user>
SMTP_PASSWORD=<smtp-password>
SMTP_FROM_EMAIL=<sender-address>

CLOUDFLARE_TURNSTILE_SECRET_KEY=<secret>
CLOUDFLARE_TURNSTILE_SITE_KEY=<site-key>
```

5. In Coolify's service settings, assign public domains to the `web` (port 3000) and `api` (port 8000) services. Coolify provisions TLS automatically via Traefik.
6. Click **Deploy**. The system user is created automatically on the first boot. Check the `api` service logs to retrieve the generated credentials if `EVENTS_QUEUE_USER_ID` and `EVENTS_QUEUE_USER_API_KEY` were not set in advance.

### Individual Dockerfile deployment

Use this approach to deploy the API and web frontend as separate Coolify applications, with PostgreSQL and Redis managed as Coolify database resources.

**API application:**

- Resource type: Application (Dockerfile)
- Dockerfile path: `api/Dockerfile`
- Build context: `api/`
- Port: `8000`

**Web application:**

- Resource type: Application (Dockerfile)
- Dockerfile path: `web/Dockerfile`
- Build context: `web/`
- Port: `3000`

Set environment variables for each application separately. Note that web frontend variables are build-time variables and must be configured as **Build Variables** (not runtime environment variables) in Coolify's application settings.

---

## Portainer

[Portainer](https://www.portainer.io/) is a container management UI for Docker and Kubernetes.

### Docker Compose stack

1. In Portainer, navigate to **Stacks > Add Stack**.
2. Select **Repository** and connect your Git repository.
3. Set the Compose file path to `docker-compose.yml`.
4. Under **Environment variables**, add all required variables.
5. Click **Deploy the stack**.

### Managing stacks

To update the deployment after pushing new commits:

1. Navigate to the stack in Portainer.
2. Click **Pull and redeploy**.

---

## CapRover

[CapRover](https://caprover.com/) provides a simple PaaS layer on top of Docker Swarm.

### Setup

Create the following applications in CapRover before deploying:

1. **Database services** — use CapRover's One-Click Apps to create PostgreSQL and Redis instances.
2. **API application** — create a new app named `httpsms-api`.
3. **Web application** — create a new app named `httpsms-web`.

### Deploy the API

1. In the `httpsms-api` app settings, go to **Deployment** and connect your Git repository.
2. Set the **Dockerfile path** to `api/Dockerfile` and the **build context** to `api/`.
3. Under **App Configs**, set all API environment variables. Use the internal service names provided by CapRover for `DATABASE_URL` and `REDIS_URL`.
4. Enable HTTPS under the **HTTP Settings** tab.
5. Deploy.

### Deploy the web frontend

1. In the `httpsms-web` app settings, connect the same repository.
2. Set the **Dockerfile path** to `web/Dockerfile` and the **build context** to `web/`.
3. Under **App Configs**, set all web build arguments as environment variables. CapRover passes them to the Docker build process as build arguments when the variable names match the `ARG` declarations in the Dockerfile.
4. Enable HTTPS and deploy.

---

## Docker Swarm

Docker Swarm is suitable for multi-node production deployments.

### Deploy as a stack

```bash
# Initialize Swarm (if not already active)
docker swarm init

# Deploy the stack from the repository root
docker stack deploy --compose-file docker-compose.yml httpsms
```

Environment variables can be provided by exporting them in the shell before running the deploy command, or by sourcing a `.env` file:

```bash
set -a && source .env && set +a
docker stack deploy --compose-file docker-compose.yml httpsms
```

### Managing secrets

For production Swarm deployments, use Docker Secrets instead of plain environment variables for sensitive values. Create secrets for `POSTGRES_PASSWORD` and `FIREBASE_CREDENTIALS`:

```bash
echo "strong-password" | docker secret create postgres_password -
cat firebase-credentials.json | tr -d '\n' | docker secret create firebase_credentials -
```

Then reference them in the Compose file using the `secrets:` key. This requires modifying `docker-compose.yml` to mount secrets as files and reading them from the API code, which is not configured by default.

### Updating a running stack

```bash
docker stack deploy --compose-file docker-compose.yml httpsms
```

Docker Swarm performs a rolling update automatically.

---

## Generic Docker Compose

For any Linux host with Docker Engine installed:

```bash
# Clone the repository
git clone https://github.com/NdoleStudio/httpsms.git
cd httpsms

# Create and configure the environment file
cp .env.example .env
# Edit .env with your values

# Build and start all services in the background
docker compose up -d --build

# View logs
docker compose logs -f api

# Stop all services
docker compose down
```

To update after pulling new code:

```bash
git pull
docker compose up -d --build
```

---

## Important Notes on the Web Frontend

The Nuxt web application is a statically generated Single Page Application (SPA). All configuration, including the API base URL and Firebase credentials, is embedded into the JavaScript bundle at image build time.

This has two practical consequences:

1. **Changing a web environment variable requires rebuilding the image.** Restarting the container without rebuilding does not apply configuration changes to the web frontend.

2. **`API_BASE_URL` must be the public URL of the API**, reachable from the end user's browser. Setting it to an internal Docker network address (such as `http://api:8000`) will cause all browser requests to fail.

When using Docker Compose locally, set `API_BASE_URL=http://localhost:8000`. In production, set it to the HTTPS URL of your API (for example, `https://api.your-domain.com`).
