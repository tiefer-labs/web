# Tiefer website

The public website of **Tiefer**, a space technology startup from Baku that builds AI software for Earth observation
satellites, at **https://tiefer.space**. One page, three legal pages and a contact form.

- Next.js 16, React 19, TypeScript and Tailwind CSS v4; every page is prerendered.
- No cookies, analytics or trackers. Fonts are served from our own origin.
- All copy lives in [`messages/en.json`](messages/en.json), ready for more languages with next-intl.

## Run it

Node 24 and pnpm 10 or later.

```sh
pnpm install
pnpm dev            # http://localhost:3000
pnpm fmt            # format
pnpm check          # TypeScript, Oxfmt, Oxlint and Vitest
pnpm build && pnpm start
```

Configuration comes only from environment variables; every one is listed, with a comment, in
[`.env.example`](.env.example). Copy it to `.env.local` for local work. Without `SMTP_HOST` the contact form is replaced
by an "Email us" button.

To try the contact form locally, run a mail catcher and point the site at it:

```sh
docker run -d -p 1025:1025 -p 8025:8025 -e MP_SMTP_TLS_CERT=sans:localhost -e MP_SMTP_TLS_KEY=sans:localhost axllent/mailpit
```

with `SMTP_HOST=localhost`, `SMTP_PORT=1025`, `SMTP_SKIP_VERIFY=true`, `CONTACT_TO` and `CONTACT_FROM` in
`.env.local`. Messages appear at http://localhost:8025.

## Edit the copy

Everything visitors read is in [`messages/en.json`](messages/en.json), grouped by section. Lists are keyed objects; the
order of the keys is the order on the page.

- **Legal pages**: `Legal.notice`, `Legal.privacy`, `Legal.acceptableUse`. Words in curly braces such as
  `{legal_name}` are filled from the `LEGAL_*` variables; words in square brackets such as `[RETENTION_PERIOD]` are
  placeholders, highlighted on the page until real text replaces them.
- **Roadmap**: the current milestone is set in `src/components/home/roadmap.tsx`.
- Tests enforce the writing rules from tiefer-labs/lab: no long dashes, no emoji, none of the banned hype words.

[`AGENTS.md`](AGENTS.md) describes the code structure and conventions; [`STYLES.md`](STYLES.md) the design system.

## Deploy

### Vercel

Import the repository, then set the variables from `.env.example` for each environment. Set `ENV=production` and
`SITE_URL=https://tiefer.space` only for production, so previews are not indexed. Add Upstash Redis
(`UPSTASH_REDIS_REST_URL`, `UPSTASH_REDIS_REST_TOKEN`) so the contact form limits hold across serverless instances.

### Azure

The [`Dockerfile`](Dockerfile) builds a standalone Node 24 server (it sets `NEXT_OUTPUT=standalone`; other builds use
the default output) on a distroless, non-root image that listens on port
8080 and answers `/healthz`, as the infrastructure in [`infra/`](infra) expects. Pages are prerendered in the image, so
the settings they show are build arguments; secrets stay runtime settings from Key Vault.

```sh
docker build -t tiefer-web:local --build-arg SITE_URL=http://localhost:8080 --build-arg ENV=development .
docker run --rm -p 8080:8080 --read-only --tmpfs /app/.next/cache:uid=65532 --cap-drop=ALL \
  --security-opt=no-new-privileges -e CSRF_SECRET="$(openssl rand -hex 32)" tiefer-web:local
```

The deploy workflow ([`.github/workflows/deploy.yml`](.github/workflows/deploy.yml)) is manual only until the Azure
cutover. Before switching it back to deploying on every push to `main`, set the `LEGAL_*`, `CONTACT_*`, `SMTP_HOST`
and `REPO_URL` variables of the GitHub `production` environment, which it passes to the image build.

## Security

- Security headers for every response are set in [`next.config.ts`](next.config.ts): a Content Security Policy, HSTS,
  `X-Frame-Options`, a restrictive Permissions-Policy and cross-origin isolation headers. The policy allows inline
  scripts and styles because Next.js needs them for prerendered pages; nonces would make every page render per request.
- Behind Azure Front Door, [`src/proxy.ts`](src/proxy.ts) rejects requests without the right `X-Azure-FDID`.
- The contact form checks a signed, single-use, time-limited token, a honeypot field and per-client and global rate
  limits; it validates every field on the server and refuses line breaks in mail headers. Server Actions also reject
  cross-origin requests.
- Report vulnerabilities as described in [`SECURITY.md`](SECURITY.md).

## Licence

The code is licensed under the Mozilla Public License 2.0 ([`LICENSE`](LICENSE)). The Tiefer name, logo and brand
assets are not; see [`NOTICE.md`](NOTICE.md).
