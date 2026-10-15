// SPDX-License-Identifier: MPL-2.0
import "server-only";
import { emptyValues, type ContactState, type ContactStatus } from "./state";
import type { Store } from "./store";
import { parseSubmission, validateSubmission, type FieldErrors, type Submission } from "./submission";
import { checkToken, issueToken, MAX_AGE_MS } from "./token";

const HOUR = 60 * 60 * 1000;
export const limits = { attempts: 20, sends: 5, globalSends: 60 } as const;

export type ContactDeps = {
  secret: string;
  store: Store;
  clientKey: string;
  send: (submission: Submission) => Promise<void>;
  now?: number;
};

/** Same order as the Go handler: the cheap checks first, the token is only spent on a valid submission. */
export async function handleContact(form: FormData, deps: ContactDeps): Promise<ContactState> {
  const { secret, store, clientKey, now = Date.now() } = deps;
  const answer = (status: ContactStatus, values = emptyValues, errors: FieldErrors = {}): ContactState => ({
    status,
    errors,
    values,
    token: issueToken(secret, now),
  });

  if (!(await store.hit(`attempts:${clientKey}`, limits.attempts, HOUR))) return answer("limited");

  const submission = parseSubmission(form);
  // A filled honeypot or a form sent faster than a person can type gets the success answer; nothing is sent.
  if (form.get("website")) return answer("sent");

  const field = form.get("t");
  const token = typeof field === "string" ? field : "";
  const check = checkToken(secret, token, now);
  if (!check.ok) return check.reason === "too-fast" ? answer("sent") : answer("expired", submission);

  const errors = validateSubmission(submission);
  if (Object.keys(errors).length > 0) return answer("invalid", submission, errors);

  if (!(await store.claim(`token:${check.nonce}`, MAX_AGE_MS + 60_000))) return answer("expired", submission);

  const allowed =
    (await store.hit(`sends:${clientKey}`, limits.sends, HOUR)) &&
    (await store.hit("sends:all", limits.globalSends, HOUR));
  if (!allowed) return answer("limited", submission);

  try {
    await deps.send(submission);
  } catch (error) {
    // Only the error itself: never the message, the sender or credentials.
    console.error("contact delivery failed", error instanceof Error ? error.message : error);
    return answer("error", submission);
  }
  return answer("sent");
}
