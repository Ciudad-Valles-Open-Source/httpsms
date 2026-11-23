# Webhooks

A webhook is a URL in **your** application that httpSMS sends a `POST` request to every time something happens: the phone receives an SMS, a message gets delivered, the phone goes offline, and so on.

Why use them? The alternative is asking the API every few seconds whether anything changed (polling), which is slow and wastes resources. With webhooks you find out instantly, and only when something actually happens.

## What you need to create a webhook

| Field                               | Required                     | Description                                                                                                                                                                                             |
| ----------------------------------- | ---------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| **Callback URL** (`url`)            | Yes                          | The address on your server that will receive the `POST` requests. Up to 255 characters. If you leave out the scheme, `https://` is added.                                                               |
| **Signing key** (`signing_key`)     | No, but strongly recommended | A secret shared only between httpSMS and your server. It is used to sign every request so you can check that it really comes from your httpSMS server and not from a third party. Up to 255 characters. |
| **Events** (`events`)               | Yes                          | The event types you want to receive (at least one). See the [event list](/docs/webhooks-events).                                                                                                        |
| **Phone numbers** (`phone_numbers`) | Yes                          | Which phones on your account you want events from. If you have several, you can pick just some of them.                                                                                                 |

The available events are:

| Event                     | When it is sent                                                            |
| ------------------------- | -------------------------------------------------------------------------- |
| `message.phone.received`  | The phone received an SMS.                                                 |
| `message.phone.sent`      | The phone sent an SMS.                                                     |
| `message.phone.delivered` | The SMS reached the recipient's phone.                                     |
| `message.send.failed`     | The phone could not send the SMS.                                          |
| `message.send.expired`    | The phone did not send the SMS in time (for example, it was switched off). |
| `message.call.missed`     | The phone received a missed call.                                          |
| `phone.heartbeat.offline` | The server has not heard from the phone for about an hour.                 |
| `phone.heartbeat.online`  | The phone came back online after being offline.                            |

The details and payload of each event are on [Webhook events](/docs/webhooks-events).

> **Note:** this self-hosted edition has no limit on webhooks. Create as many as you need, for example one for your CRM and another for your alerting system. If several webhooks subscribe to the same event and phone, all of them receive it.

## Create a webhook from the dashboard

1. Open **Settings** at `{{APP_URL}}/settings` and scroll down to the **Webhooks** section.
2. Click **Add webhook**.
3. Enter the **Callback URL**, for example `https://my-app.example.com/webhooks/httpsms`.
4. Enter a **Signing Key**. You can generate a random one with `openssl rand -hex 32`.
5. Choose the **Events** and **Phone Numbers** you care about.
6. Click **Save Webhook**.

You can edit or delete webhooks from the same table at any time.

## Create a webhook with the API

You can also manage webhooks from code with these endpoints (all of them use the `x-api-key` header):

| Method and path                   | What it does                         |
| --------------------------------- | ------------------------------------ |
| `GET /v1/webhooks`                | List your webhooks.                  |
| `POST /v1/webhooks`               | Create a webhook.                    |
| `PUT /v1/webhooks/{webhookID}`    | Update a webhook (send every field). |
| `DELETE /v1/webhooks/{webhookID}` | Delete a webhook.                    |

```bash
curl -X POST '{{API_BASE_URL}}/v1/webhooks' \
  -H 'x-api-key: YOUR_API_KEY' \
  -H 'Content-Type: application/json' \
  -d '{
    "url": "https://my-app.example.com/webhooks/httpsms",
    "signing_key": "a-long-random-secret-key",
    "events": ["message.phone.received", "message.phone.delivered", "message.send.failed"],
    "phone_numbers": ["+18005550199"]
  }'
```

The numbers in `phone_numbers` must be phones already registered on your account (that is, with the Android app connected); otherwise the API responds with `422`. All fields are listed in the API reference: `{{API_DOCS_URL}}`.

## What the request you receive looks like

Whenever an event you subscribed to happens, your httpSMS server sends a `POST` to your URL.

### Headers

| Header          | Value                                                                                                  |
| --------------- | ------------------------------------------------------------------------------------------------------ |
| `Content-Type`  | Always `application/json`.                                                                             |
| `X-Event-Type`  | The event type, for example `message.phone.received`. Handy for routing without parsing the body.      |
| `Authorization` | `Bearer <JWT token>`. **Only sent when the webhook has a signing key.** Explained in the next section. |

### Body

The body is a [CloudEvents 1.0](https://cloudevents.io/) event serialized as JSON. The fields shared by every event are:

| Field             | Description                                                                                   |
| ----------------- | --------------------------------------------------------------------------------------------- |
| `id`              | Unique event identifier (UUID). Use it to detect duplicates.                                  |
| `type`            | The event type (same as the `X-Event-Type` header).                                           |
| `source`          | The API path that produced the event. It is informational; don't rely on its value.           |
| `specversion`     | Always `1.0`.                                                                                 |
| `time`            | When the event was created (UTC).                                                             |
| `datacontenttype` | Always `application/json`.                                                                    |
| `data`            | The event-specific data. It depends on the type; see [Webhook events](/docs/webhooks-events). |

Example of a received SMS:

```json
{
  "specversion": "1.0",
  "id": "f4aed1d3-ab4f-42b9-b9dd-9fc7182f7197",
  "source": "/v1/messages/receive",
  "type": "message.phone.received",
  "datacontenttype": "application/json",
  "time": "2026-01-15T10:20:31.524331882Z",
  "data": {
    "message_id": "0b0123bb-ef2e-468f-908a-c026d51636aa",
    "user_id": "XtABz6zdeFMoBLoltz6SREDvRSh2",
    "owner": "+18005550199",
    "encrypted": false,
    "contact": "+18005550100",
    "timestamp": "2026-01-15T10:20:30.814Z",
    "content": "Hi, what time do you open tomorrow?",
    "sim": "SIM1",
    "attachments": []
  }
}
```

## Verify that the request is authentic

Your webhook URL is public, so anyone could send it fake requests. To prevent that, set a signing key: httpSMS will then include a **JWT signed with HMAC-SHA256 (`HS256`)** using that key in every request. If the signature is valid, the request came from your httpSMS server.

The token contains these claims:

| Claim | Value                                                                                          |
| ----- | ---------------------------------------------------------------------------------------------- |
| `iss` | Always `api.httpsms.com`, even if your server uses a different domain.                         |
| `aud` | The webhook URL exactly as stored (for example `https://my-app.example.com/webhooks/httpsms`). |
| `sub` | Your httpSMS user ID.                                                                          |
| `iat` | When the token was issued.                                                                     |
| `nbf` | 10 minutes before issue time (leeway for clock drift).                                         |
| `exp` | 10 minutes after issue time.                                                                   |

To validate it:

1. Read the `Authorization` header and strip the `Bearer ` prefix.
2. Verify the token with a JWT library using your signing key, accepting **only** the `HS256` algorithm.
3. Optionally check that `aud` is your URL and `iss` is `api.httpsms.com`.
4. If anything fails, respond with `401` and drop the request.

> **Important:** the token proves **who sent** the request; it does not sign the body. Always use HTTPS for your URL so nobody can read or tamper with the request in transit.

### Node.js (Express) example

```javascript
// npm install express jsonwebtoken
const express = require('express')
const jwt = require('jsonwebtoken')

const app = express()
app.use(express.json())

const SIGNING_KEY = process.env.HTTPSMS_WEBHOOK_SIGNING_KEY

app.post('/webhooks/httpsms', (req, res) => {
  const token = (req.get('Authorization') || '').replace(/^Bearer /, '')

  try {
    jwt.verify(token, SIGNING_KEY, {
      algorithms: ['HS256'],
      issuer: 'api.httpsms.com',
    })
  } catch (error) {
    return res.status(401).send('invalid signature')
  }

  // Respond quickly, process afterwards
  res.sendStatus(200)

  const event = req.body
  if (event.type === 'message.phone.received') {
    console.log(`SMS from ${event.data.contact}: ${event.data.content}`)
  }
})

app.listen(3000)
```

### Python (Flask) example

```python
# pip install flask pyjwt
import os

import jwt
from flask import Flask, abort, request

app = Flask(__name__)
SIGNING_KEY = os.environ["HTTPSMS_WEBHOOK_SIGNING_KEY"]


@app.post("/webhooks/httpsms")
def httpsms_webhook():
    token = request.headers.get("Authorization", "").removeprefix("Bearer ")
    try:
        jwt.decode(
            token,
            SIGNING_KEY,
            algorithms=["HS256"],
            issuer="api.httpsms.com",
            options={"verify_aud": False},
        )
    except jwt.InvalidTokenError:
        abort(401)

    event = request.get_json()
    if event["type"] == "message.phone.received":
        print(f"SMS from {event['data']['contact']}: {event['data']['content']}")

    return "", 200
```

> **Note:** for simplicity these examples skip the `aud` check. To enable it, pass your exact URL (`audience` in Node, `audience=` in PyJWT) and remove `verify_aud: False`.

## Responses, timeouts and retries

Your server must respond with a **2xx** status code (for example `200`) to acknowledge the event. This is what httpSMS does for each delivery:

- **Timeout:** each attempt has at most **5 seconds**. If your server takes longer, the attempt counts as failed.
- **Retries:** if the first attempt fails (network error, timeout, or any status code **400 or higher**, 4xx included), httpSMS retries it **once** after a short pause. That is at most **2 attempts** per event and webhook in total.
- **If both attempts fail,** the event is dropped (there are no further retries). If your server has email configured and **Webhook and discord emails** is turned on in your settings, you will receive an email about the failure (at most one per hour per phone, so your inbox doesn't get flooded).

Recommendations:

- **Respond fast.** Return `200` as soon as the signature checks out and process the event in the background (a queue, a job, etc.).
- **Be idempotent.** If your server was slow to respond but did process the event, the retry delivers it again. Store the event `id` and skip events you have already handled.
- **Don't rely on ordering.** Events can arrive out of order (for example `delivered` before `sent`). Use the `timestamp` field in `data` to order them.
- **Avoid redirects.** Configure the final URL; a redirect can make the request lose its `POST` method or body.
- **If an event is critical,** don't rely on the webhook alone: you can check the message status with `GET {{API_BASE_URL}}/v1/messages/{id}`.

## Forward received SMS to Discord

If the webhook URL starts with `https://discord.com/api/webhooks/` (a Discord channel webhook URL), the `message.phone.received` event is not sent in CloudEvents format but as a Discord message showing the sender, the recipient, the text and the message ID. This lets you watch incoming SMS in a channel without writing any code. Other events are still sent as regular CloudEvents, which Discord does not understand, so subscribe a Discord webhook to `message.phone.received` only.

## Testing your webhooks

- Your httpSMS server is the one making the request, so the URL must be reachable **from the server**, not from your browser. If you are developing on your own machine, use a tunnel (for example `cloudflared` or `ngrok`) to expose your local port.
- To see exactly what arrives, you can temporarily point the webhook at a request-inspection service.
- The easiest way to trigger an event is to text your gateway's number from another phone (`message.phone.received`) or send an SMS through the API (`message.phone.sent` and `message.phone.delivered`).
- If nothing arrives, check the API logs: every webhook delivery logs the URL, the event and the response code.
