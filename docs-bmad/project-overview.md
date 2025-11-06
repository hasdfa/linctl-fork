# Project Overview - linctl

## Project Identity

**Name**: linctl
**Type**: Command-Line Interface (CLI) Application
**Purpose**: Comprehensive CLI tool for interacting with Linear's API, designed for both human users and AI agents
**Repository**: https://github.com/dorkitude/linctl
**License**: MIT

## Executive Summary

linctl is a feature-rich command-line interface for Linear's project management platform. It provides comprehensive functionality for managing issues, teams, projects, users, and comments through an intuitive CLI experience. Built with Go and the Cobra framework, it emphasizes performance, flexibility, and ease of automation.

The tool supports multiple output formats (table, plaintext, JSON) making it suitable for both interactive terminal use and scripting/automation scenarios. Special attention has been paid to AI agent compatibility, with JSON output and comprehensive documentation.

## Technology Stack Summary

| Layer | Technology | Version |
|-------|-----------|---------|
| **Language** | Go | 1.23.0+ |
| **CLI Framework** | Cobra | 1.8.0 |
| **Configuration** | Viper | 1.18.2 |
| **Output** | tablewriter | 0.0.5 |
| **Terminal UI** | fatih/color | 1.16.0 |
| **API** | Linear GraphQL | - |

## Architecture Type

**Classification**: Monolithic CLI Application

**Pattern**: Layered Architecture
- **UI Layer**: Commands, flags, output formatters
- **Business Logic**: Command handlers, data processing
- **API Client**: GraphQL query/mutation handling
- **Data Layer**: Configuration and authentication storage

## Repository Structure

**Type**: Monolith
**Parts**: 1 (single cohesive CLI application)
**Entry Point**: `main.go`

```
linctl/
├── cmd/           # Command definitions (Cobra)
├── internal/      # Internal packages (API, config, output)
├── tests/         # Test files and smoke tests
├── main.go        # Application entry point
└── Makefile       # Build automation
```

## Core Features

### Issue Management
- List issues with advanced filtering
- Full-text search via Linear's search API
- Create, update, and assign issues
- View issue details with hierarchies
- Parent/child issue relationships
- Time-based filtering (default: 6 months)

### Team Management
- List all teams in workspace
- View team details and statistics
- List team members with roles
- Team-based filtering for issues/projects

### Project Tracking
- List projects with progress visualization
- View detailed project information
- Filter by team, state, and time
- Initiative and milestone tracking

### User Management
- List workspace users
- View user details and profiles
- Active/inactive user filtering
- Current user information

### Comment Operations
- List issue comments with timestamps
- Create new comments
- Time-aware formatting
- Comment author information

### Output Formats
- **Table**: Human-readable ASCII tables
- **Plaintext**: Markdown-style output
- **JSON**: Machine-readable for automation

## Key Capabilities

### Performance Optimization
- Default 6-month time filtering to reduce data transfer
- Configurable pagination limits
- Connection pooling for API requests
- Efficient GraphQL queries (selective fields)

### Flexibility
- Multiple output formats via flags
- Configurable defaults via `~/.linctl.yaml`
- Time expression support (e.g., `2_weeks_ago`, `all_time`)
- Comprehensive filtering options

### Developer Experience
- Built-in documentation via `linctl docs`
- Comprehensive help text for all commands
- Clear error messages with suggested fixes
- Smoke test suite for reliability

### Automation Support
- JSON output for scripting
- Exit codes for error handling
- Environment variable support
- Non-interactive mode for CI/CD

## Installation Methods

1. **Homebrew (Recommended)**
   ```bash
   brew tap dorkitude/linctl
   brew install linctl
   ```

2. **From Source**
   ```bash
   git clone https://github.com/dorkitude/linctl.git
   cd linctl
   make deps && make build && make install
   ```

## Configuration

**Config File**: `~/.linctl.yaml`
- Default output format
- Pagination limits
- API timeout and retry settings

**Authentication**: `~/.linctl-auth.json`
- Personal API key storage
- Secure local file permissions

## Documentation Structure

This documentation suite includes:

- **[Project Overview](./project-overview.md)** _(this file)_ - High-level project information
- **[Architecture](./architecture.md)** - Detailed technical architecture
- **[Source Tree Analysis](./source-tree-analysis.md)** - Directory structure and organization
- **[Critical Folders Summary](./critical-folders-summary.md)** - Key directories explained
- **[Development Guide](./development-guide.md)** - Setup, build, and development workflow
- **[Deployment Guide](./deployment-guide.md)** - Release process and distribution
- **[Contribution Guide](./contribution-guide.md)** - Contributing guidelines

## Getting Started

### Quick Start
```bash
# Install via Homebrew
brew tap dorkitude/linctl
brew install linctl

# Authenticate
linctl auth

# Verify installation
linctl whoami

# View documentation
linctl docs

# List your issues
linctl issue list --assignee me
```

### For AI Agents
```bash
# Use JSON output for all read operations
linctl issue list --assignee me --json
linctl project list --team ENG --json
linctl team members ENG --json
```

### Common Workflows

**Daily standup helper**:
```bash
linctl issue list --assignee me --newer-than 1_week_ago --sort updated
```

**Create and assign issue**:
```bash
linctl issue create --title "Fix bug" --team ENG --assign-me
```

**Search for specific issue**:
```bash
linctl issue search "authentication" --team ENG
```

**Get project status**:
```bash
linctl project list --team ENG --json | jq '.[] | {name, progress}'
```

## Development Status

**Current State**: Active development
**Stability**: Production-ready
**Test Coverage**: Smoke tests for all read-only commands
**Distribution**: Homebrew tap with automated updates

### Recent Features
- Full-text issue search
- Sub-issue hierarchy display
- Git branch integration
- Project tagging (in progress)
- Comprehensive time-based filtering

### Planned Enhancements
- Webhook management improvements
- Bulk operations support
- Interactive TUI mode
- Configuration profiles for multiple workspaces

## Target Audience

### Primary Users
1. **Developers**: Command-line workflow for issue tracking
2. **AI Agents**: Automation via JSON output (Claude Code, Cursor, Gemini)
3. **DevOps Engineers**: CI/CD integration and scripting
4. **Project Managers**: Quick project status checks

### Use Cases
- Issue tracking from terminal
- Automated workflows and scripts
- CI/CD pipeline integration
- Team metrics and reporting
- Personal productivity dashboards

## Links

- **Repository**: https://github.com/dorkitude/linctl
- **Issues**: https://github.com/dorkitude/linctl/issues
- **Homebrew Tap**: https://github.com/dorkitude/homebrew-linctl
- **Linear API**: https://developers.linear.app/

## Project Metrics

**Language**: Go
**Commands**: 25+ (across 7 command groups)
**Output Formats**: 3 (table, plaintext, JSON)
**Distribution Channels**: 2 (Homebrew, source)
**Core Dependencies**: 4 (cobra, viper, tablewriter, color)

---

**Last Updated**: 2025-11-06
**Documentation Generated**: BMAD document-project workflow
