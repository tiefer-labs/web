// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

// Package config reads the configuration from environment variables and
// validates it once at startup. No other package reads the environment.
package config

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"log/slog"
	"net/mail"
	"net/url"
	"strconv"
	"strings"
)

// Environments accepted in ENV.
const (
	Development = "development"
	Production  = "production"
)

// Secret holds a value that must never appear in logs, errors or
// rendered output. Its String and LogValue methods return a fixed mask;
// the value itself is only available through Reveal.
type Secret struct{ v string }

// NewSecret wraps a secret value.
func NewSecret(v string) Secret { return Secret{v: v} }

// Reveal returns the secret value, for the one place that needs it.
func (s Secret) Reveal() string { return s.v }

// Set reports whether the secret has a value.
func (s Secret) Set() bool { return s.v != "" }

// String masks the value, so fmt verbs never print it.
func (s Secret) String() string {
	if s.v == "" {
		return ""
	}
	return "[redacted]"
}

// GoString masks the value for %#v.
func (s Secret) GoString() string { return s.String() }

// LogValue masks the value for log/slog.
func (s Secret) LogValue() slog.Value { return slog.StringValue(s.String()) }

// MarshalText masks the value for encoders.
func (s Secret) MarshalText() ([]byte, error) { return []byte(s.String()), nil }

// SMTP holds the mail delivery settings of the contact form.
type SMTP struct {
	Host string
	Port int
	User string
	Pass Secret
	// SkipVerify accepts any TLS certificate. Development only, for a
	// local catcher such as Mailpit; refused in production.
	SkipVerify bool
}

// Legal holds the company details shown on the legal pages. Unset values
// default to clearly marked placeholders, never to invented data.
type Legal struct {
	Name            string // LEGAL_NAME
	Form            string // LEGAL_FORM
	Address         string // LEGAL_ADDRESS
	TaxID           string // LEGAL_TAX_ID (VOEN)
	Registration    string // LEGAL_REGISTRATION
	Director        string // LEGAL_DIRECTOR
	HostingProvider string // LEGAL_HOSTING_PROVIDER
	HostingCountry  string // LEGAL_HOSTING_COUNTRY
	SMTPProvider    string // LEGAL_SMTP_PROVIDER
	SMTPCountry     string // LEGAL_SMTP_COUNTRY
	Reviewed        bool   // LEGAL_REVIEWED
}

// Config is the validated configuration.
type Config struct {
	Port         int
	Env          string
	SiteURL      *url.URL // scheme and host only
	ContactEmail string
	LinkedInURL  string
	RepoURL      string

	SMTP        SMTP
	ContactTo   string
	ContactFrom string
	CSRFSecret  Secret

	Legal            Legal
	LogRetentionDays int

	BehindFrontDoor bool
	FrontDoorID     Secret
	HSTSPreload     bool

	// defaults lists the variables that still carry a development default
	// or a placeholder, for the startup warning.
	defaults []string
}

// Var is one environment variable with its development default.
type Var struct {
	Name    string
	Default string
	Comment string
}

// Vars lists every variable the site reads, in the order of .env.example.
// A test keeps .env.example in step with this list.
var Vars = []Var{
	{"PORT", "8080", "HTTP port."},
	{"ENV", Development, "development or production."},
	{"SITE_URL", "http://localhost:8080", "Canonical base URL. Production: https://tiefer.space"},
	{"CONTACT_EMAIL", "hello@tiefer.space", "Public contact address and mailto link."},
	{"LINKEDIN_URL", "https://www.linkedin.com/company/tiefer/", "LinkedIn company page."},
	{"REPO_URL", "", "Public source repository, linked in the footer when set."},
	{"SMTP_HOST", "", "Mail server for the contact form. Empty hides the form."},
	{"SMTP_PORT", "587", "465 uses implicit TLS; any other port must offer STARTTLS."},
	{"SMTP_USER", "", "SMTP user name."},
	{"SMTP_PASS", "", "SMTP password. In production it comes from Azure Key Vault."},
	{"CONTACT_TO", "", "Where contact form messages are delivered."},
	{"CONTACT_FROM", "", "Sender address of contact form messages."},
	{"SMTP_SKIP_VERIFY", "false", "Development only: accept the self-signed certificate of a local SMTP catcher such as Mailpit. Refused in production."},
	{"CSRF_SECRET", "", "Signing key for form tokens, at least 32 characters. Required in production; random in development."},
	{"LEGAL_NAME", "[COMPANY_LEGAL_NAME]", "Company legal name."},
	{"LEGAL_FORM", "[LEGAL_FORM]", "Legal form, for example MMC (limited liability company under the laws of the Republic of Azerbaijan)."},
	{"LEGAL_ADDRESS", "[REGISTERED_ADDRESS_BAKU]", "Registered address in Baku."},
	{"LEGAL_TAX_ID", "[VOEN]", "Tax identification number (VOEN)."},
	{"LEGAL_REGISTRATION", "[STATE_REGISTRATION_DETAILS]", "State registration details."},
	{"LEGAL_DIRECTOR", "[MANAGING_DIRECTOR]", "Managing director."},
	{"LEGAL_HOSTING_PROVIDER", "[HOSTING_PROVIDER]", "Hosting provider, named in the privacy notice."},
	{"LEGAL_HOSTING_COUNTRY", "[HOSTING_COUNTRY]", "Country where the site is hosted."},
	{"LEGAL_SMTP_PROVIDER", "[SMTP_PROVIDER]", "Email provider, named in the privacy notice."},
	{"LEGAL_SMTP_COUNTRY", "[SMTP_COUNTRY]", "Country of the email provider."},
	{"LEGAL_REVIEWED", "false", "true only after a lawyer has reviewed the legal pages."},
	{"LOG_RETENTION_DAYS", "30", "Shown on the privacy page; must match the Azure log retention."},
	{"BEHIND_FRONT_DOOR", "false", "true on Azure: requests must carry the Front Door ID, and the client address comes from Front Door."},
	{"FRONT_DOOR_ID", "", "The Front Door profile ID. Required when BEHIND_FRONT_DOOR=true in production."},
	{"HSTS_PRELOAD", "false", "Adds preload to Strict-Transport-Security. Founder decision."},
}

func defaultOf(name string) string {
	for _, v := range Vars {
		if v.Name == name {
			return v.Default
		}
	}
	panic("config: unknown variable " + name)
}

// Load reads the configuration through lookup, which has the signature
// of os.LookupEnv. All problems are reported together.
func Load(lookup func(string) (string, bool)) (*Config, error) {
	var errs []error
	c := &Config{}
	get := func(name string) string {
		v, ok := lookup(name)
		if !ok || strings.TrimSpace(v) == "" {
			return ""
		}
		return strings.TrimSpace(v)
	}
	// value returns the variable or its default, and remembers defaults
	// that are placeholders or development values.
	value := func(name string) string {
		if v := get(name); v != "" {
			return v
		}
		d := defaultOf(name)
		if strings.HasPrefix(d, "[") {
			c.defaults = append(c.defaults, name)
		}
		return d
	}
	boolean := func(name string) bool {
		s := value(name)
		b, err := strconv.ParseBool(s)
		if err != nil {
			errs = append(errs, fmt.Errorf("%s must be true or false, got %q", name, s))
		}
		return b
	}

	c.Env = value("ENV")
	if c.Env != Development && c.Env != Production {
		errs = append(errs, fmt.Errorf("ENV must be %q or %q, got %q", Development, Production, c.Env))
	}
	prod := c.Env == Production

	port, err := strconv.Atoi(value("PORT"))
	if err != nil || port < 1 || port > 65535 {
		errs = append(errs, fmt.Errorf("PORT must be a number from 1 to 65535"))
	}
	c.Port = port

	site := get("SITE_URL")
	if site == "" {
		if prod {
			errs = append(errs, errors.New("SITE_URL is required in production (https://tiefer.space)"))
		}
		site = defaultOf("SITE_URL")
		c.defaults = append(c.defaults, "SITE_URL")
	}
	u, err := url.Parse(strings.TrimRight(site, "/"))
	switch {
	case err != nil || u.Host == "" || (u.Scheme != "http" && u.Scheme != "https"):
		errs = append(errs, fmt.Errorf("SITE_URL must be an absolute http or https URL, got %q", site))
	case u.Path != "" || u.RawQuery != "" || u.Fragment != "" || u.User != nil:
		errs = append(errs, fmt.Errorf("SITE_URL must contain only a scheme and a host, got %q", site))
	case prod && u.Scheme != "https":
		errs = append(errs, errors.New("SITE_URL must use https in production"))
	}
	c.SiteURL = u

	c.ContactEmail = value("CONTACT_EMAIL")
	if !ValidEmail(c.ContactEmail) {
		errs = append(errs, fmt.Errorf("CONTACT_EMAIL is not a valid address: %q", c.ContactEmail))
	}
	c.LinkedInURL = value("LINKEDIN_URL")
	if err := checkLink("LINKEDIN_URL", c.LinkedInURL); err != nil {
		errs = append(errs, err)
	}
	c.RepoURL = get("REPO_URL")
	if c.RepoURL != "" {
		if err := checkLink("REPO_URL", c.RepoURL); err != nil {
			errs = append(errs, err)
		}
	}

	// Contact form delivery. SMTP_HOST switches it on; then every
	// setting it needs must be valid.
	c.SMTP.Host = get("SMTP_HOST")
	c.SMTP.User = get("SMTP_USER")
	if v, ok := lookup("SMTP_PASS"); ok {
		c.SMTP.Pass = NewSecret(v) // not trimmed: passwords may contain spaces
	}
	c.ContactTo = get("CONTACT_TO")
	c.ContactFrom = get("CONTACT_FROM")
	c.SMTP.Port, err = strconv.Atoi(value("SMTP_PORT"))
	if err != nil || c.SMTP.Port < 1 || c.SMTP.Port > 65535 {
		errs = append(errs, errors.New("SMTP_PORT must be a number from 1 to 65535"))
	}
	if c.SMTP.Host != "" {
		if strings.ContainsAny(c.SMTP.Host, "/: \t") {
			errs = append(errs, errors.New("SMTP_HOST must be a host name without scheme or port"))
		}
		if !ValidEmail(c.ContactTo) {
			errs = append(errs, errors.New("CONTACT_TO must be a valid address when SMTP_HOST is set"))
		}
		if !ValidEmail(c.ContactFrom) {
			errs = append(errs, errors.New("CONTACT_FROM must be a valid address when SMTP_HOST is set"))
		}
		if (c.SMTP.User == "") != !c.SMTP.Pass.Set() {
			errs = append(errs, errors.New("SMTP_USER and SMTP_PASS must be set together"))
		}
	} else if c.ContactTo != "" || c.ContactFrom != "" || c.SMTP.User != "" || c.SMTP.Pass.Set() {
		errs = append(errs, errors.New("SMTP_HOST is required when other SMTP or CONTACT_ settings are set"))
	}

	c.SMTP.SkipVerify = boolean("SMTP_SKIP_VERIFY")
	if prod && c.SMTP.SkipVerify {
		errs = append(errs, errors.New("SMTP_SKIP_VERIFY is not allowed in production"))
	}

	// An unresolved Key Vault reference reaches the app as its literal
	// text. For CSRF_SECRET that text would be a public, known key.
	for _, name := range []string{"CSRF_SECRET", "SMTP_PASS"} {
		if v, _ := lookup(name); strings.HasPrefix(strings.TrimSpace(v), "@Microsoft.KeyVault(") {
			errs = append(errs, fmt.Errorf("%s is an unresolved Key Vault reference: check the web app identity and the secret in Key Vault", name))
		}
	}

	switch secret := get("CSRF_SECRET"); {
	case secret == "" && prod:
		errs = append(errs, errors.New("CSRF_SECRET is required in production (at least 32 characters: openssl rand -hex 32)"))
	case secret == "":
		b := make([]byte, 32)
		_, _ = rand.Read(b)
		c.CSRFSecret = NewSecret(hex.EncodeToString(b))
	case len(secret) < 32:
		errs = append(errs, errors.New("CSRF_SECRET must be at least 32 characters"))
	case weakSecret(secret):
		errs = append(errs, errors.New("CSRF_SECRET looks guessable (a repeated pattern or a placeholder); generate one with: openssl rand -hex 32"))
	default:
		c.CSRFSecret = NewSecret(secret)
	}

	c.Legal = Legal{
		Name:            value("LEGAL_NAME"),
		Form:            value("LEGAL_FORM"),
		Address:         value("LEGAL_ADDRESS"),
		TaxID:           value("LEGAL_TAX_ID"),
		Registration:    value("LEGAL_REGISTRATION"),
		Director:        value("LEGAL_DIRECTOR"),
		HostingProvider: value("LEGAL_HOSTING_PROVIDER"),
		HostingCountry:  value("LEGAL_HOSTING_COUNTRY"),
		SMTPProvider:    value("LEGAL_SMTP_PROVIDER"),
		SMTPCountry:     value("LEGAL_SMTP_COUNTRY"),
		Reviewed:        boolean("LEGAL_REVIEWED"),
	}
	if !c.Legal.Reviewed {
		c.defaults = append(c.defaults, "LEGAL_REVIEWED")
	}
	days, err := strconv.Atoi(value("LOG_RETENTION_DAYS"))
	if err != nil || days < 1 || days > 3650 {
		errs = append(errs, errors.New("LOG_RETENTION_DAYS must be a number of days from 1 to 3650"))
	}
	c.LogRetentionDays = days

	c.BehindFrontDoor = boolean("BEHIND_FRONT_DOOR")
	c.FrontDoorID = NewSecret(get("FRONT_DOOR_ID"))
	if c.BehindFrontDoor && prod && !c.FrontDoorID.Set() {
		errs = append(errs, errors.New("FRONT_DOOR_ID is required when BEHIND_FRONT_DOOR=true in production"))
	}
	c.HSTSPreload = boolean("HSTS_PRELOAD")

	if len(errs) > 0 {
		return nil, fmt.Errorf("invalid configuration: %w", errors.Join(errs...))
	}
	return c, nil
}

// Production reports whether ENV is production.
func (c *Config) Production() bool { return c.Env == Production }

// ContactEnabled reports whether the contact form can deliver mail.
func (c *Config) ContactEnabled() bool { return c.SMTP.Host != "" }

// Site returns SITE_URL without a trailing slash.
func (c *Config) Site() string { return c.SiteURL.String() }

// DefaultsInUse returns the variables that still carry a placeholder or
// development default.
func (c *Config) DefaultsInUse() []string { return append([]string(nil), c.defaults...) }

// LogValue describes the configuration for the startup log line without
// any secret or personal value.
func (c *Config) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("env", c.Env),
		slog.Int("port", c.Port),
		slog.String("site", c.Site()),
		slog.Bool("contact_form", c.ContactEnabled()),
		slog.Bool("behind_front_door", c.BehindFrontDoor),
		slog.Bool("legal_reviewed", c.Legal.Reviewed),
	)
}

// ValidEmail reports whether s is a single plain address, without a
// display name and without characters that could inject mail headers.
func ValidEmail(s string) bool {
	if s == "" || len(s) > 254 || strings.ContainsAny(s, "\r\n<>\"(),;: \t") {
		return false
	}
	a, err := mail.ParseAddress(s)
	return err == nil && a.Name == "" && a.Address == s && strings.Contains(s[strings.LastIndexByte(s, '@')+1:], ".")
}

func checkLink(name, s string) error {
	u, err := url.Parse(s)
	if err != nil || u.Scheme != "https" || u.Host == "" || u.User != nil {
		return fmt.Errorf("%s must be an absolute https URL, got %q", name, s)
	}
	return nil
}

// weakSecret reports secrets that are long enough but guessable: few
// distinct characters (such as "aaaa..." or "abab..."), or an obvious
// placeholder word. Random hex from openssl always passes.
func weakSecret(s string) bool {
	distinct := map[rune]bool{}
	for _, r := range s {
		distinct[r] = true
	}
	if len(distinct) < 8 {
		return true
	}
	lower := strings.ToLower(s)
	for _, w := range []string{"changeme", "change-me", "change_me", "secret", "password", "example", "placeholder", "0123456789abcdef0123"} {
		if strings.Contains(lower, w) {
			return true
		}
	}
	return false
}
