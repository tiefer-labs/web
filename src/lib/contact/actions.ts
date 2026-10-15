// SPDX-License-Identifier: MPL-2.0
"use server";

import { getTranslations } from "next-intl/server";
import { headers } from "next/headers";
import { routing } from "@/i18n/routing";
import { env, siteUrl } from "@/lib/env.server";
import { clientAddress, clientKey } from "./client-ip";
import { handleContact } from "./handle";
import { sendMessage, smtpTransport } from "./mail";
import { isRole } from "./roles";
import type { ContactState } from "./state";
import { createStore } from "./store";
import type { Submission } from "./submission";

const store = createStore(env);
const domain = new URL(siteUrl).hostname;
let transport: ReturnType<typeof smtpTransport> | undefined;

async function deliver(submission: Submission) {
  // The mail goes to the team, so it always uses the default locale.
  const t = await getTranslations({ locale: routing.defaultLocale, namespace: "Contact" });
  const role = isRole(submission.role) ? t(`roles.${submission.role}`) : "-";
  const body = [
    `Name: ${submission.name}`,
    `Email: ${submission.email}`,
    `Organisation: ${submission.org || "-"}`,
    `Role: ${role}`,
    "Consent to reply: yes",
    "",
    t("messageLabel"),
    "",
    submission.message,
    "",
  ].join("\n");

  transport ??= smtpTransport(
    {
      host: env.SMTP_HOST,
      port: env.SMTP_PORT,
      user: env.SMTP_USER,
      pass: env.SMTP_PASS,
      skipVerify: env.SMTP_SKIP_VERIFY,
    },
    domain,
  );
  await sendMessage(transport, {
    from: env.CONTACT_FROM,
    to: env.CONTACT_TO,
    replyTo: submission.email,
    subject: `Website contact: ${submission.name}`,
    body,
    domain,
  });
}

export async function submitContact(_previous: ContactState, form: FormData): Promise<ContactState> {
  const address = clientAddress(await headers(), env.BEHIND_FRONT_DOOR);
  return handleContact(form, { secret: env.CSRF_SECRET, store, clientKey: clientKey(address), send: deliver });
}
