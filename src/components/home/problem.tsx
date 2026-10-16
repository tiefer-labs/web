// SPDX-License-Identifier: MPL-2.0
import { getMessages, getTranslations } from "next-intl/server";
import type { ReactNode } from "react";
import { Section, SectionLabel } from "@/components/ui/section";
import { cn } from "@/lib/cn";
import { keysOf } from "@/lib/messages";

const COLUMNS = 8;
const ROWS = 4;
const TILE = 48;
const GAP = 6;
// Every tile starts its cloud loop at a different point, so about two thirds are covered at any moment.
const DELAYS = [
  "[animation-delay:-0.4s]",
  "[animation-delay:-7.1s]",
  "[animation-delay:-3.3s]",
  "[animation-delay:-10.2s]",
  "[animation-delay:-5.6s]",
  "[animation-delay:-1.9s]",
  "[animation-delay:-8.8s]",
  "[animation-delay:-4.4s]",
  "[animation-delay:-11.5s]",
  "[animation-delay:-2.7s]",
  "[animation-delay:-9.6s]",
  "[animation-delay:-6.3s]",
];
// Fields, a river and a coast: three kinds of ground, so the grid reads as imagery.
const GROUND = [
  "M6 14H42M6 24H42M6 34H42M24 6V42",
  "M4 30C14 22 20 34 30 24C36 18 40 20 44 16",
  "M4 8C16 14 18 28 30 30C36 31 40 38 44 44",
];

/** A satellite's frames, with clouds drifting over most of them. */
function FrameGrid({ label }: { label: string }) {
  const width = COLUMNS * (TILE + GAP) - GAP;
  const height = ROWS * (TILE + GAP) - GAP;
  return (
    <figure>
      <svg className="h-auto w-full" viewBox={`0 0 ${width} ${height}`} aria-hidden="true" focusable="false">
        {Array.from({ length: COLUMNS * ROWS }, (_, i) => {
          const x = (i % COLUMNS) * (TILE + GAP);
          const y = Math.floor(i / COLUMNS) * (TILE + GAP);
          return (
            <g key={i} transform={`translate(${x} ${y})`}>
              <rect className="fill-ink" width={TILE} height={TILE} />
              <path
                className="fill-none stroke-on-ink-muted opacity-50"
                strokeWidth="1.25"
                vectorEffect="non-scaling-stroke"
                d={GROUND[(i * 5) % GROUND.length]}
              />
              <g className={cn("animate-cloud-pass fill-paper", DELAYS[(i * 7) % DELAYS.length])}>
                <circle cx="16" cy="26" r="11" />
                <circle cx="28" cy="20" r="13" />
                <circle cx="36" cy="30" r="9" />
                <rect x="10" y="28" width="30" height="10" />
              </g>
            </g>
          );
        })}
      </svg>
      <figcaption className="mt-3 label-caps text-ink-60">{label}</figcaption>
    </figure>
  );
}

const ICONS: Record<string, ReactNode> = {
  delay: (
    <>
      <circle cx="16" cy="16" r="12" />
      <path d="M16 9V16L21 19" />
    </>
  ),
  downlink: (
    <>
      <path d="M16 4V22M10 16L16 22L22 16" />
      <path d="M6 28H26" strokeDasharray="2 3" />
    </>
  ),
  passive: (
    <>
      <rect x="4" y="9" width="24" height="16" />
      <circle cx="16" cy="17" r="5" />
      <path d="M11 9L13 5H19L21 9" />
    </>
  ),
};

export async function Problem() {
  const [t, messages] = await Promise.all([getTranslations("Problem"), getMessages()]);

  return (
    <Section id="problem" labelledBy="problem-title">
      <SectionLabel className="lg:col-[1/3] lg:row-1">{t("label")}</SectionLabel>
      <h2 id="problem-title" className="lg:col-[3/11] lg:row-1">
        {t("title")}
      </h2>
      <p className="mt-12 lg:col-[1/6] lg:row-2 lg:self-center">
        <span className="block font-headline text-[clamp(5rem,4rem+5vw,9rem)] leading-none font-semibold tracking-[-0.04em] text-ink">
          {t("stat")}
        </span>
        <span className="mt-3 block max-w-[16em] text-lead text-ink">{t("statText")}</span>
      </p>
      <div className="mt-12 lg:col-[6/13] lg:row-2">
        <FrameGrid label={t("gridLabel")} />
      </div>
      <ul role="list" className="mt-16 grid gap-8 border-t border-line pt-10 sm:grid-cols-3 lg:col-[1/13]">
        {keysOf(messages.Problem.items).map((key) => (
          <li key={key} className="flex gap-4">
            <svg
              className="size-8 shrink-0 fill-none stroke-ink"
              viewBox="0 0 32 32"
              strokeWidth="1.5"
              aria-hidden="true"
              focusable="false"
            >
              {ICONS[key]}
            </svg>
            <div>
              <h3>{t(`items.${key}.title`)}</h3>
              <p className="mt-1 text-small">{t(`items.${key}.text`)}</p>
            </div>
          </li>
        ))}
      </ul>
    </Section>
  );
}
