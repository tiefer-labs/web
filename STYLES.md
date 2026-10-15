# STYLES.md

The design system of the Tiefer website. Tokens are defined once in `src/app/globals.css` and used as Tailwind
utilities; never hardcode a colour, font or size that has a token here.

The look is calm and technical: white paper, one deep ink blue, thin lines, uppercase labels, generous space. Dark
sections in ink blue break the page into chapters. Motion is limited to one satellite gliding into place.

## Colour

| Token          | Value                     | Utility example                    | Use                                                   |
| -------------- | ------------------------- | ---------------------------------- | ----------------------------------------------------- |
| `ink`          | `#0c003d`                 | `bg-ink`, `text-ink`               | Brand blue. Headings, primary buttons, dark sections. |
| `ink-80`       | `#3d3364`                 | `text-ink-80`                      | Body text on paper (the `body` default).              |
| `ink-60`       | `#6d6690`                 | `text-ink-60`                      | Labels, hints, secondary text, timeline lines.        |
| `line`         | `#e4e1ee`                 | `border-line`                      | Hairlines and dividers on paper.                      |
| `tint`         | `#f5f4f9`                 | `bg-tint`                          | Alternate section background, highlighted rows.       |
| `paper`        | `#fff`                    | `bg-paper`                         | Page background, cards, inputs.                       |
| `on-ink`       | `#fff`                    | `text-on-ink`                      | Headings and buttons on ink.                          |
| `on-ink-muted` | `#b9b4d1`                 | `text-on-ink-muted`                | Body text and labels on ink.                          |
| `line-on-ink`  | `rgb(185 180 209 / 0.24)` | `border-line-on-ink`               | Hairlines on ink.                                     |
| `inferred`     | `#9a6400`                 | `text-inferred`, `border-inferred` | "Inferred" evidence tag and the legal review note.    |

The theme colour (browser UI, manifest) is `#0C003D`. There is no dark mode: the site declares `color-scheme: light`.

## Typography

Two typefaces, both by Mozilla and loaded with `next/font/google` (served from our own origin):

- **Mozilla Headline** (`font-headline`), weights 500 and 600: `h1`, `h2`, `h3`, statements.
- **Mozilla Text** (`font-text`), weights 400 and 500: everything else.

Each has a metric-matched local fallback (`Headline Fallback`, `Text Fallback`) so text does not move when the fonts
arrive.

| Token            | Size                                           | Use                                             |
| ---------------- | ---------------------------------------------- | ----------------------------------------------- |
| `text-h1`        | `clamp(2.75rem, 1.82rem + 3.81vw, 4.5rem)`     | Hero title, weight 600, leading 1.02            |
| `text-h2`        | `clamp(1.875rem, 1.46rem + 1.71vw, 3rem)`      | Section titles (the `h2` default), leading 1.06 |
| `text-name`      | `clamp(2.5rem, 1.6rem + 1.9vw, 3.5rem)`        | The name section line on large screens          |
| `text-statement` | `clamp(1.5rem, 1.27rem + 0.95vw, 2.125rem)`    | Large statement paragraphs (headline font)      |
| `text-h3`        | `clamp(1.25rem, 1.2rem + 0.19vw, 1.375rem)`    | Item titles (the `h3` default), leading 1.2     |
| `text-lead`      | `clamp(1.125rem, 1.06rem + 0.29vw, 1.3125rem)` | Intro paragraphs, leading 1.5                   |
| `text-body`      | `clamp(1.0625rem, 1.04rem + 0.1vw, 1.125rem)`  | Body text (the `body` default), leading 1.55    |
| `text-small`     | `0.9375rem`                                    | Notes, nav, footer, table cells                 |
| `text-label`     | `0.75rem`                                      | Labels                                          |
| `text-tag`       | `0.6875rem`                                    | Evidence tags                                   |

Headings use `text-wrap: balance`, paragraphs `pretty`. Headings are weight 500 unless stated.

`label-caps` is the shared label style: Mozilla Text, `text-label`, weight 500, leading 1.3, `0.12em` tracking,
uppercase, tabular figures. Use it for section labels ("02 / PRODUCT"), eyebrows, table headers and the hero
coordinates.

## Space and layout

| Variable    | Below 480 px | 480 px and up | 1024 px and up | Use                                |
| ----------- | ------------ | ------------- | -------------- | ---------------------------------- |
| `--pad`     | 16px         | 24px          | 40px           | Page side padding (`px-(--pad)`)   |
| `--gap`     | 24px         | 24px          | 32px           | Grid column gap (`gap-x-(--gap)`)  |
| `--section` | 72px         | 72px          | 128px          | Section padding (`py-(--section)`) |
| `--max`     | 1120px       |               |                | Content width                      |
| `--stroke`  | 1.25px       |               |                | Line art stroke width              |

- Breakpoints: `xs` 480px (custom), `sm` 640px, `lg` 1024px. The design only changes at these three; `md` is not used.
- `wrap`: centred container, `--max` wide plus padding.
- `grid-12`: one column on small screens, twelve from `lg`. Every child spans the full row until it is placed with
  `lg:col-[a/b]` and `lg:row-n`.

### Section anatomy

Every section uses `Section` and `SectionLabel` from `src/components/ui/section.tsx`:

```tsx
<Section id="principles" labelledBy="principles-title" tone="tint">
  <SectionLabel className="lg:row-1">{t("label")}</SectionLabel>
  <h2 id="principles-title" className="lg:col-[1/6] lg:row-2">
    {t("title")}
  </h2>
  ...
</Section>
```

- `tone`: `paper` (default), `tint` or `ink`. `ink` sets `data-tone="dark"`, which switches label and focus colours.
- Two layouts repeat: label in columns 1 to 2 with the title from column 3 on the same row, or label on its own row
  with the content split into columns 1 to 5 and 7 to 12.
- Section order alternates paper, ink and tint so neighbours never share a background.

## Components and patterns

- **Buttons** (`ButtonLink`, `buttonClass` in `src/components/ui/button.tsx`): 2px radius, 1px border, weight 500.
  `primary` is filled ink, `outline` is a bordered ghost; `tone="dark"` inverts both for ink backgrounds;
  `size="sm"` is the header call to action. Text links are plain underlined links with weight 500.
- **Links**: underlined, 1px thick, 0.2em offset; 2px on hover.
- **Lists and timelines**: hairline-separated rows (`border-t border-line`). Numbered steps use `counter-steps`
  ("01", "02"). Timeline nodes are `SquareMark` squares on a 1px rail; a filled square marks the last step or the
  current milestone.
- **Evidence tags**: uppercase `text-tag` in a 1px box, `text-ink` for observed and `text-inferred` for inferred, on a
  `bg-tint` row. Observed and inferred are always labelled separately.
- **Cards**: only the alert packet: paper, 1px `line` border, a very soft ink shadow.
- **Form fields** (`Field` in `src/components/home/contact/field.tsx`): label above in `text-small` weight 500,
  1px `ink-60` border, 2px radius, ink border on hover; invalid fields get a 2px ink border and a message whose
  "Error:" prefix is only read by screen readers. Status messages sit in a box with a 2px ink left border.
- **Line art**: inline SVG, `vector-effect: non-scaling-stroke`, `--stroke` width, `on-ink-muted` on dark and `ink`
  at low opacity on light. Always `aria-hidden`.

## Accessibility and motion

- Focus ring: 2px solid ink, 3px offset; white inside `data-tone="dark"`.
- A skip link leads to `main#main`. Every section is labelled by its heading; each page has one `h1` and no skipped
  heading levels (tested).
- Motion: only `animate-satellite-pass`, behind `motion-safe:`. `prefers-reduced-motion: reduce` turns off all
  transitions, animations and smooth scrolling.
- `forced-colors`: dark sections follow the system colours; logos keep theirs (`data-forced-colors="keep"`).
- Print: `data-print="hide"` removes the header, line art and the contact panel; dark sections print black on white.

## Brand

- The name and logos are trademarks, not covered by the code licence (see [NOTICE.md](NOTICE.md)). Never edit,
  recolour or combine them with other marks.
- Files are in `public/brand/`: `tiefer-logo.svg` (ink) on light backgrounds, `tiefer-logo-white.svg` on ink;
  the four-line mark `tiefer-mark-*.svg` for decoration. The header logo is 126 by 24 px, the footer logo 84 by 16 px.
- The four lines of the mark stand for the four stages: Filter, Detect, Alert, Update.
