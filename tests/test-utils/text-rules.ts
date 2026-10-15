// SPDX-License-Identifier: MPL-2.0

/**
 * The copy rules from tiefer-labs/lab docs/STYLE.md, ported from the Go `textcheck` package: no long dashes, no
 * emoji or decorative symbols, and none of the hype words.
 */

// © ® ™ are Extended_Pictographic but are ordinary text symbols on their own.
const TEXT_SYMBOLS = new Set([0x00a9, 0x00ae, 0x2122]);

const DECORATIVE: [number, number][] = [
  [0x2190, 0x21ff], // Arrows
  [0x2500, 0x27bf], // Box Drawing to Dingbats
  [0x27f0, 0x27ff], // Supplemental Arrows-A
  [0x2900, 0x297f], // Supplemental Arrows-B
  [0x2b00, 0x2bff], // Miscellaneous Symbols and Arrows
];

const PICTOGRAPHIC = /\p{Extended_Pictographic}/u;

function describe(char: string): string | null {
  const code = char.codePointAt(0)!;
  const hex = `U+${code.toString(16).toUpperCase().padStart(4, "0")}`;
  if (code === 0x2014) return "em dash (U+2014)";
  if (code === 0x2013) return "en dash (U+2013)";
  if (code === 0x2015) return "horizontal bar (U+2015)";
  if (code === 0xfe0f) return "emoji variation selector (U+FE0F)";
  if (code === 0x20e3) return "combining keycap (U+20E3)";
  if (code >= 0x1f1e6 && code <= 0x1f1ff) return `regional indicator (${hex})`;
  if (code >= 0x1f3fb && code <= 0x1f3ff) return `skin tone modifier (${hex})`;
  if (TEXT_SYMBOLS.has(code)) return null;
  if (PICTOGRAPHIC.test(char)) return `emoji (${hex})`;
  if (DECORATIVE.some(([from, to]) => code >= from && code <= to)) return `decorative symbol (${hex})`;
  return null;
}

/** Every forbidden character in `text`, as "line:col: what". */
export function forbiddenCharacters(text: string): string[] {
  const findings: string[] = [];
  let line = 1;
  let col = 1;
  for (const char of text) {
    if (char.codePointAt(0) === 0xfffd) findings.push(`${line}:${col}: invalid UTF-8`);
    const what = describe(char);
    if (what) findings.push(`${line}:${col}: ${what}`);
    if (char === "\n") {
      line++;
      col = 1;
    } else {
      col++;
    }
  }
  return findings;
}

const BANNED =
  /\b(revolutionar(?:y|ily)|cutting[- ]edge|game[- ]chang(?:ing|er|ers)|ai[- ]powered|seamless(?:ly)?|unlock(?:s|ed|ing)?|leverag(?:e|es|ed|ing)|empower(?:s|ed|ing|ment)?|magic(?:al|ally)?)\b/gi;

export function bannedWords(text: string): string[] {
  return [...text.matchAll(BANNED)].map((match) => match[0]);
}

/** Every string leaf of a messages object, with its dotted key. */
export function messageStrings(messages: object, prefix = ""): [string, string][] {
  return Object.entries(messages).flatMap(([key, value]): [string, string][] =>
    typeof value === "string" ? [[`${prefix}${key}`, value]] : messageStrings(value as object, `${prefix}${key}.`),
  );
}

/** The text people read or hear: text nodes plus alt, title, aria-label and similar attributes. */
export function visibleText(root: Element): string {
  const clone = root.cloneNode(true) as Element;
  for (const hidden of clone.querySelectorAll("script, style, template")) hidden.remove();
  const attributes = [...clone.querySelectorAll("[alt], [title], [aria-label], [aria-description], [placeholder]")]
    .flatMap((el) => ["alt", "title", "aria-label", "aria-description", "placeholder"].map((a) => el.getAttribute(a)))
    .filter(Boolean);
  return [clone.textContent ?? "", ...attributes].join("\n");
}
