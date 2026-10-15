// SPDX-License-Identifier: MPL-2.0
import { getMessages, getTranslations } from "next-intl/server";
import { Section, SectionLabel } from "@/components/ui/section";
import { SquareMark } from "@/components/ui/square-mark";
import { keysOf } from "@/lib/messages";

/** The filter step branches off: frames it rejects stay on board. */
const BRANCH_STEP = "filter";

export async function How() {
  const [t, messages] = await Promise.all([getTranslations("How"), getMessages()]);
  const steps = keysOf(messages.How.steps);

  return (
    <Section id="how" labelledBy="how-title">
      <SectionLabel className="lg:col-[1/3] lg:row-1">{t("label")}</SectionLabel>
      <h2 id="how-title" className="max-w-[16em] lg:col-[3/10] lg:row-1">
        {t("title")}
      </h2>
      <div className="relative mt-16 lg:col-[1/13]">
        <svg
          className="absolute inset-x-0 top-[5px] hidden h-[2px] w-full overflow-visible lg:block"
          viewBox="0 0 100 2"
          preserveAspectRatio="none"
          aria-hidden="true"
          focusable="false"
        >
          <path d="M0 1H100" className="stroke-ink-60" strokeWidth="1.25" vectorEffect="non-scaling-stroke" />
        </svg>
        <ol className="counter-steps grid *:before:mb-1.5 lg:grid-cols-5 lg:gap-(--gap)">
          {steps.map((key, index) => (
            <li
              key={key}
              className="relative ml-[5px] border-l border-ink-60 pb-8 pl-9 last:border-l-transparent last:pb-0 lg:ml-0 lg:border-l-0 lg:pt-9 lg:pr-4 lg:pb-0 lg:pl-0"
            >
              <SquareMark
                className="top-[0.35em] -left-1.5 bg-paper lg:top-0 lg:left-0"
                rectClassName={index === steps.length - 1 ? "fill-ink stroke-ink" : "fill-paper stroke-ink"}
              />
              <h3>{t(`steps.${key}.title`)}</h3>
              <p className="mt-1.5 text-small">{t(`steps.${key}.text`)}</p>
              {key === BRANCH_STEP && (
                <p
                  aria-hidden="true"
                  className="mt-3.5 border-l border-dashed border-ink-60 pt-2.5 pl-3 text-label leading-[1.3] font-medium tracking-[0.12em] text-ink-60"
                >
                  {t("branch")}
                </p>
              )}
            </li>
          ))}
        </ol>
      </div>
      <p className="mt-12 max-w-[40em] lg:col-[3/10]">{t("note")}</p>
      <div className="mt-18 lg:col-[3/13]">
        <h3 className="label-caps leading-[1.2] text-ink-60">{t("groundTitle")}</h3>
        <dl className="mt-6 grid sm:grid-cols-2 sm:border-t sm:border-line">
          {keysOf(messages.How.ground).map((key) => (
            <div
              key={key}
              className="border-t border-line py-6 sm:border-t-0 sm:pt-6 sm:pr-8 sm:pb-0 sm:[&+div]:border-l sm:[&+div]:pl-8"
            >
              <dt className="font-headline text-h3 leading-[1.2] font-medium text-ink">{t(`ground.${key}.title`)}</dt>
              <dd className="mt-2 max-w-[30em]">{t(`ground.${key}.text`)}</dd>
            </div>
          ))}
        </dl>
      </div>
    </Section>
  );
}
