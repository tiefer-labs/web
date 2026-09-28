// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package server

import (
	"bytes"
	"compress/gzip"
	"html"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"regexp"
	"strings"
	"testing"

	"github.com/tiefer-labs/web/internal/content"
)

func TestIndexRenders(t *testing.T) {
	h := newHarness(t, nil)
	w := h.get("/")
	if w.Code != http.StatusOK || !strings.HasPrefix(w.Header().Get("Content-Type"), "text/html") {
		t.Fatalf("/: %d %q", w.Code, w.Header().Get("Content-Type"))
	}
	body := w.Body.String()
	if n := strings.Count(body, "<h1"); n != 1 {
		t.Errorf("index has %d h1 elements, want 1", n)
	}
	if w.Header().Get("Cache-Control") != "private, no-cache" {
		t.Errorf("Cache-Control = %q", w.Header().Get("Cache-Control"))
	}
	for _, id := range []string{`id="product"`, `id="how"`, `id="use-cases"`, `id="principles"`, `id="contact"`, `id="main"`} {
		if !strings.Contains(body, id) {
			t.Errorf("index lacks %s", id)
		}
	}
}

// TestAllCopyFromContent walks the English content of the index page and
// checks that every string appears in the rendered page, so the page is
// built from en.go and nothing is hardcoded instead.
func TestAllCopyFromContent(t *testing.T) {
	h := newHarness(t, map[string]string{"SMTP_HOST": "smtp.example.org", "CONTACT_TO": "a@example.org", "CONTACT_FROM": "b@example.org", "REPO_URL": "https://github.com/tiefer-labs/web"})
	body := html.UnescapeString(h.get("/").Body.String())
	site := content.Default()
	skip := map[string]bool{
		// Shown only in other states or on other pages.
		"Success": true, "Invalid": true, "Expired": true, "Limited": true, "ErrorBefore": true, "ErrorAfter": true,
		"EmailUs": true, "Sending": true, "Required": true, "CharactersLabel": true,
		"NameMissing": true, "NameTooLong": true, "EmailMissing": true, "EmailInvalid": true, "OrgTooLong": true,
		"RoleInvalid": true, "MessageMissing": true, "MessageTooLong": true, "ConsentMissing": true,
		"LegalLabel": true, "BackHome": true, "MenuClose": true, "Prefix": true, "Code": true,
	}
	for _, part := range []any{site.Meta, site.UI, site.Nav, site.CTA, site.Hero, site.Problem, site.Product, site.How,
		site.Demo, site.UseCases, site.Principles, site.Roadmap, site.Name, site.Contact, site.Footer} {
		eachString(reflect.ValueOf(part), "", func(field, s string) {
			if skip[field] || s == "" || strings.HasPrefix(s, "#") {
				return
			}
			if !strings.Contains(body, s) {
				t.Errorf("index does not show %s: %q", field, s)
			}
		})
	}
}

func eachString(v reflect.Value, field string, fn func(field, s string)) {
	switch v.Kind() {
	case reflect.String:
		fn(field, v.String())
	case reflect.Struct:
		for i := 0; i < v.NumField(); i++ {
			eachString(v.Field(i), v.Type().Field(i).Name, fn)
		}
	case reflect.Slice:
		for i := 0; i < v.Len(); i++ {
			eachString(v.Index(i), field, fn)
		}
	}
}

func TestNotFoundPage(t *testing.T) {
	h := newHarness(t, nil)
	w := h.get("/nothing-here")
	body := w.Body.String()
	if w.Code != http.StatusNotFound || !strings.Contains(body, content.Default().NotFound.Title) || !strings.Contains(body, `<meta name="robots" content="noindex">`) {
		t.Errorf("404: %d", w.Code)
	}
	if strings.Count(body, "<h1") != 1 || strings.Contains(body, `rel="canonical"`) {
		t.Error("404 page needs one h1 and no canonical link")
	}
	if w := h.do(httptest.NewRequest(http.MethodPost, "/nothing-here", nil)); w.Code != http.StatusNotFound || strings.Contains(w.Body.String(), "<html") {
		t.Errorf("POST to an unknown path: %d, want a short text 404", w.Code)
	}
}

func TestStaticAssets(t *testing.T) {
	h := newHarness(t, nil)
	css := h.srv.assetURL("css/site.css")
	if !regexp.MustCompile(`^/static/css/site\.[0-9a-f]{12}\.css$`).MatchString(css) {
		t.Fatalf("hashed URL %q", css)
	}
	w := h.get(css)
	if w.Code != http.StatusOK || w.Header().Get("Cache-Control") != "public, max-age=31536000, immutable" || w.Header().Get("ETag") == "" {
		t.Errorf("%s: %d %q", css, w.Code, w.Header().Get("Cache-Control"))
	}
	r := httptest.NewRequest(http.MethodGet, css, nil)
	r.Header.Set("If-None-Match", w.Header().Get("ETag"))
	if w := h.do(r); w.Code != http.StatusNotModified || w.Body.Len() != 0 {
		t.Errorf("If-None-Match: %d", w.Code)
	}
	for _, p := range []string{"/static/css/site.css", "/static/css/site.000000000000.css", "/static/../go.mod", "/static/"} {
		switch w := h.get(p); w.Code {
		case http.StatusNotFound, http.StatusMovedPermanently, http.StatusTemporaryRedirect:
			if loc := w.Header().Get("Location"); strings.HasPrefix(loc, "/static/") && strings.Contains(loc, "..") {
				t.Errorf("%s redirects to %q", p, loc)
			}
		default:
			t.Errorf("%s: %d, want 404 or a clean redirect", p, w.Code)
		}
	}
	for _, p := range []string{"/favicon.ico", "/apple-touch-icon.png"} {
		if w := h.get(p); w.Code != http.StatusOK || w.Header().Get("Cache-Control") != "public, max-age=86400" {
			t.Errorf("%s: %d %q", p, w.Code, w.Header().Get("Cache-Control"))
		}
	}
}

func TestGzip(t *testing.T) {
	h := newHarness(t, nil)
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r.Header.Set("Accept-Encoding", "br, gzip;q=0.8")
	w := h.do(r)
	if w.Header().Get("Content-Encoding") != "gzip" || !strings.Contains(w.Header().Get("Vary"), "Accept-Encoding") {
		t.Fatalf("index not gzipped: %v", w.Header())
	}
	zr, err := gzip.NewReader(w.Body)
	if err != nil {
		t.Fatal(err)
	}
	plain, _ := io.ReadAll(zr)
	if !bytes.Contains(plain, []byte("<h1")) {
		t.Error("gzipped index does not decode to the page")
	}
	r.Header.Set("Accept-Encoding", "gzip;q=0")
	if w := h.do(r); w.Header().Get("Content-Encoding") != "" {
		t.Error("gzip;q=0 must refuse gzip")
	}
	r = httptest.NewRequest(http.MethodHead, "/", nil)
	if w := h.do(r); w.Code != http.StatusOK || w.Body.Len() != 0 {
		t.Errorf("HEAD /: %d, %d bytes", w.Code, w.Body.Len())
	}
}
