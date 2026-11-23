# Real-time updates (Soketi)

The web UI updates itself when an SMS is received, sent or fails, without reloading the page. This is powered by [Soketi](https://soketi.app/), a self-hosted WebSocket server compatible with the Pusher protocol.

> **Read [⚠️ `SOKETI_KEY` must be alphanumeric](#️-soketi_key-must-be-alphanumeric) before configuring it.** It is the most common reason real-time updates do not work.

---

## How it works

```
Android phone ──▶ API (Go)
                   │  dispatches the event (phone.updated, message.phone.received, …)
                   ▼
              WebsocketListener ──HTTP──▶ Soketi (soketi:6001, internal Docker network)
                                             │
Browser ◀──WebSocket── nginx (web:3000, /app/ route) ◀─┘
(pusher-js)
```

The API publishes each event on a channel named after the **user ID**. The browser connects to `wss://<your-domain>/app/<SOKETI_KEY>` and nginx (`web` container) proxies `/app/` to Soketi. The web app opens **a single shared WebSocket connection** (`web/app/composables/useRealtime.ts`) and refreshes the thread list, the open thread, phones and API keys when `message.phone.received`, `message.phone.sent`, `message.send.failed`, `message.call.missed` or `phone.updated` arrive.

---

## Environment variables

| Variable         | Service(s)       | Description                                                                            |
| ---------------- | ---------------- | -------------------------------------------------------------------------------------- |
| `SOKETI_APP_ID`  | api, soketi      | Soketi application ID.                                                                 |
| `SOKETI_KEY`     | api, soketi, web | Public key. **Letters and numbers only** (see below).                                  |
| `SOKETI_SECRET`  | api, soketi      | Secret used by the API to sign events. When empty, the API publishes nothing.          |
| `SOKETI_HOST`    | api              | Where the API reaches Soketi. Defaults to `soketi:6001`.                               |
| `SOKETI_SECURE`  | api              | `false` when the API talks to Soketi over HTTP inside Docker.                          |
| `SOKETI_CLUSTER` | api, web         | Only used with hosted Pusher. Can stay `mt1` with Soketi.                              |
| `USE_SOKETI`     | web              | `true` to connect to the same domain (`/app/`) instead of hosted Pusher.               |

> **The web app bakes `SOKETI_KEY` in at build time.** After changing it, **rebuild the `web` container** (`docker compose up -d --build web`, or redeploy in Coolify).

---

## ⚠️ `SOKETI_KEY` must be alphanumeric

Pusher builds the WebSocket URL by appending the key as-is: `wss://your-domain/app/<SOKETI_KEY>`. A **base64** key may contain `/` (splits the path, Soketi answers **404**) and `+` (decoded as a space, so the key no longer matches).

**Symptoms:** nothing updates until you reload; nginx logs show `GET /app/… 404` instead of `101`; the browser console shows `SOKETI_KEY must only contain letters, numbers, "-" or "_"…`.

**Fix:** generate hex keys, update `.env` and Coolify, then rebuild and redeploy every service:

```bash
openssl rand -hex 32
```

---

## Also fixed: `PUSHER_SECRET`

The API only enabled its WebSocket listeners when `PUSHER_SECRET` was set, a variable from the upstream project. This fork uses `SOKETI_SECRET`, so the API never published events. `api/pkg/di/container.go` now accepts `SOKETI_SECRET` (and still `PUSHER_SECRET`).

---

## Verifying

1. API logs on startup show `registering listeners for listeners.WebsocketListener` and `creating *pusher.Client`.
2. Browser DevTools → Network → WS: a connection to `/app/<key>` with status **101** that stays open.
3. Send an SMS to the Android phone: the thread list updates within a second or two without reloading.
