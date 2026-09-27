// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

// Package web embeds the templates and static files into the binary, so
// the server needs no files on disk.
package web

import (
	"embed"
	"io/fs"
)

//go:embed templates
var templates embed.FS

//go:embed static
var static embed.FS

// Templates returns the HTML templates, rooted at web/templates.
func Templates() fs.FS { return sub(templates, "templates") }

// Static returns the static files, rooted at web/static.
func Static() fs.FS { return sub(static, "static") }

func sub(f embed.FS, dir string) fs.FS {
	s, err := fs.Sub(f, dir)
	if err != nil {
		panic(err) // the directory is embedded above, so this cannot happen
	}
	return s
}
