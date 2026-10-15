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

// FuzzContactPost sends arbitrary bodies to the contact form: the server
// never answers with a 5xx, and without a valid token nothing is sent.
func FuzzContactPost(f *testing.F) {
	f.Add("name=Ada&email=ada%40example.org&message=x&consent=yes&t=abc")
	f.Add("t=&website=x")
	f.Add("%zz&&==")
	f.Add(strings.Repeat("a=b&", 500))
	m := &fakeMailer{}
	h := newHarnessMailer(f, smtpEnv, m)
	f.Fuzz(func(t *testing.T, body string) {
		r := httptest.NewRequest(http.MethodPost, "/contact", strings.NewReader(body))
		r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		r.Header.Set("Sec-Fetch-Site", "same-origin")
		w := h.do(r)
		if w.Code >= 500 {
			t.Fatalf("status %d for %q", w.Code, body)
		}
		if m.count() != 0 {
			t.Fatalf("a message was sent for %q", body)
		}
	})
}

// FuzzAcceptsGzip parses arbitrary Accept-Encoding values without panic.
func FuzzAcceptsGzip(f *testing.F) {
	for _, s := range []string{"gzip", "br, gzip;q=0.8", "gzip;q=0", "*", ";;,,", "GZIP ; Q=0.000"} {
		f.Add(s)
	}
	f.Fuzz(func(t *testing.T, s string) {
		r := httptest.NewRequest(http.MethodGet, "/", nil)
		r.Header.Set("Accept-Encoding", s)
		acceptsGzip(r)
	})
}
