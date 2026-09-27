# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.

BINARY      := bin/tiefer-web
IMAGE       ?= tiefer-web
STATICCHECK := $(shell command -v staticcheck 2>/dev/null || echo $(shell go env GOPATH)/bin/staticcheck)

.PHONY: run build test lint check fmt docker clean

## run: start the server on PORT (default 8080), loading .env if present
run:
	@set -a; [ -f .env ] && . ./.env; set +a; go run ./cmd/tiefer-web

## build: build one self-contained binary into bin/
build:
	CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o $(BINARY) ./cmd/tiefer-web

## test: run every test, including the text check
test:
	go test ./...

## lint: gofmt check, go vet, and staticcheck if it is installed
lint:
	@test -z "$$(gofmt -l .)" || { echo "gofmt needed:"; gofmt -l .; exit 1; }
	go vet ./...
	@if [ -x "$(STATICCHECK)" ]; then echo "$(STATICCHECK) ./..."; "$(STATICCHECK)" ./...; \
	else echo "staticcheck not installed, skipped (go install honnef.co/go/tools/cmd/staticcheck@latest)"; fi

## check: the text check (dashes, emoji, banned words), with placeholder warnings
check:
	go test -count=1 -v -run 'TestTextCheck' ./...

## fmt: format the Go code
fmt:
	gofmt -w .

## docker: build the container image
docker:
	docker build -t $(IMAGE) .

clean:
	rm -rf bin
