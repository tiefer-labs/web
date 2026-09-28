# Security decisions

Every accepted risk, suppressed finding and trade-off, with the reason and
when to revisit it. `docs/HARDENING.md` was not available during the
rebuild (see `docs/PROGRESS.md`); every entry marked **check against
HARDENING.md** must be compared with that document once it exists.

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
