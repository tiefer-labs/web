# Contributing to the Tiefer website

The public website of Tiefer (Baku), an AI analyst for satellite data. Go
standard library only, server-rendered, one self-contained binary. The
original brief is `docs/PROMPT.md`; the README explains the structure.

Before changing anything, keep these rules (tests enforce most of them):

- All copy lives in `internal/content/en.go`. Templates only read from it.
- No em dash (U+2014), en dash (U+2013), U+2015 or emoji anywhere. No hype
  words (revolutionary, cutting-edge, game-changing, AI-powered, seamless,
  unlock, leverage, empower, magic). No invented customers, partners,
  numbers or testimonials. Never imply Tiefer operates its own satellites.
- Colours only from the tokens in `web/static/css/site.css`; primary `#0C003D`.
  Fonts Mozilla Headline and Mozilla Text, self-hosted.
- No cookies, analytics, trackers, third-party scripts or requests. No inline
  scripts or styles (the CSP forbids them).
- Do not modify files in `web/static/brand/`, the icons or `og-image.png`.
- Every new Go, template, CSS and JS file starts with the MPL 2.0 header.
- New Go modules need a reason in the README and an entry in `NOTICE.md`.

Before committing: `make lint test`.
