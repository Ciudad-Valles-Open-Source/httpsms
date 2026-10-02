# Docker Configuration Reference

This document describes the Docker configuration provided in this repository, explains each design decision, and serves as a reference for operators building or modifying container images.

## Table of Contents

- [File Overview](#file-overview)
- [API Dockerfile](#api-dockerfile)
- [Web Dockerfile](#web-dockerfile)
- [docker-compose.yml](#docker-composeyml)
- [Environment Variable Handling](#environment-variable-handling)
- [Health Checks](#health-checks)
- [Build Arguments vs Runtime Variables](#build-arguments-vs-runtime-variables)

---

## File Overview

| File | Purpose |
|---|---|
| `api/Dockerfile` | Multi-stage build for the Go API binary |
| `web/Dockerfile` | Multi-stage build for the Nuxt SPA served by nginx |
| `docker-compose.yml` | Orchestrates all four services for local and production deployments |
| `.env.example` | Template listing every configurable variable with documentation |

---

## API Dockerfile

**Location:** [`api/Dockerfile`](../api/Dockerfile)

### Build stages

The API image uses a two-stage build:

**Stage 1 — Builder (`golang:1.23-alpine`)**

- Downloads Go module dependencies before copying source code, exploiting Docker's layer cache. Subsequent builds reuse the dependency layer if `go.mod` and `go.sum` have not changed.
- Installs and runs `swag` to generate the Swagger specification from Go source annotations. The generated output is embedded in the binary.
- Compiles the application with `CGO_ENABLED=0` to produce a statically linked binary. The `-s -w` linker flags strip debug symbols and DWARF information, reducing the binary size.
- The `GIT_COMMIT` build argument is injected into the binary at compile time via `-X main.Version`, making the deployed version identifiable in the Swagger UI and logs.

**Stage 2 — Runtime (`alpine:3.20`)**

- Creates a non-root user and group (`httpsms`) for the process.
- Copies only the compiled binary and the `root.crt` certificate file from the builder stage.
- Installs `curl` and `ca-certificates` for health checks and TLS verification respectively, and `tzdata` for timezone support.
- Runs the binary with `--dotenv=false`, which disables `.env` file loading. All configuration must be provided via container environment variables.

### Build arguments

| Argument | Default | Description |
|---|---|---|
| `GIT_COMMIT` | `dev` | Git commit SHA injected into the binary version string. |

---

## Web Dockerfile

**Location:** [`web/Dockerfile`](../web/Dockerfile)

### Build stages

**Stage 1 — Builder (`node:22-alpine`)**

- Uses `corepack` (built into Node 22) to activate `pnpm` without a separate global install step.
- Installs dependencies with `--frozen-lockfile` to ensure reproducible builds. The build fails if `pnpm-lock.yaml` is out of sync with `package.json`.
- Accepts all application configuration as Docker build arguments (`ARG`), then promotes them to environment variables (`ENV`) so that the Nuxt build process can read them.
- Runs `pnpm run generate`, which executes `nuxi generate`. This produces a fully static HTML/CSS/JS output in `.output/public`.

**Stage 2 — Runtime (`nginx:stable-alpine`)**

- Copies the static output from the builder stage into the nginx document root.
- Uses the project-provided `nginx.conf`, which configures a single server block on port 3000 with `try_files` routing to support client-side navigation in the SPA.

### Why build arguments are required

The Nuxt configuration reads environment variables at build time using `process.env`. Because the output is fully static, there is no server process at runtime that could inject configuration. Every configuration value that must be available in the browser must be present when `pnpm run generate` runs.

A consequence of this design is that **changing any web environment variable requires rebuilding the Docker image**. There is no way to apply web configuration changes by restarting the container.

### Build arguments

All build arguments correspond directly to the environment variables described in [environment-variables.md](./environment-variables.md#web-frontend-environment-variables). They default to empty strings when not provided, which may result in features such as Firebase authentication being non-functional.

---

## docker-compose.yml

**Location:** [`docker-compose.yml`](../docker-compose.yml)

### Design principles

- **All configuration is provided via environment variables.** No `env_file` directives reference files that would be absent in a repository clone or a CI/CD runner. This makes the file compatible with orchestration platforms such as Coolify, Portainer, and Docker Swarm, which inject configuration from their own variable stores.
- **Default values are provided for non-sensitive variables.** Variables use the `${VAR:-default}` syntax. This allows the stack to start for local development without requiring every variable to be set, while still permitting full configuration override in production.
- **Health checks are defined for all services.** Dependent services wait for their dependencies to become healthy before starting, reducing startup race conditions. The `api` service waits for PostgreSQL to accept connections; the `web` service waits for the `api` service to become healthy.

### Services

#### postgres

Runs PostgreSQL 16 with a named volume for data persistence. The health check uses `pg_isready` to verify the database is accepting connections before dependent services start.

#### redis

Runs Redis 7 with append-only persistence enabled (`--appendonly yes`), which survives container restarts. The health check uses `redis-cli ping`.

#### api

Builds from `api/Dockerfile`. Receives all configuration as environment variables. The `DATABASE_URL` and `REDIS_URL` values are constructed from the PostgreSQL and Redis variable sets to ensure they stay in sync with the database service configuration.

The `EVENTS_QUEUE_ENDPOINT` variable defaults to `http://api:8000/v1/events`, using the Docker network service name `api` so that the event queue can reach the API from within the compose network.

#### web

Builds from `web/Dockerfile`. All configuration is passed as build arguments that the Dockerfile accepts via `ARG` declarations. The `web` service depends on the `api` service being healthy to avoid starting before the API is ready to serve requests.

### Volumes

Two named volumes are defined:

- `postgres_data` — persists the PostgreSQL data directory across container restarts and upgrades.
- `redis_data` — persists the Redis append-only log.

---

## Environment Variable Handling

### API (runtime)

The Go API server reads all configuration from process environment variables at startup. The `--dotenv=false` flag passed to the binary ensures it does not attempt to load a `.env` file from disk. This is the expected behavior in containerized deployments where the orchestrator is responsible for injecting environment variables.

### Web (build time)

The Nuxt `runtimeConfig.public` block in `nuxt.config.ts` reads environment variables via `process.env` at the time `nuxi generate` runs. The resulting values are embedded into the generated JavaScript bundle. This means:

- Values are visible in the browser (treat them as public).
- Secret values must never be passed as web build arguments.
- Configuration changes require a new image build and redeployment.

---

## Health Checks

Every service in `docker-compose.yml` defines a health check that Docker uses to determine readiness:

| Service | Health Check Method | Healthy Condition |
|---|---|---|
| `postgres` | `pg_isready` | Database accepting connections |
| `redis` | `redis-cli ping` | Returns `PONG` |
| `api` | `curl -f http://localhost:8000/` | HTTP 200 response |
| `web` | `wget --spider http://localhost:3000/` | HTTP 200 response |

The dependency chain is: `postgres` and `redis` must be healthy before `api` starts; `api` must be healthy before `web` starts.

---

## Build Arguments vs Runtime Variables

The following table summarizes which variables are consumed at build time versus runtime for each service:

| Variable type | API | Web |
|---|---|---|
| Database credentials | Runtime | Not applicable |
| Firebase service account | Runtime | Not applicable |
| Firebase Web SDK config | Not applicable | **Build time** |
| SMTP credentials | Runtime | Not applicable |
| API base URL | Runtime (Swagger host) | **Build time** |
| Application URL | Runtime | **Build time** |
| Pusher configuration | Runtime | **Build time** |

When reconfiguring a deployment, identify whether the changed variable is consumed at build time or runtime to determine whether a full rebuild or a simple restart is sufficient.
