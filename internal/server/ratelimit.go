// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package server

import (
	"net/netip"
	"sync"
	"time"
)

// maxClients bounds the memory the per-client limiter can use. When that
// many clients are tracked at once, new clients are refused until old
// entries expire: failing closed is safer than growing without limit.
const maxClients = 50_000

// rateLimiter allows a number of events per key within a sliding window.
// Keys (client IP addresses) live only in memory and are dropped once
// their window has passed.
type rateLimiter struct {
	mu        sync.Mutex
	limit     int
	maxKeys   int
	window    time.Duration
	now       func() time.Time
	hits      map[string][]time.Time
	lastSweep time.Time
}

func newRateLimiter(limit int, window time.Duration, now func() time.Time) *rateLimiter {
	return &rateLimiter{limit: limit, maxKeys: maxClients, window: window, now: now, hits: map[string][]time.Time{}}
}

// Allow records an event for key and reports whether it is within the
// limit.
func (l *rateLimiter) Allow(key string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	now := l.now()
	cutoff := now.Add(-l.window)
	_, known := l.hits[key]
	if now.Sub(l.lastSweep) > l.window || (!known && len(l.hits) >= l.maxKeys) {
		for k, ts := range l.hits {
			if len(ts) == 0 || !ts[len(ts)-1].After(cutoff) {
				delete(l.hits, k)
			}
		}
		l.lastSweep = now
	}
	if _, ok := l.hits[key]; !ok && len(l.hits) >= l.maxKeys {
		return false
	}
	ts := l.hits[key]
	i := 0
	for i < len(ts) && !ts[i].After(cutoff) {
		i++
	}
	ts = ts[i:]
	if len(ts) >= l.limit {
		l.hits[key] = ts
		return false
	}
	l.hits[key] = append(ts, now)
	return true
}

// clientKey turns a client address into a rate limit key. IPv6 clients
// usually control a whole /64 network, so they share one key per /64;
// otherwise changing the address would bypass the limit.
func clientKey(addr string) string {
	ip, err := netip.ParseAddr(addr)
	if err != nil {
		return addr
	}
	ip = ip.Unmap()
	if ip.Is6() {
		p, _ := ip.Prefix(64)
		return p.String()
	}
	return ip.String()
}
