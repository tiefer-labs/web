// SPDX-License-Identifier: MPL-2.0
export function cn(...classes: (string | false | null | undefined)[]) {
  return classes.filter(Boolean).join(" ");
}
