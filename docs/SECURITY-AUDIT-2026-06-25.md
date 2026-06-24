# Security audit: mochi-toolkit (2026-06-25)

Bounded security triage , Console Labs consolidation hardening pass (lighter adoption).

## Secret scan

gitleaks: no leaks found (working tree).

## Dependency audit

govulncheck added to CI (local Go toolchain could not build it; CI has a clean Go). Dependabot enabled. Findings surfaced non-blocking until triaged; remediation is deliberate (live service).



## What this PR changes

CLAUDE.md (agent guidance) + docs/ARCHITECTURE.md (reindex) + .gitleaks.toml + .github/workflows/security.yml. No source/logic change, no dependency bump.
