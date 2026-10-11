# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.

# Build stage: the Go toolchain, pinned by digest (golang:1.27.1-bookworm).
FROM golang:1.27.2-bookworm@sha256:5cf287a799e6b94384bad13d16b14904c531f51ba65792237e122ce42b392f61 AS build
WORKDIR /src
ENV CGO_ENABLED=0 GOFLAGS=-trimpath GOTOOLCHAIN=local
COPY go.mod ./
COPY cmd ./cmd
COPY internal ./internal
COPY web ./web
RUN go build -ldflags '-s -w -buildid=' -o /out/tiefer-web ./cmd/tiefer-web

# Runtime stage: no shell, no package manager, a non-root user
# (gcr.io/distroless/static-debian12:nonroot, pinned by digest).
FROM gcr.io/distroless/static-debian12:nonroot@sha256:afa5c872c891853ca7fcf1f12c3edb23f7eeef36189728842dd51042ff57f7ab
COPY --from=build /out/tiefer-web /tiefer-web
USER 65532:65532
ENV ENV=production PORT=8080
EXPOSE 8080
HEALTHCHECK --interval=30s --timeout=5s --start-period=5s --retries=3 CMD ["/tiefer-web", "healthcheck"]
ENTRYPOINT ["/tiefer-web"]
