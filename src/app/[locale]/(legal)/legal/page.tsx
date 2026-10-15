// SPDX-License-Identifier: MPL-2.0
import type { Metadata } from "next";
import { getLocale } from "next-intl/server";
import { LegalPage, legalMetadata } from "@/components/legal/legal-page";

export async function generateMetadata(): Promise<Metadata> {
  return legalMetadata(await getLocale(), "notice", "/legal");
}

export default function Page() {
  return <LegalPage page="notice" />;
}
