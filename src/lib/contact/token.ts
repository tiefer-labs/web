// SPDX-License-Identifier: MPL-2.0
import "server-only";
import { createHmac, randomBytes, timingSafeEqual } from "node:crypto";

const DOMAIN = "tiefer contact form v1\0";
const PAYLOAD = 8 + 16; // issue time (unix seconds) + random nonce
const LENGTH = PAYLOAD + 32;

export const MIN_AGE_MS = 3_000;
export const MAX_AGE_MS = 2 * 60 * 60 * 1000;
const MAX_SKEW_MS = 60_000;

export type TokenCheck = { ok: true; nonce: string } | { ok: false; reason: "invalid" | "expired" | "too-fast" };

function sign(secret: string, payload: Buffer) {
  return createHmac("sha256", secret).update(DOMAIN).update(payload).digest();
}

export function issueToken(secret: string, now = Date.now()): string {
  const bytes = Buffer.alloc(LENGTH);
  bytes.writeBigUInt64BE(BigInt(Math.floor(now / 1000)), 0);
  randomBytes(16).copy(bytes, 8);
  sign(secret, bytes.subarray(0, PAYLOAD)).copy(bytes, PAYLOAD);
  return bytes.toString("base64url");
}

export function checkToken(secret: string, token: string, now = Date.now()): TokenCheck {
  if (token.length !== Math.ceil((LENGTH * 4) / 3) || !/^[\w-]+$/.test(token)) return { ok: false, reason: "invalid" };
  const bytes = Buffer.from(token, "base64url");
  if (bytes.length !== LENGTH) return { ok: false, reason: "invalid" };
  if (!timingSafeEqual(bytes.subarray(PAYLOAD), sign(secret, bytes.subarray(0, PAYLOAD)))) {
    return { ok: false, reason: "invalid" };
  }
  const age = now - Number(bytes.readBigUInt64BE(0)) * 1000;
  if (age < -MAX_SKEW_MS || age > MAX_AGE_MS) return { ok: false, reason: "expired" };
  if (age < MIN_AGE_MS) return { ok: false, reason: "too-fast" };
  return { ok: true, nonce: bytes.subarray(8, PAYLOAD).toString("hex") };
}
