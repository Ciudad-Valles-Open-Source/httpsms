# Control SMS Send Rate

Firing off a large batch of messages all at once usually goes wrong:

- **Android throttles sending.** If an app tries to send around 30 SMS in a minute, the system can pause it and ask for manual confirmation on the phone's screen.
- **Carriers watch for bursts.** Lots of back-to-back messages from the same SIM can be flagged as spam or get the line blocked.
- **You need predictability.** With a fixed rate you know how long a campaign will take to go out.

That's why every phone has a **Messages Per Minute** setting. The server spreads sends over time so the phone never exceeds that rate, no matter how many messages hit the API.

## How the rate works

The server keeps a fixed gap between one message and the next from the same phone:

```text
interval (seconds) = 60 ÷ messages_per_minute
```

With the default of **10 messages per minute**, the phone sends one SMS every **6 seconds**:

| Message | Goes out after roughly… |
| ------- | ----------------------- |
| 1st     | 0 s (immediately)       |
| 2nd     | 6 s                     |
| 3rd     | 12 s                    |
| 10th    | 54 s                    |
| 100th   | 9 min 54 s              |

Each phone (each number) has its own rate: messages going out from one phone never delay messages from another.

## Per-phone send settings

| Dashboard setting            | API field                    | Default            | Allowed values | What it does                                                                                                                     |
| ---------------------------- | ---------------------------- | ------------------ | -------------- | -------------------------------------------------------------------------------------------------------------------------------- |
| Messages Per Minute          | `messages_per_minute`        | `10`               | 1 – 30         | Maximum number of SMS the phone sends per minute.                                                                                |
| Max Send Attempts            | `max_send_attempts`          | `2`                | 1 – 5          | Total attempts per message (the first one plus retries) when the phone doesn't confirm in time.                                  |
| Message Expiration (seconds) | `message_expiration_seconds` | `600` (10 minutes) | 60 – 3600      | How long the server waits for the phone to confirm a send before marking the attempt as expired and retrying if attempts remain. |
| Send Schedule                | `message_send_schedule_id`   | none               | A schedule ID  | Restricts sending to certain hours of the day. See [Scheduling SMS Messages](/docs/scheduling-sms-messages).                     |

Defaults are applied when a phone registers for the first time from the app. Retries and expiration are covered in depth in [Outgoing Message Queue](/docs/outgoing-message-queue).

## Change the rate from the dashboard

1. Open `{{APP_URL}}/settings` and scroll down to the **Phones** section. The **Rate** column shows each number's current rate.
2. Click **Edit** on the phone you want to adjust.
3. Change **Messages Per Minute** (and, if you like, **Max Send Attempts** and **Message Expiration**).
4. Click **Update Phone**.

The new rate applies to messages that enter the queue from then on. Messages that already have a send time keep it.

## Change the rate with the API

The `PUT /v1/phones` endpoint updates the whole phone, so it's safest to read its current settings first:

```bash
curl '{{API_BASE_URL}}/v1/phones?limit=20' \
  -H 'x-api-key: YOUR_API_KEY'
```

Then send the full object with the values you want to change:

```bash
curl -X PUT '{{API_BASE_URL}}/v1/phones' \
  -H 'x-api-key: YOUR_API_KEY' \
  -H 'Content-Type: application/json' \
  -d '{
    "phone_number": "+18005550199",
    "sim": "SIM1",
    "fcm_token": "CURRENT_PHONE_TOKEN",
    "messages_per_minute": 6,
    "max_send_attempts": 2,
    "message_expiration_seconds": 600,
    "message_send_schedule_id": null
  }'
```

> **Important:** copy `fcm_token`, `sim` and `message_send_schedule_id` exactly as `GET /v1/phones` returned them. If you leave out `fcm_token`, the server clears the token and can no longer wake the phone; if you leave out `message_send_schedule_id`, the phone loses its send schedule.

API validation rules:

- `messages_per_minute` accepts up to `30`. Sending `0` leaves the current value unchanged.
- `max_send_attempts` accepts up to `5`. Sending `0` leaves the current value unchanged.
- `message_expiration_seconds` must be between `60` and `3600`, and is required when `max_send_attempts` is greater than `0`.

## Which messages follow the rate

| Message source                                                      | Follows Messages Per Minute?          |
| ------------------------------------------------------------------- | ------------------------------------- |
| `POST /v1/messages/send` without `send_at`                          | Yes                                   |
| `POST /v1/messages/bulk-send`                                       | Yes                                   |
| CSV or Excel file on `{{APP_URL}}/bulk-messages` without `SendTime` | Yes                                   |
| Automatic retries                                                   | Yes                                   |
| Any message with an exact time (`send_at` or `SendTime`)            | **No**: it goes out at the given time |

Messages with an exact time skip the queue on purpose: if you schedule several for the same moment, they go out together. If you need them to follow the rate, stagger them yourself (for example, a few seconds apart) or send them without a time and use a [send schedule](/docs/scheduling-sms-messages).

## Picking a good value

- **Start with 5–10 per minute** for large campaigns. It's the most reliable option and the least likely to trip spam filters.
- **Go up to 20–30 per minute** only if your carrier tolerates it and you need the speed.
- **Prefer divisors of 60** (1, 2, 3, 4, 5, 6, 10, 12, 15, 20 or 30). The gap between individual messages is rounded down to whole seconds, so with values such as 25 the phone may actually send at a rate of 30 per minute.
- **Need more throughput?** Spread the load across several phones or SIMs: each number has its own limit.

> **Note:** if the **Rate** column shows **Unlimited**, that phone has `0` stored (for example, from an older version) and no spacing is applied. Edit it and set a value between 1 and 30.
