# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.

# Base images are pinned by digest, so a rebuild uses exactly the same
# bytes. Dependabot proposes updates (.github/dependabot.yml).

# Build stage: the official Go image (golang:1.27).
FROM golang:1.27@sha256:3680233e3204827fbdc66088528ae6d4b3d034f51d03a99d454f6de034888244 AS build
WORKDIR /src
ENV CGO_ENABLED=0 GOFLAGS=-mod=readonly
COPY go.mod ./
RUN go mod download
COPY . .
RUN go build -trimpath -ldflags="-s -w -buildid=" -o /out/tiefer-web ./cmd/tiefer-web

# Run stage: distroless static image (gcr.io/distroless/static-debian12:nonroot):
# no shell, no package manager, non-root user 65532.
FROM gcr.io/distroless/static-debian12@sha256:afa5c872c891853ca7fcf1f12c3edb23f7eeef36189728842dd51042ff57f7ab
COPY --from=build /out/tiefer-web /tiefer-web
ARG PORT=8080
ENV PORT=${PORT} ENV=production
EXPOSE ${PORT}
USER 65532:65532
HEALTHCHECK --interval=30s --timeout=5s --start-period=5s CMD ["/tiefer-web", "healthcheck"]
ENTRYPOINT ["/tiefer-web"]
