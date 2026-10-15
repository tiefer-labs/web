// SPDX-License-Identifier: MPL-2.0
import { expect, test } from "vitest";
import { buildMetadata } from "@/lib/seo";

test("the home page uses the site title and description", async () => {
  expect(await buildMetadata({ locale: "en", path: "/" })).toMatchObject({
    title: { absolute: "Tiefer | Edge AI in orbit" },
    alternates: { canonical: "/" },
    openGraph: { title: "Tiefer | Edge AI in orbit", url: "/", locale: "en" },
    twitter: { card: "summary_large_image" },
  });
});

test("other pages put their own title before the site name", async () => {
  const metadata = await buildMetadata({ locale: "en", path: "/privacy", title: "Privacy notice", description: "D" });
  expect(metadata).toMatchObject({
    title: { absolute: "Privacy notice | Tiefer" },
    description: "D",
    alternates: { canonical: "/privacy", languages: {} },
    openGraph: { title: "Privacy notice | Tiefer", description: "D", url: "/privacy" },
  });
});
