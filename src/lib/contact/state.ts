// SPDX-License-Identifier: MPL-2.0
import type { FieldErrors, Submission } from "./submission";

export type ContactStatus = "idle" | "sent" | "invalid" | "expired" | "limited" | "error";

export type ContactState = {
  status: ContactStatus;
  errors: FieldErrors;
  values: Submission;
  token: string;
};

export const emptyValues: Submission = { name: "", email: "", org: "", role: "", message: "", consent: false };
