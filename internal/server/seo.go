// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package server

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"encoding/xml"
	"html/template"
	"net/http"
	"strings"
)

// manifest serves the web app manifest, built once at startup.
func (s *Server) manifest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "public, max-age=86400")
	writeBody(w, r, http.StatusOK, "application/manifest+json", s.manifestJSON)
}

func (s *Server) buildManifest() error {
	site := s.locales[0]
	m := map[string]any{
		"name":             site.Meta.SiteName,
		"short_name":       site.Meta.SiteName,
		"description":      site.Meta.Description,
		"id":               "/",
		"start_url":        "/",
		"scope":            "/",
		"display":          "browser",
		"lang":             site.Locale.Code,
		"theme_color":      site.Meta.ThemeColor,
		"background_color": "#FFFFFF",
		"icons": []map[string]string{
			{"src": s.assetURL("icon-192.png"), "sizes": "192x192", "type": "image/png"},
			{"src": s.assetURL("icon-512.png"), "sizes": "512x512", "type": "image/png"},
			{"src": s.assetURL("favicon.svg"), "sizes": "any", "type": "image/svg+xml"},
		},
	}
	b, err := json.Marshal(m)
	s.manifestJSON = b
	return err
}

// Pages lists the public pages in sitemap order, as paths without the
// locale prefix.
var Pages = []string{"/", "/legal", "/privacy", "/acceptable-use"}

// buildJSONLD prepares the Organization block and returns the CSP source
// that allows exactly this script, by its SHA-256 hash.
func (s *Server) buildJSONLD() (string, error) {
	site := s.locales[0]
	org := map[string]any{
		"@context": "https://schema.org",
		"@type":    "Organization",
		"name":     site.Meta.SiteName,
		"url":      s.cfg.Site() + "/",
		"logo":     s.cfg.Site() + s.assetURL("icon-512.png"),
		"sameAs":   []string{s.cfg.LinkedInURL},
	}
	b, err := json.Marshal(org) // escapes <, > and & for embedding in HTML
	if err != nil {
		return "", err
	}
	sum := sha256.Sum256(b)
	s.jsonLD = template.HTML(`<script type="application/ld+json">` + string(b) + `</script>`) // #nosec G203 -- json.Marshal output, see TestNoUnsafeConversions
	return " 'sha256-" + base64.StdEncoding.EncodeToString(sum[:]) + "'", nil
}

// buildSitemap lists every page of every locale.
func (s *Server) buildSitemap() []byte {
	var b strings.Builder
	b.WriteString(xml.Header)
	b.WriteString(`<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">` + "\n")
	for _, site := range s.locales {
		p := &page{Site: site}
		for _, path := range Pages {
			b.WriteString("  <url><loc>")
			_ = xml.EscapeText(&b, []byte(s.cfg.Site()+p.Href(path)))
			b.WriteString("</loc></url>\n")
		}
	}
	b.WriteString("</urlset>\n")
	return []byte(b.String())
}

func (s *Server) buildRobots() []byte {
	if !s.cfg.Production() {
		// Keep development and staging copies out of search engines.
		return []byte("User-agent: *\nDisallow: /\n")
	}
	return []byte("User-agent: *\nAllow: /\n\nSitemap: " + s.cfg.Site() + "/sitemap.xml\n")
}

func (s *Server) sitemap(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "public, max-age=3600")
	writeBody(w, r, http.StatusOK, "application/xml; charset=utf-8", s.sitemapXML)
}

func (s *Server) robots(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "public, max-age=3600")
	writeBody(w, r, http.StatusOK, "text/plain; charset=utf-8", s.robotsTxt)
}
