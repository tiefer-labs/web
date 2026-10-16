// SPDX-License-Identifier: MPL-2.0
import "server-only";
import { env, siteUrl } from "@/lib/env.server";

/** RFC 9116 asks for an expiry under a year away; a test fails a month before this date. */
export const SECURITY_TXT_EXPIRES = "2027-09-01T00:00:00Z";

export function securityTxt() {
  return [
    `Contact: mailto:${env.CONTACT_EMAIL}`,
    `Expires: ${SECURITY_TXT_EXPIRES}`,
    "Preferred-Languages: en",
    `Canonical: ${siteUrl}/.well-known/security.txt`,
    ...(env.REPO_URL ? [`Policy: ${env.REPO_URL}/blob/main/SECURITY.md`] : []),
    "",
  ].join("\n");
}
