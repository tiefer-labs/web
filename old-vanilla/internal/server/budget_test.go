// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package server

import (
	"net/http"
	"net/http/httptest"
	"regexp"
	"runtime"
	"slices"
	"strings"
	"testing"
	"time"
)

// Budgets from docs/PROMPT.md section 9.
const (
	budgetHTMLGzip   = 30 << 10
	budgetCSSGzip    = 12 << 10
	budgetJSRaw      = 3 << 10
	budgetFonts      = 110 << 10
	budgetImages     = 40 << 10
	budgetFirstView  = 200 << 10
	budgetRequests   = 12
	budgetRenderP99  = 5 * time.Millisecond
	renderIterations = 400
)

// firstView lists the resources the index page loads before any
// scrolling: stylesheet, script, preloaded fonts, icons, manifest and
// every image that is not lazy.
var firstView = regexp.MustCompile(`<(?:link|script|img)\b[^>]*?\s(?:href|src)="([^"]+)"[^>]*>`)

func fetchSize(t *testing.T, h *harness, url string) (size int, ctype string) {
	t.Helper()
	r := httptest.NewRequest(http.MethodGet, url, nil)
	r.Header.Set("Accept-Encoding", "gzip")
	w := h.do(r)
	if w.Code != http.StatusOK {
		t.Errorf("%s: status %d", url, w.Code)
	}
	return w.Body.Len(), w.Header().Get("Content-Type")
}

func TestPerformanceBudgets(t *testing.T) {
	h := newHarness(t, map[string]string{"SMTP_HOST": "smtp.example.org", "CONTACT_TO": "a@example.org", "CONTACT_FROM": "b@example.org", "REPO_URL": "https://github.com/tiefer-labs/web"})
	htmlSize, _ := fetchSize(t, h, "/")
	if htmlSize > budgetHTMLGzip {
		t.Errorf("index HTML is %d bytes gzipped, budget %d", htmlSize, budgetHTMLGzip)
	}
	body := h.get("/").Body.String()

	total, requests, images, fonts := htmlSize, 1, 0, 0
	for _, m := range firstView.FindAllStringSubmatch(body, -1) {
		tag, url := m[0], m[1]
		if strings.Contains(tag, `rel="canonical"`) || strings.Contains(tag, `rel="alternate"`) || strings.Contains(tag, `loading="lazy"`) {
			continue
		}
		if !strings.HasPrefix(url, "/") || strings.HasPrefix(url, "//") {
			t.Errorf("first view loads %s from another origin", url)
			continue
		}
		size, ctype := fetchSize(t, h, url)
		requests++
		total += size
		switch {
		case strings.HasPrefix(ctype, "font/"):
			fonts += size
		case strings.HasPrefix(ctype, "image/"):
			images += size
		case strings.HasPrefix(ctype, "text/css"):
			if size > budgetCSSGzip {
				t.Errorf("CSS is %d bytes gzipped, budget %d", size, budgetCSSGzip)
			}
		}
	}
	js := h.srv.assets
	if as, _ := js.ByPath("js/site.js"); len(as.Body) > budgetJSRaw {
		t.Errorf("JavaScript is %d bytes raw, budget %d", len(as.Body), budgetJSRaw)
	}
	var allFonts int
	for _, as := range h.srv.assets.All() {
		if strings.HasSuffix(as.Path, ".woff2") {
			allFonts += len(as.Body)
		}
	}
	checks := []struct {
		name        string
		got, budget int
	}{
		{"fonts (all WOFF2 files)", allFonts, budgetFonts},
		{"first-view images and SVGs", images, budgetImages},
		{"first-view transfer", total, budgetFirstView},
		{"first-view requests", requests, budgetRequests},
	}
	for _, c := range checks {
		if c.got > c.budget {
			t.Errorf("%s: %d, budget %d", c.name, c.got, c.budget)
		}
	}
	t.Logf("index %d B gz, first view %d B in %d requests (fonts %d B, images %d B)", htmlSize, total, requests, fonts, images)
}

// TestRenderTime checks the p99 render time of the index page, including
// gzip. Other test packages may run at the same time, so it warms up and
// takes the best of three rounds. Skipped under -race, which slows
// everything down.
func TestRenderTime(t *testing.T) {
	if raceEnabled || testing.Short() {
		t.Skip("timing budget is checked without -race")
	}
	h := newHarness(t, nil)
	render := func() time.Duration {
		r := httptest.NewRequest(http.MethodGet, "/", nil)
		r.Header.Set("Accept-Encoding", "gzip")
		start := time.Now()
		h.do(r)
		return time.Since(start)
	}
	for range 50 {
		render()
	}
	best := time.Hour
	var p50 time.Duration
	for range 3 {
		runtime.GC()
		times := make([]time.Duration, renderIterations)
		for i := range times {
			times[i] = render()
		}
		slices.Sort(times)
		if p99 := times[len(times)*99/100]; p99 < best {
			best, p50 = p99, times[len(times)/2]
		}
	}
	if best > budgetRenderP99 {
		t.Errorf("index render p99 %v, budget %v", best, budgetRenderP99)
	}
	t.Logf("index render p50 %v, p99 %v", p50, best)
}

func BenchmarkIndex(b *testing.B) {
	h := newHarness(b, nil)
	b.ReportAllocs()
	for b.Loop() {
		r := httptest.NewRequest(http.MethodGet, "/", nil)
		r.Header.Set("Accept-Encoding", "gzip")
		h.do(r)
	}
}
