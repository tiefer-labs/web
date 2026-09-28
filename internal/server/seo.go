// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package server

import (
	"encoding/json"
	"net/http"
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
