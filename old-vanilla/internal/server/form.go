// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"mime"
	"net/http"
	"strings"

	"github.com/tiefer-labs/web/internal/contact"
	"github.com/tiefer-labs/web/internal/content"
)

// Form states shown in the contact section.
const (
	statusSent    = "sent"
	statusInvalid = "invalid"
	statusExpired = "expired"
	statusLimited = "limited"
	statusError   = "error"
)

// formView is the contact form as rendered.
type formView struct {
	Status string
	Token  string
	Values contact.Submission
	Errors map[string]string // field name to message
}

// Err returns the message for a field, or "".
func (f *formView) Err(field string) string { return f.Errors[field] }

func (s *Server) newForm(*http.Request) *formView {
	return &formView{Token: s.tokens.Issue()}
}

// formReply is the JSON answer for the script.
type formReply struct {
	Status  string            `json:"status"`
	Message string            `json:"message"`
	Errors  map[string]string `json:"errors,omitempty"`
	Token   string            `json:"token,omitempty"`
}

// contactForm handles POST /contact. Without JavaScript it answers with
// Post/Redirect/Get on success and with the page and its errors
// otherwise; with JavaScript (Accept: application/json) it answers JSON.
func (s *Server) contactForm(site *content.Site) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store")
		mt, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
		if err != nil || mt != "application/x-www-form-urlencoded" {
			plainError(w, http.StatusUnsupportedMediaType)
			return
		}
		if err := r.ParseForm(); err != nil {
			plainError(w, http.StatusBadRequest)
			return
		}
		key := contact.ClientKey(s.clientAddr(r))
		if !s.attempts.Allow(key) {
			s.formAnswer(w, r, site, http.StatusTooManyRequests, statusLimited, contact.Submission{}, nil)
			return
		}
		sub := contact.Parse(r.PostForm)
		// A filled honeypot or a form sent faster than a person can type
		// gets the success answer, so bots learn nothing; nothing is sent.
		if r.PostForm.Get("website") != "" {
			s.formAnswer(w, r, site, http.StatusOK, statusSent, contact.Submission{}, nil)
			return
		}
		switch err := s.tokens.Check(r.PostForm.Get("t")); {
		case errors.Is(err, contact.ErrTooFast):
			s.formAnswer(w, r, site, http.StatusOK, statusSent, contact.Submission{}, nil)
			return
		case err != nil:
			s.formAnswer(w, r, site, http.StatusForbidden, statusExpired, sub, nil)
			return
		}
		if problems := sub.Validate(roleValues(site)); len(problems) > 0 {
			s.formAnswer(w, r, site, http.StatusUnprocessableEntity, statusInvalid, sub, messages(site, problems))
			return
		}
		if err := s.tokens.Use(r.PostForm.Get("t")); err != nil {
			s.formAnswer(w, r, site, http.StatusForbidden, statusExpired, sub, nil)
			return
		}
		if !s.sends.Allow(key) || !s.globalSends.Allow("all") {
			s.formAnswer(w, r, site, http.StatusTooManyRequests, statusLimited, sub, nil)
			return
		}
		if err := s.mailer.Send(r.Context(), s.message(site, sub)); err != nil {
			// The error may name the mail server; it never holds the
			// message, the sender or the SMTP password.
			s.log.Error("contact delivery failed", "err", err)
			s.formAnswer(w, r, site, http.StatusBadGateway, statusError, sub, nil)
			return
		}
		s.log.Info("contact message sent")
		s.formAnswer(w, r, site, http.StatusOK, statusSent, contact.Submission{}, nil)
	}
}

// formAnswer replies in the form the client asked for.
func (s *Server) formAnswer(w http.ResponseWriter, r *http.Request, site *content.Site, code int, status string, sub contact.Submission, errs map[string]string) {
	if wantsJSON(r) {
		b, _ := json.Marshal(formReply{Status: status, Message: statusMessage(site, status), Errors: errs, Token: s.tokens.Issue()})
		writeBody(w, r, code, "application/json; charset=utf-8", b)
		return
	}
	if status == statusSent {
		http.Redirect(w, r, (&page{Site: site}).Href("/")+"?sent=1#contact", http.StatusSeeOther)
		return
	}
	p := s.newPage(site, PageIndex, "/")
	p.OnDark = true
	p.NoIndex = true
	p.Form = &formView{Status: status, Token: s.tokens.Issue(), Values: sub, Errors: errs}
	s.renderPage(w, r, code, p)
}

func wantsJSON(r *http.Request) bool {
	return strings.Contains(r.Header.Get("Accept"), "application/json")
}

func statusMessage(site *content.Site, status string) string {
	c := site.Contact
	switch status {
	case statusSent:
		return c.Success
	case statusInvalid:
		return c.Invalid
	case statusExpired:
		return c.Expired
	case statusLimited:
		return c.Limited
	}
	return c.ErrorBefore
}

func roleValues(site *content.Site) []string {
	out := make([]string, len(site.Contact.Roles))
	for i, o := range site.Contact.Roles {
		out[i] = o.Value
	}
	return out
}

// messages turns validation problems into the copy of the locale.
func messages(site *content.Site, problems map[string]contact.Problem) map[string]string {
	c := site.Contact
	text := map[string]map[contact.Problem]string{
		"name":    {contact.Missing: c.NameMissing, contact.TooLong: c.NameTooLong, contact.Invalid: c.NameMissing},
		"email":   {contact.Missing: c.EmailMissing, contact.Invalid: c.EmailInvalid},
		"org":     {contact.TooLong: c.OrgTooLong, contact.Invalid: c.OrgTooLong},
		"role":    {contact.Invalid: c.RoleInvalid},
		"message": {contact.Missing: c.MessageMissing, contact.TooLong: c.MessageTooLong, contact.Invalid: c.MessageMissing},
		"consent": {contact.Missing: c.ConsentMissing},
	}
	out := map[string]string{}
	for field, p := range problems {
		if m := text[field][p]; m != "" {
			out[field] = m
		} else {
			out[field] = c.Invalid
		}
	}
	return out
}

// message builds the mail for a valid submission. The sender's address
// goes into Reply-To only; From is always CONTACT_FROM.
func (s *Server) message(site *content.Site, sub contact.Submission) contact.Message {
	role := sub.Role
	for _, o := range site.Contact.Roles {
		if o.Value == sub.Role {
			role = o.Label
		}
	}
	body := fmt.Sprintf("Name: %s\nEmail: %s\nOrganisation: %s\nRole: %s\nConsent to reply: yes\n\n%s\n\n%s\n",
		sub.Name, sub.Email, orDash(sub.Org), orDash(role), site.Contact.MessageLabel, sub.Message)
	return contact.Message{
		From:    s.cfg.ContactFrom,
		To:      s.cfg.ContactTo,
		ReplyTo: sub.Email,
		Subject: "Website contact: " + sub.Name,
		Body:    body,
		Date:    s.now(),
		Domain:  s.cfg.SiteURL.Hostname(),
	}
}

func orDash(s string) string {
	if s == "" {
		return "-"
	}
	return s
}

// mailerFromConfig returns the SMTP mailer, or nil when SMTP is not set.
func (s *Server) mailerFromConfig() contact.Mailer {
	if !s.cfg.ContactEnabled() {
		return nil
	}
	return &contact.SMTPMailer{
		Host:       s.cfg.SMTP.Host,
		Port:       s.cfg.SMTP.Port,
		User:       s.cfg.SMTP.User,
		Pass:       s.cfg.SMTP.Pass.Reveal(),
		LocalName:  s.cfg.SiteURL.Hostname(),
		SkipVerify: s.cfg.SMTP.SkipVerify,
	}
}
