// SPDX-License-Identifier: MPL-2.0
import { securityTxt } from "@/lib/security-txt";

export function GET() {
  return new Response(securityTxt(), {
    headers: { "Content-Type": "text/plain; charset=utf-8", "Cache-Control": "public, max-age=86400" },
  });
}
