// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

// Package config loads and validates the configuration of the website
// from environment variables. It is the only place that reads the
// environment.
package config

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"net/mail"
	"net/url"
	"os"
	"strconv"
	"strings"
)

// Environment names accepted in ENV.
const (
	Development = "development"
	Production  = "production"
)

// ProductionSiteURL is the canonical address of the public site. It is the
// default for SITE_URL in production.
const ProductionSiteURL = "https://tiefer.space"

// Legal holds the company details shown on the legal pages. Every value
// defaults to a clearly marked placeholder, never to invented data.
type Legal struct {
	Name            string // LEGAL_NAME
	Form            string // LEGAL_FORM
	Address         string // LEGAL_ADDRESS
	TaxID           string // LEGAL_TAX_ID (VÖEN)
	Registration    string // LEGAL_REGISTRATION
	Director        string // LEGAL_DIRECTOR
	HostingProvider string // LEGAL_HOSTING_PROVIDER
	HostingCountry  string // LEGAL_HOSTING_COUNTRY
	SMTPProvider    string // LEGAL_SMTP_PROVIDER
	SMTPCountry     string // LEGAL_SMTP_COUNTRY
	Reviewed        bool   // LEGAL_REVIEWED
}

// SMTP holds the mail delivery settings for the contact form.
type SMTP struct {
	Host string
	Port int
	User string
	Pass string
	// SkipVerify disables TLS certificate checks. Only allowed in
	// development, for local SMTP catchers with self-signed certificates.
	SkipVerify bool
}

// Config is the validated configuration.
type Config struct {
	Port         string
	Env          string
	SiteURL      string // without trailing slash
	ContactEmail string
	LinkedInURL  string
	RepoURL      string
	TrustProxy   bool

	SMTP        SMTP
	ContactTo   string
	ContactFrom string
	CSRFSecret  []byte

	Legal Legal
}

// Var describes one environment variable: its name and development default.
type Var struct {
	Name    string
	Default string
}

// Defaults lists the variables whose development defaults are placeholders
// that must be replaced before launch. Values that are empty by default
// (SMTP, REPO_URL) are optional and therefore not listed here.
var Defaults = []Var{
	{"SITE_URL", "http://localhost:8080"},
	{"CONTACT_EMAIL", "hello@example.com"},
	{"LINKEDIN_URL", "https://www.linkedin.com/company/tiefer"},
	{"LEGAL_NAME", "[COMPANY_LEGAL_NAME]"},
	{"LEGAL_FORM", "[LEGAL_FORM]"},
	{"LEGAL_ADDRESS", "[REGISTERED_ADDRESS_BAKU]"},
	{"LEGAL_TAX_ID", "[VOEN]"},
	{"LEGAL_REGISTRATION", "[STATE_REGISTRATION_DETAILS]"},
	{"LEGAL_DIRECTOR", "[MANAGING_DIRECTOR]"},
	{"LEGAL_HOSTING_PROVIDER", "[HOSTING_PROVIDER]"},
	{"LEGAL_HOSTING_COUNTRY", "[HOSTING_COUNTRY]"},
	{"LEGAL_SMTP_PROVIDER", "[SMTP_PROVIDER]"},
	{"LEGAL_SMTP_COUNTRY", "[SMTP_COUNTRY]"},
}

func defaultOf(name string) string {
	for _, v := range Defaults {
		if v.Name == name {
			return v.Default
		}
	}
	return ""
}

// Load reads the configuration from the process environment.
func Load() (*Config, error) {
	return FromLookup(os.LookupEnv)
}

// FromLookup reads the configuration through lookup, which has the
// signature of os.LookupEnv. Tests pass a map-backed function.
func FromLookup(lookup func(string) (string, bool)) (*Config, error) {
	var errs []error
	get := func(name string) string {
		v, _ := lookup(name)
		return strings.TrimSpace(v)
	}
	c := &Config{}
	c.Env = get("ENV")
	if c.Env == "" {
		c.Env = Development
	}
	if c.Env != Development && c.Env != Production {
		errs = append(errs, fmt.Errorf("ENV must be %q or %q, got %q", Development, Production, c.Env))
	}
	prod := c.Env == Production
	// orDefault returns the value or its development default. In
	// production, required variables get no default.
	orDefault := func(name string, required bool) string {
		v := get(name)
		if v != "" {
			return v
		}
		if prod && required {
			errs = append(errs, fmt.Errorf("%s is required in production", name))
			return ""
		}
		return defaultOf(name)
	}

	c.Port = get("PORT")
	if c.Port == "" {
		c.Port = "8080"
	}
	if p, err := strconv.Atoi(c.Port); err != nil || p < 1 || p > 65535 {
		errs = append(errs, fmt.Errorf("PORT must be a number between 1 and 65535, got %q", c.Port))
	}

	c.SiteURL = strings.TrimRight(get("SITE_URL"), "/")
	if c.SiteURL == "" {
		c.SiteURL = defaultOf("SITE_URL")
		if prod {
			c.SiteURL = ProductionSiteURL
		}
	}
	if c.SiteURL != "" {
		u, err := url.Parse(c.SiteURL)
		switch {
		case err != nil || u.Host == "" || (u.Scheme != "http" && u.Scheme != "https"):
			errs = append(errs, fmt.Errorf("SITE_URL must be an absolute http or https URL, got %q", c.SiteURL))
		case u.Path != "" || u.RawQuery != "" || u.Fragment != "":
			errs = append(errs, fmt.Errorf("SITE_URL must not contain a path, query or fragment, got %q", c.SiteURL))
		case prod && u.Scheme != "https":
			errs = append(errs, fmt.Errorf("SITE_URL must use https in production, got %q", c.SiteURL))
		}
	}

	c.ContactEmail = orDefault("CONTACT_EMAIL", true)
	if c.ContactEmail != "" && !validEmail(c.ContactEmail) {
		errs = append(errs, fmt.Errorf("CONTACT_EMAIL is not a valid email address: %q", c.ContactEmail))
	}
	c.LinkedInURL = orDefault("LINKEDIN_URL", false)
	if err := checkHTTPURL("LINKEDIN_URL", c.LinkedInURL); err != nil {
		errs = append(errs, err)
	}
	c.RepoURL = get("REPO_URL")
	if c.RepoURL != "" {
		if err := checkHTTPURL("REPO_URL", c.RepoURL); err != nil {
			errs = append(errs, err)
		}
	}
	var err error
	if c.TrustProxy, err = parseBool("TRUST_PROXY", get("TRUST_PROXY")); err != nil {
		errs = append(errs, err)
	}

	// Contact form delivery. All or nothing: a partial setup is an error.
	c.SMTP.Host = get("SMTP_HOST")
	c.SMTP.User = get("SMTP_USER")
	c.SMTP.Pass, _ = lookup("SMTP_PASS") // passwords may contain spaces
	c.ContactTo = get("CONTACT_TO")
	c.ContactFrom = get("CONTACT_FROM")
	if c.SMTP.Host != "" || c.ContactTo != "" || c.ContactFrom != "" {
		if c.SMTP.Host == "" {
			errs = append(errs, errors.New("SMTP_HOST is required when CONTACT_TO or CONTACT_FROM is set"))
		}
		if !validEmail(c.ContactTo) {
			errs = append(errs, fmt.Errorf("CONTACT_TO must be a valid email address when SMTP is configured, got %q", c.ContactTo))
		}
		if !validEmail(c.ContactFrom) {
			errs = append(errs, fmt.Errorf("CONTACT_FROM must be a valid email address when SMTP is configured, got %q", c.ContactFrom))
		}
		port := get("SMTP_PORT")
		if port == "" {
			port = "587"
		}
		if c.SMTP.Port, err = strconv.Atoi(port); err != nil || c.SMTP.Port < 1 || c.SMTP.Port > 65535 {
			errs = append(errs, fmt.Errorf("SMTP_PORT must be a port number, got %q", port))
		}
		if (c.SMTP.User == "") != (c.SMTP.Pass == "") {
			errs = append(errs, errors.New("SMTP_USER and SMTP_PASS must be set together"))
		}
	}
	if c.SMTP.SkipVerify, err = parseBool("SMTP_SKIP_VERIFY", get("SMTP_SKIP_VERIFY")); err != nil {
		errs = append(errs, err)
	}
	if prod && c.SMTP.SkipVerify {
		errs = append(errs, errors.New("SMTP_SKIP_VERIFY is not allowed in production"))
	}

	secret := get("CSRF_SECRET")
	switch {
	case secret == "" && prod:
		errs = append(errs, errors.New("CSRF_SECRET is required in production (at least 32 characters, for example: openssl rand -hex 32)"))
	case secret == "":
		b := make([]byte, 32)
		_, _ = rand.Read(b)
		c.CSRFSecret = []byte(hex.EncodeToString(b))
	case len(secret) < 32:
		errs = append(errs, errors.New("CSRF_SECRET must be at least 32 characters"))
	default:
		c.CSRFSecret = []byte(secret)
	}

	c.Legal = Legal{
		Name:            orDefault("LEGAL_NAME", false),
		Form:            orDefault("LEGAL_FORM", false),
		Address:         orDefault("LEGAL_ADDRESS", false),
		TaxID:           orDefault("LEGAL_TAX_ID", false),
		Registration:    orDefault("LEGAL_REGISTRATION", false),
		Director:        orDefault("LEGAL_DIRECTOR", false),
		HostingProvider: orDefault("LEGAL_HOSTING_PROVIDER", false),
		HostingCountry:  orDefault("LEGAL_HOSTING_COUNTRY", false),
		SMTPProvider:    orDefault("LEGAL_SMTP_PROVIDER", false),
		SMTPCountry:     orDefault("LEGAL_SMTP_COUNTRY", false),
	}
	if c.Legal.Reviewed, err = parseBool("LEGAL_REVIEWED", get("LEGAL_REVIEWED")); err != nil {
		errs = append(errs, err)
	}

	if len(errs) > 0 {
		return nil, fmt.Errorf("invalid configuration:\n  %w", joinLines(errs))
	}
	return c, nil
}

// SiteHost returns the host name of SITE_URL, including a port if it has
// one, for example tiefer.space or localhost:8080.
func (c *Config) SiteHost() string {
	u, err := url.Parse(c.SiteURL)
	if err != nil {
		return ""
	}
	return u.Host
}

// ContactEnabled reports whether the contact form can deliver mail.
func (c *Config) ContactEnabled() bool {
	return c.SMTP.Host != "" && c.ContactTo != "" && c.ContactFrom != ""
}

// Production reports whether ENV is production.
func (c *Config) Production() bool { return c.Env == Production }

// DefaultsInUse returns the names of variables that still carry their
// development placeholder value.
func (c *Config) DefaultsInUse() []string {
	values := map[string]string{
		"SITE_URL":               c.SiteURL,
		"CONTACT_EMAIL":          c.ContactEmail,
		"LINKEDIN_URL":           c.LinkedInURL,
		"LEGAL_NAME":             c.Legal.Name,
		"LEGAL_FORM":             c.Legal.Form,
		"LEGAL_ADDRESS":          c.Legal.Address,
		"LEGAL_TAX_ID":           c.Legal.TaxID,
		"LEGAL_REGISTRATION":     c.Legal.Registration,
		"LEGAL_DIRECTOR":         c.Legal.Director,
		"LEGAL_HOSTING_PROVIDER": c.Legal.HostingProvider,
		"LEGAL_HOSTING_COUNTRY":  c.Legal.HostingCountry,
		"LEGAL_SMTP_PROVIDER":    c.Legal.SMTPProvider,
		"LEGAL_SMTP_COUNTRY":     c.Legal.SMTPCountry,
	}
	var out []string
	for _, v := range Defaults {
		if v.Name == "SITE_URL" && c.Production() {
			continue // the production default is the real domain
		}
		if values[v.Name] == v.Default {
			out = append(out, v.Name)
		}
	}
	if !c.Legal.Reviewed {
		out = append(out, "LEGAL_REVIEWED")
	}
	return out
}

func validEmail(s string) bool {
	if s == "" || len(s) > 254 || strings.ContainsAny(s, "\r\n") {
		return false
	}
	a, err := mail.ParseAddress(s)
	return err == nil && a.Address == s && a.Name == ""
}

func checkHTTPURL(name, v string) error {
	u, err := url.Parse(v)
	if err != nil || u.Host == "" || (u.Scheme != "http" && u.Scheme != "https") {
		return fmt.Errorf("%s must be an absolute http or https URL, got %q", name, v)
	}
	return nil
}

func parseBool(name, v string) (bool, error) {
	if v == "" {
		return false, nil
	}
	b, err := strconv.ParseBool(v)
	if err != nil {
		return false, fmt.Errorf("%s must be true or false, got %q", name, v)
	}
	return b, nil
}

type lines []error

func (l lines) Error() string {
	s := make([]string, len(l))
	for i, e := range l {
		s[i] = e.Error()
	}
	return strings.Join(s, "\n  ")
}

func (l lines) Unwrap() []error { return l }

func joinLines(errs []error) error { return lines(errs) }
