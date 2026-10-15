// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package server

import (
	"net/http"
	"strings"
	"time"
)

// SecurityTxtExpires is the Expires date of /.well-known/security.txt
// (RFC 9116). It is fixed on purpose: a test fails 30 days before it, so
// someone confirms the contact and moves the date forward (at most a year).
var SecurityTxtExpires = time.Date(2027, time.September, 1, 0, 0, 0, 0, time.UTC)

func (s *Server) buildSecurityTxt() []byte {
	lines := []string{
		"Contact: mailto:" + s.cfg.ContactEmail,
		"Expires: " + SecurityTxtExpires.Format(time.RFC3339),
		"Preferred-Languages: en",
		"Canonical: " + s.cfg.Site() + "/.well-known/security.txt",
	}
	if s.cfg.RepoURL != "" {
		lines = append(lines, "Policy: "+strings.TrimSuffix(s.cfg.RepoURL, "/")+"/blob/main/SECURITY.md")
	}
	return []byte(strings.Join(lines, "\n") + "\n")
}

func (s *Server) securityTxt(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "public, max-age=86400")
	writeBody(w, r, http.StatusOK, "text/plain; charset=utf-8", s.securityTxtBody)
}
