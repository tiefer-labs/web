// SPDX-License-Identifier: MPL-2.0
import { defineConfig } from "oxlint";

export default defineConfig({
  plugins: ["typescript", "nextjs", "unicorn", "import", "jsx-a11y", "react", "react-perf", "vitest"],
  categories: {
    suspicious: "warn",
  },
  options: {
    typeAware: true,
    typeCheck: true,
  },
  rules: {
    eqeqeq: "warn",
    "no-throw-literal": "warn",
    "import/no-unassigned-import": [
      "warn",
      { allow: ["**/globals.css", "server-only", "@testing-library/jest-dom/vitest", "./shared"] },
    ],
    // Safari drops list semantics from lists styled without bullets, so they keep role="list".
    "jsx-a11y/no-redundant-roles": ["warn", { ul: ["list"], ol: ["list"] }],
    "unicorn/prefer-node-protocol": "warn",
    "react/react-in-jsx-scope": "off",
    "typescript/consistent-type-imports": "warn",
  },
  overrides: [
    {
      // Fixtures and mocks cast on purpose.
      files: ["tests/**"],
      rules: { "typescript/no-unsafe-type-assertion": "off" },
    },
  ],
});
