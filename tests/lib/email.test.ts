// SPDX-License-Identifier: MPL-2.0
import { expect, test } from "vitest";
import { isEmail } from "@/lib/email";

test.each(["hello@tiefer.space", "first.last+tag@example.co.uk"])("accepts %s", (address) => {
  expect(isEmail(address)).toBe(true);
});

test.each([
  "",
  "a",
  "a@b",
  "a@b.c\n",
  "a b@c.de",
  "<a@b.de>",
  "a@b.de,c@d.de",
  '"x"@b.de',
  "Ada <a@b.de>",
  "a..b@c.de",
  "a@b..de",
  `${"a".repeat(250)}@b.de`,
])("rejects %j", (address) => {
  expect(isEmail(address)).toBe(false);
});
