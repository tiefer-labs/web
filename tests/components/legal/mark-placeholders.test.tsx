// SPDX-License-Identifier: MPL-2.0
import { render } from "@testing-library/react";
import { expect, test } from "vitest";
import { markPlaceholders } from "@/components/legal/mark-placeholders";

test("marks placeholders in text and leaves elements alone", () => {
  const { container } = render(
    <p>
      {markPlaceholders([
        "Run by [COMPANY_LEGAL_NAME] in ",
        <a key="x" href="https://example.org">
          [NOT_MARKED]
        </a>,
        ". [LEGAL_BASIS: ask counsel]",
      ])}
    </p>,
  );
  expect([...container.querySelectorAll("mark")].map((m) => m.textContent)).toEqual([
    "[COMPANY_LEGAL_NAME]",
    "[LEGAL_BASIS: ask counsel]",
  ]);
  expect(container).toHaveTextContent("Run by [COMPANY_LEGAL_NAME] in [NOT_MARKED]. [LEGAL_BASIS: ask counsel]");
});

test("ignores ordinary brackets", () => {
  const { container } = render(<p>{markPlaceholders("see [a] or [AB]")}</p>);
  expect(container.querySelector("mark")).toBeNull();
});
