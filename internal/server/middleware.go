// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package server

import (
	"crypto/subtle"
	"log/slog"
	"net"
	"net/http"
	"net/netip"
	"regexp"
	"strings"

	"github.com/tiefer-labs/web/internal/config"
)

// securityHeaders returns the headers sent with every response. scripts
// holds extra CSP sources for script-src, such as the hash of the JSON-LD
// block.
func securityHeaders(cfg *config.Config, scripts string) http.Header {
	csp := []string{
		"default-src 'none'",
		"script-src 'self'" + scripts,
		"script-src-attr 'none'",
		"style-src 'self'",
		"style-src-attr 'none'",
		"img-src 'self'",
		"font-src 'self'",
		"connect-src 'self'",
		"manifest-src 'self'",
		"form-action 'self'",
		"frame-ancestors 'none'",
		"base-uri 'none'",
		"object-src 'none'",
		"require-trusted-types-for 'script'",
		"trusted-types 'none'",
	}
	if cfg.Production() {
		csp = append(csp, "upgrade-insecure-requests")
	}
	hsts := "max-age=63072000; includeSubDomains"
	if cfg.HSTSPreload {
		hsts += "; preload"
	}
	h := http.Header{}
	h.Set("Content-Security-Policy", strings.Join(csp, "; "))
	h.Set("Strict-Transport-Security", hsts)
	h.Set("X-Content-Type-Options", "nosniff")
	h.Set("X-Frame-Options", "DENY")
	h.Set("Referrer-Policy", "no-referrer")
	h.Set("Permissions-Policy", "accelerometer=(), autoplay=(), browsing-topics=(), camera=(), display-capture=(), encrypted-media=(), fullscreen=(), geolocation=(), gyroscope=(), hid=(), idle-detection=(), magnetometer=(), microphone=(), midi=(), payment=(), picture-in-picture=(), publickey-credentials-get=(), screen-wake-lock=(), serial=(), usb=(), xr-spatial-tracking=()")
	h.Set("Cross-Origin-Opener-Policy", "same-origin")
	h.Set("Cross-Origin-Resource-Policy", "same-origin")
	h.Set("Cross-Origin-Embedder-Policy", "require-corp")
	h.Set("Origin-Agent-Cluster", "?1")
	h.Set("X-Permitted-Cross-Domain-Policies", "none")
	h.Set("X-DNS-Prefetch-Control", "off")
	return h
}

// securityHeaders sets the security headers before the handler runs, so
// errors, redirects and panics carry them too.
func (s *Server) securityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		dst := w.Header()
		for k, v := range s.headers {
			dst[k] = v
		}
		next.ServeHTTP(w, r)
	})
}

// limitRequests enforces the method list, the URL length and the body
// rules before any handler runs.
func (s *Server) limitRequests(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet, http.MethodHead, http.MethodPost:
		default:
			w.Header().Set("Allow", "GET, HEAD, POST")
			plainError(w, http.StatusMethodNotAllowed)
			return
		}
		if len(r.RequestURI) > MaxURLLength {
			plainError(w, http.StatusRequestURITooLong)
			return
		}
		if headerSize(r) > MaxHeaderBytes {
			plainError(w, http.StatusRequestHeaderFieldsTooLarge)
			return
		}
		if r.Method != http.MethodPost {
			if r.ContentLength > 0 || len(r.TransferEncoding) > 0 {
				plainError(w, http.StatusBadRequest)
				return
			}
		} else {
			if r.ContentLength > MaxFormBytes {
				plainError(w, http.StatusRequestEntityTooLarge)
				return
			}
			r.Body = http.MaxBytesReader(w, r.Body, MaxFormBytes)
		}
		next.ServeHTTP(w, r)
	})
}

// headerSize approximates the size of the request header as sent.
func headerSize(r *http.Request) int {
	n := len(r.Method) + len(r.RequestURI) + len(r.Host) + 16
	for k, vs := range r.Header {
		for _, v := range vs {
			n += len(k) + len(v) + 4
		}
	}
	return n
}

// frontDoor admits requests only from the configured Azure Front Door
// profile when the site runs behind it, and only for the site's own host
// otherwise in production. /healthz stays reachable for the platform's
// health probe, which does not pass through Front Door.
func (s *Server) frontDoor(next http.Handler) http.Handler {
	want := []byte(s.cfg.FrontDoorID.Reveal())
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/healthz" {
			next.ServeHTTP(w, r)
			return
		}
		switch {
		case s.cfg.BehindFrontDoor && len(want) > 0:
			got := []byte(r.Header.Get("X-Azure-FDID"))
			if subtle.ConstantTimeCompare(got, want) != 1 {
				plainError(w, http.StatusForbidden)
				return
			}
		case s.cfg.Production() && !s.cfg.BehindFrontDoor:
			if !strings.EqualFold(r.Host, s.cfg.SiteURL.Host) {
				plainError(w, http.StatusMisdirectedRequest)
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}

// clientAddr returns the address used for rate limiting. Behind Front
// Door it is the X-Azure-ClientIP header, which Front Door sets and which
// can only be trusted because frontDoor has checked the profile ID; without
// a configured ID the header is ignored. It is kept in memory only and
// never logged.
func (s *Server) clientAddr(r *http.Request) netip.Addr {
	if s.cfg.BehindFrontDoor && s.cfg.FrontDoorID.Set() {
		if a, err := netip.ParseAddr(strings.TrimSpace(r.Header.Get("X-Azure-ClientIP"))); err == nil {
			return a.Unmap()
		}
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		host = r.RemoteAddr
	}
	a, err := netip.ParseAddr(host)
	if err != nil {
		return netip.Addr{}
	}
	return a.Unmap()
}

// recorder captures the status and size of a response for the log.
type recorder struct {
	http.ResponseWriter
	status int
	bytes  int
}

func (r *recorder) WriteHeader(code int) {
	if r.status == 0 {
		r.status = code
	}
	r.ResponseWriter.WriteHeader(code)
}

func (r *recorder) Write(b []byte) (int, error) {
	if r.status == 0 {
		r.status = http.StatusOK
	}
	n, err := r.ResponseWriter.Write(b)
	r.bytes += n
	return n, err
}

func (r *recorder) Unwrap() http.ResponseWriter { return r.ResponseWriter }

// logRequests logs one line per request: method, path, status, size and
// duration. No IP address, query string, user agent or form content.
func (s *Server) logRequests(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := s.now()
		rec := &recorder{ResponseWriter: w}
		next.ServeHTTP(rec, r)
		if r.URL.Path == "/healthz" && rec.status == http.StatusOK {
			return // the platform probes every few seconds
		}
		path := r.URL.Path
		if len(path) > 128 {
			path = path[:128]
		}
		s.log.Info("request",
			"method", r.Method,
			"path", path,
			"status", rec.status,
			"bytes", rec.bytes,
			"ms", s.now().Sub(start).Milliseconds())
	})
}

// recoverPanics turns a panic into a 500 with the security headers
// already set, and logs it without request data.
func (s *Server) recoverPanics(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if v := recover(); v != nil {
				if v == http.ErrAbortHandler {
					panic(v)
				}
				s.log.Error("panic", "path", r.URL.Path, "err", v)
				plainError(w, http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

// Addresses in log text: IPv4 and IPv6, with an optional port.
var addressPattern = regexp.MustCompile(`\[[0-9A-Fa-f:.%a-z]+\](?::\d+)?|\b(?:\d{1,3}\.){3}\d{1,3}(?::\d+)?\b|\b[0-9A-Fa-f]{0,4}(?::[0-9A-Fa-f]{0,4}){2,7}(?:%[0-9A-Za-z]+)?\b`)

// redactingWriter receives the messages of net/http's own error log (for
// example "http: panic serving 198.51.100.7:5555") and logs them without
// client addresses.
type redactingWriter struct{ log *slog.Logger }

func (w redactingWriter) Write(p []byte) (int, error) {
	msg := addressPattern.ReplaceAllString(strings.TrimSpace(string(p)), "[address]")
	w.log.Warn("http server", "msg", msg)
	return len(p), nil
}
