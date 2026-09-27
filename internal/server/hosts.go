// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package server

import (
	"net"
	"net/http"
	"strings"
)

// canonicalHost serves the site under exactly one origin in production,
// the host of SITE_URL (https://tiefer.space):
//
//   - requests for www.<host> are redirected to the canonical host;
//   - behind a trusted proxy, plain HTTP requests are redirected to HTTPS;
//   - requests for any other Host header are refused with 421, which
//     blocks host header injection and DNS rebinding.
//
// /healthz is exempt, because container health checks call it on
// 127.0.0.1. In development every host is accepted.
func (s *Server) canonicalHost(next http.Handler) http.Handler {
	if !s.cfg.Production() {
		return next
	}
	canonical := strings.ToLower(s.cfg.SiteHost())
	alias := "www." + canonical
	if strings.HasPrefix(canonical, "www.") {
		alias = strings.TrimPrefix(canonical, "www.")
	}
	redirect := func(w http.ResponseWriter, r *http.Request) {
		code := http.StatusPermanentRedirect
		if r.Method == http.MethodGet || r.Method == http.MethodHead {
			code = http.StatusMovedPermanently
		}
		http.Redirect(w, r, "https://"+canonical+r.URL.RequestURI(), code)
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/healthz" {
			next.ServeHTTP(w, r)
			return
		}
		host := normalizeHost(r.Host)
		switch {
		case host == alias:
			redirect(w, r)
		case host != canonical:
			w.Header().Set("Cache-Control", "no-store")
			http.Error(w, "Misdirected Request", http.StatusMisdirectedRequest)
		case s.cfg.TrustProxy && !strings.EqualFold(forwardedProto(r), "https"):
			redirect(w, r)
		default:
			next.ServeHTTP(w, r)
		}
	})
}

// normalizeHost lowercases a Host header and removes a trailing dot and
// the default HTTP and HTTPS ports.
func normalizeHost(h string) string {
	h = strings.ToLower(strings.TrimSpace(h))
	if host, port, err := net.SplitHostPort(h); err == nil && (port == "80" || port == "443") {
		h = host
	}
	return strings.TrimSuffix(h, ".")
}

// forwardedProto returns the scheme the client used, as reported by the
// reverse proxy in X-Forwarded-Proto (the last value, set by our proxy).
func forwardedProto(r *http.Request) string {
	v := r.Header.Values("X-Forwarded-Proto")
	if len(v) == 0 {
		return "https" // no proxy information: TLS is terminated in front of us
	}
	parts := strings.Split(v[len(v)-1], ",")
	return strings.TrimSpace(parts[len(parts)-1])
}
