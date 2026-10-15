// SPDX-License-Identifier: MPL-2.0
import { readFileSync } from "node:fs";
import { describe, expect, test } from "vitest";
import { envNames, parseEnv } from "@/lib/env.server";

const secret = "s".repeat(32);
const smtp = { SMTP_HOST: "smtp.example.org", CONTACT_TO: "hello@tiefer.space", CONTACT_FROM: "web@tiefer.space" };

describe("parseEnv", () => {
  test("defaults to a development setup with placeholders and a generated secret", () => {
    const env = parseEnv({});
    expect(env).toMatchObject({
      ENV: "development",
      SITE_URL: "http://localhost:3000",
      CONTACT_EMAIL: "hello@tiefer.space",
      LEGAL_NAME: "[COMPANY_LEGAL_NAME]",
      LEGAL_REVIEWED: false,
      LOG_RETENTION_DAYS: 30,
    });
    expect(env.CSRF_SECRET).toHaveLength(64);
  });

  test("treats blank values as unset", () => {
    expect(parseEnv({ LEGAL_NAME: "", SMTP_HOST: "" }).LEGAL_NAME).toBe("[COMPANY_LEGAL_NAME]");
  });

  test("accepts a complete production setup", () => {
    const env = parseEnv({ ENV: "production", SITE_URL: "https://tiefer.space", CSRF_SECRET: secret, ...smtp });
    expect(env.SMTP_HOST).toBe("smtp.example.org");
  });

  test.each([
    [{ ENV: "staging" }, "ENV"],
    [{ SITE_URL: "https://tiefer.space/path" }, "SITE_URL must contain only a scheme and a host"],
    [{ ENV: "production", SITE_URL: "http://tiefer.space", CSRF_SECRET: secret }, "SITE_URL must use https"],
    [{ ENV: "production", SITE_URL: "https://tiefer.space" }, "CSRF_SECRET is required in production"],
    [{ CSRF_SECRET: "short" }, "CSRF_SECRET must be at least 32 characters"],
    [{ CSRF_SECRET: "@Microsoft.KeyVault(SecretUri=x)" }, "unresolved Key Vault reference"],
    [{ LINKEDIN_URL: "http://linkedin.com" }, "LINKEDIN_URL"],
    [{ CONTACT_TO: "hello@tiefer.space" }, "SMTP_HOST is required"],
    [{ ...smtp, CONTACT_TO: "nope" }, "CONTACT_TO must be a valid address"],
    [{ ...smtp, SMTP_USER: "user" }, "SMTP_USER and SMTP_PASS must be set together"],
    [{ ...smtp, SMTP_HOST: "smtp://example.org" }, "SMTP_HOST must be a host name"],
    [{ LOG_RETENTION_DAYS: "0" }, "LOG_RETENTION_DAYS"],
    [{ UPSTASH_REDIS_REST_URL: "https://x.upstash.io" }, "UPSTASH_REDIS_REST_TOKEN"],
    [
      { ENV: "production", SITE_URL: "https://tiefer.space", CSRF_SECRET: secret, BEHIND_FRONT_DOOR: "true" },
      "FRONT_DOOR_ID",
    ],
  ])("rejects %j", (source, problem) => {
    expect(() => parseEnv(source)).toThrow(problem);
  });
});

test(".env.example documents every setting, with no secret filled in", () => {
  const example = readFileSync(".env.example", "utf8");
  const entries = example.split(/\r?\n/).filter((line) => /^[A-Z_]+=/.test(line));
  expect(entries.map((line) => line.split("=")[0]).toSorted()).toEqual([...envNames].toSorted());
  for (const name of ["SMTP_PASS", "CSRF_SECRET", "FRONT_DOOR_ID", "UPSTASH_REDIS_REST_TOKEN"]) {
    expect(entries).toContain(`${name}=`);
  }
});
