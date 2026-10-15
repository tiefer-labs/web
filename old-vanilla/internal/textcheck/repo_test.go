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

// binary lists extensions that are not text.
var binary = map[string]bool{".png": true, ".ico": true, ".ttf": true, ".woff2": true, ".gz": true, ".zip": true}

// external lists files that are not written for this project and are
// kept as received.
var external = map[string]bool{
	"LICENSE":                true,
	"docs/Tiefer_2026_az.md": true,
	"tools/go.sum":           true,
}

// TestRepository scans every text file of the repository (code, templates,
// styles, scripts, documentation, workflows and infrastructure) for long
// dashes, emoji and decorative symbols.
func TestRepository(t *testing.T) {
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
			if strings.HasPrefix(rel, "testdata/fuzz") || strings.HasSuffix(rel, "/testdata/fuzz") {
				return filepath.SkipDir
			}
			return nil
		}
		if binary[filepath.Ext(path)] || external[rel] || strings.HasPrefix(d.Name(), ".env") && d.Name() != ".env.example" {
			return nil
		}
		b, err := os.ReadFile(path) // #nosec G304 -- walking this repository in a test
		if err != nil {
			return err
		}
		for _, f := range Characters(string(b)) {
			t.Errorf("%s:%s", rel, f)
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}
