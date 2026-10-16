# Configuration

Everything the site needs to know comes from environment variables. This page says which ones to set for each
situation, where each value comes from, and which ones only the company can decide. Every variable is also listed, with
a one-line comment, in [`.env.example`](.env.example).

## How settings work

- **Local:** put them in `.env.local` (never committed; `.env*` files are git-ignored). Start from a copy of
  `.env.example`.
- **Read at build time:** pages are prerendered, so anything a page shows is fixed when `pnpm build` runs: `ENV`,
  `SITE_URL`, `CONTACT_EMAIL`, `LINKEDIN_URL`, `REPO_URL`, `SMTP_HOST` (it decides whether the form is shown), `LEGAL_*`,
  `LOG_RETENTION_DAYS`, `HSTS_PRELOAD`. Change one of these and you must build again.
- **Read by the running server:** the mail settings, `CSRF_SECRET`, Redis and Front Door settings.
- **Secrets:** `CSRF_SECRET`, `SMTP_PASS`, `UPSTASH_REDIS_REST_TOKEN` and a `REDIS_URL` with a password. Never put them
  in the repository, a build argument, a log or a variable starting with `NEXT_PUBLIC_`.
- **Wrong values stop the site.** On start and on build it prints `Invalid configuration:` followed by each problem, for
  example `CONTACT_TO must be a valid address when SMTP_HOST is set`. Fix what it lists and run it again.

## 1. Run it locally

Nothing to set. `pnpm dev` works with the defaults: the contact section shows an "Email us" button instead of the form,
and the legal pages show `[PLACEHOLDERS]`.

## 2. Try the contact form locally

The form needs a mail server. Run Mailpit, a mail catcher that keeps every message instead of sending it. The two
certificate settings matter: the site only talks to a mail server over TLS.

```sh
docker rm -f mailpit
docker run -d --name mailpit -p 1025:1025 -p 8025:8025 -e MP_SMTP_TLS_CERT=sans:localhost -e MP_SMTP_TLS_KEY=sans:localhost axllent/mailpit
```

Its log must say `starting on [::]:1025 (STARTTLS optional)`. If it says `(no encryption)`, it was started without the
certificate settings.

Then set, in `.env.local`:

```sh
SMTP_HOST=localhost
SMTP_PORT=1025
SMTP_SKIP_VERIFY=true
CONTACT_TO=inbox@example.test
CONTACT_FROM=website@example.test
```

Leave `SMTP_USER` and `SMTP_PASS` empty (Mailpit needs no login) and `CSRF_SECRET` empty (development makes a random one
on each start). Restart `pnpm dev`. Sent messages appear at http://localhost:8025.

## 3. Production on Vercel

Set these in the Vercel project, for the **Production** environment only. Vercel passes them to both the build and the
running server.

| Variable                                             | Value                                                                                          |
| ---------------------------------------------------- | ---------------------------------------------------------------------------------------------- |
| `ENV`                                                | `production`                                                                                   |
| `SITE_URL`                                           | `https://tiefer.space`                                                                         |
| `CONTACT_EMAIL`                                      | the public address, for example `hello@tiefer.space`                                           |
| `SMTP_HOST`, `SMTP_PORT`, `SMTP_USER`, `SMTP_PASS`   | from the email provider's SMTP settings; port 587 (STARTTLS) or 465 (TLS)                      |
| `CONTACT_TO`                                         | the inbox that receives form messages                                                          |
| `CONTACT_FROM`                                       | a sender address the email provider allows for this account, on a domain set up for sending    |
| `CSRF_SECRET`                                        | 32 or more random characters: `openssl rand -hex 32`; mark it sensitive                        |
| `UPSTASH_REDIS_REST_URL`, `UPSTASH_REDIS_REST_TOKEN` | from an Upstash Redis database (Vercel Marketplace); needed because Vercel runs many instances |
| `LEGAL_*`, `LEGAL_REVIEWED`, `HSTS_PRELOAD`          | see "Values only the company can give" below                                                   |

Optional: `LINKEDIN_URL`, `REPO_URL`, `LOG_RETENTION_DAYS`. Do not set `SMTP_SKIP_VERIFY` (refused in production).

For **Preview** deployments set only what you need to test. Leave `ENV` unset there, so previews are never indexed by
search engines.

## 4. Production on Azure

The same values, in three places (details in [`infra/README.md`](infra/README.md)):

- **Shown on pages:** GitHub, repository settings, environment `production`, variables: `CONTACT_EMAIL`,
  `LINKEDIN_URL`, `REPO_URL`, `SMTP_HOST`, `CONTACT_TO`, `CONTACT_FROM`, `LEGAL_*`, `LEGAL_REVIEWED`,
  `LOG_RETENTION_DAYS`, `HSTS_PRELOAD`. The deploy workflow passes them to the image build.
- **Read by the server:** App Service settings, written by the Bicep templates in `infra/`.
- **Secrets:** Key Vault, as `csrf-secret` and `smtp-pass`. The app reads them through Key Vault references.

`BEHIND_FRONT_DOOR=true` and `FRONT_DOOR_ID` (the Front Door profile ID) are Azure only.

## Values only the company can give

An assistant must not invent these. Leave the placeholder and ask the founder.

- **`LEGAL_NAME`, `LEGAL_FORM`, `LEGAL_ADDRESS`, `LEGAL_TAX_ID`, `LEGAL_REGISTRATION`, `LEGAL_DIRECTOR`:** from the
  company registration documents.
- **`LEGAL_HOSTING_PROVIDER`, `LEGAL_HOSTING_COUNTRY`:** where the site runs, for example `Vercel Inc.` and `United
  States`.
- **`LEGAL_SMTP_PROVIDER`, `LEGAL_SMTP_COUNTRY`:** the email provider and its country. The privacy notice names them, so
  they must match the real provider.
- **`LEGAL_REVIEWED`:** `true` only after a lawyer has reviewed the legal pages. Until then a note on each legal page
  says they are not reviewed.
- **`LOG_RETENTION_DAYS`:** must match how long the host keeps request logs.
- **`HSTS_PRELOAD`:** a founder decision. Once the domain is on the browsers' preload list, it is hard to undo.
- **`CONTACT_TO`, `CONTACT_EMAIL`, the SMTP account:** which inbox and which email provider the company uses.

Until a `LEGAL_*` value is set, its `[PLACEHOLDER]` is highlighted on the legal pages, so a missing value is easy to
spot.
