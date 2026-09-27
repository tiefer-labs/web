// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package server

import (
	"net/http"
	"net/http/httptest"
	"regexp"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/tiefer-labs/web/internal/config"
	"github.com/tiefer-labs/web/internal/textcheck"
)

// landmark finds elements that name a section of a page.
var (
	landmark = regexp.MustCompile(`<(?:section|header|footer|main|article|nav|figure|form)\b[^>]*>`)
	idAttr   = regexp.MustCompile(`\bid="([^"]+)"`)
	clsAttr  = regexp.MustCompile(`\bclass="([^" ]+)`)
)

// sectionAt returns the id (or first class) of the landmark element that
// opens last before byte offset off.
func sectionAt(doc string, off int) string {
	name := "head"
	for _, m := range landmark.FindAllStringIndex(doc, -1) {
		if m[0] > off {
			break
		}
		tag := doc[m[0]:m[1]]
		if a := idAttr.FindStringSubmatch(tag); a != nil {
			name = a[1]
		} else if a := clsAttr.FindStringSubmatch(tag); a != nil {
			name = a[1]
		}
	}
	return name
}

// TestTextCheckRendered renders every page in every state and fails on em
// and en dashes, the horizontal bar, emoji and banned words, reporting the
// page, section and line. It also lists, as warnings that do not fail the
// test, every [PLACEHOLDER] still on the site and every configuration
// value that still has its default.
func TestTextCheckRendered(t *testing.T) {
	type render struct {
		name string
		req  func(h *harness) *http.Request
	}
	get := func(p string) func(*harness) *http.Request {
		return func(*harness) *http.Request { return httptest.NewRequest(http.MethodGet, p, nil) }
	}
	renders := []render{
		{"/", get("/")},
		{"/?sent=1", get("/?sent=1")},
		{"/legal", get("/legal")},
		{"/privacy", get("/privacy")},
		{"/acceptable-use", get("/acceptable-use")},
		{"404", get("/missing")},
		{"/ (form with errors)", func(h *harness) *http.Request {
			tok := h.formToken(t)
			h.advance(10 * time.Second)
			return post("/contact", map[string][]string{"t": {tok}}, "")
		}},
		{"/ (form expired)", func(h *harness) *http.Request {
			tok := h.formToken(t)
			h.advance(24 * time.Hour)
			return post("/contact", validValues(tok), "")
		}},
		{"/ (form disabled)", nil},
	}

	placeholders := map[string][]string{}
	for _, rd := range renders {
		env := smtpEnv
		req := rd.req
		if req == nil {
			env, req = nil, get("/")
		}
		h := newHarness(t, env)
		w := h.do(t, req(h))
		doc := w.Body.String()
		if !strings.Contains(doc, "</html>") {
			t.Errorf("%s: incomplete page (status %d)", rd.name, w.Code)
			continue
		}
		for _, f := range textcheck.Characters(doc) {
			t.Errorf("page %s, section %s, line %d: %s", rd.name, sectionAt(doc, f.Offset), f.Line, f.What)
		}
		for _, seg := range textcheck.VisibleText(doc) {
			for _, f := range textcheck.BannedWords(seg.Text) {
				line, _ := textcheck.Position(doc, seg.Offset)
				t.Errorf("page %s, section %s, line %d: %s", rd.name, sectionAt(doc, seg.Offset), line+f.Line-1, f.What)
			}
		}
		for _, p := range textcheck.Placeholders(doc) {
			placeholders[p] = append(placeholders[p], rd.name)
		}
	}

	keys := make([]string, 0, len(placeholders))
	for k := range placeholders {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		t.Logf("WARNING: placeholder %s on %s", k, strings.Join(unique(placeholders[k]), ", "))
	}

	// The defaults warning reflects the environment the test runs in.
	cfg, err := config.Load()
	if err != nil {
		t.Logf("WARNING: configuration in this environment is invalid: %v", err)
		return
	}
	for _, name := range cfg.DefaultsInUse() {
		t.Logf("WARNING: %s still uses its default; set it before launch", name)
	}
}

func unique(s []string) []string {
	var out []string
	seen := map[string]bool{}
	for _, v := range s {
		if !seen[v] {
			seen[v] = true
			out = append(out, v)
		}
	}
	return out
}

// TestTextCheckTemplates checks the template sources: forbidden characters
// anywhere and banned words in the text outside template actions.
func TestTextCheckTemplates(t *testing.T) {
	files, err := templateFiles()
	if err != nil {
		t.Fatal(err)
	}
	action := regexp.MustCompile(`(?s)\{\{.*?\}\}`)
	for name, src := range files {
		for _, f := range textcheck.Characters(src) {
			t.Errorf("web/templates/%s:%s", name, f)
		}
		// Blank out actions but keep offsets, so line numbers stay right.
		plain := action.ReplaceAllStringFunc(src, func(m string) string {
			return strings.Map(func(r rune) rune {
				if r == '\n' {
					return r
				}
				return ' '
			}, m)
		})
		for _, seg := range textcheck.VisibleText(plain) {
			for _, f := range textcheck.BannedWords(seg.Text) {
				line, _ := textcheck.Position(plain, seg.Offset)
				t.Errorf("web/templates/%s:%d: %s", name, line+f.Line-1, f.What)
			}
		}
	}
}
