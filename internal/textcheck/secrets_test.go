// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package textcheck

import (
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// secretPatterns match credentials that must never be committed. The
// patterns are written so that this file does not match itself.
var secretPatterns = map[string]*regexp.Regexp{
	"private key":             regexp.MustCompile(`-----BEGIN [A-Z ]*PRIVATE KEY-----`),
	"AWS access key":          regexp.MustCompile(`\bAKIA[0-9A-Z]{16}\b`),
	"GitHub token":            regexp.MustCompile(`\bgh[pousr]_[A-Za-z0-9]{36,}\b|\bgithub_pat_[A-Za-z0-9_]{40,}\b`),
	"Azure storage key":       regexp.MustCompile(`AccountKey=[A-Za-z0-9+/]{40,}={0,2}`),
	"Azure SAS signature":     regexp.MustCompile(`[?&]sig=[A-Za-z0-9%]{30,}`),
	"Slack token":             regexp.MustCompile(`\bxox[abprs]-[A-Za-z0-9-]{10,}\b`),
	"JSON Web Token":          regexp.MustCompile(`\beyJ[A-Za-z0-9_-]{10,}\.eyJ[A-Za-z0-9_-]{10,}\.[A-Za-z0-9_-]{10,}`),
	"secret in an assignment": regexp.MustCompile(`(?im)^\s*(?:SMTP_PASS|CSRF_SECRET|FRONT_DOOR_ID|[A-Z_]*API_KEY)\s*[=:]\s*["']?[A-Za-z0-9+/_\-]{12,}`),
}

// TestNoSecrets scans every committed text file for credentials. Local
// .env files are ignored by git and skipped here.
func TestNoSecrets(t *testing.T) {
	root := filepath.Join("..", "..")
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, _ := filepath.Rel(root, path)
		rel = filepath.ToSlash(rel)
		if d.IsDir() {
			switch rel {
			case ".git", "bin", "dist", "node_modules":
				return filepath.SkipDir
			}
			return nil
		}
		if binary[filepath.Ext(path)] || rel == "tools/go.sum" || (strings.HasPrefix(d.Name(), ".env") && d.Name() != ".env.example") {
			return nil
		}
		b, err := os.ReadFile(path) // #nosec G304 -- walking this repository in a test
		if err != nil {
			return err
		}
		for what, re := range secretPatterns {
			if loc := re.FindIndex(b); loc != nil {
				t.Errorf("%s:%d: looks like a %s", rel, strings.Count(string(b[:loc[0]]), "\n")+1, what)
			}
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}

// TestSecretPatterns makes sure the patterns catch what they are for.
func TestSecretPatterns(t *testing.T) {
	samples := map[string]string{
		"private key":             "-----BEGIN " + "RSA PRIVATE KEY-----",
		"AWS access key":          "AKIA" + "ABCDEFGHIJKLMNOP",
		"GitHub token":            "ghp_" + strings.Repeat("a1B2", 9),
		"Azure storage key":       "AccountKey=" + strings.Repeat("Ab1+", 11) + "==",
		"Azure SAS signature":     "https://x.blob.core.windows.net/c?sv=1&sig=" + strings.Repeat("aB3", 12),
		"Slack token":             "xoxb-" + "1234567890-abcdef",
		"JSON Web Token":          "eyJ" + "hbGciOiJIUzI1NiJ9.eyJ" + "zdWIiOiIxMjM0NTY3ODkwIn0.abcdefghijklmnop",
		"secret in an assignment": "SMTP_PASS=" + "hunter2hunter2hunter2",
	}
	for what, s := range samples {
		if !secretPatterns[what].MatchString(s) {
			t.Errorf("%s pattern misses %q", what, s)
		}
	}
	for _, ok := range []string{"SMTP_PASS=", "CSRF_SECRET=", "CSRF_SECRET: required in production"} {
		if secretPatterns["secret in an assignment"].MatchString(ok) {
			t.Errorf("false positive on %q", ok)
		}
	}
}
