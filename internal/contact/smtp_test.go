// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package contact

import (
	"bufio"
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"math/big"
	"net"
	"strconv"
	"strings"
	"testing"
	"time"
)

// testCert returns a self-signed certificate for 127.0.0.1 and a pool
// that trusts it.
func testCert(t *testing.T) (tls.Certificate, *x509.CertPool) {
	t.Helper()
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	tpl := &x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject:      pkix.Name{CommonName: "smtp.test"},
		DNSNames:     []string{"smtp.test"},
		NotBefore:    time.Now().Add(-time.Hour),
		NotAfter:     time.Now().Add(time.Hour),
		KeyUsage:     x509.KeyUsageDigitalSignature,
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
	}
	der, err := x509.CreateCertificate(rand.Reader, tpl, tpl, &key.PublicKey, key)
	if err != nil {
		t.Fatal(err)
	}
	cert, _ := x509.ParseCertificate(der)
	pool := x509.NewCertPool()
	pool.AddCert(cert)
	return tls.Certificate{Certificate: [][]byte{der}, PrivateKey: key}, pool
}

// fakeSMTP is a minimal SMTP server with STARTTLS. It records the
// commands and the message it receives.
type fakeSMTP struct {
	ln       net.Listener
	cert     tls.Certificate
	starttls bool
	got      chan string
}

func startSMTP(t *testing.T, starttls bool) (*fakeSMTP, *x509.CertPool) {
	cert, pool := testCert(t)
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	f := &fakeSMTP{ln: ln, cert: cert, starttls: starttls, got: make(chan string, 1)}
	go f.serve()
	t.Cleanup(func() { ln.Close() })
	return f, pool
}

func (f *fakeSMTP) serve() {
	conn, err := f.ln.Accept()
	if err != nil {
		return
	}
	defer conn.Close()
	var log strings.Builder
	rw := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))
	say := func(s string) { rw.WriteString(s + "\r\n"); rw.Flush() }
	say("220 smtp.test ready")
	tlsOn := false
	for {
		line, err := rw.ReadString('\n')
		if err != nil {
			f.got <- log.String()
			return
		}
		cmd := strings.TrimSpace(line)
		log.WriteString("C: " + cmd + "\n")
		switch up := strings.ToUpper(cmd); {
		case strings.HasPrefix(up, "EHLO"):
			rw.WriteString("250-smtp.test\r\n")
			if f.starttls && !tlsOn {
				rw.WriteString("250-STARTTLS\r\n")
			}
			say("250 AUTH PLAIN")
		case up == "STARTTLS":
			say("220 go ahead")
			tc := tls.Server(conn, &tls.Config{Certificates: []tls.Certificate{f.cert}, MinVersion: tls.VersionTLS12})
			if err := tc.Handshake(); err != nil {
				f.got <- log.String() + "handshake: " + err.Error()
				return
			}
			tlsOn = true
			log.WriteString("TLS\n")
			rw = bufio.NewReadWriter(bufio.NewReader(tc), bufio.NewWriter(tc))
		case strings.HasPrefix(up, "AUTH PLAIN"):
			say("235 ok")
		case strings.HasPrefix(up, "MAIL FROM"), strings.HasPrefix(up, "RCPT TO"):
			say("250 ok")
		case up == "DATA":
			say("354 go on")
			for {
				l, err := rw.ReadString('\n')
				if err != nil || l == ".\r\n" {
					break
				}
				log.WriteString("D: " + l)
			}
			say("250 queued")
		case up == "QUIT":
			say("221 bye")
			f.got <- log.String()
			return
		default:
			say("500 unknown")
		}
	}
}

func (f *fakeSMTP) port() int {
	_, p, _ := net.SplitHostPort(f.ln.Addr().String())
	n, _ := strconv.Atoi(p)
	return n
}

func testMessage() Message {
	return Message{From: "web@tiefer.space", To: "hello@tiefer.space", ReplyTo: "ada@example.org",
		Subject: "Website contact: Ada", Body: "Hello", Date: time.Now(), Domain: "tiefer.space"}
}

func TestSMTPStartTLS(t *testing.T) {
	f, pool := startSMTP(t, true)
	m := &SMTPMailer{Host: "127.0.0.1", Port: f.port(), User: "web", Pass: "secret", LocalName: "tiefer.space",
		TLSConfig: &tls.Config{ServerName: "smtp.test", RootCAs: pool, MinVersion: tls.VersionTLS12}}
	if err := m.Send(context.Background(), testMessage()); err != nil {
		t.Fatal(err)
	}
	log := <-f.got
	tlsAt, authAt := strings.Index(log, "TLS\n"), strings.Index(log, "AUTH PLAIN")
	if tlsAt < 0 || authAt < tlsAt {
		t.Errorf("credentials must only be sent after STARTTLS:\n%s", log)
	}
	for _, want := range []string{"MAIL FROM:<web@tiefer.space>", "RCPT TO:<hello@tiefer.space>", "D: Reply-To: ada@example.org"} {
		if !strings.Contains(log, want) {
			t.Errorf("session lacks %q:\n%s", want, log)
		}
	}
}

func TestSMTPRefusesPlaintext(t *testing.T) {
	f, pool := startSMTP(t, false)
	m := &SMTPMailer{Host: "127.0.0.1", Port: f.port(), User: "web", Pass: "secret",
		TLSConfig: &tls.Config{ServerName: "smtp.test", RootCAs: pool, MinVersion: tls.VersionTLS12}}
	err := m.Send(context.Background(), testMessage())
	if err == nil || !strings.Contains(err.Error(), "STARTTLS") {
		t.Fatalf("server without STARTTLS: %v", err)
	}
	f.ln.Close()
	if log := <-f.got; strings.Contains(log, "AUTH") || strings.Contains(log, "MAIL FROM") {
		t.Errorf("nothing may be sent without TLS:\n%s", log)
	}
}

func TestSMTPRejectsUntrustedCertificate(t *testing.T) {
	f, _ := startSMTP(t, true)
	m := &SMTPMailer{Host: "127.0.0.1", Port: f.port(), TLSConfig: &tls.Config{ServerName: "smtp.test", MinVersion: tls.VersionTLS12}}
	if err := m.Send(context.Background(), testMessage()); err == nil {
		t.Fatal("an untrusted certificate must fail")
	}
}
