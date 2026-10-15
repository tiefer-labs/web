// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

// Package web embeds the templates and static assets, so the binary needs
// no files on disk. The source TTF fonts are not embedded; only the WOFF2
// files built from them are.
package web

import (
	"embed"
	"io/fs"
)

//go:embed templates
var templates embed.FS

//go:embed static/css static/js static/brand static/fonts/*.woff2
//go:embed static/favicon.svg static/favicon.ico static/apple-touch-icon.png
//go:embed static/icon-192.png static/icon-512.png static/og-image.png
var static embed.FS

// Templates returns the template tree (layout.html, partials/, pages/).
func Templates() fs.FS { return sub(templates, "templates") }

// Static returns the static tree (css/, js/, brand/, fonts/ and icons).
func Static() fs.FS { return sub(static, "static") }

func sub(f embed.FS, dir string) fs.FS {
	s, err := fs.Sub(f, dir)
	if err != nil {
		panic(err) // the directory is embedded at build time
	}
	return s
}
