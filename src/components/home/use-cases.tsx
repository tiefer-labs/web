// SPDX-License-Identifier: MPL-2.0
import { getMessages, getTranslations } from "next-intl/server";
import type { ReactNode } from "react";
import { Section, SectionLabel } from "@/components/ui/section";
import { keysOf } from "@/lib/messages";

const line = "fill-none stroke-ink";
const box = "[transform-box:fill-box]";

/** One small looping picture per event, drawn in a 120 by 80 box. */
const ART: Record<string, ReactNode> = {
  wildfires: (
    <>
      <path className={line} d="M8 66H112M20 66V52L28 44L36 52V66M84 66V48L92 40L100 48V66" />
      <g className={`origin-bottom animate-flicker fill-signal ${box}`}>
        <path d="M52 66C46 58 46 50 54 40C56 48 60 50 60 42C68 50 70 58 64 66Z" />
      </g>
    </>
  ),
  oilSpills: (
    <>
      <path className="fill-line" d="M0 40H120V80H0Z" />
      <ellipse className={`origin-center animate-spread fill-ink opacity-80 ${box}`} cx="64" cy="58" rx="34" ry="10" />
      <path className={line} d="M84 40V16H104V40M84 24H104M94 16V8" />
    </>
  ),
  maritime: (
    <>
      <path className="fill-line" d="M0 44H120V80H0Z" />
      <circle className={`origin-center animate-signal fill-none stroke-ink ${box}`} cx="34" cy="50" r="16" />
      <path className="fill-ink" d="M22 50H46L42 58H26Z" />
      <path className="fill-ink-60" d="M74 50H98L94 58H78Z" />
      <path
        className="fill-none stroke-signal"
        strokeWidth="1.5"
        d="M70 40V36H74M98 36H102V40M102 62V66H98M74 66H70V62"
      />
    </>
  ),
  floods: (
    <>
      <path className={line} d="M40 52V32L56 20L72 32V52" />
      <path className={line} d="M8 64H112" strokeDasharray="3 4" />
      <path
        className="animate-tide fill-ink opacity-80"
        d="M0 56C10 52 20 60 30 56C40 52 50 60 60 56C70 52 80 60 90 56C100 52 110 60 120 56V80H0Z"
      />
    </>
  ),
};

export async function UseCases() {
  const [t, messages] = await Promise.all([getTranslations("UseCases"), getMessages()]);

  return (
    <Section id="use-cases" labelledBy="use-cases-title">
      <SectionLabel className="lg:col-[1/3] lg:row-1">{t("label")}</SectionLabel>
      <h2 id="use-cases-title" className="lg:col-[3/11] lg:row-1">
        {t("title")}
      </h2>
      <ul role="list" className="mt-14 grid gap-10 sm:grid-cols-2 lg:col-[1/13] lg:grid-cols-4 lg:gap-(--gap)">
        {keysOf(messages.UseCases.cases).map((key) => (
          <li key={key} className="bg-tint p-6">
            <svg
              className="h-auto w-full overflow-hidden bg-paper"
              viewBox="0 0 120 80"
              strokeWidth="1.25"
              aria-hidden="true"
              focusable="false"
              data-print="hide"
            >
              {ART[key]}
            </svg>
            <h3 className="mt-6 text-h3">{t(`cases.${key}.title`)}</h3>
            <p className="mt-1.5 text-small">{t(`cases.${key}.text`)}</p>
          </li>
        ))}
      </ul>
      <p className="mt-12 text-small text-ink-60 lg:col-[1/13]">{t("note")}</p>
    </Section>
  );
}
