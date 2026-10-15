// SPDX-License-Identifier: MPL-2.0
import { getMessages, getTranslations } from "next-intl/server";
import { Section, SectionLabel } from "@/components/ui/section";
import { keysOf } from "@/lib/messages";

export async function UseCases() {
  const [t, messages] = await Promise.all([getTranslations("UseCases"), getMessages()]);

  return (
    <Section id="use-cases" labelledBy="use-cases-title">
      <SectionLabel className="lg:col-[1/3] lg:row-1">{t("label")}</SectionLabel>
      <h2 id="use-cases-title" className="lg:col-[3/9] lg:row-1">
        {t("title")}
      </h2>
      <ul role="list" className="mt-14 grid sm:grid-cols-2 lg:col-[3/13]">
        {keysOf(messages.UseCases.cases).map((key) => (
          <li
            key={key}
            className="border-t border-line py-7 sm:border-t-0 sm:pt-8 sm:pr-8 sm:pb-10 sm:odd:border-r sm:even:pr-0 sm:even:pl-8 sm:nth-[-n+2]:border-b sm:nth-[n+3]:pt-10"
          >
            <p className="mb-3 label-caps text-ink-60">{t(`cases.${key}.label`)}</p>
            <h3>{t(`cases.${key}.title`)}</h3>
            <p className="mt-2 max-w-[28em]">{t(`cases.${key}.text`)}</p>
          </li>
        ))}
      </ul>
      <p className="mt-10 text-small text-ink-60 lg:col-[3/13]">{t("note")}</p>
    </Section>
  );
}
