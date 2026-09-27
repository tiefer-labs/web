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

// Templates holds one parsed template set per page. Each set contains the
// layout, every partial and one file of pages/.
type Templates struct {
	pages map[string]*template.Template
}

// LoadTemplates parses layout.html, partials/*.html and every page in
// pages/*.html once, at startup.
func LoadTemplates(fsys fs.FS, assets *Assets) (*Templates, error) {
	funcs := template.FuncMap{
		"asset":     assets.URL,
		"integrity": assets.Integrity,
		"inc":       func(i int) int { return i + 1 },
		"pad2":      func(i int) string { return fmt.Sprintf("%02d", i) },
	}
	base, err := template.New("layout.html").Funcs(funcs).ParseFS(fsys, "layout.html", "partials/*.html")
	if err != nil {
		return nil, fmt.Errorf("render: %w", err)
	}
	pages, err := fs.Glob(fsys, "pages/*.html")
	if err != nil {
		return nil, err
	}
	t := &Templates{pages: map[string]*template.Template{}}
	for _, p := range pages {
		set, err := base.Clone()
		if err != nil {
			return nil, err
		}
		if _, err := set.ParseFS(fsys, p); err != nil {
			return nil, fmt.Errorf("render: %w", err)
		}
		t.pages[strings.TrimSuffix(path.Base(p), ".html")] = set
	}
	return t, nil
}

// Render executes the layout for page (the file name in pages/ without
// .html) with data.
func (t *Templates) Render(w io.Writer, page string, data any) error {
	set, ok := t.pages[page]
	if !ok {
		return fmt.Errorf("render: unknown page %q", page)
	}
	return set.ExecuteTemplate(w, "layout", data)
}
