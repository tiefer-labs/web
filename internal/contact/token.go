// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

// Package contact validates contact form submissions, protects the form
// against forgery and spam without cookies or third parties, and delivers
// messages by SMTP.
package contact

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"sync"
	"time"
)

// Token errors. ErrTooFast marks a form sent faster than a person can
// fill it in.
var (
	ErrTokenInvalid = errors.New("contact: invalid form token")
	ErrTokenExpired = errors.New("contact: form token expired or already used")
	ErrTooFast      = errors.New("contact: form sent too fast")
)

const (
	tokenPayload = 8 + 16 // issue time (unix seconds) + random nonce
	tokenLen     = tokenPayload + sha256.Size
	maxUsed      = 100_000 // bound on remembered nonces
)

// Tokens issues and checks signed form tokens. A token carries the time
// the form was rendered and a random nonce, signed with HMAC-SHA256. It
// replaces a session-bound CSRF token (the site sets no cookies): with the
// Origin and Sec-Fetch-Site check in the server, it proves that a
// submission comes from a form this server rendered, and its time gives
// the minimum fill time.
type Tokens struct {
	secret []byte
	now    func() time.Time
	MinAge time.Duration
	MaxAge time.Duration

	mu   sync.Mutex
	used map[[16]byte]time.Time // nonce to expiry
}

// NewTokens returns an issuer. now may be nil for time.Now.
func NewTokens(secret []byte, now func() time.Time) *Tokens {
	if now == nil {
		now = time.Now
	}
	return &Tokens{secret: secret, now: now, MinAge: 3 * time.Second, MaxAge: 2 * time.Hour, used: map[[16]byte]time.Time{}}
}

func (t *Tokens) sign(payload []byte) []byte {
	m := hmac.New(sha256.New, t.secret)
	m.Write([]byte("tiefer contact form v1\x00"))
	m.Write(payload)
	return m.Sum(nil)
}

// Issue returns a new token.
func (t *Tokens) Issue() string {
	b := make([]byte, tokenLen)
	binary.BigEndian.PutUint64(b[:8], uint64(t.now().Unix())) // #nosec G115 -- Unix time is positive
	_, _ = rand.Read(b[8:tokenPayload])
	copy(b[tokenPayload:], t.sign(b[:tokenPayload]))
	return base64.RawURLEncoding.EncodeToString(b)
}

// Check verifies a token without using it up.
func (t *Tokens) Check(s string) error {
	_, err := t.parse(s)
	return err
}

// Use verifies a token and marks it as used, so it cannot be replayed.
func (t *Tokens) Use(s string) error {
	nonce, err := t.parse(s)
	if err != nil {
		return err
	}
	now := t.now()
	t.mu.Lock()
	defer t.mu.Unlock()
	if _, seen := t.used[nonce]; seen {
		return ErrTokenExpired
	}
	if len(t.used) >= maxUsed {
		for n, exp := range t.used {
			if now.After(exp) {
				delete(t.used, n)
			}
		}
		if len(t.used) >= maxUsed {
			return ErrTokenExpired // fail closed under a flood
		}
	}
	t.used[nonce] = now.Add(t.MaxAge + time.Minute)
	return nil
}

func (t *Tokens) parse(s string) ([16]byte, error) {
	var nonce [16]byte
	if len(s) != base64.RawURLEncoding.EncodedLen(tokenLen) {
		return nonce, ErrTokenInvalid
	}
	b, err := base64.RawURLEncoding.DecodeString(s)
	if err != nil || len(b) != tokenLen {
		return nonce, ErrTokenInvalid
	}
	if !hmac.Equal(b[tokenPayload:], t.sign(b[:tokenPayload])) {
		return nonce, ErrTokenInvalid
	}
	issued := time.Unix(int64(binary.BigEndian.Uint64(b[:8])), 0) // #nosec G115 -- signed by us
	age := t.now().Sub(issued)
	switch {
	case age < -time.Minute || age > t.MaxAge:
		return nonce, ErrTokenExpired
	case age < t.MinAge:
		return nonce, ErrTooFast
	}
	copy(nonce[:], b[8:tokenPayload])
	return nonce, nil
}
