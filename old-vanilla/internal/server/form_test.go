// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package server

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"html"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"net/url"
	"regexp"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/tiefer-labs/web/internal/contact"
)

// fakeMailer records messages instead of sending them.
type fakeMailer struct {
	mu   sync.Mutex
	sent []contact.Message
	err  error
}

func (f *fakeMailer) Send(_ context.Context, m contact.Message) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.err != nil {
		return f.err
	}
	if _, err := contact.Build(m); err != nil {
		return err
	}
	f.sent = append(f.sent, m)
	return nil
}

func (f *fakeMailer) count() int {
	f.mu.Lock()
	defer f.mu.Unlock()
	return len(f.sent)
}

var smtpEnv = map[string]string{
	"SMTP_HOST": "smtp.example.org", "CONTACT_TO": "hello@tiefer.space", "CONTACT_FROM": "website@tiefer.space",
}

func formHarness(t *testing.T) (*harness, *fakeMailer) {
	m := &fakeMailer{}
	return newHarnessMailer(t, smtpEnv, m), m
}

var tokenField = regexp.MustCompile(`name="t" value="([^"]+)"`)

// token loads the index page and returns its form token.
func (h *harness) token(t *testing.T) string {
	t.Helper()
	m := tokenField.FindStringSubmatch(h.get("/").Body.String())
	if m == nil {
		t.Fatal("no form token on the page")
	}
	return html.UnescapeString(m[1])
}

func validForm(tok string) url.Values {
	return url.Values{
		"t": {tok}, "name": {"Ada Lovelace"}, "email": {"ada@example.org"}, "org": {"Orbit Lab"},
		"role": {"operator"}, "message": {"Wildfires first."}, "consent": {"yes"}, "website": {""},
	}
}

func post(v url.Values, json bool) *http.Request {
	r := httptest.NewRequest(http.MethodPost, "/contact", strings.NewReader(v.Encode()))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Set("Sec-Fetch-Site", "same-origin")
	r.RemoteAddr = "198.51.100.23:40000"
	if json {
		r.Header.Set("Accept", "application/json")
	}
	return r
}

func TestContactWithoutJavaScript(t *testing.T) {
	h, m := formHarness(t)
	tok := h.token(t)
	h.advance(10 * time.Second)
	w := h.do(post(validForm(tok), false))
	if w.Code != http.StatusSeeOther || w.Header().Get("Location") != "/?sent=1#contact" {
		t.Fatalf("status %d, Location %q", w.Code, w.Header().Get("Location"))
	}
	if w.Header().Get("Cache-Control") != "no-store" {
		t.Errorf("Cache-Control = %q", w.Header().Get("Cache-Control"))
	}
	if m.count() != 1 {
		t.Fatalf("sent %d messages", m.count())
	}
	msg := m.sent[0]
	if msg.From != "website@tiefer.space" || msg.To != "hello@tiefer.space" || msg.ReplyTo != "ada@example.org" {
		t.Errorf("addresses: %+v", msg)
	}
	for _, want := range []string{"Ada Lovelace", "Orbit Lab", "Satellite operator", "Wildfires first."} {
		if !strings.Contains(msg.Body, want) {
			t.Errorf("body lacks %q", want)
		}
	}
	if !strings.Contains(h.get("/?sent=1").Body.String(), "Thank you. We will reply within two working days.") {
		t.Error("success message missing after the redirect")
	}
	// The token is single use.
	if w := h.do(post(validForm(tok), false)); w.Code != http.StatusForbidden || m.count() != 1 {
		t.Errorf("replayed token: %d, %d messages", w.Code, m.count())
	}
}

func TestContactWithJavaScript(t *testing.T) {
	h, m := formHarness(t)
	tok := h.token(t)
	h.advance(10 * time.Second)
	w := h.do(post(validForm(tok), true))
	var r formReply
	if err := json.Unmarshal(w.Body.Bytes(), &r); err != nil {
		t.Fatal(err)
	}
	if w.Code != http.StatusOK || r.Status != statusSent || r.Message == "" || r.Token == "" || r.Token == tok || m.count() != 1 {
		t.Errorf("JSON reply %d %+v, %d messages", w.Code, r, m.count())
	}

	v := validForm(r.Token)
	v.Set("email", "not-an-address")
	v.Set("consent", "")
	h.advance(10 * time.Second)
	w = h.do(post(v, true))
	r = formReply{}
	_ = json.Unmarshal(w.Body.Bytes(), &r)
	if w.Code != http.StatusUnprocessableEntity || r.Status != statusInvalid || r.Errors["email"] == "" || r.Errors["consent"] == "" {
		t.Errorf("invalid JSON reply %d %+v", w.Code, r)
	}
}

func TestContactValidationPage(t *testing.T) {
	h, m := formHarness(t)
	tok := h.token(t)
	h.advance(10 * time.Second)
	v := validForm(tok)
	v.Set("name", "")
	v.Set("message", strings.Repeat("x", 2001))
	v.Set("org", `<script>alert(1)</script>`)
	w := h.do(post(v, false))
	body := w.Body.String()
	if w.Code != http.StatusUnprocessableEntity || m.count() != 0 {
		t.Fatalf("status %d, %d messages", w.Code, m.count())
	}
	for _, want := range []string{
		`aria-invalid="true" aria-describedby="f-name-err"`,
		`id="f-name-err">Please enter your name.`,
		`f-message-hint f-message-err`,
		`value="ada@example.org"`,
		`&lt;script&gt;alert(1)&lt;/script&gt;`,
		`<meta name="robots" content="noindex">`,
	} {
		if !strings.Contains(body, want) {
			t.Errorf("error page lacks %q", want)
		}
	}
	if strings.Contains(body, "<script>alert(1)") {
		t.Error("submitted markup is not escaped")
	}
	// A validation error does not use up the token.
	v = validForm(tok)
	if w := h.do(post(v, false)); w.Code != http.StatusSeeOther {
		t.Errorf("corrected form with the same token: %d", w.Code)
	}
}

func TestContactSpamDefences(t *testing.T) {
	h, m := formHarness(t)
	tok := h.token(t)
	// Sent faster than a person can type: success answer, nothing sent.
	if w := h.do(post(validForm(tok), false)); w.Code != http.StatusSeeOther || m.count() != 0 {
		t.Errorf("too fast: %d, %d messages", w.Code, m.count())
	}
	h.advance(10 * time.Second)
	v := validForm(tok)
	v.Set("website", "https://spam.example")
	if w := h.do(post(v, false)); w.Code != http.StatusSeeOther || m.count() != 0 {
		t.Errorf("honeypot: %d, %d messages", w.Code, m.count())
	}
	for name, tk := range map[string]string{"missing": "", "forged": strings.Repeat("A", 64), "tampered": tok[:10] + "x" + tok[11:]} {
		v := validForm(tk)
		if w := h.do(post(v, false)); w.Code != http.StatusForbidden || m.count() != 0 {
			t.Errorf("%s token: %d", name, w.Code)
		}
	}
	h.advance(3 * time.Hour)
	if w := h.do(post(validForm(tok), false)); w.Code != http.StatusForbidden {
		t.Errorf("expired token: %d", w.Code)
	}
}

func TestContactCrossOrigin(t *testing.T) {
	h, m := formHarness(t)
	tok := h.token(t)
	h.advance(10 * time.Second)
	r := post(validForm(tok), false)
	r.Header.Set("Sec-Fetch-Site", "cross-site")
	if w := h.do(r); w.Code != http.StatusForbidden || m.count() != 0 {
		t.Errorf("Sec-Fetch-Site cross-site: %d", w.Code)
	}
	r = post(validForm(tok), false)
	r.Header.Del("Sec-Fetch-Site")
	r.Header.Set("Origin", "https://attacker.example")
	if w := h.do(r); w.Code != http.StatusForbidden || m.count() != 0 {
		t.Errorf("foreign Origin: %d", w.Code)
	}
	r = post(validForm(tok), false)
	r.Header.Del("Sec-Fetch-Site")
	r.Header.Set("Origin", "http://localhost:8080")
	r.Host = "origin-app.internal"
	if w := h.do(r); w.Code != http.StatusSeeOther {
		t.Errorf("SITE_URL origin through a proxy host: %d", w.Code)
	}
}

func TestContactRequestShape(t *testing.T) {
	h, _ := formHarness(t)
	r := httptest.NewRequest(http.MethodPost, "/contact", strings.NewReader(`{"name":"x"}`))
	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("Sec-Fetch-Site", "same-origin")
	if w := h.do(r); w.Code != http.StatusUnsupportedMediaType {
		t.Errorf("JSON body: %d, want 415", w.Code)
	}
	r = httptest.NewRequest(http.MethodPost, "/contact", strings.NewReader("message="+strings.Repeat("a", MaxFormBytes)))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Set("Sec-Fetch-Site", "same-origin")
	if w := h.do(r); w.Code != http.StatusRequestEntityTooLarge {
		t.Errorf("oversized body: %d, want 413", w.Code)
	}
	if w := h.get("/contact"); w.Code != http.StatusMethodNotAllowed {
		t.Errorf("GET /contact: %d, want 405", w.Code)
	}
}

func TestContactRateLimits(t *testing.T) {
	h, m := formHarness(t)
	for i := 0; i < MaxSendsPerClient; i++ {
		tok := h.token(t)
		h.advance(5 * time.Second)
		if w := h.do(post(validForm(tok), false)); w.Code != http.StatusSeeOther {
			t.Fatalf("message %d: %d", i+1, w.Code)
		}
	}
	tok := h.token(t)
	h.advance(5 * time.Second)
	if w := h.do(post(validForm(tok), false)); w.Code != http.StatusTooManyRequests || m.count() != MaxSendsPerClient {
		t.Errorf("message %d: %d, %d sent", MaxSendsPerClient+1, w.Code, m.count())
	}
	// Another address in the same IPv6 /64 shares the limit.
	h2, m2 := formHarness(t)
	for i := 0; i <= MaxSendsPerClient; i++ {
		tok := h2.token(t)
		h2.advance(5 * time.Second)
		r := post(validForm(tok), false)
		r.RemoteAddr = "[2001:db8:1:2::" + string(rune('a'+i)) + "]:5000"
		h2.do(r)
	}
	if m2.count() != MaxSendsPerClient {
		t.Errorf("IPv6 /64: %d sent, want %d", m2.count(), MaxSendsPerClient)
	}
	h.advance(time.Hour)
	tok = h.token(t)
	h.advance(5 * time.Second)
	if w := h.do(post(validForm(tok), false)); w.Code != http.StatusSeeOther {
		t.Errorf("after an hour: %d", w.Code)
	}
}

func TestContactDeliveryFailure(t *testing.T) {
	m := &fakeMailer{err: errors.New("dial tcp: connection refused")}
	h := newHarnessMailer(t, smtpEnv, m)
	var buf bytes.Buffer
	h.srv.log = slog.New(slog.NewJSONHandler(&buf, nil))
	tok := h.token(t)
	h.advance(10 * time.Second)
	w := h.do(post(validForm(tok), false))
	body := w.Body.String()
	if w.Code != http.StatusBadGateway || !strings.Contains(body, "Something went wrong. Please email us at") || !strings.Contains(body, `href="mailto:hello@tiefer.space"`) {
		t.Errorf("failure page: %d", w.Code)
	}
	for _, s := range []string{"Ada", "ada@example.org", "Wildfires", "198.51.100.23"} {
		if strings.Contains(buf.String(), s) {
			t.Errorf("log contains %q", s)
		}
	}
}

func TestContactDisabled(t *testing.T) {
	h := newHarness(t, nil)
	body := h.get("/").Body.String()
	if strings.Contains(body, "<form") || !strings.Contains(body, ">Email us</a>") {
		t.Error("without SMTP the form must be replaced by the Email us button")
	}
	for _, want := range []string{"Prefer email?", "Follow Tiefer on LinkedIn", `href="https://www.linkedin.com/company/tiefer/"`} {
		if !strings.Contains(body, want) {
			t.Errorf("contact section lacks %q", want)
		}
	}
	r := post(validForm("x"), false)
	if w := h.do(r); w.Code != http.StatusNotFound {
		t.Errorf("POST /contact without SMTP: %d", w.Code)
	}
}
