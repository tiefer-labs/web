// SPDX-License-Identifier: MPL-2.0
import { SiteHeader } from "@/components/layout/site-header";

export default function LegalLayout({ children }: LayoutProps<"/[locale]">) {
  return (
    <>
      <SiteHeader />
      <main id="main" tabIndex={-1}>
        {children}
      </main>
    </>
  );
}
