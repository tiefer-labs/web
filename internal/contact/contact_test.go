// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package contact

import (
	"net/url"
	"strings"
	"testing"
	"time"
)

var roles = []string{"trader-analyst", "insurer", "other"}

func validForm() url.Values {
	return url.Values{
		FieldName:     {"Aysel Mammadova"},
		FieldEmail:    {"aysel@example.org"},
		FieldOrg:      {"Example Trading"},
		FieldRole:     {"trader-analyst"},
		FieldQuestion: {"Which ports saw more vessel activity?\r\nAnd why?"},
		FieldConsent:  {"yes"},
	}
}

func TestParseValid(t *testing.T) {
	s, errs := Parse(validForm(), roles)
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if s.Question != "Which ports saw more vessel activity?\nAnd why?" {
		t.Errorf("newlines not normalised: %q", s.Question)
	}
	if !s.Consent || s.Role != "trader-analyst" {
		t.Errorf("unexpected submission: %+v", s)
	}
}

func TestParseOptionalFieldsMayBeEmpty(t *testing.T) {
	f := validForm()
	f.Del(FieldOrg)
	f.Del(FieldRole)
	if _, errs := Parse(f, roles); len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
}

func TestParseMissingFields(t *testing.T) {
	_, errs := Parse(url.Values{}, roles)
	for _, field := range []string{FieldName, FieldEmail, FieldQuestion, FieldConsent} {
		if errs[field] != Required {
			t.Errorf("%s: got %q, want %q", field, errs[field], Required)
		}
	}
	if _, ok := errs[FieldOrg]; ok {
		t.Error("organisation is optional")
	}
}

func TestParseRejects(t *testing.T) {
	cases := []struct {
		field, value string
		want         Problem
	}{
		{FieldName, "Eve\r\nBcc: victim@example.com", NotAllowed},
		{FieldName, "Eve\nBcc: victim@example.com", NotAllowed},
		{FieldEmail, "eve@example.com\r\nBcc: victim@example.com", NotAllowed},
		{FieldEmail, "eve@example.com\nBcc: victim@example.com", NotAllowed},
		{FieldEmail, "Eve <eve@example.com>", InvalidEmail},
		{FieldEmail, "eve@localhost", InvalidEmail},
		{FieldEmail, "not an email", InvalidEmail},
		{FieldEmail, strings.Repeat("a", 250) + "@example.com", InvalidEmail},
		{FieldOrg, "Org\r\nX-Header: 1", NotAllowed},
		{FieldName, strings.Repeat("n", MaxName+1), TooLong},
		{FieldOrg, strings.Repeat("o", MaxOrg+1), TooLong},
		{FieldQuestion, strings.Repeat("q", MaxQuestion+1), TooLong},
		{FieldQuestion, "null\x00byte", NotAllowed},
		{FieldRole, "admin", UnknownRole},
	}
	for _, c := range cases {
		f := validForm()
		f.Set(c.field, c.value)
		_, errs := Parse(f, roles)
		if errs[c.field] != c.want {
			t.Errorf("%s=%q: got %q, want %q", c.field, c.value, errs[c.field], c.want)
		}
	}
}

func TestParseCountsCharactersNotBytes(t *testing.T) {
	f := validForm()
	f.Set(FieldQuestion, strings.Repeat("ə", MaxQuestion)) // two bytes each
	if _, errs := Parse(f, roles); len(errs) != 0 {
		t.Fatalf("2,000 characters must be accepted: %v", errs)
	}
}

func TestTokenRoundTrip(t *testing.T) {
	now := time.Unix(1_800_000_000, 0)
	clock := func() time.Time { return now }
	tokens := NewTokens([]byte("0123456789abcdef0123456789abcdef"), clock)

	tok := tokens.Issue()
	if _, err := tokens.Check(tok); err != ErrTooFast {
		t.Fatalf("immediate submit: got %v, want ErrTooFast", err)
	}
	now = now.Add(5 * time.Second)
	nonce, err := tokens.Check(tok)
	if err != nil {
		t.Fatalf("valid token rejected: %v", err)
	}
	if !tokens.Consume(nonce) {
		t.Fatal("first use must succeed")
	}
	if tokens.Consume(nonce) {
		t.Fatal("second use must fail")
	}

	now = now.Add(13 * time.Hour)
	if _, err := tokens.Check(tok); err != ErrTokenExpired {
		t.Fatalf("old token: got %v, want ErrTokenExpired", err)
	}
}

func TestTokenTampering(t *testing.T) {
	now := time.Unix(1_800_000_000, 0)
	a := NewTokens([]byte("0123456789abcdef0123456789abcdef"), func() time.Time { return now })
	b := NewTokens([]byte("fedcba9876543210fedcba9876543210"), func() time.Time { return now.Add(time.Minute) })
	tok := a.Issue()
	if _, err := b.Check(tok); err != ErrTokenInvalid {
		t.Errorf("token from another key: got %v, want ErrTokenInvalid", err)
	}
	flipped := []byte(tok)
	if flipped[3] == 'A' {
		flipped[3] = 'B'
	} else {
		flipped[3] = 'A'
	}
	now = now.Add(time.Minute)
	for _, bad := range []string{"", "abc", string(flipped), tok + "x"} {
		if _, err := a.Check(bad); err != ErrTokenInvalid {
			t.Errorf("Check(%q): got %v, want ErrTokenInvalid", bad, err)
		}
	}
}

func TestBuildMessage(t *testing.T) {
	s, _ := Parse(validForm(), roles)
	s.Name = "Aysel Məmmədova"
	msg := string(BuildMessage(Envelope{
		From: "web@tiefer.example", To: "team@tiefer.example", SiteURL: "https://tiefer.example",
		RoleLabel: "Trader or analyst", Now: time.Unix(1_800_000_000, 0),
	}, s))
	head, body, ok := strings.Cut(msg, "\r\n\r\n")
	if !ok {
		t.Fatal("no header/body separator")
	}
	for _, want := range []string{
		"From: \"Tiefer website\" <web@tiefer.example>\r\n",
		"To: team@tiefer.example\r\n",
		"Reply-To: =?utf-8?q?Aysel_M=C9=99mm=C9=99dova?= <aysel@example.org>\r\n",
		"Subject: =?utf-8?q?",
		"Content-Type: text/plain; charset=utf-8\r\n",
	} {
		if !strings.Contains(head, want) {
			t.Errorf("header missing %q in:\n%s", want, head)
		}
	}
	if !strings.Contains(body, "Trader or analyst") || !strings.Contains(body, "vessel activity") {
		t.Errorf("body incomplete:\n%s", body)
	}
}
