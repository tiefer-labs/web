// SPDX-License-Identifier: MPL-2.0
import type { NextConfig } from "next";
import createNextIntlPlugin from "next-intl/plugin";

const isDev = process.env.NODE_ENV === "development";
const isProduction = process.env.ENV === "production";

// Next.js needs inline scripts and styles; a nonce would force every page to render per request.
const csp = [
  "default-src 'none'",
  `script-src 'self' 'unsafe-inline'${isDev ? " 'unsafe-eval'" : ""}`,
  "script-src-attr 'none'",
  "style-src 'self' 'unsafe-inline'",
  "img-src 'self'",
  "font-src 'self'",
  `connect-src 'self'${isDev ? " ws:" : ""}`,
  "manifest-src 'self'",
  "form-action 'self'",
  "frame-ancestors 'none'",
  "base-uri 'none'",
  "object-src 'none'",
  ...(isProduction ? ["upgrade-insecure-requests"] : []),
].join("; ");

const permissionsPolicy = [
  "accelerometer",
  "autoplay",
  "browsing-topics",
  "camera",
  "display-capture",
  "encrypted-media",
  "fullscreen",
  "geolocation",
  "gyroscope",
  "hid",
  "idle-detection",
  "magnetometer",
  "microphone",
  "midi",
  "payment",
  "picture-in-picture",
  "publickey-credentials-get",
  "screen-wake-lock",
  "serial",
  "usb",
  "xr-spatial-tracking",
]
  .map((feature) => `${feature}=()`)
  .join(", ");

const securityHeaders = [
  { key: "Content-Security-Policy", value: csp },
  {
    key: "Strict-Transport-Security",
    value: `max-age=63072000; includeSubDomains${process.env.HSTS_PRELOAD === "true" ? "; preload" : ""}`,
  },
  { key: "X-Content-Type-Options", value: "nosniff" },
  { key: "X-Frame-Options", value: "DENY" },
  { key: "Referrer-Policy", value: "no-referrer" },
  { key: "Permissions-Policy", value: permissionsPolicy },
  { key: "Cross-Origin-Opener-Policy", value: "same-origin" },
  { key: "Cross-Origin-Resource-Policy", value: "same-origin" },
  { key: "Cross-Origin-Embedder-Policy", value: "require-corp" },
  { key: "Origin-Agent-Cluster", value: "?1" },
  { key: "X-Permitted-Cross-Domain-Policies", value: "none" },
  { key: "X-DNS-Prefetch-Control", value: "off" },
];

const nextConfig: NextConfig = {
  // The Docker image sets NEXT_OUTPUT=standalone; local runs and Vercel use the default output.
  output: process.env.NEXT_OUTPUT === "standalone" ? "standalone" : undefined,
  poweredByHeader: false,
  reactCompiler: true,
  cacheComponents: true,
  partialPrefetching: true,
  typedRoutes: true,
  experimental: {
    globalNotFound: true,
    turbopackRustReactCompiler: true,
    useTypeScriptCli: true,
    agentUpgrade: "latest",
    serverActions: { bodySizeLimit: "16kb" },
  },
  turbopack: {
    // Next has no metrics for the Mozilla fonts yet; globals.css defines hand-tuned fallbacks instead.
    ignoreIssue: [{ path: /.*/, title: /Failed to find font override values/ }],
    rules: {
      "*.css": {
        loaders: ["@tailwindcss/turbopack"],
        as: "*.css",
      },
    },
  },
  async headers() {
    return [{ source: "/:path*", headers: securityHeaders }];
  },
  async redirects() {
    return [{ source: "/security.txt", destination: "/.well-known/security.txt", permanent: true }];
  },
};

const withNextIntl = createNextIntlPlugin({
  experimental: {
    messages: {
      format: "json",
      locales: "infer",
      path: "./messages",
      precompile: true,
    },
  },
});

export default withNextIntl(nextConfig);
