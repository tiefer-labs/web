// SPDX-License-Identifier: MPL-2.0
import { getTranslations } from "next-intl/server";

export async function SkipLink() {
  const t = await getTranslations("UI");
  return (
    <a
      href="#main"
      data-print="hide"
      className="absolute top-2 left-(--pad) z-10 -translate-y-[160%] bg-paper px-4 py-2 font-medium text-ink focus:translate-y-0"
    >
      {t("skipLink")}
    </a>
  );
}
