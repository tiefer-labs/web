// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package contact

import (
	"net/url"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/tiefer-labs/web/internal/config"
)

// Field limits, in characters.
const (
	MaxName    = 100
	MaxEmail   = 254
	MaxOrg     = 150
	MaxMessage = 2000
)

// Submission is one contact form message.
type Submission struct {
	Name    string
	Email   string
	Org     string
	Role    string // one of the role values, or ""
	Message string
	Consent bool
}

// Problem names what is wrong with a field. The server maps it to the
// copy of the current locale.
type Problem string

// Problems a field can have.
const (
	Missing Problem = "missing"
	TooLong Problem = "too-long"
	Invalid Problem = "invalid"
)

// Parse reads a submission from posted form values. Values are trimmed,
// line endings normalised and control characters removed from single-line
// fields.
func Parse(v url.Values) Submission {
	return Submission{
		Name:    singleLine(v.Get("name")),
		Email:   strings.TrimSpace(v.Get("email")),
		Org:     singleLine(v.Get("org")),
		Role:    strings.TrimSpace(v.Get("role")),
		Message: multiLine(v.Get("message")),
		Consent: v.Get("consent") == "yes",
	}
}

// Validate returns the problems of each field, keyed by form field name.
// roles lists the accepted role values.
func (s Submission) Validate(roles []string) map[string]Problem {
	p := map[string]Problem{}
	check := func(field, value string, max int, required bool) {
		switch n := utf8.RuneCountInString(value); {
		case n == 0 && required:
			p[field] = Missing
		case n > max:
			p[field] = TooLong
		case !utf8.ValidString(value):
			p[field] = Invalid
		}
	}
	check("name", s.Name, MaxName, true)
	check("org", s.Org, MaxOrg, false)
	check("message", s.Message, MaxMessage, true)
	switch {
	case s.Email == "":
		p["email"] = Missing
	case len(s.Email) > MaxEmail || !config.ValidEmail(s.Email):
		p["email"] = Invalid
	}
	if s.Role != "" && !contains(roles, s.Role) {
		p["role"] = Invalid
	}
	if !s.Consent {
		p["consent"] = Missing
	}
	return p
}

func contains(list []string, s string) bool {
	for _, v := range list {
		if v == s {
			return true
		}
	}
	return false
}

// singleLine trims a value and turns every control character, including
// CR and LF, into a space, so it can never break a mail header.
func singleLine(s string) string {
	s = strings.Map(func(r rune) rune {
		if unicode.IsControl(r) || r == ' ' || r == ' ' {
			return ' '
		}
		return r
	}, s)
	return strings.Join(strings.Fields(s), " ")
}

// multiLine normalises line endings to LF and removes control characters
// other than LF and tab.
func multiLine(s string) string {
	s = strings.ReplaceAll(s, "\r\n", "\n")
	s = strings.ReplaceAll(s, "\r", "\n")
	s = strings.Map(func(r rune) rune {
		if r == '\n' || r == '\t' {
			return r
		}
		if unicode.IsControl(r) || r == ' ' || r == ' ' {
			return -1
		}
		return r
	}, s)
	return strings.TrimSpace(s)
}
