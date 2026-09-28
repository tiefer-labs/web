# Security decisions

Every accepted risk, suppressed finding and trade-off, with the reason and
when to revisit it. `docs/HARDENING.md` was not available during the
rebuild (see `docs/PROGRESS.md`); every entry marked **check against
HARDENING.md** must be compared with that document once it exists.

## Headers and limits

Chosen without `HARDENING.md` sections 3 and 4: **check against
HARDENING.md**.

- **One header set on every response**, applied before the handler runs,
  so errors (400, 404, 405, 413, 414, 431), redirects and panics carry it
  too. The 8 KB header limit is enforced in the middleware; `net/http`'s
  own limit (32 KB) only answers far larger headers, without the security
  headers and without a body.
- **CSP**: `default-src 'none'` and only what the site uses; the JSON-LD
  block by SHA-256 hash; `script-src-attr` and `style-src-attr 'none'`;
  Trusted Types required with no policy (the script only assigns text);
  `upgrade-insecure-requests` in production.
- **HSTS** is sent in every environment (browsers ignore it over plain
  HTTP), two years with `includeSubDomains`; `preload` only with
  `HSTS_PRELOAD=true`, a founder decision.
- **COEP `require-corp`** works because every resource is same-origin; it
  must be revisited before embedding anything from elsewhere.
- **`Referrer-Policy: no-referrer`**: the LinkedIn and GitHub links do not
  learn which page a visitor came from.
- **Limits**: ReadHeaderTimeout 5 s, ReadTimeout 10 s, WriteTimeout 30 s
  (covers SMTP), IdleTimeout 60 s, headers 8 KB, URLs 2 KB, form bodies
  16 KB, bodies refused on GET and HEAD, methods GET, HEAD and POST only,
  graceful shutdown 15 s.
- **Host check**: in production without Front Door only the `SITE_URL`
  host is answered (421 otherwise). Behind Front Door the host is the
  origin's own name, so the Front Door ID check replaces it.
- **Gzip and BREACH**: HTML is compressed per response. Pages contain a
  form token and error pages repeat submitted values, but the token is new
  on every response and cross-site posts are refused, so an attacker cannot
  collect many responses with the same secret. Revisit if a page ever
  carries a long-lived secret.

## Hardening after review

- **Proxy trust needs the Front Door ID.** `X-Azure-ClientIP` is used for
  rate limiting only when `BEHIND_FRONT_DOOR=true` and `FRONT_DOOR_ID` is
  set (so the ID check ran). Before, a development setup with
  `BEHIND_FRONT_DOOR=true` and no ID trusted a header anyone could send.
- **No addresses in net/http's error log.** Messages such as "http: panic
  serving <address>" go through a writer that replaces IPv4 and IPv6
  addresses with `[address]`.
- **Guessable CSRF keys are refused**: fewer than 8 distinct characters or
  placeholder words such as "changeme", even at 32 characters.
- **Subresource Integrity** (sha384) on the stylesheet and the script.
  Front Door caches `/static/*` at the edge; a copy altered there is refused
  by the browser. The fonts have no integrity attribute, because `@font-face`
  cannot carry one and a preload with it would download them twice; a
  changed font can only change rendering, not run code.
- **Secret scan** in the tests: private keys, cloud and GitHub tokens,
  JWTs and secret assignments fail the build.
- **Open redirects** tested: cleaned paths and fixed redirects stay on the
  site.

## Suppressed findings

| Where | Finding | Why it is safe | Revisit |
|---|---|---|---|
| `internal/server/static.go` | gosec G705 (XSS via taint) on writing a static asset | The bytes are an embedded file found by an exact map lookup on the hashed URL; the request only selects which file, it never becomes content. Content types are fixed per file and `nosniff` is set. | When static serving changes |
| `internal/contact/mail.go` | gosec G402 (`InsecureSkipVerify`) | Only set by `SMTP_SKIP_VERIFY`, which the configuration refuses in production; it exists to test against Mailpit locally. | If a local CA is used for testing instead |

## Contact form

- **Tokens instead of cookies.** The form carries an HMAC-SHA256 signed
  token with the issue time and a random nonce (`internal/contact/token.go`),
  verified with a constant-time comparison. Together with Go's
  `http.CrossOriginProtection` (Sec-Fetch-Site and Origin, with `SITE_URL`
  trusted explicitly because behind Front Door the Host header is the
  origin's name) it stops cross-site submissions without any cookie.
  Minimum age 3 s, maximum 2 h, single use.
- **Used tokens are remembered in memory.** A restart forgets them, and
  several instances would not share them, so a token could be replayed
  once per instance within two hours. Accepted: the site runs as one
  instance and the send limits still apply. Revisit before scaling out.
- **Spam answers look like success.** A filled honeypot or a form sent in
  under 3 s gets the normal success answer and nothing is sent, so bots
  learn nothing.
- **Rate limits** (in memory, never logged): 20 posts and 5 delivered
  messages per client per hour, 60 delivered messages per hour in total.
  IPv6 clients are counted per /64. At most 50,000 clients are tracked;
  when the table is full of active clients the form fails closed. Values
  chosen without HARDENING.md: **check against HARDENING.md**.
- **Mail only over TLS.** Port 465 uses implicit TLS, any other port must
  offer STARTTLS or nothing is sent; certificates are verified; TLS 1.2 or
  later. Credentials are only sent after TLS is up.
- **Header injection.** Name and organisation lose all control characters;
  the email address must be a plain address without display name or
  separators; `Build` refuses any header value with a line break; the
  sender's address goes into `Reply-To` only. Fuzz tests cover this.
- **No content in logs.** Delivery failures log the error only, never
  the message, the sender or the address of the client.

## Azure (infra/)

`docs/AZURE.md` was not available; the layout follows the brief (App
Service, Front Door with WAF, Key Vault, Container Registry, GitHub OIDC,
Bicep). **Check against AZURE.md.**

- **API versions** were taken from the stable folders of
  `Azure/azure-rest-api-specs` (each resource type checked in its version):
  App Service 2024-11-01, Front Door (Microsoft.Cdn) 2024-09-01, WAF policy
  2024-02-01, Key Vault 2024-11-01, Container Registry 2023-07-01, managed
  identities 2023-01-31, role assignments 2022-04-01, Log Analytics
  2023-09-01. Diagnostic settings use 2021-05-01-preview, the only version
  the specification has. Role IDs come from the Azure built-in roles
  documentation; the WAF rule set names from Azure quickstart templates
  and their versions (default rule set 2.1, bot manager 1.1) from the
  Front Door DRS documentation.
- **Key Vault public endpoint.** The web app has no virtual network, so
  it reaches Key Vault over the public endpoint; access is by Azure RBAC
  only (the app identity may read secrets), soft delete and purge
  protection are on. Revisit with VNet integration and a private endpoint
  if the budget allows.
- **Registry Basic tier**, public endpoint, no admin user; push and pull
  only with the two managed identities.
- **Front Door Premium by default** for the managed rule sets. Standard is
  a founder decision on cost; then only the custom rules apply.
- **Logs**: only the WAF log and the application console log go to Log
  Analytics, kept 30 days (`logRetentionDays`, which must equal
  `LOG_RETENTION_DAYS`). The Front Door access log and the App Service HTTP
  log, which hold IP addresses, are not enabled.
- **Unresolved Key Vault references** make the application refuse to
  start, so a literal reference can never become the CSRF key.
- **Deployment**: CI and deployment run in one workflow on pushes to
  `main`; the deploy job needs the `production` environment (with a
  required reviewer) and Azure login by OpenID Connect. No credentials are
  stored in GitHub.
- **The default Front Door domain** is linked at first so the site can be
  checked before DNS exists; the runbook switches it off afterwards.

## Process

- **Kit incomplete.** `HARDENING.md`, `AZURE.md`, `START.md`, `CLAUDE.md`,
  `.claude/settings.json` and the current `Tiefer_2026_az.md` were not
  available. `START.md` and the brief were saved from the session; the
  others must be added and this file checked against them.
- **govulncheck** could not reach vuln.go.dev from the build environment;
  it runs in CI (`make vuln`) on every push and weekly.
- **Unicode data** for the emoji check is Unicode 15.0 (via rivo/uniseg);
  unicode.org was not reachable. Regenerate with `gen.go` from the current
  `emoji-data.txt` when possible.
- **Company status.** The brief describes a company registered in
  Azerbaijan; the legal notice therefore shows the company fields with
  placeholders until they are configured.
