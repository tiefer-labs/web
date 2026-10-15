// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package contact

import (
	"net/netip"
	"sync"
	"time"
)

// Limiter counts events per key in a sliding window, in memory only.
type Limiter struct {
	max    int
	window time.Duration
	keys   int // at most this many keys are tracked
	now    func() time.Time

	mu     sync.Mutex
	events map[string][]time.Time
}

// NewLimiter allows max events per key within window, for at most keys
// distinct keys. now may be nil for time.Now.
func NewLimiter(max int, window time.Duration, keys int, now func() time.Time) *Limiter {
	if now == nil {
		now = time.Now
	}
	return &Limiter{max: max, window: window, keys: keys, now: now, events: map[string][]time.Time{}}
}

// Allow records an event for key and reports whether it is within the
// limit. When the table is full of active keys it fails closed.
func (l *Limiter) Allow(key string) bool {
	now := l.now()
	cutoff := now.Add(-l.window)
	l.mu.Lock()
	defer l.mu.Unlock()
	ev := l.events[key]
	i := 0
	for i < len(ev) && !ev[i].After(cutoff) {
		i++
	}
	ev = ev[i:]
	if len(ev) >= l.max {
		l.events[key] = ev
		return false
	}
	if _, known := l.events[key]; !known && len(l.events) >= l.keys {
		for k, e := range l.events {
			if len(e) == 0 || !e[len(e)-1].After(cutoff) {
				delete(l.events, k)
			}
		}
		if len(l.events) >= l.keys {
			return false
		}
	}
	l.events[key] = append(ev, now)
	return true
}

// ClientKey returns the rate limit key of an address: the address for
// IPv4, its /64 network for IPv6, since one host usually owns a /64.
func ClientKey(a netip.Addr) string {
	if !a.IsValid() {
		return "unknown"
	}
	a = a.Unmap()
	if a.Is6() {
		p, _ := a.Prefix(64)
		return p.String()
	}
	return a.String()
}
