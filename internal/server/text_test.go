// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package server

import (
	"net/http"
	"os"
	"sort"
	"strings"
	"testing"

	"github.com/tiefer-labs/web/internal/config"
	"github.com/tiefer-labs/web/internal/textcheck"
)

// renderedPages returns the body of every GET route that answers with
// HTML, plus the 404 page and any extra state the site can render.
func renderedPages(t *testing.T, h *harness) map[string]string {
	t.Helper()
	pages := map[string]string{}
	for _, rt := range h.srv.routes {
		if rt.method != http.MethodGet {
			continue
		}
		w := h.get(rt.sample)
		if strings.HasPrefix(w.Header().Get("Content-Type"), "text/html") {
			pages[rt.sample] = w.Body.String()
		}
	}
	for _, extra := range textStates {
		w := h.get(extra)
		if strings.HasPrefix(w.Header().Get("Content-Type"), "text/html") {
			pages[extra] = w.Body.String()
		}
	}
	return pages
}

// textStates are extra URLs whose pages differ from the plain routes.
var textStates = []string{"/does-not-exist", "/?sent=1"}

// TestRenderedText checks every rendered page, with and without the
// contact form, for forbidden characters anywhere in the HTML and for
// banned words in the visible text. Findings name the page and section.
func TestRenderedText(t *testing.T) {
	envs := []map[string]string{nil, {
		"SMTP_HOST": "smtp.example.org", "CONTACT_TO": "team@example.org", "CONTACT_FROM": "web@example.org",
		"REPO_URL": "https://github.com/tiefer-labs/web",
	}}
	for _, env := range envs {
		h := newHarness(t, env)
		for page, body := range renderedPages(t, h) {
			for _, f := range textcheck.Characters(body) {
				t.Errorf("page %s: line %s", page, f)
			}
			for _, seg := range textcheck.VisibleText(body) {
				for _, f := range textcheck.BannedWords(seg.Text) {
					t.Errorf("page %s, section %s: %s", page, seg.Section, f.What)
				}
			}
		}
	}
}

// TestPlaceholderReport lists, as warnings that do not fail, every
// placeholder left on a rendered page and every configuration value that
// still has its default. Run it with "make placeholders".
func TestPlaceholderReport(t *testing.T) {
	h := newHarness(t, nil)
	found := map[string][]string{}
	for page, body := range renderedPages(t, h) {
		for _, seg := range textcheck.VisibleText(body) {
			for _, p := range textcheck.Placeholders(seg.Text) {
				found[p] = appendUnique(found[p], page)
			}
		}
	}
	keys := make([]string, 0, len(found))
	for k := range found {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		t.Logf("WARNING: placeholder %s on %s", k, strings.Join(found[k], ", "))
	}
	cfg, err := config.Load(os.LookupEnv)
	if err != nil {
		t.Logf("WARNING: the configuration of this environment is invalid: %v", err)
		return
	}
	for _, name := range cfg.DefaultsInUse() {
		t.Logf("WARNING: %s still has its default", name)
	}
}

func appendUnique(list []string, s string) []string {
	for _, v := range list {
		if v == s {
			return list
		}
	}
	return append(list, s)
}
