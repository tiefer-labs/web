// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package textcheck

import (
	"slices"
	"strings"
	"testing"
)

// The forbidden characters are written as escapes so that this file
// passes the repository scan.
func TestCharacters(t *testing.T) {
	cases := []struct {
		in   string
		want int
	}{
		{"plain text, a hyphen-minus and 2027-2029", 0},
		{"\u00a9 2026 Tiefer", 0},
		{"T\u00fcrkiye, V\u00d6EN, \u018f\u0259, \u0411\u0430\u043a\u0443", 0},
		{"48 KB \u00b7 12% \u2022 40.41 N", 0},
		{"a \u2014 b", 1},
		{"a \u2013 b", 1},
		{"a \u2015 b", 1},
		{"launch \U0001F680", 1},
		{"star \u2605", 1},
		{"check \u2714", 1},
		{"arrow \u2192 next", 1},
		{"box \u25A0", 1},
		{"heart \u2764\ufe0f", 2},
		{"\u00a9\ufe0f", 1},
		{"flag \U0001F1E6\U0001F1FF", 2},
		{"keycap 1\ufe0f\u20e3", 2},
		{"thumb \U0001F44D\U0001F3FD", 2},
		{"bad \xff utf8", 1},
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

func TestExtPictTable(t *testing.T) {
	n := 0
	for i, r := range extPict {
		if r[0] > r[1] || (i > 0 && r[0] <= extPict[i-1][1]+1) {
			t.Fatalf("range %d (%X..%X) is not sorted and merged", i, r[0], r[1])
		}
		n += int(r[1]-r[0]) + 1
	}
	if n != 3537 { // Extended_Pictographic in Unicode 15.0
		t.Errorf("table holds %d code points, want 3537", n)
	}
}

func TestBannedWords(t *testing.T) {
	bad := []string{"A revolutionary tool", "Cutting-edge", "cutting edge", "game-changing", "AI-powered analyst",
		"AI powered", "seamless", "Seamlessly", "unlock value", "We leverage", "leveraging", "empowering teams",
		"empowerment", "like magic", "magical"}
	for _, s := range bad {
		if len(BannedWords(s)) != 1 {
			t.Errorf("BannedWords(%q) found %d, want 1", s, len(BannedWords(s)))
		}
	}
	good := []string{"AI software that runs on the satellite", "unlockable", "magicians", "imagery", "powered by the bus"}
	for _, s := range good {
		if f := BannedWords(s); len(f) != 0 {
			t.Errorf("BannedWords(%q) = %v, want none", s, f)
		}
	}
	for _, w := range Banned {
		if len(BannedWords(w)) != 1 {
			t.Errorf("the pattern does not match the listed word %q", w)
		}
	}
	f := BannedWords("ok\nwe unlock it")
	if len(f) != 1 || f[0].Line != 2 || f[0].Col != 4 {
		t.Errorf("position = %v, want 2:4", f)
	}
}

func TestPlaceholders(t *testing.T) {
	got := Placeholders("[COMPANY_LEGAL_NAME] MMC, [VOEN], [LEGAL_BASIS: ask counsel], [a], [X], [VOEN]")
	want := []string{"[COMPANY_LEGAL_NAME]", "[VOEN]", "[LEGAL_BASIS]"}
	if !slices.Equal(got, want) {
		t.Errorf("Placeholders = %v, want %v", got, want)
	}
}

func TestVisibleText(t *testing.T) {
	doc := `<html><head><title>Tiefer | Edge AI</title><meta name="description" content="Kilobyte alerts">
<script type="application/ld+json">{"unlock":"hidden"}</script><style>.magic{}</style></head>
<body><!-- seamless --><section id="hero"><h1>See deeper.</h1><img alt="Orbit &amp; swath" src="x.svg"></section>
<section id="contact"><input placeholder="you@example.org"><p>Talk to us</p></section></body></html>`
	segs := VisibleText(doc)
	all := ""
	for _, s := range segs {
		all += "[" + s.Section + "] " + s.Text + "\n"
	}
	for _, want := range []string{"[page] Kilobyte alerts", "Tiefer | Edge AI", "[hero] Orbit & swath", "See deeper.", "[contact] you@example.org", "Talk to us"} {
		if !strings.Contains(all, want) {
			t.Errorf("visible text lacks %q:\n%s", want, all)
		}
	}
	for _, hidden := range []string{"unlock", "magic", "seamless", "x.svg"} {
		if strings.Contains(all, hidden) {
			t.Errorf("visible text contains hidden %q:\n%s", hidden, all)
		}
	}
}
