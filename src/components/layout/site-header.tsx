// SPDX-License-Identifier: MPL-2.0
import { getTranslations } from "next-intl/server";
import Image from "next/image";
import { Link } from "@/i18n/navigation";
import { cn } from "@/lib/cn";
import { SiteNav } from "./site-nav";

const sections = [
  ["how", "/#how"],
  ["useCases", "/#use-cases"],
  ["principles", "/#principles"],
  ["roadmap", "/#roadmap"],
  ["contact", "/#contact"],
] as const;

/** `dark` sits on top of the home page hero; every other page uses the light header. */
export async function SiteHeader({ dark = false }: { dark?: boolean }) {
  const [t, nav, meta] = await Promise.all([getTranslations("UI"), getTranslations("Nav"), getTranslations("Meta")]);

  return (
    <header
      id="top"
      data-tone={dark ? "dark" : undefined}
      data-print="hide"
      className={cn("z-5", dark ? "absolute inset-x-0 top-0 text-on-ink" : "relative border-b border-line")}
    >
      <div className="wrap flex min-h-18 items-center justify-between gap-6">
        <Link href="/#top" className="block py-2" data-forced-colors="keep">
          <Image
            src={dark ? "/brand/tiefer-logo-white.svg" : "/brand/tiefer-logo.svg"}
            alt={meta("siteName")}
            width={126}
            height={24}
            priority
          />
        </Link>
        <SiteNav
          dark={dark}
          label={t("navLabel")}
          openLabel={t("menuOpen")}
          closeLabel={t("menuClose")}
          items={sections.map(([key, href]) => ({ href, label: nav(key) }))}
          cta={{ href: "/#contact", label: nav("cta") }}
        />
      </div>
    </header>
  );
}
