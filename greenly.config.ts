// SPDX-License-Identifier: MPL-2.0
import { defineConfig } from "greenly";

export default defineConfig({
  name: "Tiefer website",
  checks: [
    { name: "TypeScript", command: "pnpm tsc --noEmit --incremental false" },
    { name: "Oxfmt", command: "pnpm fmt:check", onFail: "pnpm fmt" },
    { name: "Oxlint", command: "pnpm lint" },
    { name: "Vitest", command: "pnpm test" },
  ],
});
