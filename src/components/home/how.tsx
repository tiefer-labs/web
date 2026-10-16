// SPDX-License-Identifier: MPL-2.0
import { getMessages, getTranslations } from "next-intl/server";
import type { ReactNode } from "react";
import { Section, SectionLabel } from "@/components/ui/section";
import { keysOf } from "@/lib/messages";

const line = "fill-none stroke-on-ink-muted";
const tag = "fill-on-ink-muted font-text text-[9px] font-medium tracking-[0.14em]";
const svgBox = "origin-center [transform-box:fill-box]";

type Words = { kept: string; packet: string; signed: string };

/** One looping picture per stage, drawn in a 160 by 120 box. */
const STAGE_ART: Record<string, (words: Words) => ReactNode> = {
  // A cloudy frame drops through the filter into the frames kept on board.
  filter: ({ kept }) => (
    <>
      <path className={line} strokeDasharray="3 4" d="M20 56H140" />
      <g className="animate-frame-drop">
        <rect className="fill-ink stroke-on-ink" x="58" y="12" width="44" height="30" />
        <g className="fill-paper opacity-90">
          <circle cx="72" cy="28" r="7" />
          <circle cx="82" cy="24" r="8" />
          <circle cx="90" cy="30" r="6" />
        </g>
      </g>
      <rect className={line} x="56" y="74" width="48" height="8" />
      <rect className={line} x="56" y="84" width="48" height="8" />
      <rect className={line} x="56" y="94" width="48" height="8" />
      <text className={tag} x="80" y="116" textAnchor="middle">
        {kept}
      </text>
    </>
  ),
  // The detector locks onto a fire in the frame.
  detect: () => (
    <>
      <rect className="fill-ink stroke-on-ink-muted" x="30" y="16" width="100" height="78" />
      <path className={line} d="M30 70C50 60 64 76 84 66C100 58 112 64 130 56M54 16V94" />
      <path
        className="fill-signal"
        d="M80 38C83 43 88 46 88 51C88 55 84 58 80 58C76 58 72 55 72 51C72 47 76 45 77 41C78 43 79 44 80 44Z"
      />
      <path
        className={`animate-lock-on fill-none stroke-on-ink ${svgBox}`}
        strokeWidth="1.5"
        d="M62 34V28H68M92 28H98V34M98 62V68H92M68 68H62V62"
      />
    </>
  ),
  // A small packet leaves the frame for the ground antenna.
  alert: ({ packet }) => (
    <>
      <rect className="fill-ink stroke-on-ink-muted" x="14" y="62" width="40" height="40" />
      <path
        className="fill-signal"
        d="M34 72C36 75 39 77 39 80C39 83 37 85 34 85C31 85 29 83 29 80C29 77 31 76 32 74C32.5 75 33.3 75.6 34 75.6Z"
      />
      <path className={line} d="M118 46L128 18L138 46M120 40H136M128 18V10" />
      <path className={line} d="M116 10C122 4 134 4 140 10" />
      <path className={line} strokeDasharray="3 4" d="M48 70L116 28" />
      <rect className="animate-packet-out fill-signal" x="44" y="68" width="8" height="8" />
      <text
        className="fill-signal font-text text-[10px] font-medium tracking-[0.14em]"
        x="128"
        y="62"
        textAnchor="middle"
      >
        {packet}
      </text>
    </>
  ),
  // A signed model comes down into the satellite.
  update: ({ signed }) => (
    <>
      <path className={line} strokeDasharray="3 4" d="M80 4V60" />
      <g className="animate-update-in">
        <rect className="fill-ink stroke-on-ink" x="62" y="22" width="36" height="26" />
        <path className="fill-none stroke-signal" strokeWidth="1.5" d="M71 35L77 41L89 29" />
      </g>
      <rect className="fill-on-ink" x="68" y="72" width="24" height="22" />
      <rect className="fill-none stroke-on-ink" strokeWidth="1.5" x="28" y="78" width="36" height="10" />
      <rect className="fill-none stroke-on-ink" strokeWidth="1.5" x="96" y="78" width="36" height="10" />
      <text className={tag} x="80" y="114" textAnchor="middle">
        {signed}
      </text>
    </>
  ),
};

export async function How() {
  const [t, messages] = await Promise.all([getTranslations("How"), getMessages()]);
  const words = { kept: t("kept"), packet: t("packet"), signed: t("signed") };

  return (
    <Section id="how" labelledBy="how-title" tone="ink">
      <SectionLabel className="lg:col-[1/3] lg:row-1">{t("label")}</SectionLabel>
      <h2 id="how-title" className="text-on-ink lg:col-[3/11] lg:row-1">
        {t("title")}
      </h2>
      <div className="relative mt-14 lg:col-[1/13]">
        {/* A frame travelling through the stages. */}
        <div
          aria-hidden="true"
          className="absolute inset-x-0 top-0 hidden h-px bg-line-on-ink lg:block"
          data-print="hide"
        >
          <span className="absolute -top-[3px] size-[7px] -translate-x-1/2 animate-travel bg-signal" />
        </div>
        <ol role="list" className="grid gap-10 sm:grid-cols-2 lg:grid-cols-4 lg:gap-(--gap) lg:pt-8">
          {keysOf(messages.How.stages).map((key) => (
            <li key={key}>
              <svg
                className="h-auto w-full max-w-60"
                viewBox="0 0 160 120"
                strokeWidth="1.25"
                aria-hidden="true"
                focusable="false"
                data-print="hide"
              >
                {STAGE_ART[key]?.(words)}
              </svg>
              <p className="mt-6 label-caps text-on-ink-muted">{t(`stages.${key}.number`)}</p>
              <h3 className="mt-2 text-h3 text-on-ink">{t(`stages.${key}.name`)}</h3>
              <p className="mt-1.5 max-w-[18em] text-small">{t(`stages.${key}.text`)}</p>
            </li>
          ))}
        </ol>
      </div>
      <div className="mt-16 flex flex-wrap items-center gap-x-4 gap-y-3 border-t border-line-on-ink pt-8 lg:col-[1/13]">
        <p className="text-on-ink">{t("hardware")}</p>
        <ul role="list" className="flex gap-2">
          {keysOf(messages.How.chips).map((key) => (
            <li key={key} className="border border-line-on-ink px-3 py-1 label-caps text-on-ink">
              {t(`chips.${key}`)}
            </li>
          ))}
        </ul>
      </div>
      <div className="mt-12 lg:col-[1/13]">
        <h3 className="label-caps leading-[1.2] text-on-ink-muted">{t("groundTitle")}</h3>
        <dl className="mt-4 grid gap-6 sm:grid-cols-2">
          {keysOf(messages.How.ground).map((key) => (
            <div key={key} className="border-l-2 border-signal pl-4">
              <dt className="font-headline text-h3 leading-[1.2] font-medium text-on-ink">
                {t(`ground.${key}.title`)}
              </dt>
              <dd className="mt-1 text-small">{t(`ground.${key}.text`)}</dd>
            </div>
          ))}
        </dl>
      </div>
    </Section>
  );
}
