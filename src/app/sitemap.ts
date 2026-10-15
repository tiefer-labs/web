// SPDX-License-Identifier: MPL-2.0
import type { MetadataRoute } from "next";
import { getPathname } from "@/i18n/navigation";
import { routing } from "@/i18n/routing";
import { siteUrl } from "@/lib/env.server";

const paths = ["/", "/legal", "/privacy", "/acceptable-use"];

export default function sitemap(): MetadataRoute.Sitemap {
  return routing.locales.flatMap((locale) =>
    paths.map((href) => ({ url: `${siteUrl}${getPathname({ href, locale })}` })),
  );
}
