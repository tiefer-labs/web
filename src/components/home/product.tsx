// SPDX-License-Identifier: MPL-2.0
import { getMessages, getTranslations } from "next-intl/server";
import { Section, SectionLabel } from "@/components/ui/section";
import { cn } from "@/lib/cn";
import { keysOf } from "@/lib/messages";

const cell = "py-[22px] text-left align-baseline font-normal max-sm:p-0";

export async function Product() {
  const [t, messages] = await Promise.all([getTranslations("Product"), getMessages()]);

  return (
    <Section id="product" labelledBy="product-title" tone="ink">
      <SectionLabel className="lg:row-1">{t("label")}</SectionLabel>
      <div className="lg:col-[1/6] lg:row-2">
        <h2 id="product-title" className="text-on-ink">
          {t("title")}
        </h2>
        <p className="mt-6 max-w-[30em]">{t("body")}</p>
      </div>
      <table className="mt-12 w-full border-b border-line-on-ink lg:col-[7/13] lg:row-2 lg:mt-2">
        <caption className="sr-only">{t("title")}</caption>
        <tbody>
          {keysOf(messages.Product.stages).map((key) => (
            <tr key={key} className="border-t border-line-on-ink max-sm:grid max-sm:grid-cols-[3.5rem_1fr] max-sm:py-5">
              <td
                className={cn(cell, "w-14 text-label tracking-[0.12em] text-on-ink-muted tabular-nums max-sm:w-auto")}
              >
                {t(`stages.${key}.number`)}
              </td>
              <th scope="row" className={cn(cell, "w-28 label-caps text-on-ink max-sm:w-auto")}>
                {t(`stages.${key}.name`)}
              </th>
              <td className={cn(cell, "text-on-ink-muted max-sm:col-start-2 max-sm:mt-2")}>
                {t(`stages.${key}.text`)}
              </td>
            </tr>
          ))}
        </tbody>
      </table>
      <p className="mt-8 text-small text-on-ink lg:col-[7/13] lg:row-3">{t("note")}</p>
    </Section>
  );
}
