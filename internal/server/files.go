// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package server

import (
	"encoding/json"
	"encoding/xml"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/tiefer-labs/web/internal/render"
)

// static serves embedded assets. Hashed URLs are cached for a year;
// plain paths for an hour.
func (s *Server) static(w http.ResponseWriter, r *http.Request) {
	a, immutable := s.assets.Lookup(r.URL.Path)
	if a == nil {
		s.notFound(w, r)
		return
	}
	cache := "public, max-age=3600"
	if immutable {
		cache = "public, max-age=31536000, immutable"
	}
	s.serveAsset(w, r, a, cache)
}

// rootAsset serves a file that browsers request at a fixed root path.
func (s *Server) rootAsset(path string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.serveAsset(w, r, s.assets.Get(path), "public, max-age=86400")
	}
}

func (s *Server) serveAsset(w http.ResponseWriter, r *http.Request, a *render.Asset, cache string) {
	h := w.Header()
	h.Set("Content-Type", a.ContentType)
	h.Set("Cache-Control", cache)
	body, etag := a.Body, a.ETag
	if a.Gzip != nil {
		h.Add("Vary", "Accept-Encoding")
		if acceptsGzip(r) {
			body, etag = a.Gzip, strings.TrimSuffix(a.ETag, `"`)+`-gz"`
			h.Set("Content-Encoding", "gzip")
		}
	}
	h.Set("ETag", etag)
	if match := r.Header.Get("If-None-Match"); match != "" && strings.Contains(match, etag) {
		h.Del("Content-Encoding")
		w.WriteHeader(http.StatusNotModified)
		return
	}
	h.Set("Content-Length", strconv.Itoa(len(body)))
	if r.Method != http.MethodHead {
		_, _ = w.Write(body)
	} else {
		w.WriteHeader(http.StatusOK)
	}
}

// securityTxt serves /.well-known/security.txt (RFC 9116), so that people
// who find a vulnerability know where to report it. Expires is kept about
// six months ahead of the time the server started, as the RFC asks for a
// date less than a year away.
func (s *Server) securityTxt(w http.ResponseWriter, r *http.Request) {
	var b strings.Builder
	b.WriteString("Contact: mailto:" + s.cfg.SecurityEmail + "\n")
	b.WriteString("Expires: " + s.started.AddDate(0, 6, 0).UTC().Format(time.RFC3339) + "\n")
	b.WriteString("Preferred-Languages: en, az, tr, ru\n")
	b.WriteString("Canonical: " + s.abs("/.well-known/security.txt") + "\n")
	if s.cfg.RepoURL != "" && strings.HasPrefix(s.cfg.RepoURL, "https://github.com/") {
		b.WriteString("Policy: " + strings.TrimSuffix(s.cfg.RepoURL, "/") + "/security/policy\n")
	}
	w.Header().Set("Cache-Control", "public, max-age=86400")
	writeBody(w, r, http.StatusOK, "text/plain; charset=utf-8", []byte(b.String()))
}

func (s *Server) robots(w http.ResponseWriter, r *http.Request) {
	body := "User-agent: *\nAllow: /\n\nSitemap: " + s.abs("/sitemap.xml") + "\n"
	w.Header().Set("Cache-Control", "public, max-age=3600")
	writeBody(w, r, http.StatusOK, "text/plain; charset=utf-8", []byte(body))
}

type sitemapLink struct {
	Rel      string `xml:"rel,attr"`
	Hreflang string `xml:"hreflang,attr"`
	Href     string `xml:"href,attr"`
}

type sitemapURL struct {
	Loc   string        `xml:"loc"`
	Links []sitemapLink `xml:"xhtml:link"`
}

type sitemapSet struct {
	XMLName xml.Name     `xml:"urlset"`
	NS      string       `xml:"xmlns,attr"`
	XHTML   string       `xml:"xmlns:xhtml,attr"`
	URLs    []sitemapURL `xml:"url"`
}

// sitemap lists every page in every locale, with hreflang alternates,
// generated from the route list and SITE_URL.
func (s *Server) sitemap(w http.ResponseWriter, r *http.Request) {
	set := sitemapSet{NS: "http://www.sitemaps.org/schemas/sitemap/0.9", XHTML: "http://www.w3.org/1999/xhtml"}
	for _, site := range s.locales {
		for _, pg := range Pages {
			p := s.newPage(site, pg.Name, pg.Path)
			u := sitemapURL{Loc: p.Canonical}
			if len(s.locales) > 1 {
				for _, alt := range p.Alternates {
					u.Links = append(u.Links, sitemapLink{"alternate", alt.Lang, alt.Href})
				}
			}
			set.URLs = append(set.URLs, u)
		}
	}
	b, err := xml.MarshalIndent(set, "", "  ")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Cache-Control", "public, max-age=3600")
	writeBody(w, r, http.StatusOK, "application/xml; charset=utf-8", append([]byte(xml.Header), b...))
}

func (s *Server) manifest(w http.ResponseWriter, r *http.Request) {
	site := s.locales[0]
	type icon struct {
		Src   string `json:"src"`
		Sizes string `json:"sizes"`
		Type  string `json:"type"`
	}
	m := struct {
		ID              string `json:"id"`
		Name            string `json:"name"`
		ShortName       string `json:"short_name"`
		Description     string `json:"description"`
		Lang            string `json:"lang"`
		StartURL        string `json:"start_url"`
		Scope           string `json:"scope"`
		Display         string `json:"display"`
		ThemeColor      string `json:"theme_color"`
		BackgroundColor string `json:"background_color"`
		Icons           []icon `json:"icons"`
	}{
		Name:            site.Meta.SiteName,
		ShortName:       site.Meta.SiteName,
		Description:     site.Meta.Description,
		Lang:            site.Locale.Code,
		ID:              "/",
		StartURL:        "/",
		Scope:           "/",
		Display:         "browser",
		ThemeColor:      "#0C003D",
		BackgroundColor: "#FFFFFF",
		Icons: []icon{
			{s.assetURL("icon-192.png"), "192x192", "image/png"},
			{s.assetURL("icon-512.png"), "512x512", "image/png"},
		},
	}
	b, _ := json.MarshalIndent(m, "", "  ")
	w.Header().Set("Cache-Control", "public, max-age=3600")
	writeBody(w, r, http.StatusOK, "application/manifest+json", b)
}
