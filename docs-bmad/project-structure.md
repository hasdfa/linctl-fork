# Project Structure

- Repository Type: Monolith
- Parts: 1
- Detected Parts:
  - cli: linctl CLI (Go)

Key Directories
- cmd/        Command definitions (Cobra)
- pkg/        Shared packages/utilities (if present)
- bmad/       BMAD workflows and configs (meta tooling)
- docs-bmad/  Generated documentation output
- docs/       Project docs (if any)
- .github/    CI workflows (GitHub Actions)

Key Files
- main.go     Entry point
- go.mod      Module metadata
- Makefile    Dev/build tasks
- README.md   Product documentation
