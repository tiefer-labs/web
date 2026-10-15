// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package textcheck

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// mplLine is the line every source file written for this project carries
// in its header, in the comment syntax of its type.
const mplLine = "This Source Code Form is subject to the terms of the Mozilla Public"

// sourceExt lists the file types that must carry the header.
var sourceExt = map[string]bool{
	".go": true, ".html": true, ".css": true, ".js": true, ".py": true,
	".bicep": true, ".bicepparam": true, ".yml": true, ".yaml": true,
}

// TestMPLHeaders checks that every source file starts with the MPL 2.0
// header (within its first five lines) and that the full licence text is
// in LICENSE.
func TestMPLHeaders(t *testing.T) {
	root := filepath.Join("..", "..")
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, _ := filepath.Rel(root, path)
		rel = filepath.ToSlash(rel)
		if d.IsDir() {
			switch rel {
			case ".git", "bin", "dist", "node_modules":
				return filepath.SkipDir
			}
			return nil
		}
		if !sourceExt[filepath.Ext(path)] && d.Name() != "Makefile" && d.Name() != "Dockerfile" {
			return nil
		}
		b, err := os.ReadFile(path) // #nosec G304 -- walking this repository in a test
		if err != nil {
			return err
		}
		head := strings.Join(strings.SplitN(string(b), "\n", 6)[:min(5, strings.Count(string(b), "\n"))], "\n")
		if !strings.Contains(head, mplLine) {
			t.Errorf("%s: missing the MPL 2.0 header", rel)
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	lic, err := os.ReadFile(filepath.Join(root, "LICENSE"))
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"Mozilla Public License Version 2.0", "Exhibit B - \"Incompatible With Secondary Licenses\" Notice"} {
		if !strings.Contains(string(lic), want) {
			t.Errorf("LICENSE does not look like the full MPL 2.0 text (missing %q)", want)
		}
	}
}
