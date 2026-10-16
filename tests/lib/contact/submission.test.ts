// SPDX-License-Identifier: MPL-2.0
import { describe, expect, test } from "vitest";
import { parseSubmission, validateSubmission, type Submission } from "@/lib/contact/submission";

function form(fields: Record<string, string>) {
  const data = new FormData();
  for (const [name, value] of Object.entries(fields)) data.set(name, value);
  return data;
}

const valid: Submission = {
  name: "Ada",
  email: "ada@example.org",
  org: "",
  role: "",
  message: "Fires",
  consent: true,
};

describe("parseSubmission", () => {
  test("flattens single-line fields and cleans the message", () => {
    const parsed = parseSubmission(
      form({
        name: "  Ada\r\nBcc: x@example.org  ",
        email: " ada@example.org ",
        org: "Orbit\tLab",
        role: "operator",
        message: `Line one\r\nLine two${String.fromCodePoint(0, 7, 0x2028)}`,
        consent: "yes",
      }),
    );
    expect(parsed).toEqual({
      name: "Ada Bcc: x@example.org",
      email: "ada@example.org",
      org: "Orbit Lab",
      role: "operator",
      message: "Line one\nLine two",
      consent: true,
    });
  });

  test("only the exact value yes counts as consent", () => {
    expect(parseSubmission(form({ consent: "on" })).consent).toBe(false);
  });
});

describe("validateSubmission", () => {
  test("accepts a complete submission", () => {
    expect(validateSubmission(valid)).toEqual({});
    expect(validateSubmission({ ...valid, role: "hosted-compute", org: "Orbit Lab" })).toEqual({});
  });

  test("reports each field with its own message key", () => {
    expect(
      validateSubmission({
        name: "a".repeat(101),
        email: "ada@example.org\nBcc: x@y.z",
        org: "o".repeat(151),
        role: "admin",
        message: "",
        consent: false,
      }),
    ).toEqual({
      name: "nameTooLong",
      email: "emailInvalid",
      org: "orgTooLong",
      role: "roleInvalid",
      message: "messageMissing",
      consent: "consentMissing",
    });
    expect(validateSubmission({ ...valid, name: "", email: "" })).toEqual({
      name: "nameMissing",
      email: "emailMissing",
    });
    expect(validateSubmission({ ...valid, message: "m".repeat(2001) })).toEqual({ message: "messageTooLong" });
  });

  test("counts characters, not bytes", () => {
    expect(validateSubmission({ ...valid, message: "é".repeat(2000) })).toEqual({});
  });
});
