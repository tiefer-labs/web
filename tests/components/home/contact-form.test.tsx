// SPDX-License-Identifier: MPL-2.0
import { screen, waitFor } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { beforeEach, expect, test, vi } from "vitest";
import { ContactForm } from "@/components/home/contact/contact-form";
import { submitContact } from "@/lib/contact/actions";
import { emptyValues, type ContactState } from "@/lib/contact/state";
import { renderClient } from "../../test-utils/render";

vi.mock(import("@/lib/contact/actions"), () => ({ submitContact: vi.fn<typeof submitContact>() }));

const answer = vi.mocked(submitContact);

beforeEach(() => {
  answer.mockReset();
});

function renderForm() {
  return renderClient(
    <ContactForm email="hello@tiefer.space" tokenField={<input type="hidden" name="t" value="first" />} />,
  );
}

async function submit() {
  await userEvent.type(screen.getByLabelText("Name"), "Ada");
  await userEvent.type(screen.getByLabelText("Work email"), "ada@example.org");
  await userEvent.type(screen.getByLabelText("What would you want your satellite to detect first?"), "Fires");
  await userEvent.click(screen.getByRole("checkbox"));
  await userEvent.click(screen.getByRole("button", { name: "Send message" }));
}

test("ties field errors to their inputs and focuses the status message", async () => {
  const state: ContactState = {
    status: "invalid",
    errors: { email: "emailInvalid" },
    values: { ...emptyValues, name: "Ada", email: "a@b" },
    token: "second",
  };
  answer.mockResolvedValue(state);
  renderForm();
  await submit();

  const status = await screen.findByText("Please check the fields marked below.");
  await waitFor(() => expect(status.parentElement).toHaveFocus());
  const email = screen.getByLabelText("Work email");
  expect(email).toHaveAttribute("aria-invalid", "true");
  expect(email).toHaveAccessibleDescription("Error: Please enter a valid email address, for example name@company.com.");
  expect(email).toHaveValue("a@b");
  expect(screen.getByLabelText("Name")).toHaveValue("Ada");
  expect(screen.getByLabelText("What would you want your satellite to detect first?")).toHaveAccessibleDescription(
    "Up to 2,000 characters.",
  );
  expect(document.querySelector<HTMLInputElement>('input[name="t"]')!.value).toBe("second");
});

test("sends the token, clears the form on success and thanks the sender", async () => {
  answer.mockResolvedValue({ status: "sent", errors: {}, values: emptyValues, token: "next" });
  renderForm();
  await submit();

  expect(await screen.findByText("Thank you. We will reply within two working days.")).toBeInTheDocument();
  expect(answer.mock.calls[0]?.[1].get("t")).toBe("first");
  expect(screen.getByLabelText("Name")).toHaveValue("");
});

test("offers the email address when sending fails", async () => {
  answer.mockResolvedValue({ status: "error", errors: {}, values: emptyValues, token: "next" });
  renderForm();
  await submit();

  expect(await screen.findByRole("link", { name: "hello@tiefer.space" })).toHaveAttribute(
    "href",
    "mailto:hello@tiefer.space",
  );
});
