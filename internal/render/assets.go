// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

// Package render loads the embedded templates and static assets, gives
// every asset a content-hashed URL, and renders pages.
package render

import (
	"bytes"
	"compress/gzip"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io/fs"
	"path"
	"regexp"
	"sort"
	"strings"
)

// StaticPrefix is the URL path under which assets are served.
const StaticPrefix = "/static/"

// Asset is one static file, held in memory.
type Asset struct {
	Path        string // relative to web/static, for example css/site.css
	URL         string // content-hashed URL, for example /static/css/site.0123456789.css
	ContentType string
	Body        []byte
	Gzip        []byte // gzip-compressed body, nil when compression does not help
	ETag        string
	Integrity   string // Subresource Integrity value, sha384-...
}

// Assets is the set of static files.
type Assets struct {
	byPath map[string]*Asset
	byURL  map[string]*Asset
}

var contentTypes = map[string]string{
	".css":         "text/css; charset=utf-8",
	".js":          "text/javascript; charset=utf-8",
	".svg":         "image/svg+xml",
	".png":         "image/png",
	".ico":         "image/x-icon",
	".woff2":       "font/woff2",
	".txt":         "text/plain; charset=utf-8",
	".md":          "text/markdown; charset=utf-8",
	".json":        "application/json",
	".webmanifest": "application/manifest+json",
}

var compressible = map[string]bool{".css": true, ".js": true, ".svg": true, ".txt": true, ".md": true, ".json": true, ".ico": true}

// cssURL matches url(...) references in stylesheets.
var cssURL = regexp.MustCompile(`url\(\s*(['"]?)([^'")]+)(['"]?)\s*\)`)

// LoadAssets reads every file of fsys. Stylesheets are processed last:
// their relative url(...) references are rewritten to the hashed URLs of
// the files they point to, so fonts and images are cached forever too.
func LoadAssets(fsys fs.FS) (*Assets, error) {
	a := &Assets{byPath: map[string]*Asset{}, byURL: map[string]*Asset{}}
	var css []string
	err := fs.WalkDir(fsys, ".", func(p string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		if path.Ext(p) == ".css" {
			css = append(css, p)
			return nil
		}
		body, err := fs.ReadFile(fsys, p)
		if err != nil {
			return err
		}
		return a.add(p, body)
	})
	if err != nil {
		return nil, err
	}
	sort.Strings(css)
	for _, p := range css {
		body, err := fs.ReadFile(fsys, p)
		if err != nil {
			return nil, err
		}
		var missing error
		body = cssURL.ReplaceAllFunc(body, func(m []byte) []byte {
			ref := string(cssURL.FindSubmatch(m)[2])
			if strings.Contains(ref, ":") || strings.HasPrefix(ref, "/") || strings.HasPrefix(ref, "#") {
				return m
			}
			target, ok := a.byPath[path.Join(path.Dir(p), ref)]
			if !ok {
				missing = fmt.Errorf("render: %s references missing file %q", p, ref)
				return m
			}
			return []byte(`url("` + target.URL + `")`)
		})
		if missing != nil {
			return nil, missing
		}
		if err := a.add(p, body); err != nil {
			return nil, err
		}
	}
	return a, nil
}

func (a *Assets) add(p string, body []byte) error {
	ext := path.Ext(p)
	ct, ok := contentTypes[ext]
	if !ok {
		return fmt.Errorf("render: no content type for %s", p)
	}
	sri := sha512.Sum384(body)
	sum := sha256.Sum256(body)
	hash := hex.EncodeToString(sum[:])[:10]
	asset := &Asset{
		Path:        p,
		URL:         StaticPrefix + strings.TrimSuffix(p, ext) + "." + hash + ext,
		ContentType: ct,
		Body:        body,
		ETag:        `"` + hash + `"`,
		Integrity:   "sha384-" + base64.StdEncoding.EncodeToString(sri[:]),
	}
	if compressible[ext] {
		var buf bytes.Buffer
		zw, _ := gzip.NewWriterLevel(&buf, gzip.BestCompression)
		_, _ = zw.Write(body)
		_ = zw.Close()
		if buf.Len() < len(body)*9/10 {
			asset.Gzip = buf.Bytes()
		}
	}
	a.byPath[p] = asset
	a.byURL[asset.URL] = asset
	return nil
}

// URL returns the hashed URL of the asset at path p (relative to
// web/static). It fails for unknown files, so a typo in a template is
// caught when the page is rendered in tests.
func (a *Assets) URL(p string) (string, error) {
	asset, ok := a.byPath[strings.TrimPrefix(p, "/")]
	if !ok {
		return "", fmt.Errorf("render: unknown asset %q", p)
	}
	return asset.URL, nil
}

// Integrity returns the Subresource Integrity value of the asset at p.
func (a *Assets) Integrity(p string) (string, error) {
	asset, ok := a.byPath[strings.TrimPrefix(p, "/")]
	if !ok {
		return "", fmt.Errorf("render: unknown asset %q", p)
	}
	return asset.Integrity, nil
}

// Get returns the asset at path p (relative to web/static), or nil.
func (a *Assets) Get(p string) *Asset { return a.byPath[p] }

// Lookup finds the asset for a request path under StaticPrefix. Hashed
// URLs are immutable; plain paths (for example /static/icon-512.png) are
// also served, with a short cache lifetime, for links that must stay
// stable.
func (a *Assets) Lookup(urlPath string) (asset *Asset, immutable bool) {
	if asset, ok := a.byURL[urlPath]; ok {
		return asset, true
	}
	if p, ok := strings.CutPrefix(urlPath, StaticPrefix); ok {
		return a.byPath[p], false
	}
	return nil, false
}
