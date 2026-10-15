// SPDX-License-Identifier: MPL-2.0
import type { Metadata, Viewport } from "next";
import { NextIntlClientProvider } from "next-intl";
import { SiteFooter } from "@/components/layout/site-footer";
import { SiteHeader } from "@/components/layout/site-header";
import { SkipLink } from "@/components/layout/skip-link";
import { NotFoundContent, notFoundMetadata } from "@/components/not-found-content";
import { routing } from "@/i18n/routing";
import { cn } from "@/lib/cn";
import "./globals.css";
import { headline, text } from "@/lib/fonts";

export function generateMetadata(): Promise<Metadata> {
  return notFoundMetadata();
}

export const viewport: Viewport = {
  themeColor: "#0C003D",
  colorScheme: "light",
};

/** Paths that match no route at all skip the locale layout, so this page brings its own document. */
export default function GlobalNotFound() {
  return (
    <html lang={routing.defaultLocale} className={cn(headline.variable, text.variable)}>
      <body>
        <NextIntlClientProvider locale={routing.defaultLocale} messages={{}}>
          <SkipLink />
          <SiteHeader />
          <main id="main" tabIndex={-1}>
            <NotFoundContent />
          </main>
          <SiteFooter />
        </NextIntlClientProvider>
      </body>
    </html>
  );
}
