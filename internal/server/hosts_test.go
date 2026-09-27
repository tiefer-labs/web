// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package server

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var prodEnv = map[string]string{
	"ENV":           "production",
	"CONTACT_EMAIL": "hello@tiefer.space",
	"CSRF_SECRET":   strings.Repeat("k", 64),
	"TRUST_PROXY":   "true",
}

func TestCanonicalHost(t *testing.T) {
	h := newHarness(t, prodEnv)
	cases := []struct {
		method, host, target, proto string
		status                      int
		location                    string
	}{
		{"GET", "tiefer.space", "/privacy", "https", 200, ""},
		{"GET", "TIEFER.SPACE.", "/", "https", 200, ""},
		{"GET", "tiefer.space:443", "/", "https", 200, ""},
		{"GET", "www.tiefer.space", "/legal?x=1", "https", 301, "https://tiefer.space/legal?x=1"},
		{"POST", "www.tiefer.space", "/contact", "https", 308, "https://tiefer.space/contact"},
		{"GET", "tiefer.space", "/privacy", "http", 301, "https://tiefer.space/privacy"},
		{"GET", "tiefer.space", "//evil.example/x", "http", 301, "https://tiefer.space//evil.example/x"},
		{"GET", "evil.example", "/", "https", 421, ""},
		{"GET", "192.0.2.10", "/", "https", 421, ""},
		{"GET", "", "/", "https", 421, ""},
		{"GET", "127.0.0.1:8080", "/healthz", "http", 200, ""},
	}
	for _, c := range cases {
		r := httptest.NewRequest(c.method, c.target, nil)
		r.Host = c.host
		r.Header.Set("X-Forwarded-Proto", c.proto)
		w := h.do(t, r)
		if w.Code != c.status || w.Header().Get("Location") != c.location {
			t.Errorf("%s %s%s (%s): got %d %q, want %d %q", c.method, c.host, c.target, c.proto,
				w.Code, w.Header().Get("Location"), c.status, c.location)
		}
		if w.Header().Get("Content-Security-Policy") == "" {
			t.Errorf("%s %s%s: security headers missing", c.method, c.host, c.target)
		}
	}
}

func TestProductionCanonicalLinks(t *testing.T) {
	h := newHarness(t, prodEnv)
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r.Host = "tiefer.space"
	body := h.do(t, r).Body.String()
	for _, want := range []string{
		`<link rel="canonical" href="https://tiefer.space/">`,
		`property="og:image" content="https://tiefer.space/static/og-image.`,
		`"url":"https://tiefer.space/"`,
		`"email":"hello@tiefer.space"`,
		`href="mailto:hello@tiefer.space"`,
	} {
		if !strings.Contains(body, want) {
			t.Errorf("production page lacks %q", want)
		}
	}
}

func TestFetchMetadataIsolation(t *testing.T) {
	h := newHarness(t, smtpEnv)
	cases := []struct {
		method, target, site, mode, dest string
		status                           int
	}{
		{"GET", "/", "", "", "", 200},                            // no metadata: crawler, old browser
		{"GET", "/", "none", "navigate", "document", 200},        // typed address or bookmark
		{"GET", "/", "cross-site", "navigate", "document", 200},  // link from another site
		{"GET", "/privacy", "same-origin", "cors", "empty", 200}, // our own fetch
		{"GET", "/", "cross-site", "no-cors", "script", 403},     // <script src> on another site
		{"GET", "/static/og-image.png", "cross-site", "no-cors", "image", 403},
		{"GET", "/", "cross-site", "cors", "empty", 403},                // fetch() from another site
		{"GET", "/", "cross-site", "navigate", "iframe", 403},           // framing
		{"POST", "/contact", "cross-site", "navigate", "document", 403}, // cross-site form post
	}
	for _, c := range cases {
		r := httptest.NewRequest(c.method, c.target, nil)
		if c.site != "" {
			r.Header.Set("Sec-Fetch-Site", c.site)
			r.Header.Set("Sec-Fetch-Mode", c.mode)
			r.Header.Set("Sec-Fetch-Dest", c.dest)
		}
		if w := h.do(t, r); w.Code != c.status {
			t.Errorf("%s %s (%s, %s, %s): %d, want %d", c.method, c.target, c.site, c.mode, c.dest, w.Code, c.status)
		}
	}
}
