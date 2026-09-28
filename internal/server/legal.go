// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package server

import (
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/tiefer-labs/web/internal/content"
	"github.com/tiefer-labs/web/internal/textcheck"
)

// part is a piece of legal text: plain text, a link, or a highlighted
// placeholder. Templates render parts, so no HTML is built in Go.
type part struct {
	Text string
	Href string // link target, if the part is a link
	Mark bool   // a placeholder a lawyer must replace
}

// legalView is a legal page with its tokens filled in.
type legalView struct {
	Title    string
	Intro    []part
	Sections []legalSection
	Updated  []part
}

type legalSection struct {
	ID    string
	Title string
	Facts []legalFact
	Paras [][]part
	List  [][]part
}

type legalFact struct {
	Label string
	Value []part
}

var tokenPattern = regexp.MustCompile(`\{[a-z_]+\}`)

// legalTokens maps the {tokens} of the legal copy to configuration.
func (s *Server) legalTokens() map[string]string {
	l := s.cfg.Legal
	return map[string]string{
		"legal_name":         l.Name,
		"legal_form":         l.Form,
		"address":            l.Address,
		"tax_id":             l.TaxID,
		"registration":       l.Registration,
		"director":           l.Director,
		"hosting_provider":   l.HostingProvider,
		"hosting_country":    l.HostingCountry,
		"smtp_provider":      l.SMTPProvider,
		"smtp_country":       l.SMTPCountry,
		"log_retention_days": strconv.Itoa(s.cfg.LogRetentionDays),
		"contact_email":      s.cfg.ContactEmail,
		"site_url":           s.cfg.Site(),
	}
}

// fill replaces {tokens} with configuration values, turns the email
// address and site URL into links, and marks [PLACEHOLDERS].
func fill(s string, tokens map[string]string) []part {
	var out []part
	text := func(t string) {
		for t != "" {
			loc := textcheck.PlaceholderPattern.FindStringIndex(t)
			if loc == nil {
				out = append(out, part{Text: t})
				return
			}
			if loc[0] > 0 {
				out = append(out, part{Text: t[:loc[0]]})
			}
			out = append(out, part{Text: t[loc[0]:loc[1]], Mark: true})
			t = t[loc[1]:]
		}
	}
	for s != "" {
		loc := tokenPattern.FindStringIndex(s)
		if loc == nil {
			text(s)
			break
		}
		text(s[:loc[0]])
		name := s[loc[0]+1 : loc[1]-1]
		v, ok := tokens[name]
		switch {
		case !ok:
			panic("server: unknown legal token {" + name + "}") // a test renders every page
		case name == "contact_email":
			out = append(out, part{Text: v, Href: "mailto:" + v})
		case name == "site_url":
			out = append(out, part{Text: v, Href: v + "/"})
		default:
			text(v) // a config value may itself be a placeholder
		}
		s = s[loc[1]:]
	}
	return out
}

func (s *Server) legalPage(lp *content.LegalPage) *legalView {
	tok := s.legalTokens()
	v := &legalView{Title: lp.Title, Intro: fill(lp.Intro, tok), Updated: fill(lp.Updated, tok)}
	for _, sec := range lp.Sections {
		ls := legalSection{ID: slug(sec.Title), Title: sec.Title}
		for _, f := range sec.Facts {
			ls.Facts = append(ls.Facts, legalFact{Label: f.Label, Value: fill(f.Value, tok)})
		}
		for _, p := range sec.Paras {
			ls.Paras = append(ls.Paras, fill(p, tok))
		}
		for _, li := range sec.List {
			ls.List = append(ls.List, fill(li, tok))
		}
		v.Sections = append(v.Sections, ls)
	}
	return v
}

// slug turns a heading into a fragment id: "How long we keep it" becomes
// "how-long-we-keep-it".
func slug(s string) string {
	var b strings.Builder
	dash := false
	for _, r := range strings.ToLower(s) {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
			b.WriteRune(r)
			dash = false
		} else if !dash && b.Len() > 0 {
			b.WriteByte('-')
			dash = true
		}
	}
	return strings.TrimSuffix(b.String(), "-")
}

func (s *Server) legal(site *content.Site, name string, lp *content.LegalPage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		p := s.newPage(site, name, "/"+name)
		p.Title = lp.Title + " | " + site.Meta.SiteName
		p.Description = lp.Description
		p.Legal = s.legalPage(lp)
		s.renderPage(w, r, http.StatusOK, p)
	}
}
