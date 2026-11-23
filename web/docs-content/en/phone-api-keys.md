# Phone API Keys

A **phone API key** is a credential made only for the Android app. Instead of putting your account API key (which has full access to the API) on every phone, you create a separate key for each device.

Why bother?

- **Less risk:** if a phone is lost or stolen, you delete just its key. Your account key, which your integrations and servers rely on, keeps working untouched.
- **Least privilege:** a phone key can only do what the app needs (heartbeats, received messages, sent and delivery reports). It cannot send SMS or read your history through the API.
- **Per-number control:** each key knows which phone numbers belong to it and rejects requests for any other number.

You can create as many phone API keys as you need on this server.

## Account key vs. phone key

|                       | Account API key                                   | Phone API key                                     |
| --------------------- | ------------------------------------------------- | ------------------------------------------------- |
| Where to find it      | Settings page: `{{APP_URL}}/settings`             | Phone API keys page: `{{APP_URL}}/phone-api-keys` |
| Format                | Random string                                     | Always starts with `pk_`                          |
| What it can do        | The whole API: send messages, webhooks, contacts… | Only the endpoints the Android app uses           |
| How many you can have | One per account                                   | As many as you like (for example, one per phone)  |
| Revoking it affects…  | Every integration and phone using it              | Only the phones signed in with it                 |

> **Note:** to send messages, configure webhooks or call any other API endpoint (`{{API_DOCS_URL}}`), use your account key. A `pk_` key on those endpoints gets a `401 Unauthorized` error.

## What a phone key can do

A `pk_` key is accepted only by these endpoints, which are exactly the ones the app calls:

| Method | Endpoint                          | Used by the app to…                              |
| ------ | --------------------------------- | ------------------------------------------------ |
| `PUT`  | `/v1/phones/fcm-token`            | Register the phone when you sign in              |
| `POST` | `/v1/heartbeats`                  | Send heartbeats showing the phone is online      |
| `GET`  | `/v1/messages/outstanding`        | Download the next message it has to send         |
| `POST` | `/v1/messages/{messageID}/events` | Report that an SMS was sent, failed or delivered |
| `POST` | `/v1/messages/receive`            | Store an incoming SMS                            |
| `POST` | `/v1/messages/calls/missed`       | Store a missed call                              |

On top of that, for heartbeats, received messages and send reports the server checks that the phone number involved is linked to the key. If it isn't, it answers `401` with _"You are not authorized to carry out the request for this phone number"_.

## Create a phone API key

### From the dashboard

1. Open `{{APP_URL}}/phone-api-keys`.
2. Click **Create API Key**.
3. Enter a name that helps you recognise the device (1 to 60 characters), for example `Office phone`.
4. Click **Create Phone API Key**. The new key shows up in the table.

### With the API

Use your **account** key to create phone keys:

```bash
curl -X POST '{{API_BASE_URL}}/v1/phone-api-keys' \
  -H 'x-api-key: YOUR_API_KEY' \
  -H 'Content-Type: application/json' \
  -d '{ "name": "Office phone" }'
```

The response contains the full key in `data.api_key` (it starts with `pk_`), plus `phone_numbers` and `phone_ids`, which stay empty until a phone signs in with it.

## Sign in to the Android app

Before you start:

- Install the latest version of the app from `{{APP_DOWNLOAD_URL}}`.
- The phone needs Google Play Services: the app gets a push notification (FCM) token and cannot sign in without it.

Steps:

1. On `{{APP_URL}}/phone-api-keys`, click **View** on the key you want to use. A dialog opens with the key as text, a copy button and a QR code.
2. In the app, tap the QR scanner icon next to the **API Key** field and scan the code. The key fills in automatically (you can also paste it by hand).
3. Enter the SIM's phone number in international E.164 format, for example `+18005550199`. If the phone has two SIMs, the app also asks for the SIM2 number.
4. In **Server URL**, enter your API address: `{{API_BASE_URL}}`.
5. Tap **Login With API Key**.

> **Important:** the QR code only holds the key, not the server address. Always check the **Server URL** field: if it shows a different address, replace it with `{{API_BASE_URL}}`. The app only accepts `https://` addresses.

When you sign in, the app registers each number with the server, which links it to the key automatically. Within a few seconds the number appears in the **Phone Numbers** column of the table (the page refreshes on its own).

> **Note:** the app also accepts your account key for signing in, but it is not recommended: anyone with access to the phone would get full access to your API.

## How numbers are linked to a key

- **A key can hold several numbers.** For example, both SIMs of the same phone are linked to the same key. You can also share one key between several phones, although one key per device is the safest setup.
- **A number can belong to only one key at a time.** If you sign in with a number using a different key, the server removes it from the old key and moves it to the new one automatically.
- **If you delete a phone** from your account, its number is removed from the key as well.

## Remove a number from a key

Use this when you want a phone to stop working with a key without deleting the whole key.

1. On `{{APP_URL}}/phone-api-keys`, click **Remove** next to the number.
2. Confirm with **Remove Phone from key**.

From then on, the server rejects heartbeats, received messages and reports for that number made with that key. To get the phone working again, **log out of the app and sign in again** with a valid key.

With the API you need the key ID and the phone ID (get it from `GET {{API_BASE_URL}}/v1/phones`):

```bash
curl -X DELETE '{{API_BASE_URL}}/v1/phone-api-keys/KEY_ID/phones/PHONE_ID' \
  -H 'x-api-key: YOUR_API_KEY'
```

## Delete a key

1. On `{{APP_URL}}/phone-api-keys`, click **Delete** in the **Actions** column.
2. Confirm with **Delete API Key**.

Deleting is permanent. Every phone using that key can no longer talk to the server: **log out of the app and sign in with a new key** on each of them.

```bash
curl -X DELETE '{{API_BASE_URL}}/v1/phone-api-keys/KEY_ID' \
  -H 'x-api-key: YOUR_API_KEY'
```

## List your keys with the API

```bash
curl '{{API_BASE_URL}}/v1/phone-api-keys?limit=100' \
  -H 'x-api-key: YOUR_API_KEY'
```

| Parameter | Default | Range     | Description                             |
| --------- | ------- | --------- | --------------------------------------- |
| `limit`   | `1`     | 1 – 100   | How many keys to return                 |
| `skip`    | `0`     | 0 or more | How many keys to skip (pagination)      |
| `query`   | —       | max. 100  | Only keys whose name contains this text |

> **Note:** if you don't send `limit`, the API returns **only one** key. Always pass `limit` to see the full list.

## Troubleshooting

| Symptom                                                       | Likely cause                                    | What to do                                    |
| ------------------------------------------------------------- | ----------------------------------------------- | --------------------------------------------- |
| The app shows _"Cannot validate the API key"_ when signing in | The key was copied wrongly or has been deleted  | Scan the QR code again or create a new key    |
| The app says the server URL _"must be HTTPS"_                 | The address starts with `http://`               | Use `{{API_BASE_URL}}` with `https://`        |
| The app cannot find the FCM token                             | The phone doesn't have Google Play Services     | Install or update Google Play Services        |
| `401` responses for a number after clicking **Remove**        | The number is no longer linked to the key       | Log out and sign in again in the app          |
| `401` when trying to send an SMS with a `pk_` key             | Phone keys cannot send messages through the API | Use your account key (`{{APP_URL}}/settings`) |
