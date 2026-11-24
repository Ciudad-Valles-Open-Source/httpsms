# Account approval and administrators

> Versión en español: [admin_ES.md](./admin_ES.md)

By default, anyone who can open your httpSMS instance can create an account with Firebase and start using it: generate API keys, register phones and send SMS through them. On a private installation that is rarely what you want.

**Account approval** fixes this: new accounts stay **pending** until an administrator approves them in the dashboard. Until then they can do nothing except view their own profile or delete their account.

---

## Table of Contents

- [How it works](#how-it-works)
- [Configuring `ADMIN_EMAILS`](#configuring-admin_emails)
- [Why the admin email must be verified](#why-the-admin-email-must-be-verified)
- [Managing users in the dashboard](#managing-users-in-the-dashboard)
- [What pending or suspended users see](#what-pending-or-suspended-users-see)
- [Effect on API keys and phones](#effect-on-api-keys-and-phones)
- [Upgrading from an earlier version](#upgrading-from-an-earlier-version)
- [Recommended Firebase settings](#recommended-firebase-settings)
- [Turning approval off](#turning-approval-off)
- [Security notes](#security-notes)

---

## How it works

Every user has a **status**:

| Status | Meaning |
| --- | --- |
| `active` | Approved account. Works normally. |
| `pending` | New account waiting for approval. Cannot use the API. |
| `suspended` | Account blocked by an administrator. Cannot use the API. |

Approval is enabled as soon as the API's `ADMIN_EMAILS` variable contains at least one email. From then on:

1. Every new account is created as `pending`, unless its email is in `ADMIN_EMAILS` (and verified): that one is created as `active`.
2. Any authenticated request from a `pending` or `suspended` user gets `403` with a code that explains why:

   ```json
   {
     "status": "error",
     "message": "...",
     "data": { "code": "account_pending" }
   }
   ```

   The code is either `account_pending` or `account_suspended`.
3. Those users have only two exceptions: `GET /v1/users/me` (so the dashboard knows their status) and `DELETE /v1/users/me` (so they can delete their account if they want to).
4. Administrators and the internal system user (`EVENTS_QUEUE_USER_ID`, which processes the event queue) are never blocked.

When `ADMIN_EMAILS` is empty, the instance behaves exactly as before: there is no approval and every account works.

A user is an **administrator** only when both conditions hold:

- The request comes from the web dashboard, authenticated with a Firebase session token (not an API key).
- The email of that session is listed in `ADMIN_EMAILS` **and Firebase marks it as verified** (`email_verified`).

---

## Configuring `ADMIN_EMAILS`

`ADMIN_EMAILS` is a comma-separated list of emails. It is case-insensitive. It is an **API** runtime variable; the web frontend does not need it.

```env
ADMIN_EMAILS=you@example.com,someone-else@example.com
```

### Local (`docker compose`)

1. Open the `.env` file at the repository root (your copy of `.env.example`).
2. Add the line `ADMIN_EMAILS=...` with your email.
3. Restart the API so it picks up the new value:

   ```bash
   docker compose up -d api
   ```

### Coolify

1. Open the httpSMS resource in Coolify and go to **Environment Variables**.
2. Add `ADMIN_EMAILS` with your email (or several, comma-separated). It is a normal runtime variable, **not** a *Build Variable*.
3. Save and redeploy (or restart) the `api` service.

If you deploy the API as a separate Coolify application, add the variable to that application.

> Use exactly the same email you sign in to the dashboard with (for example, your Google account's email).

---

## Why the admin email must be verified

Firebase lets anyone sign up with an email and password **without proving the email belongs to them**. If the API trusted the email alone, anyone could register with your address (say `you@example.com`) before you, or alongside you, and would automatically get administrator rights.

That is why the API requires Firebase to mark the email as verified. A verified address means whoever signed in controls that mailbox (or the Google account behind it), so typing your email into a form is not enough.

In practice:

- **Google sign-in:** Google already provides a verified email. It is the simplest option for administrators.
- **Email and password:** the email has to be verified (through the link Firebase sends) before you get administrator rights. After verifying it, sign out and back in so the dashboard gets a fresh token with the verified email.
- **GitHub:** Firebase does not always mark emails coming from GitHub as verified. Do not use it for the administrator account.

If you sign in with an `ADMIN_EMAILS` address that is not verified yet, the account is treated like any other user (it stays pending). As soon as the email is verified and you sign in again, you get administrator access even if the account was created as pending.

---

## Managing users in the dashboard

Administrators get a new section in the dashboard: **`/admin/users`**.

There you can:

- **List users**, filter by status (pending, active, suspended) and search by email.
- **Approve** a pending account: it becomes `active` and can use the API right away.
- **Suspend** an account: it becomes `suspended` and stops working, including its API keys and phones. You can reactivate it later.
- **Reactivate** a suspended account: it goes back to `active`.
- **Delete** an account together with all its data (phones, messages, keys, webhooks…). **This cannot be undone.**

A few rules prevent accidents:

- An administrator cannot change their own status or delete themselves from this screen.
- An administrator cannot suspend or delete another administrator. To take someone's admin rights away, remove their email from `ADMIN_EMAILS` and restart the API.

The same actions are available through the API (only with an administrator's session, never with API keys):

| Method | Path | Purpose |
| --- | --- | --- |
| `GET` | `/v1/admin/users?status=&query=&skip=&limit=` | List and search users |
| `PUT` | `/v1/admin/users/{userID}/status` | Change the status: `{"status": "active"}` or `{"status": "suspended"}` |
| `DELETE` | `/v1/admin/users/{userID}` | Delete the user and their data |

Responses: `401` when unauthenticated, `403` when the caller is not an administrator (or uses an API key), `422` when you target yourself or another administrator.

> **Suspend or delete?** If someone should not use your instance, **suspend** them. If you delete them and they sign in again, a new (pending) account is created and shows up in your list again. Only delete when you want their data gone.

---

## What pending or suspended users see

When a pending or suspended user signs in, the dashboard does not show the usual pages. Instead they see a notice:

- **Pending:** their account is waiting for an administrator's approval. They do not need to do anything; once approved they can use the dashboard normally (they may need to reload the page).
- **Suspended:** their account was suspended by an administrator.

In both cases they can sign out or delete their account. If someone is waiting for your approval, it is a good idea to let them know when you approve them.

---

## Effect on API keys and phones

Approval applies to **every** way of authenticating, not only the dashboard:

- **User API key** (`x-api-key` or `Authorization: Bearer <key>`): gets `403` with `account_pending` or `account_suspended`. Only `GET` and `DELETE /v1/users/me` work.
- **Phone API keys** (used by the Android app, they start with `pk_`): also get `403`. The app cannot register its push token, send heartbeats, fetch outstanding messages or report received SMS.

In practice:

- A pending or suspended user cannot send or receive SMS through httpSMS.
- SMS that reach the phone while the account is suspended are **not stored** in httpSMS.
- Suspending an account that already had working phones makes those phones stop syncing. Reactivating it makes them work again without reconfiguring anything.
- An API key **never** grants administrator rights, even when it belongs to a user whose email is in `ADMIN_EMAILS`. A leaked key therefore cannot be used to approve accounts or delete users.

---

## Upgrading from an earlier version

When you deploy this version, the database migration marks **every existing account as `active`**. Nobody loses access by upgrading, and the internal system user keeps working.

It also means that any account someone created before you enabled approval is still active. After you configure `ADMIN_EMAILS`:

1. Open `/admin/users` with your administrator account.
2. Review the list of active users.
3. **Suspend any account you do not recognize.** If you are sure you do not need them, you can delete them.

While `ADMIN_EMAILS` stays empty, account status has no effect at all.

---

## Recommended Firebase settings

In the Firebase console, under **Authentication**:

- **Sign-in method:** keep only the providers you actually use (Google, email/password, GitHub). Every extra provider is another way in.
- **Email/password:** if you use it, ask your users to verify their email (the dashboard has a button under **Settings** to resend the verification link). For the administrator account, verification is mandatory.
- **Settings → User account linking:** keep the default, *Link accounts that use the same email* (one account per email). It prevents the same address from ending up in several accounts with different providers.
- **Settings → Authorized domains:** keep only `localhost` (if you develop locally) and your dashboard's real domains.
- **Administrator's Google account:** turn on 2-step verification. Whoever controls that account controls your instance.

Firebase will still let anyone sign up; approval does not prevent that, it only stops those accounts from using the instance. If you want to close sign-ups completely, you can disable account creation in the Firebase/Identity Platform settings (*User actions → Enable create (sign-up)*), but then you will have to create every new account yourself in the Firebase console.

---

## Turning approval off

Leave `ADMIN_EMAILS` empty (or remove it) and restart the API. The instance goes back to the previous behavior: every account works, including pending and suspended ones, and `/admin/users` is no longer available because there are no administrators.

Each account's status is kept in the database. If you configure `ADMIN_EMAILS` again, pending and suspended accounts are blocked again.

---

## Security notes

- **Admin rights are decided on every request**, from the verified Firebase token. It is not a stored flag someone could tamper with: it depends only on `ADMIN_EMAILS` and on the email being verified.
- **Protect wherever `ADMIN_EMAILS` is configured.** Anyone who can edit the API's environment variables (your `.env`, the Coolify dashboard) can make themselves an administrator. Restrict access to Coolify and your server.
- **Keep the list short.** Every email in `ADMIN_EMAILS` can approve, suspend and delete accounts.
- **API keys are never administrators.** A leaked API key (even an administrator's) cannot use `/v1/admin/...`. Rotate it from **Settings** anyway.
- **Suspending blocks API usage from the next requests on**, but it does not end the user's Firebase session: they can still sign in and will see the suspended notice.
- **Deleting is irreversible.** Consider suspending before deleting an account.
- **Do not share your administrator account.** If someone else needs to administer the instance, add their email to `ADMIN_EMAILS` instead of lending them your session.
- Review the pending list from time to time: accounts nobody is going to approve can be deleted.
