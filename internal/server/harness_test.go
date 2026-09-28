// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package server

import (
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/tiefer-labs/web/internal/config"
	"github.com/tiefer-labs/web/internal/contact"
	"github.com/tiefer-labs/web/web"
)

// harness runs the full handler with a controllable clock.
type harness struct {
	t   testing.TB
	srv *Server
	cfg *config.Config

	mu  sync.Mutex
	now time.Time
}

func newHarness(t testing.TB, set map[string]string) *harness {
	t.Helper()
	return newHarnessMailer(t, set, nil)
}

// newHarnessMailer builds a harness whose contact form delivers to m.
func newHarnessMailer(t testing.TB, set map[string]string, m contact.Mailer) *harness {
	t.Helper()
	cfg, err := config.Load(func(k string) (string, bool) {
		v, ok := set[k]
		return v, ok
	})
	if err != nil {
		t.Fatal(err)
	}
	h := &harness{t: t, cfg: cfg, now: time.Date(2026, 10, 1, 9, 0, 0, 0, time.UTC)}
	h.srv, err = New(Options{
		Config:    cfg,
		Logger:    slog.New(slog.NewTextHandler(io.Discard, nil)),
		Now:       h.clock,
		Templates: web.Templates(),
		Static:    web.Static(),
		Mailer:    m,
	})
	if err != nil {
		t.Fatal(err)
	}
	return h
}

func (h *harness) clock() time.Time {
	h.mu.Lock()
	defer h.mu.Unlock()
	return h.now
}

func (h *harness) advance(d time.Duration) {
	h.mu.Lock()
	h.now = h.now.Add(d)
	h.mu.Unlock()
}

func (h *harness) do(r *http.Request) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	h.srv.Handler().ServeHTTP(w, r)
	return w
}

func (h *harness) get(path string) *httptest.ResponseRecorder {
	return h.do(httptest.NewRequest(http.MethodGet, path, nil))
}

// prod is a valid production configuration, optionally extended.
func prod(extra map[string]string) map[string]string {
	m := map[string]string{
		"ENV":         "production",
		"SITE_URL":    "https://tiefer.space",
		"CSRF_SECRET": "5e8a1f07c93b24d6af10e7c58b3d92f4",
	}
	for k, v := range extra {
		m[k] = v
	}
	return m
}
