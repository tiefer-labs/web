// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package render

import (
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"path"
	"strings"
)

// Funcs returns the template functions. They are the only functions
// templates may call besides the comparison builtins; a test enforces
// this list.
func Funcs(a *Assets) template.FuncMap {
	return template.FuncMap{
		// asset returns the hashed URL of a static file. An unknown path
		// fails the render, so a typo cannot ship.
		"asset": a.URL,
		// integrity returns the Subresource Integrity value of a file.
		"integrity": a.Integrity,
	}
}

// Templates holds one parsed template set per page.
type Templates struct {
	pages map[string]*template.Template
}

// LoadTemplates parses layout.html and every partial once, then clones
// the set for each file in pages/, which defines the "content" block.
func LoadTemplates(fsys fs.FS, a *Assets) (*Templates, error) {
	base := template.New("layout.html").Funcs(Funcs(a))
	base, err := base.ParseFS(fsys, "layout.html", "partials/*.html")
	if err != nil {
		return nil, fmt.Errorf("render: %w", err)
	}
	files, err := fs.Glob(fsys, "pages/*.html")
	if err != nil {
		return nil, err
	}
	t := &Templates{pages: map[string]*template.Template{}}
	for _, f := range files {
		clone, err := base.Clone()
		if err != nil {
			return nil, err
		}
		if _, err := clone.ParseFS(fsys, f); err != nil {
			return nil, fmt.Errorf("render: %s: %w", f, err)
		}
		t.pages[strings.TrimSuffix(path.Base(f), ".html")] = clone
	}
	return t, nil
}

// Render writes page (a file name in pages/ without .html) with data.
func (t *Templates) Render(w io.Writer, page string, data any) error {
	tpl, ok := t.pages[page]
	if !ok {
		return fmt.Errorf("render: unknown page %q", page)
	}
	return tpl.ExecuteTemplate(w, "layout.html", data)
}

// Pages returns the names of all page templates.
func (t *Templates) Pages() []string {
	out := make([]string, 0, len(t.pages))
	for name := range t.pages {
		out = append(out, name)
	}
	return out
}
