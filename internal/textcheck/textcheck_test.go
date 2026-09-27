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

// The forbidden characters are written as escapes so that this file
// passes its own repository scan.
func TestCharacters(t *testing.T) {
	cases := []struct {
		in   string
		want int
	}{
		{"plain text, with a hyphen-minus and 2027-2029", 0},
		{"\u00a9 2026 Tiefer", 0},
		{"T\u00fcrkiye and V\u00d6EN", 0},
		{"a \u2014 b", 1},
		{"a \u2013 b", 1},
		{"a \u2015 b", 1},
		{"launch \U0001F680", 1},
		{"star \u2605", 1},
		{"heart \u2764\ufe0f", 2},
		{"\u00a9\ufe0f", 1},
		{"flag \U0001F1E6\U0001F1FF", 2},
		{"keycap 1\ufe0f\u20e3", 2},
	}
	for _, c := range cases {
		if got := Characters(c.in); len(got) != c.want {
			t.Errorf("Characters(%q) = %v, want %d findings", c.in, got, c.want)
		}
	}
	f := Characters("one\ntwo \u2014")
	if len(f) != 1 || f[0].Line != 2 || f[0].Col != 5 {
		t.Errorf("position = %v, want 2:5", f)
	}
}

func TestBannedWords(t *testing.T) {
	bad := []string{"A revolutionary tool", "Cutting-edge", "game-changing", "AI-powered analyst",
		"seamless", "Seamlessly", "unlock value", "We leverage", "empowering teams", "like magic"}
	for _, s := range bad {
		if len(BannedWords(s)) != 1 {
			t.Errorf("BannedWords(%q) found nothing", s)
		}
	}
	good := []string{"An AI analyst for satellite data", "leverets", "unlockable", "imagery", "powered by radar"}
	for _, s := range good {
		if f := BannedWords(s); len(f) != 0 {
			t.Errorf("BannedWords(%q) = %v, want none", s, f)
		}
	}
}

func TestPlaceholders(t *testing.T) {
	got := Placeholders("[COMPANY_LEGAL_NAME] at [VOEN] and [LEGAL_BASIS: ask counsel] [a] [COMPANY_LEGAL_NAME]")
	want := []string{"[COMPANY_LEGAL_NAME]", "[VOEN]", "[LEGAL_BASIS]"}
	if strings.Join(got, ",") != strings.Join(want, ",") {
		t.Errorf("Placeholders = %v, want %v", got, want)
	}
}

func TestVisibleText(t *testing.T) {
	doc := `<html><head><title>T &amp; U</title><meta name="description" content="Desc">` +
		`<style>.magic{}</style><script>var seamless=1</script></head>` +
		`<body><!-- unlock --><img alt="Logo" src="x.svg"><p>Hello <b>world</b></p></body></html>`
	var texts []string
	for _, s := range VisibleText(doc) {
		texts = append(texts, s.Text)
	}
	got := strings.Join(texts, "|")
	if got != "T & U|Desc|Logo|Hello|world" {
		t.Errorf("VisibleText = %q", got)
	}
}

// TestTextCheckRepository scans every text file in the repository for em
// and en dashes, the horizontal bar and emoji, and reports file and line.
func TestTextCheckRepository(t *testing.T) {
	root := moduleRoot(t)
	binary := map[string]bool{".png": true, ".ico": true, ".woff2": true, ".woff": true, ".ttf": true, ".gz": true}
	checked := 0
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		name := d.Name()
		if d.IsDir() {
			if path != root && (strings.HasPrefix(name, ".") && name != ".github" || name == "bin" || name == "dist") {
				return filepath.SkipDir
			}
			return nil
		}
		if binary[strings.ToLower(filepath.Ext(name))] || name == "tiefer-web" {
			return nil
		}
		b, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		checked++
		rel, _ := filepath.Rel(root, path)
		for _, f := range Characters(string(b)) {
			t.Errorf("%s:%s", rel, f)
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if checked < 20 {
		t.Fatalf("only %d files checked; is the module root right?", checked)
	}
}

func moduleRoot(t *testing.T) string {
	t.Helper()
	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			t.Fatal("go.mod not found")
		}
		dir = parent
	}
}
