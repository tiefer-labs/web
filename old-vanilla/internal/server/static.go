// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package server

import (
	"net/http"
	"strconv"

	"github.com/tiefer-labs/web/internal/render"
)

// static serves hashed assets. Only exact hashed URLs exist; they never
// change, so they are cached for a year.
func (s *Server) static(w http.ResponseWriter, r *http.Request) {
	as, ok := s.assets.ByURL(r.URL.Path)
	if !ok {
		s.notFound(w, r)
		return
	}
	s.serveAsset(w, r, as, "public, max-age=31536000, immutable")
}

// rootAsset serves a file that browsers request at a fixed path, such as
// /favicon.ico, with a shorter cache time.
func (s *Server) rootAsset(p string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		as, ok := s.assets.ByPath(p)
		if !ok {
			s.notFound(w, r)
			return
		}
		s.serveAsset(w, r, as, "public, max-age=86400")
	}
}

func (s *Server) serveAsset(w http.ResponseWriter, r *http.Request, as *render.Asset, cache string) {
	h := w.Header()
	h.Set("Cache-Control", cache)
	h.Set("ETag", as.ETag)
	h.Set("Content-Type", as.ContentType)
	if as.Gzip != nil {
		h.Add("Vary", "Accept-Encoding")
	}
	if match := r.Header.Get("If-None-Match"); match != "" && (match == as.ETag || match == "W/"+as.ETag) {
		w.WriteHeader(http.StatusNotModified)
		return
	}
	body := as.Body
	if as.Gzip != nil && acceptsGzip(r) {
		h.Set("Content-Encoding", "gzip")
		body = as.Gzip
	}
	h.Set("Content-Length", strconv.Itoa(len(body)))
	if r.Method != http.MethodHead {
		// The body is an embedded file chosen by an exact map lookup;
		// no request data reaches it.
		_, _ = w.Write(body) // #nosec G705 -- embedded asset bytes, see docs/SECURITY-DECISIONS.md
	}
}
