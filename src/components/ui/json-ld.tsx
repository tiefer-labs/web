// SPDX-License-Identifier: MPL-2.0

/** Escapes the JSON so it cannot close the script element; also escapes the line and paragraph separators. */
function toSafeJsonLd(data: object): string {
  return JSON.stringify(data).replace(
    /[<>&\p{Zl}\p{Zp}]/gu,
    (c) => `\\u${c.charCodeAt(0).toString(16).padStart(4, "0")}`,
  );
}

export function JsonLd({ data }: { data: object }) {
  // oxlint-disable-next-line react/no-danger
  return <script type="application/ld+json" dangerouslySetInnerHTML={{ __html: toSafeJsonLd(data) }} />;
}
