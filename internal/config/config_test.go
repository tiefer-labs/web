// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package config

import (
	"bytes"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"testing"
)

func env(m map[string]string) func(string) (string, bool) {
	return func(k string) (string, bool) {
		v, ok := m[k]
		return v, ok
	}
}

const secret = "0123456789abcdef0123456789abcdef-test"

func prodEnv() map[string]string {
	return map[string]string{
		"ENV":         "production",
		"SITE_URL":    "https://tiefer.space",
		"CSRF_SECRET": secret,
	}
}

func TestDevelopmentDefaults(t *testing.T) {
	c, err := Load(env(nil))
	if err != nil {
		t.Fatal(err)
	}
	if c.Production() || c.Port != 8080 || c.Site() != "http://localhost:8080" {
		t.Errorf("unexpected defaults: env %q port %d site %q", c.Env, c.Port, c.Site())
	}
	if c.ContactEmail != "hello@tiefer.space" || c.LinkedInURL != "https://www.linkedin.com/company/tiefer/" {
		t.Errorf("contact %q, linkedin %q", c.ContactEmail, c.LinkedInURL)
	}
	if len(c.CSRFSecret.Reveal()) < 32 {
		t.Error("development needs a random CSRF secret")
	}
	if c.ContactEnabled() || c.RepoURL != "" || c.LogRetentionDays != 30 {
		t.Error("contact form must be off, REPO_URL empty and retention 30 days by default")
	}
	if c.Legal.Name != "[COMPANY_LEGAL_NAME]" || c.Legal.TaxID != "[VOEN]" || c.Legal.Reviewed {
		t.Errorf("legal defaults: %+v", c.Legal)
	}
	got := strings.Join(c.DefaultsInUse(), ",")
	for _, want := range []string{"SITE_URL", "LEGAL_NAME", "LEGAL_TAX_ID", "LEGAL_REVIEWED"} {
		if !strings.Contains(got, want) {
			t.Errorf("DefaultsInUse = %s, missing %s", got, want)
		}
	}
}

func TestProductionRequires(t *testing.T) {
	_, err := Load(env(map[string]string{"ENV": "production"}))
	if err == nil {
		t.Fatal("production without SITE_URL and CSRF_SECRET must fail")
	}
	for _, name := range []string{"SITE_URL", "CSRF_SECRET"} {
		if !strings.Contains(err.Error(), name) {
			t.Errorf("error does not name %s: %v", name, err)
		}
	}
	if _, err := Load(env(prodEnv())); err != nil {
		t.Errorf("minimal production config: %v", err)
	}
}

func TestInvalid(t *testing.T) {
	cases := []struct {
		set  map[string]string
		want string
	}{
		{map[string]string{"ENV": "staging"}, "ENV"},
		{map[string]string{"PORT": "http"}, "PORT"},
		{map[string]string{"PORT": "70000"}, "PORT"},
		{map[string]string{"SITE_URL": "tiefer.space"}, "SITE_URL"},
		{map[string]string{"SITE_URL": "https://tiefer.space/path"}, "SITE_URL"},
		{map[string]string{"SITE_URL": "https://user@tiefer.space"}, "SITE_URL"},
		{map[string]string{"SITE_URL": "http://tiefer.space"}, "https"},
		{map[string]string{"CONTACT_EMAIL": "not an address"}, "CONTACT_EMAIL"},
		{map[string]string{"CONTACT_EMAIL": "a@b.example\r\nBcc: c@d.example"}, "CONTACT_EMAIL"},
		{map[string]string{"CONTACT_EMAIL": "Name <a@b.example>"}, "CONTACT_EMAIL"},
		{map[string]string{"LINKEDIN_URL": "javascript:alert(1)"}, "LINKEDIN_URL"},
		{map[string]string{"REPO_URL": "http://example.org/repo"}, "REPO_URL"},
		{map[string]string{"CSRF_SECRET": "short"}, "CSRF_SECRET"},
		{map[string]string{"LEGAL_REVIEWED": "maybe"}, "LEGAL_REVIEWED"},
		{map[string]string{"LOG_RETENTION_DAYS": "0"}, "LOG_RETENTION_DAYS"},
		{map[string]string{"HSTS_PRELOAD": "yes please"}, "HSTS_PRELOAD"},
		{map[string]string{"CONTACT_TO": "a@b.example"}, "SMTP_HOST"},
		{map[string]string{"SMTP_HOST": "smtp.example.org"}, "CONTACT_TO"},
		{map[string]string{"SMTP_HOST": "smtp.example.org:25", "CONTACT_TO": "a@b.example", "CONTACT_FROM": "c@d.example"}, "SMTP_HOST"},
		{map[string]string{"SMTP_HOST": "smtp.example.org", "CONTACT_TO": "a@b.example", "CONTACT_FROM": "c@d.example", "SMTP_USER": "u"}, "SMTP_PASS"},
		{map[string]string{"SMTP_HOST": "smtp.example.org", "CONTACT_TO": "a@b.example", "CONTACT_FROM": "c@d.example", "SMTP_PORT": "x"}, "SMTP_PORT"},
		{map[string]string{"BEHIND_FRONT_DOOR": "true"}, "FRONT_DOOR_ID"},
	}
	for _, tc := range cases {
		set := map[string]string{}
		if tc.want == "https" || tc.want == "FRONT_DOOR_ID" {
			set = prodEnv()
		}
		for k, v := range tc.set {
			set[k] = v
		}
		_, err := Load(env(set))
		if err == nil || !strings.Contains(err.Error(), tc.want) {
			t.Errorf("%v: want error mentioning %s, got %v", tc.set, tc.want, err)
		}
	}
}

func TestSecretsAreRedacted(t *testing.T) {
	set := prodEnv()
	set["SMTP_HOST"] = "smtp.example.org"
	set["SMTP_USER"] = "web"
	set["SMTP_PASS"] = "smtp-password-value"
	set["CONTACT_TO"] = "team@example.org"
	set["CONTACT_FROM"] = "web@example.org"
	set["BEHIND_FRONT_DOOR"] = "true"
	set["FRONT_DOOR_ID"] = "front-door-id-value"
	c, err := Load(env(set))
	if err != nil {
		t.Fatal(err)
	}
	var buf bytes.Buffer
	log := slog.New(slog.NewJSONHandler(&buf, nil))
	log.Info("config", "config", c, "pass", c.SMTP.Pass, "smtp", c.SMTP)
	out := buf.String() + fmt.Sprintf("%v %+v %#v %s", c, c, c.SMTP, c.CSRFSecret)
	for _, s := range []string{secret, "smtp-password-value", "front-door-id-value"} {
		if strings.Contains(out, s) {
			t.Errorf("secret %q leaked into output", s)
		}
	}
	if c.SMTP.Pass.Reveal() != "smtp-password-value" {
		t.Error("Reveal must return the value")
	}

	// A failing configuration must not echo secrets either.
	set["CONTACT_TO"] = "invalid"
	_, err = Load(env(set))
	if err == nil || strings.Contains(err.Error(), "smtp-password-value") || strings.Contains(err.Error(), secret) {
		t.Errorf("error leaks a secret or is missing: %v", err)
	}
}

// TestEnvExample keeps .env.example in step with the variables the code
// reads, and makes sure it contains no real secret.
func TestEnvExample(t *testing.T) {
	b, err := os.ReadFile("../../.env.example")
	if err != nil {
		t.Fatal(err)
	}
	text := string(b)
	for _, v := range Vars {
		if !strings.Contains(text, "\n"+v.Name+"=") {
			t.Errorf(".env.example lacks %s", v.Name)
		}
	}
	for _, line := range strings.Split(text, "\n") {
		for _, name := range []string{"SMTP_PASS=", "CSRF_SECRET=", "FRONT_DOOR_ID="} {
			if strings.HasPrefix(line, name) && line != name {
				t.Errorf(".env.example must leave %s empty", name)
			}
		}
	}
}

func TestValidEmail(t *testing.T) {
	for _, s := range []string{"hello@tiefer.space", "first.last+tag@example.co.uk"} {
		if !ValidEmail(s) {
			t.Errorf("ValidEmail(%q) = false", s)
		}
	}
	for _, s := range []string{"", "a", "a@b", "a@b.c\n", "a b@c.de", "<a@b.de>", "a@b.de,c@d.de", "\"x\"@b.de", strings.Repeat("a", 250) + "@b.de"} {
		if ValidEmail(s) {
			t.Errorf("ValidEmail(%q) = true", s)
		}
	}
}
