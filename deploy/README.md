# Deploying tiefer.space

This folder runs the site in production: Caddy terminates TLS for
`tiefer.space` and passes requests to the Go server. Some of the security
depends on settings outside the code (DNS, mail, the server itself). They
are listed here, with the exact records to create.

## 1. Server

Any Linux host with Docker Engine and the Compose plugin, in Azerbaijan or
the EU (see the main README on the legal side of that choice).

- Keep the system updated automatically (for example `unattended-upgrades`).
- SSH with keys only: `PasswordAuthentication no`, `PermitRootLogin no`.
- Firewall: allow only 22 (SSH, ideally from known addresses), 80 and 443
  (TCP), and 443 (UDP, for HTTP/3). Everything else closed.

## 2. DNS for tiefer.space

| Name | Type | Value | Why |
|---|---|---|---|
| `tiefer.space` | `A` / `AAAA` | the server's IPv4 / IPv6 address | the site |
| `www.tiefer.space` | `A` / `AAAA` | the same addresses | redirects to the bare domain |
| `tiefer.space` | `CAA` | `0 issue "letsencrypt.org"` | only Let's Encrypt may issue certificates |
| `tiefer.space` | `CAA` | `0 issuewild ";"` | no wildcard certificates at all |
| `tiefer.space` | `CAA` | `0 iodef "mailto:hello@tiefer.space"` | where CAs report refused requests |

Also:

- **Turn on DNSSEC** at the registrar (the `.space` zone is signed), so DNS
  answers for the domain cannot be forged.
- **Lock the domain** at the registrar (transfer lock) and protect the
  registrar account with two-factor authentication. Whoever controls the
  registrar account controls the site.

## 3. Mail for @tiefer.space

The contact form sends mail from `website@tiefer.space` (`CONTACT_FROM`) to
`hello@tiefer.space` (`CONTACT_EMAIL`) through the SMTP provider. So that nobody else can
send mail in the name of tiefer.space, and so that the form's messages are
not taken for spam:

| Name | Type | Value |
|---|---|---|
| `tiefer.space` | `TXT` | `v=spf1 include:<SPF record of your mail provider> -all` |
| `<selector>._domainkey.tiefer.space` | `TXT` or `CNAME` | the DKIM key from your mail provider |
| `_dmarc.tiefer.space` | `TXT` | `v=DMARC1; p=reject; adkim=s; aspf=s; rua=mailto:<reports address>` |
| `_smtp._tls.tiefer.space` | `TXT` | `v=TLSRPTv1; rua=mailto:<reports address>` |

Start DMARC with `p=none` for a week, read the reports, then switch to
`p=reject`.

**MTA-STS** (optional, recommended): it makes other mail servers use TLS
when they deliver to @tiefer.space. Add `A`/`AAAA` records for
`mta-sts.tiefer.space`, uncomment the `mta-sts.tiefer.space` block in the
`Caddyfile` with your provider's MX host, and add:

| Name | Type | Value |
|---|---|---|
| `_mta-sts.tiefer.space` | `TXT` | `v=STSv1; id=<date, for example 20261001>` |

If the domain sends no mail at all, publish `v=spf1 -all` and
`v=DMARC1; p=reject;` instead.

## 4. Configuration

```sh
cp .env.example deploy/.env
```

In `deploy/.env` set at least:

- `CSRF_SECRET`: `openssl rand -hex 32`;
- `SMTP_HOST`, `SMTP_PORT`, `SMTP_USER`, `SMTP_PASS` for the contact form;
- `LINKEDIN_URL` and the `LEGAL_*` values.

The addresses default to tiefer.space: `hello@tiefer.space` for contact,
security reports and Let's Encrypt notices (`CONTACT_EMAIL`,
`SECURITY_EMAIL`, `ACME_EMAIL`), `website@tiefer.space` as the sender of
form messages (`CONTACT_FROM`). Set those variables only to use other
addresses.

`ENV=production`, `SITE_URL=https://tiefer.space` and `TRUST_PROXY=true`
are set by `docker-compose.yml`. Keep `deploy/.env` readable only by the
deploying user (`chmod 600 deploy/.env`); it is ignored by git.

## 5. Start

From the repository root:

```sh
docker compose -f deploy/docker-compose.yml up -d --build
```

Caddy obtains the certificates on the first request. Check:

```sh
curl -sI https://tiefer.space/ | head -20
curl -sI https://www.tiefer.space/      # 301 to https://tiefer.space/
curl -s https://tiefer.space/.well-known/security.txt
```

External checks worth running once: an SSL Labs test (expect A+), the
Mozilla HTTP Observatory, and `internet.nl` for DNSSEC, mail and TLS.

## 6. HSTS preload

The site sends `Strict-Transport-Security: max-age=63072000;
includeSubDomains; preload`. Once HTTPS works for `tiefer.space` and
`www.tiefer.space`, submit the domain at <https://hstspreload.org>. Browsers
then never contact any part of tiefer.space over plain HTTP.

This is hard to undo: every current and future subdomain must serve HTTPS.
Do not submit before that is certain.

## 7. Operating

- **Updates**: Dependabot proposes new base images and Go versions; CI runs
  govulncheck on every push and weekly. Rebuild and restart after merging:
  `docker compose -f deploy/docker-compose.yml up -d --build`.
- **Secret rotation**: changing `CSRF_SECRET` only invalidates forms that
  are open at that moment. Rotate it and the SMTP password if the server
  or `deploy/.env` may have been exposed.
- **Logs**: the Go server logs method, path, status, size and duration, no
  IP addresses; Caddy keeps no access log. Container logs are rotated
  (5 files of 10 MB).
- **Backups**: the site has no data to back up. Keep the `caddy_data`
  volume to avoid requesting new certificates after a rebuild.
