// SPDX-License-Identifier: MPL-2.0
import type { Metadata } from "next";
import type { Locale } from "next-intl";
import { getMessages, getTranslations } from "next-intl/server";
import type { ReactNode } from "react";
import { Section, SectionLabel } from "@/components/ui/section";
import { env, siteUrl } from "@/lib/env.server";
import { keysOf } from "@/lib/messages";
import { buildMetadata } from "@/lib/seo";
import { markPlaceholders } from "./mark-placeholders";

export type LegalPageKey = "notice" | "privacy" | "acceptableUse";

type LegalSection = {
  facts?: Record<string, unknown>;
  paras?: Record<string, unknown>;
  list?: Record<string, unknown>;
};

type LegalTranslator = {
  (key: string): string;
  rich(key: string, values: ReturnType<typeof legalValues>): ReactNode;
};

/** Config values for the `{tokens}` in the legal copy, and the tags that turn them into links. */
function legalValues() {
  return {
    legal_name: env.LEGAL_NAME,
    legal_form: env.LEGAL_FORM,
    address: env.LEGAL_ADDRESS,
    tax_id: env.LEGAL_TAX_ID,
    registration: env.LEGAL_REGISTRATION,
    director: env.LEGAL_DIRECTOR,
    hosting_provider: env.LEGAL_HOSTING_PROVIDER,
    hosting_country: env.LEGAL_HOSTING_COUNTRY,
    smtp_provider: env.LEGAL_SMTP_PROVIDER,
    smtp_country: env.LEGAL_SMTP_COUNTRY,
    log_retention_days: String(env.LOG_RETENTION_DAYS),
    contact_email: env.CONTACT_EMAIL,
    site_url: siteUrl,
    email: (chunks: ReactNode) => <a href={`mailto:${env.CONTACT_EMAIL}`}>{chunks}</a>,
    site: (chunks: ReactNode) => <a href={`${siteUrl}/`}>{chunks}</a>,
  };
}

export async function legalMetadata(locale: Locale, page: LegalPageKey, path: string): Promise<Metadata> {
  const t = await getTranslations({ locale, namespace: `Legal.${page}` });
  return buildMetadata({ locale, path, title: t("title"), description: t("description") });
}

/** Pages with more sections than this get a table of contents. */
const TOC_AFTER = 4;

export async function LegalPage({ page }: { page: LegalPageKey }) {
  const [pageT, ui, legal, messages] = await Promise.all([
    getTranslations(`Legal.${page}`),
    getTranslations("UI"),
    getTranslations("Legal"),
    getMessages(),
  ]);
  // Each page has its own sections, so their keys are only known from the messages at run time.
  // oxlint-disable-next-line typescript/no-unsafe-type-assertion
  const t = pageT as unknown as LegalTranslator;
  const values = legalValues();
  const rich = (key: string) => markPlaceholders(t.rich(key, values));
  const sections: Record<string, LegalSection> = messages.Legal[page].sections;
  const keys = keysOf(sections);

  return (
    <Section as="article" labelledBy="legal-title" padding="pt-18 pb-(--section)">
      <SectionLabel className="lg:col-[1/3] lg:row-1">{ui("legalLabel")}</SectionLabel>
      <div className="lg:col-[3/11]">
        {!env.LEGAL_REVIEWED && (
          <p role="note" className="mb-10 border-l-2 border-inferred bg-tint px-4 py-3 font-medium text-ink">
            {legal("reviewNote")}
          </p>
        )}
        <h1 id="legal-title" className="text-h2 leading-[1.06] font-medium tracking-[-0.018em]">
          {t("title")}
        </h1>
        <p className="mt-6 max-w-[36em] text-lead leading-normal">{rich("intro")}</p>
        {keys.length > TOC_AFTER && (
          <nav aria-labelledby="legal-toc-title" className="mt-10 border-t border-line pt-4">
            <h2 id="legal-toc-title" className="label-caps leading-[1.06] text-ink-60">
              {legal("contents")}
            </h2>
            <ol className="mt-3 list-decimal columns-[2_14rem] gap-x-(--gap) pl-[1.4em] text-small">
              {keys.map((key) => (
                <li key={key} className="break-inside-avoid py-0.5">
                  <a href={`#${key}`}>{t(`sections.${key}.title`)}</a>
                </li>
              ))}
            </ol>
          </nav>
        )}
        {keys.map((key) => {
          const section: LegalSection = sections[key];
          return (
            <section
              key={key}
              id={key}
              aria-labelledby={`${key}-title`}
              className="mt-12 max-w-176 border-t border-line pt-6"
            >
              <h2 id={`${key}-title`} className="text-h3 leading-[1.2] tracking-[-0.005em]">
                {t(`sections.${key}.title`)}
              </h2>
              {section.facts && (
                <dl className="mt-4">
                  {keysOf(section.facts).map((fact) => (
                    <div
                      key={fact}
                      className="grid gap-x-6 gap-y-0.5 border-b border-line py-2.5 sm:grid-cols-[14rem_1fr]"
                    >
                      <dt className="text-small text-ink-60">{t(`sections.${key}.facts.${fact}.label`)}</dt>
                      <dd className="text-ink">{rich(`sections.${key}.facts.${fact}.value`)}</dd>
                    </div>
                  ))}
                </dl>
              )}
              {section.paras &&
                keysOf(section.paras).map((para) => (
                  <p key={para} className="mt-4">
                    {rich(`sections.${key}.paras.${para}`)}
                  </p>
                ))}
              {section.list && (
                <ul className="mt-4 list-disc space-y-2 pl-[1.2em]">
                  {keysOf(section.list).map((item) => (
                    <li key={item}>{rich(`sections.${key}.list.${item}`)}</li>
                  ))}
                </ul>
              )}
            </section>
          );
        })}
        <p className="mt-12 text-small text-ink-60">{rich("updated")}</p>
      </div>
    </Section>
  );
}
