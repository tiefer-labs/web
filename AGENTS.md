<!-- BEGIN:nextjs-agent-rules -->

## This is NOT the Next.js you know

This version has breaking changes — APIs, conventions, and file structure may all differ from your training data. Read the relevant guide in `node_modules/next/dist/docs/` (resolved from this file's directory; in monorepos the `next` package may not be visible from the repo root) before writing any code. Heed deprecation notices.

This block is written and re-added by `next dev` — verify at `node_modules/next/dist/server/lib/generate-agent-files.js`. Removing it from a diff only re-creates the uncommitted change; committing it with your work keeps the tree clean.

<!-- END:nextjs-agent-rules -->

# AGENTS.md

The single source of truth for anyone, human or AI agent, working on this repository.

## Project

The public website of **Tiefer** (https://tiefer.space), a space technology startup from Baku that builds AI software
for Earth observation satellites: filter useless pixels in orbit, detect events on board, send small alerts instead of
raw imagery. The research behind it lives in [tiefer-labs/lab](https://github.com/tiefer-labs/lab) (motto: "Measured,
not claimed").

The site is one long home page, three legal pages and a 404 page, plus a contact form. It used to be a Go server;
that version is in the git history before the Next.js migration.

## Commands

- `pnpm dev`: development server on http://localhost:3000
- `pnpm build` / `pnpm start`: production build and server
- `pnpm fmt`, then `pnpm check`: always run both before you finish. `check` runs greenly: TypeScript, Oxfmt, Oxlint
  and Vitest. CI runs `pnpm check --strict`, which also fails on lint warnings.
- `pnpm test` / `pnpm test:watch`: Vitest only
- `pnpm install:skills`: install the agent skills listed in `skills-lock.json`

pnpm only, never npm or npx. Run local binaries as `pnpm <bin>`.

## Stack

- Next.js 16 App Router with React 19, React Compiler, Cache Components, typed routes, `output: "standalone"`
- TypeScript, Tailwind CSS v4 (design tokens in `src/app/globals.css`, documented in [STYLES.md](STYLES.md))
- next-intl for every visible string, `next/font/google` for the Mozilla Headline and Mozilla Text fonts
- zod for configuration and form validation, nodemailer for SMTP, optional Redis (Upstash or any Redis URL)
- Oxfmt, Oxlint, Vitest (+ Testing Library), greenly

## Structure

- `messages/en.json`: all copy, by section (`Hero`, `Problem`, ..., `Contact`, `Legal`, `NotFound`)
- `src/app/[locale]/`: the localised routes. `layout.tsx` is the root layout (shared Open Graph fields),
  `opengraph-image.png` the social image for every page (Next reuses it for Twitter), `page.tsx` the home page,
  `(legal)/` the three legal pages, `not-found.tsx` the 404 page
- `src/app/global-not-found.tsx`: the 404 for URLs that match no route
- `src/app/`: `robots.ts`, `sitemap.ts`, `manifest.ts`, icons, `healthz/` (Azure and Docker probe),
  `.well-known/security.txt/`
- `src/components/ui/`: shared building blocks (`Section`, `SectionLabel`, `ButtonLink`/`buttonClass`, `SquareMark`,
  `JsonLd`)
- `src/components/layout/`: header, mobile nav, footer, skip link
- `src/components/home/`: one component per home page section; `contact/` holds the form
- `src/components/legal/`: the shared legal page and placeholder highlighting
- `src/lib/contact/`: the contact form backend (token, validation, rate limits, mail, the Server Action)
- `src/lib/env.server.ts`: every setting, validated with zod; `.env.example` documents each one
- `src/i18n/`: next-intl routing, request config and navigation helpers
- `src/proxy.ts`: the Front Door check, then next-intl routing
- `tests/`: mirrors `src/`; shared helpers in `tests/test-utils/` (import via `@test/*`)
- `public/brand/`: logos and the mark (never edit them); `public/icon-*.png` for the manifest
- `infra/`: Azure Bicep (App Service behind Front Door), `Dockerfile` for that deployment

## Conventions

- **Server Components first.** Only the mobile nav, the footer links and the contact form are Client Components.
- **No hardcoded copy.** Every visible string is a message in `messages/en.json`. Lists are keyed objects (next-intl has
  no arrays) iterated with `keysOf()` from `src/lib/messages.ts`; the order of keys is the order on the page.
- **Tailwind tokens only.** Use the colours, type scale and spacing from [STYLES.md](STYLES.md). No inline styles, no new
  CSS files; `globals.css` only holds tokens, base styles and the few utilities Tailwind cannot express.
- **Compose classes with `cn` from `@/lib/cn`** (the `cn` package, taught our type scale tokens), never with template
  strings.
- **Reuse the building blocks.** A new home page section is a `Section` with a `SectionLabel` and a heading whose id
  matches `labelledBy`.
- **Every source file starts with** `// SPDX-License-Identifier: MPL-2.0` (`/* ... */` in CSS). A test checks this.
- **Comments are short** (one line, two at most) and explain why, never what.
- **Icons:** the design uses its own inline SVGs. If an icon library is ever needed, use `react-icons`.
- Path alias `@/*` maps to `src/*`. Oxfmt enforces 120 columns, Tailwind class order and import order.

## Copy rules (STRICT)

Tiefer's writing rules come from tiefer-labs/lab `docs/STYLE.md` and `CLAIMS.md`. Tests fail the build on the
mechanical ones, in messages, rendered pages and every text file in the repository:

- No em dashes, en dashes or horizontal bars. No emoji, arrows or decorative symbols.
- Banned words: revolutionary, cutting-edge, game-changing, AI-powered, seamless, unlock, leverage, empower, magic.
- British spelling, short sentences, facts over adjectives. Write "percent", dates as "1 October 2026".
- Anything planned or not yet built is written as an aim, never in the present tense. A public claim needs a row in
  lab's `CLAIMS.md`; numbers only under its `POLICY.md`.
- `[UPPER_CASE]` text is a placeholder still waiting for real details; legal pages highlight it.

## Configuration

Settings come from environment variables only (`.env.local` locally, never committed). Pages are prerendered at build
time, so the values they show (`SITE_URL`, `CONTACT_EMAIL`, `LEGAL_*`, whether `SMTP_HOST` is set) must be present
when `pnpm build` runs. Secrets (`CSRF_SECRET`, `SMTP_PASS`, Redis tokens) are only read by the running server.

## Contact form

`src/lib/contact/handle.ts` runs these steps, in order: attempts limit (20 per hour per
client), honeypot, signed token (at least 3 seconds and at most 2 hours old), validation, single-use token, send limits
(5 per hour per client, 60 per hour in total), then SMTP. Bots get a fake success. The per-request token is the only
dynamic part of the home page: a `<Suspense>` hole that calls `connection()`.

Rate limits and used tokens live in memory unless `UPSTASH_REDIS_REST_URL` or `REDIS_URL` is set. Memory is only correct
with one server instance; use Redis on Vercel.

## Tests

Vitest: `*.test.ts` run in Node, `*.test.tsx` in jsdom. Render Server Components with `renderServer()` and client
components with `renderClient()` from `tests/test-utils/render.tsx`; `next-intl/server` is mocked to read
`messages/en.json`. Port behaviour, not markup: assert what a visitor or screen reader gets.

## Skills

Read these before writing or editing components:

- `.agents/skills/next-best-practices/SKILL.md`
- `.agents/skills/next-cache-components/SKILL.md`
- `.agents/skills/vercel-composition-patterns/`
- `.agents/skills/vercel-react-best-practices/`
- `.agents/skills/web-design-guidelines/SKILL.md`
