# Security decisions

Every accepted risk, suppressed finding and trade-off, with the reason and
when to revisit it. `docs/HARDENING.md` was not available during the
rebuild (see `docs/PROGRESS.md`); every entry marked **check against
HARDENING.md** must be compared with that document once it exists.

## Suppressed findings

| Where | Finding | Why it is safe | Revisit |
|---|---|---|---|
| `internal/server/static.go` | gosec G705 (XSS via taint) on writing a static asset | The bytes are an embedded file found by an exact map lookup on the hashed URL; the request only selects which file, it never becomes content. Content types are fixed per file and `nosniff` is set. | When static serving changes |
