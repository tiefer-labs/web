// SPDX-License-Identifier: MPL-2.0
import { getTranslations } from "next-intl/server";
import { ButtonLink } from "@/components/ui/button";
import { HeroGlobe } from "./hero-globe";
import { HeroOrbit } from "./hero-orbit";

export async function Hero() {
  const t = await getTranslations("Hero");

  return (
    <section
      aria-labelledby="hero-title"
      data-tone="dark"
      className="relative flex min-h-[clamp(660px,92vh,860px)] items-end overflow-hidden bg-ink pt-32 pb-26 text-on-ink-muted max-lg:min-h-0 max-lg:items-start max-lg:pt-30 max-lg:pb-[300px]"
    >
      <HeroGlobe />
      <HeroOrbit
        labels={{
          scanning: t("scanning"),
          detected: t("detected"),
          alerted: t("alerted"),
          fire: t("fire"),
          responder: t("responder"),
          stored: t("stored"),
          sent: t("sent"),
          alertSize: t("alertSize"),
          alertTag: t("alertTag"),
        }}
      />
      <div className="relative wrap grid-12">
        <div className="max-w-[36rem] lg:col-[1/7]">
          <p className="mb-6 label-caps text-on-ink-muted">{t("eyebrow")}</p>
          <h1 id="hero-title" className="text-h1 leading-[1.02] font-semibold tracking-[-0.025em] text-on-ink">
            {t("title")}
          </h1>
          <p className="mt-7 max-w-[34em] text-lead leading-[1.5] text-on-ink-muted">{t("lead")}</p>
          <p className="mt-10 flex flex-wrap items-center gap-x-8 gap-y-4">
            <ButtonLink href="#contact" tone="dark">
              {t("primary")}
            </ButtonLink>
            <a href="#how" className="font-medium text-on-ink">
              {t("secondary")}
            </a>
          </p>
        </div>
      </div>
    </section>
  );
}
