// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package textcheck

import (
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

var (
	actions  = regexp.MustCompile(`(?s)\{\{.*?\}\}`)
	anyTag   = regexp.MustCompile(`(?s)<[^>]*>`)
	newlines = regexp.MustCompile(`[^\n]`)
)

// TestTemplates checks the template sources: forbidden characters
// anywhere, and banned words in the literal text outside template actions
// and tags. Copy belongs in internal/content, but a stray word in a
// template must not slip through.
func TestTemplates(t *testing.T) {
	root := filepath.Join("..", "..", "web", "templates")
	if _, err := os.Stat(root); os.IsNotExist(err) {
		return
	}
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() || filepath.Ext(path) != ".html" {
			return err
		}
		b, err := os.ReadFile(path) // #nosec G304 -- template files of this repository
		if err != nil {
			return err
		}
		rel, _ := filepath.Rel(root, path)
		src := string(b)
		for _, f := range Characters(src) {
			t.Errorf("web/templates/%s:%s", rel, f)
		}
		// Blank out actions and tags but keep line breaks, so the
		// reported lines match the source.
		keep := func(m string) string { return newlines.ReplaceAllString(m, " ") }
		text := anyTag.ReplaceAllStringFunc(actions.ReplaceAllStringFunc(src, keep), keep)
		for _, f := range BannedWords(text) {
			t.Errorf("web/templates/%s:%d: %s", rel, f.Line, f.What)
		}
		if strings.Contains(text, "&mdash;") || strings.Contains(text, "&ndash;") {
			t.Errorf("web/templates/%s: dash entity", rel)
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}
