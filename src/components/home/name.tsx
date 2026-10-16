// SPDX-License-Identifier: MPL-2.0
import { getTranslations } from "next-intl/server";
import Image from "next/image";
import { Section, SectionLabel } from "@/components/ui/section";

export async function Name() {
  const t = await getTranslations("Name");

  return (
    <Section
      id="name"
      labelledBy="name-title"
      tone="ink"
      padding="pt-28 pb-40"
      className="overflow-hidden"
      decoration={
        <>
          <Image
            src="/brand/tiefer-mark-white.svg"
            alt=""
            width={690}
            height={690}
            data-print="hide"
            data-forced-colors="keep"
            className="pointer-events-none absolute top-1/2 right-(--pad) h-auto w-[min(28vw,240px)] -translate-y-1/2 opacity-10 max-lg:top-26 max-lg:w-16 max-lg:translate-y-0 max-lg:opacity-16 lg:right-[max(var(--pad),calc((100%-var(--max))/2))] lg:w-[280px]"
          />
          <svg
            className="absolute inset-x-0 bottom-0 h-24 w-full"
            viewBox="0 0 1440 120"
            preserveAspectRatio="none"
            aria-hidden="true"
            focusable="false"
            data-print="hide"
          >
            <path
              d="M-20 120C320 26 1120 26 1460 120"
              className="fill-none stroke-on-ink-muted opacity-40"
              strokeWidth="1.25"
              vectorEffect="non-scaling-stroke"
            />
          </svg>
        </>
      }
    >
      <SectionLabel className="lg:col-[1/3] lg:row-1">{t("label")}</SectionLabel>
      <h2
        id="name-title"
        className="text-h2 leading-[1.06] font-medium tracking-[-0.018em] text-on-ink lg:col-[3/9] lg:row-1 lg:text-name"
      >
        {t("line")}
      </h2>
      <p className="mt-6 max-w-[30em] text-lead leading-[1.5] lg:col-[3/8]">{t("body")}</p>
    </Section>
  );
}
