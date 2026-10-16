// SPDX-License-Identifier: MPL-2.0
import { readdirSync, readFileSync } from "node:fs";
import { createTranslator, IntlErrorCode } from "next-intl";
import { describe, expect, test } from "vitest";
import { roles } from "@/lib/contact/roles";
import en from "../messages/en.json";
import { bannedWords, forbiddenCharacters, messageStrings } from "./test-utils/text-rules";

const locales = readdirSync("messages").filter((file) => file.endsWith(".json"));
const enKeys = messageStrings(en).map(([key]) => key);

describe.each(locales)("messages/%s", (file) => {
  const messages = JSON.parse(readFileSync(`messages/${file}`, "utf8")) as object;

  test("follows the copy rules: no long dashes, emoji or hype words", () => {
    const problems = messageStrings(messages).flatMap(([key, text]) =>
      [...forbiddenCharacters(text), ...bannedWords(text)].map((what) => `${key}: ${what}`),
    );
    expect(problems).toEqual([]);
  });

  test("has exactly the keys of en.json", () => {
    expect(messageStrings(messages).map(([key]) => key)).toEqual(enKeys);
  });

  test("every message is valid ICU", () => {
    const invalid: string[] = [];
    const t = createTranslator({
      locale: "en",
      messages: messages as typeof en,
      onError: (error) => error.code === IntlErrorCode.INVALID_MESSAGE && invalid.push(error.message),
    });
    for (const [key] of messageStrings(messages)) t(key as never);
    expect(invalid).toEqual([]);
  });
});

test("has a label for every contact role", () => {
  expect(Object.keys(en.Contact.roles)).toEqual([...roles]);
});
