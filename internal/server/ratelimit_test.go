// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package server

import (
	"strconv"
	"testing"
	"time"
)

func TestClientKey(t *testing.T) {
	cases := map[string]string{
		"192.0.2.7":                 "192.0.2.7",
		"::ffff:192.0.2.7":          "192.0.2.7",
		"2001:db8:1:2:3:4:5:6":      "2001:db8:1:2::/64",
		"2001:db8:1:2:ffff:ffff::1": "2001:db8:1:2::/64",
		"not an ip":                 "not an ip",
	}
	for in, want := range cases {
		if got := clientKey(in); got != want {
			t.Errorf("clientKey(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestRateLimiterBoundsMemory(t *testing.T) {
	now := time.Unix(1_800_000_000, 0)
	l := newRateLimiter(5, time.Hour, func() time.Time { return now })
	l.maxKeys = 3
	for i := range 3 {
		if !l.Allow(strconv.Itoa(i)) {
			t.Fatalf("client %d refused", i)
		}
	}
	if l.Allow("new") {
		t.Error("a new client must be refused while the table is full")
	}
	if !l.Allow("1") {
		t.Error("a known client must still be served")
	}
	now = now.Add(time.Hour + time.Second)
	if !l.Allow("new") {
		t.Error("expired entries must make room again")
	}
	if len(l.hits) != 1 {
		t.Errorf("%d entries kept, want 1", len(l.hits))
	}
}
