# mochi-toolkit architecture

Go service for Console Labs / Mochi. Module `github.com/consolelabs/mochi-toolkit`. Runs on EKS `mochi-prod`.

```
mochi-toolkit/
├── (single binary)
├── pkg/ (or internal/)  service + handler + config layers

└── docs/                this architecture + the security audit
```

## Notes for agents

- Live prod service: treat changes as production; verify with `go test ./...`; prefer additive.
- Config + secrets come from env / the platform, not from source defaults.

