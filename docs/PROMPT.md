# Build the Tiefer website

You are building the first public website for **Tiefer**, an early-stage B2B startup. This is a business-card website: one main page plus two legal pages, served by a small **Go** backend. There is no product dashboard and no login yet. Build it so it is fast, accessible, easy to edit and easy to extend later with more languages and a product area (the future dashboard will live in the same Go codebase).

The repository is licensed under the **Mozilla Public License 2.0** (see section 13).

Read this whole brief before you start. Where the brief gives exact copy, use it word for word. Do not invent claims, customers, partners, numbers or testimonials.

---

## 1. Context: what Tiefer does

Tiefer is a **space technology startup** based in **Baku, Azerbaijan**. It builds an **AI analyst for satellite data**: a user asks a question in plain language, and Tiefer finds and fuses optical, radar (SAR), thermal and night-light satellite data, runs change and object detection, tasks new commercial imagery when the archive is not enough, and returns a structured, sourced intelligence brief.

Example questions:

- "Which ports on the eastern Caspian coast saw more vessel activity in the last 14 days?"
- "Map flooded buildings in this district after last week's storm."
- "Find new clearings in this forest area in the last 10 days and rank them by risk."

How it works, in five steps: **Ask** (understand the question), **Plan** (pick sensors, use radar when it is cloudy), **Fuse** (find, clean and combine imagery, run models), **Task** (order new imagery within budget rules, human approval beyond them), **Brief** (map, findings, sources, confidence, and what could not be seen).

Customers, in order: commodity traders and analysts, then insurers, then governments (sovereign, on-premise deployments). Sales start in the Caspian, Caucasus, Black Sea and Central Asia region, with interfaces in English, Azerbaijani, Turkish and Russian; the product itself works worldwide.

Product principles that the site must reflect:

1. Observed and inferred are always labelled separately.
2. Generated pixels are never shown as evidence. If it is cloudy, the brief says so.
3. Every finding links to its source image, sensor, date and processing steps.
4. Spending decisions (tasking) stay under human control.
5. Limits are stated openly.
6. Tiefer is not used to track individuals.

"Tiefer" is German for "deeper". Brand face: **See deeper.** Product promise: **Space intelligence you can ask.** Visitors will be traders, risk analysts, investors, accelerator reviewers and potential partners, mostly reading in English.

---

## 2. Hard rules (non-negotiable)

1. **Language:** English only for now. Structure the project so more locales can be added later without rewriting templates (likely candidates: Azerbaijani `/az/`, German `/de/`, Italian `/it/`): all copy lives in one content file per locale.
2. **No emoji** anywhere: copy, alt text, meta tags, code comments that ship, icons made of characters.
3. **No em dash (U+2014) and no en dash (U+2013)** anywhere in the site. Also avoid U+2015. Use commas, colons, full stops or parentheses instead. Plain hyphens (U+002D) in compound words and ranges like "2027-2029" are fine.
4. Add an automated check (see section 10) that fails the build if rule 2 or 3 is broken.
5. **Primary colour:** `#0C003D`.
6. **Fonts:** **Mozilla Headline** for headings, **Mozilla Text** for body text. Both are SIL Open Font License 1.1.
7. **Self-host the fonts.** Do not load fonts from Google's servers at runtime (most visitors are in the EU, where courts have ruled that loading Google Fonts from Google's CDN transmits visitor IP addresses without consent; it also keeps the site free of third-party data transfers under Azerbaijani law). The font files are **already provided** in `web/static/fonts/`: `MozillaHeadline-VF.ttf` (variable, axes `wght` 200-700 and `wdth` 75-125) and `MozillaText-VF.ttf` (variable, axis `wght` 200-700), taken from the official Google Fonts repository, with their `OFL-*.txt` licences. Convert both to WOFF2 (for example with `fonttools` and `brotli`, or `woff2_compress`), serve only the WOFF2 files, keep the licences next to them, and use `font-weight` ranges in `@font-face` so the variable axes work. If you cannot convert in your environment, serve the TTF files and tell me.
7a. **Glyph coverage.** Both Mozilla fonts cover Latin including Turkish, Polish, Czech, Hungarian and Romanian, but they do **not** contain the Azerbaijani letters `Ə ə` and have **no Cyrillic or Greek**. The site is English now, but Azerbaijani and Russian versions are planned. Prepare a self-hosted fallback font with a calm, neutral style and full Latin Extended, Cyrillic and `Ə ə` coverage (for example IBM Plex Sans or Noto Sans, both open licence), declared with `unicode-range` so it only loads for characters the Mozilla fonts lack. Do not ship the fallback files until a non-English locale exists; document the plan in the README.
8. **No trackers, no cookies, no third-party scripts, no third-party requests from the browser.** Therefore no cookie banner is needed. Do not add analytics. Do not add reCAPTCHA or any other third-party anti-spam service.
9. **Honesty in copy:** no customer logos, no "trusted by", no partner names, no "ESA-backed", no "certified", no usage statistics, no testimonials, no claims of accuracy percentages. No hype words: "revolutionary", "cutting-edge", "game-changing", "AI-powered", "seamless", "unlock", "leverage", "empower", "magic".
11. **No real sensitive locations.** Demo visuals use abstract or clearly fictional places. Never show real military sites, real borders in dispute, or real people.
10. Do not use logos, names or images of real energy companies or real facilities.

---

## 3. Brand

### Colour tokens

Define all colours as CSS custom properties on `:root`. Start from these and adjust only if a contrast check fails:

| Token | Value | Use |
|---|---|---|
| `--ink` | `#0C003D` | Primary brand colour, headings, primary buttons, dark sections |
| `--ink-80` | `#3D3364` | Body text on light backgrounds |
| `--ink-60` | `#6D6690` | Secondary text, captions |
| `--line` | `#E4E1EE` | Borders, dividers |
| `--tint` | `#F5F4F9` | Alternate section background |
| `--paper` | `#FFFFFF` | Page background |
| `--on-ink` | `#FFFFFF` | Text on `--ink` |
| `--on-ink-muted` | `#B9B4D1` | Secondary text on `--ink` |
| `--observed` | `#FFFFFF` on dark, `--ink` on light | Label for observed facts in the demo brief |
| `--inferred` | `#9A6400` | Label for model inferences in the demo brief |

All text must meet WCAG 2.2 AA contrast. Do not introduce other accent colours. The brand is calm and precise: an intelligence briefing, not a sci-fi poster.

### Typography

- Headings: Mozilla Headline. Body, UI and captions: Mozilla Text.
- Fluid type scale with `clamp()`. H1 roughly 40px on mobile to 72px on wide screens. Body 17-18px, line height about 1.55, measure 60-72 characters.
- Numbers in tables and the demo brief should use tabular figures if the font supports them.

### Logo and motif

- The logo is a mark of four parallel wave lines plus the wordmark "Tiefer". The four lines stand for the four data layers Tiefer fuses: optical, radar, thermal and night lights. Use this meaning in the product section. Use the SVG files provided (section 7). Never redraw, recolour (other than the dark and white versions provided), stretch or add effects.
- Minimum clear space around the logo: the height of one wave line.
- The four-wave motif may be used as a subtle, very low-contrast background element in the contact section (for example, long horizontal waves at 4-8% opacity).
- **Orbital motif (space identity).** Dark sections (`--ink`) carry a restrained orbital visual language, drawn as inline SVG: thin dashed orbit arcs, a small satellite glyph, a narrow observation swath (a thin wedge from the satellite to a ground point), the curve of the Earth's limb at the bottom edge, and at most a few dozen tiny static star points at low opacity. Lines 1-1.5 px in `--on-ink-muted` at 15-45% opacity. No glow, no gradients, no lens flares, no rockets, no planets other than Earth, no stock photography.
- Small technical labels in uppercase with letter spacing (for example `SAR`, `OPTICAL`, `CLOUD 82%`, `PASS 034`, `BAKU 40.41 N 49.87 E`) may annotate these visuals, like an instrument readout. Never use them for claims or real facility data.
- Motion: at most one slow, subtle animation in the hero (for example the satellite glyph moving a short distance along its orbit arc once, then resting). No animation that loops forever. Respect `prefers-reduced-motion` (no motion at all).

### Tone of voice

Like an intelligence briefing: calm, precise, sourced. Short sentences. Space vocabulary where it is accurate (orbit, pass, revisit, tasking, SAR), never where it would exaggerate. Never imply that Tiefer operates its own satellites, and never promise that the AI sees everything.

---

## 4. Tech stack

### Principle

One language, one binary, no Node toolchain. The backend is Go, and the frontend is server-rendered HTML from the same Go program. This keeps the project small enough for one founder to maintain, and the future dashboard can grow inside the same codebase.

### Backend

- **Go** (latest stable, at least 1.23). Standard library first: `net/http` with the Go 1.22+ pattern router (`mux.HandleFunc("GET /", ...)`), `html/template`, `embed`, `log/slog`, `net/smtp` or `crypto/tls` as needed.
- No web framework (no Gin, Echo, Fiber). Add a third-party module only when the standard library clearly cannot do the job, and justify each one in the README.
- All templates, CSS, fonts and brand assets are embedded with `//go:embed`, so the result is **one self-contained binary**.
- Templates are parsed once at startup. Pages are rendered from the content structs, not from hardcoded strings.
- Configuration only through environment variables, loaded and validated at startup in `internal/config` (see section 11). Fail fast with a clear message if a required variable is missing in production mode.
- Graceful shutdown on SIGINT and SIGTERM. Sensible server timeouts (`ReadHeaderTimeout`, `ReadTimeout`, `WriteTimeout`, `IdleTimeout`).
- `GET /healthz` returns `200 ok` for the hosting platform.
- Structured logging with `log/slog`. **Do not log full IP addresses** (data minimisation under the GDPR and Azerbaijan's Law on Personal Data). If an IP is needed for rate limiting, keep it in memory only and never write it to logs.

### Frontend

- Server-rendered HTML with Go `html/template`. Split into a base layout and one partial per section (header, hero, question box, problem, product layers, how, demo brief, use cases, principles, roadmap, name, contact, footer).
- Plain CSS with custom properties in one stylesheet (`web/static/css/site.css`). No Tailwind, no CSS framework, no build step for CSS.
- Vanilla JavaScript only where required: the mobile navigation toggle and progressive enhancement of the contact form. Under 3 KB total. Everything must work with JavaScript disabled.
- Static files served with long cache headers and content-hashed file names (compute the hash at startup from the embedded file and expose a template function such as `{{ asset "css/site.css" }}`).

### Content

- All copy lives in `internal/content/en.go` as typed Go structs (for example `content.Site`, `content.Hero`, `content.TimelineItem`). Templates only read from these structs.
- Adding a language later means adding for example `internal/content/az.go` and an `/az/` route that renders the same templates with that struct. Design the routing, the `lang` attribute and `hreflang` links for this now.

### HTTP and security

- Security headers on every response: a strict `Content-Security-Policy` (`default-src 'self'`; no inline scripts; no external origins), `Strict-Transport-Security`, `X-Content-Type-Options: nosniff`, `Referrer-Policy: strict-origin-when-cross-origin`, `Permissions-Policy` disabling unused features, `X-Frame-Options: DENY`.
- Gzip or Brotli compression for HTML, CSS, JS and SVG.
- Custom 404 page in the site design.
- `robots.txt`, `sitemap.xml` and `site.webmanifest` served by Go (sitemap generated from the route list and `SITE_URL`).

### Build and run

- `Makefile` with: `make run`, `make build`, `make test`, `make lint` (`go vet` plus `staticcheck` if available), `make check` (runs the text check in section 10), `make docker`.
- Multi-stage `Dockerfile`: build with the official Go image, run on a minimal image (distroless static or scratch) as a non-root user, exposing one port from the `PORT` variable.
- The binary must run anywhere a container runs, on a host in Azerbaijan or in the EU. The code must not depend on any provider. In the README, list both options neutrally and note that the hosting choice (and whether personal data from the contact form may be processed abroad) must be confirmed with a local lawyer, because Azerbaijan's Law on Personal Data has registration and cross-border transfer rules.

### Suggested structure

```
/
  cmd/tiefer-web/main.go
  internal/
    config/           (env loading and validation)
    content/          (en.go, later de.go)
    server/           (routes, handlers, middleware, security headers, rate limiting)
    contact/          (form validation, spam protection, mail sending)
    render/           (template loading, asset hashing, template funcs)
  web/
    templates/        (layout.html, partials/*.html, pages/index.html, legal.html, privacy.html, acceptable-use.html, 404.html)
    static/
      css/site.css
      js/site.js
      fonts/          (WOFF2 files + OFL.txt)
      brand/          (logo SVGs)
      favicon.svg, favicon.ico, apple-touch-icon.png, icon-192.png, icon-512.png, og-image.png
    embed.go          (//go:embed directives)
  scripts/            (only if something cannot live in Go tests)
  LICENSE
  NOTICE.md
  Makefile
  Dockerfile
  .env.example
  README.md
```

---

## 5. Page structure and copy (index page)

Single scrolling page. Sticky header. Sections in this order. Use the copy exactly.

### 5.1 Header

- Left: logo (`tiefer-logo.svg` on light, `tiefer-logo-white.svg` while over the dark hero), links to top.
- Right navigation: `Product` (#product), `How it works` (#how), `Use cases` (#use-cases), `Principles` (#principles), `Contact` (#contact).
- Primary button: `Request access` (links to #contact).
- On mobile: logo plus a menu button that opens an accessible disclosure menu.

### 5.2 Hero (dark)

- Background `--ink` with the orbital motif (section 3): an orbit arc across the upper right, a small satellite glyph, a thin observation swath down to the Earth's limb along the bottom edge, sparse static stars. Text on the left, fully legible over the visual.
- Small technical label, top right on desktop only: `BAKU 40.41 N 49.87 E`
- Eyebrow: `Space intelligence you can ask`
- H1: `See deeper.`
- Lead: `Ask a question about any place on Earth in plain language. Tiefer finds and fuses optical, radar, thermal and night-light satellite data, tasks new imagery when the archive is not enough, and returns a sourced intelligence brief.`
- Below the lead, a **question box** styled like a real prompt field but not interactive (a `<figure>` with a caption, not an `<input>`), with this example typed in: `Which ports on the eastern Caspian coast saw more vessel activity in the last 14 days?` and a small caption under it: `Example question`.
- Primary button: `Request access` (#contact). Secondary link: `See how it works` (#how).

### 5.3 Problem

- H2: `Satellite data is everywhere. Answers are not.`
- Intro: `Getting one answer from space still takes days: find the right images, buy them, clean them, align optical, radar and thermal layers by hand, look for change, write it up. The people who need the answer rarely have the specialists to do it.`
- Four items (grid of 4 on desktop, 2 on tablet, 1 on mobile):
  1. `Scattered catalogues` / `Every satellite has its own archive, format, resolution and revisit.`
  2. `Clouds and darkness` / `Optical images fail under cloud and at night. Radar and thermal fill the gap, but few teams can combine them.`
  3. `Specialist work` / `Cleaning, aligning and analysing imagery needs remote sensing experts most teams do not have.`
  4. `Slow tasking` / `When no recent image exists, ordering a new one is a manual process that takes days.`

### 5.4 Product (id `product`, dark)

- H2: `Four layers. One answer.`
- Body: `The four lines in our logo are the four layers Tiefer fuses. Each keeps its own source. Together they answer the question.`
- Four columns (2 by 2 on mobile), each with a thin line icon drawn in inline SVG and a label in the technical label style:
  1. `OPTICAL` / `What the surface looks like. Sentinel-2, Landsat and commercial imagery.`
  2. `RADAR` / `Sees through cloud and at night. Sentinel-1 and commercial SAR.`
  3. `THERMAL` / `Heat: fires, industrial activity, surface temperature.`
  4. `NIGHT LIGHTS` / `Where light appears or disappears after dark.`
- Keep these entries in `internal/content/en.go`.

### 5.5 How it works (id `how`)

- H2: `From question to brief in five steps`
- Five numbered steps, visually connected left to right on desktop (a simple line, not arrows made of characters), vertical on mobile:
  1. `Ask` / `Tiefer turns your question into a place, a time window and the signs to look for. If something is unclear, it asks.`
  2. `Plan` / `It picks the sensors that can answer. Cloudy area? Radar first.`
  3. `Fuse` / `It finds, cleans and aligns the imagery, then runs change and object detection.`
  4. `Task` / `No recent image? It prepares an order for a commercial satellite, within your budget rules. Anything beyond them waits for your approval.`
  5. `Brief` / `You get a map, the findings, a confidence level for each, the source images, and what could not be seen.`

### 5.6 Demo brief

- H2: `Every finding shows its source.`
- Body: `An answer you cannot check is not intelligence. Tiefer labels what a satellite observed separately from what a model inferred, links every finding to its source image, and never presents generated pixels as evidence.`
- Next to the text, an **illustrative brief card** built in HTML and CSS (not an image). Visible label at the top: `Illustrative example. Not real data.` Content:

| Field | Value |
|---|---|
| Question | `Vessel activity at Port A, last 14 days` |
| Sensors | `Radar (SAR), optical` |
| Passes used | `9 radar, 3 optical (4 optical rejected: cloud)` |
| Finding 1 | `Observed` / `Vessels at berth rose from 6 to 11 between day 1 and day 14.` |
| Finding 2 | `Inferred` / `Likely increase in loading activity. Confidence: medium.` |
| Not seen | `No optical view on days 5 to 9 (cloud). Radar only.` |

  Style the `Observed` and `Inferred` labels as small pills (`--observed` and `--inferred` tokens). Under the table, a thin 14-day strip of small marks: filled for radar passes, outlined for optical passes, crossed for rejected cloudy passes. Decorative only (`aria-hidden="true"`); the table carries the information. `Port A` is fictional; do not name or map any real port.

### 5.7 Use cases (id `use-cases`)

- H2: `Built for people who decide on what happens on the ground`
- Three cards:
  1. `Commodity traders` / `Ports, terminals, storage and crops: see supply move before the statistics do.` / example question in small text: `How full are the tanks at this terminal compared with last month?`
  2. `Insurers` / `Floods, fires and damage mapped in hours, with the images to back each claim.` / example: `Which insured sites in this district were under water on Tuesday?`
  3. `Governments` / `Sovereign deployments that keep sensitive questions and data inside the country.` / example: `What changed at this infrastructure site since January?`
- Under the cards, one line: `We start in the Caspian, Caucasus, Black Sea and Central Asia region. Tiefer answers in English, Azerbaijani, Turkish and Russian.`

### 5.8 Principles (id `principles`)

- H2: `How we work`
- Six short statements in a 3 by 2 grid (1 column on mobile):
  1. `Observed is not inferred.` / `Every brief labels the two separately.`
  2. `No invented pixels.` / `Generated imagery is never shown as evidence. If it was cloudy, we say so.`
  3. `Every finding is traceable.` / `Sensor, date, source image and processing steps travel with it.`
  4. `You control the spend.` / `Automatic tasking runs only within your rules. Anything else waits for approval.`
  5. `Limits are stated.` / `Every brief says what could not be seen.`
  6. `People are not targets.` / `Tiefer is not used to track individuals. Customers are screened.`

### 5.9 Roadmap

- H2: `Where we are`
- A horizontal track (vertical on mobile) with five nodes; the first marked `Now`:
  1. `Demo` / `Three question types on free satellite data: port activity, flood mapping, land change.`
  2. `Pilots` / `First paying teams, watchlists and weekly briefs.`
  3. `Tasking` / `Commercial imagery ordered by Tiefer, with human approval first.`
  4. `Sovereign` / `On-premise deployments for governments and large companies.`
  5. `More regions` / `The same analyst, new places and new sensors.`

### 5.10 Name

- Dark section (`--ink` background, white text), short, with a quiet Earth-limb curve along the bottom.
- Large line: `Tiefer is German for deeper.`
- Body: `We look past the surface of an image to what is actually happening on the ground, and we show you how we know.`
- The four-wave mark (white) may appear here, large and quiet.

### 5.11 Contact (id `contact`)

- H2: `We are early, and we are listening.`
- Body: `We are looking for our first pilot teams: traders, insurers, analysts and public institutions with questions about the physical world. Tell us the question you would ask first.`
- A short contact form handled by the Go backend (`POST /contact`):
  - Fields: `Name` (required), `Work email` (required), `Organisation` (optional), `Your role` (optional select: `Trader or analyst`, `Insurer`, `Public sector`, `Investor`, `Other`), `The first question you would ask Tiefer` (required, max 2,000 characters).
  - Checkbox (required): `I agree that Tiefer may use these details to reply to my message. See the privacy notice.` with a link to `/privacy`.
  - Submit button: `Request access`.
  - Success message: `Thank you. We will reply within two working days.`
  - Error message: `Something went wrong. Please email us at` followed by the contact email as a link.
- The form must work without JavaScript (normal POST, then redirect to `/?sent=1#contact` using the Post/Redirect/Get pattern). With JavaScript, submit via `fetch` and show the message inline without reloading.
- Server side: validate and length-limit every field, reject header injection in name and email, CSRF protection (a signed token in a hidden field, verified on submit), a hidden honeypot field, a minimum fill time check, and an in-memory rate limit per client (for example 5 submissions per hour). No third-party anti-spam.
- Delivery: send the message by email through SMTP settings from environment variables (section 11). Set `Reply-To` to the sender. **Do not store submissions in a database.** Do not log message content.
- If SMTP is not configured, the form is not rendered and the section shows only the `Email us` button.
- Below the form, always show: `Prefer email?` with a `mailto:` link to `CONTACT_EMAIL`, and `Follow Tiefer on LinkedIn` linking to `LINKEDIN_URL`.

### 5.12 Footer

- Logo (small), `© 2026 Tiefer`, links: `Legal notice`, `Privacy`, `Acceptable use`, `LinkedIn`, `Source code` (link to `REPO_URL`, only rendered when it is set).
- One quiet line: `See deeper. Built in Baku.`

---

## 6. Legal pages

The company is registered in Azerbaijan and serves EU customers. Create `/legal`, `/privacy` and `/acceptable-use` with a clean text layout that uses the same header and footer. Do not use German-specific terms such as "Impressum" or references to German law.

- **Legal notice (`/legal`):** company information that EU business visitors expect: legal name and legal form (for example `[COMPANY_LEGAL_NAME] MMC`, a limited liability company under the laws of the Republic of Azerbaijan), registered address in Baku, tax identification number (VÖEN), state registration details, managing director, contact email. Fill every value from config, with clearly marked placeholders like `[COMPANY_LEGAL_NAME]` as defaults. Do not invent any data.
- **Privacy (`/privacy`):** a short, plain-English privacy notice that is true for this site as built: server-rendered site, no cookies, no analytics, no third-party requests, fonts self-hosted, server logs without full IP addresses, the contact form (which fields, purpose, that data is sent by email and not stored in a database, retention placeholder), the SMTP provider and hosting provider (with their country) as placeholders. Structure it so it can satisfy both Azerbaijan's Law on Personal Data (11 May 2010) and, because the site is aimed at people in the EU, the GDPR: controller identity and contact, purposes, legal basis placeholder, recipients, international transfers placeholder, retention, data subject rights, and how to complain. Do not write legal conclusions; use placeholders where a lawyer must decide.
- **Acceptable use (`/acceptable-use`):** a short, plain-English policy based on the product principles: no tracking or identifying individuals, no use that violates human rights or applicable sanctions and export controls, customer screening, the right to suspend accounts. Placeholders where a lawyer must decide.
- Add a visible note at the top of all three pages: `Placeholder text. To be reviewed before launch.` Render this note whenever `LEGAL_REVIEWED` is not `true`, so it cannot go live unnoticed.

---

## 7. Assets (provided)

These files are **already in the repository** at the paths shown (under `web/static/`). Do not modify them. The folder `docs/` contains this brief (`PROMPT.md`) and the founder's product document in Azerbaijani (`Tiefer_2026_az.md`) for context; the site copy in this brief takes priority.

| File | Use |
|---|---|
| `brand/tiefer-logo.svg` | Header and footer on light backgrounds (`#0C003D`) |
| `brand/tiefer-logo-white.svg` | Logo on dark backgrounds |
| `brand/tiefer-mark.svg` | Four-wave mark on light backgrounds |
| `brand/tiefer-mark-white.svg` | Four-wave mark on dark backgrounds, name section |
| `brand/tiefer-logo-black.svg`, `brand/tiefer-mark-black.svg` | Monochrome versions for print and press. Not used on the website by default |
| `brand/tiefer-logo-gray.svg`, `brand/tiefer-mark-gray.svg` | Neutral grey (`#808080`) versions for press kits and partner lists. Not used on the website by default |
| `favicon.svg` | Modern browsers |
| `favicon.ico` | Legacy browsers (16, 32, 48 px), also served at `/favicon.ico` |
| `apple-touch-icon.png` | 180 x 180, also served at `/apple-touch-icon.png` |
| `icon-192.png`, `icon-512.png` | Web app manifest |
| `og-image.png` | 1200 x 630 social preview (absolute URL in meta tags) |

All logo files are the original vector artwork (viewBox `0 0 3669 700` for the full logo, `0 0 690 690` for the mark). Use them inline or as `<img>`; do not re-export or trace them.

Always give the logo an accessible name ("Tiefer"). Decorative uses of the mark get `aria-hidden="true"` and empty alt text.

---

## 8. SEO and metadata

- `<title>`: `Tiefer | Space intelligence you can ask`
- Meta description: `Tiefer is an AI analyst for satellite data. Ask a question in plain language and get a sourced intelligence brief from optical, radar, thermal and night-light imagery. Built in Baku.`
- Open Graph and Twitter card tags using `og-image.png`.
- `lang="en"`, canonical URL built from `SITE_URL`.
- `sitemap.xml`, `robots.txt`, `site.webmanifest` (theme colour `#0C003D`), all served by Go.
- JSON-LD `Organization` with name, url, logo and `sameAs` (LinkedIn). Serve it as a `<script type="application/ld+json">` block that is allowed by the CSP (use a hash or nonce, not `unsafe-inline`). No invented fields.
- One H1 per page, logical heading order.

---

## 9. Accessibility, performance, responsiveness

- WCAG 2.2 AA. Visible focus styles in `--ink` with sufficient offset. Skip link to main content. All interactive elements reachable and usable by keyboard. Mobile menu uses `aria-expanded` and closes on Escape. Form fields have visible labels, error messages are linked with `aria-describedby`, and the success or error message is announced (`role="status"` or `aria-live`).
- Respect `prefers-reduced-motion`. Any motion is subtle and short.
- Light theme is the default and the only required theme. Do not add a dark mode toggle.
- Test at 360, 390, 768, 1024, 1280 and 1440 px widths. No horizontal scrolling at any width. Side padding at least 16 px on mobile.
- Lighthouse targets on mobile: Performance 95+, Accessibility 100, Best Practices 100, SEO 100.
- Preload the two main font files. Use `font-display: swap`. Subset only if it does not remove characters needed for "Türkiye" and other copy.
- Total page weight for the index page under 300 KB excluding fonts. Time to first byte from the Go server under 50 ms locally.

---

## 10. Automated text check

Implement the check as a **Go test** (for example `internal/content/text_test.go` plus a test that renders every page), so `go test ./...` and `make check` run it. It must fail with file name and line number (or page and section for rendered output) if it finds, in the content files, the templates or the rendered HTML of every page:

- U+2014 (em dash), U+2013 (en dash), U+2015 (horizontal bar)
- Any emoji: characters with the Unicode property Extended_Pictographic (use a maintained table or `golang.org/x/text` if needed), plus U+FE0F
- The banned words from rule 9 in section 2 (case-insensitive, whole words), in visible copy only

Also report, as a warning that does not fail the test, every remaining `[PLACEHOLDER]`-style token and every config value that still has its default. Run the check in CI (see section 12).

---

## 11. Configuration (environment variables)

Load and validate in `internal/config`. Provide `.env.example` with every variable, a comment for each, and safe defaults for local development. Never commit real secrets.

| Variable | Default (dev) | Purpose |
|---|---|---|
| `PORT` | `8080` | HTTP port |
| `ENV` | `development` | `development` or `production` |
| `SITE_URL` | `http://localhost:8080` | Canonical base URL |
| `CONTACT_EMAIL` | `hello@example.com` | Public contact address and `mailto:` link |
| `LINKEDIN_URL` | `https://www.linkedin.com/company/tiefer` | LinkedIn page |
| `REPO_URL` | empty | Public source repository, footer link |
| `SMTP_HOST`, `SMTP_PORT`, `SMTP_USER`, `SMTP_PASS` | empty | Mail delivery for the contact form (TLS required) |
| `CONTACT_TO` | empty | Where form messages are delivered |
| `CONTACT_FROM` | empty | Sender address for form mails |
| `CSRF_SECRET` | random in dev | Signing key for form tokens, required in production |
| `LEGAL_*` | placeholders | Company legal name, legal form, address in Baku, VÖEN (tax ID), state registration, managing director, hosting provider and country, SMTP provider and country |
| `LEGAL_REVIEWED` | `false` | Hides the placeholder warning on legal pages when `true` |

---

## 12. Repository hygiene

- `go.mod` with module path from a `MODULE_PATH` you ask me for, or use `github.com/tiefer-labs/web`.
- `gofmt`, `go vet` and `staticcheck` clean.
- Tests for: config validation, contact form validation (valid, missing fields, header injection, honeypot, too fast, rate limit), CSRF token round trip, security headers present on every route, 404 page, every page renders with status 200, and the text check.
- GitHub Actions workflow `.github/workflows/ci.yml`: run `gofmt` check, `go vet`, `staticcheck`, `go test ./...`, and `go build` on every push and pull request.
- `.gitignore`, `.editorconfig`, `.env.example`.

---

## 13. Licence: Mozilla Public License 2.0

- Add the full, unmodified text of the **Mozilla Public License, v. 2.0** as `LICENSE` in the repository root (from `https://www.mozilla.org/en-US/MPL/2.0/`). Set the licence to `MPL-2.0` wherever a machine-readable identifier is used.
- At the top of every Go source file, template, CSS and JS file written for this project, add the standard MPL 2.0 header in that file type's comment syntax:

  ```
  This Source Code Form is subject to the terms of the Mozilla Public
  License, v. 2.0. If a copy of the MPL was not distributed with this
  file, You can obtain one at https://mozilla.org/MPL/2.0/.
  ```

  In HTML templates, use a Go template comment (`{{/* ... */}}`) so the header does not appear in the rendered page.
- **Not covered by the MPL:** the Tiefer name, logos and brand assets in `web/static/brand/`, the favicons and the social image. MPL 2.0 does not grant trademark rights (section 2.3 of the licence). Add `NOTICE.md` explaining this in plain English: the code is MPL 2.0; the Tiefer name and logo are trademarks of the company and may not be used to suggest endorsement or to brand a fork.
- **Fonts:** Mozilla Headline and Mozilla Text stay under the SIL Open Font License 1.1. Keep `OFL.txt` next to the font files and mention it in `NOTICE.md`.
- **Website copy:** state in `NOTICE.md` that the text content in `internal/content/` is also covered by MPL 2.0 unless the founder decides otherwise, and flag this as a decision for me in your final summary.
- List every third-party Go module with its licence in `NOTICE.md`. Only use modules whose licences are compatible with MPL 2.0 (MIT, BSD, Apache 2.0, MPL 2.0 are fine; ask before adding anything else).

---

## 14. Deliverables and acceptance criteria

When you finish:

1. `make run`, `make build`, `make test`, `make lint`, `make check` and `make docker` all work with no errors.
2. `go build` produces one binary that serves the whole site with no files on disk next to it.
3. The index page, `/legal`, `/privacy`, `/acceptable-use` and the 404 page render correctly at all widths listed in section 9.
4. The contact form works with and without JavaScript against a local SMTP catcher (document how to test it, for example with Mailpit in Docker).
5. No em dash, en dash or emoji anywhere (the Go test passes).
6. Fonts are self-hosted, the licence files are present, and there are zero browser requests to third-party domains (verify with a headless browser).
7. All copy comes from `internal/content/en.go`; changing a sentence there changes the site.
8. `LICENSE` contains the full MPL 2.0 text, every source file has the MPL header, and `NOTICE.md` covers brand assets, fonts and dependencies.
9. `README.md` explains: what the project is, how to run it, how to edit copy, how to update the example questions and the roadmap, how to add another language, how to configure SMTP, how to deploy with Docker (Azerbaijani or EU host), the licence, and the list of placeholders still to fill.
10. Take screenshots of the index page at 390 px and 1440 px and review them yourself before telling me you are done. Fix anything that looks unbalanced, cramped or off-brand.
11. At the end, give me a short list of anything you could not do (for example, fonts you could not download) and exactly what I need to provide.
