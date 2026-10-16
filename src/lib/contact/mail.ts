// SPDX-License-Identifier: MPL-2.0
import "server-only";
import { randomBytes } from "node:crypto";
import { createTransport, type Transporter } from "nodemailer";

export type Message = {
  from: string;
  to: string;
  replyTo: string;
  subject: string;
  body: string;
  /** Host name used in the Message-ID and in EHLO. */
  domain: string;
};

export type SmtpConfig = {
  host: string;
  port: number;
  user: string;
  pass: string;
  skipVerify: boolean;
};

export class HeaderInjectionError extends Error {
  constructor() {
    super("contact: header value contains a line break");
  }
}

export function smtpTransport(config: SmtpConfig, domain: string): Transporter {
  return createTransport({
    host: config.host,
    port: config.port,
    // 465 is implicit TLS; every other port must upgrade with STARTTLS or the send fails.
    secure: config.port === 465,
    requireTLS: config.port !== 465,
    name: domain,
    auth: config.user ? { user: config.user, pass: config.pass } : undefined,
    tls: { minVersion: "TLSv1.2", rejectUnauthorized: !config.skipVerify, servername: config.host },
    connectionTimeout: 20_000,
    greetingTimeout: 20_000,
    socketTimeout: 20_000,
  });
}

export async function sendMessage(transport: Transporter, message: Message) {
  const headers = [message.from, message.to, message.replyTo, message.subject, message.domain];
  if (headers.some((value) => /[\r\n\0]/.test(value))) throw new HeaderInjectionError();

  return transport.sendMail({
    from: message.from,
    to: message.to,
    replyTo: message.replyTo,
    subject: message.subject,
    text: message.body,
    date: new Date(),
    messageId: `<${randomBytes(16).toString("hex")}@${message.domain}>`,
    textEncoding: "quoted-printable",
    headers: { "Auto-Submitted": "auto-generated" },
  });
}
