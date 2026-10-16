// SPDX-License-Identifier: MPL-2.0
import { timingSafeEqual } from "node:crypto";
import createMiddleware from "next-intl/middleware";
import type { NextRequest } from "next/server";
import { routing } from "@/i18n/routing";

const handleI18nRouting = createMiddleware(routing);

/** Behind Azure Front Door, only requests carrying our Front Door profile id are served. */
function fromOurFrontDoor(request: NextRequest) {
  if (process.env.BEHIND_FRONT_DOOR !== "true" || !process.env.FRONT_DOOR_ID) return true;
  const expected = Buffer.from(process.env.FRONT_DOOR_ID);
  const actual = Buffer.from(request.headers.get("x-azure-fdid") ?? "");
  return actual.length === expected.length && timingSafeEqual(actual, expected);
}

export default function proxy(request: NextRequest) {
  if (!fromOurFrontDoor(request)) return new Response("Forbidden\n", { status: 403 });
  return handleI18nRouting(request);
}

export const config = {
  // Everything except Next internals, the health probe and files with an extension (robots.txt, icons...).
  matcher: "/((?!_next|_vercel|healthz|.*\\..*).*)",
};
