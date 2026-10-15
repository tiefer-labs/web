// SPDX-License-Identifier: MPL-2.0
import { describe, expect, test } from "vitest";
import { checkToken, issueToken } from "@/lib/contact/token";

const secret = "0123456789abcdef0123456789abcdef";
const start = Date.UTC(2026, 9, 1, 9);
const seconds = (n: number) => n * 1000;

describe("form token", () => {
  test("is too fast at first, valid after a few seconds, expired after two hours", () => {
    const token = issueToken(secret, start);
    expect(token).toHaveLength(75);
    expect(checkToken(secret, token, start)).toEqual({ ok: false, reason: "too-fast" });
    expect(checkToken(secret, token, start + seconds(5))).toMatchObject({ ok: true });
    expect(checkToken(secret, token, start + seconds(3 * 3600))).toEqual({ ok: false, reason: "expired" });
  });

  test("rejects a token issued more than a minute in the future", () => {
    expect(checkToken(secret, issueToken(secret, start + seconds(120)), start)).toEqual({
      ok: false,
      reason: "expired",
    });
  });

  test("gives each token its own nonce", () => {
    const a = checkToken(secret, issueToken(secret, start), start + seconds(5));
    const b = checkToken(secret, issueToken(secret, start), start + seconds(5));
    expect(a.ok && b.ok && a.nonce !== b.nonce).toBe(true);
  });

  test("rejects tampering, other secrets and garbage", () => {
    const token = issueToken(secret, start);
    const flipped = token.slice(0, 37) + (token[37] === "A" ? "B" : "A") + token.slice(38);
    const other = issueToken("another-secret-another-secret-xx", start);
    for (const bad of [flipped, "", "abc", "A".repeat(75), other, `${token.slice(0, 74)}!`]) {
      expect(checkToken(secret, bad, start + seconds(5))).toEqual({ ok: false, reason: "invalid" });
    }
  });
});
