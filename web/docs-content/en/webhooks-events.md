# Webhook events

This page describes every event your httpSMS server can send to your [webhooks](/docs/webhooks-introduction): when it happens, what data it carries, and a realistic example of the request body.

## Common structure

Every event arrives as a JSON [CloudEvent](https://cloudevents.io/). What changes from one event to another is the content of `data`:

```json
{
  "specversion": "1.0",
  "id": "unique event identifier",
  "source": "API path that produced it",
  "type": "event type, for example message.phone.received",
  "datacontenttype": "application/json",
  "time": "when the event was created (UTC)",
  "data": {}
}
```

The same type is also sent in the `X-Event-Type` header, so you can route the request without parsing the body.

Fields that show up repeatedly in `data`:

| Field        | Description                                                                                                     |
| ------------ | --------------------------------------------------------------------------------------------------------------- |
| `user_id`    | Your httpSMS user ID.                                                                                           |
| `owner`      | **Your** phone's number (the gateway), in E.164 format.                                                         |
| `contact`    | The **other** party's number: the sender for received SMS, the recipient for sent SMS.                          |
| `sim`        | The SIM that was used: `SIM1` or `SIM2`.                                                                        |
| `timestamp`  | When the thing happened (according to the phone or the server, depending on the event). Use it to order events. |
| `encrypted`  | `true` when `content` is end-to-end encrypted (see [Getting started](/docs/getting-started)).                   |
| `request_id` | The `request_id` you sent when creating the message, or `null` if you didn't send one.                          |

Event summary:

| Event                     | When it is sent                         | Message ID in     |
| ------------------------- | --------------------------------------- | ----------------- |
| `message.phone.received`  | The phone received an SMS.              | `data.message_id` |
| `message.phone.sent`      | The phone sent an SMS.                  | `data.id`         |
| `message.phone.delivered` | The recipient received the SMS.         | `data.id`         |
| `message.send.failed`     | The SMS could not be sent.              | `data.id`         |
| `message.send.expired`    | The phone did not send the SMS in time. | `data.message_id` |
| `message.call.missed`     | The phone received a missed call.       | `data.message_id` |
| `phone.heartbeat.offline` | The phone stopped checking in.          | —                 |
| `phone.heartbeat.online`  | The phone came back online.             | —                 |

> **Note:** some events call the message ID `id` and others `message_id`. That is not a typo in these docs; it is how the API sends them.

## `message.phone.received`

Sent when your phone **receives an SMS**. This is the event you need to auto-reply, open tickets, store conversations in your CRM, and so on.

| Field         | Description                                                        |
| ------------- | ------------------------------------------------------------------ |
| `message_id`  | The message ID in httpSMS.                                         |
| `contact`     | Who sent the SMS.                                                  |
| `owner`       | Your number, the one that received the SMS.                        |
| `content`     | The message text (encrypted when `encrypted` is `true`).           |
| `attachments` | List of attachment URLs if it was an MMS; an empty list otherwise. |

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

> **Note:** `attachments` URLs point to your own server, in the form `{{API_BASE_URL}}/v1/attachments/...`.

## `message.phone.sent`

Sent when the app on the phone **hands the SMS over** to the carrier network. Being sent does not guarantee it arrived; that is what `message.phone.delivered` is for.

```json
{
  "specversion": "1.0",
  "id": "cf241cb9-0762-43d0-8b74-4f763f788c93",
  "source": "/v1/messages/ff313c14-17a3-4f74-bcb2-ca77213a64af/events",
  "type": "message.phone.sent",
  "datacontenttype": "application/json",
  "time": "2026-01-15T10:25:45.877240866Z",
  "data": {
    "id": "ff313c14-17a3-4f74-bcb2-ca77213a64af",
    "user_id": "XtABz6zdeFMoBLoltz6SREDvRSh2",
    "request_id": "order-1234",
    "owner": "+18005550199",
    "contact": "+18005550100",
    "encrypted": false,
    "timestamp": "2026-01-15T10:25:45.166Z",
    "content": "Your order 1234 is on its way",
    "sim": "SIM1"
  }
}
```

## `message.phone.delivered`

Sent when the carrier confirms the SMS **reached the recipient's phone** (delivery report).

> **Note:** not every carrier sends delivery reports. If yours doesn't, you will get `message.phone.sent` but never `message.phone.delivered`. Don't use it as your only success signal.

```json
{
  "specversion": "1.0",
  "id": "f4b86c6d-8c90-4ba7-a599-516a27f1a1e9",
  "source": "/v1/messages/5be8f09e-7007-4fe9-86b6-591d63fd38ad/events",
  "type": "message.phone.delivered",
  "datacontenttype": "application/json",
  "time": "2026-01-15T10:25:52.012355134Z",
  "data": {
    "id": "5be8f09e-7007-4fe9-86b6-591d63fd38ad",
    "owner": "+18005550199",
    "contact": "+18005550100",
    "request_id": null,
    "user_id": "XtABz6zdeFMoBLoltz6SREDvRSh2",
    "encrypted": false,
    "timestamp": "2026-01-15T10:25:51.262Z",
    "content": "Your order 1234 is on its way",
    "sim": "SIM2"
  }
}
```

## `message.send.failed`

Sent when the SMS **could not be sent**. The reason is in `error_message`.

Common `error_message` values:

| Value                                                                | What it usually means                                                                                                           |
| -------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------- |
| `NO_SERVICE`                                                         | The phone had no signal.                                                                                                        |
| `RADIO_OFF`                                                          | The radio was off (for example, airplane mode).                                                                                 |
| `GENERIC_FAILURE`                                                    | A generic carrier or phone error; often no credit or an invalid number.                                                         |
| `LIMIT_EXCEEDED`                                                     | Android blocked the send because too many SMS were sent in a row. See [Control the SMS send rate](/docs/control-sms-send-rate). |
| `NULL_PDU`                                                           | The message could not be encoded.                                                                                               |
| `UNKNOWN:<code>`                                                     | Android returned an unrecognized error code.                                                                                    |
| `Outgoing messages have been disabled on the mobile app`             | Outgoing messages are turned off for that SIM in the app.                                                                       |
| `Outgoing message is encrypted but mobile app has no encryption key` | You sent `encrypted: true` but the phone has no encryption key configured.                                                      |

It can also contain server-side errors, for example when the server could not reach the phone through the push notification.

```json
{
  "specversion": "1.0",
  "id": "d69c96fb-9f46-498b-811d-3adfb785e439",
  "source": "/v1/messages/508d783f-df33-4eaf-85e8-0fcab9958654/events",
  "type": "message.send.failed",
  "datacontenttype": "application/json",
  "time": "2026-01-15T10:30:58.979647773Z",
  "data": {
    "id": "508d783f-df33-4eaf-85e8-0fcab9958654",
    "error_message": "NO_SERVICE",
    "user_id": "XtABz6zdeFMoBLoltz6SREDvRSh2",
    "owner": "+18005550199",
    "request_id": null,
    "contact": "+18005550100",
    "timestamp": "2026-01-15T10:30:58.218Z",
    "encrypted": false,
    "content": "Your order 1234 is on its way",
    "sim": "SIM2"
  }
}
```

## `message.send.expired`

Sent when the phone **did not send the SMS within the phone's expiration window** (10 minutes by default, configurable per phone at `{{APP_URL}}/settings`). It usually happens when the phone was switched off, had no internet, or the app was put to sleep by battery saving.

httpSMS automatically retries expired messages until the phone's configured attempts are used up (2 by default). That is why this event can arrive more than once for the same message:

| Field                | Description                                                                                          |
| -------------------- | ---------------------------------------------------------------------------------------------------- |
| `message_id`         | The message ID.                                                                                      |
| `send_attempt_count` | How many attempts have been made so far.                                                             |
| `is_final`           | `false`: httpSMS will try again. `true`: that was the last attempt and the message will not be sent. |

```json
{
  "specversion": "1.0",
  "id": "8f5a9eaf-c495-4487-b11c-1824b8df6da2",
  "source": "/v1/messages/send",
  "type": "message.send.expired",
  "datacontenttype": "application/json",
  "time": "2026-01-15T10:45:43.461263458Z",
  "data": {
    "message_id": "ff313c14-17a3-4f74-bcb2-ca77213a64af",
    "owner": "+18005550199",
    "send_attempt_count": 2,
    "is_final": true,
    "request_id": null,
    "contact": "+18005550100",
    "encrypted": false,
    "user_id": "XtABz6zdeFMoBLoltz6SREDvRSh2",
    "timestamp": "2026-01-15T10:45:43.461254738Z",
    "content": "Your order 1234 is on its way",
    "sim": "SIM1"
  }
}
```

> **Tip:** if you only care about messages that are definitely lost, handle only events with `is_final: true`.

## `message.call.missed`

Sent when your phone **receives a missed call**. Useful, for example, to auto-reply with an SMS saying the number only handles text messages.

To receive it, turn on **Enable SIM1 missed call events** (or the SIM2 equivalent) in the Android app's settings.

```json
{
  "specversion": "1.0",
  "id": "f39d4b09-e2b7-42e2-b08f-d45820a18b05",
  "source": "/v1/messages/calls/missed",
  "type": "message.call.missed",
  "datacontenttype": "application/json",
  "time": "2026-01-15T11:02:35.3664741Z",
  "data": {
    "message_id": "682a47cd-adfd-4b5c-b4f2-65302746ac52",
    "user_id": "XtABz6zdeFMoBLoltz6SREDvRSh2",
    "owner": "+18005550199",
    "contact": "+18005550111",
    "timestamp": "2026-01-15T11:02:34.527Z",
    "sim": "SIM1"
  }
}
```

## `phone.heartbeat.offline`

The Android app sends a "heartbeat" to the server every 15 minutes to say "I'm still here". If the server goes about **an hour** without receiving one (more than 64 minutes, to be precise), it considers the phone offline and sends this event.

This is the alert worth watching: while the phone is offline, your messages will sit in the queue or expire.

| Field                      | Description                               |
| -------------------------- | ----------------------------------------- |
| `phone_id`                 | The phone ID in httpSMS.                  |
| `monitor_id`               | The ID of that phone's heartbeat monitor. |
| `owner`                    | The number of the offline phone.          |
| `last_heartbeat_timestamp` | When the last heartbeat was received.     |
| `timestamp`                | When the outage was detected.             |

```json
{
  "specversion": "1.0",
  "id": "5a7e1d0c-2f4b-4c8e-9a61-3b2d7c9e8f10",
  "source": "/v1/heartbeats",
  "type": "phone.heartbeat.offline",
  "datacontenttype": "application/json",
  "time": "2026-01-15T12:20:23.461263458Z",
  "data": {
    "phone_id": "a1b2c3d4-e5f6-4a7b-8c9d-0e1f2a3b4c5d",
    "user_id": "XtABz6zdeFMoBLoltz6SREDvRSh2",
    "last_heartbeat_timestamp": "2026-01-15T11:15:43.461254738Z",
    "timestamp": "2026-01-15T12:20:23.461254738Z",
    "monitor_id": "9f8e7d6c-5b4a-4392-8170-6e5d4c3b2a19",
    "owner": "+18005550199"
  }
}
```

> **Note:** the event is sent once per outage, not every hour. The next notification about that phone will be `phone.heartbeat.online` when it comes back.

## `phone.heartbeat.online`

Sent when the server receives a heartbeat from a phone that was marked **offline**. Use it to close the alert you opened on `phone.heartbeat.offline`. It has the same fields; here `last_heartbeat_timestamp` is the heartbeat that just arrived.

```json
{
  "specversion": "1.0",
  "id": "0c9b8a7f-6e5d-4c3b-a291-8f7e6d5c4b3a",
  "source": "/v1/heartbeats",
  "type": "phone.heartbeat.online",
  "datacontenttype": "application/json",
  "time": "2026-01-15T13:05:10.112233445Z",
  "data": {
    "phone_id": "a1b2c3d4-e5f6-4a7b-8c9d-0e1f2a3b4c5d",
    "user_id": "XtABz6zdeFMoBLoltz6SREDvRSh2",
    "last_heartbeat_timestamp": "2026-01-15T13:05:09.871Z",
    "timestamp": "2026-01-15T13:05:10.112230112Z",
    "monitor_id": "9f8e7d6c-5b4a-4392-8170-6e5d4c3b2a19",
    "owner": "+18005550199"
  }
}
```
