# Scheduling SMS Messages

Scheduling lets you decide **when** a message arrives: an appointment reminder one hour before, a promotion first thing on Monday, or simply making sure your customers never get texts in the middle of the night.

There are two tools, and you can combine them:

| Tool                                    | Configured on… | Use it when…                                                                 |
| --------------------------------------- | -------------- | ---------------------------------------------------------------------------- |
| **Exact time** (`send_at` / `SendTime`) | Each message   | A specific message must go out at a precise moment.                          |
| **Send schedule**                       | Each phone     | You want a phone to send only within certain hours, e.g. Mon–Fri 9:00–18:00. |

> **Note:** a message with an exact time **ignores** the phone's send schedule and [send rate](/docs/control-sms-send-rate): it goes out at the time you set. Messages without an exact time follow both.

## Schedule a message with the API

Add the optional `send_at` field to `POST /v1/messages/send`:

```bash
curl -X POST '{{API_BASE_URL}}/v1/messages/send' \
  -H 'x-api-key: YOUR_API_KEY' \
  -H 'Content-Type: application/json' \
  -d '{
    "from": "+18005550199",
    "to": "+18005550100",
    "content": "Hi! A reminder about your appointment tomorrow at 10:00.",
    "send_at": "2026-12-01T09:30:00-06:00"
  }'
```

`send_at` rules:

| Rule             | Details                                                                                                                                                                                    |
| ---------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| Format           | [RFC 3339](https://datatracker.ietf.org/doc/html/rfc3339) **with a time zone**: `2026-12-01T09:30:00-06:00` or in UTC `2026-12-01T15:30:00Z`. Without a time zone the request is rejected. |
| Maximum          | Up to **480 hours (20 days)** from the time of the request. Beyond that, the API responds `422` with an error on `send_at`.                                                                |
| Time in the past | The message is sent immediately.                                                                                                                                                           |
| Content          | Up to 2048 characters, same as a regular send.                                                                                                                                             |

The response returns the message with status `pending` and a `scheduled_send_time` field holding the scheduled time (in UTC). Keep the `id` if you want to check on it or cancel it later.

> **Note:** `POST /v1/messages/bulk-send` does not accept `send_at`. To schedule sends to many recipients, upload a CSV or Excel file (next section) or call `/v1/messages/send` once per recipient.

## Schedule messages from a CSV or Excel file

On `{{APP_URL}}/bulk-messages` you can upload a file with up to **1000 messages**. Download a template to get started:

- CSV: `{{APP_URL}}/templates/httpsms-bulk.csv`
- Excel: `{{APP_URL}}/templates/httpsms-bulk.xlsx`

Columns, in this order:

| Column                     | Required | Description                                                                      |
| -------------------------- | -------- | -------------------------------------------------------------------------------- |
| `FromPhoneNumber`          | Yes      | Your phone number that will send the SMS. It must be registered on your account. |
| `ToPhoneNumber`            | Yes      | The recipient's number.                                                          |
| `Content`                  | Yes      | Message text, up to 1024 characters.                                             |
| `SendTime(optional)`       | No       | Date and time to send. Leave it empty to send as soon as the queue allows.       |
| `AttachmentURLs(optional)` | No       | Comma-separated attachment URLs (up to 10), sent as an MMS.                      |

Accepted `SendTime(optional)` formats:

| Example                     | How it's interpreted                         |
| --------------------------- | -------------------------------------------- |
| `2026-12-01T09:30:00`       | Local time in your **profile time zone**.    |
| `2026-12-01 09:30:00`       | Same as above (with a space instead of `T`). |
| `2026-12-01T09:30:00-06:00` | Uses the time zone included in the value.    |

You choose your profile time zone on `{{APP_URL}}/settings`. If you haven't set one, UTC is used.

> **Important:** in files, `SendTime` can be at most **420 hours (17.5 days)** in the future, less than the API's 480 hours. Rows beyond that are rejected with an error naming the row number.

> **Excel tip:** enter the date as text (for example, start the cell with an apostrophe: `'2026-12-01T09:30:00`). If Excel converts it into its own date format, the file may be rejected for an invalid format.

You can also upload the file through the API, in a form field named `document`:

```bash
curl -X POST '{{API_BASE_URL}}/v1/bulk-messages' \
  -H 'x-api-key: YOUR_API_KEY' \
  -F 'document=@messages.csv'
```

Rows with `SendTime` go out at their exact time. Rows without `SendTime` are spread out according to each sending phone's send rate.

## Send schedules

A send schedule defines **time windows** during which a phone may send SMS. Messages without an exact time that arrive outside those windows wait in the queue and go out when the next window opens, following the phone's send rate.

It's useful to:

- Respect your customers' rest (no texts at 3 AM).
- Comply with local rules on when commercial messages may be sent.
- Concentrate sends in business hours, when someone can handle replies.

### Create a schedule from the dashboard

1. Open `{{APP_URL}}/settings` and scroll down to the **Send Schedules** section.
2. Click **Create Send Schedule**.
3. Enter a name, for example `Business hours`.
4. Pick the **Timezone** the hours are interpreted in (defaults to your profile time zone).
5. Turn on the days when sending is allowed and set one or more windows per day with a **Start** and **End** time. A new schedule comes pre-filled with Monday to Friday, 09:00 to 17:00.
6. Click **Save Schedule**.

### Attach the schedule to a phone

1. On `{{APP_URL}}/settings`, **Phones** section, click **Edit** on the phone.
2. Choose the schedule in the **Send Schedule** field.
3. Click **Update Phone**.

Each phone can have at most one schedule, and the same schedule can be used on as many phones as you like. To remove it, clear the **Send Schedule** field and update the phone.

### Schedule rules

| Field                    | Allowed values                                     | Description                                                                                      |
| ------------------------ | -------------------------------------------------- | ------------------------------------------------------------------------------------------------ |
| `name`                   | 2 – 100 characters                                 | A name to identify it.                                                                           |
| `timezone`               | A valid IANA time zone, e.g. `America/Mexico_City` | Windows are interpreted in this zone; daylight-saving changes are handled automatically.         |
| `windows`                | At least 1 window                                  | List of time windows.                                                                            |
| `windows[].day_of_week`  | `0` – `6`                                          | Day of the week: `0` = Sunday, `1` = Monday … `6` = Saturday.                                    |
| `windows[].start_minute` | `0` – `1439`                                       | Start, in minutes after midnight (09:00 = `540`).                                                |
| `windows[].end_minute`   | `1` – `1440`                                       | End, in minutes after midnight (17:00 = `1020`, 24:00 = `1440`). Must be greater than the start. |

Also:

- Each day allows **up to 6 windows**, and they can't overlap.
- A window can't cross midnight. To send from 22:00 to 02:00, create two windows: `22:00–24:00` on one day and `00:00–02:00` on the next.
- The end time is exclusive: with a 09:00–17:00 window, a message arriving at exactly 17:00 waits for the next window.

### Manage schedules with the API

Create a Monday-to-Friday schedule, 9:00 to 18:00 in Mexico City:

```bash
curl -X POST '{{API_BASE_URL}}/v1/send-schedules' \
  -H 'x-api-key: YOUR_API_KEY' \
  -H 'Content-Type: application/json' \
  -d '{
    "name": "Business hours",
    "timezone": "America/Mexico_City",
    "windows": [
      { "day_of_week": 1, "start_minute": 540, "end_minute": 1080 },
      { "day_of_week": 2, "start_minute": 540, "end_minute": 1080 },
      { "day_of_week": 3, "start_minute": 540, "end_minute": 1080 },
      { "day_of_week": 4, "start_minute": 540, "end_minute": 1080 },
      { "day_of_week": 5, "start_minute": 540, "end_minute": 1080 }
    ]
  }'
```

| Action            | Request                                                     |
| ----------------- | ----------------------------------------------------------- |
| List schedules    | `GET /v1/send-schedules`                                    |
| Create a schedule | `POST /v1/send-schedules`                                   |
| Update            | `PUT /v1/send-schedules/{scheduleID}` (same body as create) |
| Delete            | `DELETE /v1/send-schedules/{scheduleID}`                    |

To attach a schedule to a phone through the API, send its `id` as `message_send_schedule_id` with `PUT /v1/phones` (see the precautions in [Control SMS Send Rate](/docs/control-sms-send-rate)).

> **Note:** deleting a schedule automatically detaches it from every phone that used it, and those phones go back to sending at any time.

### How the server picks the send time

| When the message arrives      | What happens                                                      |
| ----------------------------- | ----------------------------------------------------------------- |
| Inside a window               | Goes out right away (following the send rate).                    |
| Before the day's first window | Waits until that window starts.                                   |
| After the day's last window   | Waits for the next available window, even if it's on another day. |

The schedule is evaluated when the message enters the queue. If you change the schedule later, messages that already have a send time keep it; new messages use the updated schedule.

## View and cancel scheduled messages

- On `{{APP_URL}}/bulk-messages`, the history shows how many messages from each upload are still **Scheduled**.
- With the API you can look up a message with `GET {{API_BASE_URL}}/v1/messages/{messageID}`: `scheduled_send_time` shows the exact time you requested and `status` its current state.

To cancel a message that hasn't gone out yet (status `pending` or `scheduled`), delete it:

```bash
curl -X DELETE '{{API_BASE_URL}}/v1/messages/MESSAGE_ID' \
  -H 'x-api-key: YOUR_API_KEY'
```

The phone can no longer download it, so it won't be sent. Keep in mind the message also disappears from your history.

To understand the full journey of a message, from the moment it reaches the API until it's delivered, see [Outgoing Message Queue](/docs/outgoing-message-queue).
