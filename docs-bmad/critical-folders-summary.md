# Critical Folders Summary

- cmd/: Cobra command hierarchy; each file defines a command group with flags and `Run` handlers.
- pkg/: Shared helpers (formatting, output, API helpers) used across commands.
- bmad/: Project automation and documentation workflows (non-runtime).
- docs/: API snapshot references for documentation and development context.
- .github/workflows/: CI workflows; currently tap bump automation.
- Formula/: Homebrew packaging metadata.
