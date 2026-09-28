// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package contact

import (
	"bytes"
	"net/mail"
	"net/url"
	"strings"
	"testing"
	"time"
)

// FuzzToken: no input panics, and no forged token is ever accepted.
func FuzzToken(f *testing.F) {
	tk := NewTokens([]byte("0123456789abcdef0123456789abcdef"), func() time.Time { return time.Unix(2_000_000_000, 0) })
	f.Add(tk.Issue())
	f.Add("")
	f.Add("AAAA")
	f.Fuzz(func(t *testing.T, s string) {
		err := tk.Check(s)
		if err == nil {
			t.Fatalf("accepted a token that was issued at the same second: %q", s)
		}
	})
}

// FuzzSubmission: parsing never panics, and a valid submission never
// carries a line break in a value that ends up in a mail header.
func FuzzSubmission(f *testing.F) {
	f.Add("Ada", "ada@example.org", "Lab", "operator", "Fires", "yes")
	f.Add("A\r\nBcc: x@y.z", "a@b.de\nBcc: c@d.de", "\x00", "admin", "", "no")
	f.Fuzz(func(t *testing.T, name, email, org, role, message, consent string) {
		s := Parse(url.Values{"name": {name}, "email": {email}, "org": {org}, "role": {role}, "message": {message}, "consent": {consent}})
		if len(s.Validate([]string{"operator"})) != 0 {
			return
		}
		for _, v := range []string{s.Name, s.Email, s.Org} {
			if strings.ContainsAny(v, "\r\n\x00") {
				t.Fatalf("valid submission with a line break: %q", v)
			}
		}
		raw, err := Build(Message{From: "web@tiefer.space", To: "hello@tiefer.space", ReplyTo: s.Email,
			Subject: "Website contact: " + s.Name, Body: s.Message, Domain: "tiefer.space"})
		if err != nil {
			t.Fatalf("Build: %v", err)
		}
		m, err := mail.ReadMessage(bytes.NewReader(raw))
		if err != nil {
			t.Fatalf("unparsable message: %v", err)
		}
		if m.Header.Get("Bcc") != "" || m.Header.Get("Cc") != "" || len(m.Header["To"]) != 1 {
			t.Fatalf("injected header: %v", m.Header)
		}
	})
}
