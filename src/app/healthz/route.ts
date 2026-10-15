// SPDX-License-Identifier: MPL-2.0
import { connection } from "next/server";

/** Liveness probe for Azure App Service and the Docker HEALTHCHECK. */
export async function GET() {
  await connection();
  return new Response("ok", { headers: { "Content-Type": "text/plain; charset=utf-8", "Cache-Control": "no-store" } });
}
