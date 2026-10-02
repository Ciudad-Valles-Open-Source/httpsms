# =============================================================================
# httpSMS API – Dockerfile (Coolify-ready)
# =============================================================================
# Multi-stage build: compiles the Go binary with Swagger docs, then runs in a
# minimal Alpine container. All runtime configuration is via environment
# variables (set in Coolify's UI or your .env file).
# =============================================================================

# ── Stage 1: Build ───────────────────────────────────────────────────────────
FROM golang:1.26-alpine AS builder

ARG GIT_COMMIT=dev
ENV GIT_COMMIT=$GIT_COMMIT

# Install git (needed by swag for parsing)
RUN apk add --no-cache git

WORKDIR /http-sms

# Cache dependency downloads
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY . .

# Generate Swagger docs
RUN go install github.com/swaggo/swag/cmd/swag@latest \
    && swag init --requiredByDefault --parseDependency --parseInternal

# Build a statically linked binary
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags "-s -w -X main.Version=$GIT_COMMIT" \
    -o /bin/http-sms .

# ── Stage 2: Production ─────────────────────────────────────────────────────
FROM alpine:3.20

RUN apk add --no-cache curl ca-certificates tzdata \
    && addgroup -S httpsms \
    && adduser -S httpsms -G httpsms

USER httpsms
WORKDIR /home/httpsms

COPY --from=builder /bin/http-sms ./
COPY --from=builder /http-sms/root.crt ./

EXPOSE 8000

# --dotenv=false tells the app to read config from environment variables
# (not from a .env file), which is exactly what Coolify provides.
ENTRYPOINT ["./http-sms", "--dotenv=false"]
