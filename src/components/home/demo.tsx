// SPDX-License-Identifier: MPL-2.0
import { getMessages, getTranslations } from "next-intl/server";
import { Section, SectionLabel } from "@/components/ui/section";
import { cn } from "@/lib/cn";
import { keysOf } from "@/lib/messages";

/** Rows that carry an evidence tag, so readers can tell what was seen from what was guessed. */
const tagged = { observed: "text-ink", inferred: "text-inferred" } as const;

const isTagged = (key: string): key is keyof typeof tagged => key in tagged;

const barLabel = "fill-ink-60 font-text text-[10px] leading-none font-medium tracking-[0.12em] uppercase";

/** The whole scene the satellite captured; the box marks the part that becomes the alert. */
function Scene({ label, note }: { label: string; note: string }) {
  return (
    <figure>
      <svg className="h-auto w-full" viewBox="0 0 320 220" aria-hidden="true" focusable="false">
        <rect className="fill-ink" width="320" height="220" />
        <path
          className="fill-none stroke-on-ink-muted opacity-40"
          strokeWidth="1.25"
          vectorEffect="non-scaling-stroke"
          d="M0 150C60 130 90 170 150 150C210 130 250 160 320 140M0 60H320M0 100H320M90 0V220M200 0V220M260 0V220"
        />
        <path
          className="fill-on-ink-muted opacity-15"
          d="M200 80C220 50 260 40 290 60C300 80 280 100 250 96C230 94 210 100 200 80Z"
        />
        {/* Smoke drifting north-east from the fire front. */}
        <path className="fill-on-ink-muted opacity-25" d="M150 108C170 80 200 70 230 50C236 64 214 92 180 110Z" />
        <g className="origin-bottom animate-flicker fill-signal [transform-box:fill-box]">
          <path d="M132 120C136 112 142 108 142 100C146 106 150 112 150 120Z" />
          <path d="M146 122C150 112 156 106 156 96C162 104 166 112 164 122Z" />
          <path d="M160 124C163 117 168 113 168 106C172 112 175 118 174 124Z" />
        </g>
        <rect
          className="origin-center animate-beacon fill-none stroke-signal [transform-box:fill-box]"
          strokeWidth="1.5"
          vectorEffect="non-scaling-stroke"
          x="124"
          y="90"
          width="58"
          height="44"
        />
        <rect
          className="fill-none stroke-on-ink"
          strokeWidth="1.5"
          vectorEffect="non-scaling-stroke"
          x="124"
          y="90"
          width="58"
          height="44"
        />
      </svg>
      <figcaption className="mt-3 flex justify-between label-caps text-ink-60">
        <span>{label}</span>
        <span>{note}</span>
      </figcaption>
    </figure>
  );
}

export async function Demo() {
  const [t, messages] = await Promise.all([getTranslations("Demo"), getMessages()]);

  return (
    <Section id="alert" labelledBy="demo-title" tone="tint">
      <SectionLabel className="lg:row-1">{t("label")}</SectionLabel>
      <div className="lg:col-[1/7] lg:row-2">
        <h2 id="demo-title">{t("title")}</h2>
        <p className="mt-4 max-w-[30em] text-lead">{t("body")}</p>
        <div className="mt-10">
          <Scene label={t("sceneLabel")} note={t("sceneNote")} />
        </div>
      </div>
      <figure className="mt-12 self-end border border-line bg-paper shadow-[0_1px_2px_rgb(12_0_61/0.04),0_8px_24px_rgb(12_0_61/0.05)] lg:col-[8/13] lg:row-2 lg:mt-0">
        <figcaption className="flex flex-wrap justify-between gap-x-6 gap-y-2 border-b border-line px-6 py-4">
          <span className="label-caps text-ink">{t("cardTitle")}</span>
          <span className="label-caps tracking-[0.08em] text-ink-60">{t("cardLabel")}</span>
        </figcaption>
        <table className="w-full text-small tabular-nums">
          <tbody>
            {keysOf(messages.Demo.fields).map((key) => (
              <tr key={key} className={cn("not-first:border-t not-first:border-line", isTagged(key) && "bg-tint")}>
                <th scope="row" className="w-[38%] px-6 py-[11px] text-left align-baseline font-normal text-ink-60">
                  {isTagged(key) ? (
                    <span
                      className={cn(
                        "inline-block border border-current px-[7px] pt-[3px] pb-[2px] text-tag leading-[1.3] font-medium tracking-[0.1em] uppercase forced-colors:border-[CanvasText]",
                        tagged[key],
                      )}
                    >
                      {t(`fields.${key}.name`)}
                    </span>
                  ) : (
                    t(`fields.${key}.name`)
                  )}
                </th>
                <td className="py-[11px] pr-6 text-left align-baseline text-ink">{t(`fields.${key}.value`)}</td>
              </tr>
            ))}
          </tbody>
        </table>
        <svg
          className="mx-6 mt-5 mb-6 h-auto w-[calc(100%-48px)] overflow-visible"
          viewBox="0 0 400 56"
          aria-hidden="true"
          focusable="false"
        >
          <text className={barLabel} x="0" y="11">
            {t("sceneLabel")}
          </text>
          <rect className="fill-line" x="0" y="17" width="400" height="8" />
          <text className={barLabel} x="0" y="43">
            {t("alertLabel")}
          </text>
          <rect className="fill-signal" x="0" y="49" width="3" height="8" />
        </svg>
      </figure>
    </Section>
  );
}
