// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

// Command tiefer-web serves the Tiefer website. Everything it needs is
// embedded in the binary; configuration comes from environment variables.
package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/tiefer-labs/web/internal/config"
	"github.com/tiefer-labs/web/internal/server"
)

func main() {
	// "tiefer-web healthcheck" probes /healthz, for container health
	// checks in an image without a shell or curl.
	if len(os.Args) > 1 && os.Args[1] == "healthcheck" {
		os.Exit(healthcheck())
	}
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "tiefer-web:", err)
		os.Exit(1)
	}
}

func run() error {
	cfg, err := config.Load(os.LookupEnv)
	if err != nil {
		return err
	}
	var handler slog.Handler = slog.NewTextHandler(os.Stderr, nil)
	if cfg.Production() {
		handler = slog.NewJSONHandler(os.Stderr, nil)
	}
	log := slog.New(handler)
	slog.SetDefault(log)
	if d := cfg.DefaultsInUse(); len(d) > 0 {
		log.Warn("placeholder or development defaults in use", "vars", strings.Join(d, ","))
	}

	srv, err := server.New(server.Options{Config: cfg, Logger: log})
	if err != nil {
		return err
	}
	hs := srv.HTTPServer(net.JoinHostPort("", strconv.Itoa(cfg.Port)))

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	errc := make(chan error, 1)
	go func() {
		log.Info("listening", "config", cfg)
		errc <- hs.ListenAndServe()
	}()
	select {
	case err := <-errc:
		return err
	case <-ctx.Done():
	}
	log.Info("shutting down")
	sctx, cancel := context.WithTimeout(context.Background(), server.ShutdownTimeout)
	defer cancel()
	if err := hs.Shutdown(sctx); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}

func healthcheck() int {
	// The probe only ever talks to the loopback interface; the port is the
	// only input and must be a plain number.
	port := 8080
	if v := os.Getenv("PORT"); v != "" {
		p, err := strconv.Atoi(v)
		if err != nil || p < 1 || p > 65535 {
			fmt.Fprintln(os.Stderr, "healthcheck: invalid PORT")
			return 1
		}
		port = p
	}
	// The URL is constant; the dialer alone decides the port.
	addr := net.JoinHostPort("127.0.0.1", strconv.Itoa(port))
	client := &http.Client{
		Timeout: 3 * time.Second,
		Transport: &http.Transport{
			Proxy: nil,
			DialContext: func(ctx context.Context, network, _ string) (net.Conn, error) {
				var d net.Dialer
				return d.DialContext(ctx, network, addr)
			},
		},
	}
	res, err := client.Get("http://localhost/healthz")
	if err != nil {
		fmt.Fprintln(os.Stderr, "healthcheck:", err)
		return 1
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		fmt.Fprintln(os.Stderr, "healthcheck: status", res.StatusCode)
		return 1
	}
	return 0
}
