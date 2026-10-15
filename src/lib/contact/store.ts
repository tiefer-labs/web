// SPDX-License-Identifier: MPL-2.0
import "server-only";
import { Redis as UpstashRedis } from "@upstash/redis";
import { createClient } from "redis";

/** Shared state for rate limits and used form tokens. */
export interface Store {
  /** Records one event for `key`; false when `limit` events already happened within `windowMs`. */
  hit(key: string, limit: number, windowMs: number): Promise<boolean>;
  /** Marks `key` as taken for `ttlMs`; false when it was already taken. */
  claim(key: string, ttlMs: number): Promise<boolean>;
}

/**
 * Per-process store, fine for a single instance. Both tables are bounded and fail closed when full of live
 * entries, so a flood cannot exhaust memory.
 */
export class MemoryStore implements Store {
  private readonly events = new Map<string, number[]>();
  private readonly claims = new Map<string, number>();

  constructor(
    private readonly maxKeys = 50_000,
    private readonly maxClaims = 100_000,
    private readonly now: () => number = Date.now,
  ) {}

  async hit(key: string, limit: number, windowMs: number) {
    const now = this.now();
    const cutoff = now - windowMs;
    const recent = (this.events.get(key) ?? []).filter((t) => t > cutoff);
    if (recent.length >= limit) {
      this.events.set(key, recent);
      return false;
    }
    if (!this.events.has(key) && this.events.size >= this.maxKeys) {
      for (const [k, times] of this.events) if (!times.some((t) => t > cutoff)) this.events.delete(k);
      if (this.events.size >= this.maxKeys) return false;
    }
    this.events.set(key, [...recent, now]);
    return true;
  }

  async claim(key: string, ttlMs: number) {
    const now = this.now();
    const expiry = this.claims.get(key);
    if (expiry !== undefined && expiry > now) return false;
    if (this.claims.size >= this.maxClaims) {
      for (const [k, exp] of this.claims) if (exp <= now) this.claims.delete(k);
      if (this.claims.size >= this.maxClaims) return false;
    }
    this.claims.set(key, now + ttlMs);
    return true;
  }
}

type RedisCommands = {
  incr(key: string): Promise<number>;
  pexpire(key: string, ms: number): Promise<unknown>;
  setNx(key: string, ttlMs: number): Promise<boolean>;
};

/** Fixed-window counters and SET NX claims, shared by every instance. */
export class RedisStore implements Store {
  constructor(
    private readonly redis: RedisCommands,
    private readonly prefix = "tiefer:contact:",
  ) {}

  async hit(key: string, limit: number, windowMs: number) {
    const window = Math.floor(Date.now() / windowMs);
    const id = `${this.prefix}hit:${key}:${window}`;
    const count = await this.redis.incr(id);
    if (count === 1) await this.redis.pexpire(id, windowMs);
    return count <= limit;
  }

  claim(key: string, ttlMs: number) {
    return this.redis.setNx(`${this.prefix}claim:${key}`, ttlMs);
  }
}

function upstash(url: string, token: string): RedisCommands {
  const redis = new UpstashRedis({ url, token });
  return {
    incr: (key) => redis.incr(key),
    pexpire: (key, ms) => redis.pexpire(key, ms),
    setNx: async (key, ttlMs) => (await redis.set(key, 1, { nx: true, px: ttlMs })) === "OK",
  };
}

function nodeRedis(url: string): RedisCommands {
  const client = createClient({ url });
  client.on("error", (error: unknown) => console.error("redis error", error));
  const ready = client.connect();
  return {
    incr: async (key) => (await ready).incr(key),
    pexpire: async (key, ms) => (await ready).pExpire(key, ms),
    setNx: async (key, ttlMs) =>
      (await (await ready).set(key, "1", { condition: "NX", expiration: { type: "PX", value: ttlMs } })) === "OK",
  };
}

type StoreEnv = { REDIS_URL: string; UPSTASH_REDIS_REST_URL: string; UPSTASH_REDIS_REST_TOKEN: string };

/** Upstash when its REST URL is set, then a plain Redis URL, otherwise memory. */
export function createStore(env: StoreEnv): Store {
  if (env.UPSTASH_REDIS_REST_URL)
    return new RedisStore(upstash(env.UPSTASH_REDIS_REST_URL, env.UPSTASH_REDIS_REST_TOKEN));
  if (env.REDIS_URL) return new RedisStore(nodeRedis(env.REDIS_URL));
  return new MemoryStore();
}
