// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package contact

import (
	"bytes"
	"errors"
	"mime"
	"net/mail"
	"net/netip"
	"net/url"
	"strings"
	"testing"
	"time"
)

type clock struct{ t time.Time }

func (c *clock) now() time.Time { return c.t }

func TestTokens(t *testing.T) {
	c := &clock{time.Date(2026, 10, 1, 9, 0, 0, 0, time.UTC)}
	tk := NewTokens([]byte("0123456789abcdef0123456789abcdef"), c.now)
	tok := tk.Issue()
	if err := tk.Check(tok); !errors.Is(err, ErrTooFast) {
		t.Errorf("fresh token: %v, want ErrTooFast", err)
	}
	c.t = c.t.Add(5 * time.Second)
	if err := tk.Check(tok); err != nil {
		t.Errorf("after 5 s: %v", err)
	}
	if err := tk.Use(tok); err != nil {
		t.Errorf("first use: %v", err)
	}
	if err := tk.Use(tok); !errors.Is(err, ErrTokenExpired) {
		t.Errorf("replay: %v, want ErrTokenExpired", err)
	}
	old := tk.Issue()
	c.t = c.t.Add(3 * time.Hour)
	if err := tk.Check(old); !errors.Is(err, ErrTokenExpired) {
		t.Errorf("3 h old token: %v", err)
	}

	// Tampering, other keys and garbage are invalid.
	fresh := tk.Issue()
	c.t = c.t.Add(10 * time.Second)
	b := []byte(fresh)
	b[len(b)/2] ^= 1
	other := NewTokens([]byte("another-secret-another-secret-xx"), c.now)
	for _, s := range []string{string(b), "", "abc", strings.Repeat("A", len(fresh)), other.Issue()} {
		if err := tk.Check(s); !errors.Is(err, ErrTokenInvalid) && !errors.Is(err, ErrTooFast) {
			t.Errorf("Check(%q) = %v, want invalid", s, err)
		}
	}
	if err := tk.Check(fresh); err != nil {
		t.Errorf("untampered token: %v", err)
	}
}

func TestParseAndValidate(t *testing.T) {
	roles := []string{"operator", "other"}
	v := url.Values{
		"name":    {"  Ada\r\nBcc: x@example.org  "},
		"email":   {" ada@example.org "},
		"org":     {"Orbit\tLab"},
		"role":    {"operator"},
		"message": {"Line one\r\nLine two\x00\x07"},
		"consent": {"yes"},
	}
	s := Parse(v)
	if s.Name != "Ada Bcc: x@example.org" || s.Email != "ada@example.org" || s.Org != "Orbit Lab" {
		t.Errorf("parsed: %+v", s)
	}
	if s.Message != "Line one\nLine two" {
		t.Errorf("message = %q", s.Message)
	}
	if p := s.Validate(roles); len(p) != 0 {
		t.Errorf("valid submission: %v", p)
	}
	bad := Submission{Name: strings.Repeat("a", MaxName+1), Email: "ada@example.org\nBcc: x@y.z", Role: "admin", Message: ""}
	want := map[string]Problem{"name": TooLong, "email": Invalid, "role": Invalid, "message": Missing, "consent": Missing}
	got := bad.Validate(roles)
	for k, v := range want {
		if got[k] != v {
			t.Errorf("%s: %q, want %q", k, got[k], v)
		}
	}
	if p := (Submission{Name: "A", Email: "a@b.de", Message: strings.Repeat("é", MaxMessage), Consent: true}).Validate(roles); len(p) != 0 {
		t.Errorf("2,000 characters (not bytes) must pass: %v", p)
	}
}

func TestBuild(t *testing.T) {
	m := Message{From: "web@tiefer.space", To: "hello@tiefer.space", ReplyTo: "ada@example.org",
		Subject: "Website contact: Ada Łukasiewicz", Body: "Hello\nsecond line with = sign", Date: time.Unix(0, 0), Domain: "tiefer.space"}
	raw, err := Build(m)
	if err != nil {
		t.Fatal(err)
	}
	msg, err := mail.ReadMessage(bytes.NewReader(raw))
	if err != nil {
		t.Fatal(err)
	}
	if msg.Header.Get("Reply-To") != "ada@example.org" || msg.Header.Get("From") != "web@tiefer.space" {
		t.Errorf("headers: %v", msg.Header)
	}
	dec := new(mime.WordDecoder)
	if subj, _ := dec.DecodeHeader(msg.Header.Get("Subject")); subj != m.Subject {
		t.Errorf("subject = %q", subj)
	}
	for _, v := range []string{"a\r\nBcc: x@y.z", "a\nb", "a\x00"} {
		if _, err := Build(Message{From: "a@b.de", To: "c@d.de", Subject: v, Domain: "x"}); !errors.Is(err, ErrHeader) {
			t.Errorf("Build accepted header %q", v)
		}
	}
}

func TestLimiter(t *testing.T) {
	c := &clock{time.Unix(1000, 0)}
	l := NewLimiter(2, time.Hour, 2, c.now)
	for i, want := range []bool{true, true, false} {
		if got := l.Allow("a"); got != want {
			t.Errorf("event %d: Allow = %v, want %v", i+1, got, want)
		}
	}
	c.t = c.t.Add(61 * time.Minute)
	if !l.Allow("a") {
		t.Error("the window must slide")
	}
	if !l.Allow("b") {
		t.Error("second key")
	}
	if l.Allow("c") {
		t.Error("a full table of active keys must fail closed")
	}
	c.t = c.t.Add(2 * time.Hour)
	if !l.Allow("c") {
		t.Error("expired keys must be pruned")
	}
}

func TestClientKey(t *testing.T) {
	cases := map[string]string{
		"203.0.113.7":             "203.0.113.7",
		"::ffff:203.0.113.7":      "203.0.113.7",
		"2001:db8:1:2:3:4:5:6":    "2001:db8:1:2::/64",
		"2001:db8:1:2:ffff::1234": "2001:db8:1:2::/64",
	}
	for in, want := range cases {
		if got := ClientKey(netip.MustParseAddr(in)); got != want {
			t.Errorf("ClientKey(%s) = %s, want %s", in, got, want)
		}
	}
	if ClientKey(netip.Addr{}) != "unknown" {
		t.Error("invalid address")
	}
}
