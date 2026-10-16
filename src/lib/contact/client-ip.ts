// SPDX-License-Identifier: MPL-2.0
import { isIPv4, isIPv6 } from "node:net";

/** The first four groups of an IPv6 address, written out in full. */
function network64(address: string): string {
  const [head = "", tail = ""] = address.split("::");
  const left = head ? head.split(":") : [];
  const right = tail ? tail.split(":") : [];
  const groups = address.includes("::")
    ? [...left, ...Array(8 - left.length - right.length).fill("0"), ...right]
    : left;
  return `${groups
    .slice(0, 4)
    .map((g) => Number.parseInt(g, 16).toString(16))
    .join(":")}::/64`;
}

/** Rate-limit key for an address: IPv4 as is, IPv6 by its /64 network, since one host usually owns a /64. */
export function clientKey(address: string | null | undefined): string {
  const ip = address?.trim().replace(/^::ffff:(?=\d+\.\d+\.\d+\.\d+$)/i, "") ?? "";
  if (isIPv4(ip)) return ip;
  if (isIPv6(ip) && !ip.includes(".")) return network64(ip.toLowerCase());
  return "unknown";
}

/** Behind Front Door only its own header can be trusted; elsewhere (Vercel) the platform sets X-Forwarded-For. */
export function clientAddress(headers: Headers, behindFrontDoor: boolean): string | null {
  if (behindFrontDoor) return headers.get("x-azure-clientip");
  return headers.get("x-forwarded-for")?.split(",")[0] ?? headers.get("x-real-ip");
}
