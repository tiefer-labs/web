// SPDX-License-Identifier: MPL-2.0
import { createTransport } from "nodemailer";
import { describe, expect, test } from "vitest";
import { HeaderInjectionError, sendMessage, smtpTransport, type Message } from "@/lib/contact/mail";

const message: Message = {
  from: "web@tiefer.space",
  to: "hello@tiefer.space",
  replyTo: "ada@example.org",
  subject: "Website contact: Ada Łukasiewicz",
  body: "Hello\nsecond line with = sign, from Bakı",
  domain: "tiefer.space",
};

async function raw(input: Message) {
  const transport = createTransport({ streamTransport: true, buffer: true, newline: "windows" });
  const info = await sendMessage(transport, input);
  return String((info as unknown as { message: Buffer }).message);
}

describe("sendMessage", () => {
  test("writes plain text mail with Reply-To, a Message-ID and Auto-Submitted", async () => {
    const mail = await raw(message);
    expect(mail).toMatch(/^From: web@tiefer\.space$/m);
    expect(mail).toMatch(/^To: hello@tiefer\.space$/m);
    expect(mail).toMatch(/^Reply-To: ada@example\.org$/m);
    expect(mail).toMatch(/^Subject: =\?UTF-8\?Q\?/m);
    expect(mail).toMatch(/^Message-ID: <[0-9a-f]{32}@tiefer\.space>$/m);
    expect(mail).toMatch(/^Auto-Submitted: auto-generated$/m);
    expect(mail).toMatch(/^Content-Type: text\/plain; charset=utf-8$/m);
    expect(mail).toMatch(/^Content-Transfer-Encoding: quoted-printable$/m);
    expect(mail).toContain("second line with =3D sign");
  });

  test.each(["a\r\nBcc: x@y.z", "a\nb", "a\0"])("refuses a header value with a line break (%j)", async (subject) => {
    await expect(raw({ ...message, subject })).rejects.toBeInstanceOf(HeaderInjectionError);
  });
});

describe("smtpTransport", () => {
  test("uses implicit TLS on 465 and requires STARTTLS elsewhere", () => {
    const config = { host: "smtp.example.org", user: "u", pass: "p", skipVerify: false };
    const implicit = smtpTransport({ ...config, port: 465 }, "tiefer.space").options as Record<string, unknown>;
    const starttls = smtpTransport({ ...config, port: 587 }, "tiefer.space").options as Record<string, unknown>;
    expect(implicit).toMatchObject({ secure: true, requireTLS: false, name: "tiefer.space" });
    expect(starttls).toMatchObject({
      secure: false,
      requireTLS: true,
      tls: { minVersion: "TLSv1.2", rejectUnauthorized: true },
    });
  });
});
