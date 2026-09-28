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
	"github.com/tiefer-labs/web/web"
)

// harness runs the full handler with a controllable clock.
type harness struct {
	t   *testing.T
	srv *Server
	cfg *config.Config

	mu  sync.Mutex
	now time.Time
}

func newHarness(t *testing.T, set map[string]string) *harness {
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
		"CSRF_SECRET": "0123456789abcdef0123456789abcdef-test",
	}
	for k, v := range extra {
		m[k] = v
	}
	return m
}
