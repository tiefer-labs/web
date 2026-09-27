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
	"net/mail"
	"net/smtp"
	"strconv"
	"strings"
	"time"
)

// Mailer delivers a message.
type Mailer interface {
	Send(ctx context.Context, from string, to []string, msg []byte) error
}

// Envelope holds everything needed to build the notification email.
type Envelope struct {
	From      string // CONTACT_FROM
	To        string // CONTACT_TO
	SiteURL   string
	RoleLabel string // human readable role, or empty
	Now       time.Time
}

// BuildMessage renders a submission as a plain text email with Reply-To
// set to the sender. All header values are encoded, and Parse has already
// rejected line breaks in single-line fields.
func BuildMessage(e Envelope, s Submission) []byte {
	var b bytes.Buffer
	host := "localhost"
	if i := strings.LastIndexByte(e.From, '@'); i >= 0 {
		host = e.From[i+1:]
	}
	id := make([]byte, 12)
	_, _ = rand.Read(id)

	from := mail.Address{Name: "Tiefer website", Address: e.From}
	reply := mail.Address{Name: s.Name, Address: s.Email}
	header := [][2]string{
		{"From", from.String()},
		{"To", e.To},
		{"Reply-To", reply.String()},
		{"Subject", mime.QEncoding.Encode("utf-8", "Access request from "+s.Name)},
		{"Date", e.Now.Format(time.RFC1123Z)},
		{"Message-ID", "<" + hex.EncodeToString(id) + "@" + host + ">"},
		{"MIME-Version", "1.0"},
		{"Content-Type", "text/plain; charset=utf-8"},
		{"Content-Transfer-Encoding", "quoted-printable"},
		{"Auto-Submitted", "auto-generated"},
	}
	for _, h := range header {
		fmt.Fprintf(&b, "%s: %s\r\n", h[0], h[1])
	}
	b.WriteString("\r\n")

	orEmpty := func(v string) string {
		if v == "" {
			return "(not given)"
		}
		return v
	}
	var body bytes.Buffer
	fmt.Fprintf(&body, "Name: %s\n", s.Name)
	fmt.Fprintf(&body, "Work email: %s\n", s.Email)
	fmt.Fprintf(&body, "Organisation: %s\n", orEmpty(s.Org))
	fmt.Fprintf(&body, "Role: %s\n\n", orEmpty(e.RoleLabel))
	fmt.Fprintf(&body, "The first question they would ask Tiefer:\n\n%s\n\n", s.Question)
	fmt.Fprintf(&body, "The sender agreed that we may use these details to reply.\n")
	fmt.Fprintf(&body, "Sent from the contact form at %s. Reply to this email to answer.\n", e.SiteURL)

	qp := quotedprintable.NewWriter(&b)
	_, _ = qp.Write(bytes.ReplaceAll(body.Bytes(), []byte("\n"), []byte("\r\n")))
	_ = qp.Close()
	return b.Bytes()
}

// SMTPMailer sends mail through an SMTP server. TLS is required: port 465
// uses implicit TLS, every other port must offer STARTTLS.
type SMTPMailer struct {
	Host       string
	Port       int
	User       string
	Pass       string
	SkipVerify bool   // development only, enforced by config
	LocalName  string // name sent in EHLO
	Timeout    time.Duration
}

// Send delivers msg. It honours ctx cancellation and the Timeout.
func (m *SMTPMailer) Send(ctx context.Context, from string, to []string, msg []byte) error {
	timeout := m.Timeout
	if timeout == 0 {
		timeout = 20 * time.Second
	}
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	addr := net.JoinHostPort(m.Host, strconv.Itoa(m.Port))
	tlsConf := &tls.Config{ServerName: m.Host, MinVersion: tls.VersionTLS12, InsecureSkipVerify: m.SkipVerify}
	dialer := &net.Dialer{}
	var conn net.Conn
	var err error
	if m.Port == 465 {
		conn, err = (&tls.Dialer{NetDialer: dialer, Config: tlsConf}).DialContext(ctx, "tcp", addr)
	} else {
		conn, err = dialer.DialContext(ctx, "tcp", addr)
	}
	if err != nil {
		return fmt.Errorf("smtp dial: %w", err)
	}
	deadline, _ := ctx.Deadline()
	_ = conn.SetDeadline(deadline)
	stop := context.AfterFunc(ctx, func() { _ = conn.Close() })
	defer stop()

	c, err := smtp.NewClient(conn, m.Host)
	if err != nil {
		_ = conn.Close()
		return fmt.Errorf("smtp greeting: %w", err)
	}
	defer c.Close()
	if m.LocalName != "" {
		if err := c.Hello(m.LocalName); err != nil {
			return fmt.Errorf("smtp hello: %w", err)
		}
	}
	if m.Port != 465 {
		if ok, _ := c.Extension("STARTTLS"); !ok {
			return errors.New("smtp: server does not offer STARTTLS, and TLS is required")
		}
		if err := c.StartTLS(tlsConf); err != nil {
			return fmt.Errorf("smtp starttls: %w", err)
		}
	}
	if m.User != "" {
		if err := c.Auth(smtp.PlainAuth("", m.User, m.Pass, m.Host)); err != nil {
			return fmt.Errorf("smtp auth: %w", err)
		}
	}
	if err := c.Mail(from); err != nil {
		return fmt.Errorf("smtp mail from: %w", err)
	}
	for _, rcpt := range to {
		if err := c.Rcpt(rcpt); err != nil {
			return fmt.Errorf("smtp rcpt to: %w", err)
		}
	}
	w, err := c.Data()
	if err != nil {
		return fmt.Errorf("smtp data: %w", err)
	}
	if _, err := w.Write(msg); err != nil {
		return fmt.Errorf("smtp write: %w", err)
	}
	if err := w.Close(); err != nil {
		return fmt.Errorf("smtp data end: %w", err)
	}
	return c.Quit()
}
