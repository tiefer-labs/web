// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

// Package server wires the routes, handlers and middleware of the
// website.
package server

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"html/template"
	"io/fs"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/tiefer-labs/web/internal/config"
	"github.com/tiefer-labs/web/internal/contact"
	"github.com/tiefer-labs/web/internal/content"
	"github.com/tiefer-labs/web/internal/render"
)

// Page names, matching the files in web/templates/pages, and their paths
// relative to the locale prefix.
const (
	PageIndex         = "index"
	PageLegal         = "legal"
	PagePrivacy       = "privacy"
	PageAcceptableUse = "acceptable-use"
	PageNotFound      = "404"
)

// Pages lists every public page in sitemap order.
var Pages = []struct{ Name, Path string }{
	{PageIndex, "/"},
	{PageLegal, "/legal"},
	{PagePrivacy, "/privacy"},
	{PageAcceptableUse, "/acceptable-use"},
}

// Options holds the dependencies of the server.
type Options struct {
	Config    *config.Config
	Logger    *slog.Logger
	Templates fs.FS
	Static    fs.FS
	Locales   []*content.Site // defaults to content.Locales
	Mailer    contact.Mailer  // defaults to SMTP from the config
	Now       func() time.Time
}

// Server serves the website.
type Server struct {
	cfg       *config.Config
	log       *slog.Logger
	assets    *render.Assets
	templates *render.Templates
	locales   []*content.Site
	mailer    contact.Mailer
	tokens    *contact.Tokens
	limiter   *rateLimiter
	now       func() time.Time
	jsonLD    template.HTML
	csp       string
	handler   http.Handler
}

// New builds the server: it loads assets and templates, prepares the
// JSON-LD block and its CSP hash, and registers every route.
func New(o Options) (*Server, error) {
	if o.Now == nil {
		o.Now = time.Now
	}
	if o.Locales == nil {
		o.Locales = content.Locales
	}
	if o.Logger == nil {
		o.Logger = slog.Default()
	}
	assets, err := render.LoadAssets(o.Static)
	if err != nil {
		return nil, err
	}
	templates, err := render.LoadTemplates(o.Templates, assets)
	if err != nil {
		return nil, err
	}
	s := &Server{
		cfg:       o.Config,
		log:       o.Logger,
		assets:    assets,
		templates: templates,
		locales:   o.Locales,
		mailer:    o.Mailer,
		tokens:    contact.NewTokens(o.Config.CSRFSecret, o.Now),
		limiter:   newRateLimiter(5, time.Hour, o.Now),
		now:       o.Now,
	}
	if s.mailer == nil && o.Config.ContactEnabled() {
		s.mailer = &contact.SMTPMailer{
			Host:       o.Config.SMTP.Host,
			Port:       o.Config.SMTP.Port,
			User:       o.Config.SMTP.User,
			Pass:       o.Config.SMTP.Pass,
			SkipVerify: o.Config.SMTP.SkipVerify,
			LocalName:  hostOf(o.Config.SiteURL),
		}
	}
	if err := s.buildJSONLD(); err != nil {
		return nil, err
	}
	s.handler = s.routes()
	return s, nil
}

// Handler returns the root HTTP handler with all middleware applied.
func (s *Server) Handler() http.Handler { return s.handler }

func (s *Server) routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.Header().Set("Cache-Control", "no-store")
		_, _ = w.Write([]byte("ok"))
	})
	mux.HandleFunc("GET /robots.txt", s.robots)
	mux.HandleFunc("GET /sitemap.xml", s.sitemap)
	mux.HandleFunc("GET /site.webmanifest", s.manifest)
	mux.HandleFunc("GET /favicon.ico", s.rootAsset("favicon.ico"))
	mux.HandleFunc("GET /apple-touch-icon.png", s.rootAsset("apple-touch-icon.png"))
	mux.HandleFunc("GET "+render.StaticPrefix+"{path...}", s.static)

	csrf := http.NewCrossOriginProtection()
	for _, site := range s.locales {
		prefix := site.Locale.Prefix
		mux.HandleFunc("GET "+prefix+"/{$}", s.index(site))
		mux.HandleFunc("GET "+prefix+"/legal", s.legal(site, PageLegal, &site.Legal.Notice))
		mux.HandleFunc("GET "+prefix+"/privacy", s.legal(site, PagePrivacy, &site.Legal.Privacy))
		mux.HandleFunc("GET "+prefix+"/acceptable-use", s.legal(site, PageAcceptableUse, &site.Legal.AcceptableUse))
		mux.Handle("POST "+prefix+"/contact", csrf.Handler(s.contact(site)))
		if prefix != "" {
			mux.Handle("GET "+prefix, http.RedirectHandler(prefix+"/", http.StatusMovedPermanently))
		}
	}
	mux.HandleFunc("/", s.notFound)

	var h http.Handler = mux
	h = s.canonicalHost(h)
	h = s.securityHeaders(h)
	h = s.logRequests(h)
	h = s.recoverPanics(h)
	return h
}

// buildJSONLD prepares the Organization block and the CSP that allows
// exactly this script by its hash, without 'unsafe-inline'.
func (s *Server) buildJSONLD() error {
	org := map[string]any{
		"@context": "https://schema.org",
		"@type":    "Organization",
		"name":     content.Default().Meta.SiteName,
		"url":      s.cfg.SiteURL + "/",
		"logo":     s.cfg.SiteURL + render.StaticPrefix + "icon-512.png",
		"sameAs":   []string{s.cfg.LinkedInURL},
	}
	b, err := json.Marshal(org) // escapes <, > and & for safe embedding
	if err != nil {
		return err
	}
	sum := sha256.Sum256(b)
	hash := "'sha256-" + base64.StdEncoding.EncodeToString(sum[:]) + "'"
	s.jsonLD = template.HTML(`<script type="application/ld+json">` + string(b) + `</script>`)

	csp := []string{
		"default-src 'self'",
		"script-src 'self' " + hash,
		"style-src 'self'",
		"img-src 'self'",
		"font-src 'self'",
		"connect-src 'self'",
		"manifest-src 'self'",
		"form-action 'self'",
		"frame-ancestors 'none'",
		"base-uri 'none'",
		"object-src 'none'",
	}
	if s.cfg.Production() {
		csp = append(csp, "upgrade-insecure-requests")
	}
	s.csp = strings.Join(csp, "; ")
	return nil
}

func hostOf(siteURL string) string {
	_, rest, _ := strings.Cut(siteURL, "://")
	host, _, _ := strings.Cut(rest, ":")
	return host
}

// localeFor returns the locale whose prefix matches the request path.
func (s *Server) localeFor(path string) *content.Site {
	for _, site := range s.locales {
		p := site.Locale.Prefix
		if p != "" && (path == p || strings.HasPrefix(path, p+"/")) {
			return site
		}
	}
	return s.locales[0]
}

// abs returns the absolute URL of a site path.
func (s *Server) abs(path string) string { return s.cfg.SiteURL + path }

func (s *Server) assetURL(path string) string {
	u, err := s.assets.URL(path)
	if err != nil {
		panic(fmt.Sprintf("server: %v", err)) // embedded at build time; a test covers it
	}
	return u
}
