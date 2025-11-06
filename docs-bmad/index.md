# linctl - Project Documentation Index

> **Primary AI Retrieval Source**: This index provides comprehensive navigation to all project documentation for AI-assisted development and human reference.

## Project Overview

- **Type**: Monolith (single cohesive CLI application)
- **Primary Language**: Go 1.23.0+
- **Architecture**: Layered CLI Application
- **Framework**: Cobra CLI Framework
- **Purpose**: Comprehensive command-line interface for Linear's API

## Quick Reference

- **Tech Stack**: Go, Cobra, Viper, GraphQL
- **Entry Point**: `main.go`
- **Architecture Pattern**: Layered CLI (UI → Business Logic → API Client → Data)
- **Distribution**: Homebrew tap + Source installation
- **Project Root**: `/Users/johnpates/Documents/github_clones/linctl`

## Generated Documentation

### Core Documentation
- **[Project Overview](./project-overview.md)** - High-level project information, features, and capabilities
- **[Architecture](./architecture.md)** - Detailed technical architecture, component design, and system structure
- **[Technology Stack](./technology-stack.md)** - Complete technology inventory with versions and purposes
- **[Project Structure](./project-structure.md)** - Repository organization and structural metadata

### Source Code Analysis
- **[Source Tree Analysis](./source-tree-analysis.md)** - Annotated directory tree with purpose descriptions
- **[Critical Folders Summary](./critical-folders-summary.md)** - Key directories and their roles explained
- **[Architecture Patterns](./architecture-patterns.md)** - Design patterns and architectural decisions
- **[Comprehensive Analysis](./comprehensive-analysis-cli.md)** - Detailed codebase analysis

### Development & Operations
- **[Development Guide](./development-guide.md)** - Setup, build process, development workflow, and common tasks
- **[Development Instructions](./development-instructions.md)** - Detailed development procedures and best practices
- **[Deployment Guide](./deployment-guide.md)** - Release process, CI/CD pipeline, and distribution channels
- **[Contribution Guide](./contribution-guide.md)** - Contributing guidelines, code style, and PR process

## Existing Project Documentation

### User Documentation
- **[README.md](../README.md)** - Primary user documentation with features, installation, and usage examples
- **[master_api_ref.md](../master_api_ref.md)** - API reference documentation

### Contributor Documentation
- **[CONTRIBUTING.md](../CONTRIBUTING.md)** - Official contribution guidelines and release checklist

## Key Features

### Issue Management
- List, create, update, and search issues
- Sub-issue hierarchy and relationships
- Full-text search via Linear's API
- Time-based filtering and pagination

### Team & Project Management
- Team listing and member management
- Project tracking with progress visualization
- User management and profiles
- Comment operations

### Output Formats
- **Table**: ASCII tables for terminal display
- **Plaintext**: Markdown-style output
- **JSON**: Machine-readable for AI agents and automation

### Performance & Flexibility
- Default 6-month time filtering for performance
- Configurable pagination (default: 50 items)
- Multiple sorting options (linear, created, updated)
- Environment variable and config file support

## Technology Stack Summary

| Category | Technology | Version | Purpose |
|----------|-----------|---------|---------|
| **Language** | Go | 1.23.0+ | Core implementation |
| **CLI Framework** | Cobra | 1.8.0 | Command structure |
| **Configuration** | Viper | 1.18.2 | Config management |
| **Output** | tablewriter | 0.0.5 | Table formatting |
| **Terminal UI** | fatih/color | 1.16.0 | Colored output |
| **API** | Linear GraphQL | - | Backend integration |

## Architecture Highlights

### Command Structure
```
linctl (root)
├── auth            # Authentication
├── issue           # Issue management (list, get, create, update, search)
├── team            # Team operations (list, get, members)
├── project         # Project tracking (list, get)
├── user            # User management (list, get, me)
├── comment         # Comments (list, create)
└── docs            # Built-in documentation
```

### Directory Structure
```
linctl/
├── cmd/           # Cobra command definitions
├── internal/      # Internal packages (API, config, output)
│   ├── api/       # Linear API client
│   ├── config/    # Configuration management
│   └── output/    # Output formatters
├── tests/         # Test files and smoke tests
├── main.go        # Application entry point
└── Makefile       # Build automation
```

## Getting Started

### Installation
```bash
# Homebrew (recommended)
brew tap dorkitude/linctl
brew install linctl

# From source
git clone https://github.com/dorkitude/linctl.git
cd linctl
make deps && make build && make install
```

### Quick Start
```bash
linctl auth                              # Authenticate
linctl whoami                            # Verify
linctl issue list --assignee me          # List your issues
linctl issue list --assignee me --json   # JSON output for AI agents
```

### For AI Agents
**IMPORTANT**: Always use `--json` flag for read operations
```bash
linctl issue list --json
linctl team list --json
linctl project list --team ENG --json
```

## Development Workflow

### Build & Test
```bash
make deps          # Install dependencies
make build         # Build binary
make test          # Run smoke tests
make fmt           # Format code
make lint          # Run linter
```

### Common Tasks
- **Local development**: `go run main.go [command]`
- **Install locally**: `make install` (installs to /usr/local/bin)
- **Cross-compile**: `make build-all`
- **Release**: Tag → GitHub Release → Homebrew auto-update

## Configuration

### Files
- **Config**: `~/.linctl.yaml` (output format, limits, API settings)
- **Auth**: `~/.linctl-auth.json` (Personal API key storage)

### Environment Variables
- `LINEAR_API_KEY`: Override config file authentication
- `LINEAR_TEST_API_KEY`: For integration testing

## Testing

- **Smoke Tests**: `make test` (all read-only commands)
- **Manual Testing**: Test against real Linear workspace
- **Integration Tests**: Use `LINEAR_TEST_API_KEY` environment variable

## CI/CD Pipeline

### GitHub Actions
- **Workflow**: `.github/workflows/bump-tap.yml`
- **Trigger**: On GitHub release publication
- **Actions**: Computes SHA256, updates Homebrew formula, opens PR
- **Required Secret**: `HOMEBREW_TAP_TOKEN`

## Use Cases

### Primary Audiences
1. **Developers**: Terminal-based issue tracking
2. **AI Agents**: JSON output for automation (Claude Code, Cursor, Gemini)
3. **DevOps**: CI/CD integration and scripting
4. **Project Managers**: Quick status checks

### Common Workflows
- Daily standup: `linctl issue list --assignee me --newer-than 1_week_ago`
- Create issue: `linctl issue create --title "Fix bug" --team ENG --assign-me`
- Search: `linctl issue search "authentication" --team ENG`
- Project status: `linctl project list --team ENG --json | jq '.[] | {name, progress}'`

## Links

- **Repository**: https://github.com/dorkitude/linctl
- **Issues**: https://github.com/dorkitude/linctl/issues
- **Homebrew Tap**: https://github.com/dorkitude/homebrew-linctl
- **Linear API**: https://developers.linear.app/

## Documentation Metadata

- **Generated**: 2025-11-06
- **Workflow**: BMAD document-project (v1.2.0)
- **Mode**: initial_scan
- **Scan Level**: deep
- **Documentation Root**: `/Users/johnpates/Documents/github_clones/linctl/docs-bmad`

## For Brownfield PRD Development

When creating a brownfield PRD for new features:

1. **Start here**: Use this index as your primary reference
2. **Architecture**: Review [architecture.md](./architecture.md) for system design
3. **Codebase**: Check [source-tree-analysis.md](./source-tree-analysis.md) for file locations
4. **Development**: Follow [development-guide.md](./development-guide.md) for implementation patterns
5. **Testing**: Reference [development-guide.md](./development-guide.md#testing) for test strategy

### Feature-Specific References

**For new commands**:
- See `cmd/` directory structure in [source-tree-analysis.md](./source-tree-analysis.md)
- Review command patterns in [architecture.md](./architecture.md#command-structure)

**For API changes**:
- Review API client in [architecture.md](./architecture.md#api-design)
- Check `internal/api/` patterns in [source-tree-analysis.md](./source-tree-analysis.md)

**For output format changes**:
- Review formatters in [architecture.md](./architecture.md#component-overview)
- Check `internal/output/` patterns

**For configuration changes**:
- Review config management in [architecture.md](./architecture.md#data-architecture)
- Check `internal/config/` patterns

---

**Last Updated**: 2025-11-06
**Status**: ✅ Complete and ready for AI-assisted development
