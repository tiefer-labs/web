// SPDX-License-Identifier: MPL-2.0
import "server-only";
import { randomBytes } from "node:crypto";
import * as z from "zod";
import { isEmail } from "@/lib/email";

const httpsUrl = z.url({ protocol: /^https$/, error: "must be an absolute https URL" });
const bool = z
  .enum(["true", "false"])
  .default("false")
  .transform((v) => v === "true");
const optional = z.string().trim().default("");

const variables = z.object({
  ENV: z.enum(["development", "production"]).default("development"),
  SITE_URL: z.url({ protocol: /^https?$/ }).default("http://localhost:3000"),
  CONTACT_EMAIL: z.string().trim().default("hello@tiefer.space").refine(isEmail, "is not a valid address"),
  LINKEDIN_URL: httpsUrl.default("https://www.linkedin.com/company/tiefer/"),
  REPO_URL: z.union([z.literal(""), httpsUrl]).default(""),

  SMTP_HOST: optional,
  SMTP_PORT: z.coerce.number().int().min(1).max(65535).default(587),
  SMTP_USER: optional,
  SMTP_PASS: z.string().default(""),
  SMTP_SKIP_VERIFY: bool,
  CONTACT_TO: optional,
  CONTACT_FROM: optional,
  CSRF_SECRET: z.string().default(""),

  REDIS_URL: optional,
  UPSTASH_REDIS_REST_URL: optional,
  UPSTASH_REDIS_REST_TOKEN: optional,

  LEGAL_NAME: z.string().default("[COMPANY_LEGAL_NAME]"),
  LEGAL_FORM: z.string().default("[LEGAL_FORM]"),
  LEGAL_ADDRESS: z.string().default("[REGISTERED_ADDRESS_BAKU]"),
  LEGAL_TAX_ID: z.string().default("[VOEN]"),
  LEGAL_REGISTRATION: z.string().default("[STATE_REGISTRATION_DETAILS]"),
  LEGAL_DIRECTOR: z.string().default("[MANAGING_DIRECTOR]"),
  LEGAL_HOSTING_PROVIDER: z.string().default("[HOSTING_PROVIDER]"),
  LEGAL_HOSTING_COUNTRY: z.string().default("[HOSTING_COUNTRY]"),
  LEGAL_SMTP_PROVIDER: z.string().default("[SMTP_PROVIDER]"),
  LEGAL_SMTP_COUNTRY: z.string().default("[SMTP_COUNTRY]"),
  LEGAL_REVIEWED: bool,
  LOG_RETENTION_DAYS: z.coerce.number().int().min(1).max(3650).default(30),

  BEHIND_FRONT_DOOR: bool,
  FRONT_DOOR_ID: optional,
  // Read by next.config.ts for the HSTS header; listed here so it is validated too.
  HSTS_PRELOAD: bool,
});

/** Every setting the site reads; `.env.example` documents each one. */
export const envNames = Object.keys(variables.shape);

const schema = variables.superRefine((env, ctx) => {
  const fail = (path: string, message: string) => ctx.addIssue({ code: "custom", path: [path], message });
  const production = env.ENV === "production";
  const site = new URL(env.SITE_URL);

  if (site.pathname !== "/" || site.search || site.hash) fail("SITE_URL", "must contain only a scheme and a host");
  if (production && site.protocol !== "https:") fail("SITE_URL", "must use https in production");

  const smtpRelated = [env.SMTP_USER, env.SMTP_PASS, env.CONTACT_TO, env.CONTACT_FROM].some(Boolean);
  if (env.SMTP_HOST) {
    if (/[/:]/.test(env.SMTP_HOST)) fail("SMTP_HOST", "must be a host name without scheme or port");
    if (!isEmail(env.CONTACT_TO)) fail("CONTACT_TO", "must be a valid address when SMTP_HOST is set");
    if (!isEmail(env.CONTACT_FROM)) fail("CONTACT_FROM", "must be a valid address when SMTP_HOST is set");
    if (Boolean(env.SMTP_USER) !== Boolean(env.SMTP_PASS))
      fail("SMTP_USER", "SMTP_USER and SMTP_PASS must be set together");
  } else if (smtpRelated) {
    fail("SMTP_HOST", "is required when other SMTP or CONTACT_ settings are set");
  }
  if (production && env.SMTP_SKIP_VERIFY) fail("SMTP_SKIP_VERIFY", "is not allowed in production");

  for (const name of ["CSRF_SECRET", "SMTP_PASS"] as const) {
    if (env[name].startsWith("@Microsoft.KeyVault(")) fail(name, "is an unresolved Key Vault reference");
  }
  if (env.CSRF_SECRET && env.CSRF_SECRET.length < 32) fail("CSRF_SECRET", "must be at least 32 characters");
  // Pages are prerendered at build time without secrets; the running server must have the secret.
  const building = process.env.NEXT_PHASE === "phase-production-build";
  if (production && !building && !env.CSRF_SECRET) fail("CSRF_SECRET", "is required in production");
  if (env.UPSTASH_REDIS_REST_URL && !env.UPSTASH_REDIS_REST_TOKEN) {
    fail("UPSTASH_REDIS_REST_TOKEN", "is required with UPSTASH_REDIS_REST_URL");
  }
  if (production && env.BEHIND_FRONT_DOOR && !env.FRONT_DOOR_ID) {
    fail("FRONT_DOOR_ID", "is required when BEHIND_FRONT_DOOR=true in production");
  }
});

export type Env = z.infer<typeof schema>;

export function parseEnv(source: Record<string, string | undefined>): Env {
  const blankAsUnset = Object.fromEntries(Object.entries(source).filter(([, value]) => value !== ""));
  const result = schema.safeParse(blankAsUnset);
  if (!result.success) {
    const problems = result.error.issues.map((issue) => `${issue.path.join(".")} ${issue.message}`);
    throw new Error(`Invalid configuration:\n- ${problems.join("\n- ")}`);
  }
  // Development only: a fresh secret per process is enough when nothing needs to verify tokens across restarts.
  return { ...result.data, CSRF_SECRET: result.data.CSRF_SECRET || randomBytes(32).toString("hex") };
}

export const env = parseEnv(process.env);

export const siteUrl = new URL(env.SITE_URL).origin;
export const contactEnabled = env.SMTP_HOST !== "";
