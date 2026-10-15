// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

// Package textcheck finds text that breaks the house rules: long dashes,
// emoji and decorative symbols anywhere, and hype words in visible copy.
// Tests in this package and in the content and server packages use it
// on the content files, the templates and every rendered page.
package textcheck

import (
	"fmt"
	"html"
	"regexp"
	"sort"
	"strings"
	"unicode/utf8"
)

// Finding is one problem at a 1-based line and column.
type Finding struct {
	Line, Col int
	What      string
}

func (f Finding) String() string { return fmt.Sprintf("%d:%d: %s", f.Line, f.Col, f.What) }

// textSymbols are Extended_Pictographic characters that default to text
// presentation and are ordinary typographic symbols. They are allowed on
// their own ("© 2026 Tiefer"); followed by U+FE0F they would turn into
// emoji, and U+FE0F is always reported.
var textSymbols = map[rune]bool{
	0x00A9: true, // copyright sign
	0x00AE: true, // registered sign
	0x2122: true, // trade mark sign
}

// decorative lists blocks of symbols that are used as icons or ornaments
// rather than text: arrows, box drawing, geometric shapes, miscellaneous
// symbols, dingbats and supplemental arrows.
var decorative = [][2]rune{
	{0x2190, 0x21FF}, // Arrows
	{0x2500, 0x27BF}, // Box Drawing to Dingbats
	{0x27F0, 0x27FF}, // Supplemental Arrows-A
	{0x2900, 0x297F}, // Supplemental Arrows-B
	{0x2B00, 0x2BFF}, // Miscellaneous Symbols and Arrows
}

func inRanges(r rune, ranges [][2]rune) bool {
	i := sort.Search(len(ranges), func(i int) bool { return ranges[i][1] >= r })
	return i < len(ranges) && ranges[i][0] <= r
}

// Rune describes why r is forbidden, or returns "" if it is allowed.
func Rune(r rune) string {
	switch {
	case r == 0x2014:
		return "em dash (U+2014)"
	case r == 0x2013:
		return "en dash (U+2013)"
	case r == 0x2015:
		return "horizontal bar (U+2015)"
	case r == 0xFE0F:
		return "emoji variation selector (U+FE0F)"
	case r == 0x20E3:
		return "combining keycap (U+20E3)"
	case r >= 0x1F1E6 && r <= 0x1F1FF:
		return fmt.Sprintf("regional indicator (U+%04X)", r)
	case r >= 0x1F3FB && r <= 0x1F3FF:
		return fmt.Sprintf("skin tone modifier (U+%04X)", r)
	case textSymbols[r]:
		return ""
	case inRanges(r, extPict):
		return fmt.Sprintf("emoji (U+%04X)", r)
	case inRanges(r, decorative):
		return fmt.Sprintf("decorative symbol (U+%04X)", r)
	}
	return ""
}

// Characters reports every forbidden character in s, and invalid UTF-8.
func Characters(s string) []Finding {
	var out []Finding
	line, col := 1, 1
	for i, r := range s {
		if r == utf8.RuneError {
			if _, size := utf8.DecodeRuneInString(s[i:]); size == 1 {
				out = append(out, Finding{line, col, "invalid UTF-8"})
			}
		} else if what := Rune(r); what != "" {
			out = append(out, Finding{line, col, what})
		}
		if r == '\n' {
			line, col = line+1, 1
		} else {
			col++
		}
	}
	return out
}

// Banned lists the hype words that must not appear in visible copy.
var Banned = []string{"revolutionary", "cutting-edge", "game-changing", "AI-powered", "seamless", "unlock", "leverage", "empower", "magic"}

// bannedWords matches the banned words as whole words, case-insensitive,
// with a space or hyphen where the word has a hyphen, and with common
// inflections ("seamlessly", "unlocks", "leveraging", "magical").
var bannedWords = regexp.MustCompile(`(?i)\b(revolutionar(?:y|ily)|cutting[- ]edge|game[- ]chang(?:ing|er|ers)|ai[- ]powered|seamless(?:ly)?|unlock(?:s|ed|ing)?|leverag(?:e|es|ed|ing)|empower(?:s|ed|ing|ment)?|magic(?:al|ally)?)\b`)

// BannedWords reports every banned word in s.
func BannedWords(s string) []Finding {
	var out []Finding
	for _, m := range bannedWords.FindAllStringIndex(s, -1) {
		line := strings.Count(s[:m[0]], "\n") + 1
		col := utf8.RuneCountInString(s[strings.LastIndexByte(s[:m[0]], '\n')+1:m[0]]) + 1
		out = append(out, Finding{line, col, fmt.Sprintf("banned word %q", s[m[0]:m[1]])})
	}
	return out
}

// PlaceholderPattern matches placeholders such as [COMPANY_LEGAL_NAME] or
// [LEGAL_BASIS: ask counsel].
var PlaceholderPattern = regexp.MustCompile(`\[[A-Z][A-Z0-9_]{2,}(?::[^\]]*)?\]`)

// Placeholders returns the distinct placeholder names in s, such as
// [LEGAL_BASIS], in order of appearance.
func Placeholders(s string) []string {
	var out []string
	seen := map[string]bool{}
	for _, m := range PlaceholderPattern.FindAllString(s, -1) {
		name, _, _ := strings.Cut(strings.TrimSuffix(m, "]"), ":")
		name += "]"
		if !seen[name] {
			seen[name] = true
			out = append(out, name)
		}
	}
	return out
}

var (
	hiddenBlocks = regexp.MustCompile(`(?is)<(script|style|template)\b[^>]*>.*?</(?:script|style|template)>`)
	comments     = regexp.MustCompile(`(?s)<!--.*?-->`)
	tags         = regexp.MustCompile(`(?s)<[^>]*>`)
	// Attributes whose values people read or hear.
	visibleAttrs = regexp.MustCompile(`(?i)\s(?:alt|title|aria-label|aria-description|placeholder|content)\s*=\s*"([^"]*)"`)
	sectionStart = regexp.MustCompile(`(?i)<(?:section|header|footer|main|article)\b[^>]*\bid="([^"]+)"`)
)

// Segment is a piece of visible text with the section it belongs to.
type Segment struct {
	Section string // id of the enclosing section, or "page"
	Text    string
}

// VisibleText splits rendered HTML into the text a visitor can read or
// hear, grouped by section: text nodes plus alt, title, aria-label,
// placeholder and meta content values. Scripts, styles and comments are
// dropped. It is a check helper for pages this server renders, not a
// general HTML parser.
func VisibleText(doc string) []Segment {
	doc = hiddenBlocks.ReplaceAllString(doc, " ")
	doc = comments.ReplaceAllString(doc, " ")
	starts := sectionStart.FindAllStringSubmatchIndex(doc, -1)
	var out []Segment
	add := func(section, part string) {
		var b strings.Builder
		for _, m := range visibleAttrs.FindAllStringSubmatch(part, -1) {
			b.WriteString(html.UnescapeString(m[1]))
			b.WriteString("\n")
		}
		b.WriteString(html.UnescapeString(tags.ReplaceAllString(part, " ")))
		if t := strings.TrimSpace(b.String()); t != "" {
			out = append(out, Segment{section, t})
		}
	}
	prev, section := 0, "page"
	for _, m := range starts {
		add(section, doc[prev:m[0]])
		prev, section = m[0], doc[m[2]:m[3]]
	}
	add(section, doc[prev:])
	return out
}
