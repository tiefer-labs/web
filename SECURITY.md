# Security policy

## Reporting a vulnerability

Please report security problems privately, not in public issues:

- by email to the address in <https://tiefer.space/.well-known/security.txt>;
- or through GitHub: **Security**, then **Report a vulnerability** in this
  repository (private vulnerability reporting).

Please include what you found, where, and how to reproduce it. You may write
in English, Azerbaijani, Turkish or Russian.

We aim to confirm receipt within five working days, keep you informed while
we fix the problem, and credit you when the fix is published if you wish.

## Scope

In scope: the website at `https://tiefer.space` and the code in this
repository.

Out of scope: denial of service by volume, reports from automated scanners
without a demonstrated impact, missing best practices without a security
consequence, and social engineering of people.

Please do not access, change or delete data that is not yours, and do not
degrade the service for other visitors while testing.

## Supported versions

Only the current `main` branch, which is what runs on `https://tiefer.space`,
receives security fixes.

## How the site is protected

The measures are described in the README under "Security". In short: no
cookies, no third-party requests, a strict Content Security Policy with
Trusted Types, HSTS with preload, cross-origin isolation, one canonical
host, a contact form without storage, and no third-party Go modules.
