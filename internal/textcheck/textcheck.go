// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

// Package textcheck finds characters and words that must not appear on
// the site: em and en dashes, the horizontal bar, emoji and hype words.
// The checks run as Go tests over the content files, the templates, the
// repository and the rendered HTML of every page.
package textcheck

import (
	"fmt"
	"html"
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"
)

// Finding is one problem in a text. Line and Col are 1-based; Offset is
// the byte offset in the checked text.
type Finding struct {
	Line   int
	Col    int
	Offset int
	What   string
}

func (f Finding) String() string { return fmt.Sprintf("%d:%d: %s", f.Line, f.Col, f.What) }

// textSymbols are Extended_Pictographic characters that render as plain
// text symbols unless followed by U+FE0F, which is flagged on its own.
// The footer needs the copyright sign.
var textSymbols = map[rune]bool{
	0x00A9: true, // copyright sign
	0x00AE: true, // registered sign
	0x2122: true, // trade mark sign
}

// ForbiddenRune describes r if it must not appear anywhere, or returns "".
func ForbiddenRune(r rune) string {
	switch {
	case r == 0x2014:
		return "em dash U+2014"
	case r == 0x2013:
		return "en dash U+2013"
	case r == 0x2015:
		return "horizontal bar U+2015"
	case r == 0xFE0F:
		return "emoji variation selector U+FE0F"
	case r == 0x20E3:
		return "combining keycap U+20E3"
	case r >= 0x1F1E6 && r <= 0x1F1FF:
		return fmt.Sprintf("regional indicator (flag emoji) U+%04X", r)
	case textSymbols[r]:
		return ""
	case unicode.Is(extendedPictographic, r):
		return fmt.Sprintf("emoji U+%04X", r)
	}
	return ""
}

// Characters reports every forbidden character in s.
func Characters(s string) []Finding {
	var out []Finding
	line, col := 1, 1
	for i, r := range s {
		if r == utf8.RuneError {
			if _, size := utf8.DecodeRuneInString(s[i:]); size == 1 {
				out = append(out, Finding{line, col, i, "invalid UTF-8"})
			}
		}
		if what := ForbiddenRune(r); what != "" {
			out = append(out, Finding{line, col, i, what})
		}
		if r == '\n' {
			line, col = line+1, 1
		} else {
			col++
		}
	}
	return out
}

// bannedWords matches the hype words of the brand rules as whole words,
// including their common inflections.
var bannedWords = regexp.MustCompile(`(?i)\b(revolutionar(?:y|ies)|cutting[- ]edge|game[- ]chang(?:ing|er|ers)|ai[- ]powered|seamless(?:ly)?|unlock(?:s|ed|ing)?|leverag(?:e|es|ed|ing)|empower(?:s|ed|ing|ment)?|magic(?:al|ally)?)\b`)

// BannedWords reports every banned word in s, which should be visible
// copy only.
func BannedWords(s string) []Finding {
	var out []Finding
	for _, m := range bannedWords.FindAllStringIndex(s, -1) {
		line, col := Position(s, m[0])
		out = append(out, Finding{line, col, m[0], fmt.Sprintf("banned word %q", s[m[0]:m[1]])})
	}
	return out
}

// PlaceholderPattern matches [PLACEHOLDER] tokens, optionally with an
// explanation after a colon: [LEGAL_BASIS: to be decided].
var PlaceholderPattern = regexp.MustCompile(`\[[A-Z][A-Z0-9_]{2,}(?::[^\]]*)?\]`)

// Placeholders returns the distinct [PLACEHOLDER] tokens in s, without
// their explanations.
func Placeholders(s string) []string {
	seen := map[string]bool{}
	var out []string
	for _, m := range PlaceholderPattern.FindAllString(s, -1) {
		if i := strings.IndexByte(m, ':'); i > 0 {
			m = m[:i] + "]"
		}
		if !seen[m] {
			seen[m] = true
			out = append(out, m)
		}
	}
	return out
}

// Position returns the 1-based line and column of byte offset off in s.
func Position(s string, off int) (line, col int) {
	before := s[:off]
	line = strings.Count(before, "\n") + 1
	col = utf8.RuneCountInString(before[strings.LastIndexByte(before, '\n')+1:]) + 1
	return line, col
}

// Segment is a piece of visible text and its byte offset in the HTML.
type Segment struct {
	Offset int
	Text   string
}

var (
	tagRE  = regexp.MustCompile(`(?s)<!--.*?-->|<(script|style)\b[^>]*>.*?</(?:script|style)\s*>|<[^>]*>`)
	attrRE = regexp.MustCompile(`(?i)\s(alt|title|aria-label|placeholder|content|value)\s*=\s*"([^"]*)"`)
)

// VisibleText extracts the text a visitor can see or hear from an HTML
// document: text nodes plus the alt, title, aria-label, placeholder,
// content and value attributes. Scripts, styles and comments are skipped.
// It is meant for HTML this project produces, not arbitrary markup.
func VisibleText(doc string) []Segment {
	var out []Segment
	add := func(off int, raw string) {
		if t := strings.TrimSpace(html.UnescapeString(raw)); t != "" {
			out = append(out, Segment{off, t})
		}
	}
	last := 0
	for _, m := range tagRE.FindAllStringSubmatchIndex(doc, -1) {
		add(last, doc[last:m[0]])
		last = m[1]
		tag := doc[m[0]:m[1]]
		if strings.HasPrefix(tag, "<!--") || m[2] >= 0 {
			continue
		}
		for _, a := range attrRE.FindAllStringSubmatchIndex(tag, -1) {
			add(m[0]+a[4], tag[a[4]:a[5]])
		}
	}
	add(last, doc[last:])
	return out
}
