// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package server

import (
	"bytes"
	"html/template"
	"net/http"
	"strings"

	"github.com/tiefer-labs/web/internal/content"
)

// Page names: files in web/templates/pages without .html.
const (
	PageIndex         = "index"
	PageLegal         = "legal"
	PagePrivacy       = "privacy"
	PageAcceptableUse = "acceptable-use"
	PageNotFound      = "404"
)

// page is the data every template receives.
type page struct {
	Site        *content.Site
	Name        string
	Path        string // path without the locale prefix
	Title       string
	Description string
	Canonical   string
	Alternates  []alternate
	OGImage     string
	JSONLD      template.HTML
	OnDark      bool // the header sits on the dark hero
	NoIndex     bool
	Cfg         publicConfig
	Form        *formView
	Legal       *legalView
}

type alternate struct{ Lang, Href string }

// publicConfig is the part of the configuration templates may read.
type publicConfig struct {
	ContactEmail   string
	LinkedInURL    string
	RepoURL        string
	ContactEnabled bool
	LegalReviewed  bool
}

// Href returns a site path with the locale prefix.
func (p *page) Href(path string) string {
	if path == "/" {
		return p.Site.Locale.Prefix + "/"
	}
	return p.Site.Locale.Prefix + path
}

// Anchor turns a fragment link from the content ("#contact") into a link
// that works from every page: the fragment alone on the index page, the
// index path with the fragment elsewhere.
func (p *page) Anchor(href string) string {
	if !strings.HasPrefix(href, "#") {
		return href
	}
	if p.Name == PageIndex {
		return href
	}
	return p.Href("/") + href
}

func (s *Server) newPage(site *content.Site, name, path string) *page {
	p := &page{
		Site:        site,
		Name:        name,
		Path:        path,
		Title:       site.Meta.Title,
		Description: site.Meta.Description,
		OGImage:     s.cfg.Site() + s.assetURL("og-image.png"),
		JSONLD:      s.jsonLD,
		Cfg: publicConfig{
			ContactEmail:   s.cfg.ContactEmail,
			LinkedInURL:    s.cfg.LinkedInURL,
			RepoURL:        s.cfg.RepoURL,
			ContactEnabled: s.cfg.ContactEnabled(),
			LegalReviewed:  s.cfg.Legal.Reviewed,
		},
	}
	p.Canonical = s.cfg.Site() + p.Href(path)
	for _, l := range s.locales {
		lp := &page{Site: l}
		p.Alternates = append(p.Alternates, alternate{l.Locale.Code, s.cfg.Site() + lp.Href(path)})
	}
	if len(s.locales) > 1 {
		p.Alternates = append(p.Alternates, alternate{"x-default", s.cfg.Site() + path})
	}
	return p
}

func (s *Server) assetURL(p string) string {
	u, err := s.assets.URL(p)
	if err != nil {
		panic(err) // embedded at build time; a test renders every page
	}
	return u
}

func (s *Server) index(site *content.Site) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		p := s.newPage(site, PageIndex, "/")
		p.OnDark = true
		p.Form = s.newForm(r)
		if r.URL.Query().Get("sent") == "1" {
			p.Form.Status = statusSent
		}
		s.renderPage(w, r, http.StatusOK, p)
	}
}

// notFoundPage renders the 404 page in the locale of the path.
func (s *Server) notFoundPage(w http.ResponseWriter, r *http.Request) {
	site := s.localeFor(r.URL.Path)
	p := s.newPage(site, PageNotFound, "/")
	p.Title = site.NotFound.Title + " | " + site.Meta.SiteName
	p.NoIndex = true
	p.Alternates = nil
	p.Canonical = ""
	s.renderPage(w, r, http.StatusNotFound, p)
}

// localeFor returns the locale whose prefix starts the path.
func (s *Server) localeFor(path string) *content.Site {
	for _, l := range s.locales {
		if pre := l.Locale.Prefix; pre != "" && (path == pre || strings.HasPrefix(path, pre+"/")) {
			return l
		}
	}
	return s.locales[0]
}

// renderPage renders into a buffer first, so a template error becomes a
// clean 500 instead of half a page.
func (s *Server) renderPage(w http.ResponseWriter, r *http.Request, status int, p *page) {
	var buf bytes.Buffer
	if err := s.templates.Render(&buf, p.Name, p); err != nil {
		s.log.Error("render", "page", p.Name, "err", err)
		plainError(w, http.StatusInternalServerError)
		return
	}
	// Pages carry a single-use form token, so no shared cache may keep
	// them; the browser revalidates.
	if w.Header().Get("Cache-Control") == "" {
		w.Header().Set("Cache-Control", "private, no-cache")
	}
	writeBody(w, r, status, "text/html; charset=utf-8", buf.Bytes())
}
