# Tiefer website

The public website of Tiefer at https://tiefer.space. This file is
completed in a later phase of the rebuild; see `docs/PROGRESS.md`.

## Test the contact form locally with Mailpit

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
