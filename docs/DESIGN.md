# Design

The page reads like a mission document: a strict grid, hairlines, type
that does the work, and one colour (`--ink`, `#0C003D`) with its tints.

## Grid

- Content width 1,120 px, centred. Side padding 16 px up to 480 px wide,
  24 px up to 1,024 px, 40 px above.
- 12 columns with 32 px gaps from 1,024 px; 6 columns with 24 px gaps from
  640 px; one column below.
- Columns 1 and 2 on desktop are the **label column**: every section starts
  with its document label (`02 / PRODUCT`) there, aligned to the first line
  of the heading. Content starts at column 3, 7 or 1 depending on the
  section, never at the same column in two neighbouring sections.

## Type scale

Mozilla Headline (500 and 600) for headings, Mozilla Text (400 and 500) for
everything else. Sizes are fluid between 390 px and 1,440 px.

| Role | Font, weight | Size (390 to 1,440 px) | Line height | Tracking |
|---|---|---|---|---|
| H1 | Headline 600 | 44 to 72 px | 1.02 | -0.025em |
| H2 | Headline 500 | 30 to 48 px | 1.06 | -0.018em |
| Statement (problem intro, name line) | Headline 500 | 24 to 34 px | 1.25 | -0.01em |
| H3 | Headline 500 | 20 to 22 px | 1.2 | -0.005em |
| Lead | Text 400 | 18 to 21 px | 1.5 | 0 |
| Body | Text 400 | 17 to 18 px | 1.55 | 0 |
| Small | Text 400 | 15 px | 1.5 | 0 |
| Technical label | Text 500, uppercase, tabular figures | 12 px | 1.3 | 0.12em |

Headings use `text-wrap: balance`, paragraphs `text-wrap: pretty`; text
columns stay under 72 characters.

## Spacing

An 8 px scale: 8, 16, 24, 32, 48, 64, 96, 128. Section padding is 128 px on
desktop and 72 px on mobile, except the hero (taller) and the name section
(shorter), so the rhythm is not uniform. Hairlines are 1 px in `--line`
(on dark bands: `--on-ink-muted` at 24% opacity).

## Sections

| Section | Background | Composition | Why it differs from its neighbours |
|---|---|---|---|
| Header | on the hero | Logo left, five links and one outlined button right; a disclosure menu on mobile. Not sticky. | The only horizontal bar; the button is outlined so only one solid button is in view. |
| Hero | `--ink` | H1 and lead in columns 1 to 6; the orbital drawing fills the right and bottom edge; position label top right. | The only full-bleed drawing and the only H1. |
| 01 Problem | white | H2 in columns 3 to 10; the statement large across columns 3 to 12; three numbered text columns under one hairline. | Horizontal reading: one statement, then three narrow columns side by side. |
| 02 Product | `--ink` | H2 and body in columns 1 to 5; a four-row specification table in columns 7 to 12, hairline rows. | Rows in a right-hand table, after the columns above; the second dark band. |
| 03 How it works | white | H2 in columns 3 to 9; one technical drawing across all 12 columns (rail with five square nodes and a branch at Filter); note; two ground columns split by a vertical hairline. | The only full-width diagram. |
| 04 Alert packet | `--tint` | Text in columns 1 to 5; the packet card in columns 7 to 12, table rows with data labels and two size bars. | The only object with a border: a data record, not text. |
| 05 Use cases | white | H2 in columns 3 to 8; four blocks in a 2 by 2 grid whose hairlines form a cross. | Square grid, where the neighbours are a card and a list. |
| 06 Principles | `--tint` | H2 in columns 1 to 5; six clauses in two columns, `P1` to `P6` in a narrow number column, hairline above each clause. | A numbered specification list, denser than anything around it. |
| 07 Roadmap | white | H2 in columns 3 to 9; a horizontal track on one hairline with five square nodes, the first filled and marked `Now`. | The only timeline; reads left to right along a single line. |
| 08 Name | `--ink` | One large line in columns 3 to 10 and a short body; the white mark large and quiet on the right; an Earth-limb curve along the bottom edge. | Shortest section, the third and last dark band. |
| 09 Contact | `--tint` with four waves at 6% | H2 and body in columns 1 to 5; the form in columns 7 to 12; email and LinkedIn under the form. | The only form; the waves are the only texture on a light band. |
| Footer | white | Small logo, links in one row, copyright and tagline. | Plain, small type. |

## Visual system

- One line weight (1.25 px, non-scaling) for every drawing: orbits, swath,
  flow rail, roadmap track, size bars.
- Square nodes (10 px) everywhere: flow, roadmap, satellite glyph parts.
- Technical labels share one style in drawings, the packet card and
  section labels.
- Motion: the satellite glyph moves a short way along its orbit once, over
  2.4 s, then rests. Nothing else moves. With `prefers-reduced-motion`
  nothing moves at all.
