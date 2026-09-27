// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package server

import (
	"io/fs"

	"github.com/tiefer-labs/web/web"
)

// templateFiles returns the source of every embedded template by path.
func templateFiles() (map[string]string, error) {
	out := map[string]string{}
	err := fs.WalkDir(web.Templates(), ".", func(p string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		b, err := fs.ReadFile(web.Templates(), p)
		out[p] = string(b)
		return err
	})
	return out, err
}
