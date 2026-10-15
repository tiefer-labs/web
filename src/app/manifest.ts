// SPDX-License-Identifier: MPL-2.0
import type { MetadataRoute } from "next";
import { getTranslations } from "next-intl/server";
import { routing } from "@/i18n/routing";

export default async function manifest(): Promise<MetadataRoute.Manifest> {
  const t = await getTranslations({ locale: routing.defaultLocale, namespace: "Meta" });
  return {
    id: "/",
    name: t("siteName"),
    short_name: t("siteName"),
    description: t("description"),
    start_url: "/",
    scope: "/",
    display: "browser",
    lang: routing.defaultLocale,
    theme_color: "#0C003D",
    background_color: "#FFFFFF",
    icons: [
      { src: "/icon-192.png", sizes: "192x192", type: "image/png" },
      { src: "/icon-512.png", sizes: "512x512", type: "image/png" },
      { src: "/icon.svg", sizes: "any", type: "image/svg+xml" },
    ],
  };
}
