// SPDX-License-Identifier: MPL-2.0
import "server-only";
import { env } from "@/lib/env.server";
import { issueToken } from "./token";

export function freshToken() {
  return issueToken(env.CSRF_SECRET);
}
