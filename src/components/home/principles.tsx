// SPDX-License-Identifier: MPL-2.0
import { getMessages, getTranslations } from "next-intl/server";
import { Section, SectionLabel } from "@/components/ui/section";
import { keysOf } from "@/lib/messages";

export async function Principles() {
  const [t, messages] = await Promise.all([getTranslations("Principles"), getMessages()]);

  return (
    <Section id="principles" labelledBy="principles-title" tone="tint">
      <SectionLabel className="lg:row-1">{t("label")}</SectionLabel>
      <h2 id="principles-title" className="lg:col-[1/6] lg:row-2">
        {t("title")}
      </h2>
      <ol role="list" className="mt-14 grid gap-x-16 lg:col-[1/13] lg:grid-cols-2">
        {keysOf(messages.Principles.items).map((key) => (
          <li
            key={key}
            className="grid grid-cols-[3.5rem_minmax(0,1fr)] border-t border-line pt-6 pb-7 lg:nth-last-[-n+2]:border-b"
          >
            <span className="pt-[0.45em] label-caps text-ink-60">{t(`items.${key}.number`)}</span>
            <div>
              <h3>{t(`items.${key}.statement`)}</h3>
              <p className="mt-2">{t(`items.${key}.text`)}</p>
            </div>
          </li>
        ))}
      </ol>
    </Section>
  );
}
