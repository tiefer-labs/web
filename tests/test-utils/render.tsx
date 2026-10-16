// SPDX-License-Identifier: MPL-2.0
import { render, type RenderOptions } from "@testing-library/react";
import { NextIntlClientProvider } from "next-intl";
import { AppRouterContext, type AppRouterInstance } from "next/dist/shared/lib/app-router-context.shared-runtime";
import { PathnameContext, SearchParamsContext } from "next/dist/shared/lib/hooks-client-context.shared-runtime";
import type { ReactNode } from "react";
import { prerender } from "react-dom/static";
import { vi } from "vitest";
import messages from "../../messages/en.json";

function createTestRouter(): AppRouterInstance {
  return {
    back: vi.fn<AppRouterInstance["back"]>(),
    forward: vi.fn<AppRouterInstance["forward"]>(),
    refresh: vi.fn<AppRouterInstance["refresh"]>(),
    push: vi.fn<AppRouterInstance["push"]>(),
    replace: vi.fn<AppRouterInstance["replace"]>(),
    prefetch: vi.fn<AppRouterInstance["prefetch"]>(),
    bfcacheId: "test",
  };
}

/** The App Router and next-intl contexts that links and translations read. */
export function Providers({ pathname = "/", children }: { pathname?: string; children: ReactNode }) {
  return (
    <AppRouterContext.Provider value={createTestRouter()}>
      <PathnameContext.Provider value={pathname}>
        <SearchParamsContext.Provider value={new URLSearchParams()}>
          <NextIntlClientProvider locale="en" messages={messages}>
            {children}
          </NextIntlClientProvider>
        </SearchParamsContext.Provider>
      </PathnameContext.Provider>
    </AppRouterContext.Provider>
  );
}

export function renderClient(ui: ReactNode, options: { pathname?: string } & Omit<RenderOptions, "wrapper"> = {}) {
  const { pathname, ...rest } = options;
  return render(ui, { ...rest, wrapper: ({ children }) => <Providers pathname={pathname}>{children}</Providers> });
}

/**
 * Renders a Server Component tree (async components included) to static HTML, like Next's prerender, and mounts
 * it into the document; query it with `screen`. The markup is not hydrated.
 */
export async function renderServer(ui: ReactNode, { pathname = "/" }: { pathname?: string } = {}) {
  const { prelude } = await prerender(<Providers pathname={pathname}>{ui}</Providers>);
  const container = document.createElement("div");
  container.innerHTML = await new Response(prelude).text();
  document.body.append(container);
  return { container };
}
