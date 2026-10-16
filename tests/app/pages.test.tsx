// SPDX-License-Identifier: MPL-2.0
import { screen } from "@testing-library/react";
import type { ReactNode } from "react";
import { describe, expect, test } from "vitest";
import AcceptableUsePage from "@/app/[locale]/(legal)/acceptable-use/page";
import LegalLayout from "@/app/[locale]/(legal)/layout";
import LegalNoticePage from "@/app/[locale]/(legal)/legal/page";
import PrivacyPage from "@/app/[locale]/(legal)/privacy/page";
import NotFound from "@/app/[locale]/not-found";
import HomePage from "@/app/[locale]/page";
import messages from "../../messages/en.json";
import { renderServer } from "../test-utils/render";
import { bannedWords, forbiddenCharacters, messageStrings, visibleText } from "../test-utils/text-rules";

const legal = (page: ReactNode) => <LegalLayout params={Promise.resolve({ locale: "en" })}>{page}</LegalLayout>;

const pages: [string, () => ReactNode][] = [
  ["home", () => <HomePage />],
  ["legal notice", () => legal(<LegalNoticePage />)],
  ["privacy", () => legal(<PrivacyPage />)],
  ["acceptable use", () => legal(<AcceptableUsePage />)],
  ["404", () => <NotFound />],
];

describe.each(pages)("%s page", (_name, page) => {
  test("has one h1 and never skips a heading level", async () => {
    const { container } = await renderServer(page());
    const levels = [...container.querySelectorAll("h1, h2, h3, h4")].map((h) => Number(h.tagName[1]));
    expect(levels.filter((level) => level === 1)).toHaveLength(1);
    const jumps = levels.slice(1).filter((level, i) => level - (levels[i] ?? 0) > 1);
    expect(jumps).toEqual([]);
  });

  test("follows the copy rules in every visible word", async () => {
    const { container } = await renderServer(page());
    const text = visibleText(container);
    expect([...forbiddenCharacters(text), ...bannedWords(text)]).toEqual([]);
  });
});

describe("home page", () => {
  test("shows every piece of home page copy", async () => {
    const { container } = await renderServer(<HomePage />);
    const text = container.textContent.replace(/\s+/g, " ");
    const home = ["Hero", "Problem", "How", "Demo", "UseCases", "Principles", "Roadmap", "Name"] as const;
    const missing = home
      .flatMap((section) => messageStrings(messages[section]))
      // A {placeholder} is filled in on the page, so only the text around it is checked.
      .filter(([, copy]) => !copy.split(/\{\w+\}/).every((part) => text.includes(part.trim())));
    expect(missing).toEqual([]);
  });

  test("marks the current roadmap step and tags observed and inferred evidence", async () => {
    await renderServer(<HomePage />);
    expect(screen.getByRole("listitem", { current: "step" })).toHaveTextContent("Benchmark");
    expect(screen.getByRole("rowheader", { name: "Observed" }).closest("tr")).toHaveClass("bg-tint");
  });

  test("links the Organization structured data to LinkedIn", async () => {
    const { container } = await renderServer(<HomePage />);
    const jsonLd = JSON.parse(container.querySelector('script[type="application/ld+json"]')!.textContent);
    expect(jsonLd).toEqual({
      "@context": "https://schema.org",
      "@type": "Organization",
      logo: "https://tiefer.space/icon-512.png",
      name: "Tiefer",
      sameAs: ["https://www.linkedin.com/company/tiefer/"],
      url: "https://tiefer.space/",
    });
  });

  test("offers email instead of the form when SMTP is not configured", async () => {
    await renderServer(<HomePage />);
    expect(screen.queryByRole("form")).not.toBeInTheDocument();
    expect(screen.getByRole("link", { name: "Email us" })).toHaveAttribute("href", "mailto:hello@tiefer.space");
  });
});

describe("legal pages", () => {
  test("fill company details from config, marking the ones still missing", async () => {
    const { container } = await renderServer(legal(<LegalNoticePage />));
    const marks = [...container.querySelectorAll("mark")].map((mark) => mark.textContent);
    expect(marks).toContain("[COMPANY_LEGAL_NAME]");
    expect(marks).toContain("[LAST_UPDATED]");
    expect(screen.getByRole("link", { name: "hello@tiefer.space" })).toHaveAttribute(
      "href",
      "mailto:hello@tiefer.space",
    );
    expect(screen.getByRole("link", { name: "https://tiefer.space" })).toHaveAttribute("href", "https://tiefer.space/");
    expect(screen.getByRole("note")).toHaveTextContent("Placeholder text. To be reviewed before launch.");
  });

  test("list sections in a table of contents only on long pages", async () => {
    await renderServer(legal(<PrivacyPage />));
    const toc = screen.getByRole("navigation", { name: "On this page" });
    expect(toc.querySelectorAll("a")).toHaveLength(11);
    expect(toc.querySelector("a")).toHaveAttribute("href", "#who-is-responsible");
    document.body.innerHTML = "";
    await renderServer(legal(<LegalNoticePage />));
    expect(screen.queryByRole("navigation", { name: "On this page" })).not.toBeInTheDocument();
  });
});
