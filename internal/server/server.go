// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

// Package server holds the routes, handlers and middleware of the site.
package server

import (
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/tiefer-labs/web/internal/config"
)

// Limits of the HTTP server. They are exported so that cmd/tiefer-web
// and the tests use the same values.
const (
	ReadHeaderTimeout = 5 * time.Second
	ReadTimeout       = 10 * time.Second
	WriteTimeout      = 30 * time.Second // covers SMTP delivery of the contact form
	IdleTimeout       = 60 * time.Second
	ShutdownTimeout   = 15 * time.Second
	MaxHeaderBytes    = 8 << 10
	MaxURLLength      = 2048
	MaxFormBytes      = 16 << 10 // the contact form, the only request with a body
)

// Options are the dependencies of the server.
type Options struct {
	Config *config.Config
	Logger *slog.Logger
	Now    func() time.Time
}

// Server serves the website.
type Server struct {
	cfg     *config.Config
	log     *slog.Logger
	now     func() time.Time
	headers http.Header // security headers set on every response
	handler http.Handler
	routes  []route
	allow   map[string][]string // exact path to its methods, for 405 replies
}

// route is one registered endpoint. The list drives the header tests,
// so every route is covered.
type route struct {
	method, pattern, sample string
}

// New builds the server and registers every route.
func New(o Options) (*Server, error) {
	if o.Logger == nil {
		o.Logger = slog.Default()
	}
	if o.Now == nil {
		o.Now = time.Now
	}
	s := &Server{cfg: o.Config, log: o.Logger, now: o.Now, allow: map[string][]string{}}
	s.headers = securityHeaders(s.cfg, "")
	mux := http.NewServeMux()
	s.handle(mux, "GET", "/healthz", "/healthz", http.HandlerFunc(s.healthz))
	mux.HandleFunc("/", s.notFound)

	var h http.Handler = mux
	h = s.limitRequests(h)
	h = s.frontDoor(h)
	h = s.securityHeaders(h)
	h = s.logRequests(h)
	h = s.recoverPanics(h)
	s.handler = h
	return s, nil
}

// Handler returns the root handler with all middleware.
func (s *Server) Handler() http.Handler { return s.handler }

// HTTPServer returns an http.Server with the limits of this package.
func (s *Server) HTTPServer(addr string) *http.Server {
	return &http.Server{
		Addr:              addr,
		Handler:           s.handler,
		ReadHeaderTimeout: ReadHeaderTimeout,
		ReadTimeout:       ReadTimeout,
		WriteTimeout:      WriteTimeout,
		IdleTimeout:       IdleTimeout,
		MaxHeaderBytes:    MaxHeaderBytes,
		// Answer "OPTIONS *" through the handler (405) instead of
		// Go's built-in reply.
		DisableGeneralOptionsHandler: true,
		ErrorLog:                     slog.NewLogLogger(s.log.Handler(), slog.LevelWarn),
	}
}

func (s *Server) handle(mux *http.ServeMux, method, pattern, sample string, h http.Handler) {
	mux.Handle(method+" "+pattern, h)
	s.routes = append(s.routes, route{method, pattern, sample})
	if pattern == sample {
		s.allow[pattern] = append(s.allow[pattern], method)
		if method == http.MethodGet {
			s.allow[pattern] = append(s.allow[pattern], http.MethodHead)
		}
	}
}

func (s *Server) healthz(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("Cache-Control", "no-store")
	_, _ = w.Write([]byte("ok"))
}

// notFound answers every request no route matched. A known path with
// the wrong method gets 405, anything else 404.
func (s *Server) notFound(w http.ResponseWriter, r *http.Request) {
	if methods, ok := s.allow[r.URL.Path]; ok {
		w.Header().Set("Allow", strings.Join(methods, ", "))
		plainError(w, http.StatusMethodNotAllowed)
		return
	}
	plainError(w, http.StatusNotFound)
}

// plainError writes a short text error. It never includes request data.
func plainError(w http.ResponseWriter, status int) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("Cache-Control", "no-store")
	w.WriteHeader(status)
	_, _ = w.Write([]byte(http.StatusText(status) + "\n"))
}
