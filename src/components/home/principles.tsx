// SPDX-License-Identifier: MPL-2.0
import { getMessages, getTranslations } from "next-intl/server";
import { Section, SectionLabel } from "@/components/ui/section";
import { keysOf } from "@/lib/messages";

export async function Principles() {
  const [t, messages] = await Promise.all([getTranslations("Principles"), getMessages()]);

  return (
    <Section id="principles" labelledBy="principles-title" tone="ink">
      <SectionLabel className="lg:col-[1/3] lg:row-1">{t("label")}</SectionLabel>
      <h2 id="principles-title" className="text-on-ink lg:col-[3/11] lg:row-1">
        {t("title")}
      </h2>
      <ol role="list" className="mt-14 grid gap-x-(--gap) sm:grid-cols-2 lg:col-[1/13] lg:grid-cols-3">
        {keysOf(messages.Principles.items).map((key) => (
          <li key={key} className="border-t border-line-on-ink pt-7 pb-12">
            <span className="label-caps text-signal">{t(`items.${key}.number`)}</span>
            <h3 className="mt-4 max-w-[12em] text-statement leading-[1.15] tracking-[-0.01em] text-on-ink">
              {t(`items.${key}.statement`)}
            </h3>
          </li>
        ))}
      </ol>
    </Section>
  );
}
