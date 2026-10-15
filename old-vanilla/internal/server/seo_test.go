// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package server

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"encoding/xml"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"testing"
)

var jsonLDBlock = regexp.MustCompile(`<script type="application/ld\+json">(.*?)</script>`)

// TestJSONLD checks the Organization block and that the CSP allows exactly
// its bytes by hash.
func TestJSONLD(t *testing.T) {
	h := newHarness(t, prod(nil))
	w := h.do(hostReq(h, "/"))
	m := jsonLDBlock.FindStringSubmatch(w.Body.String())
	if m == nil {
		t.Fatal("no JSON-LD block on the index page")
	}
	var org map[string]any
	if err := json.Unmarshal([]byte(m[1]), &org); err != nil {
		t.Fatal(err)
	}
	if org["@type"] != "Organization" || org["name"] != "Tiefer" || org["url"] != "https://tiefer.space/" {
		t.Errorf("organization: %v", org)
	}
	if same, _ := org["sameAs"].([]any); len(same) != 1 || same[0] != "https://www.linkedin.com/company/tiefer/" {
		t.Errorf("sameAs: %v", org["sameAs"])
	}
	if logo, _ := org["logo"].(string); !strings.HasPrefix(logo, "https://tiefer.space/static/icon-512.") {
		t.Errorf("logo: %v", org["logo"])
	}
	if len(org) != 6 {
		t.Errorf("unexpected fields: %v", org)
	}
	sum := sha256.Sum256([]byte(m[1]))
	hash := "'sha256-" + base64.StdEncoding.EncodeToString(sum[:]) + "'"
	csp := w.Header().Get("Content-Security-Policy")
	if !strings.Contains(csp, "script-src 'self' "+hash) || strings.Contains(csp, "unsafe-inline") {
		t.Errorf("CSP does not allow the JSON-LD block by hash: %s", csp)
	}
	for _, p := range legalPaths {
		if strings.Contains(h.do(hostReq(h, p)).Body.String(), "application/ld+json") {
			t.Errorf("%s: JSON-LD belongs on the index page only", p)
		}
	}
}

func hostReq(h *harness, path string) *http.Request {
	r, _ := http.NewRequest(http.MethodGet, path, nil)
	r.Host = h.cfg.SiteURL.Host
	r.RequestURI = path
	return r
}

func TestMetadata(t *testing.T) {
	h := newHarness(t, prod(nil))
	body := h.do(hostReq(h, "/")).Body.String()
	for _, want := range []string{
		`<html lang="en">`,
		`<title>Tiefer | Edge AI in orbit</title>`,
		`<meta name="description" content="Tiefer builds AI software that runs on Earth observation satellites: it filters useless pixels in orbit, detects events on board and sends kilobyte-sized alerts instead of gigabytes of raw imagery. Built in Baku.">`,
		`<link rel="canonical" href="https://tiefer.space/">`,
		`<meta property="og:url" content="https://tiefer.space/">`,
		`<meta property="og:image" content="https://tiefer.space/static/og-image.`,
		`<meta name="twitter:card" content="summary_large_image">`,
		`<meta name="theme-color" content="#0C003D">`,
		`<link rel="manifest" href="/site.webmanifest">`,
	} {
		if !strings.Contains(body, want) {
			t.Errorf("index lacks %s", want)
		}
	}
	privacy := h.do(hostReq(h, "/privacy")).Body.String()
	if !strings.Contains(privacy, `<link rel="canonical" href="https://tiefer.space/privacy">`) || !strings.Contains(privacy, "<title>Privacy notice | Tiefer</title>") {
		t.Error("privacy page metadata")
	}
}

var headingTag = regexp.MustCompile(`<h([1-6])\b`)

// TestHeadingOrder checks one h1 per page and no skipped heading levels.
func TestHeadingOrder(t *testing.T) {
	h := newHarness(t, smtpEnv)
	for _, p := range append([]string{"/", "/missing"}, legalPaths...) {
		levels := headingTag.FindAllStringSubmatch(h.get(p).Body.String(), -1)
		prev := 0
		ones := 0
		for _, l := range levels {
			n, _ := strconv.Atoi(l[1])
			if n == 1 {
				ones++
			}
			if n > prev+1 {
				t.Errorf("%s: h%d after h%d", p, n, prev)
			}
			prev = n
		}
		if ones != 1 {
			t.Errorf("%s: %d h1 elements", p, ones)
		}
	}
}

func TestSitemapAndRobots(t *testing.T) {
	h := newHarness(t, prod(nil))
	w := h.do(hostReq(h, "/sitemap.xml"))
	var set struct {
		URLs []struct {
			Loc string `xml:"loc"`
		} `xml:"url"`
	}
	if err := xml.Unmarshal(w.Body.Bytes(), &set); err != nil {
		t.Fatal(err)
	}
	var locs []string
	for _, u := range set.URLs {
		locs = append(locs, u.Loc)
	}
	if strings.Join(locs, " ") != "https://tiefer.space/ https://tiefer.space/legal https://tiefer.space/privacy https://tiefer.space/acceptable-use" {
		t.Errorf("sitemap: %v", locs)
	}
	for _, loc := range locs {
		path := strings.TrimPrefix(loc, "https://tiefer.space")
		if w := h.do(hostReq(h, path)); w.Code != http.StatusOK {
			t.Errorf("sitemap lists %s, which answers %d", loc, w.Code)
		}
	}
	if robots := h.do(hostReq(h, "/robots.txt")).Body.String(); !strings.Contains(robots, "Sitemap: https://tiefer.space/sitemap.xml") || strings.Contains(robots, "Disallow: /\n") {
		t.Errorf("production robots.txt: %q", robots)
	}
	if robots := newHarness(t, nil).get("/robots.txt").Body.String(); !strings.Contains(robots, "Disallow: /\n") {
		t.Errorf("development robots.txt must keep crawlers out: %q", robots)
	}
}

func TestManifest(t *testing.T) {
	h := newHarness(t, nil)
	w := h.get("/site.webmanifest")
	var m struct {
		Name       string `json:"name"`
		ThemeColor string `json:"theme_color"`
		Icons      []struct {
			Src string `json:"src"`
		} `json:"icons"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &m); err != nil {
		t.Fatal(err)
	}
	if m.Name != "Tiefer" || m.ThemeColor != "#0C003D" || len(m.Icons) != 3 || w.Header().Get("Content-Type") != "application/manifest+json" {
		t.Errorf("manifest: %+v", m)
	}
	for _, ic := range m.Icons {
		if w := h.get(ic.Src); w.Code != http.StatusOK {
			t.Errorf("manifest icon %s: %d", ic.Src, w.Code)
		}
	}
}
