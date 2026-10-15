// SPDX-License-Identifier: MPL-2.0
import { createTranslator } from "next-intl";
import { vi } from "vitest";
import messages from "../../../messages/en.json";

// The real `next-intl/server` reads the request config through Next's runtime; tests use the English messages.
vi.mock(import("next-intl/server"), async (importOriginal) => ({
  ...(await importOriginal()),
  getTranslations: (async (options?: string | { namespace?: string }) => {
    const namespace = typeof options === "string" ? options : options?.namespace;
    return createTranslator({ locale: "en", messages, namespace: namespace as never });
  }) as never,
  getMessages: (async () => messages) as never,
  getLocale: async () => "en" as const,
}));
