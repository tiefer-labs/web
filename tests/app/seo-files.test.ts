// SPDX-License-Identifier: MPL-2.0
import { describe, expect, test, vi } from "vitest";

/** Config is read when a module loads, so each case reloads the modules under its own environment. */
function withEnv(env: Record<string, string> = {}) {
  vi.resetModules();
  for (const [key, value] of Object.entries(env)) vi.stubEnv(key, value);
}

const production = { ENV: "production", CSRF_SECRET: "s".repeat(32) };

describe("robots", () => {
  test("asks crawlers to stay away outside production", async () => {
    withEnv();
    const { default: robots } = await import("@/app/robots");
    expect(robots()).toEqual({ rules: { userAgent: "*", disallow: "/" } });
  });

  test("allows everything in production and points to the sitemap", async () => {
    withEnv(production);
    const { default: robots } = await import("@/app/robots");
    expect(robots()).toEqual({ rules: { userAgent: "*", allow: "/" }, sitemap: "https://tiefer.space/sitemap.xml" });
  });
});

test("sitemap lists the home page and the three legal pages without a locale prefix", async () => {
  withEnv();
  const { default: sitemap } = await import("@/app/sitemap");
  expect(sitemap().map((entry) => entry.url)).toEqual([
    "https://tiefer.space/",
    "https://tiefer.space/legal",
    "https://tiefer.space/privacy",
    "https://tiefer.space/acceptable-use",
  ]);
});

test("manifest uses the brand colours and icons", async () => {
  withEnv();
  const { default: manifest } = await import("@/app/manifest");
  expect(await manifest()).toMatchObject({
    name: "Tiefer",
    theme_color: "#0C003D",
    background_color: "#FFFFFF",
    icons: [{ src: "/icon-192.png" }, { src: "/icon-512.png" }, { src: "/icon.svg" }],
  });
});

describe("security.txt", () => {
  test("names the contact and links the policy when the repository is known", async () => {
    withEnv({ REPO_URL: "https://github.com/tiefer-labs/web" });
    const { securityTxt } = await import("@/lib/security-txt");
    expect(securityTxt()).toBe(
      [
        "Contact: mailto:hello@tiefer.space",
        "Expires: 2027-09-01T00:00:00Z",
        "Preferred-Languages: en",
        "Canonical: https://tiefer.space/.well-known/security.txt",
        "Policy: https://github.com/tiefer-labs/web/blob/main/SECURITY.md",
        "",
      ].join("\n"),
    );
  });

  test("expires within a year and not within the next month", async () => {
    withEnv();
    const { SECURITY_TXT_EXPIRES } = await import("@/lib/security-txt");
    const daysLeft = (Date.parse(SECURITY_TXT_EXPIRES) - Date.now()) / 86_400_000;
    expect(daysLeft, "renew Expires in src/lib/security-txt.ts").toBeGreaterThan(30);
    expect(daysLeft).toBeLessThanOrEqual(366);
  });
});
