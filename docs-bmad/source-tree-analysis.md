# Source Tree Analysis

```
linctl/
├── main.go                  # Entry point; embeds README and runs cmd.Execute()
├── go.mod                   # Go module and toolchain version
├── Makefile                 # Build, test, lint, fmt targets
├── README.md                # Primary project documentation
├── cmd/                     # Cobra command tree (top-level commands and subcommands)
│   ├── root.go              # Root command configuration and common flags
│   ├── docs.go              # `linctl docs` renderer
│   ├── auth.go              # `linctl auth` (login/status/logout)
│   ├── issue.go             # `linctl issue` (list/search/get/create/update/assign)
│   ├── project.go           # `linctl project` (list/get)
│   ├── team.go              # `linctl team` (list/get/members)
│   ├── user.go              # `linctl user` (list/get/me)
│   └── comment.go           # `linctl comment` (list/create)
├── pkg/                     # Shared packages/utilities (supporting code)
├── bmad/                    # BMAD workflows, tasks, and docs (meta tooling)
├── docs/                    # Snapshot docs (upstream API references)
│   └── snapshots/           # Various Linear API snapshots
├── .github/
│   └── workflows/
│       └── bump-tap.yml     # Brew tap bump workflow
└── Formula/                 # Homebrew formula
```
