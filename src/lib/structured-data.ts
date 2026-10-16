// SPDX-License-Identifier: MPL-2.0
import "server-only";
import { env, siteUrl } from "@/lib/env.server";

export function organizationSchema(name: string) {
  return {
    "@context": "https://schema.org",
    "@type": "Organization",
    logo: `${siteUrl}/icon-512.png`,
    name,
    sameAs: [env.LINKEDIN_URL],
    url: `${siteUrl}/`,
  };
}
