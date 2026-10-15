// SPDX-License-Identifier: MPL-2.0
import { execFileSync } from "node:child_process";
import { existsSync, readFileSync } from "node:fs";
import { expect, test } from "vitest";
import { forbiddenCharacters } from "./test-utils/text-rules";

const files = execFileSync("git", ["ls-files", "--cached", "--others", "--exclude-standard"], { encoding: "utf8" })
  .split("\n")
  .filter((file) => file && !file.startsWith("old-vanilla/") && existsSync(file));

const BINARY = /\.(png|ico|jpe?g|webp|woff2?|ttf)$/;
const SOURCE = /\.(tsx?|css)$/;
const HEADER = "SPDX-License-Identifier: MPL-2.0";

// `next dev` writes and re-adds this block in AGENTS.md, so it is not ours to edit.
const GENERATED = /<!-- BEGIN:nextjs-agent-rules -->[\s\S]*?<!-- END:nextjs-agent-rules -->/;

test("every source file starts with the MPL-2.0 SPDX header", () => {
  const missing = files.filter(
    (file) =>
      SOURCE.test(file) && !file.endsWith(".d.ts") && !readFileSync(file, "utf8").split("\n", 1)[0]?.includes(HEADER),
  );
  expect(missing).toEqual([]);
});

test("no text file contains long dashes, emoji or decorative symbols", () => {
  const problems = files
    .filter((file) => !BINARY.test(file) && file !== "LICENSE" && file !== "pnpm-lock.yaml")
    .flatMap((file) =>
      forbiddenCharacters(readFileSync(file, "utf8").replace(GENERATED, "")).map((what) => `${file}:${what}`),
    );
  expect(problems).toEqual([]);
});
