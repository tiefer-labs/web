# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.

# Build stage: the official Go image.
FROM golang:1.27 AS build
WORKDIR /src
COPY go.mod ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/tiefer-web ./cmd/tiefer-web

# Run stage: distroless static image, no shell, non-root user.
FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=build /out/tiefer-web /tiefer-web
ARG PORT=8080
ENV PORT=${PORT} ENV=production
EXPOSE ${PORT}
USER nonroot:nonroot
HEALTHCHECK --interval=30s --timeout=5s --start-period=5s CMD ["/tiefer-web", "healthcheck"]
ENTRYPOINT ["/tiefer-web"]
