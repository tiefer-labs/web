// SPDX-License-Identifier: MPL-2.0
import { defineConfig } from "oxfmt";

export default defineConfig({
  printWidth: 120,
  singleQuote: false,
  insertFinalNewline: true,
  ignorePatterns: ["old-vanilla/**"],
  sortTailwindcss: {
    stylesheet: "./src/app/globals.css",
    functions: ["clsx", "cn"],
    attributes: ["classNames", "tw"],
  },
  sortImports: {
    newlinesBetween: false,
    internalPattern: ["~/", "@/", "#", "@test/"],
  },
});
