// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package server

import (
	"encoding/json"
	"errors"
	"net"
	"net/http"
	"strings"

	"github.com/tiefer-labs/web/internal/contact"
	"github.com/tiefer-labs/web/internal/content"
)

// maxFormBytes limits the request body of the contact form.
const maxFormBytes = 32 << 10

// contactReply is the JSON answer for submissions made with JavaScript.
type contactReply struct {
	Status string            `json:"status"` // sent, invalid, expired or error
	Errors map[string]string `json:"errors,omitempty"`
	Token  string            `json:"token,omitempty"`
}

// contact handles POST /contact. Without JavaScript it follows the
// Post/Redirect/Get pattern on success and re-renders the page with the
// entered values on error. With JavaScript (Accept: application/json) it
// answers with JSON.
//
// Spam protection, in order: a per-client rate limit, a hidden honeypot
// field, a signed form token (see contact.Tokens) with a minimum fill
// time, and single use of each token. Cross-origin posts are rejected
// before this handler by http.CrossOriginProtection.
func (s *Server) contact(site *content.Site) http.HandlerFunc {
	var roles []string
	for _, o := range site.Contact.Roles {
		roles = append(roles, o.Value)
	}
	home := site.Locale.Prefix + "/"

	return func(w http.ResponseWriter, r *http.Request) {
		wantsJSON := strings.Contains(r.Header.Get("Accept"), "application/json")
		if !s.cfg.ContactEnabled() {
			http.Redirect(w, r, home+"#contact", http.StatusSeeOther)
			return
		}

		reply := func(status int, st *formState) {
			// Answers to a post may hold what the visitor typed: never cache them.
			w.Header().Set("Cache-Control", "no-store")
			if st.Status == "sent" && !wantsJSON {
				http.Redirect(w, r, home+"?sent=1#contact", http.StatusSeeOther)
				return
			}
			if st.Token == "" {
				st.Token = s.tokens.Issue()
			}
			if wantsJSON {
				b, _ := json.Marshal(contactReply{Status: st.Status, Errors: st.Errors, Token: st.Token})
				writeBody(w, r, status, "application/json", b)
				return
			}
			p := s.newPage(site, PageIndex, "/")
			p.HeaderDark = true
			p.Form = st
			s.renderPage(w, r, status, p)
		}

		if !s.limiter.Allow(clientKey(s.clientIP(r))) {
			s.log.Warn("contact: rate limited")
			reply(http.StatusTooManyRequests, &formState{Status: "error"})
			return
		}

		r.Body = http.MaxBytesReader(w, r.Body, maxFormBytes)
		if err := r.ParseForm(); err != nil {
			reply(http.StatusBadRequest, &formState{Status: "error"})
			return
		}
		sub, problems := contact.Parse(r.PostForm, roles)
		st := &formState{Values: sub, Token: r.PostForm.Get(contact.FieldToken)}

		// Bots that fill the hidden field get a normal looking success.
		if r.PostForm.Get(contact.FieldHoneypot) != "" {
			s.log.Info("contact: honeypot filled, message dropped")
			reply(http.StatusOK, &formState{Status: "sent"})
			return
		}

		nonce, err := s.tokens.Check(st.Token)
		switch {
		case errors.Is(err, contact.ErrTokenExpired):
			st.Status, st.Token = "expired", ""
			reply(http.StatusUnprocessableEntity, st)
			return
		case errors.Is(err, contact.ErrTooFast):
			s.log.Info("contact: submitted too fast, message dropped")
			st.Status = "error"
			reply(http.StatusUnprocessableEntity, st)
			return
		case err != nil:
			s.log.Info("contact: invalid form token")
			st.Status, st.Token = "error", ""
			reply(http.StatusBadRequest, st)
			return
		}

		if len(problems) > 0 {
			st.Status = "invalid"
			st.Errors = messages(site.Contact.Errors, problems)
			reply(http.StatusUnprocessableEntity, st)
			return
		}

		// A second post of the same form (a double click without
		// JavaScript, or the back button) is answered but not sent again.
		if !s.tokens.Consume(nonce) {
			reply(http.StatusOK, &formState{Status: "sent"})
			return
		}

		// A global cap on delivered messages protects the mailbox and the
		// sender reputation from spam spread over many addresses.
		if !s.sendLimiter.Allow("all") {
			s.log.Warn("contact: global send limit reached, message not sent")
			st.Status, st.Token = "error", ""
			reply(http.StatusServiceUnavailable, st)
			return
		}

		roleLabel := ""
		for _, o := range site.Contact.Roles {
			if o.Value == sub.Role {
				roleLabel = o.Label
			}
		}
		msg := contact.BuildMessage(contact.Envelope{
			From:      s.cfg.ContactFrom,
			To:        s.cfg.ContactTo,
			SiteURL:   s.cfg.SiteURL,
			RoleLabel: roleLabel,
			Now:       s.now(),
		}, sub)
		if err := s.mailer.Send(r.Context(), s.cfg.ContactFrom, []string{s.cfg.ContactTo}, msg); err != nil {
			// The error comes from the SMTP client and holds no message content.
			s.log.Error("contact: delivery failed", "err", err)
			st.Status, st.Token = "error", ""
			reply(http.StatusServiceUnavailable, st)
			return
		}
		s.log.Info("contact: message delivered")
		reply(http.StatusOK, &formState{Status: "sent"})
	}
}

// messages maps validation problems to the localised messages.
func messages(e content.FormErrors, problems contact.Errors) map[string]string {
	out := map[string]string{}
	for field, p := range problems {
		var m string
		switch field + "/" + string(p) {
		case "name/required":
			m = e.NameRequired
		case "name/too_long":
			m = e.NameTooLong
		case "email/required":
			m = e.EmailRequired
		case "email/invalid_email":
			m = e.EmailInvalid
		case "org/too_long":
			m = e.OrgTooLong
		case "role/unknown_role":
			m = e.RoleInvalid
		case "question/required":
			m = e.QuestionRequired
		case "question/too_long":
			m = e.QuestionTooLong
		case "consent/required":
			m = e.ConsentRequired
		default:
			m = e.Invalid
		}
		out[field] = m
	}
	return out
}

// clientIP returns the address used for rate limiting. It stays in
// memory and is never logged. Behind a reverse proxy (TRUST_PROXY=true)
// the right-most X-Forwarded-For entry is used, which is the one the
// proxy appended and the client cannot forge.
func (s *Server) clientIP(r *http.Request) string {
	if s.cfg.TrustProxy {
		if xff := r.Header.Values("X-Forwarded-For"); len(xff) > 0 {
			parts := strings.Split(xff[len(xff)-1], ",")
			if ip := strings.TrimSpace(parts[len(parts)-1]); ip != "" {
				return ip
			}
		}
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}
