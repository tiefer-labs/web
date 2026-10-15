// SPDX-License-Identifier: MPL-2.0
import { defineRouting } from "next-intl/routing";

export const locales = ["en"] as const;

export const routing = defineRouting({
  locales,
  defaultLocale: "en",
  // The default locale keeps the unprefixed URLs (/legal); other locales get /az/legal.
  localePrefix: "as-needed",
  // The site sets no cookies, and the privacy notice says so.
  localeCookie: false,
});
