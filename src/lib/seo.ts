// SPDX-License-Identifier: MPL-2.0
import "server-only";
import type { Metadata } from "next";
import type { Locale } from "next-intl";
import { getTranslations } from "next-intl/server";
import { getPathname } from "@/i18n/navigation";
import { routing } from "@/i18n/routing";

type BuildMetadataArgs = {
  locale: Locale;
  path: string;
  /** Page title without the site name; omit on the home page. */
  title?: string;
  description?: string;
};

/** Title, description, canonical and hreflang for one page; Next derives the social tags from them. */
export async function buildMetadata({ locale, path, title, description }: BuildMetadataArgs): Promise<Metadata> {
  const meta = await getTranslations({ locale, namespace: "Meta" });
  const fullTitle = title ? `${title} | ${meta("siteName")}` : meta("title");
  const summary = description ?? meta("description");
  const canonical = getPathname({ href: path, locale });

  const languages: Record<string, string> = {};
  if (routing.locales.length > 1) {
    for (const l of routing.locales) languages[l] = getPathname({ href: path, locale: l });
    languages["x-default"] = getPathname({ href: path, locale: routing.defaultLocale });
  }

  return {
    title: { absolute: fullTitle },
    description: summary,
    alternates: { canonical, languages },
  };
}
