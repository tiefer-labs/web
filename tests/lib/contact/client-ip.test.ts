// SPDX-License-Identifier: MPL-2.0
import { describe, expect, test } from "vitest";
import { clientAddress, clientKey } from "@/lib/contact/client-ip";

describe("clientKey", () => {
  test.each([
    ["203.0.113.7", "203.0.113.7"],
    ["::ffff:203.0.113.7", "203.0.113.7"],
    ["2001:db8:1:2:3:4:5:6", "2001:db8:1:2::/64"],
    ["2001:db8:1:2:ffff::1234", "2001:db8:1:2::/64"],
    ["2001:DB8::1", "2001:db8:0:0::/64"],
    ["not an address", "unknown"],
    [null, "unknown"],
  ])("%s -> %s", (address, key) => {
    expect(clientKey(address)).toBe(key);
  });
});

describe("clientAddress", () => {
  const headers = new Headers({ "x-azure-clientip": "198.51.100.1", "x-forwarded-for": "203.0.113.7, 10.0.0.1" });

  test("trusts only the Front Door header behind Front Door", () => {
    expect(clientAddress(headers, true)).toBe("198.51.100.1");
    expect(clientAddress(new Headers({ "x-forwarded-for": "203.0.113.7" }), true)).toBeNull();
  });

  test("uses the first X-Forwarded-For entry elsewhere", () => {
    expect(clientAddress(headers, false)).toBe("203.0.113.7");
  });
});
