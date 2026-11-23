# Outgoing Message Queue

Your Android phone isn't permanently connected to the server waiting for orders. When you send an SMS (through the API, with `bulk-send` or by uploading a file), the message first lands in a **queue** on the server. The server decides when to wake the phone up with a push notification, and the phone then downloads the message and sends it.

Thanks to this queue you can send hundreds of messages in a single request without overwhelming the phone, respect sending hours, and retry automatically when the phone doesn't respond.

## The journey of a message

1. **The API accepts the message** and stores it with status `pending`.
2. **The server works out when it should go out** (see the next section) and moves it to `scheduled`.
3. **When the time comes, it sends a push notification to the phone.** Each notification counts as one send attempt and starts the expiration timer.
4. **The phone downloads the message**, which moves to `sending`.
5. **The phone sends it** and reports it as `sent`. If the carrier returns a delivery receipt, it moves to `delivered`.
6. **If something goes wrong:** if the phone reports an error or can't be notified, the message ends up `failed`. If the phone doesn't confirm in time, it becomes `expired` and, if attempts remain, it's scheduled again.

### Outgoing message statuses

| Status      | Meaning                                                                                                   |
| ----------- | --------------------------------------------------------------------------------------------------------- |
| `pending`   | The server received the message and hasn't assigned a send time yet.                                      |
| `scheduled` | It has a send time and is waiting its turn (exact time, send rate or send schedule).                      |
| `sending`   | The phone has downloaded it and is sending it.                                                            |
| `sent`      | The phone handed it to the carrier network.                                                               |
| `delivered` | The carrier confirmed it reached the recipient.                                                           |
| `failed`    | The phone reported an error, or the server couldn't notify the phone (for example, it has no push token). |
| `expired`   | The phone didn't confirm the send within the expiration time. If attempts remain, it's retried.           |

## How the queue decides when each message goes out

The server decides using these rules, in order:

```text
Does the message have an exact time (send_at or SendTime)?
├── Yes → it goes out at that time. If the time has passed, it goes out now.
│         Neither the send rate nor the phone's schedule applies.
└── No  → 1. If the phone has a send schedule and it's currently closed,
             move to the start of the next window.
          2. Wait at least (60 ÷ messages per minute) seconds
             after the same phone's previous message.
          3. If that lands outside the schedule, move again
             to the next open window.
```

### 1. Exact time

If you pass `send_at` to `POST /v1/messages/send`, or fill the `SendTime(optional)` column in a file, the message goes out at exactly that time and **skips** every other rule. It's ideal for reminders and alerts that can't move. Formats and limits are in [Scheduling SMS Messages](/docs/scheduling-sms-messages).

> **Note:** because they skip the rate, many messages scheduled for the same minute all go out together. If there are a lot of them, stagger them.

### 2. Send rate

Messages without an exact time are spaced out using the sending phone's **Messages Per Minute** setting. With the default (10), one goes out every 6 seconds. Each number keeps its own pace, so one phone's messages never delay another's. Details in [Control SMS Send Rate](/docs/control-sms-send-rate).

For bulk sends (`/v1/messages/bulk-send` or a file), the server also numbers each sending phone's messages and delays each one according to its position:

```text
delay = position × (60 ÷ messages_per_minute)
```

| Message position | Delay at 10 per minute |
| ---------------- | ---------------------- |
| 0 (the first)    | 0 s                    |
| 1                | 6 s                    |
| 2                | 12 s                   |
| 9                | 54 s                   |

If a file uses several sender numbers (`FromPhoneNumber` column), each one gets its own numbering. Rows with `SendTime` don't count towards it.

Bulk send example through the API (up to 1000 recipients and 1024 characters):

```bash
curl -X POST '{{API_BASE_URL}}/v1/messages/bulk-send' \
  -H 'x-api-key: YOUR_API_KEY' \
  -H 'Content-Type: application/json' \
  -d '{
    "from": "+18005550199",
    "to": ["+18005550100", "+18005550101", "+18005550102"],
    "content": "Hi! This is a notice for all our customers."
  }'
```

At 10 messages per minute, the first goes out immediately, the second after 6 seconds and the third after 12.

### 3. Send schedule

If the phone has a send schedule attached, messages without an exact time only go out inside its windows. Anything that arrives outside waits as `scheduled` until the next window opens, and from then on goes out at the configured rate.

| When the message arrives      | What happens                                                |
| ----------------------------- | ----------------------------------------------------------- |
| Inside a window               | Goes out right away, following the send rate.               |
| Before the day's window opens | Waits until the window starts.                              |
| After the day's last window   | Waits for the next available window (possibly another day). |

How to create schedules and attach them to a phone: [Scheduling SMS Messages](/docs/scheduling-sms-messages).

> **Note:** the send time is worked out when the message enters the queue. If you later change the phone's rate or schedule, messages that are already scheduled keep their time.

## Expiration and retries

Sometimes the phone doesn't respond: the battery is dead, there's no connection, or the system has stopped the app. Two per-phone settings handle this:

| Setting                      | Default            | Range     | What it controls                                                                                         |
| ---------------------------- | ------------------ | --------- | -------------------------------------------------------------------------------------------------------- |
| Message Expiration (seconds) | `600` (10 minutes) | 60 – 3600 | How long the server waits, from the moment it sends the notification, for the phone to confirm the send. |
| Max Send Attempts            | `2`                | 1 – 5     | **Total** number of attempts per message, including the first.                                           |

How it works:

1. The server sends the notification and starts the expiration timer.
2. If, when the time is up, the message is still `pending`, `scheduled` or `sending`, it becomes `expired`.
3. If it still has attempts left, it goes back into the queue. This retry **does follow** the send rate and send schedule, even if the original message had an exact time.
4. Once the attempts run out, the message stays `expired` for good.

With the defaults, a phone that doesn't respond gets two attempts and the message expires for good about 20 minutes after the first notification.

> **Note:** the push notification also lapses when the expiration time is reached. If the phone stays offline longer than that, it won't receive that particular notification, but it will receive the retry's.

> **Tip:** don't make the expiration too short. If the phone sends the SMS but is slow to confirm it, a retry could produce a duplicate message. If the confirmation arrives late, the message still moves to `sent`.

## Track your messages

- **Webhooks:** receive the `message.phone.sent`, `message.phone.delivered`, `message.send.failed` and `message.send.expired` events on your server. The last one fires on every expiration; its `is_final` field is `true` when no more retries will happen. See [Webhook Events](/docs/webhooks-events).
- **API:** `GET {{API_BASE_URL}}/v1/messages/{messageID}` returns the status (`status`), attempts made (`send_attempt_count`), the maximum allowed (`max_send_attempts`) and the exact time requested (`scheduled_send_time`), if any.
- **Dashboard:** the history on `{{APP_URL}}/bulk-messages` shows how many messages from each upload are pending, scheduled, sent, delivered, failed or expired.

## Notes for self-hosted servers

With the default Docker configuration (`EVENTS_QUEUE_TYPE=emulator`), the queue's timers live in the memory of the API process. If you restart the API while messages are waiting (scheduled for later, bulk sends in progress or pending retries), those timers are lost and the affected messages may stay `pending` or `scheduled` without going out.

To avoid this:

- Restart or upgrade the API when there are no scheduled sends or running campaigns.
- After a restart, check for messages still `pending` or `scheduled` and send them again if needed.
- If you need a persistent queue, set `EVENTS_QUEUE_TYPE=cloud_tasks` and use Google Cloud Tasks.
