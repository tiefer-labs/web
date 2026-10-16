// SPDX-License-Identifier: MPL-2.0
import { getTranslations } from "next-intl/server";
import { Section, SectionLabel } from "@/components/ui/section";
import { Link } from "@/i18n/navigation";
import { siteUrl } from "@/lib/env.server";

export async function NotFoundContent() {
  const t = await getTranslations("NotFound");

  return (
    <Section labelledBy="nf-title" padding="pt-30 pb-40">
      <SectionLabel className="lg:col-[1/3] lg:row-1">{t("code")}</SectionLabel>
      <h1 id="nf-title" className="text-h2 leading-[1.06] font-medium tracking-[-0.018em] lg:col-[3/11]">
        {t("title")}
      </h1>
      <p className="mt-6 text-lead lg:col-[3/11]">{t("body")}</p>
      <p className="mt-10 lg:col-[3/11]">
        <Link href="/" className="font-medium">
          {t("home")}
        </Link>
      </p>
    </Section>
  );
}

export async function notFoundMetadata() {
  const [t, meta] = await Promise.all([getTranslations("NotFound"), getTranslations("Meta")]);
  return {
    metadataBase: new URL(siteUrl),
    title: { absolute: `${t("title")} | ${meta("siteName")}` },
  };
}
