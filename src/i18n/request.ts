// SPDX-License-Identifier: MPL-2.0
import { hasLocale } from "next-intl";
import { getRequestConfig } from "next-intl/server";
import { notFound } from "next/navigation";
import * as rootParams from "next/root-params";
import { routing } from "./routing";

export default getRequestConfig(async ({ locale }) => {
  if (!locale) {
    const paramValue = await rootParams.locale();
    // No locale segment at all (the global 404 page) falls back to the default locale.
    if (paramValue === undefined) locale = routing.defaultLocale;
    else if (hasLocale(routing.locales, paramValue)) locale = paramValue;
    else notFound();
  }

  return {
    locale,
    messages: (await import(`../../messages/${locale}.json`)).default,
  };
});
