// SPDX-License-Identifier: MPL-2.0
import { screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { expect, test } from "vitest";
import { SiteNav } from "@/components/layout/site-nav";
import { renderClient } from "../../test-utils/render";

function renderNav() {
  return renderClient(
    <SiteNav
      dark={false}
      label="Main"
      openLabel="Menu"
      closeLabel="Close"
      items={[{ href: "/#product", label: "Product" }]}
      cta={{ href: "/#contact", label: "Talk to us" }}
    />,
  );
}

test("toggles the menu and swaps the button label", async () => {
  renderNav();
  const toggle = screen.getByRole("button", { name: "Menu" });
  expect(toggle).toHaveAttribute("aria-expanded", "false");
  await userEvent.click(toggle);
  expect(toggle).toHaveAttribute("aria-expanded", "true");
  expect(toggle).toHaveTextContent("Close");
});

test("Escape closes the menu and returns focus to the button", async () => {
  renderNav();
  const toggle = screen.getByRole("button", { name: "Menu" });
  await userEvent.click(toggle);
  await userEvent.tab();
  await userEvent.keyboard("{Escape}");
  expect(toggle).toHaveAttribute("aria-expanded", "false");
  expect(toggle).toHaveFocus();
});

test("following a link closes the menu", async () => {
  renderNav();
  await userEvent.click(screen.getByRole("button", { name: "Menu" }));
  await userEvent.click(screen.getByRole("link", { name: "Product" }));
  expect(screen.getByRole("button", { name: "Menu" })).toHaveAttribute("aria-expanded", "false");
});
