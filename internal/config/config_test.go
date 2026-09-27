// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package config

import (
	"strings"
	"testing"
)

func load(env map[string]string) (*Config, error) {
	return FromLookup(func(k string) (string, bool) {
		v, ok := env[k]
		return v, ok
	})
}

func TestDevelopmentDefaults(t *testing.T) {
	c, err := load(nil)
	if err != nil {
		t.Fatal(err)
	}
	if c.Env != Development || c.Port != "8080" || c.SiteURL != "http://localhost:8080" {
		t.Errorf("unexpected defaults: %+v", c)
	}
	if len(c.CSRFSecret) < 32 {
		t.Error("development must get a random CSRF secret")
	}
	if c.ContactEnabled() {
		t.Error("contact form must be off without SMTP")
	}
	if c.Legal.Name != "[COMPANY_LEGAL_NAME]" || c.Legal.Reviewed {
		t.Errorf("legal defaults: %+v", c.Legal)
	}
	if got := strings.Join(c.DefaultsInUse(), ","); !strings.Contains(got, "LEGAL_TAX_ID") || !strings.Contains(got, "LEGAL_REVIEWED") {
		t.Errorf("DefaultsInUse = %s", got)
	}
}

func TestProductionRequiresValues(t *testing.T) {
	_, err := load(map[string]string{"ENV": "production"})
	if err == nil {
		t.Fatal("production without CSRF_SECRET must fail")
	}
	for _, name := range []string{"CSRF_SECRET"} {
		if !strings.Contains(err.Error(), name) {
			t.Errorf("error does not mention %s: %v", name, err)
		}
	}
}

func TestProductionDefaultsToTheRealDomain(t *testing.T) {
	c, err := load(map[string]string{
		"ENV":           "production",
		"CONTACT_EMAIL": "hello@tiefer.space",
		"CSRF_SECRET":   strings.Repeat("s", 64),
	})
	if err != nil {
		t.Fatal(err)
	}
	if c.SiteURL != "https://tiefer.space" || c.SiteHost() != "tiefer.space" {
		t.Errorf("SiteURL = %q, SiteHost = %q", c.SiteURL, c.SiteHost())
	}
	for _, name := range c.DefaultsInUse() {
		if name == "SITE_URL" {
			t.Error("the production domain must not be reported as a placeholder")
		}
	}
}

func TestTieferSpaceDefaults(t *testing.T) {
	c, err := load(map[string]string{"SMTP_HOST": "smtp.example.org"})
	if err != nil {
		t.Fatal(err)
	}
	if c.ContactEmail != "hello@tiefer.space" || c.SecurityEmail != "hello@tiefer.space" {
		t.Errorf("contact %q, security %q", c.ContactEmail, c.SecurityEmail)
	}
	if !c.ContactEnabled() || c.ContactTo != "hello@tiefer.space" || c.ContactFrom != "website@tiefer.space" {
		t.Errorf("contact form defaults: to %q, from %q", c.ContactTo, c.ContactFrom)
	}
	if c.RepoURL != "https://github.com/tiefer-labs/web" {
		t.Errorf("RepoURL = %q", c.RepoURL)
	}
	c, err = load(map[string]string{"REPO_URL": ""})
	if err != nil || c.RepoURL != "" {
		t.Errorf("REPO_URL set to empty must hide the link: %q, %v", c.RepoURL, err)
	}
	for _, name := range c.DefaultsInUse() {
		if name == "CONTACT_EMAIL" {
			t.Error("hello@tiefer.space is not a placeholder")
		}
	}
}

func TestProductionValid(t *testing.T) {
	c, err := load(map[string]string{
		"ENV":           "production",
		"SITE_URL":      "https://tiefer.example/",
		"CONTACT_EMAIL": "hello@tiefer.example",
		"CSRF_SECRET":   strings.Repeat("s", 64),
		"SMTP_HOST":     "smtp.tiefer.example",
		"SMTP_PORT":     "465",
		"SMTP_USER":     "web",
		"SMTP_PASS":     "secret with spaces",
		"CONTACT_TO":    "team@tiefer.example",
		"CONTACT_FROM":  "web@tiefer.example",
	})
	if err != nil {
		t.Fatal(err)
	}
	if c.SiteURL != "https://tiefer.example" {
		t.Errorf("trailing slash not trimmed: %q", c.SiteURL)
	}
	if !c.ContactEnabled() || c.SMTP.Port != 465 || c.SMTP.Pass != "secret with spaces" {
		t.Errorf("SMTP: %+v", c.SMTP)
	}
}

func TestInvalidValues(t *testing.T) {
	cases := []struct {
		env  map[string]string
		want string
	}{
		{map[string]string{"ENV": "staging"}, "ENV"},
		{map[string]string{"PORT": "http"}, "PORT"},
		{map[string]string{"SITE_URL": "tiefer.example"}, "SITE_URL"},
		{map[string]string{"SITE_URL": "https://tiefer.example/path"}, "SITE_URL"},
		{map[string]string{"ENV": "production", "SITE_URL": "http://tiefer.example", "CONTACT_EMAIL": "a@b.example", "CSRF_SECRET": strings.Repeat("s", 32)}, "https"},
		{map[string]string{"CONTACT_EMAIL": "not an email"}, "CONTACT_EMAIL"},
		{map[string]string{"CONTACT_EMAIL": "a@b.example\r\nBcc: c@d.example"}, "CONTACT_EMAIL"},
		{map[string]string{"LINKEDIN_URL": "javascript:alert(1)"}, "LINKEDIN_URL"},
		{map[string]string{"SECURITY_EMAIL": "security at tiefer"}, "SECURITY_EMAIL"},
		{map[string]string{"REPO_URL": "ftp://example.org"}, "REPO_URL"},
		{map[string]string{"CSRF_SECRET": "short"}, "CSRF_SECRET"},
		{map[string]string{"LEGAL_REVIEWED": "yes please"}, "LEGAL_REVIEWED"},
		{map[string]string{"CONTACT_TO": "team@example.org"}, "SMTP_HOST"},
		{map[string]string{"SMTP_HOST": "smtp.example.org", "CONTACT_TO": "a@b.example", "CONTACT_FROM": "c@d.example", "SMTP_USER": "u"}, "SMTP_PASS"},
		{map[string]string{"SMTP_HOST": "smtp.example.org", "CONTACT_TO": "a@b.example", "CONTACT_FROM": "c@d.example", "SMTP_PORT": "x"}, "SMTP_PORT"},
		{map[string]string{"ENV": "production", "SITE_URL": "https://a.example", "CONTACT_EMAIL": "a@b.example", "CSRF_SECRET": strings.Repeat("s", 32), "SMTP_SKIP_VERIFY": "true"}, "SMTP_SKIP_VERIFY"},
	}
	for _, c := range cases {
		_, err := load(c.env)
		if err == nil || !strings.Contains(err.Error(), c.want) {
			t.Errorf("%v: got %v, want an error about %s", c.env, err, c.want)
		}
	}
}
