// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

// Package contact validates contact form submissions, protects the form
// against forgery and spam, and delivers messages by email. Submissions
// are never stored and their content is never logged.
package contact

import (
	"net/mail"
	"net/url"
	"strings"
	"unicode"
	"unicode/utf8"
)

// Field names used in the HTML form.
const (
	FieldName     = "name"
	FieldEmail    = "email"
	FieldOrg      = "org"
	FieldRole     = "role"
	FieldQuestion = "question"
	FieldConsent  = "consent"
	FieldToken    = "t"
	FieldHoneypot = "website"
)

// Length limits, in characters (email in bytes, per RFC 5321).
const (
	MaxName     = 100
	MaxEmail    = 254
	MaxOrg      = 150
	MaxQuestion = 2000
)

// Problem is the reason a field was rejected. The server maps it to a
// localised message.
type Problem string

// Problems a field can have.
const (
	Required     Problem = "required"
	TooLong      Problem = "too_long"
	InvalidEmail Problem = "invalid_email"
	NotAllowed   Problem = "not_allowed" // control characters, line breaks in single-line fields
	UnknownRole  Problem = "unknown_role"
)

// Submission is a validated contact form message.
type Submission struct {
	Name     string
	Email    string
	Org      string
	Role     string // one of the allowed role values, or empty
	Question string
	Consent  bool
}

// Errors maps field names to their problem. It is empty when the
// submission is valid.
type Errors map[string]Problem

// Parse reads and validates a submission. roles lists the allowed values
// of the role field. The returned Submission holds the cleaned values
// even when there are errors, so the form can be shown again.
func Parse(form url.Values, roles []string) (Submission, Errors) {
	errs := Errors{}
	s := Submission{
		Name:     strings.TrimSpace(form.Get(FieldName)),
		Email:    strings.TrimSpace(form.Get(FieldEmail)),
		Org:      strings.TrimSpace(form.Get(FieldOrg)),
		Role:     strings.TrimSpace(form.Get(FieldRole)),
		Question: strings.TrimSpace(normalizeNewlines(form.Get(FieldQuestion))),
		Consent:  form.Get(FieldConsent) != "",
	}

	singleLine(errs, FieldName, s.Name, MaxName, true)
	singleLine(errs, FieldOrg, s.Org, MaxOrg, false)

	switch {
	case s.Email == "":
		errs[FieldEmail] = Required
	case hasControl(s.Email, false):
		errs[FieldEmail] = NotAllowed
	case len(s.Email) > MaxEmail || !ValidEmail(s.Email):
		errs[FieldEmail] = InvalidEmail
	}

	if s.Role != "" {
		known := false
		for _, r := range roles {
			known = known || r == s.Role
		}
		if !known {
			errs[FieldRole] = UnknownRole
		}
	}

	switch {
	case s.Question == "":
		errs[FieldQuestion] = Required
	case !utf8.ValidString(s.Question) || hasControl(s.Question, true):
		errs[FieldQuestion] = NotAllowed
	case utf8.RuneCountInString(s.Question) > MaxQuestion:
		errs[FieldQuestion] = TooLong
	}

	if !s.Consent {
		errs[FieldConsent] = Required
	}
	return s, errs
}

// ValidEmail reports whether s is a plain email address (no display name)
// with a dotted domain.
func ValidEmail(s string) bool {
	a, err := mail.ParseAddress(s)
	if err != nil || a.Address != s || a.Name != "" {
		return false
	}
	at := strings.LastIndexByte(s, '@')
	domain := s[at+1:]
	return at > 0 && strings.Contains(domain, ".") && !strings.HasPrefix(domain, ".") && !strings.HasSuffix(domain, ".")
}

func singleLine(errs Errors, field, v string, limit int, required bool) {
	switch {
	case v == "" && required:
		errs[field] = Required
	case !utf8.ValidString(v) || hasControl(v, false):
		// Also rejects CR and LF, which would allow mail header injection.
		errs[field] = NotAllowed
	case utf8.RuneCountInString(v) > limit:
		errs[field] = TooLong
	}
}

// hasControl reports control and invisible format characters. Tabs and,
// if multiline is set, line feeds are allowed.
func hasControl(s string, multiline bool) bool {
	for _, r := range s {
		if r == '\t' || (multiline && r == '\n') {
			continue
		}
		if unicode.IsControl(r) || r == 0x2028 || r == 0x2029 || unicode.Is(unicode.Bidi_Control, r) {
			return true
		}
	}
	return false
}

func normalizeNewlines(s string) string {
	return strings.ReplaceAll(strings.ReplaceAll(s, "\r\n", "\n"), "\r", "\n")
}
