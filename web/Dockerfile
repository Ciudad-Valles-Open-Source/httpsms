# =============================================================================
# httpSMS Web – Dockerfile (Coolify-ready)
# =============================================================================
# Nuxt SPA with SSG (pnpm run generate). The final image serves static files
# via nginx on port 3000.
#
# Since Nuxt bakes environment variables at BUILD TIME, all app configuration
# must be passed as build ARGs. In Coolify, set these in the "Build Variables"
# or "Environment Variables" section.
# =============================================================================

# ── Stage 1: Build ───────────────────────────────────────────────────────────
FROM node:22-alpine AS build

WORKDIR /app

# Accept all config as build args so Nuxt can bake them into the static output
ARG API_BASE_URL=http://localhost:8000
ARG APP_URL=http://localhost:3000
ARG APP_NAME=httpSMS
ARG APP_GITHUB_URL=https://github.com/NdoleStudio/httpsms
ARG APP_DOCUMENTATION_URL=https://docs.httpsms.com
ARG APP_DOWNLOAD_URL=https://github.com/NdoleStudio/httpsms/releases/latest/download/HttpSms.apk
ARG APP_ENV=production
ARG FIREBASE_API_KEY=
ARG FIREBASE_AUTH_DOMAIN=
ARG FIREBASE_PROJECT_ID=
ARG FIREBASE_STORAGE_BUCKET=
ARG FIREBASE_MESSAGING_SENDER_ID=
ARG FIREBASE_APP_ID=
ARG FIREBASE_MEASUREMENT_ID=
ARG CLOUDFLARE_TURNSTILE_SITE_KEY=
ARG SOKETI_KEY=
ARG SOKETI_CLUSTER=
ARG USE_SOKETI=
ARG ENABLE_BILLING=

# Promote build args to env vars so Nuxt picks them up during generate
ENV API_BASE_URL=$API_BASE_URL \
    APP_URL=$APP_URL \
    APP_NAME=$APP_NAME \
    APP_GITHUB_URL=$APP_GITHUB_URL \
    APP_DOCUMENTATION_URL=$APP_DOCUMENTATION_URL \
    APP_DOWNLOAD_URL=$APP_DOWNLOAD_URL \
    APP_ENV=$APP_ENV \
    FIREBASE_API_KEY=$FIREBASE_API_KEY \
    FIREBASE_AUTH_DOMAIN=$FIREBASE_AUTH_DOMAIN \
    FIREBASE_PROJECT_ID=$FIREBASE_PROJECT_ID \
    FIREBASE_STORAGE_BUCKET=$FIREBASE_STORAGE_BUCKET \
    FIREBASE_MESSAGING_SENDER_ID=$FIREBASE_MESSAGING_SENDER_ID \
    FIREBASE_APP_ID=$FIREBASE_APP_ID \
    FIREBASE_MEASUREMENT_ID=$FIREBASE_MEASUREMENT_ID \
    CLOUDFLARE_TURNSTILE_SITE_KEY=$CLOUDFLARE_TURNSTILE_SITE_KEY \
    SOKETI_KEY=$SOKETI_KEY \
    SOKETI_CLUSTER=$SOKETI_CLUSTER \
    USE_SOKETI=$USE_SOKETI \
    ENABLE_BILLING=$ENABLE_BILLING

# Install pnpm via corepack (built into Node 22)
RUN corepack enable && corepack prepare pnpm@latest --activate

# Install dependencies (layer cache optimization)
COPY package.json pnpm-lock.yaml pnpm-workspace.yaml ./
RUN pnpm install --frozen-lockfile

# Copy source and generate static output
COPY . .
RUN pnpm run generate

# ── Stage 2: Production (nginx) ──────────────────────────────────────────────
FROM nginx:stable-alpine AS production

# Copy the generated static site
COPY --from=build /app/.output/public /usr/share/nginx/html

# Copy the nginx configuration
COPY ./nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 3000

CMD ["nginx", "-g", "daemon off;"]
