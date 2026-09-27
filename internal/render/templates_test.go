// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package render

import "testing"

func TestSlug(t *testing.T) {
	cases := map[string]string{
		"Who is responsible":                  "who-is-responsible",
		"EU representative (Article 27 GDPR)": "eu-representative-article-27-gdpr",
		"  Transfers to other countries  ":    "transfers-to-other-countries",
		"Tax identification number (VÖEN)":    "tax-identification-number-v-en",
	}
	for in, want := range cases {
		if got := Slug(in); got != want {
			t.Errorf("Slug(%q) = %q, want %q", in, got, want)
		}
	}
}
