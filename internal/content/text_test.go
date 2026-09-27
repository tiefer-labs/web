// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package content

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/tiefer-labs/web/internal/textcheck"
)

// TestTextCheckContent checks every content file for forbidden characters
// (anywhere in the file) and banned words (in string literals, which is
// where the visible copy lives). It reports file and line.
func TestTextCheckContent(t *testing.T) {
	files, err := filepath.Glob("*.go")
	if err != nil {
		t.Fatal(err)
	}
	fset := token.NewFileSet()
	for _, name := range files {
		if strings.HasSuffix(name, "_test.go") {
			continue
		}
		src, err := os.ReadFile(name)
		if err != nil {
			t.Fatal(err)
		}
		for _, f := range textcheck.Characters(string(src)) {
			t.Errorf("internal/content/%s:%s", name, f)
		}
		file, err := parser.ParseFile(fset, name, src, parser.SkipObjectResolution)
		if err != nil {
			t.Fatal(err)
		}
		ast.Inspect(file, func(n ast.Node) bool {
			lit, ok := n.(*ast.BasicLit)
			if !ok || lit.Kind != token.STRING {
				return true
			}
			s, err := strconv.Unquote(lit.Value)
			if err != nil {
				t.Fatalf("%s: %v", fset.Position(lit.Pos()), err)
			}
			pos := fset.Position(lit.Pos())
			// Escaped characters are invisible to the raw scan above.
			if strings.Contains(lit.Value, `\`) {
				for _, f := range textcheck.Characters(s) {
					t.Errorf("internal/content/%s:%d: %s (escaped)", name, pos.Line, f.What)
				}
			}
			for _, f := range textcheck.BannedWords(s) {
				t.Errorf("internal/content/%s:%d: %s", name, pos.Line+f.Line-1, f.What)
			}
			return true
		})
	}
}

// TestContentComplete catches copy that was left empty by mistake.
func TestContentComplete(t *testing.T) {
	for _, s := range Locales {
		if s.Locale.Code == "" || s.Meta.Title == "" || s.Hero.Title == "" || s.Contact.Submit == "" {
			t.Errorf("locale %q: essential copy missing", s.Locale.Code)
		}
		if n := len(s.Problem.Items); n != 4 {
			t.Errorf("locale %q: %d problem items, want 4", s.Locale.Code, n)
		}
		if n := len(s.Product.Layers); n != 4 {
			t.Errorf("locale %q: %d product layers, want 4", s.Locale.Code, n)
		}
		if n := len(s.Demo.Strip.Days); n != 14 {
			t.Errorf("locale %q: %d strip days, want 14", s.Locale.Code, n)
		}
		if s.Locale.Prefix != "" && (!strings.HasPrefix(s.Locale.Prefix, "/") || strings.HasSuffix(s.Locale.Prefix, "/")) {
			t.Errorf("locale %q: prefix %q must look like /xx", s.Locale.Code, s.Locale.Prefix)
		}
	}
	if Default().Locale.Prefix != "" {
		t.Error("the default locale must not have a URL prefix")
	}
}
