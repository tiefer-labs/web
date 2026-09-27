# Tiefer website

The public website of **Tiefer**, a space technology startup from Baku that
builds an AI analyst for satellite data, at **<https://tiefer.space>**. It is
a business-card site: one main page, three legal pages and a contact form,
served by one small Go program.

- One language (Go), one self-contained binary, no Node toolchain, no CSS build step.
- Server-rendered HTML from `html/template`; about 3 KB of optional JavaScript.
- No cookies, no analytics, no third-party requests. Fonts are self-hosted.
- All copy lives in [`internal/content/en.go`](internal/content/en.go).
- Code under the Mozilla Public License 2.0; see [Licence](#licence).

## Contents

- [Run it](#run-it)
- [Edit the copy](#edit-the-copy)
- [Add another language](#add-another-language)
- [Configure the contact form (SMTP)](#configure-the-contact-form-smtp)
- [Deploy with Docker](#deploy-with-docker)
- [Security](#security)
- [Checks and tests](#checks-and-tests)
- [How it is built](#how-it-is-built)
- [Placeholders still to fill](#placeholders-still-to-fill)
- [Licence](#licence)

## Run it

You need Go 1.27 or newer (`go.mod` pins the toolchain, and `go` fetches it
automatically if yours is older).

```sh
cp .env.example .env   # optional; make run loads it
make run               # http://localhost:8080
```

| Command | What it does |
|---|---|
| `make run` | Start the server (loads `.env` if present) |
| `make build` | Build `bin/tiefer-web`, one binary with everything embedded |
| `make test` | Run all tests, including the text check |
| `make lint` | `gofmt` check, `go vet`, and `staticcheck` if installed |
| `make check` | Run only the text check, with placeholder warnings |
| `make docker` | Build the container image `tiefer-web` |

The binary needs no files next to it: templates, CSS, JavaScript, fonts and
images are embedded with `//go:embed`. Configuration comes only from
environment variables, listed with comments in [`.env.example`](.env.example).

To install staticcheck: `go install honnef.co/go/tools/cmd/staticcheck@latest`.

## Edit the copy

Every sentence on the site is in [`internal/content/en.go`](internal/content/en.go),
as typed Go values (the types are in `content.go`). Change the text there,
restart the server, and the site changes. The templates in `web/templates/`
only decide where text appears.

Rules the tests enforce (see [Checks and tests](#checks-and-tests)):

- no em dash (U+2014), en dash (U+2013) or horizontal bar (U+2015): use a
  comma, colon, full stop or parentheses; a plain hyphen is fine;
- no emoji;
- none of these words: revolutionary, cutting-edge, game-changing,
  AI-powered, seamless, unlock, leverage, empower, magic.

**Example questions.** The hero question is `Hero.Question`; the three use
case questions are `UseCases.Cards[i].Question`. Keep them about places and
activity, never about people, and never name real military sites or real
facilities of real companies.

**Demo brief.** The rows of the illustrative brief are `Demo.Rows`. The
14-day strip under it is `Demo.Strip.Days` (one entry per day, with the
number of radar, usable optical and rejected optical passes). If you change
the rows, keep the strip consistent with them. `Port A` is fictional; keep it
that way.

**Roadmap.** `Roadmap.Items` is the list of milestones in order, each with
its planned period in `When`. Set `Now: true` on the milestone you are at
(exactly one).

**Data layers.** Each entry of `Product.Layers` has `Specs`: its open data
sources with their resolution, what can be ordered on demand, and its limit.
Keep the figures to what the missions publish.

**Limits and questions.** `Limits.Items` pairs each limit of satellite data
with what Tiefer does about it. `FAQ.Items` holds the questions and answers;
an answer can end with a link (`Link`, a site path such as
`/acceptable-use`). Keep answers to what is true today, and say "on our
roadmap" for what is planned.

**Legal pages.** Their text is in `Legal` in the same file. Words in curly
braces such as `{operator}` are filled from `LEGAL_*` environment
variables. Tiefer is not yet registered as a company, so the pages name the
founder (`LEGAL_FOUNDER`) as the operator of the site; sections marked
`Stage: Founding` are shown until then. Once the company is registered, set
`LEGAL_NAME` and the other company details: the sections marked
`Stage: Company` replace them, with no code change. Words in square brackets such as `[RETENTION_PERIOD]` are
placeholders for a lawyer to replace in the file; they are highlighted on the
page and listed by `make check`.

## Add another language

The site is built for more locales (likely Azerbaijani `/az/`, German `/de/`,
Italian `/it/`). To add one:

1. Copy `internal/content/en.go` to, for example, `internal/content/az.go`,
   rename the variable to `Azerbaijani` and translate every string.
2. Set its locale: `Locale{Code: "az", Prefix: "/az", Name: "Azərbaycanca", OGLocale: "az_AZ"}`.
3. Add it to `Locales` in `internal/content/content.go`:
   `var Locales = []*Site{&English, &Azerbaijani}` (the first entry is the
   default and has no prefix).

That is all the code needed: the routes (`/az/`, `/az/legal`, ...,
`POST /az/contact`), the `lang` attribute, the `hreflang` links, the sitemap
alternates and the localised 404 page follow from the list. Two things to do
by hand:

- **Fonts.** Mozilla Headline and Mozilla Text have no `Ə ə` (needed for
  Azerbaijani) and no Cyrillic (needed for Russian). Add a self-hosted
  fallback font with an `@font-face` rule and a `unicode-range` for those
  characters in `web/static/css/site.css`, and ship its licence next to it.
- **Language switcher.** There is none yet, because there is one language.
  The page data already has `Alternates` (language code and URL of every
  version), so a switcher in `web/templates/partials/header.html` is a short
  `range` over it.

## Configure the contact form (SMTP)

The form is shown when `SMTP_HOST` is set; otherwise the contact section
shows an "Email us" button for `hello@tiefer.space` instead.
Messages are sent by email and never stored; their content is never logged.

| Variable | Meaning |
|---|---|
| `SMTP_HOST`, `SMTP_PORT` | Mail server. TLS is required: port 465 uses implicit TLS, any other port must offer STARTTLS |
| `SMTP_USER`, `SMTP_PASS` | Login (optional, but set both or neither) |
| `CONTACT_TO` | Where messages go; default `CONTACT_EMAIL` (`hello@tiefer.space`) |
| `CONTACT_FROM` | Sender address, default `website@tiefer.space`; `Reply-To` is set to the visitor |
| `CSRF_SECRET` | Signing key for form tokens, 32+ characters, required in production (`openssl rand -hex 32`) |
| `TRUST_PROXY` | `true` behind one reverse proxy, so rate limiting sees the real client |

Spam protection without third parties: a hidden honeypot field, a signed
form token with a minimum fill time of 3 seconds and a maximum age of 12
hours, single use of each token, a check that rejects cross-origin posts
(`Sec-Fetch-Site` and `Origin`), and an in-memory limit of 5 submissions per
hour per client IP address. The IP address is kept in memory only.

### Test it locally with Mailpit

[Mailpit](https://mailpit.axllent.org/) catches mail locally. Start it with
a self-signed certificate so STARTTLS works:

```sh
docker run -d --name mailpit -p 1025:1025 -p 8025:8025 \
  -e MP_SMTP_TLS_CERT=sans:localhost -e MP_SMTP_TLS_KEY=sans:localhost \
  -e MP_SMTP_REQUIRE_STARTTLS=true -e MP_SMTP_AUTH_ACCEPT_ANY=true \
  axllent/mailpit
```

Then in `.env`:

```sh
SMTP_HOST=localhost
SMTP_PORT=1025
SMTP_USER=test
SMTP_PASS=test
SMTP_SKIP_VERIFY=true   # development only; refused in production
CONTACT_TO=team@tiefer.test
CONTACT_FROM=web@tiefer.test
```

Run `make run`, send the form at http://localhost:8080/#contact, and read the
message at http://localhost:8025. Test both ways:

- **With JavaScript**: the form is sent with `fetch` and the result appears
  in place.
- **Without JavaScript** (disable it in the browser's developer tools): the
  browser posts the form, and on success is redirected to `/?sent=1#contact`
  (Post/Redirect/Get). Errors are shown next to the fields, with the values
  kept.

## Deploy with Docker

The production setup for `https://tiefer.space` is in [`deploy/`](deploy):
Caddy for TLS in front of the Go server, with Docker Compose. Follow
[`deploy/README.md`](deploy/README.md), which also lists the DNS, mail and
server settings the security depends on:

```sh
cp .env.example deploy/.env    # fill in production values, chmod 600
docker compose -f deploy/docker-compose.yml up -d --build
```

The image on its own:

```sh
make docker
docker run -p 8080:8080 \
  -e CSRF_SECRET="$(openssl rand -hex 32)" \
  tiefer-web
```

The image is built in two stages (the official Go image, then
`gcr.io/distroless/static-debian12:nonroot`, both pinned by digest), is about
20 MB, runs as user 65532, exposes `PORT` (default 8080) and has a health
check (`/healthz`, also available as `tiefer-web healthcheck`). It sets
`ENV=production`: `SITE_URL` defaults to `https://tiefer.space` and
`CONTACT_EMAIL` to `hello@tiefer.space`, it refuses to start without
`CSRF_SECRET`, answers only for the host `tiefer.space`, and logs JSON to
stderr. Log lines hold method, path,
status, size and duration, never IP addresses.

The Go server speaks plain HTTP and must sit behind a TLS-terminating proxy
(Caddy in `deploy/`). Keep the `Host` header intact, and set
`TRUST_PROXY=true` when the proxy sets `X-Forwarded-For` and
`X-Forwarded-Proto`.

The code does not depend on any provider. Two neutral options:

- **A host in Azerbaijan**: any provider that runs containers or a Linux
  virtual machine with Docker.
- **A host in the European Union**: any provider that runs containers, for
  example a managed container service or a virtual machine with Docker.

**Confirm the hosting choice with a local lawyer before launch.**
Azerbaijan's Law on Personal Data has registration and cross-border transfer
rules, and the contact form processes personal data (name, email,
organisation, message). Whether that data may be processed abroad, by the
host and by the mail provider, must be decided with counsel. Record the
decision in `LEGAL_HOSTING_*` and `LEGAL_SMTP_*` and in the privacy notice.

## Security

What protects the site, layer by layer. Report vulnerabilities as described
in [`SECURITY.md`](SECURITY.md) or `https://tiefer.space/.well-known/security.txt`.

**Transport and host**

- TLS 1.2 and 1.3 only, certificates from Let's Encrypt (Caddy), HTTP/2 and
  HTTP/3. CAA records restrict who may issue certificates; DNSSEC protects
  the DNS answers (`deploy/README.md`).
- HSTS for two years with `includeSubDomains; preload`.
- One canonical origin: `www.tiefer.space` and plain HTTP redirect to
  `https://tiefer.space`; any other `Host` header is refused (421), in the Go
  server and in Caddy, which also requires the Host header to match the TLS
  server name. This blocks host header injection, DNS rebinding and domain
  fronting.

**Browser**

- Content Security Policy starting from `default-src 'none'`: only our own
  scripts, styles, fonts and images; no inline scripts, event handlers or
  style attributes; the JSON-LD block allowed by its SHA-256 hash; no frames,
  objects, workers or plugins; `form-action 'self'`, `base-uri 'none'`,
  `frame-ancestors 'none'`; Trusted Types enforced.
- Subresource Integrity (sha384) on the stylesheet and the script.
- Cross-origin isolation: COOP `same-origin`, COEP `require-corp`, CORP
  `same-origin`, `Origin-Agent-Cluster`.
- Fetch Metadata resource isolation: other sites cannot fetch, embed or
  frame our pages or files; only following a link is allowed.
- `X-Frame-Options: DENY`, `X-Content-Type-Options: nosniff`,
  `Referrer-Policy`, a `Permissions-Policy` that turns off every browser
  feature the site does not use, `X-Permitted-Cross-Domain-Policies: none`.
- No cookies, no storage, no third-party requests, external links with
  `noopener noreferrer`.

**Contact form**

- Only URL-encoded posts up to 32 KB; every field validated and length
  limited; line breaks and control characters refused in single-line fields
  (no mail header injection).
- Cross-origin posts refused (`Sec-Fetch-Site` and `Origin`), an
  HMAC-signed single-use token with a minimum fill time and a maximum age, a
  honeypot field.
- Rate limits: 5 posts per hour per client (per /64 network for IPv6), at
  most 50 delivered messages per hour in total, a bounded limiter table that
  fails closed.
- Mail only over TLS with certificate verification; nothing stored; message
  content never logged; answers marked `no-store`, pages kept out of shared
  caches.

**Server and supply chain**

- Go standard library only (no third-party modules), one static binary in a
  distroless image, non-root, read-only file system, no Linux capabilities,
  `no-new-privileges`, memory and process limits.
- Timeouts on every phase of a request, 16 KB header limit, no built-in
  `OPTIONS *` handler, panics recovered, no software names in headers.
- Base images pinned by digest, GitHub Actions pinned by commit SHA,
  Dependabot for updates, govulncheck (also weekly), CodeQL, `go test -race`
  and staticcheck in CI.

## Checks and tests

`go test ./...` (or `make test`) runs:

- **the text check** (`make check` runs only this part): fails, with file and
  line or with page and section, on any em dash, en dash, horizontal bar or
  emoji (Unicode Extended_Pictographic, U+FE0F, keycaps, flags) in the
  content files, the templates, every text file of the repository and the
  rendered HTML of every page; and on the banned words in visible copy. It
  prints, as warnings that do not fail, every remaining `[PLACEHOLDER]` and
  every configuration value that still has its default. The copyright,
  registered and trade mark signs are allowed as text symbols.
- config validation; contact form validation (valid, missing fields, header
  injection, honeypot, too fast, expired, rate limit, cross-origin, replay,
  delivery failure, JSON mode); token round trip and tampering; security
  headers on every route; the 404 page; every page rendering with status 200;
  the CSP hash of the JSON-LD block; hashed and compressed assets; sitemap,
  robots.txt and web manifest.

CI (`.github/workflows/ci.yml`) runs `gofmt`, `go vet`, `staticcheck`,
`go test`, `make check`, `go build` and a Docker build on every push and pull
request.

## How it is built

```
cmd/tiefer-web/        main: config, logging, HTTP server, graceful shutdown
internal/config/       environment variables, loaded and validated at start
internal/content/      all copy: content.go (types), en.go (English)
internal/server/       routes, handlers, middleware, security headers, rate limit, SEO files
internal/contact/      form validation, signed tokens, email building, SMTP
internal/render/       template loading, asset hashing, template functions
internal/textcheck/    the dash, emoji and banned word checks
web/templates/         layout.html, partials/*.html, pages/*.html
web/static/            css/site.css, js/site.js, fonts/, brand/, icons, og-image.png
web/embed.go           //go:embed of templates and static files
docs/PROMPT.md         the original brief
```

Rules for changes (copy, colours, fonts, licence headers) are in
[`CONTRIBUTING.md`](CONTRIBUTING.md).

Notes on decisions:

- **Assets** get a content hash in their URL at startup
  (`{{asset "css/site.css"}}` gives `/static/css/site.0123456789.css`) and
  are served with a one-year immutable cache, gzip-compressed where it helps.
  `url(...)` references in the stylesheet are rewritten to hashed URLs too.
- **Security headers** on every response; see [Security](#security).
- **No cookies**, so the form token is not bound to a session. It is an
  HMAC-signed token with a timestamp and a nonce, combined with the
  cross-origin check in Go's `http.CrossOriginProtection`.
- **JavaScript** only enhances: the mobile menu, the header colour over the
  hero, and sending the form without a reload. Without JavaScript the
  navigation shows as a row of links and the form posts normally.
- **Motion**: the satellite in the hero moves a short way along its orbit
  once, then rests. With `prefers-reduced-motion` there is no motion.
- **Fonts**: both are variable fonts (weights 200 to 700; Mozilla Headline
  also has a width axis 75 to 125). They have no tabular figures feature, so
  `font-variant-numeric: tabular-nums` in the stylesheet has no effect with
  them.
- **No third-party Go modules.** Add one only when the standard library
  cannot do the job, justify it here, and list it in `NOTICE.md` (MIT, BSD,
  Apache 2.0 or MPL 2.0 only).

## Placeholders still to fill

Environment variables (the server logs a warning at start while any of these
has its default):

- The site already uses `https://tiefer.space`, `hello@tiefer.space` (contact
  and security reports), `website@tiefer.space` (sender of form messages),
  `https://www.linkedin.com/company/tiefer` and
  `https://github.com/tiefer-labs/web`; create these mailboxes, or set
  `CONTACT_EMAIL`, `SECURITY_EMAIL` and `CONTACT_FROM` to other addresses
- `LEGAL_FOUNDER`, the founder who operates the site until the company is
  registered
- After registration: `LEGAL_NAME`, `LEGAL_FORM`, `LEGAL_ADDRESS`,
  `LEGAL_TAX_ID` (VÖEN), `LEGAL_REGISTRATION`, `LEGAL_DIRECTOR`, all
  together
- `LEGAL_HOSTING_PROVIDER`, `LEGAL_HOSTING_COUNTRY`, `LEGAL_SMTP_PROVIDER`,
  `LEGAL_SMTP_COUNTRY`
- `LEGAL_REVIEWED=true`, only after the legal review
- SMTP settings (`SMTP_HOST`, `SMTP_PORT`, `SMTP_USER`, `SMTP_PASS`) and
  `CSRF_SECRET` for production
- the DNS records in `deploy/README.md`

Placeholders in `internal/content/en.go` for a lawyer (`make check` lists
them with the pages they appear on):

- Legal notice: `[POSTAL_ADDRESS]`, `[CONTENT_RESPONSIBILITY]`,
  `[LAST_UPDATED]`
- Privacy: `[EU_REPRESENTATIVE]`, `[HOSTING_PROVIDER_LOGS]`, `[LEGAL_BASIS]`,
  `[RECIPIENTS]`, `[INTERNATIONAL_TRANSFERS]`, `[RETENTION_PERIOD]`,
  `[LOG_RETENTION_PERIOD]`, `[DATA_SUBJECT_RIGHTS]`,
  `[SUPERVISORY_AUTHORITY_AZ]`, `[LAST_UPDATED]`
- Acceptable use: `[AUP_SCOPE]`, `[SANCTIONS_REGIMES]`, `[SCREENING_PROCESS]`,
  `[ENFORCEMENT_TERMS]`, `[LAST_UPDATED]`

## Licence

The code is licensed under the [Mozilla Public License 2.0](LICENSE)
(`MPL-2.0`). The Tiefer name, logos, icons and social image are **not**
covered and may not be used to suggest endorsement or to brand a fork. The
fonts are under the SIL Open Font License 1.1. Details, including the status
of the website copy, are in [`NOTICE.md`](NOTICE.md).
