# Security Policy

## Supported Versions

Only the latest version of httpSMS on the `main` branch is actively maintained. Security fixes are applied to `main` and are not backported to older tags or branches.

| Branch / Version | Supported |
|---|---|
| `main` | Yes |
| Any tagged release older than `main` | No |

Users running a self-hosted deployment are responsible for keeping their installation up to date by pulling the latest changes from `main` and rebuilding their Docker images.

## Reporting a Vulnerability

Security vulnerabilities must be reported privately. Do not open a public GitHub issue for a security problem, as this may expose users before a fix is available.

Report vulnerabilities using any of the following channels:

- **Email:** arnold@httpsms.com
- **Discord:** Send a private message on the [httpSMS Discord server](https://discord.gg/kGk8HVqeEZ)

Include as much detail as possible in your report:

- A description of the vulnerability and the potential impact
- Steps to reproduce the issue or a proof-of-concept
- The version or commit hash you tested against
- Any suggested mitigation or fix, if applicable

All reports are reviewed promptly. You will receive an acknowledgement within 48 hours. If the vulnerability is confirmed, a fix will be developed and released as quickly as possible, and you will be credited in the release notes unless you prefer to remain anonymous.

## Security Considerations for Self-Hosted Deployments

If you are running httpSMS on your own infrastructure, review the following security guidance:

### Environment variables

- The `.env.example` file at the repository root contains placeholder values. Never commit real credentials to version control.
- `POSTGRES_PASSWORD` must be changed from the default before any internet-facing deployment.
- `FIREBASE_CREDENTIALS` contains the full Firebase service account JSON. Treat this as a high-value secret and restrict access accordingly.
- `EVENTS_QUEUE_USER_API_KEY` is the API key used internally by the event queue. If auto-generated on first boot, retrieve it from the container logs immediately and persist it in your secret store.

### Network exposure

- The PostgreSQL and Redis services defined in `docker-compose.yml` expose ports to the host for local development convenience. In production, remove the `ports:` entries for `postgres` and `redis` so they are only accessible within the Docker network.
- Place the web frontend and API behind a TLS-terminating reverse proxy. Platforms such as Coolify and CapRover handle this automatically via Traefik.

### Firebase Authentication

- Enable only the sign-in providers you intend to use. The application requires Email/Password authentication.
- Review the [Firebase email enumeration protection](https://cloud.google.com/identity-platform/docs/admin/email-enumeration-protection) settings for your project.

### API keys

- User API keys are prefixed with `uk_` and are stored in the database. They grant full access to the user's account. Users should be advised to treat their API key as a password.
- The system user API key (`EVENTS_QUEUE_USER_API_KEY`) grants access to the `/v1/events` internal endpoint. Do not expose this key to end users.
