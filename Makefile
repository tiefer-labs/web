# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.

# Development tools are pinned in tools/go.mod and run with "go tool".
TOOL    := go tool -modfile=tools/go.mod
BIN     := bin/tiefer-web
IMAGE   := tiefer-web:local
LDFLAGS := -s -w -buildid=
FUZZTIME ?= 20s

.PHONY: run build test check fmt vet lint sec vuln fuzz docker placeholders clean

run:
	go run ./cmd/tiefer-web

build:
	CGO_ENABLED=0 go build -trimpath -ldflags '$(LDFLAGS)' -o $(BIN) ./cmd/tiefer-web

test:
	go test -race -count=1 ./...

# check runs everything a commit must pass. govulncheck needs the network
# (vuln.go.dev), so it runs separately as "make vuln" and in CI.
check: fmt vet lint sec test

fmt:
	@out="$$(gofmt -l cmd internal web)"; if [ -n "$$out" ]; then echo "gofmt needed:"; echo "$$out"; exit 1; fi

vet:
	go vet ./...

lint:
	$(TOOL) staticcheck ./...

sec:
	$(TOOL) gosec -quiet ./...

vuln:
	$(TOOL) govulncheck ./...

# Each fuzz target runs for FUZZTIME.
fuzz:
	@for pkg in $$(go list ./...); do \
	  for t in $$(go test -list '^Fuzz' $$pkg | grep '^Fuzz'); do \
	    echo "fuzz $$pkg $$t"; go test -run '^$$' -fuzz "^$$t$$" -fuzztime $(FUZZTIME) $$pkg || exit 1; \
	  done; \
	done

docker:
	docker build -t $(IMAGE) .

# Lists placeholders and configuration defaults still in use (warnings).
placeholders:
	@go test -count=1 -v -run 'TestPlaceholderReport' ./internal/server | grep WARNING || true

clean:
	rm -rf bin
