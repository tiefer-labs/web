// SPDX-License-Identifier: MPL-2.0
import { screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { expect, test, vi } from "vitest";
import { CopyEmail } from "@/components/home/contact/copy-email";
import { renderClient } from "../../test-utils/render";

const labels = { copy: "Copy email address", copied: "Copied", selected: "Address selected" };

test("copies the address and says so", async () => {
  const writeText = vi.fn<(text: string) => Promise<void>>().mockResolvedValue(undefined);
  Object.defineProperty(navigator, "clipboard", { value: { writeText }, configurable: true });
  renderClient(<CopyEmail email="hello@tiefer.space" addressId="address" labels={labels} />);

  await userEvent.click(screen.getByRole("button", { name: "Copy email address" }));

  expect(writeText).toHaveBeenCalledWith("hello@tiefer.space");
  expect(screen.getByRole("button", { name: "Copied" })).toBeInTheDocument();
});

test("selects the address for the visitor when the browser refuses clipboard access", async () => {
  const writeText = vi.fn<(text: string) => Promise<void>>().mockRejectedValue(new DOMException("", "NotAllowedError"));
  Object.defineProperty(navigator, "clipboard", { value: { writeText }, configurable: true });
  renderClient(
    <>
      <a id="address" href="mailto:hello@tiefer.space">
        hello@tiefer.space
      </a>
      <CopyEmail email="hello@tiefer.space" addressId="address" labels={labels} />
    </>,
  );

  await userEvent.click(screen.getByRole("button", { name: "Copy email address" }));

  expect(window.getSelection()?.toString()).toBe("hello@tiefer.space");
  expect(screen.getByRole("button", { name: "Address selected" })).toBeInTheDocument();
});
