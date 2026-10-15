// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package render

import (
	"bytes"
	"compress/gzip"
	"io"
	"strings"
	"testing"
	"testing/fstest"
)

func TestAssets(t *testing.T) {
	fsys := fstest.MapFS{
		"fonts/a.woff2": {Data: []byte("font")},
		"css/site.css":  {Data: []byte(`@font-face{src:url(../fonts/a.woff2) format("woff2")}` + strings.Repeat("a{color:red}", 50))},
		"js/site.js":    {Data: []byte("x")},
		"logo.svg":      {Data: []byte(`<svg xmlns="http://www.w3.org/2000/svg"></svg>`)},
	}
	a, err := LoadAssets(fsys)
	if err != nil {
		t.Fatal(err)
	}
	font, _ := a.ByPath("fonts/a.woff2")
	css, _ := a.ByPath("css/site.css")
	if !strings.HasPrefix(font.URL, "/static/fonts/a.") || !strings.HasSuffix(font.URL, ".woff2") || font.ContentType != "font/woff2" {
		t.Errorf("font asset: %+v", font)
	}
	if !bytes.Contains(css.Body, []byte("url("+font.URL+")")) {
		t.Errorf("css url() not rewritten: %s", css.Body[:80])
	}
	if css.ContentType != "text/css; charset=utf-8" || css.Gzip == nil {
		t.Errorf("css type %q, gzip %v", css.ContentType, css.Gzip != nil)
	}
	zr, err := gzip.NewReader(bytes.NewReader(css.Gzip))
	if err != nil {
		t.Fatal(err)
	}
	plain, _ := io.ReadAll(zr)
	if !bytes.Equal(plain, css.Body) {
		t.Error("gzip body differs from the plain body")
	}
	if js, _ := a.ByPath("js/site.js"); js.Gzip != nil {
		t.Error("tiny files must not be gzipped when it does not pay")
	}
	if got, ok := a.ByURL(font.URL); !ok || got != font {
		t.Error("ByURL")
	}
	if _, err := a.URL("missing.css"); err == nil {
		t.Error("an unknown asset must be an error")
	}

	// Changing the font changes the stylesheet's URL too.
	fsys["fonts/a.woff2"] = &fstest.MapFile{Data: []byte("font v2")}
	b, _ := LoadAssets(fsys)
	css2, _ := b.ByPath("css/site.css")
	if css2.URL == css.URL {
		t.Error("stylesheet hash must follow the fonts it references")
	}
}

func TestAssetsRejectForeignURLs(t *testing.T) {
	for _, css := range []string{`a{background:url(https://example.org/x.png)}`, `a{background:url(data:image/png;base64,AA)}`, `a{background:url(missing.png)}`} {
		if _, err := LoadAssets(fstest.MapFS{"site.css": {Data: []byte(css)}}); err == nil {
			t.Errorf("LoadAssets accepted %s", css)
		}
	}
}
