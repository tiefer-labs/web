// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package server

import (
	"bytes"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// wantHeaders are required on every response, whatever its status.
var wantHeaders = map[string]string{
	"Content-Security-Policy":           "default-src 'none'",
	"Strict-Transport-Security":         "max-age=63072000; includeSubDomains",
	"X-Content-Type-Options":            "nosniff",
	"X-Frame-Options":                   "DENY",
	"Referrer-Policy":                   "no-referrer",
	"Permissions-Policy":                "camera=()",
	"Cross-Origin-Opener-Policy":        "same-origin",
	"Cross-Origin-Resource-Policy":      "same-origin",
	"Cross-Origin-Embedder-Policy":      "require-corp",
	"Origin-Agent-Cluster":              "?1",
	"X-Permitted-Cross-Domain-Policies": "none",
}

func checkHeaders(t *testing.T, name string, w *httptest.ResponseRecorder) {
	t.Helper()
	for k, v := range wantHeaders {
		if got := w.Header().Get(k); !strings.Contains(got, v) {
			t.Errorf("%s (%d): %s = %q, want it to contain %q", name, w.Code, k, got, v)
		}
	}
	for _, k := range []string{"Server", "X-Powered-By", "Set-Cookie"} {
		if w.Header().Get(k) != "" {
			t.Errorf("%s: unexpected header %s", name, k)
		}
	}
}

// TestHeadersOnEveryRoute requests every registered route plus the error
// paths (404, 405, 414, 400, 413) and checks the full header set.
func TestHeadersOnEveryRoute(t *testing.T) {
	for _, env := range []map[string]string{nil, prod(nil)} {
		h := newHarness(t, env)
		for _, rt := range h.srv.routes {
			var r *http.Request
			if rt.method == http.MethodPost {
				r = httptest.NewRequest(rt.method, rt.sample, strings.NewReader("x=1"))
				r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			} else {
				r = httptest.NewRequest(rt.method, rt.sample, nil)
			}
			r.Host = h.cfg.SiteURL.Host
			checkHeaders(t, rt.method+" "+rt.sample, h.do(r))
		}
		extra := []*http.Request{
			httptest.NewRequest(http.MethodGet, "/does-not-exist", nil),
			httptest.NewRequest(http.MethodDelete, "/", nil),
			httptest.NewRequest(http.MethodGet, "/"+strings.Repeat("a", MaxURLLength), nil),
			httptest.NewRequest(http.MethodGet, "/", strings.NewReader("body")),
			httptest.NewRequest(http.MethodPost, "/", strings.NewReader(strings.Repeat("a", MaxFormBytes+1))),
		}
		for _, r := range extra {
			r.Host = h.cfg.SiteURL.Host
			w := h.do(r)
			checkHeaders(t, r.Method+" "+r.URL.Path[:min(len(r.URL.Path), 20)], w)
			if w.Code < 400 {
				t.Errorf("%s %s: status %d, want an error", r.Method, r.URL.Path[:min(len(r.URL.Path), 20)], w.Code)
			}
		}
	}
}

func TestProductionCSP(t *testing.T) {
	csp := newHarness(t, prod(nil)).get("/healthz").Header().Get("Content-Security-Policy")
	for _, want := range []string{"upgrade-insecure-requests", "frame-ancestors 'none'", "base-uri 'none'", "form-action 'self'", "trusted-types 'none'"} {
		if !strings.Contains(csp, want) {
			t.Errorf("CSP lacks %q: %s", want, csp)
		}
	}
	for _, bad := range []string{"unsafe-inline", "unsafe-eval", "*", "data:", "http:"} {
		if strings.Contains(csp, bad) {
			t.Errorf("CSP contains %q: %s", bad, csp)
		}
	}
	hsts := newHarness(t, prod(map[string]string{"HSTS_PRELOAD": "true"})).get("/healthz").Header().Get("Strict-Transport-Security")
	if !strings.HasSuffix(hsts, "; preload") {
		t.Errorf("HSTS_PRELOAD=true: %q", hsts)
	}
}

func TestHealthz(t *testing.T) {
	w := newHarness(t, nil).get("/healthz")
	if w.Code != http.StatusOK || w.Body.String() != "ok" || w.Header().Get("Cache-Control") != "no-store" {
		t.Errorf("/healthz: %d %q %q", w.Code, w.Body.String(), w.Header().Get("Cache-Control"))
	}
}

func TestMethods(t *testing.T) {
	h := newHarness(t, nil)
	for _, m := range []string{"PUT", "DELETE", "PATCH", "OPTIONS", "TRACE", "CONNECT", "PROPFIND"} {
		w := h.do(httptest.NewRequest(m, "/healthz", nil))
		if w.Code != http.StatusMethodNotAllowed || w.Header().Get("Allow") != "GET, HEAD, POST" {
			t.Errorf("%s: %d, Allow %q", m, w.Code, w.Header().Get("Allow"))
		}
	}
	if w := h.do(httptest.NewRequest(http.MethodPost, "/healthz", nil)); w.Code != http.StatusMethodNotAllowed || w.Header().Get("Allow") != "GET, HEAD" {
		t.Errorf("POST /healthz: %d, Allow %q", w.Code, w.Header().Get("Allow"))
	}
	if w := h.do(httptest.NewRequest(http.MethodHead, "/healthz", nil)); w.Code != http.StatusOK {
		t.Errorf("HEAD /healthz: %d", w.Code)
	}
}

func TestProductionHost(t *testing.T) {
	h := newHarness(t, prod(nil))
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r.Host = "attacker.example"
	if w := h.do(r); w.Code != http.StatusMisdirectedRequest {
		t.Errorf("foreign host: %d, want 421", w.Code)
	}
	r = httptest.NewRequest(http.MethodGet, "/healthz", nil)
	r.Host = "10.0.0.4:8080"
	if w := h.do(r); w.Code != http.StatusOK {
		t.Errorf("/healthz on internal host: %d", w.Code)
	}
}

func TestFrontDoor(t *testing.T) {
	const id = "00000000-1111-2222-3333-444444444444"
	h := newHarness(t, prod(map[string]string{"BEHIND_FRONT_DOOR": "true", "FRONT_DOOR_ID": id}))
	r := httptest.NewRequest(http.MethodGet, "/does-not-exist", nil)
	if w := h.do(r); w.Code != http.StatusForbidden {
		t.Errorf("without X-Azure-FDID: %d, want 403", w.Code)
	}
	r.Header.Set("X-Azure-FDID", "00000000-1111-2222-3333-555555555555")
	if w := h.do(r); w.Code != http.StatusForbidden {
		t.Errorf("wrong X-Azure-FDID: %d, want 403", w.Code)
	}
	r.Header.Set("X-Azure-FDID", id)
	r.Host = "tiefer-web-app.azurewebsites.net" // the origin host Front Door sends
	if w := h.do(r); w.Code != http.StatusNotFound {
		t.Errorf("right X-Azure-FDID: %d, want the route's own 404", w.Code)
	}
	if w := h.get("/healthz"); w.Code != http.StatusOK {
		t.Errorf("/healthz without Front Door: %d", w.Code)
	}

	r.Header.Set("X-Azure-ClientIP", "203.0.113.7")
	r.RemoteAddr = "10.0.0.1:1234"
	if got := h.srv.clientAddr(r).String(); got != "203.0.113.7" {
		t.Errorf("clientAddr behind Front Door = %s", got)
	}
	direct := newHarness(t, nil)
	if got := direct.srv.clientAddr(r).String(); got != "10.0.0.1" {
		t.Errorf("clientAddr without Front Door must ignore the header, got %s", got)
	}
}

// TestLogsWithoutPersonalData checks that the request log holds no IP
// address, query string or user agent.
func TestLogsWithoutPersonalData(t *testing.T) {
	var buf bytes.Buffer
	h := newHarness(t, nil)
	h.srv.log = slog.New(slog.NewJSONHandler(&buf, nil))
	h.srv.handler = h.srv.logRequests(h.srv.handler)
	r := httptest.NewRequest(http.MethodGet, "/missing?email=someone@example.org", nil)
	r.RemoteAddr = "198.51.100.23:5555"
	r.Header.Set("User-Agent", "UniqueAgent/1.0")
	h.do(r)
	out := buf.String()
	if !strings.Contains(out, `"path":"/missing"`) {
		t.Fatalf("request not logged: %s", out)
	}
	for _, s := range []string{"198.51.100.23", "someone@example.org", "UniqueAgent"} {
		if strings.Contains(out, s) {
			t.Errorf("log contains %q: %s", s, out)
		}
	}
}

func TestRecoverPanics(t *testing.T) {
	h := newHarness(t, nil)
	var buf bytes.Buffer
	h.srv.log = slog.New(slog.NewJSONHandler(&buf, nil))
	panicky := h.srv.recoverPanics(h.srv.securityHeaders(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		panic("boom")
	})))
	w := httptest.NewRecorder()
	panicky.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/", nil))
	if w.Code != http.StatusInternalServerError || strings.Contains(w.Body.String(), "boom") {
		t.Errorf("panic: %d %q", w.Code, w.Body.String())
	}
	checkHeaders(t, "panic", w)
}

func TestServerLimits(t *testing.T) {
	hs := newHarness(t, nil).srv.HTTPServer(":0")
	if hs.ReadHeaderTimeout == 0 || hs.ReadTimeout == 0 || hs.WriteTimeout == 0 || hs.IdleTimeout == 0 {
		t.Error("every server timeout must be set")
	}
	if hs.MaxHeaderBytes != MaxHeaderBytes || !hs.DisableGeneralOptionsHandler {
		t.Error("header limit and OPTIONS handling")
	}
	// A real server rejects oversized headers before any handler runs.
	ts := httptest.NewUnstartedServer(hs.Handler)
	ts.Config.MaxHeaderBytes = MaxHeaderBytes
	ts.Start()
	defer ts.Close()
	req, _ := http.NewRequest(http.MethodGet, ts.URL+"/healthz", nil)
	req.Header.Set("X-Big", strings.Repeat("a", 2*MaxHeaderBytes))
	res, err := ts.Client().Do(req)
	if err != nil {
		t.Fatal(err)
	}
	_, _ = io.Copy(io.Discard, res.Body)
	res.Body.Close()
	if res.StatusCode != http.StatusRequestHeaderFieldsTooLarge {
		t.Errorf("oversized headers: %d, want 431", res.StatusCode)
	}
}
