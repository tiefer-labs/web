// SPDX-License-Identifier: MPL-2.0
import type { Metadata } from "next";
import { getLocale, getTranslations } from "next-intl/server";
import { ContactSection } from "@/components/home/contact/contact-section";
import { Demo } from "@/components/home/demo";
import { Hero } from "@/components/home/hero";
import { How } from "@/components/home/how";
import { Name } from "@/components/home/name";
import { Principles } from "@/components/home/principles";
import { Problem } from "@/components/home/problem";
import { Product } from "@/components/home/product";
import { Roadmap } from "@/components/home/roadmap";
import { UseCases } from "@/components/home/use-cases";
import { SiteHeader } from "@/components/layout/site-header";
import { JsonLd } from "@/components/ui/json-ld";
import { buildMetadata } from "@/lib/seo";
import { organizationSchema } from "@/lib/structured-data";

export async function generateMetadata(): Promise<Metadata> {
  return buildMetadata({ locale: await getLocale(), path: "/" });
}

export default async function HomePage() {
  const meta = await getTranslations("Meta");

  return (
    <>
      <JsonLd data={organizationSchema(meta("siteName"))} />
      <SiteHeader dark />
      <main id="main" tabIndex={-1}>
        <Hero />
        <Problem />
        <Product />
        <How />
        <Demo />
        <UseCases />
        <Principles />
        <Roadmap />
        <Name />
        <ContactSection />
      </main>
    </>
  );
}
