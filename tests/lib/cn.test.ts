// SPDX-License-Identifier: MPL-2.0
import { expect, test } from "vitest";
import { cn } from "@/lib/cn";

test("keeps a type scale size next to a colour", () => {
  expect(cn("text-small text-ink-60")).toBe("text-small text-ink-60");
});

test("lets a later size or colour win", () => {
  expect(cn("text-small text-ink", "text-label text-on-ink")).toBe("text-label text-on-ink");
  expect(cn("mb-6", false, "mb-3")).toBe("mb-3");
});
