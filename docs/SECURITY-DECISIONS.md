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
