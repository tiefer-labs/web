// SPDX-License-Identifier: MPL-2.0
import { getTranslations } from "next-intl/server";
import { ButtonLink } from "@/components/ui/button";
import { HeroGlobe } from "./hero-globe";
import { type EventLabels, HeroOrbit } from "./hero-orbit";
import type { EventId } from "./hero-scene";

export async function Hero() {
  const t = await getTranslations("Hero");
  const eventLabels = (event: EventId): EventLabels => ({
    scanning: t(`events.${event}.scanning`),
    detected: t(`events.${event}.detected`),
    alerted: t(`events.${event}.alerted`),
    event: t(`events.${event}.event`),
    responder: t(`events.${event}.responder`),
  });

  return (
    <section
      aria-labelledby="hero-title"
      data-tone="dark"
      className="relative flex min-h-[clamp(660px,92vh,860px)] items-end overflow-hidden bg-ink pt-32 pb-26 text-on-ink-muted max-lg:min-h-0 max-lg:items-start max-lg:pt-30 max-lg:pb-[300px]"
    >
      <HeroGlobe />
      <HeroOrbit
        labels={{
          stored: t("stored"),
          sent: t("sent"),
          sentSizes: [t("sentSize", { size: 48 }), t("sentSize", { size: 96 })],
          alertTag: t("alertTag"),
          ais: t("ais"),
          events: {
            wildfire: eventLabels("wildfire"),
            oilSpill: eventLabels("oilSpill"),
            vessels: eventLabels("vessels"),
          },
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
