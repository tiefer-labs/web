// SPDX-License-Identifier: MPL-2.0
import { getTranslations } from "next-intl/server";
import Image from "next/image";
import { Link } from "@/i18n/navigation";
import { env } from "@/lib/env.server";
import { ExternalFooterLink, FooterLink } from "./footer-link";

export async function SiteFooter() {
  const [t, ui, meta] = await Promise.all([getTranslations("Footer"), getTranslations("UI"), getTranslations("Meta")]);

  return (
    <footer className="border-t border-line pt-12 pb-10 text-small">
      <div className="wrap grid gap-5 lg:grid-cols-12 lg:items-baseline lg:gap-x-(--gap)">
        <Link href="/" className="block lg:col-[1/3]" data-forced-colors="keep">
          <Image src="/brand/tiefer-logo.svg" alt={meta("siteName")} width={84} height={16} />
        </Link>
        <nav aria-label={ui("footerNav")} className="lg:col-[3/10]">
          <ul role="list" className="flex flex-wrap gap-x-6 gap-y-2">
            <li>
              <FooterLink href="/legal">{t("legalNotice")}</FooterLink>
            </li>
            <li>
              <FooterLink href="/privacy">{t("privacy")}</FooterLink>
            </li>
            <li>
              <FooterLink href="/acceptable-use">{t("acceptableUse")}</FooterLink>
            </li>
            <li>
              <ExternalFooterLink href={env.LINKEDIN_URL}>{t("linkedIn")}</ExternalFooterLink>
            </li>
            {env.REPO_URL && (
              <li>
                <ExternalFooterLink href={env.REPO_URL}>{t("sourceCode")}</ExternalFooterLink>
              </li>
            )}
          </ul>
        </nav>
        <p className="text-ink-60 lg:col-[10/13] lg:text-right">{t("copyright")}</p>
        <p className="text-ink-60 lg:col-[3/13]">{t("tagline")}</p>
      </div>
    </footer>
  );
}
