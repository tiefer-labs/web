# syntax=docker/dockerfile:1
# SPDX-License-Identifier: MPL-2.0

# Build stage: pnpm manages the Node runtime, as in our other projects. Node 24 matches the runtime image.
FROM ghcr.io/pnpm/pnpm:12 AS build
RUN pnpm runtime set node 24 -g
WORKDIR /app
ENV NEXT_TELEMETRY_DISABLED=1 NEXT_OUTPUT=standalone
COPY package.json pnpm-lock.yaml pnpm-workspace.yaml ./
# The store lives outside /pnpm/store so the cache mount does not hide the Node runtime.
RUN --mount=type=cache,id=pnpm,target=/var/cache/pnpm pnpm install --store-dir /var/cache/pnpm --frozen-lockfile

# Pages are prerendered here, so the public settings they show are build arguments.
# Secrets (CSRF_SECRET, SMTP_PASS) are runtime only and never passed to the build.
ARG ENV=production
ARG SITE_URL=https://tiefer.space
ARG CONTACT_EMAIL
ARG LINKEDIN_URL
ARG REPO_URL
ARG SMTP_HOST
ARG CONTACT_TO
ARG CONTACT_FROM
ARG LEGAL_NAME
ARG LEGAL_FORM
ARG LEGAL_ADDRESS
ARG LEGAL_TAX_ID
ARG LEGAL_REGISTRATION
ARG LEGAL_DIRECTOR
ARG LEGAL_HOSTING_PROVIDER
ARG LEGAL_HOSTING_COUNTRY
ARG LEGAL_SMTP_PROVIDER
ARG LEGAL_SMTP_COUNTRY
ARG LEGAL_REVIEWED
ARG LOG_RETENTION_DAYS
ARG HSTS_PRELOAD
COPY messages ./messages
COPY public ./public
COPY src ./src
COPY next.config.ts tsconfig.json ./
RUN pnpm next build

# Runtime stage: no shell, no package manager, a non-root user
# (gcr.io/distroless/nodejs24-debian12:nonroot, pinned by digest).
FROM gcr.io/distroless/nodejs24-debian12:nonroot@sha256:14d42e2511532589a7c7e01a753667a74fcc96266e137e8125006b87b0c32d0a
WORKDIR /app
ENV NODE_ENV=production NEXT_TELEMETRY_DISABLED=1 ENV=production PORT=8080 HOSTNAME=0.0.0.0
COPY --from=build --chown=65532:65532 /app/.next/standalone ./
COPY --from=build --chown=65532:65532 /app/.next/static ./.next/static
COPY --from=build --chown=65532:65532 /app/public ./public
USER 65532:65532
EXPOSE 8080
HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 \
  CMD ["/nodejs/bin/node", "-e", "fetch('http://127.0.0.1:8080/healthz').then((r) => process.exit(r.ok ? 0 : 1), () => process.exit(1))"]
CMD ["server.js"]
