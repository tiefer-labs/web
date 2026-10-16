// SPDX-License-Identifier: MPL-2.0
const FORBIDDEN = /[\r\n<>"(),;: \t]/;

/** A single plain address: no display name, no list, and a dot in the domain. */
export function isEmail(value: string): boolean {
  if (value === "" || value.length > 254 || FORBIDDEN.test(value)) return false;
  const at = value.lastIndexOf("@");
  const local = value.slice(0, at);
  const domain = value.slice(at + 1);
  if (at <= 0 || local.includes("@") || !domain.includes(".")) return false;
  const atoms = [...local.split("."), ...domain.split(".")];
  return atoms.every((atom) => atom !== "" && !/[[\]]/.test(atom));
}
