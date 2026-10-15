// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package contact

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/tls"
	"encoding/hex"
	"errors"
	"fmt"
	"mime"
	"mime/quotedprintable"
	"net"
	"net/smtp"
	"strconv"
	"strings"
	"time"
)

// Message is one mail to send. Every header value has been validated
// before it gets here; Build checks again.
type Message struct {
	From    string
	To      string
	ReplyTo string
	Subject string
	Body    string
	Date    time.Time
	Domain  string // domain for the Message-ID
}

// Mailer delivers messages. Tests use a fake.
type Mailer interface {
	Send(ctx context.Context, m Message) error
}

// ErrHeader reports a header value that could inject further headers.
var ErrHeader = errors.New("contact: header value contains a line break")

// Build renders a message as RFC 5322 bytes with a quoted-printable UTF-8
// body. It refuses header values with line breaks.
func Build(m Message) ([]byte, error) {
	for _, v := range []string{m.From, m.To, m.ReplyTo, m.Subject, m.Domain} {
		if strings.ContainsAny(v, "\r\n\x00") {
			return nil, ErrHeader
		}
	}
	id := make([]byte, 16)
	_, _ = rand.Read(id)
	var b bytes.Buffer
	header := func(k, v string) { fmt.Fprintf(&b, "%s: %s\r\n", k, v) }
	header("From", m.From)
	header("To", m.To)
	if m.ReplyTo != "" {
		header("Reply-To", m.ReplyTo)
	}
	header("Subject", mime.QEncoding.Encode("utf-8", m.Subject))
	header("Date", m.Date.UTC().Format(time.RFC1123Z))
	header("Message-ID", "<"+hex.EncodeToString(id)+"@"+m.Domain+">")
	header("MIME-Version", "1.0")
	header("Content-Type", "text/plain; charset=utf-8")
	header("Content-Transfer-Encoding", "quoted-printable")
	header("Auto-Submitted", "auto-generated")
	b.WriteString("\r\n")
	qp := quotedprintable.NewWriter(&b)
	_, _ = qp.Write([]byte(strings.ReplaceAll(m.Body, "\n", "\r\n")))
	_ = qp.Close()
	return b.Bytes(), nil
}

// SMTPMailer sends through an SMTP server with TLS: implicit TLS on port
// 465, STARTTLS on any other port. It never sends without TLS.
type SMTPMailer struct {
	Host      string
	Port      int
	User      string
	Pass      string
	LocalName string // name used in EHLO
	Timeout   time.Duration
	// SkipVerify accepts any certificate (development only; the
	// configuration refuses it in production).
	SkipVerify bool
	// TLSConfig is for tests only; nil means verified TLS 1.2 or later.
	TLSConfig *tls.Config
}

// Send delivers m. The whole exchange is bounded by the context and by
// Timeout.
func (s *SMTPMailer) Send(ctx context.Context, m Message) error {
	raw, err := Build(m)
	if err != nil {
		return err
	}
	timeout := s.Timeout
	if timeout == 0 {
		timeout = 20 * time.Second
	}
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	cfg := s.TLSConfig
	if cfg == nil {
		cfg = &tls.Config{ServerName: s.Host, MinVersion: tls.VersionTLS12, InsecureSkipVerify: s.SkipVerify} // #nosec G402 -- SkipVerify is development only, refused in production by config
	}
	addr := net.JoinHostPort(s.Host, strconv.Itoa(s.Port))
	var d net.Dialer
	conn, err := d.DialContext(ctx, "tcp", addr)
	if err != nil {
		return err
	}
	if dl, ok := ctx.Deadline(); ok {
		_ = conn.SetDeadline(dl)
	}
	if s.Port == 465 {
		tc := tls.Client(conn, cfg)
		if err := tc.HandshakeContext(ctx); err != nil {
			_ = conn.Close()
			return err
		}
		conn = tc
	}
	c, err := smtp.NewClient(conn, s.Host)
	if err != nil {
		_ = conn.Close()
		return err
	}
	defer c.Close()
	if s.LocalName != "" {
		if err := c.Hello(s.LocalName); err != nil {
			return err
		}
	}
	if s.Port != 465 {
		if ok, _ := c.Extension("STARTTLS"); !ok {
			return errors.New("contact: SMTP server does not offer STARTTLS")
		}
		if err := c.StartTLS(cfg); err != nil {
			return err
		}
	}
	if s.User != "" {
		if err := c.Auth(smtp.PlainAuth("", s.User, s.Pass, s.Host)); err != nil {
			return err
		}
	}
	if err := c.Mail(m.From); err != nil {
		return err
	}
	if err := c.Rcpt(m.To); err != nil {
		return err
	}
	w, err := c.Data()
	if err != nil {
		return err
	}
	if _, err := w.Write(raw); err != nil {
		return err
	}
	if err := w.Close(); err != nil {
		return err
	}
	return c.Quit()
}
