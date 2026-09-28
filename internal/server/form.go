// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package server

import "net/http"

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
	Values formValues
	Errors map[string]string // field name to message
}

// formValues are the submitted values shown again after an error.
type formValues struct {
	Name, Email, Org, Role, Message string
	Consent                         bool
}

// Err returns the message for a field, or "".
func (f *formView) Err(field string) string { return f.Errors[field] }

func (s *Server) newForm(*http.Request) *formView { return &formView{} }
