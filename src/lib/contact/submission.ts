// SPDX-License-Identifier: MPL-2.0
import { isEmail } from "@/lib/email";
import { isRole } from "./roles";

export const limits = { name: 100, email: 254, org: 150, message: 2000 } as const;

export type Submission = {
  name: string;
  email: string;
  org: string;
  role: string;
  message: string;
  consent: boolean;
};

export type ErrorKey =
  | "nameMissing"
  | "nameTooLong"
  | "emailMissing"
  | "emailInvalid"
  | "orgTooLong"
  | "roleInvalid"
  | "messageMissing"
  | "messageTooLong"
  | "consentMissing";

export type FieldErrors = Partial<Record<keyof Submission, ErrorKey>>;

// Control characters plus the Unicode line and paragraph separators.
const CONTROL = /[\p{Cc}\p{Zl}\p{Zp}]/gu;

function singleLine(value: string) {
  return value.replace(CONTROL, " ").split(/\s+/u).filter(Boolean).join(" ");
}

function multiLine(value: string) {
  return value
    .replace(/\r\n?/g, "\n")
    .replace(CONTROL, (c) => (c === "\n" || c === "\t" ? c : ""))
    .trim();
}

const text = (form: FormData, name: string) => {
  const value = form.get(name);
  return typeof value === "string" ? value : "";
};

export function parseSubmission(form: FormData): Submission {
  return {
    name: singleLine(text(form, "name")),
    email: text(form, "email").trim(),
    org: singleLine(text(form, "org")),
    role: text(form, "role").trim(),
    message: multiLine(text(form, "message")),
    consent: text(form, "consent") === "yes",
  };
}

// Limits count code points, not UTF-16 units.
// oxlint-disable-next-line typescript/no-misused-spread
const length = (value: string) => [...value].length;

export function validateSubmission(s: Submission): FieldErrors {
  const errors: FieldErrors = {};
  if (length(s.name) === 0) errors.name = "nameMissing";
  else if (length(s.name) > limits.name) errors.name = "nameTooLong";
  if (s.email === "") errors.email = "emailMissing";
  else if (s.email.length > limits.email || !isEmail(s.email)) errors.email = "emailInvalid";
  if (length(s.org) > limits.org) errors.org = "orgTooLong";
  if (s.role !== "" && !isRole(s.role)) errors.role = "roleInvalid";
  if (length(s.message) === 0) errors.message = "messageMissing";
  else if (length(s.message) > limits.message) errors.message = "messageTooLong";
  if (!s.consent) errors.consent = "consentMissing";
  return errors;
}
