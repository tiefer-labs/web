// SPDX-License-Identifier: MPL-2.0
import { describe, expect, test, vi } from "vitest";
import { createStore, MemoryStore, RedisStore } from "@/lib/contact/store";

const HOUR = 3_600_000;

describe("MemoryStore", () => {
  test("slides the window and fails closed when full of live keys", async () => {
    let now = 1_000_000;
    const store = new MemoryStore(2, 10, () => now);
    expect([await store.hit("a", 2, HOUR), await store.hit("a", 2, HOUR), await store.hit("a", 2, HOUR)]).toEqual([
      true,
      true,
      false,
    ]);
    now += 61 * 60_000;
    expect(await store.hit("a", 2, HOUR)).toBe(true);
    expect(await store.hit("b", 2, HOUR)).toBe(true);
    expect(await store.hit("c", 2, HOUR)).toBe(false);
    now += 2 * HOUR;
    expect(await store.hit("c", 2, HOUR)).toBe(true);
  });

  test("claims a key once until it expires", async () => {
    let now = 0;
    const store = new MemoryStore(10, 1, () => now);
    expect(await store.claim("t1", 1000)).toBe(true);
    expect(await store.claim("t1", 1000)).toBe(false);
    expect(await store.claim("t2", 1000)).toBe(false);
    now = 1001;
    expect(await store.claim("t2", 1000)).toBe(true);
  });
});

function fakeRedis() {
  const counts = new Map<string, number>();
  const taken = new Set<string>();
  return {
    incr: vi.fn<(key: string) => Promise<number>>(async (key) => {
      const count = (counts.get(key) ?? 0) + 1;
      counts.set(key, count);
      return count;
    }),
    pexpire: vi.fn<(key: string, ms: number) => Promise<number>>(async () => 1),
    setNx: vi.fn<(key: string) => Promise<boolean>>(async (key) => !taken.has(key) && Boolean(taken.add(key))),
  };
}

describe("RedisStore", () => {
  test("counts per window and sets the expiry once", async () => {
    const redis = fakeRedis();
    const store = new RedisStore(redis);
    expect([await store.hit("a", 2, HOUR), await store.hit("a", 2, HOUR), await store.hit("a", 2, HOUR)]).toEqual([
      true,
      true,
      false,
    ]);
    expect(redis.pexpire).toHaveBeenCalledOnce();
  });

  test("claims through SET NX", async () => {
    const store = new RedisStore(fakeRedis());
    expect(await store.claim("t", 1000)).toBe(true);
    expect(await store.claim("t", 1000)).toBe(false);
  });
});

test("falls back to memory without a Redis URL", () => {
  expect(createStore({ REDIS_URL: "", UPSTASH_REDIS_REST_URL: "", UPSTASH_REDIS_REST_TOKEN: "" })).toBeInstanceOf(
    MemoryStore,
  );
  expect(
    createStore({ REDIS_URL: "", UPSTASH_REDIS_REST_URL: "https://x.upstash.io", UPSTASH_REDIS_REST_TOKEN: "t" }),
  ).toBeInstanceOf(RedisStore);
});
