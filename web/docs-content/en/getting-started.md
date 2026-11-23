# Getting started

httpSMS turns an Android phone into an **SMS gateway**: your application makes an HTTP request and the phone sends the SMS using its own SIM card. It also works the other way around: when the phone receives an SMS, httpSMS forwards it to your server.

This edition of httpSMS is **self-hosted and 100% free**: it runs on your own server, has no paid plans and enforces no limits on messages, webhooks or contacts. The only thing that limits how many SMS you can send is your carrier and your SIM plan.

In this guide you will connect your phone to your server and send your first SMS in a few minutes.

## How it works

It helps to understand the journey of a message first, because it explains why each piece you are about to configure exists:

1. Your application calls your server's API (`{{API_BASE_URL}}`) with your API key.
2. The server stores the message in a queue and notifies the phone with a push notification.
3. The httpSMS app on the phone fetches the message and sends it as a regular SMS.
4. The phone reports back to the server whether the SMS was sent, delivered or failed.
5. If you configured [webhooks](/docs/webhooks-introduction), the server notifies you of each of those changes in real time.

## Before you start

You need:

- **Your httpSMS server up and running.** If you don't have one yet, follow the self-hosting guide `docs/self-hosting.md` in the repository.
- **An account on the web dashboard** (`{{APP_URL}}`). Sign up from the login page.
- **An Android phone with a SIM card** that can send SMS. It will be your gateway, so ideally keep it powered on, charging and connected to the internet.

## Step 1: get your API key

Every API request is authenticated with an API key sent in the `x-api-key` HTTP header. This is how the server knows which account each request belongs to.

1. Sign in to the web dashboard and open **Settings** (`{{APP_URL}}/settings`).
2. In the **API Key** section, copy the key with the copy button.
3. If you are going to sign in on the phone, you can also click **Show QR Code** and scan it from the app instead of typing it.

If a request has no key, or the key is invalid, the API responds with `401 Unauthorized`:

```json
{
  "status": "error",
  "message": "You are not authorized to carry out this request.",
  "data": "Make sure your API key is set in the [x-api-key] header in the request"
}
```

> **Important:** your API key gives full access to your account (sending and reading messages, managing webhooks…). Never publish it or ship it in browser or mobile app code. If it leaks, use **Rotate API Key** on the same settings page; the old key stops working immediately.

> **Note:** if you have several phones, you can create a separate **phone API key** for each one (they start with `pk_`). They only let the Android app talk to the server and cannot be used for the rest of the API. See [Phone API keys](/docs/phone-api-keys).

## Step 2: install the Android app

1. On the phone, download the APK from `{{APP_DOWNLOAD_URL}}` and install it. Android will ask you to allow installing apps from unknown sources.
2. Open the app and grant the **SMS permission** (the **Enable SMS Permission** button).
3. Tap **Disable Battery Optimization**. Without this, Android may put the app to sleep and your messages will sit in the queue without being sent.

> **Note for server administrators:** the app is told "there is a new message" through Firebase Cloud Messaging from **your** Firebase project. That is why the APK must be built with your project's `google-services.json` (step 8 of `docs/self-hosting.md`). An APK built for a different project can still sign in, but it will never receive the messages to send.

## Step 3: connect the phone to your server

The app doesn't know where your server is until you tell it. Fill in these fields on the login screen:

| Field                   | What to enter                                                                                                           |
| ----------------------- | ----------------------------------------------------------------------------------------------------------------------- |
| **API Key**             | The key from step 1 (or a `pk_…` phone API key). Paste it, or tap the QR icon and scan the code shown on the dashboard. |
| **Phone number (SIM1)** | The SIM's number in international E.164 format, for example `+18005550199`. The app tries to detect it automatically.   |
| **Phone number (SIM2)** | Only shown on dual-SIM phones.                                                                                          |
| **Server URL**          | Your API base URL: `{{API_BASE_URL}}` (without a trailing `/v1`).                                                       |

Tap **Login With API Key**. If everything is correct, the phone gets registered and shows up in the **Phones** section of `{{APP_URL}}/settings`.

> **Important:** the app only accepts URLs starting with `https://`. If your server runs on a local network without a certificate, put it behind an HTTPS proxy (for example Caddy, Traefik or a tunnel) before connecting the phone.

> **Note:** you can change the server URL later from **App Settings** inside the app, for example if you move to a new domain.

## Step 4: send your first SMS

To send an SMS, make a `POST` request to `{{API_BASE_URL}}/v1/messages/send` with a JSON body.

### Parameters

| Field         | Type          | Required | Description                                                                                                                           |
| ------------- | ------------- | -------- | ------------------------------------------------------------------------------------------------------------------------------------- |
| `from`        | string        | Yes      | The number of the phone that will send the SMS, in E.164 format. It must be a phone registered on your account.                       |
| `to`          | string        | Yes      | The recipient's number.                                                                                                               |
| `content`     | string        | Yes      | The message text, **up to 2048 characters**. If it doesn't fit in a single SMS, the phone sends it as a multipart message.            |
| `request_id`  | string        | No       | Your own identifier (max 255 characters). It is echoed back in webhooks, which makes it easy to match the message with your database. |
| `send_at`     | ISO 8601 date | No       | Schedules the message for later, up to 20 days (480 hours) ahead. See [Scheduling SMS](/docs/scheduling-sms-messages).                |
| `encrypted`   | boolean       | No       | Tells the server that `content` is end-to-end encrypted. See the "End-to-end encryption" section below.                               |
| `attachments` | list of URLs  | No       | Up to 10 file URLs. When present, the message is sent as an MMS.                                                                      |

### Examples

Replace `YOUR_API_KEY` with your key and the numbers with your own.

```bash
curl -X POST '{{API_BASE_URL}}/v1/messages/send' \
  -H 'x-api-key: YOUR_API_KEY' \
  -H 'Content-Type: application/json' \
  -d '{
    "from": "+18005550199",
    "to": "+18005550100",
    "content": "Hello, this is my first SMS with httpSMS"
  }'
```

```javascript
const response = await fetch('{{API_BASE_URL}}/v1/messages/send', {
  method: 'POST',
  headers: {
    'x-api-key': process.env.HTTPSMS_API_KEY,
    'Content-Type': 'application/json',
  },
  body: JSON.stringify({
    from: '+18005550199',
    to: '+18005550100',
    content: 'Hello, this is my first SMS with httpSMS',
  }),
})

console.log(await response.json())
```

```python
import os
import requests

response = requests.post(
    "{{API_BASE_URL}}/v1/messages/send",
    headers={"x-api-key": os.environ["HTTPSMS_API_KEY"]},
    json={
        "from": "+18005550199",
        "to": "+18005550100",
        "content": "Hello, this is my first SMS with httpSMS",
    },
    timeout=10,
)

print(response.status_code, response.json())
```

```php
<?php
// Requires Guzzle: composer require guzzlehttp/guzzle
$client = new GuzzleHttp\Client();

$response = $client->post('{{API_BASE_URL}}/v1/messages/send', [
    'headers' => ['x-api-key' => getenv('HTTPSMS_API_KEY')],
    'json' => [
        'from' => '+18005550199',
        'to' => '+18005550100',
        'content' => 'Hello, this is my first SMS with httpSMS',
    ],
]);

echo $response->getBody();
```

If you use Go, the official `github.com/NdoleStudio/httpsms-go` SDK works with your server once you give it the base URL:

```go
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/NdoleStudio/httpsms-go"
)

func main() {
	client := httpsms.New(
		httpsms.WithBaseURL("{{API_BASE_URL}}"),
		httpsms.WithAPIKey(os.Getenv("HTTPSMS_API_KEY")),
	)

	message, _, err := client.Messages.Send(context.Background(), &httpsms.MessageSendParams{
		From:    "+18005550199",
		To:      "+18005550100",
		Content: "Hello, this is my first SMS with httpSMS",
	})
	if err != nil {
		panic(err)
	}

	fmt.Println(message.Data.ID)
}
```

### Response

A `200` response means the message **was added to the queue**, not that it has already been sent. The actual send happens a few seconds later, when the phone picks it up:

```json
{
  "status": "success",
  "message": "message added to queue",
  "data": {
    "id": "32343a19-da5e-4b1b-a767-3298a73703cb",
    "request_id": null,
    "owner": "+18005550199",
    "contact": "+18005550100",
    "content": "Hello, this is my first SMS with httpSMS",
    "encrypted": false,
    "type": "mobile-terminated",
    "status": "pending",
    "sim": "SIM1",
    "send_attempt_count": 0,
    "max_send_attempts": 2,
    "created_at": "2026-01-15T10:20:30.123456Z"
  }
}
```

(The real response includes more date fields, such as `sent_at` or `delivered_at`, which are filled in as the message progresses.)

Keep the `id`: you can use it to look the message up with `GET {{API_BASE_URL}}/v1/messages/{id}`. The `status` field moves through `pending` → `sending` → `sent` → `delivered`, or ends in `failed` or `expired` if something goes wrong.

### Common errors

| Code  | What it means                                     | What to do                                                                                                                                  |
| ----- | ------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------- |
| `400` | The body is not valid JSON.                       | Check the JSON and the `Content-Type: application/json` header.                                                                             |
| `401` | The API key is missing or invalid.                | Check the `x-api-key` header.                                                                                                               |
| `422` | A field failed validation; details are in `data`. | Most often: `from` is not a phone registered on your account, a number is not in E.164 format, or `content` is longer than 2048 characters. |
| `500` | Internal server error.                            | Check your API logs.                                                                                                                        |

## Step 5: receive SMS

When someone texts your phone's number, the app forwards the SMS to your server and it shows up on the web dashboard. To get it into **your application** too, you have two options:

- **Webhooks (recommended):** the server sends you a `POST` in real time with the `message.phone.received` event. See [Webhooks](/docs/webhooks-introduction).
- **Query the API:** list the conversation between two numbers with `GET {{API_BASE_URL}}/v1/messages?owner=+18005550199&contact=+18005550100`.

## End-to-end encryption

If you don't want message content to be readable in your database (for example because someone else runs the server), you can encrypt it. Encryption and decryption happen only in your application and on the phone; the server only stores and forwards ciphertext.

1. In the Android app, open **App Settings** and enter an **Encryption Key** (a passphrase known only to your application and the phone).
2. To **send** encrypted messages: encrypt `content` in your application and send `"encrypted": true`. The phone decrypts it right before sending the SMS.
3. To **receive** encrypted messages: turn on **Encrypt received messages** in the app. `message.phone.received` events will arrive with `"encrypted": true` and encrypted `content`.

The format is AES-256 in CFB mode: the AES key is the SHA-256 hash of your encryption key, a random 16-byte IV is generated, and the result is `base64(IV + ciphertext)`. Python example (requires `pip install cryptography`):

```python
import base64
import hashlib
import os

from cryptography.hazmat.primitives.ciphers import Cipher, algorithms, modes


def _key(encryption_key: str) -> bytes:
    return hashlib.sha256(encryption_key.encode()).digest()


def encrypt(encryption_key: str, text: str) -> str:
    iv = os.urandom(16)
    encryptor = Cipher(algorithms.AES(_key(encryption_key)), modes.CFB(iv)).encryptor()
    data = encryptor.update(text.encode()) + encryptor.finalize()
    return base64.b64encode(iv + data).decode()


def decrypt(encryption_key: str, payload: str) -> str:
    raw = base64.b64decode(payload)
    decryptor = Cipher(algorithms.AES(_key(encryption_key)), modes.CFB(raw[:16])).decryptor()
    return (decryptor.update(raw[16:]) + decryptor.finalize()).decode()
```

> **Note:** the 2048-character limit applies to the **encrypted** text, which is longer than the original (the IV and Base64 add roughly a third). If the phone gets an encrypted message but has no encryption key configured, it marks the message as failed.

## Next steps

- [Webhooks](/docs/webhooks-introduction): receive events on your server in real time.
- [Webhook events](/docs/webhooks-events): the full list of events and their data.
- [Control the SMS send rate](/docs/control-sms-send-rate): avoid having your carrier block the SIM for sending too fast.
- [Outgoing message queue](/docs/outgoing-message-queue): what happens to messages while they wait.
- Full API reference (Swagger): `{{API_DOCS_URL}}`.
