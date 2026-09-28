# Build the Tiefer website

You are rebuilding, from zero, the public website for **Tiefer**, an early-stage B2B startup. Nothing from any earlier version is reused. This is a business-card website: one main page plus two legal pages, served by a small **Go** backend. There is no product dashboard and no login yet. Build it so it is fast, accessible, easy to edit and easy to extend later with more languages and a product area (the future dashboard will live in the same Go codebase).

The repository is licensed under the **Mozilla Public License 2.0** (see section 13).

Read this whole brief before you start. Where the brief gives exact copy, use it word for word. Do not invent claims, customers, partners, numbers or testimonials.

Two companion documents are part of this specification and are equally binding:

- `docs/HARDENING.md`: security requirements, threat model and security tests. On any security question it wins over this brief.
- `docs/AZURE.md`: hosting on Microsoft Azure (App Service, Front Door with WAF, Key Vault, Container Registry, GitHub OIDC deployment, Bicep).

Three qualities matter more than features: **secure, fast, calm**. When in doubt, remove rather than add.

---

## 1. Context: what Tiefer does

Tiefer is a **space technology startup** based in **Baku, Azerbaijan**. It builds **AI software that runs on Earth observation satellites** (edge AI in orbit). Instead of storing every image and sending gigabytes of raw data to the ground hours later, the satellite analyses its own images in orbit: it filters out cloudy and empty frames, detects events such as wildfires, oil spills, vessels and floods, and sends a small alert packet (tens of kilobytes: coordinates, event type, confidence, a small image chip) through the first available link.

Tiefer does not build hardware. Its software runs on the onboard compute that operators already fly (for example NVIDIA Jetson, Intel Movidius or FPGA-based processors).

The onboard pipeline has four stages, matching the four lines in the logo:

1. **Filter**: check each frame for cloud and quality the moment it is captured. Frames that are not useful are compressed and kept on board, not sent.
2. **Detect**: run computer vision models for specific events on the useful frames.
3. **Alert**: turn a detection into a small alert packet and send it to the ground through the first available channel.
4. **Update**: replace or add models in orbit with small, signed update packages.

Two ground products complete the system: **Tiefer Lab** (train, quantise and test models in a "virtual orbit" simulator and on flight-like hardware) and **Tiefer Ground** (the same models running at the ground station, so operators get value before any satellite carries Tiefer).

Customers: satellite operators and space agencies (starting with Azerbaijan and the region), satellite manufacturers and integrators, and hosted-compute platforms that run third-party apps in orbit. End users (emergency services, environmental agencies, energy and maritime authorities) benefit through the operators.

Product principles that the site must reflect:

1. Nothing is deleted blindly. Filtered frames are compressed and kept; the operator sets the rules.
2. Observed and inferred are always labelled separately, with a confidence value.
3. Generated pixels are never sent as evidence.
4. Every model update is signed, versioned and reversible.
5. Performance claims are measured and published with their method.
6. Alerts support human decisions; they do not replace them.
7. Tiefer is not used to track individuals. Customers are screened; export control and sanctions rules are followed.

The domain is **tiefer.space** (production URL `https://tiefer.space`; redirect `www.tiefer.space` to it). "Tiefer" is German for "deeper". Brand face: **See deeper.** Product promise: **AI that runs on the satellite.** Pitch line: **From passive cameras to orbital analysts.** Visitors will be satellite operators, space agency engineers, integrators, investors and accelerator reviewers, mostly reading in English.

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
10. Do not use logos, names or images of real energy companies or real facilities.
11. **No real sensitive locations.** Demo visuals use abstract or clearly fictional places. Never show real military sites, real borders in dispute, or real people.
12. **No unmeasured performance numbers.** Do not print accuracy, latency, bandwidth savings or cost reductions as facts. Where the brief shows numbers, they sit in a card labelled as illustrative.
13. **Restraint in design.** The site must look quiet and expensive, never loud. See "Design restraint" in section 3. If an effect draws attention to itself, remove it.
14. **Not a template, not AI-looking.** The site must look designed by a person for this company, not generated. See "Not a template" in section 3. This is an acceptance criterion, not a preference.
15. **Rebuild from zero.** Any earlier site in the repository is removed and nothing from it is reused (see `docs/START.md`, phase 0).

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
| `--observed` | `#FFFFFF` on dark, `--ink` on light | Label for observed facts in the demo alert |
| `--inferred` | `#9A6400` | Label for model inferences in the demo alert |

All text must meet WCAG 2.2 AA contrast. Do not introduce other accent colours. The brand is calm and precise: flight software, not a sci-fi poster.

### Typography

- Headings: Mozilla Headline. Body, UI and captions: Mozilla Text.
- Fluid type scale with `clamp()`. H1 roughly 40px on mobile to 72px on wide screens. Body 17-18px, line height about 1.55, measure 60-72 characters.
- Numbers in tables and the demo alert should use tabular figures if the font supports them.

### Logo and motif

- The logo is a mark of four parallel wave lines plus the wordmark "Tiefer". The four lines stand for the four onboard stages: Filter, Detect, Alert, Update. Use this meaning in the product section. Use the SVG files provided (section 7). Never redraw, recolour (other than the versions provided), stretch or add effects.
- Minimum clear space around the logo: the height of one wave line.
- The four-wave motif may be used as a subtle, very low-contrast background element in the contact section (for example, long horizontal waves at 4-8% opacity).
- **Orbital motif (space identity).** The hero (and only the hero) carries a restrained orbital visual language, drawn as inline SVG: thin dashed orbit arcs, a small satellite glyph, a narrow observation swath (a thin wedge from the satellite to a ground point), the curve of the Earth's limb at the bottom edge, and at most a few dozen tiny static star points at low opacity. Lines 1-1.5 px in `--on-ink-muted` at 15-45% opacity. No glow, no gradients, no lens flares, no rockets, no planets other than Earth, no stock photography.
- Small technical labels in uppercase with letter spacing (for example `FRAME 0412`, `CLOUD 82%`, `KEPT ON BOARD`, `ALERT 48 KB`, `BAKU 40.41 N 49.87 E`) may annotate these visuals, like telemetry. Never use them for claims or real facility data.
- Motion: at most one slow, subtle animation in the hero (for example the satellite glyph moving a short distance along its orbit arc once, then resting). No animation that loops forever. Respect `prefers-reduced-motion` (no motion at all).

### Design restraint

The audience is engineers at space agencies and satellite operators. They trust sites that look like good technical documentation, not like a launch event. Rules:

- White and `--tint` backgrounds carry most of the page. At most three dark (`--ink`) bands: the hero, the product stages and the name section.
- The orbital motif appears in the hero only, plus the quiet four-wave mark in the name section and the waves in the contact section. Nowhere else.
- One accent colour (`--ink`) and its tints. No gradients, glows, blurs, glass effects, neon, shadows (except a 1 px border or a very soft shadow on the demo alert card), noise textures or background videos.
- No parallax, no scroll-triggered reveal animations, no scroll hijacking, no animated counters, no carousels, no marquee, no typewriter effects, no custom cursors, no particles, no 3D, no globe.
- Motion: only the single hero animation from the motif rules, and hover or focus transitions of colour, opacity or underline under 200 ms.
- Corners: square or at most 4 px radius. Borders 1 px in `--line`.
- Buttons: one primary style (solid `--ink`, white text) and one secondary style (text link with underline or 1 px outline). Never more than one primary button in view.
- Type does the work: large Mozilla Headline headings with generous spacing, body in Mozilla Text at a comfortable measure. Use at most two weights per family on the page (for example Headline 500 and 600, Text 400 and 500).
- Spacing on an 8 px grid. Section padding generous (roughly 96 to 144 px on desktop, 64 to 88 px on mobile). Content width at most 1,120 px, text columns at most 72 characters.
- Icons: none, or simple 1.5 px line icons drawn in inline SVG, one size, one colour.
- No stock photos, no illustrations of people, no AI-generated imagery.

### Not a template

Visitors from space agencies have seen hundreds of generated startup sites. The moment a page looks like one, it loses their trust. Avoid every pattern below, and review the screenshots against this list.

**Patterns that make a site look AI-generated (never use):**

- Purple or blue gradients, gradient text, glowing orbs, blurred colour blobs, aurora or mesh backgrounds, dotted or grid backgrounds.
- Glass cards, cards with large radius and soft drop shadow everywhere, gradient borders, bento grids.
- The repeated block "centred heading, grey subtitle, three cards with an icon in a rounded square". Any section that could be swapped with another site's section without anyone noticing.
- Pill badges ("New", "Beta", "Now in orbit"), sparkles, stars as decoration, emoji or icon bullets.
- Everything centred. Every section with the same height, padding and layout rhythm.
- Generic icon sets used as decoration. Icons that only repeat what the label says.
- Hover effects that scale, lift or glow. Animated gradients. Buttons with arrows that slide.
- Stats bars ("10x faster"), logo walls, testimonial sliders, "How it works" as three numbered circles.
- Stock phrases in UI text: "Get started", "Learn more", "Discover", "Transform", "Next-generation", "Built for the future". Use the copy in section 5 only.
- Inter, system UI fonts or any font other than the two Mozilla fonts for visible text.

**What to do instead:**

- Think of a well-made mission document or a precision instrument manual, not a SaaS landing page. International Typographic Style: a strict grid, strong type hierarchy, hairline rules, generous white space, left-aligned text.
- A 12-column grid on desktop. Use asymmetry on purpose: for example a heading in columns 1 to 5 and body text in columns 7 to 12, a narrow label column for section numbers, content that does not always start at the same column.
- Give each section its own composition that fits its content: the problem as one strong statement in large type, the four stages as a precise table-like sequence, how it works as a single technical diagram, the alert card as a real data object, principles as a numbered list set like a specification, the roadmap as a timeline on hairlines. No two neighbouring sections share a layout.
- Section labels like a document: small uppercase labels with a number (`02 / PRODUCT`), aligned to the grid.
- Real details carry the identity: tabular figures, precise units, telemetry labels on the hero visual, careful line lengths, optical alignment of the logo, consistent 1 px hairlines in `--line`.
- Typography does the heavy lifting. Big headlines set tight (line height about 1.05, slight negative tracking), body text relaxed. Hanging punctuation where the browser supports it. No orphans in headlines (use `text-wrap: balance` for headings and `text-wrap: pretty` for paragraphs).
- Colour: mostly white, `--tint` and `--ink`. Contrast comes from scale and weight, not from colour.
- Every visual (hero orbit, capture-to-alert flow, size comparison bars) is drawn for this site in inline SVG, and the alert card (HTML and CSS) uses the same label style, so they read as one system with one line weight.

**Review question for every screenshot:** could this section appear unchanged on another company's site? If yes, redesign it.

### Tone of voice

Like a flight software engineer: calm, precise, measurable. Short sentences. Space vocabulary where it is accurate (orbit, pass, downlink, onboard, payload), never where it would exaggerate. Never imply that Tiefer builds or operates satellites, and never promise real-time alerts without saying they depend on the communication link.

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

- Server-rendered HTML with Go `html/template`. Split into a base layout and one partial per section (header, hero, problem, product stages, how it works, demo alert, use cases, principles, roadmap, name, contact, footer).
- Plain CSS with custom properties in one stylesheet (`web/static/css/site.css`). No Tailwind, no CSS framework, no build step for CSS.
- Vanilla JavaScript only where required: the mobile navigation toggle and progressive enhancement of the contact form. Under 3 KB total. Everything must work with JavaScript disabled.
- Static files served with long cache headers and content-hashed file names (compute the hash at startup from the embedded file and expose a template function such as `{{ asset "css/site.css" }}`).

### Content

- All copy lives in `internal/content/en.go` as typed Go structs (for example `content.Site`, `content.Hero`, `content.TimelineItem`). Templates only read from these structs.
- Adding a language later means adding for example `internal/content/az.go` and an `/az/` route that renders the same templates with that struct. Design the routing, the `lang` attribute and `hreflang` links for this now.

### HTTP and security

- Everything in `docs/HARDENING.md`: server limits, exact security headers and CSP, routing rules, template safety, contact form defences, proxy trust, secrets, logging, supply chain, container and security tests. Treat it as part of this section.
- Gzip compression from the Go server (standard library `compress/gzip`) for HTML, CSS, JS, SVG, XML, JSON and plain text. Precompress embedded static text files once at startup and serve the stored bytes; compress rendered HTML per response. Do not add a Brotli module.
- Custom 404 page in the site design.
- `robots.txt`, `sitemap.xml` and `site.webmanifest` served by Go (sitemap generated from the route list and `SITE_URL`).

### Build and run

- `Makefile` with the targets in `docs/HARDENING.md` section 13, plus `make run` and `make build`.
- Multi-stage `Dockerfile` as specified in `docs/HARDENING.md` section 12.
- Production hosting is **Microsoft Azure** as specified in `docs/AZURE.md`. The code itself stays provider neutral: Azure behaviour is enabled only through configuration, and the container runs on any host. In the README, note that the hosting region (and whether personal data from the contact form may be processed outside Azerbaijan) must be confirmed with a local lawyer, because Azerbaijan's Law on Personal Data has registration and cross-border transfer rules.

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
  infra/              (Bicep, see docs/AZURE.md)
  tools/go.mod        (pinned development tools, see docs/HARDENING.md section 11)
  .github/workflows/  (ci.yml, codeql.yml, deploy.yml)
  .github/dependabot.yml
  docs/SECURITY-DECISIONS.md
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

Single scrolling page. Sections in this order. Use the copy exactly.

### 5.1 Header

- Left: logo, links to top. Keep the header simple: it sits on the dark hero at the top of the page (white logo, `tiefer-logo-white.svg`) and is not sticky; a sticky, colour-changing header needs scroll scripts and adds visual noise. Legal pages use a light header with `tiefer-logo.svg`.
- Right navigation: `Product` (#product), `How it works` (#how), `Use cases` (#use-cases), `Principles` (#principles), `Contact` (#contact).
- Primary button: `Talk to us` (links to #contact).
- On mobile: logo plus a menu button that opens an accessible disclosure menu.

### 5.2 Hero (dark)

- Background `--ink` with the orbital motif (section 3): an orbit arc across the upper right, a small satellite glyph, a thin observation swath down to the Earth's limb along the bottom edge, sparse static stars. Near the satellite glyph, a tiny label `ALERT 48 KB` travelling down a thin dashed line to the Earth's limb (static, or one short animation that respects reduced motion). Text on the left, fully legible over the visual.
- Small technical label, top right on desktop only: `BAKU 40.41 N 49.87 E`
- Eyebrow: `Edge AI in orbit`
- H1: `See deeper.`
- Lead: `Tiefer builds AI software that runs on Earth observation satellites. It filters out useless pixels in orbit, detects events on board, and sends kilobyte-sized alerts to the ground instead of gigabytes of raw imagery.`
- Primary button: `Talk to us` (#contact). Secondary link: `See how it works` (#how).

### 5.3 Problem

- H2: `Satellites see more than they can send.`
- Intro: `An Earth observation satellite stores what it captures and sends it down only when it passes over a ground station. By the time the image is processed, a fire has spread or a spill has drifted. And about two thirds of Earth's surface is under cloud at any moment, so much of what is sent is never useful.`
- The intro is the section's one strong statement, set larger than body text across most of the grid width.
- Below it, three items as three narrow text columns under a single hairline, each starting with its number; stacked on mobile. Plain typography: no cards, no boxes, no icons. (The product section below uses rows, so the two neighbours differ.)
  1. `Hours of delay` / `Images wait on board for the next ground contact, then wait again for processing.`
  2. `Wasted downlink` / `Cloudy and empty frames use the same scarce bandwidth as the ones that matter.`
  3. `Passive cameras` / `The satellite captures everything and understands nothing until humans look.`

### 5.4 Product (id `product`, dark)

- H2: `Four stages. One onboard system.`
- Body: `The four lines in our logo are the four stages Tiefer runs on the satellite. Each can be tested on its own. Together they turn a camera into an analyst.`
- Four rows set like a specification table, separated by 1 px hairlines: a number (`01` to `04`), the stage name in the technical label style, and the sentence. On wide screens the rows sit in the right-hand columns of the grid, with the H2 and body on the left. No icons, no cards:
  1. `FILTER` / `Checks every frame for cloud and quality as it is captured. Useless frames are compressed and kept, not sent.`
  2. `DETECT` / `Runs computer vision models for the events that matter: fires, oil spills, vessels, floods.`
  3. `ALERT` / `Sends a small packet with location, event type, confidence and an image chip through the first available link.`
  4. `UPDATE` / `New or better models reach orbit as small, signed updates. No new satellite needed.`
- Under the rows, one line: `Tiefer is software. It runs on the onboard compute operators already fly: GPU, VPU or FPGA.`
- Keep these entries in `internal/content/en.go`.

### 5.5 How it works (id `how`)

- H2: `From capture to alert, on board`
- A horizontal flow (vertical on mobile) with five nodes connected by a thin line, drawn in inline SVG, not ASCII:
  1. `Capture` / `The sensor takes a frame.`
  2. `Filter` / `Cloudy or empty? Compress and keep on board.`
  3. `Detect` / `Useful frame? Run the event models.`
  4. `Alert` / `Event found? Build a small alert packet.`
  5. `Downlink` / `Send it through the first available channel to the operator.`
- Note under the flow: `How fast an alert reaches the ground depends on the communication link: the next ground station pass, a relay, or a wider ground network. Tiefer makes the alert small enough to use whichever comes first.`
- Then a short sub-block, `On the ground too`, as two text columns divided by a vertical hairline (stacked with a horizontal hairline on mobile). No card boxes:
  1. `Tiefer Lab` / `Train, quantise and test models in a virtual orbit simulator and on flight-like hardware before anything flies.`
  2. `Tiefer Ground` / `Run the same models at the ground station today. When your satellites carry Tiefer, the same code moves to orbit.`

### 5.6 Demo alert

- H2: `An alert, not an archive.`
- Body: `When Tiefer finds an event, the satellite does not send the whole scene. It sends what an operator needs to act, and keeps the rest on board until someone asks for it. Every packet says what was observed and what was inferred.`
- Next to the text, an **illustrative alert packet card** built in HTML and CSS (not an image). Visible label at the top: `Illustrative example. Not real data.` Content:

| Field | Value |
|---|---|
| Event | `Wildfire` |
| Location | `Sector 12` (fictional) |
| Detected | `On board, seconds after capture` |
| Confidence | `High` |
| Packet size | `48 KB` (image chip and metadata) |
| Full scene | `Kept on board. Downlink on request.` |
| Cloud in frame | `12%` |
| Observed | `Active fire front, about 1.2 km long.` |
| Inferred | `Likely spread to the north-east. Confidence: medium.` |

  Style `Observed` and `Inferred` as small square-cornered data labels (`--observed` and `--inferred` tokens), like field tags on a technical form, not rounded pills. Beside or under the table, a simple size comparison drawn as two bars: a long bar labelled `Full scene` and a tiny bar labelled `Alert`, with no numbers on the full-scene bar. Decorative only (`aria-hidden="true"`); the table carries the information. `Sector 12` is fictional; do not map any real place.

### 5.7 Use cases (id `use-cases`)

- H2: `Built for events where hours matter`
- Four text blocks in a 2 by 2 arrangement on desktop (1 column on mobile), divided by hairlines that form a cross, each with a small label (`EVENT 01` to `EVENT 04`) above its title. No card boxes, no icons, no images:
  1. `Wildfires` / `Spot active fires in forests and grasslands and alert emergency services while the fire is still small.`
  2. `Oil spills` / `Flag slicks around offshore platforms and shipping lanes for rapid response.`
  3. `Maritime awareness` / `Detect vessels at sea; match them with AIS on the ground to find ships that are not broadcasting.`
  4. `Floods and reservoirs` / `Track rising water and reservoir levels during flood season.`
- Under the blocks, one line: `We start with operators in Azerbaijan and the wider region, and work in English, Azerbaijani, Turkish and Russian.`

### 5.8 Principles (id `principles`)

- H2: `How we build flight software`
- Six statements set like clauses in a specification: labels `P1` to `P6` in a narrow left column, the statement in Mozilla Headline, the explanation below it in Mozilla Text. Two columns of clauses on wide screens, one on mobile. No cards, no icons:
  1. `Nothing is deleted blindly.` / `Filtered frames are compressed and kept. The operator sets the rules.`
  2. `Observed is not inferred.` / `Every alert labels the two separately, with a confidence value.`
  3. `No invented pixels.` / `Generated imagery is never sent as evidence.`
  4. `Signed updates only.` / `Every model in orbit is signed, versioned and can be rolled back.`
  5. `Measured, not claimed.` / `We publish accuracy, latency and power figures with the method behind them.`
  6. `People are not targets.` / `Tiefer is not used to track individuals. Customers are screened.`

### 5.9 Roadmap

- H2: `Ground first. Then orbit.`
- A horizontal track (vertical on mobile) with five nodes; the first marked `Now`:
  1. `Benchmark` / `Cloud filter and three event models, measured on flight-like hardware in a virtual orbit.`
  2. `Ground pilot` / `Tiefer Ground at an operator's ground segment.`
  3. `First flight` / `Tiefer models running in orbit on a hosted compute platform.`
  4. `National satellites` / `Tiefer on board operator satellites, agreed with the manufacturer.`
  5. `Region` / `More operators, more events, radar data on board.`

### 5.10 Name

- Dark section (`--ink` background, white text), short, with a quiet Earth-limb curve along the bottom.
- Large line: `Tiefer is German for deeper.`
- Body: `A camera records the surface. An analyst looks deeper. We are putting the analyst in orbit.`
- The four-wave mark (white) may appear here, large and quiet.

### 5.11 Contact (id `contact`)

- H2: `We are early, and we are listening.`
- Body: `We are looking for our first partners: satellite operators, space agencies, integrators and hosted-compute platforms. Tell us which events you need to catch first.`
- A short contact form handled by the Go backend (`POST /contact`):
  - Fields: `Name` (required), `Work email` (required), `Organisation` (optional), `Your role` (optional select: `Satellite operator`, `Space agency`, `Manufacturer or integrator`, `Hosted compute platform`, `Investor`, `Other`), `What would you want your satellite to detect first?` (required, max 2,000 characters).
  - Checkbox (required): `I agree that Tiefer may use these details to reply to my message. See the privacy notice.` with a link to `/privacy`.
  - Submit button: `Send message`.
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
- **Privacy (`/privacy`):** a short, plain-English privacy notice that is true for this site as built: server-rendered site, no cookies, no analytics, no third-party requests from the browser, fonts self-hosted, application logs without IP addresses. Be honest about the hosting layer: Microsoft Azure (Front Door, its web application firewall and App Service) processes visitors' IP addresses to deliver the site and protect it from attacks, and firewall logs are kept for `[LOG_RETENTION_DAYS]` days (placeholder, default 30), with the legal basis as a placeholder. Then cover the contact form (which fields, purpose, that data is sent by email and not stored in a database, retention placeholder), the SMTP provider and hosting provider (with their country) as placeholders. Structure it so it can satisfy both Azerbaijan's Law on Personal Data (11 May 2010) and, because the site is aimed at people in the EU, the GDPR: controller identity and contact, purposes, legal basis placeholder, recipients, international transfers placeholder, retention, data subject rights, and how to complain. Do not write legal conclusions; use placeholders where a lawyer must decide.
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

- `<title>`: `Tiefer | Edge AI in orbit`
- Meta description: `Tiefer builds AI software that runs on Earth observation satellites: it filters useless pixels in orbit, detects events on board and sends kilobyte-sized alerts instead of gigabytes of raw imagery. Built in Baku.`
- Open Graph and Twitter card tags using `og-image.png`.
- `lang="en"`, canonical URL built from `SITE_URL` (`https://tiefer.space` in production). Redirect `http://` and `www.` to `https://tiefer.space` at the hosting or proxy level, and document this in the README.
- `sitemap.xml`, `robots.txt`, `site.webmanifest` (theme colour `#0C003D`), all served by Go.
- JSON-LD `Organization` with name, url, logo and `sameAs` (`https://www.linkedin.com/company/tiefer/`). Serve it as a `<script type="application/ld+json">` block that is allowed by the CSP (use a hash or nonce, not `unsafe-inline`). No invented fields.
- One H1 per page, logical heading order.

---

## 9. Accessibility, performance, responsiveness

- WCAG 2.2 AA. Visible focus styles in `--ink` with sufficient offset. Skip link to main content. All interactive elements reachable and usable by keyboard. Mobile menu uses `aria-expanded` and closes on Escape. Form fields have visible labels, error messages are linked with `aria-describedby`, and the success or error message is announced (`role="status"` or `aria-live`).
- Respect `prefers-reduced-motion`. Any motion is subtle and short.
- Light theme is the default and the only required theme. Do not add a dark mode toggle.
- Test at 360, 390, 768, 1024, 1280 and 1440 px widths. No horizontal scrolling at any width. Side padding at least 16 px on mobile.
- Lighthouse targets on mobile: Performance 100 (at least 98), Accessibility 100, Best Practices 100, SEO 100.
- Core Web Vitals on Lighthouse mobile: LCP under 1.5 s, CLS under 0.02, TBT 0 ms.

### Performance budgets (enforced by Go tests where possible)

| Item | Budget |
|---|---|
| Index HTML, gzip | 30 KB or less |
| CSS, one file, gzip | 12 KB or less |
| JavaScript, one file, raw | 3 KB or less (and the page works without it) |
| Fonts, both WOFF2 files together | 110 KB or less |
| Images and SVGs loaded on first view, compressed | 40 KB or less |
| Total first view transfer | 200 KB or less |
| Requests on first view | 12 or fewer, all to the same origin |
| Server render time of the index page | p99 under 5 ms on a laptop (Go benchmark) |

Techniques:

- Fonts: convert to WOFF2 and subset to Basic Latin, Latin-1 Supplement, Latin Extended-A and the punctuation used (keeps "Türkiye", Czech, Polish and similar names). Keep the `wght` axis. Pin the `wdth` axis of Mozilla Headline to 100 with `fonttools varLib.instancer` unless the design uses it. Preload both files with `crossorigin`. Use `font-display: swap` and a metric-matched local fallback (`size-adjust`, `ascent-override`, `descent-override`) so swapping does not shift the layout.
- Give every image explicit `width` and `height`. The hero has no image files; its visual is inline SVG.
- No render-blocking resources except the single stylesheet. The script uses `defer`.
- Templates are parsed once. Precompute everything that does not depend on the request (hashed asset names, JSON-LD, CSP header value, sitemap).
- Static files: hashed names, `Cache-Control: public, max-age=31536000, immutable`, `ETag`, gzip precompressed.
- Measure, do not guess: run Lighthouse against the local container in the verification phase (with `npx lighthouse` only if Node is already installed on the machine; never add Node to the repository), and report the numbers.

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
| `SITE_URL` | `http://localhost:8080` | Canonical base URL. Production value: `https://tiefer.space` |
| `CONTACT_EMAIL` | `hello@tiefer.space` | Public contact address and `mailto:` link |
| `LINKEDIN_URL` | `https://www.linkedin.com/company/tiefer/` | LinkedIn company page (final value) |
| `REPO_URL` | empty | Public source repository, footer link |
| `SMTP_HOST`, `SMTP_PORT`, `SMTP_USER`, `SMTP_PASS` | empty | Mail delivery for the contact form (TLS required) |
| `CONTACT_TO` | empty | Where form messages are delivered |
| `CONTACT_FROM` | empty | Sender address for form mails |
| `CSRF_SECRET` | random in dev | Signing key for form tokens, required in production |
| `LEGAL_*` | placeholders | Company legal name, legal form, address in Baku, VÖEN (tax ID), state registration, managing director, hosting provider and country, SMTP provider and country |
| `LEGAL_REVIEWED` | `false` | Hides the placeholder warning on legal pages when `true` |
| `LOG_RETENTION_DAYS` | `30` | Shown on the privacy page; must match the Azure log retention |
| `BEHIND_FRONT_DOOR` | `false` | `true` on Azure: enables the Front Door check and client IP rules in `docs/HARDENING.md` section 8 |
| `FRONT_DOOR_ID` | empty | The Front Door profile ID; required when `BEHIND_FRONT_DOOR=true` in production |
| `HSTS_PRELOAD` | `false` | Adds `preload` to HSTS; founder decision, see `docs/HARDENING.md` section 16 |

In production, `CSRF_SECRET` and `SMTP_PASS` come from Azure Key Vault references (see `docs/AZURE.md`). They never appear in the repository, in logs or in error messages.

---

## 12. Repository hygiene

- `go.mod` with module path from a `MODULE_PATH` you ask me for, or use `github.com/tiefer-labs/web`.
- `gofmt`, `go vet`, `staticcheck`, `gosec` and `govulncheck` clean.
- Tests: config validation, every page renders with status 200, 404 page, the text check, the performance budgets from section 9, and every security test listed in `docs/HARDENING.md` section 14. Fuzz targets as listed there.
- Workflows `ci.yml`, `codeql.yml` and `deploy.yml`, plus `dependabot.yml`, exactly as described in `docs/HARDENING.md` section 13 and `docs/AZURE.md` section 5.
- `.gitignore`, `.dockerignore`, `.editorconfig`, `.env.example`.
- `docs/SECURITY-DECISIONS.md`: every accepted risk, suppressed finding and trade-off, with the reason and when to revisit it.

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

1. `make run`, `make build`, `make check`, `make fuzz` and `make docker` all work with no errors, and `gosec`, `govulncheck` and `staticcheck` report zero findings.
2. `go build` produces one binary that serves the whole site with no files on disk next to it.
3. The index page, `/legal`, `/privacy`, `/acceptable-use` and the 404 page render correctly at all widths listed in section 9.
4. The contact form works with and without JavaScript against a local SMTP catcher (document how to test it, for example with Mailpit in Docker).
5. No em dash, en dash or emoji anywhere (the Go test passes).
6. Fonts are self-hosted, the licence files are present, and there are zero browser requests to third-party domains (verify with a headless browser).
7. All copy comes from `internal/content/en.go`; changing a sentence there changes the site.
8. `LICENSE` contains the full MPL 2.0 text, every source file has the MPL header, and `NOTICE.md` covers brand assets, fonts and dependencies.
9. `README.md` explains: what the project is, how to run it, how to edit copy, how to update the use cases and the roadmap, how to add another language, how to configure SMTP, how the Azure deployment works (short, linking to `infra/README.md`), the licence, and the list of placeholders still to fill.
9a. Every security test in `docs/HARDENING.md` section 14 exists and passes. Every header in section 4 of that document is present on every route.
9b. The performance budgets in section 9 pass as Go tests, and the Lighthouse results are reported.
9c. `infra/` builds with the Bicep CLI (or the report says the CLI was not available), and nothing was deployed.
9d. The design follows "Design restraint" in section 3.
10. Take full-page screenshots of every page at 390 px and 1440 px and review them against "Not a template" in section 3 for at least two rounds before telling me you are done. No section may look like it could belong to another company's site. Fix anything unbalanced, cramped, loud or generic.
11. At the end, give me a short list of anything you could not do (for example, fonts you could not download) and exactly what I need to provide.
