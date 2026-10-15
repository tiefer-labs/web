// SPDX-License-Identifier: MPL-2.0
import { getMessages, getTranslations } from "next-intl/server";
import { Section, SectionLabel } from "@/components/ui/section";
import { keysOf } from "@/lib/messages";

export async function Problem() {
  const [t, messages] = await Promise.all([getTranslations("Problem"), getMessages()]);

  return (
    <Section id="problem" labelledBy="problem-title">
      <SectionLabel className="lg:col-[1/3] lg:row-1">{t("label")}</SectionLabel>
      <h2 id="problem-title" className="lg:col-[3/11] lg:row-1">
        {t("title")}
      </h2>
      <p className="mt-10 max-w-[34em] font-headline text-statement leading-[1.25] font-medium tracking-[-0.01em] text-ink lg:col-[3/13]">
        {t("intro")}
      </p>
      <ol className="mt-16 counter-steps grid gap-8 border-t border-line pt-8 *:before:mb-4 sm:grid-cols-3 lg:col-[3/13]">
        {keysOf(messages.Problem.items).map((key) => (
          <li key={key}>
            <h3>{t(`items.${key}.title`)}</h3>
            <p className="mt-2">{t(`items.${key}.text`)}</p>
          </li>
        ))}
      </ol>
    </Section>
  );
}
