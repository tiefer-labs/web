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

// TestContentText checks every content file: forbidden characters
// anywhere in the file, and banned words in every string literal (the
// copy). It reports file and line.
func TestContentText(t *testing.T) {
	files, err := filepath.Glob("*.go")
	if err != nil {
		t.Fatal(err)
	}
	fset := token.NewFileSet()
	for _, name := range files {
		if strings.HasSuffix(name, "_test.go") {
			continue
		}
		src, err := os.ReadFile(name) // #nosec G304 -- files of this package
		if err != nil {
			t.Fatal(err)
		}
		for _, f := range textcheck.Characters(string(src)) {
			t.Errorf("internal/content/%s:%s", name, f)
		}
		file, err := parser.ParseFile(fset, name, src, 0)
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
			// Characters written as escapes in the source are invisible to the raw
			// scan above.
			for _, f := range textcheck.Characters(s) {
				t.Errorf("internal/content/%s:%d: %s (in a string literal)", name, pos.Line, f.What)
			}
			for _, f := range textcheck.BannedWords(s) {
				t.Errorf("internal/content/%s:%d: %s", name, pos.Line+f.Line-1, f.What)
			}
			return true
		})
	}
}
