# Architecture Patterns

- Pattern: Command-based CLI using Cobra
- Entry Point: `main.go` embeds README and invokes `cmd.Execute()`
- Command Tree: Defined under `cmd/` with subcommands for `issue`, `project`, `team`, `user`, `comment`, `auth`, and `docs`.
- Responsibility Split: Business/API logic encapsulated behind command handlers; output adapters provide table/plaintext/JSON modes.
- External Integration: Linear GraphQL API (per README and API references).
