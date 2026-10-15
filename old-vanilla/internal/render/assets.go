// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

// Package render loads the embedded templates and static assets once at
// startup: it hashes asset names, precompresses text assets and parses
// every page template.
package render

import (
	"bytes"
	"compress/gzip"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/fs"
	"mime"
	"path"
	"regexp"
	"sort"
	"strings"
)

// StaticPrefix is the URL path under which assets are served.
const StaticPrefix = "/static/"

// Asset is one embedded static file, ready to serve.
type Asset struct {
	Path        string // path inside the static tree, for example css/site.css
	URL         string // hashed URL, for example /static/css/site.3f2a9c1d.css
	ContentType string
	Body        []byte
	Gzip        []byte // precompressed body, nil when compression does not pay
	ETag        string // strong ETag from the content hash
}

// Assets holds every static file by path and by hashed URL.
type Assets struct {
	byPath map[string]*Asset
	byURL  map[string]*Asset
}

// compressible lists the content types worth compressing.
var compressible = map[string]bool{
	".css": true, ".js": true, ".svg": true, ".txt": true, ".json": true, ".xml": true, ".webmanifest": true,
}

// cssURL matches url() references in a stylesheet.
var cssURL = regexp.MustCompile(`url\(\s*["']?([^"')]+)["']?\s*\)`)

// LoadAssets reads every file of fsys (the static tree). Stylesheets are
// loaded last, with their url() references rewritten to hashed URLs, so a
// changed font also changes the stylesheet's hash.
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
		b, err := fs.ReadFile(fsys, p)
		if err != nil {
			return err
		}
		return a.add(p, b)
	})
	if err != nil {
		return nil, err
	}
	sort.Strings(css)
	for _, p := range css {
		b, err := fs.ReadFile(fsys, p)
		if err != nil {
			return nil, err
		}
		var missing error
		b = cssURL.ReplaceAllFunc(b, func(m []byte) []byte {
			ref := string(cssURL.FindSubmatch(m)[1])
			if strings.HasPrefix(ref, "data:") || strings.Contains(ref, "://") || strings.HasPrefix(ref, "#") {
				missing = fmt.Errorf("render: %s: only local url() references are allowed, got %q", p, ref)
				return m
			}
			target, ok := a.byPath[path.Join(path.Dir(p), ref)]
			if !ok {
				missing = fmt.Errorf("render: %s references missing asset %q", p, ref)
				return m
			}
			return []byte("url(" + target.URL + ")")
		})
		if missing != nil {
			return nil, missing
		}
		if err := a.add(p, b); err != nil {
			return nil, err
		}
	}
	return a, nil
}

func (a *Assets) add(p string, b []byte) error {
	sum := sha256.Sum256(b)
	hash := hex.EncodeToString(sum[:])[:12]
	ext := path.Ext(p)
	ctype := mime.TypeByExtension(ext)
	switch ext {
	case ".woff2":
		ctype = "font/woff2"
	case ".ico":
		ctype = "image/x-icon"
	case ".webmanifest":
		ctype = "application/manifest+json"
	}
	if ctype == "" {
		return fmt.Errorf("render: no content type for %s", p)
	}
	if strings.HasPrefix(ctype, "text/") && !strings.Contains(ctype, "charset") {
		ctype += "; charset=utf-8"
	}
	as := &Asset{
		Path:        p,
		URL:         StaticPrefix + strings.TrimSuffix(p, ext) + "." + hash + ext,
		ContentType: ctype,
		Body:        b,
		ETag:        `"` + hash + `"`,
	}
	if compressible[ext] {
		var buf bytes.Buffer
		zw, _ := gzip.NewWriterLevel(&buf, gzip.BestCompression)
		_, _ = zw.Write(b)
		_ = zw.Close()
		if buf.Len() < len(b) {
			as.Gzip = buf.Bytes()
		}
	}
	a.byPath[p] = as
	a.byURL[as.URL] = as
	return nil
}

// ByPath returns the asset at a path of the static tree.
func (a *Assets) ByPath(p string) (*Asset, bool) {
	as, ok := a.byPath[p]
	return as, ok
}

// ByURL returns the asset served at a hashed URL.
func (a *Assets) ByURL(u string) (*Asset, bool) {
	as, ok := a.byURL[u]
	return as, ok
}

// URL returns the hashed URL of the asset at path p.
func (a *Assets) URL(p string) (string, error) {
	as, ok := a.byPath[p]
	if !ok {
		return "", fmt.Errorf("render: unknown asset %q", p)
	}
	return as.URL, nil
}

// All returns every asset, sorted by path.
func (a *Assets) All() []*Asset {
	out := make([]*Asset, 0, len(a.byPath))
	for _, as := range a.byPath {
		out = append(out, as)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Path < out[j].Path })
	return out
}
