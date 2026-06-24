# CLAUDE.md

Guidance for AI agents (and humans) working in `mochi-toolkit`.

## What this is

`mochi-toolkit` is a Go service in the Console Labs / Mochi backend (runs on EKS `mochi-prod`; see Console Labs MAP.md). Module `github.com/consolelabs/mochi-toolkit`, Go 1.20, gitflow.

Entrypoints:
- (single binary / no cmd/ dir)

## Commands

- Build: `go build ./...` · Test: `go test ./...`
- Run: `go run ./cmd/server`


## Conventions

- Config via env / the platform; NEVER hardcode secrets as defaults (a sibling service, mochi-api, had committed API keys , do not repeat that).
- gitflow: feature branches off the default branch.

## Security / quality (consolidation hardening pass, 2026-06-25, lighter adoption)

- gitleaks: no leaks found (working tree).
- Secret scan in CI: `.github/workflows/security.yml` runs gitleaks (with `.gitleaks.toml` allowlist) + govulncheck on PRs.
- Dependency audit: govulncheck in CI (clean Go toolchain there); Dependabot enabled. Bump deliberately , live service.

