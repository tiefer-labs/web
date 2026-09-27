# Notices

This repository contains the source code of the Tiefer website. Different
parts of it are under different terms. This file explains which is which,
in plain English. It is not legal advice.

## Code: Mozilla Public License 2.0

The Go code, templates, stylesheet, JavaScript and build files are licensed
under the Mozilla Public License, version 2.0 (SPDX identifier `MPL-2.0`).
The full text is in [`LICENSE`](LICENSE). Each source file carries the
standard MPL 2.0 header.

In short: you may use, change and redistribute the code. If you distribute a
changed version of an MPL-licensed file, you must make the source of that
file available under the MPL too.

## Website copy

The text content in `internal/content/` (the website copy and the legal page
drafts) is also covered by the MPL 2.0, unless the company decides otherwise.
This is an open decision for the founder; see the README.

## The Tiefer name and logo are not covered

The MPL 2.0 does not grant any rights to trademarks (section 2.3 of the
licence). The following are **not** licensed under the MPL and remain the
property of the company that operates Tiefer:

- the name "Tiefer" as a brand, and the slogans used on the site;
- the logos and marks in `web/static/brand/`;
- the favicons and app icons (`web/static/favicon.svg`, `favicon.ico`,
  `apple-touch-icon.png`, `icon-192.png`, `icon-512.png`);
- the social preview image `web/static/og-image.png`.

You may not use them to suggest that Tiefer endorses you or your product,
and you may not use them to brand a fork of this code. If you publish a
fork, replace the name, logos, icons and social image with your own.

## Fonts: SIL Open Font License 1.1

The site uses two typefaces by Mozilla, self-hosted in `web/static/fonts/`:

| Font | Files | Licence |
|---|---|---|
| Mozilla Headline | `MozillaHeadline-VF.woff2` | SIL Open Font License 1.1, see `OFL-MozillaHeadline.txt` |
| Mozilla Text | `MozillaText-VF.woff2` | SIL Open Font License 1.1, see `OFL-MozillaText.txt` |

The WOFF2 files were converted, unchanged in glyphs and features, from the
variable TTF files (version 1.000) distributed by the Mozilla Headline and
Mozilla Text projects (`github.com/mozilla/mozilla-headline-type`,
`github.com/mozilla/mozilla-text-type`). The fonts stay under the OFL; the
MPL does not apply to them.

## Third-party Go modules

None. The website uses only the Go standard library. The Go standard
library is compiled into the binary and is licensed under the BSD 3-Clause
licence (Copyright The Go Authors).

Development tools that are not part of the binary:

| Tool | Used for | Licence |
|---|---|---|
| staticcheck (`honnef.co/go/tools`) | `make lint` and CI | MIT |

## Unicode data

`internal/textcheck/extpict.go` contains the list of code points with the
Unicode property Extended_Pictographic, derived from `emoji-data.txt` of
Unicode 15.0.0. Unicode data files are provided under the Unicode License V3
(Copyright Unicode, Inc.), see https://www.unicode.org/license.txt.
