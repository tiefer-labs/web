# Progress

This file records the rebuild, phase by phase, so the work can resume if a
session is restarted. The specification is `docs/START.md` and
`docs/PROMPT.md`.

## Kit status

The kit named in `docs/START.md` phase 0 was only partly available:

- Present: `LICENSE`, `web/static/brand/`, `web/static/fonts/` (the original
  `MozillaHeadline-VF.ttf`, `MozillaText-VF.ttf` and their `OFL-*.txt`),
  the favicons, the app icons and `og-image.png`.
- Added from the brief as given in the session: `docs/START.md` and
  `docs/PROMPT.md` (the new brief, replacing the earlier one).
- Not available: `CLAUDE.md`, `.claude/settings.json`, `docs/HARDENING.md`,
  `docs/AZURE.md` and the current `docs/Tiefer_2026_az.md`. The founder
  decided to go ahead with the brief alone. Where the brief points to
  `HARDENING.md` or `AZURE.md`, the choices made are recorded in
  `docs/SECURITY-DECISIONS.md` and must be checked against those documents
  once they exist.

## Removed

Files tracked on `main` before the rebuild (`git ls-files`). All were
removed except the kit files listed above; the old site stays in the history
of `main`.

- `.dockerignore`
- `.editorconfig`
- `.env.example`
- `.github/dependabot.yml`
- `.github/workflows/ci.yml`
- `.github/workflows/codeql.yml`
- `.gitignore`
- `CONTRIBUTING.md`
- `Dockerfile`
- `LICENSE`
- `Makefile`
- `NOTICE.md`
- `README.md`
- `SECURITY.md`
- `cmd/tiefer-web/main.go`
- `deploy/Caddyfile`
- `deploy/README.md`
- `deploy/docker-compose.yml`
- `docs/PROMPT.md`
- `go.mod`
- `internal/config/config.go`
- `internal/config/config_test.go`
- `internal/contact/contact.go`
- `internal/contact/contact_test.go`
- `internal/contact/mail.go`
- `internal/contact/token.go`
- `internal/content/content.go`
- `internal/content/en.go`
- `internal/content/text_test.go`
- `internal/render/assets.go`
- `internal/render/templates.go`
- `internal/render/templates_test.go`
- `internal/server/contact.go`
- `internal/server/files.go`
- `internal/server/helpers_test.go`
- `internal/server/hosts.go`
- `internal/server/hosts_test.go`
- `internal/server/middleware.go`
- `internal/server/pages.go`
- `internal/server/ratelimit.go`
- `internal/server/ratelimit_test.go`
- `internal/server/server.go`
- `internal/server/server_test.go`
- `internal/server/text_test.go`
- `internal/textcheck/extpict.go`
- `internal/textcheck/textcheck.go`
- `internal/textcheck/textcheck_test.go`
- `web/embed.go`
- `web/static/apple-touch-icon.png`
- `web/static/brand/tiefer-logo-black.svg`
- `web/static/brand/tiefer-logo-gray.svg`
- `web/static/brand/tiefer-logo-white.svg`
- `web/static/brand/tiefer-logo.svg`
- `web/static/brand/tiefer-mark-black.svg`
- `web/static/brand/tiefer-mark-gray.svg`
- `web/static/brand/tiefer-mark-white.svg`
- `web/static/brand/tiefer-mark.svg`
- `web/static/css/site.css`
- `web/static/favicon.ico`
- `web/static/favicon.svg`
- `web/static/fonts/MozillaHeadline-VF.woff2`
- `web/static/fonts/MozillaText-VF.woff2`
- `web/static/fonts/OFL-MozillaHeadline.txt`
- `web/static/fonts/OFL-MozillaText.txt`
- `web/static/icon-192.png`
- `web/static/icon-512.png`
- `web/static/icons/alert.svg`
- `web/static/icons/chevron-down.svg`
- `web/static/js/site.js`
- `web/static/og-image.png`
- `web/templates/layout.html`
- `web/templates/pages/404.html`
- `web/templates/pages/acceptable-use.html`
- `web/templates/pages/ai-policy.html`
- `web/templates/pages/index.html`
- `web/templates/pages/legal.html`
- `web/templates/pages/privacy.html`
- `web/templates/partials/contact.html`
- `web/templates/partials/demo.html`
- `web/templates/partials/faq.html`
- `web/templates/partials/footer.html`
- `web/templates/partials/header.html`
- `web/templates/partials/hero.html`
- `web/templates/partials/how.html`
- `web/templates/partials/legal-body.html`
- `web/templates/partials/limits.html`
- `web/templates/partials/name.html`
- `web/templates/partials/principles.html`
- `web/templates/partials/problem.html`
- `web/templates/partials/product.html`
- `web/templates/partials/question.html`
- `web/templates/partials/roadmap.html`
- `web/templates/partials/use-cases.html`

## Phase log

- **Phase 0** (done): previous site removed on `site-v1`; kit kept; brief
  saved as `docs/PROMPT.md` and `docs/START.md`.
- **Phase 1** (done): Go module (standard library only), `tools/go.mod`
  with staticcheck v0.8.1, gosec v2.29.0 and govulncheck v1.8.0 as `go tool`
  entries; config with validation and redacted secrets; server limits,
  method, URL and body rules, Front Door ID check, security headers on every
  response, `/healthz`; Makefile, Dockerfile (digests pinned), CI, CodeQL and
  Dependabot (actions pinned by SHA). Next: the text check. Open:
  `vuln.go.dev` is blocked in this environment, so `make vuln` runs in CI.
