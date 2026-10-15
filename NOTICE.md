# Notices

## Code: Mozilla Public License 2.0

The source code of this website is licensed under the Mozilla Public
License, v. 2.0 (`MPL-2.0`); the full text is in [`LICENSE`](LICENSE). Every
source file carries an `SPDX-License-Identifier: MPL-2.0` header, and a test
checks this.

## Website copy

The text content in `messages/` is also covered by MPL 2.0, unless the
founder decides otherwise. This is an open decision.

## The Tiefer name and logo

The MPL does not grant trademark rights (section 2.3 of the licence). The
Tiefer name, the logos and the four-wave mark in `public/brand/`, the
icons (`src/app/icon.svg`, `src/app/favicon.ico`, `src/app/apple-icon.png`,
`public/icon-192.png`, `public/icon-512.png`) and the social image
(`src/app/[locale]/opengraph-image.png`) are not covered
by the MPL. They are trademarks and brand assets of the company. You may not
use them to suggest endorsement by Tiefer, or to brand a fork or a modified
version of this website.

## Fonts

Mozilla Headline and Mozilla Text are licensed under the SIL Open Font
License 1.1 (Copyright 2025 The Mozilla Headline Project Authors and The
Mozilla Text Project Authors). They are downloaded from Google Fonts by
`next/font` when the site is built and served from the site's own origin.

## Dependencies

The npm packages the site uses are listed with their versions in
[`package.json`](package.json) and [`pnpm-lock.yaml`](pnpm-lock.yaml); each
is distributed under its own licence. The main runtime dependencies are
Next.js, React and React DOM, next-intl, zod and nodemailer (all MIT), and
the Redis clients `@upstash/redis` (MIT) and `redis` (MIT).

## Build and deployment components

Downloaded when the site is built or deployed; not part of this
repository:

| Component                              | Used for          | Licence                                                          |
| -------------------------------------- | ----------------- | ---------------------------------------------------------------- |
| `node` Docker image                    | build stage only  | MIT (Node.js); Debian packages under their own licences          |
| `gcr.io/distroless/nodejs24-debian12`  | runtime image     | Apache 2.0; Node.js and Debian packages under their own licences |
| GitHub Actions in `.github/workflows/` | CI and deployment | each under its own licence, pinned by commit                     |
