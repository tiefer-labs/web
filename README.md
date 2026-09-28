# Tiefer website

The public website of **Tiefer**, a space technology startup from Baku that
builds AI software for Earth observation satellites, at
**https://tiefer.space**. One page and three legal pages, served by a small
Go program: secure, fast and calm.

- One binary. Templates, styles, fonts and images are embedded; nothing is
  read from disk at runtime.
- Go standard library only, no Node toolchain, no CSS build step, 2.9 KB of
  optional JavaScript. Everything works with JavaScript turned off.
- No cookies, analytics, trackers or requests to other servers. Fonts are
  self-hosted.
- All copy lives in [`internal/content/en.go`](internal/content/en.go).

The specification is in [`docs/START.md`](docs/START.md) and
[`docs/PROMPT.md`](docs/PROMPT.md); the design in
[`docs/DESIGN.md`](docs/DESIGN.md); progress in
[`docs/PROGRESS.md`](docs/PROGRESS.md).

## Run it

Go 1.27 or later (the toolchain is pinned in `go.mod`).

```sh
make run            # http://localhost:8080
make build          # bin/tiefer-web, one static binary
make check          # gofmt, vet, staticcheck, gosec, tests (text check, budgets, security)
make vuln           # govulncheck (needs network access to vuln.go.dev)
make fuzz           # every fuzz target, FUZZTIME=20s each by default
make docker         # the container image, tiefer-web:local
make placeholders   # placeholders and configuration defaults still to fill
```

Run the image the way production does:

```sh
docker run --rm -p 8080:8080 --read-only --cap-drop=ALL \
  --security-opt=no-new-privileges -e ENV=development tiefer-web:local
```

Configuration comes only from environment variables; every one is listed,
with a comment, in [`.env.example`](.env.example). In production
(`ENV=production`) the site refuses to start without `SITE_URL` and
`CSRF_SECRET`, and it never prints a secret.

## Edit the copy

Everything visitors read is in [`internal/content/en.go`](internal/content/en.go),
as typed Go structs (`content.Site`). Change a sentence there and the site
changes; the templates only read these structs, and a test checks that
every string of the index page comes from this file.

- **Use cases**: `UseCases.Cases`, one entry per event (`Label`, `Title`,
  `Text`).
- **Roadmap**: `Roadmap.Items`, one entry per milestone; set `Now: true` on
  the current one (the square fills and the `Now` label moves).
- **Legal pages**: `Legal.Notice`, `Legal.Privacy`, `Legal.AcceptableUse`.
  Words in curly braces such as `{legal_name}` are filled from `LEGAL_*`
  variables; words in square brackets such as `[RETENTION_PERIOD]` are
  placeholders for a lawyer to replace in the file. They are highlighted on
  the page and listed by `make placeholders`.

House rules, enforced by `make check`: no long dashes (U+2013, U+2014,
U+2015), no emoji or decorative symbols, none of the banned hype words, and
the MPL header on every source file.

## Add another language

1. Copy `internal/content/en.go` to, for example, `az.go`, translate it, and
   set `Locale: Locale{Code: "az", Prefix: "/az"}`.
2. Add it to `Locales` in `internal/content/content.go`.

The routes (`/az/`, `/az/privacy`, ...), the `lang` attribute, the
`hreflang` links, `x-default` and the sitemap follow automatically.

**Fonts for Azerbaijani and Russian.** Mozilla Headline and Mozilla Text
cover Latin (including Turkish, Polish, Czech, Hungarian and Romanian), but
not `Ə ə`, Cyrillic or Greek. Before the first such locale goes live, add a
self-hosted fallback with a calm, neutral design and full Latin Extended,
Cyrillic and `Ə ə` coverage, for example IBM Plex Sans or Noto Sans (both
under the OFL), subset it with `scripts/fonts.py`, and declare it with a
`unicode-range` so browsers load it only for characters the Mozilla fonts
lack:

```css
@font-face {
  font-family: "Tiefer Fallback";
  src: url(../fonts/Fallback.woff2) format("woff2");
  font-weight: 300 700;
  font-display: swap;
  unicode-range: U+018F, U+0259, U+0400-04FF, U+0370-03FF;
}
```

and add `"Tiefer Fallback"` after the Mozilla font in `--headline` and
`--text`. The fallback files are not shipped until then.

## Contact form and SMTP

The form (`POST /contact`) is shown only when `SMTP_HOST` is set; otherwise
the section shows an "Email us" button. Set `SMTP_HOST`, `SMTP_PORT` (465
for implicit TLS, any other port must offer STARTTLS), `SMTP_USER`,
`SMTP_PASS`, `CONTACT_TO` and `CONTACT_FROM`. Mail is never sent without TLS
and certificates are verified. The visitor's address goes into `Reply-To`.
Nothing is stored, and no message content is logged.

Protection without cookies or third parties: a signed single-use token,
Go's cross-origin protection (Sec-Fetch-Site and Origin), a honeypot field,
a three-second minimum fill time and in-memory rate limits.

### Test it locally with Mailpit

The form sends mail only over TLS, so run Mailpit with STARTTLS and a
self-signed certificate, and allow that certificate in development:

```sh
docker run -d --name mailpit -p 127.0.0.1:1025:1025 -p 127.0.0.1:8025:8025 \
  axllent/mailpit --smtp-tls-cert sans:localhost --smtp-tls-key sans:localhost \
  --smtp-require-starttls --smtp-auth-accept-any

SMTP_HOST=localhost SMTP_PORT=1025 SMTP_USER=dev SMTP_PASS=dev \
SMTP_SKIP_VERIFY=true CONTACT_TO=hello@tiefer.space \
CONTACT_FROM=website@tiefer.space make run
```

Open http://localhost:8080/#contact, wait a few seconds (faster posts are
treated as spam), send the form, and read the message at
http://localhost:8025. Try it with JavaScript turned off as well: the form
then posts normally and redirects to `/?sent=1#contact`. Stop Mailpit with
`docker rm -f mailpit`. `SMTP_SKIP_VERIFY` is refused in production.

## Hosting on Azure

Production runs on Microsoft Azure: Front Door with its web application
firewall in front of a Linux App Service container, images in Container
Registry, `CSRF_SECRET` and `SMTP_PASS` in Key Vault, logs in Log Analytics.
Front Door also redirects `http://` and `www.tiefer.space` to
`https://tiefer.space`. GitHub Actions deploys every push to `main` after
CI passes, logging in to Azure with OpenID Connect (no stored credentials).
The infrastructure is Bicep in [`infra/`](infra/); the step-by-step runbook
is [`infra/README.md`](infra/README.md).

The code stays provider neutral: Azure behaviour is switched on only by
configuration (`BEHIND_FRONT_DOOR`, `FRONT_DOOR_ID`), and the container runs
on any host.

**Confirm the hosting region with a local lawyer before launch.**
Azerbaijan's Law on Personal Data has registration and cross-border
transfer rules, and the contact form processes personal data. Whether that
data may be processed outside Azerbaijan (the default region is
`germanywestcentral`, and the mail provider's country) must be decided with
counsel and recorded in the privacy notice.

## Security

Summary; the reasons and accepted risks are in
[`docs/SECURITY-DECISIONS.md`](docs/SECURITY-DECISIONS.md), reporting in
[`SECURITY.md`](SECURITY.md).

- **Every response**, errors and panics included, carries a strict CSP
  (`default-src 'none'`, no `unsafe-inline`, the JSON-LD block allowed by
  its hash, Trusted Types), HSTS, `nosniff`, `X-Frame-Options: DENY`,
  `Referrer-Policy: no-referrer`, a Permissions-Policy that turns every
  feature off, and COOP, CORP and COEP.
- **Server limits**: header, read, write and idle timeouts, 8 KB of
  headers, 2 KB URLs, 16 KB form bodies, only GET, HEAD and POST.
- **Behind Front Door** the app only answers requests with the right
  `X-Azure-FDID` and takes the client address from Front Door; elsewhere in
  production it only answers its own host name.
- **Templates** call no unsafe functions and build no HTML in Go, which
  tests check on the parse trees and the Go syntax trees.
- **Logs** hold method, path, status, size and duration: no IP addresses,
  query strings, user agents or form content.
- **Supply chain**: no runtime dependencies, tools pinned in
  `tools/go.mod`, actions pinned by commit SHA, base images by digest,
  Dependabot, CodeQL, govulncheck, gosec and staticcheck in CI.
- **Container**: distroless, non-root, read-only file system, no
  capabilities; a built-in `healthcheck` command instead of a shell.

## Performance

Measured by Go tests (`internal/server/budget_test.go`) against the budgets
in `docs/PROMPT.md`: the index page is about 7 KB gzipped, the first view
about 84 KB in 10 requests (fonts 50 KB, images 17 KB), all from the same
origin; the page renders in well under 5 ms. Static files have hashed
names and are cached for a year; HTML is gzipped per response.

The fonts have no tabular-figures feature, so `font-variant-numeric:
tabular-nums` has no effect with them.

## Placeholders still to fill

Configuration (`make placeholders` lists what is still at its default):

- `SITE_URL=https://tiefer.space` and `CSRF_SECRET` (Key Vault) in
  production.
- SMTP: `SMTP_HOST`, `SMTP_PORT`, `SMTP_USER`, `SMTP_PASS` (Key Vault),
  `CONTACT_TO`, `CONTACT_FROM`.
- Company: `LEGAL_NAME`, `LEGAL_FORM`, `LEGAL_ADDRESS`, `LEGAL_TAX_ID`
  (VÖEN), `LEGAL_REGISTRATION`, `LEGAL_DIRECTOR`.
- Providers: `LEGAL_HOSTING_PROVIDER`, `LEGAL_HOSTING_COUNTRY`,
  `LEGAL_SMTP_PROVIDER`, `LEGAL_SMTP_COUNTRY`.
- `LEGAL_REVIEWED=true`, only after the legal review.
- `HSTS_PRELOAD`: a founder decision (see `infra/README.md`).

In `internal/content/en.go`, for a lawyer:

- Legal notice: `[CONTENT_RESPONSIBILITY]`, `[LAST_UPDATED]`.
- Privacy: `[EU_REPRESENTATIVE]`, `[LEGAL_BASIS_HOSTING]`,
  `[LEGAL_BASIS_CONTACT]`, `[LEGAL_BASIS]`, `[RECIPIENTS]`,
  `[INTERNATIONAL_TRANSFERS]`, `[APP_LOG_RETENTION]`, `[RETENTION_PERIOD]`,
  `[DATA_SUBJECT_RIGHTS]`, `[SUPERVISORY_AUTHORITY_AZ]`, `[LAST_UPDATED]`.
- Acceptable use: `[AUP_SCOPE]`, `[SANCTIONS_REGIMES]`,
  `[SCREENING_PROCESS]`, `[ENFORCEMENT_TERMS]`, `[LAST_UPDATED]`.

## Third-party modules

None at runtime; the standard library does everything the site needs. Any
module added later must be justified here and listed in
[`NOTICE.md`](NOTICE.md) (MIT, BSD, Apache 2.0 or MPL 2.0 only).

## Licence

The code is licensed under the [Mozilla Public License 2.0](LICENSE)
(`MPL-2.0`). The Tiefer name, logos, favicons and social image are **not**
covered and may not be used to suggest endorsement or to brand a fork. The
fonts are under the SIL Open Font License 1.1. Details in
[`NOTICE.md`](NOTICE.md).
