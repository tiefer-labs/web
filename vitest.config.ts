// SPDX-License-Identifier: MPL-2.0
import { fileURLToPath } from "node:url";
import babel from "@rolldown/plugin-babel";
import react, { reactCompilerPreset } from "@vitejs/plugin-react";
import { defineConfig } from "vitest/config";

const emptyModule = fileURLToPath(new URL("./tests/test-utils/empty-module.ts", import.meta.url));

export default defineConfig({
  plugins: [react(), babel({ presets: [reactCompilerPreset()] })],
  resolve: {
    tsconfigPaths: true,
    // `server-only` throws outside the react-server condition; Next strips it at build time.
    alias: { "server-only": emptyModule },
  },
  test: {
    env: { NODE_ENV: "test", TZ: "UTC", SITE_URL: "https://tiefer.space" },
    fsModuleCache: true,
    clearMocks: true,
    restoreMocks: true,
    unstubEnvs: true,
    unstubGlobals: true,
    css: false,
    // next-intl imports `next/navigation` without an extension, which only resolves through Vite.
    server: { deps: { inline: ["next-intl"] } },
    projects: [
      {
        extends: true,
        test: {
          name: "node",
          environment: "node",
          include: ["tests/**/*.test.ts"],
          setupFiles: ["./tests/test-utils/setup/shared.ts"],
        },
      },
      {
        extends: true,
        test: {
          name: "dom",
          environment: "./tests/test-utils/jsdom-vm-environment.ts",
          // One jsdom per worker instead of per file; each file still gets its own VM context.
          pool: "vmThreads",
          testTimeout: 15_000,
          environmentOptions: { jsdom: { url: "http://localhost:3000/" } },
          include: ["tests/**/*.test.tsx"],
          setupFiles: ["./tests/test-utils/setup/dom.ts"],
        },
      },
    ],
  },
});
