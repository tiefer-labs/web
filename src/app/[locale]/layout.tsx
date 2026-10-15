// SPDX-License-Identifier: MPL-2.0
import type { Metadata, Viewport } from "next";
import { NextIntlClientProvider } from "next-intl";
import { getLocale, getMessages, getTranslations } from "next-intl/server";
import { SiteFooter } from "@/components/layout/site-footer";
import { SkipLink } from "@/components/layout/skip-link";
import { locales } from "@/i18n/routing";
import { cn } from "@/lib/cn";
import { siteUrl } from "@/lib/env.server";
import "../globals.css";
import { headline, text } from "@/lib/fonts";

// Pages set only title, description and canonical; Next derives the og:/twitter: title, description and image.
export async function generateMetadata(): Promise<Metadata> {
  const [locale, meta] = await Promise.all([getLocale(), getTranslations("Meta")]);
  return {
    metadataBase: new URL(siteUrl),
    openGraph: { type: "website", siteName: meta("siteName"), locale },
    twitter: { card: "summary_large_image" },
  };
}

export const viewport: Viewport = {
  themeColor: "#0C003D",
  colorScheme: "light",
};

export function generateStaticParams() {
  return locales.map((locale) => ({ locale }));
}

export default async function LocaleLayout({ children }: LayoutProps<"/[locale]">) {
  const [locale, messages] = await Promise.all([getLocale(), getMessages()]);

  return (
    <html lang={locale} data-scroll-behavior="smooth" className={cn(headline.variable, text.variable)}>
      <body>
        <NextIntlClientProvider messages={{ Contact: messages.Contact }}>
          <SkipLink />
          {children}
          <SiteFooter />
        </NextIntlClientProvider>
      </body>
    </html>
  );
}
