# Start here

Build the public website for Tiefer at https://tiefer.space: secure, fast and calm.

Tiefer is a space technology startup from Baku, Azerbaijan. It builds AI software that runs on board Earth observation satellites (edge AI in orbit): Filter, Detect, Alert, Update.

**This is a rebuild from zero.** The repository already contains an earlier site. Remove all of it and write everything again from scratch (Phase 0). Do not reuse, copy or adapt old code, templates, styles or copy, and do not read the old code for ideas.

## Hard rules (all of them apply to everything you write)

These apply to the site, the code, comments, commit messages, the README, the pull request and your messages to me.

1. **No em dash (U+2014), no en dash (U+2013), no U+2015.** Use commas, colons, full stops or parentheses. Plain hyphens are fine.
2. **No emoji** anywhere, including U+FE0F, icons made of characters and decorative symbols.
3. A Go test enforces rules 1 and 2 on content, templates and every rendered page, and fails the build.
4. English only for now, structured so more locales can be added later.
5. Primary colour `#0C003D`. No other accent colours.
6. Mozilla Headline for headings, Mozilla Text for body, self-hosted WOFF2. Never load fonts from Google.
7. No cookies, analytics, trackers, third-party scripts or third-party requests. No reCAPTCHA.
8. Honest copy: no customer logos, "trusted by", partner names, certifications, usage statistics, testimonials or accuracy claims. Banned words: "revolutionary", "cutting-edge", "game-changing", "AI-powered", "seamless", "unlock", "leverage", "empower", "magic".
9. No real sensitive locations, no real facilities, no logos or names of real companies.
10. No unmeasured performance numbers. Numbers only inside the card labelled as illustrative.
11. Quiet design, and it must **not look AI-generated or like a template** (see `docs/PROMPT.md` section 3, "Not a template").
12. Go standard library at runtime, one self-contained binary, no Node toolchain in the repository.
13. Security as in `docs/HARDENING.md`: strict CSP without `unsafe-inline`, every header on every route, zero findings from the security tools.
14. Brand assets in `web/static/` and `LICENSE` are never modified.
15. Use the copy in `docs/PROMPT.md` word for word.

## Step 0: read, check, plan

1. Read `CLAUDE.md`, then these four files in full. Together they are the specification. Follow them exactly and use the copy word for word.
   - `docs/PROMPT.md`: content, design, Go stack, legal pages, performance, licence, acceptance criteria.
   - `docs/HARDENING.md`: security requirements and security tests. Wins on any security question.
   - `docs/AZURE.md`: Azure infrastructure as code and the deployment workflow.
   - `docs/Tiefer_2026_az.md`: background only, in Azerbaijani. If it differs, the three documents above win.
2. Check the local toolchain and report versions: `go version`, `git --version`, `docker version`, `gh --version`, `python3 -c "import fontTools, brotli"`, `az bicep version` or `bicep --version`, `node --version`. Missing tools are not a reason to stop: note them, and skip only the step that needs them.
3. Write a short plan (phases below, files, the few development tools you will pin, open questions) and show it to me. If I started you in plan mode, wait for my approval; otherwise start right after showing it, unless something in the specification is contradictory or impossible.
4. Create a todo list with one item per phase and keep it updated.

## Fixed values

| Setting | Value |
|---|---|
| Go module path | `github.com/tiefer-labs/web` |
| Production URL | `https://tiefer.space` (`www` and `http` redirect to it at Azure Front Door) |
| Contact email | `hello@tiefer.space` |
| LinkedIn URL | `https://www.linkedin.com/company/tiefer/` (footer, contact section, JSON-LD `sameAs`) |
| Source code link | `https://github.com/tiefer-labs/web` |
| Azure region (parameter default) | `germanywestcentral` |
| Azure resource name prefix | `tiefer-web` |
| GitHub deployment environment | `production` |

## Work in phases

0. **Clean slate.** Create the branch `site-v1` from `main`. Run `git ls-files` and write the list of existing files to `docs/PROGRESS.md` under "Removed". Then remove, with `git rm`, every tracked file and folder **except** this kit:
   - `CLAUDE.md`, `.claude/settings.json`, `LICENSE`
   - `docs/START.md`, `docs/PROMPT.md`, `docs/HARDENING.md`, `docs/AZURE.md`, `docs/Tiefer_2026_az.md`
   - `web/static/brand/`, `web/static/fonts/`, `web/static/favicon.svg`, `web/static/favicon.ico`, `web/static/apple-touch-icon.png`, `web/static/icon-192.png`, `web/static/icon-512.png`, `web/static/og-image.png`

   Everything else goes, including the old `README.md`, `.github/`, `go.mod`, `Dockerfile`, templates, CSS and scripts; they will be written again. Also delete untracked build output (for example `bin/`, `dist/`) but never touch `.git/` or any `.env` file. Stage the kit files with `git add` (they may be new or changed), then commit as `chore: remove previous site before rebuild`. The old site stays in the history of `main`, so this is reversible. If you find something that looks like it must not be lost (for example a real secret committed by mistake, or files unrelated to the website), stop and tell me before deleting.

From phase 1 on: commit after each phase with a clear conventional message (for example `feat(contact): signed single-use form token`). Every commit must pass `make check`. After each phase, append two or three lines to `docs/PROGRESS.md`: what is done, what is next, anything open. If this session is restarted, that file is how you resume.

1. **Scaffold with security from the first line:** Go module (standard library only), `tools/go.mod` with pinned `staticcheck`, `gosec`, `govulncheck`; `cmd/tiefer-web`; config loading and validation with secret redaction; server with the limits and timeouts from `docs/HARDENING.md` section 3; recover, security headers and method middleware; `/healthz`; header tests for every route; `Makefile`; `Dockerfile` and `.dockerignore`; `.env.example`; `ci.yml`, `codeql.yml`, `dependabot.yml` with actions pinned by commit SHA.
2. **Text check:** the Go test that fails on em dash, en dash, U+2015, emoji and banned words. Write it before any copy exists.
3. **Content and templates:** `internal/content/en.go` with all copy, layout and one partial per section, asset hashing, the template safety AST test.
4. **Design and performance:** first write `docs/DESIGN.md` (one page): the grid, the type scale with exact sizes, the spacing scale, and for **each section** a one-line description of its composition and why it differs from its neighbours. Show it to me in your message, then continue. Then build: colour tokens, WOFF2 fonts (subset and instanced as in `docs/PROMPT.md` section 9), layout following "Design restraint", hero orbital motif in inline SVG, illustrative alert card, responsive behaviour, reduced motion, gzip, cache headers, and the performance budget tests.
5. **Contact form:** everything in `docs/HARDENING.md` section 7, SMTP delivery with a fake mailer in tests, works with and without JavaScript, fuzz targets. Document local testing with Mailpit.
6. **Legal pages and security.txt:** `/legal`, `/privacy` (including the honest Azure paragraph), `/acceptable-use`, the `LEGAL_REVIEWED` banner, `/.well-known/security.txt` with its expiry test.
7. **SEO and metadata:** title, description, Open Graph, JSON-LD with its CSP hash and test, sitemap, robots, manifest, 404 page.
8. **Azure infrastructure:** `infra/` Bicep modules, `infra/README.md` runbook, `deploy.yml` with OIDC, exactly as in `docs/AZURE.md`. Validate offline with the Bicep CLI. Deploy nothing.
9. **Licence, notices and decisions:** MPL 2.0 headers in every source file, `NOTICE.md` (brand assets, fonts, Unicode data, tools), `docs/SECURITY-DECISIONS.md`.
10. **Verification:**
    - `make check`, `make fuzz`, `make docker`.
    - Run the container with `--read-only --cap-drop=ALL --security-opt=no-new-privileges` in the background, then stop it when done.
    - Headless browser (Playwright or Chromium if present): zero requests to other origins, zero console errors, zero CSP violations.
    - OWASP ZAP baseline scan if Docker is available.
    - Lighthouse on mobile if Node is already installed.
    - Screenshots of every page at 390 px and 1440 px (full page). Review them against the "Not a template" checklist in `docs/PROMPT.md` section 3, section by section, and write the findings to `docs/PROGRESS.md`. Fix, take new screenshots and review again. Do at least two rounds. Stop only when you would be comfortable showing the page to a senior designer at a space agency. Calm beats clever.

## Rules while you work

- The hard rules above apply at all times, including to commit messages, comments, the README, the pull request and this conversation.
- Standard library at runtime. Ask me before adding any runtime module. Development tools only through `tools/go.mod`, pinned.
- Do not invent customers, partners, performance numbers or testimonials. Numbers appear only inside the card labelled as illustrative.
- **Never touch anything outside this repository:** no `az login`, no Azure deployments, no DNS changes, no GitHub settings changes, no `sudo`, no global installs without asking me.
- **Secrets:** never read `.env` or other secret files, never print secrets, never put real values in any file. Use obvious placeholders.
- Never invent a commit SHA, image digest, API version or tool version. Look it up with `git ls-remote`, `docker buildx imagetools inspect` or the official documentation; if you cannot, leave a clearly marked `TODO(pin)` and list it in the report.
- Long-running processes (the server, containers) run in the background and are stopped afterwards. Do not leave the terminal blocked.
- If a network request or install fails, say so and move on; do not look for workarounds around blocked sites.
- Do not push to `main` and never force push. When all phases pass, push `site-v1` and open a pull request with `gh` if it is installed and authenticated; otherwise tell me the exact commands.

## When you finish

Give me a short report:

1. What was built, phase by phase.
2. Results: tests, fuzzing, `gosec`, `govulncheck`, `staticcheck`, CodeQL (if it ran), ZAP, Lighthouse scores, page weight against the budgets, and the screenshot paths.
3. Anything you could not do or decided yourself, and why (also recorded in `docs/SECURITY-DECISIONS.md`).
4. Every `TODO(pin)` left.
5. The exact list of what I still need to provide or do: SMTP provider and settings, legal details, Azure subscription and the runbook in `infra/README.md`, GitHub organisation settings, mail DNS records, and the HSTS preload decision.
