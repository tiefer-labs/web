// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package server

import (
	"net/http"
	"strings"
	"testing"
	"time"
)

var legalPaths = []string{"/legal", "/privacy", "/acceptable-use"}

func TestLegalPages(t *testing.T) {
	note := "Placeholder text. To be reviewed before launch."
	h := newHarness(t, nil)
	for _, p := range legalPaths {
		w := h.get(p)
		body := w.Body.String()
		if w.Code != http.StatusOK || strings.Count(body, "<h1") != 1 || !strings.Contains(body, note) {
			t.Errorf("%s: %d, h1 %d, note %v", p, w.Code, strings.Count(body, "<h1"), strings.Contains(body, note))
		}
		if !strings.Contains(body, `class="site-header"`) || !strings.Contains(body, "tiefer-logo.") {
			t.Errorf("%s: legal pages use the light header", p)
		}
		if !strings.Contains(body, `aria-current="page"`) {
			t.Errorf("%s: footer does not mark the current page", p)
		}
		if strings.Contains(body, "Impressum") {
			t.Errorf("%s: no German-specific terms", p)
		}
	}
	notice := h.get("/legal").Body.String()
	for _, want := range []string{
		`<mark class="placeholder">[COMPANY_LEGAL_NAME]</mark>`,
		`<mark class="placeholder">[VOEN]</mark>`,
		`<a href="mailto:hello@tiefer.space">hello@tiefer.space</a>`,
		`<a href="http://localhost:8080/">http://localhost:8080</a>`,
	} {
		if !strings.Contains(notice, want) {
			t.Errorf("/legal lacks %q", want)
		}
	}
	privacy := h.get("/privacy").Body.String()
	for _, want := range []string{"Microsoft Azure", "Azure Front Door", "web application firewall", "App Service",
		"for 30 days", "Law of the Republic of Azerbaijan on Personal Data of 11 May 2010", "GDPR",
		`[LEGAL_BASIS_HOSTING]`, `[RETENTION_PERIOD]`, `href="#how-long-we-keep-it"`} {
		if !strings.Contains(privacy, want) {
			t.Errorf("/privacy lacks %q", want)
		}
	}

	h = newHarness(t, map[string]string{"LEGAL_REVIEWED": "true", "LEGAL_NAME": "Tiefer MMC", "LEGAL_TAX_ID": "1234567890", "LOG_RETENTION_DAYS": "14"})
	notice = h.get("/legal").Body.String()
	if strings.Contains(notice, note) || !strings.Contains(notice, "Tiefer MMC") || !strings.Contains(notice, "1234567890") {
		t.Error("configured legal values or LEGAL_REVIEWED=true not applied")
	}
	if !strings.Contains(h.get("/privacy").Body.String(), "for 14 days") {
		t.Error("LOG_RETENTION_DAYS not shown")
	}
}

func TestLegalValuesAreEscaped(t *testing.T) {
	h := newHarness(t, map[string]string{"LEGAL_NAME": `<img src=x onerror=alert(1)>`})
	body := h.get("/legal").Body.String()
	if strings.Contains(body, "<img src=x") || !strings.Contains(body, "&lt;img src=x onerror=alert(1)&gt;") {
		t.Error("configuration values must be escaped on legal pages")
	}
}

func TestFill(t *testing.T) {
	parts := fill("Ask {legal_name} at {contact_email} [RETENTION_PERIOD].", map[string]string{"legal_name": "[COMPANY_LEGAL_NAME]", "contact_email": "a@b.de"})
	var got []string
	for _, p := range parts {
		switch {
		case p.Href != "":
			got = append(got, "link:"+p.Text)
		case p.Mark:
			got = append(got, "mark:"+p.Text)
		default:
			got = append(got, p.Text)
		}
	}
	want := "Ask |mark:[COMPANY_LEGAL_NAME]| at |link:a@b.de| |mark:[RETENTION_PERIOD]|."
	if strings.Join(got, "|") != want {
		t.Errorf("fill = %q, want %q", strings.Join(got, "|"), want)
	}
}

func TestSecurityTxt(t *testing.T) {
	h := newHarness(t, map[string]string{"SITE_URL": "https://tiefer.space", "REPO_URL": "https://github.com/tiefer-labs/web"})
	w := h.get("/.well-known/security.txt")
	body := w.Body.String()
	if w.Code != http.StatusOK || !strings.HasPrefix(w.Header().Get("Content-Type"), "text/plain") {
		t.Fatalf("security.txt: %d %q", w.Code, w.Header().Get("Content-Type"))
	}
	for _, want := range []string{"Contact: mailto:hello@tiefer.space\n", "Canonical: https://tiefer.space/.well-known/security.txt\n", "Preferred-Languages: en\n", "Policy: https://github.com/tiefer-labs/web/blob/main/SECURITY.md\n"} {
		if !strings.Contains(body, want) {
			t.Errorf("security.txt lacks %q", want)
		}
	}
	if w := h.get("/security.txt"); w.Code != http.StatusMovedPermanently || w.Header().Get("Location") != "/.well-known/security.txt" {
		t.Errorf("/security.txt: %d", w.Code)
	}
}

// TestSecurityTxtExpiry fails 30 days before Expires, and when Expires is
// more than a year ahead (RFC 9116 recommends less than a year).
func TestSecurityTxtExpiry(t *testing.T) {
	now := time.Now()
	if SecurityTxtExpires.Before(now.Add(30 * 24 * time.Hour)) {
		t.Errorf("security.txt expires on %s: confirm the contact and move SecurityTxtExpires forward", SecurityTxtExpires.Format(time.DateOnly))
	}
	if SecurityTxtExpires.After(now.Add(366 * 24 * time.Hour)) {
		t.Errorf("security.txt Expires %s is more than a year ahead", SecurityTxtExpires.Format(time.DateOnly))
	}
}
