// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

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

// Token errors.
var (
	ErrTokenInvalid = errors.New("contact: invalid form token")
	ErrTokenExpired = errors.New("contact: form token expired")
	ErrTooFast      = errors.New("contact: form submitted too fast")
)

// Tokens issues and verifies signed form tokens. A token carries the time
// the form was rendered and a random nonce, signed with HMAC-SHA256. It
// replaces a session-bound CSRF token, because the site sets no cookies:
// together with the Origin and Sec-Fetch-Site checks in the server, it
// proves that a submission comes from a form this server rendered, and
// the embedded time gives the minimum fill time check.
type Tokens struct {
	secret []byte
	now    func() time.Time
	MinAge time.Duration // submissions faster than this are rejected
	MaxAge time.Duration // older tokens are expired

	mu   sync.Mutex
	used map[string]time.Time // nonce to expiry, for single use
}

// NewTokens returns a token issuer. now may be nil to use time.Now.
func NewTokens(secret []byte, now func() time.Time) *Tokens {
	if now == nil {
		now = time.Now
	}
	return &Tokens{
		secret: secret,
		now:    now,
		MinAge: 3 * time.Second,
		MaxAge: 12 * time.Hour,
		used:   map[string]time.Time{},
	}
}

const (
	nonceLen   = 12
	payloadLen = 8 + nonceLen
	macLen     = 16 // truncated HMAC-SHA256
)

// Issue returns a new token for a freshly rendered form.
func (t *Tokens) Issue() string {
	b := make([]byte, payloadLen, payloadLen+macLen)
	binary.BigEndian.PutUint64(b, uint64(t.now().Unix()))
	_, _ = rand.Read(b[8:])
	b = append(b, t.mac(b)...)
	return base64.RawURLEncoding.EncodeToString(b)
}

// Check verifies a token's signature and age and returns its nonce.
// The nonce is not marked as used; call Consume after a successful send.
func (t *Tokens) Check(token string) (nonce string, err error) {
	b, err := base64.RawURLEncoding.DecodeString(token)
	if err != nil || len(b) != payloadLen+macLen {
		return "", ErrTokenInvalid
	}
	if !hmac.Equal(b[payloadLen:], t.mac(b[:payloadLen])) {
		return "", ErrTokenInvalid
	}
	issued := time.Unix(int64(binary.BigEndian.Uint64(b)), 0)
	age := t.now().Sub(issued)
	switch {
	case age < -time.Minute:
		return "", ErrTokenInvalid
	case age > t.MaxAge:
		return "", ErrTokenExpired
	case age < t.MinAge:
		return "", ErrTooFast
	}
	return string(b[8:payloadLen]), nil
}

// Consume marks a nonce as used. It returns false if it was used before,
// which means the same form was submitted twice.
func (t *Tokens) Consume(nonce string) bool {
	t.mu.Lock()
	defer t.mu.Unlock()
	now := t.now()
	for n, exp := range t.used {
		if now.After(exp) {
			delete(t.used, n)
		}
	}
	if _, ok := t.used[nonce]; ok {
		return false
	}
	t.used[nonce] = now.Add(t.MaxAge + time.Minute)
	return true
}

func (t *Tokens) mac(payload []byte) []byte {
	m := hmac.New(sha256.New, t.secret)
	m.Write([]byte("tiefer-contact-v1\x00"))
	m.Write(payload)
	return m.Sum(nil)[:macLen]
}
