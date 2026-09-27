// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package server

import (
	"bytes"
	"html/template"
	"net/http"
	"strings"

	"github.com/tiefer-labs/web/internal/contact"
	"github.com/tiefer-labs/web/internal/content"
	"github.com/tiefer-labs/web/internal/textcheck"
)

// page is the data every template receives.
type page struct {
	Site        *content.Site
	Name        string // template name in web/templates/pages
	Path        string // path relative to the locale prefix
	Title       string
	Description string
	Canonical   string
	Alternates  []alternate
	OGImage     string
	JSONLD      template.HTML
	HeaderDark  bool
	NoIndex     bool
	Cfg         publicConfig
	Legal       *content.LegalPage
	Form        *formState

	tokens map[string]string
}

type alternate struct{ Lang, Href string }

// publicConfig is the part of the configuration templates may use.
type publicConfig struct {
	SiteURL        string
	ContactEmail   string
	LinkedInURL    string
	RepoURL        string
	ContactEnabled bool
	LegalReviewed  bool
}

// formState is the contact form as rendered: status, values and errors.
type formState struct {
	Status string // "", "sent", "invalid", "expired" or "error"
	Token  string
	Values contact.Submission
	Errors map[string]string // field name to message
}

// Err returns the error message for a field, or "".
func (f *formState) Err(field string) string { return f.Errors[field] }

// Href returns a site path with the locale prefix, for example /privacy
// or /az/privacy.
func (p *page) Href(path string) string {
	if path == "/" {
		return p.Site.Locale.Prefix + "/"
	}
	return p.Site.Locale.Prefix + path
}

// Anchor links to a section of the index page: a plain fragment on the
// index page itself, the full path elsewhere. "top" on other pages links
// to the start page.
func (p *page) Anchor(id string) string {
	switch {
	case p.Name == PageIndex:
		return "#" + id
	case id == "top":
		return p.Href("/")
	}
	return p.Href("/") + "#" + id
}

// Fill escapes s, replaces {tokens} with configuration values and
// highlights [PLACEHOLDER] tokens, for the legal pages.
func (p *page) Fill(s string) template.HTML {
	out := template.HTMLEscapeString(s)
	for key, value := range p.tokens {
		v := template.HTMLEscapeString(value)
		switch key {
		case "contact_email":
			v = `<a href="mailto:` + v + `">` + v + `</a>`
		case "site_url":
			v = `<a href="` + v + `/">` + v + `</a>`
		}
		out = strings.ReplaceAll(out, "{"+key+"}", v)
	}
	out = textcheck.PlaceholderPattern.ReplaceAllString(out, `<mark class="placeholder">$0</mark>`)
	return template.HTML(out)
}

func (s *Server) newPage(site *content.Site, name, path string) *page {
	p := &page{
		Site:        site,
		Name:        name,
		Path:        path,
		Title:       site.Meta.Title,
		Description: site.Meta.Description,
		OGImage:     s.abs(s.assetURL("og-image.png")),
		JSONLD:      s.jsonLD,
		Cfg: publicConfig{
			SiteURL:        s.cfg.SiteURL,
			ContactEmail:   s.cfg.ContactEmail,
			LinkedInURL:    s.cfg.LinkedInURL,
			RepoURL:        s.cfg.RepoURL,
			ContactEnabled: s.cfg.ContactEnabled(),
			LegalReviewed:  s.cfg.Legal.Reviewed,
		},
	}
	p.Canonical = s.abs(p.Href(path))
	for _, l := range s.locales {
		lp := &page{Site: l}
		p.Alternates = append(p.Alternates, alternate{l.Locale.Code, s.abs(lp.Href(path))})
	}
	p.Alternates = append(p.Alternates, alternate{"x-default", s.abs(path)})
	return p
}

func (s *Server) legalTokens() map[string]string {
	l := s.cfg.Legal
	return map[string]string{
		"legal_name":       l.Name,
		"legal_form":       l.Form,
		"address":          l.Address,
		"tax_id":           l.TaxID,
		"registration":     l.Registration,
		"director":         l.Director,
		"hosting_provider": l.HostingProvider,
		"hosting_country":  l.HostingCountry,
		"smtp_provider":    l.SMTPProvider,
		"smtp_country":     l.SMTPCountry,
		"contact_email":    s.cfg.ContactEmail,
		"site_url":         s.cfg.SiteURL,
	}
}

func (s *Server) index(site *content.Site) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		p := s.newPage(site, PageIndex, "/")
		p.HeaderDark = true
		p.Form = &formState{Token: s.tokens.Issue()}
		if r.URL.Query().Get("sent") == "1" {
			p.Form.Status = "sent"
		}
		s.renderPage(w, r, http.StatusOK, p)
	}
}

func (s *Server) legal(site *content.Site, name string, lp *content.LegalPage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		p := s.newPage(site, name, "/"+name)
		p.Legal = lp
		p.Title = lp.Title + " | " + site.Meta.SiteName
		p.Description = lp.Description
		p.tokens = s.legalTokens()
		s.renderPage(w, r, http.StatusOK, p)
	}
}

func (s *Server) notFound(w http.ResponseWriter, r *http.Request) {
	site := s.localeFor(r.URL.Path)
	p := s.newPage(site, PageNotFound, "/")
	p.Title = site.NotFound.Title + " | " + site.Meta.SiteName
	p.NoIndex = true
	p.Alternates = nil
	s.renderPage(w, r, http.StatusNotFound, p)
}

// renderPage renders into a buffer first, so a template error becomes a
// clean 500 instead of a half-written page.
func (s *Server) renderPage(w http.ResponseWriter, r *http.Request, status int, p *page) {
	var buf bytes.Buffer
	if err := s.templates.Render(&buf, p.Name, p); err != nil {
		s.log.Error("render failed", "page", p.Name, "err", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	// Pages carry a single-use form token: a shared cache (CDN or proxy)
	// must not hand the same page to several visitors.
	if w.Header().Get("Cache-Control") == "" {
		w.Header().Set("Cache-Control", "private, no-cache")
	}
	writeBody(w, r, status, "text/html; charset=utf-8", buf.Bytes())
}
