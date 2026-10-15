// SPDX-License-Identifier: MPL-2.0
export const roles = ["operator", "agency", "manufacturer", "hosted-compute", "investor", "other"] as const;

export type Role = (typeof roles)[number];

export function isRole(value: string): value is Role {
  return (roles as readonly string[]).includes(value);
}
