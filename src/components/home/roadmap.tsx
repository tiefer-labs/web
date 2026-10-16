// SPDX-License-Identifier: MPL-2.0
import { getMessages, getTranslations } from "next-intl/server";
import { Section, SectionLabel } from "@/components/ui/section";
import { SquareMark } from "@/components/ui/square-mark";
import { keysOf } from "@/lib/messages";

/** The milestone being worked on now. */
const CURRENT = "benchmark";

export async function Roadmap() {
  const [t, messages] = await Promise.all([getTranslations("Roadmap"), getMessages()]);

  return (
    <Section id="roadmap" labelledBy="roadmap-title" tone="tint">
      <SectionLabel className="lg:col-[1/3] lg:row-1">{t("label")}</SectionLabel>
      <h2 id="roadmap-title" className="lg:col-[3/10] lg:row-1">
        {t("title")}
      </h2>
      <ol
        role="list"
        className="relative mt-16 grid lg:col-[1/13] lg:grid-cols-5 lg:gap-(--gap) lg:before:absolute lg:before:inset-x-0 lg:before:top-[5px] lg:before:border-t lg:before:border-line lg:before:content-['']"
      >
        {keysOf(messages.Roadmap.items).map((key) => {
          const now = key === CURRENT;
          return (
            <li
              key={key}
              aria-current={now ? "step" : undefined}
              className="relative ml-[5px] border-l border-line pb-8 pl-9 last:border-l-transparent last:pb-0 lg:ml-0 lg:border-l-0 lg:pt-9 lg:pr-4 lg:pb-0 lg:pl-0"
            >
              {now && (
                <span
                  aria-hidden="true"
                  className="absolute top-[0.3em] -left-1.5 size-3 animate-beacon border border-ink lg:top-0 lg:left-0"
                />
              )}
              <SquareMark
                className="top-[0.3em] -left-1.5 lg:top-0 lg:left-0"
                rectClassName={now ? "fill-ink stroke-ink" : "fill-paper stroke-ink-60"}
              />
              {now && (
                <p className="mb-1.5 label-caps text-ink lg:absolute lg:-top-0.5 lg:left-[22px] lg:m-0">
                  {t("nowLabel")}
                </p>
              )}
              <h3>{t(`items.${key}.title`)}</h3>
              <p className="mt-1.5 text-small text-ink-60">{t(`items.${key}.text`)}</p>
            </li>
          );
        })}
      </ol>
    </Section>
  );
}
