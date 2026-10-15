// SPDX-License-Identifier: MPL-2.0
import { getMessages, getTranslations } from "next-intl/server";
import { Section, SectionLabel } from "@/components/ui/section";
import { cn } from "@/lib/cn";
import { keysOf } from "@/lib/messages";

/** Rows that carry an evidence tag, so readers can tell what was seen from what was guessed. */
const tagged = { observed: "text-ink", inferred: "text-inferred" } as const;

const isTagged = (key: string): key is keyof typeof tagged => key in tagged;

const barLabel = "fill-ink-60 font-text text-[10px] leading-none font-medium tracking-[0.12em] uppercase";

export async function Demo() {
  const [t, messages] = await Promise.all([getTranslations("Demo"), getMessages()]);

  return (
    <Section id="alert" labelledBy="demo-title" tone="tint">
      <SectionLabel className="lg:row-1">{t("label")}</SectionLabel>
      <div className="lg:col-[1/6] lg:row-2">
        <h2 id="demo-title">{t("title")}</h2>
        <p className="mt-6 max-w-[30em]">{t("body")}</p>
      </div>
      <figure className="mt-12 border border-line bg-paper shadow-[0_1px_2px_rgb(12_0_61/0.04),0_8px_24px_rgb(12_0_61/0.05)] lg:col-[7/13] lg:row-2 lg:mt-0">
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
          <rect className="fill-ink" x="0" y="49" width="3" height="8" />
        </svg>
      </figure>
    </Section>
  );
}
