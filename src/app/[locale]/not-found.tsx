// SPDX-License-Identifier: MPL-2.0
import type { Metadata } from "next";
import { SiteHeader } from "@/components/layout/site-header";
import { NotFoundContent, notFoundMetadata } from "@/components/not-found-content";

export function generateMetadata(): Promise<Metadata> {
  return notFoundMetadata();
}

export default function NotFound() {
  return (
    <>
      <SiteHeader />
      <main id="main" tabIndex={-1}>
        <NotFoundContent />
      </main>
    </>
  );
}
