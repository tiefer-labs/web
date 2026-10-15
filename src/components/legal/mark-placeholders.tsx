// SPDX-License-Identifier: MPL-2.0
import { Fragment, isValidElement, type ReactNode } from "react";

/** `[UPPER_CASE]` or `[UPPER_CASE: note]`: text still waiting for real details. */
export const PLACEHOLDER = /\[[A-Z][A-Z0-9_]{2,}(?::[^\]]*)?\]/g;

function markText(text: string, key: string): ReactNode[] {
  const parts: ReactNode[] = [];
  let last = 0;
  for (const match of text.matchAll(PLACEHOLDER)) {
    if (match.index > last) parts.push(text.slice(last, match.index));
    parts.push(
      <mark key={`${key}-${match.index}`} className="bg-[rgb(154_100_0/0.12)] px-[3px] text-ink">
        {match[0]}
      </mark>,
    );
    last = match.index + match[0].length;
  }
  if (last < text.length) parts.push(text.slice(last));
  return parts;
}

/** Highlights placeholders in the text parts of a rich translation, both in the copy and in config values. */
export function markPlaceholders(node: ReactNode): ReactNode {
  const nodes = Array.isArray(node) ? node : [node];
  return nodes.map((part, index) =>
    typeof part === "string" ? (
      <Fragment key={index}>{markText(part, String(index))}</Fragment>
    ) : isValidElement(part) ? (
      <Fragment key={index}>{part}</Fragment>
    ) : (
      part
    ),
  );
}
