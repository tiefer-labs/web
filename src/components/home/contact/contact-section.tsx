// SPDX-License-Identifier: MPL-2.0
import { getTranslations } from "next-intl/server";
import { connection } from "next/server";
import { Suspense } from "react";
import { ButtonLink } from "@/components/ui/button";
import { Section, SectionLabel } from "@/components/ui/section";
import { freshToken } from "@/lib/contact/issue-token";
import { contactEnabled, env } from "@/lib/env.server";
import { ContactForm } from "./contact-form";

/** The only per-request part of the page: a fresh signed token for the form. */
async function TokenField() {
  await connection();
  return <input type="hidden" name="t" value={freshToken()} />;
}

const waves = [150, 200, 250, 300];

export async function ContactSection() {
  const [t, footer] = await Promise.all([getTranslations("Contact"), getTranslations("Footer")]);
  const email = env.CONTACT_EMAIL;

  return (
    <Section
      id="contact"
      labelledBy="contact-title"
      tone="tint"
      className="overflow-hidden"
      decoration={
        <svg
          className="pointer-events-none absolute inset-x-0 bottom-0 h-[60%] w-full"
          viewBox="0 0 1440 400"
          preserveAspectRatio="none"
          aria-hidden="true"
          focusable="false"
          data-print="hide"
        >
          {waves.map((y) => (
            <path
              key={y}
              d={`M0 ${y}C240 ${y - 40} 480 ${y - 40} 720 ${y}S1200 ${y + 40} 1440 ${y}`}
              className="fill-none stroke-ink opacity-6"
              strokeWidth="1.5"
              vectorEffect="non-scaling-stroke"
            />
          ))}
        </svg>
      }
    >
      <SectionLabel className="lg:row-1">{t("label")}</SectionLabel>
      <div className="lg:col-[1/6] lg:row-2">
        <h2 id="contact-title">{t("title")}</h2>
        <p className="mt-6 max-w-[30em]">{t("body")}</p>
      </div>
      <div className="mt-12 lg:col-[7/13] lg:row-2 lg:mt-0" data-print="hide">
        {contactEnabled ? (
          <ContactForm
            email={email}
            tokenField={
              <Suspense>
                <TokenField />
              </Suspense>
            }
          />
        ) : (
          <p className="mt-1">
            <ButtonLink href={`mailto:${email}`}>{t("emailUs")}</ButtonLink>
          </p>
        )}
        <dl className="mt-10 grid gap-2 border-t border-line pt-6 text-small">
          <div className="flex flex-wrap gap-x-2 gap-y-1">
            <dt className="text-ink-60">{t("preferEmail")}</dt>
            <dd>
              <a href={`mailto:${email}`} className="font-medium text-ink">
                {email}
              </a>
            </dd>
          </div>
          <div className="flex flex-wrap gap-x-2 gap-y-1">
            <dt className="sr-only">{footer("linkedIn")}</dt>
            <dd>
              <a href={env.LINKEDIN_URL} rel="noopener noreferrer" className="font-medium text-ink">
                {t("linkedIn")}
              </a>
            </dd>
          </div>
        </dl>
      </div>
    </Section>
  );
}
