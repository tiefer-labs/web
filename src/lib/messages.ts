// SPDX-License-Identifier: MPL-2.0

/** Message lists are keyed objects (next-intl has no arrays); this keeps their order and key types. */
export function keysOf<T extends object>(group: T): (keyof T & string)[] {
  // oxlint-disable-next-line typescript/no-unsafe-type-assertion
  return Object.keys(group) as (keyof T & string)[];
}
