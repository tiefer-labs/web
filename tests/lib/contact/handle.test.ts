// SPDX-License-Identifier: MPL-2.0
import { beforeEach, describe, expect, test, vi } from "vitest";
import { handleContact, type ContactDeps } from "@/lib/contact/handle";
import { MemoryStore } from "@/lib/contact/store";
import { checkToken, issueToken } from "@/lib/contact/token";

const secret = "0123456789abcdef0123456789abcdef";
const issuedAt = Date.UTC(2026, 9, 1, 9);
const later = issuedAt + 10_000;

function submission(overrides: Record<string, string> = {}, token = issueToken(secret, issuedAt)) {
  const data = new FormData();
  const fields = { t: token, name: "Ada", email: "ada@example.org", message: "Fires", consent: "yes", ...overrides };
  for (const [name, value] of Object.entries(fields)) data.set(name, value);
  return data;
}

let deps: ContactDeps & { send: ReturnType<typeof vi.fn<ContactDeps["send"]>> };

beforeEach(() => {
  deps = {
    secret,
    store: new MemoryStore(),
    clientKey: "203.0.113.7",
    send: vi.fn<ContactDeps["send"]>(async () => {}),
    now: later,
  };
});

describe("handleContact", () => {
  test("sends a valid message, clears the values and issues a fresh token", async () => {
    const state = await handleContact(submission(), deps);
    expect(state).toMatchObject({ status: "sent", errors: {}, values: { name: "" } });
    expect(checkToken(secret, state.token, later + 5000).ok).toBe(true);
    expect(deps.send).toHaveBeenCalledWith(expect.objectContaining({ name: "Ada", email: "ada@example.org" }));
  });

  test("answers bots as if sent but sends nothing", async () => {
    expect((await handleContact(submission({ website: "spam" }), deps)).status).toBe("sent");
    expect((await handleContact(submission(), { ...deps, now: issuedAt + 1000 })).status).toBe("sent");
    expect(deps.send).not.toHaveBeenCalled();
  });

  test("keeps the values and reports field errors without spending the token", async () => {
    const token = issueToken(secret, issuedAt);
    const invalid = await handleContact(submission({ email: "a@b" }, token), deps);
    expect(invalid).toMatchObject({ status: "invalid", errors: { email: "emailInvalid" }, values: { name: "Ada" } });
    expect((await handleContact(submission({}, token), deps)).status).toBe("sent");
  });

  test("rejects a missing, forged or replayed token as expired", async () => {
    const token = issueToken(secret, issuedAt);
    expect((await handleContact(submission({}, ""), deps)).status).toBe("expired");
    expect((await handleContact(submission({}, issueToken("x".repeat(32), issuedAt)), deps)).status).toBe("expired");
    await handleContact(submission({}, token), deps);
    const replay = await handleContact(submission({}, token), deps);
    expect(replay).toMatchObject({ status: "expired", values: { name: "Ada" } });
    expect(deps.send).toHaveBeenCalledOnce();
  });

  test("limits attempts and sends per client", async () => {
    const statuses = [];
    for (let i = 0; i < 6; i++) statuses.push((await handleContact(submission(), deps)).status);
    expect(statuses).toEqual(["sent", "sent", "sent", "sent", "sent", "limited"]);

    for (let i = 0; i < 14; i++) await handleContact(submission({ name: "" }), deps);
    expect(await handleContact(submission(), deps)).toMatchObject({ status: "limited", values: { name: "" } });
  });

  test("reports a delivery failure without logging the message", async () => {
    const log = vi.spyOn(console, "error").mockImplementation(() => {});
    deps.send.mockRejectedValueOnce(new Error("connection refused"));
    expect(await handleContact(submission({ message: "secret plans" }), deps)).toMatchObject({
      status: "error",
      values: { message: "secret plans" },
    });
    expect(JSON.stringify(log.mock.calls)).not.toContain("secret plans");
  });
});
