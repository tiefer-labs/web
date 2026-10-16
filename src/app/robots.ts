// SPDX-License-Identifier: MPL-2.0
import type { MetadataRoute } from "next";
import { env, siteUrl } from "@/lib/env.server";

/** Only production is indexed; previews and local servers ask every crawler to stay away. */
export default function robots(): MetadataRoute.Robots {
  if (env.ENV !== "production") return { rules: { userAgent: "*", disallow: "/" } };
  return { rules: { userAgent: "*", allow: "/" }, sitemap: `${siteUrl}/sitemap.xml` };
}
