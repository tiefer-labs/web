// SPDX-License-Identifier: MPL-2.0
import { expect, test } from "vitest";
import { buildMetadata } from "@/lib/seo";

test("the home page uses the site title and description", async () => {
  expect(await buildMetadata({ locale: "en", path: "/" })).toMatchObject({
    title: { absolute: "Tiefer | Edge AI in orbit" },
    description: expect.stringContaining("Earth observation satellites"),
    alternates: { canonical: "/" },
  });
});

test("other pages put their own title before the site name", async () => {
  const metadata = await buildMetadata({ locale: "en", path: "/privacy", title: "Privacy notice", description: "D" });
  expect(metadata).toMatchObject({
    title: { absolute: "Privacy notice | Tiefer" },
    description: "D",
    alternates: { canonical: "/privacy", languages: {} },
  });
});
