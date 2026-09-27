// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package server

import (
	"compress/gzip"
	"context"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"html"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/tiefer-labs/web/internal/config"
	"github.com/tiefer-labs/web/web"
)

// fakeMailer records messages instead of sending them.
type fakeMailer struct {
	mu   sync.Mutex
	msgs []string
	err  error
}

func (m *fakeMailer) Send(_ context.Context, _ string, _ []string, msg []byte) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.err != nil {
		return m.err
	}
	m.msgs = append(m.msgs, string(msg))
	return nil
}

func (m *fakeMailer) count() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return len(m.msgs)
}

type harness struct {
	srv    *Server
	h      http.Handler
	mailer *fakeMailer
	now    time.Time
	mu     sync.Mutex
}

func (h *harness) clock() time.Time {
	h.mu.Lock()
	defer h.mu.Unlock()
	return h.now
}

func (h *harness) advance(d time.Duration) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.now = h.now.Add(d)
}

var smtpEnv = map[string]string{
	"SMTP_HOST":    "smtp.example.org",
	"SMTP_PORT":    "587",
	"CONTACT_TO":   "team@example.org",
	"CONTACT_FROM": "web@example.org",
	"REPO_URL":     "https://github.com/tiefer-labs/web",
}

func newHarness(t *testing.T, env map[string]string) *harness {
	t.Helper()
	cfg, err := config.FromLookup(func(k string) (string, bool) {
		v, ok := env[k]
		return v, ok
	})
	if err != nil {
		t.Fatal(err)
	}
	h := &harness{mailer: &fakeMailer{}, now: time.Unix(1_800_000_000, 0)}
	h.srv, err = New(Options{
		Config:    cfg,
		Logger:    slog.New(slog.NewTextHandler(io.Discard, nil)),
		Templates: web.Templates(),
		Static:    web.Static(),
		Mailer:    h.mailer,
		Now:       h.clock,
	})
	if err != nil {
		t.Fatal(err)
	}
	h.h = h.srv.Handler()
	return h
}

func (h *harness) do(t *testing.T, r *http.Request) *httptest.ResponseRecorder {
	t.Helper()
	w := httptest.NewRecorder()
	h.h.ServeHTTP(w, r)
	return w
}

func (h *harness) get(t *testing.T, path string) *httptest.ResponseRecorder {
	return h.do(t, httptest.NewRequest(http.MethodGet, path, nil))
}

var tokenRE = regexp.MustCompile(`name="t" value="([^"]+)"`)

// formToken loads the index page and returns the form token in it.
func (h *harness) formToken(t *testing.T) string {
	t.Helper()
	m := tokenRE.FindStringSubmatch(h.get(t, "/").Body.String())
	if m == nil {
		t.Fatal("no form token on the index page")
	}
	return m[1]
}

func validValues(token string) url.Values {
	return url.Values{
		"t":        {token},
		"name":     {"Aysel Mammadova"},
		"email":    {"aysel@example.org"},
		"org":      {"Example Trading"},
		"role":     {"trader-analyst"},
		"question": {"Which ports saw more vessel activity?"},
		"consent":  {"yes"},
		"website":  {""},
	}
}

func post(path string, v url.Values, remote string) *http.Request {
	r := httptest.NewRequest(http.MethodPost, path, strings.NewReader(v.Encode()))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Set("Sec-Fetch-Site", "same-origin")
	if remote != "" {
		r.RemoteAddr = remote
	}
	return r
}

// routes lists a request for every route, for header and status checks.
var routes = []struct {
	path   string
	status int
}{
	{"/", 200},
	{"/?sent=1", 200},
	{"/legal", 200},
	{"/privacy", 200},
	{"/acceptable-use", 200},
	{"/healthz", 200},
	{"/robots.txt", 200},
	{"/.well-known/security.txt", 200},
	{"/sitemap.xml", 200},
	{"/site.webmanifest", 200},
	{"/favicon.ico", 200},
	{"/apple-touch-icon.png", 200},
	{"/static/css/site.css", 200},
	{"/static/og-image.png", 200},
	{"/does-not-exist", 404},
	{"/legal/", 404},
	{"/static/nope.css", 404},
}

func TestRoutesStatus(t *testing.T) {
	h := newHarness(t, smtpEnv)
	for _, rt := range routes {
		if w := h.get(t, rt.path); w.Code != rt.status {
			t.Errorf("GET %s = %d, want %d", rt.path, w.Code, rt.status)
		}
	}
}

func TestSecurityHeadersOnEveryRoute(t *testing.T) {
	h := newHarness(t, smtpEnv)
	want := map[string]string{
		"X-Content-Type-Options":            "nosniff",
		"Referrer-Policy":                   "strict-origin-when-cross-origin",
		"X-Frame-Options":                   "DENY",
		"Cross-Origin-Opener-Policy":        "same-origin",
		"Cross-Origin-Resource-Policy":      "same-origin",
		"Cross-Origin-Embedder-Policy":      "require-corp",
		"X-Permitted-Cross-Domain-Policies": "none",
		"Origin-Agent-Cluster":              "?1",
	}
	check := func(name string, w *httptest.ResponseRecorder) {
		for k, v := range want {
			if got := w.Header().Get(k); got != v {
				t.Errorf("%s: %s = %q, want %q", name, k, got, v)
			}
		}
		csp := w.Header().Get("Content-Security-Policy")
		for _, part := range []string{"default-src 'none'", "script-src 'self' 'sha256-", "script-src-attr 'none'",
			"style-src-attr 'none'", "frame-ancestors 'none'", "base-uri 'none'", "object-src 'none'",
			"form-action 'self'", "require-trusted-types-for 'script'", "trusted-types 'none'"} {
			if !strings.Contains(csp, part) {
				t.Errorf("%s: CSP %q lacks %q", name, csp, part)
			}
		}
		if strings.Contains(csp, "unsafe-inline") || strings.Contains(csp, "http") {
			t.Errorf("%s: CSP allows inline code or external origins: %q", name, csp)
		}
		if !strings.HasPrefix(w.Header().Get("Strict-Transport-Security"), "max-age=") {
			t.Errorf("%s: no Strict-Transport-Security", name)
		}
		if !strings.Contains(w.Header().Get("Permissions-Policy"), "camera=()") {
			t.Errorf("%s: no Permissions-Policy", name)
		}
	}
	for _, rt := range routes {
		check("GET "+rt.path, h.get(t, rt.path))
	}
	check("POST /contact", h.do(t, post("/contact", url.Values{}, "")))
}

func TestProductionHSTSPreload(t *testing.T) {
	h := newHarness(t, prodEnv)
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r.Host = "tiefer.space"
	w := h.do(t, r)
	if got := w.Header().Get("Strict-Transport-Security"); got != "max-age=63072000; includeSubDomains; preload" {
		t.Errorf("HSTS = %q", got)
	}
	if !strings.Contains(w.Header().Get("Content-Security-Policy"), "upgrade-insecure-requests") {
		t.Error("production CSP lacks upgrade-insecure-requests")
	}
}

func TestCSPHashMatchesJSONLD(t *testing.T) {
	h := newHarness(t, smtpEnv)
	w := h.get(t, "/")
	m := regexp.MustCompile(`<script type="application/ld\+json">(.*?)</script>`).FindStringSubmatch(w.Body.String())
	if m == nil {
		t.Fatal("no JSON-LD block")
	}
	var org map[string]any
	if err := json.Unmarshal([]byte(m[1]), &org); err != nil || org["@type"] != "Organization" {
		t.Fatalf("bad JSON-LD %q: %v", m[1], err)
	}
	sum := sha256.Sum256([]byte(m[1]))
	hash := "'sha256-" + base64.StdEncoding.EncodeToString(sum[:]) + "'"
	if !strings.Contains(w.Header().Get("Content-Security-Policy"), hash) {
		t.Errorf("CSP does not allow the JSON-LD block (want %s)", hash)
	}
}

func TestCacheControl(t *testing.T) {
	h := newHarness(t, smtpEnv)
	if got := h.get(t, "/").Header().Get("Cache-Control"); got != "private, no-cache" {
		t.Errorf("index Cache-Control = %q, want private, no-cache", got)
	}
	tok := h.formToken(t)
	h.advance(10 * time.Second)
	v := validValues(tok)
	v.Del("name")
	if got := h.do(t, post("/contact", v, "")).Header().Get("Cache-Control"); got != "no-store" {
		t.Errorf("form error page Cache-Control = %q, want no-store", got)
	}
	if got := h.do(t, post("/contact", validValues(tok), "")).Header().Get("Cache-Control"); got != "no-store" {
		t.Errorf("form redirect Cache-Control = %q, want no-store", got)
	}
}

func TestNotFoundPage(t *testing.T) {
	h := newHarness(t, nil)
	w := h.get(t, "/no/such/page")
	if w.Code != http.StatusNotFound {
		t.Fatalf("status = %d", w.Code)
	}
	body := w.Body.String()
	for _, want := range []string{"This page does not exist.", `<meta name="robots" content="noindex">`, "site-header", "site-footer"} {
		if !strings.Contains(body, want) {
			t.Errorf("404 page lacks %q", want)
		}
	}
}

func TestIndexPage(t *testing.T) {
	h := newHarness(t, smtpEnv)
	body := h.get(t, "/").Body.String()
	for _, want := range []string{
		`<html lang="en">`,
		"<title>Tiefer | Space intelligence you can ask</title>",
		`<link rel="canonical" href="http://localhost:8080/">`,
		`hreflang="x-default"`,
		`property="og:image" content="http://localhost:8080/static/og-image.`,
		"See deeper.",
		"Illustrative example. Not real data.",
		`action="/contact#contact"`,
		"Source code",
	} {
		if !strings.Contains(body, want) {
			t.Errorf("index lacks %q", want)
		}
	}
	if n := strings.Count(body, "<h1"); n != 1 {
		t.Errorf("index has %d h1 elements, want 1", n)
	}
	if strings.Contains(body, "http://") && strings.Contains(body, "<script src=\"http") {
		t.Error("external script found")
	}
}

func TestContactDisabledShowsEmailButton(t *testing.T) {
	h := newHarness(t, nil)
	body := h.get(t, "/").Body.String()
	if strings.Contains(body, "<form") {
		t.Error("form rendered although SMTP is not configured")
	}
	if !strings.Contains(body, `href="mailto:hello@tiefer.space">Email us</a>`) {
		t.Error("Email us button missing")
	}
	if !strings.Contains(body, `href="https://github.com/tiefer-labs/web"`) {
		t.Error("the public repository is not linked by default")
	}
	if strings.Contains(newHarness(t, map[string]string{"REPO_URL": ""}).get(t, "/").Body.String(), "Source code") {
		t.Error("source code link rendered although REPO_URL is set to empty")
	}
	w := h.do(t, post("/contact", validValues("x"), ""))
	if w.Code != http.StatusSeeOther || h.mailer.count() != 0 {
		t.Errorf("POST with contact disabled: %d, %d mails", w.Code, h.mailer.count())
	}
}

func TestLegalPlaceholderNote(t *testing.T) {
	note := "Placeholder text. To be reviewed before launch."
	h := newHarness(t, nil)
	for _, p := range []string{"/legal", "/privacy", "/acceptable-use"} {
		body := h.get(t, p).Body.String()
		if !strings.Contains(body, note) {
			t.Errorf("%s: placeholder note missing", p)
		}
		if strings.Count(body, "<h1") != 1 {
			t.Errorf("%s: want exactly one h1", p)
		}
	}
	if !strings.Contains(h.get(t, "/legal").Body.String(), `<mark class="placeholder">[FOUNDER_NAME]</mark>`) {
		t.Error("founder placeholder not rendered")
	}

	h = newHarness(t, map[string]string{"LEGAL_REVIEWED": "true", "LEGAL_FOUNDER": "Aysel Example"})
	body := h.get(t, "/legal").Body.String()
	if strings.Contains(body, note) {
		t.Error("placeholder note shown although LEGAL_REVIEWED=true")
	}
	if !strings.Contains(body, "Aysel Example") {
		t.Error("LEGAL_FOUNDER not rendered")
	}
}

// TestLegalStages checks that the legal pages describe the founder as the
// operator until the company details are set, and the company after.
func TestLegalStages(t *testing.T) {
	h := newHarness(t, map[string]string{"LEGAL_FOUNDER": "Aysel Example"})
	notice := h.get(t, "/legal").Body.String()
	privacy := h.get(t, "/privacy").Body.String()
	for _, want := range []string{
		"is not yet registered as a company",
		`<dt>Name</dt><dd>Aysel Example</dd>`,
		"Aysel Example is responsible for the content of this website.",
	} {
		if !strings.Contains(notice, want) {
			t.Errorf("founding stage: legal notice lacks %q", want)
		}
	}
	if !strings.Contains(privacy, `<dd>Aysel Example, founder of Tiefer</dd>`) {
		t.Error("founding stage: the founder is not named as controller")
	}
	for _, unwanted := range []string{"VÖEN", "Legal form", `id="company"`} {
		if strings.Contains(notice+privacy, unwanted) {
			t.Errorf("founding stage: company detail %q shown", unwanted)
		}
	}

	h = newHarness(t, map[string]string{
		"LEGAL_FOUNDER":      "Aysel Example",
		"LEGAL_NAME":         "Tiefer MMC",
		"LEGAL_FORM":         "Limited liability company",
		"LEGAL_ADDRESS":      "Baku",
		"LEGAL_TAX_ID":       "1234567890",
		"LEGAL_REGISTRATION": "Registered in Baku",
		"LEGAL_DIRECTOR":     "Aysel Example",
	})
	notice = h.get(t, "/legal").Body.String()
	privacy = h.get(t, "/privacy").Body.String()
	for _, want := range []string{"1234567890", "Tiefer MMC is responsible for the content of this website."} {
		if !strings.Contains(notice, want) {
			t.Errorf("company stage: legal notice lacks %q", want)
		}
	}
	if !strings.Contains(privacy, `<dd>Tiefer MMC</dd>`) {
		t.Error("company stage: the company is not named as controller")
	}
	if strings.Contains(notice+privacy, "not yet registered") {
		t.Error("company stage: founding text shown")
	}
	if strings.Count(privacy, `id="who-is-responsible"`) != 1 {
		t.Error("company stage: controller section must appear once")
	}
}

func TestFAQ(t *testing.T) {
	h := newHarness(t, nil)
	body := h.get(t, "/").Body.String()
	for _, want := range []string{
		`<section class="section faq" id="faq"`,
		`<summary><h3>Does Tiefer operate its own satellites?</h3>`,
		`<a class="text-link" href="/acceptable-use">`,
		`<li><a href="#faq">FAQ</a></li>`,
	} {
		if !strings.Contains(body, want) {
			t.Errorf("index lacks %q", want)
		}
	}
}

func TestLegalContents(t *testing.T) {
	h := newHarness(t, nil)
	body := h.get(t, "/privacy").Body.String()
	for _, want := range []string{
		`<nav class="legal-toc" aria-labelledby="legal-toc-title">`,
		`<a href="#who-is-responsible">Who is responsible</a>`,
		`<section class="legal-section" id="who-is-responsible">`,
		`<a href="#complaints">Complaints</a>`,
	} {
		if !strings.Contains(body, want) {
			t.Errorf("privacy page lacks %q", want)
		}
	}
}

func TestContactValidSubmission(t *testing.T) {
	h := newHarness(t, smtpEnv)
	tok := h.formToken(t)
	h.advance(10 * time.Second)
	w := h.do(t, post("/contact", validValues(tok), "192.0.2.1:1234"))
	if w.Code != http.StatusSeeOther || w.Header().Get("Location") != "/?sent=1#contact" {
		t.Fatalf("got %d %q, want 303 to /?sent=1#contact", w.Code, w.Header().Get("Location"))
	}
	if h.mailer.count() != 1 {
		t.Fatalf("%d mails sent, want 1", h.mailer.count())
	}
	msg := h.mailer.msgs[0]
	for _, want := range []string{"Reply-To: \"Aysel Mammadova\" <aysel@example.org>", "To: team@example.org", "Trader or analyst"} {
		if !strings.Contains(msg, want) {
			t.Errorf("mail lacks %q", want)
		}
	}
	body := h.get(t, "/?sent=1").Body.String()
	if !strings.Contains(body, "Thank you. We will reply within two working days.") || strings.Contains(body, "<form") {
		t.Error("success page wrong")
	}

	// The same form posted again is not sent twice.
	w = h.do(t, post("/contact", validValues(tok), "192.0.2.1:1234"))
	if w.Code != http.StatusSeeOther || h.mailer.count() != 1 {
		t.Errorf("replay: %d, %d mails", w.Code, h.mailer.count())
	}
}

func TestContactJSON(t *testing.T) {
	h := newHarness(t, smtpEnv)
	tok := h.formToken(t)
	h.advance(10 * time.Second)

	v := validValues(tok)
	v.Del("name")
	r := post("/contact", v, "")
	r.Header.Set("Accept", "application/json")
	w := h.do(t, r)
	var reply contactReply
	if err := json.Unmarshal(w.Body.Bytes(), &reply); err != nil {
		t.Fatalf("not JSON: %v", err)
	}
	if w.Code != http.StatusUnprocessableEntity || reply.Status != "invalid" || reply.Errors["name"] != "Please enter your name." {
		t.Fatalf("got %d %+v", w.Code, reply)
	}

	r = post("/contact", validValues(tok), "")
	r.Header.Set("Accept", "application/json")
	w = h.do(t, r)
	if err := json.Unmarshal(w.Body.Bytes(), &reply); err != nil || reply.Status != "sent" || w.Code != 200 {
		t.Fatalf("got %d %s", w.Code, w.Body.String())
	}
}

func TestContactRejections(t *testing.T) {
	cases := []struct {
		name   string
		edit   func(url.Values)
		status int
		errMsg string
	}{
		{"missing fields", func(v url.Values) { v.Del("name"); v.Del("email"); v.Del("question"); v.Del("consent") }, 422, "Please enter your name."},
		{"header injection in name", func(v url.Values) { v.Set("name", "Eve\r\nBcc: x@example.com") }, 422, "characters that are not allowed"},
		{"header injection in email", func(v url.Values) { v.Set("email", "eve@example.com\nBcc: x@example.com") }, 422, "characters that are not allowed"},
		{"invalid email", func(v url.Values) { v.Set("email", "not-an-email") }, 422, "valid email address"},
		{"question too long", func(v url.Values) { v.Set("question", strings.Repeat("a", 2001)) }, 422, "2,000 characters or fewer"},
		{"tampered token", func(v url.Values) { v.Set("t", flip(v.Get("t"), 15)) }, 400, "Something went wrong"},
		{"missing token", func(v url.Values) { v.Del("t") }, 400, "Something went wrong"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			h := newHarness(t, smtpEnv)
			tok := h.formToken(t)
			h.advance(10 * time.Second)
			v := validValues(tok)
			c.edit(v)
			w := h.do(t, post("/contact", v, ""))
			if w.Code != c.status {
				t.Errorf("status = %d, want %d", w.Code, c.status)
			}
			body := w.Body.String()
			if !strings.Contains(body, c.errMsg) {
				t.Errorf("response lacks %q", c.errMsg)
			}
			if c.status == 422 && !strings.Contains(body, `value="Example Trading"`) && c.name != "missing fields" {
				t.Error("entered values not kept")
			}
			if h.mailer.count() != 0 {
				t.Error("mail sent for a rejected submission")
			}
		})
	}
}

// flip changes the character at i of a base64url string.
func flip(s string, i int) string {
	b := []byte(s)
	if b[i] == 'A' {
		b[i] = 'B'
	} else {
		b[i] = 'A'
	}
	return string(b)
}

func TestContactRejectsOtherEncodings(t *testing.T) {
	h := newHarness(t, smtpEnv)
	tok := h.formToken(t)
	h.advance(10 * time.Second)
	for _, ct := range []string{"multipart/form-data; boundary=x", "application/json", "text/plain", ""} {
		r := post("/contact", validValues(tok), "")
		r.Header.Set("Content-Type", ct)
		if w := h.do(t, r); w.Code != http.StatusUnsupportedMediaType {
			t.Errorf("Content-Type %q: status %d, want 415", ct, w.Code)
		}
	}
	r := post("/contact", validValues(tok), "")
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded;charset=UTF-8")
	if w := h.do(t, r); w.Code != http.StatusSeeOther || h.mailer.count() != 1 {
		t.Errorf("URL-encoded with charset: status %d, %d mails", w.Code, h.mailer.count())
	}
}

func TestContactHoneypot(t *testing.T) {
	h := newHarness(t, smtpEnv)
	tok := h.formToken(t)
	h.advance(10 * time.Second)
	v := validValues(tok)
	v.Set("website", "http://spam.example")
	w := h.do(t, post("/contact", v, ""))
	if w.Code != http.StatusSeeOther || h.mailer.count() != 0 {
		t.Errorf("honeypot: %d, %d mails; want a quiet 303 and no mail", w.Code, h.mailer.count())
	}
}

func TestContactTooFast(t *testing.T) {
	h := newHarness(t, smtpEnv)
	tok := h.formToken(t)
	h.advance(1 * time.Second)
	w := h.do(t, post("/contact", validValues(tok), ""))
	if w.Code != http.StatusUnprocessableEntity || h.mailer.count() != 0 {
		t.Errorf("too fast: %d, %d mails", w.Code, h.mailer.count())
	}
	if !strings.Contains(w.Body.String(), "Something went wrong. Please email us at") {
		t.Error("error message missing")
	}
}

func TestContactExpired(t *testing.T) {
	h := newHarness(t, smtpEnv)
	tok := h.formToken(t)
	h.advance(13 * time.Hour)
	w := h.do(t, post("/contact", validValues(tok), ""))
	if w.Code != http.StatusUnprocessableEntity || !strings.Contains(w.Body.String(), "This form has expired.") {
		t.Errorf("expired: %d", w.Code)
	}
	if tokenRE.FindStringSubmatch(w.Body.String())[1] == tok {
		t.Error("expired form must get a fresh token")
	}
}

func TestContactRateLimit(t *testing.T) {
	h := newHarness(t, smtpEnv)
	for i := 1; i <= 6; i++ {
		tok := h.formToken(t)
		h.advance(10 * time.Second)
		w := h.do(t, post("/contact", validValues(tok), "198.51.100.7:4000"))
		if i <= 5 && w.Code != http.StatusSeeOther {
			t.Fatalf("submission %d: status %d", i, w.Code)
		}
		if i == 6 && w.Code != http.StatusTooManyRequests {
			t.Fatalf("submission 6: status %d, want 429", w.Code)
		}
	}
	if h.mailer.count() != 5 {
		t.Errorf("%d mails, want 5", h.mailer.count())
	}
	// Another client is not affected, and the limit resets after an hour.
	tok := h.formToken(t)
	h.advance(10 * time.Second)
	if w := h.do(t, post("/contact", validValues(tok), "203.0.113.9:1")); w.Code != http.StatusSeeOther {
		t.Errorf("other client: %d", w.Code)
	}
	h.advance(time.Hour)
	tok = h.formToken(t)
	h.advance(10 * time.Second)
	if w := h.do(t, post("/contact", validValues(tok), "198.51.100.7:4000")); w.Code != http.StatusSeeOther {
		t.Errorf("after an hour: %d", w.Code)
	}
}

func TestContactRateLimitIPv6Prefix(t *testing.T) {
	h := newHarness(t, smtpEnv)
	for i := 1; i <= 6; i++ {
		tok := h.formToken(t)
		h.advance(10 * time.Second)
		// A different address in the same /64 each time.
		w := h.do(t, post("/contact", validValues(tok), "[2001:db8:0:1::"+strconv.Itoa(i)+"]:4000"))
		if i == 6 && w.Code != http.StatusTooManyRequests {
			t.Fatalf("sixth address in the same /64: status %d, want 429", w.Code)
		}
	}
}

func TestContactGlobalSendLimit(t *testing.T) {
	h := newHarness(t, smtpEnv)
	h.srv.sendLimiter.limit = 2
	for i := 1; i <= 3; i++ {
		tok := h.formToken(t)
		h.advance(10 * time.Second)
		w := h.do(t, post("/contact", validValues(tok), "198.51.100."+strconv.Itoa(i)+":1"))
		if i == 3 && w.Code != http.StatusServiceUnavailable {
			t.Fatalf("third message: status %d, want 503", w.Code)
		}
	}
	if h.mailer.count() != 2 {
		t.Errorf("%d mails sent, want 2", h.mailer.count())
	}
}

func TestContactCrossOrigin(t *testing.T) {
	h := newHarness(t, smtpEnv)
	tok := h.formToken(t)
	h.advance(10 * time.Second)
	r := post("/contact", validValues(tok), "")
	r.Header.Set("Sec-Fetch-Site", "cross-site")
	if w := h.do(t, r); w.Code != http.StatusForbidden || h.mailer.count() != 0 {
		t.Errorf("cross-site post: %d, %d mails", w.Code, h.mailer.count())
	}
}

func TestContactDeliveryFailure(t *testing.T) {
	h := newHarness(t, smtpEnv)
	h.mailer.err = io.ErrUnexpectedEOF
	tok := h.formToken(t)
	h.advance(10 * time.Second)
	w := h.do(t, post("/contact", validValues(tok), ""))
	body := w.Body.String()
	if w.Code != http.StatusServiceUnavailable || !strings.Contains(body, `Something went wrong. Please email us at <a href="mailto:hello@tiefer.space">`) {
		t.Errorf("delivery failure: %d", w.Code)
	}
	if !strings.Contains(body, "Which ports saw more vessel activity?") {
		t.Error("question lost after delivery failure")
	}
}

func TestStaticAssets(t *testing.T) {
	h := newHarness(t, nil)
	css := h.srv.assetURL("css/site.css")
	r := httptest.NewRequest(http.MethodGet, css, nil)
	r.Header.Set("Accept-Encoding", "gzip")
	w := h.do(t, r)
	if w.Code != 200 || !strings.Contains(w.Header().Get("Cache-Control"), "immutable") || w.Header().Get("Content-Encoding") != "gzip" {
		t.Fatalf("hashed CSS: %d %v", w.Code, w.Header())
	}
	zr, err := gzip.NewReader(w.Body)
	if err != nil {
		t.Fatal(err)
	}
	b, _ := io.ReadAll(zr)
	if !strings.Contains(string(b), h.srv.assetURL("fonts/MozillaText-VF.woff2")) {
		t.Error("stylesheet does not reference the hashed font URL")
	}

	r = httptest.NewRequest(http.MethodGet, css, nil)
	r.Header.Set("If-None-Match", w.Header().Get("ETag"))
	r.Header.Set("Accept-Encoding", "gzip")
	if w := h.do(t, r); w.Code != http.StatusNotModified {
		t.Errorf("conditional request: %d, want 304", w.Code)
	}
}

func TestSubresourceIntegrity(t *testing.T) {
	h := newHarness(t, nil)
	body := html.UnescapeString(h.get(t, "/").Body.String())
	for _, p := range []string{"css/site.css", "js/site.js"} {
		served := h.get(t, h.srv.assetURL(p)).Body.Bytes()
		sum := sha512.Sum384(served)
		want := `integrity="sha384-` + base64.StdEncoding.EncodeToString(sum[:]) + `"`
		if !strings.Contains(body, want) {
			t.Errorf("%s: page lacks %s", p, want)
		}
	}
}

func TestHTMLIsCompressed(t *testing.T) {
	h := newHarness(t, nil)
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r.Header.Set("Accept-Encoding", "gzip, deflate, br")
	if w := h.do(t, r); w.Header().Get("Content-Encoding") != "gzip" {
		t.Error("HTML not gzip-compressed")
	}
}

func TestSecurityTxt(t *testing.T) {
	h := newHarness(t, map[string]string{
		"SITE_URL":       "https://tiefer.space",
		"SECURITY_EMAIL": "security@tiefer.space",
		"REPO_URL":       "https://github.com/tiefer-labs/web",
	})
	w := h.get(t, "/.well-known/security.txt")
	body := w.Body.String()
	for _, want := range []string{
		"Contact: mailto:security@tiefer.space\n",
		"Expires: 2027-",
		"Canonical: https://tiefer.space/.well-known/security.txt\n",
		"Policy: https://github.com/tiefer-labs/web/security/policy\n",
	} {
		if !strings.Contains(body, want) {
			t.Errorf("security.txt lacks %q:\n%s", want, body)
		}
	}
	if w.Header().Get("Content-Type") != "text/plain; charset=utf-8" {
		t.Errorf("content type %q", w.Header().Get("Content-Type"))
	}
	if w := h.get(t, "/security.txt"); w.Code != http.StatusMovedPermanently || w.Header().Get("Location") != "/.well-known/security.txt" {
		t.Errorf("/security.txt: %d %q", w.Code, w.Header().Get("Location"))
	}
}

func TestSEOFiles(t *testing.T) {
	h := newHarness(t, map[string]string{"SITE_URL": "https://tiefer.example"})
	sitemap := h.get(t, "/sitemap.xml").Body.String()
	for _, p := range []string{"https://tiefer.example/", "https://tiefer.example/legal", "https://tiefer.example/privacy", "https://tiefer.example/acceptable-use"} {
		if !strings.Contains(sitemap, "<loc>"+p+"</loc>") {
			t.Errorf("sitemap lacks %s", p)
		}
	}
	if robots := h.get(t, "/robots.txt").Body.String(); !strings.Contains(robots, "Sitemap: https://tiefer.example/sitemap.xml") {
		t.Errorf("robots.txt: %q", robots)
	}
	var m map[string]any
	w := h.get(t, "/site.webmanifest")
	if err := json.Unmarshal(w.Body.Bytes(), &m); err != nil || m["theme_color"] != "#0C003D" || m["id"] != "/" || m["scope"] != "/" {
		t.Errorf("manifest: %v %v", err, m)
	}
	if w.Header().Get("Content-Type") != "application/manifest+json" {
		t.Errorf("manifest content type %q", w.Header().Get("Content-Type"))
	}
}
