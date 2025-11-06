# Architecture Documentation - linctl

## Executive Summary

**linctl** is a command-line interface (CLI) tool for interacting with Linear's API. Built with Go and the Cobra CLI framework, it provides comprehensive functionality for issue management, team operations, project tracking, and user management. The tool is designed with both human users and AI agents in mind, offering multiple output formats (table, plaintext, JSON) for different use cases.

**Architecture Type**: Command-Line Application (Monolithic)
**Primary Language**: Go 1.23.0+
**Core Framework**: Cobra CLI Framework
**Distribution**: Homebrew tap and source installation

## Technology Stack

| Category | Technology | Version | Purpose |
|----------|-----------|---------|---------|
| **Language** | Go | 1.23.0+ (toolchain 1.24.5) | Core implementation language |
| **CLI Framework** | Cobra | 1.8.0 | Command structure and routing |
| **Configuration** | Viper | 1.18.2 | Configuration management |
| **Output Formatting** | tablewriter | 0.0.5 | Table-based output formatting |
| **Terminal Colors** | fatih/color | 1.16.0 | Colored terminal output |
| **API** | Linear GraphQL API | N/A | Backend data source |

### Key Dependencies
- **spf13/cobra**: Powerful CLI framework with command structure
- **spf13/viper**: Configuration file and environment variable management
- **olekukonko/tablewriter**: ASCII table generation for terminal output
- **fatih/color**: Cross-platform colored terminal text

## Architecture Pattern

### CLI Application Architecture

linctl follows a **layered CLI architecture** pattern:

```
┌─────────────────────────────────────────┐
│          User Interface Layer           │
│  (Commands, Flags, Output Formatters)   │
└──────────────┬──────────────────────────┘
               │
┌──────────────▼──────────────────────────┐
│         Business Logic Layer            │
│   (Command Handlers, Data Processing)   │
└──────────────┬──────────────────────────┘
               │
┌──────────────▼──────────────────────────┐
│          API Client Layer               │
│    (Linear GraphQL API Integration)     │
└──────────────┬──────────────────────────┘
               │
┌──────────────▼──────────────────────────┐
│           Data Layer                    │
│  (Config Files, Auth Storage, Cache)    │
└─────────────────────────────────────────┘
```

### Key Architectural Characteristics

1. **Command-Based Structure**: Each feature area (issue, team, project, user, comment) is organized as a command group
2. **Output Format Abstraction**: Support for multiple output formats through a consistent interface
3. **Configuration Management**: Centralized configuration via Viper with file and environment variable support
4. **Stateless Operations**: Each command execution is independent, with authentication persisted locally
5. **API Client Abstraction**: GraphQL API interactions isolated in dedicated layer

## Component Overview

### Command Structure

```
linctl (root)
├── auth            # Authentication management
│   ├── login       # Interactive login
│   ├── status      # Check auth status
│   └── logout      # Clear credentials
├── whoami          # Current user info
├── issue           # Issue management
│   ├── list        # List issues with filters
│   ├── get         # Get issue details
│   ├── create      # Create new issue
│   ├── update      # Update issue fields
│   ├── assign      # Assign issue to user
│   └── search      # Full-text search
├── team            # Team management
│   ├── list        # List teams
│   ├── get         # Get team details
│   └── members     # List team members
├── project         # Project tracking
│   ├── list        # List projects
│   └── get         # Get project details
├── user            # User management
│   ├── list        # List users
│   ├── get         # Get user details
│   └── me          # Current user profile
├── comment         # Comment operations
│   ├── list        # List issue comments
│   └── create      # Add comment
└── docs            # Built-in documentation
```

### Core Packages

#### `/cmd` - Command Definitions
- **root.go**: Root command, global flags, initialization
- **auth.go**: Authentication commands
- **issue.go**: Issue management commands
- **team.go**: Team operation commands
- **project.go**: Project tracking commands
- **user.go**: User management commands
- **comment.go**: Comment commands
- **version.go**: Version information

#### `/internal/api` - API Client Layer
- GraphQL query construction
- HTTP client configuration
- Request/response handling
- Error handling and retries
- Rate limiting awareness

#### `/internal/config` - Configuration Management
- Config file loading (`~/.linctl.yaml`)
- Authentication storage (`~/.linctl-auth.json`)
- Default value management
- Environment variable support

#### `/internal/output` - Output Formatting
- **Table Formatter**: ASCII table generation for terminal
- **Plaintext Formatter**: Markdown-style plain text
- **JSON Formatter**: Machine-readable JSON output
- Format selection via `--json` or `--plaintext` flags

## Data Architecture

### Configuration Files

**~/.linctl.yaml** (User Configuration):
```yaml
output: table           # Default output format
limit: 50              # Default pagination limit
api:
  timeout: 30s         # API request timeout
  retries: 3           # Retry attempts
```

**~/.linctl-auth.json** (Authentication):
```json
{
  "api_key": "lin_api_..."
}
```

### Data Flow

1. **Command Invocation** → User executes command with flags
2. **Configuration Loading** → Viper loads config and auth
3. **API Request** → GraphQL query sent to Linear API
4. **Response Processing** → Data parsed and transformed
5. **Output Formatting** → Data formatted per user preference
6. **Display** → Output rendered to terminal

## API Design

### Linear GraphQL API Integration

**Endpoint**: `https://api.linear.app/graphql`

**Authentication**: Personal API Key (Bearer token)

**Request Pattern**:
```graphql
query {
  issues(filter: {...}, first: 50) {
    nodes {
      id
      identifier
      title
      state { name }
      assignee { name email }
      team { key name }
      priority
      createdAt
      updatedAt
    }
  }
}
```

### Key API Operations

| Operation | GraphQL Type | Description |
|-----------|-------------|-------------|
| List Issues | Query | Fetch issues with filters |
| Get Issue | Query | Fetch single issue details |
| Create Issue | Mutation | Create new issue |
| Update Issue | Mutation | Modify issue fields |
| List Teams | Query | Fetch all teams |
| Get Team | Query | Fetch team details |
| List Projects | Query | Fetch projects with filters |
| Get Project | Query | Fetch project details |
| List Users | Query | Fetch workspace users |
| Get User | Query | Fetch user details |
| List Comments | Query | Fetch issue comments |
| Create Comment | Mutation | Add issue comment |

## Source Tree Structure

```
linctl/
├── cmd/                      # Cobra command definitions
│   ├── root.go               # Root command + global flags
│   ├── auth.go               # Authentication commands
│   ├── issue.go              # Issue management
│   ├── team.go               # Team operations
│   ├── project.go            # Project tracking
│   ├── user.go               # User management
│   ├── comment.go            # Comment operations
│   └── version.go            # Version info
├── internal/                 # Internal packages
│   ├── api/                  # Linear API client
│   │   ├── client.go         # HTTP client setup
│   │   ├── queries.go        # GraphQL queries
│   │   └── mutations.go      # GraphQL mutations
│   ├── config/               # Configuration management
│   │   ├── config.go         # Config loading/saving
│   │   └── auth.go           # Auth file management
│   └── output/               # Output formatters
│       ├── table.go          # Table formatter
│       ├── plaintext.go      # Plaintext formatter
│       └── json.go           # JSON formatter
├── tests/                    # Test files
│   ├── smoke_test.sh         # Smoke tests for commands
│   └── README.md             # Test documentation
├── .github/                  # GitHub workflows
│   └── workflows/
│       └── bump-tap.yml      # Homebrew tap auto-update
├── main.go                   # Application entry point
├── Makefile                  # Build automation
├── go.mod                    # Go module definition
├── go.sum                    # Dependency checksums
├── README.md                 # User documentation
├── CONTRIBUTING.md           # Contribution guide
└── LICENSE                   # MIT License
```

## Development Workflow

### Build Process
```
Source Files → Go Build → Binary (linctl) → Installation
```

**Build Commands**:
- `make build`: Compile for current platform
- `make build-all`: Cross-compile for all platforms
- `make release`: Prepare release artifacts

**Version Injection**:
```makefile
VERSION=$(shell git describe --tags --exact-match 2>/dev/null || git rev-parse --short HEAD)
LDFLAGS=-ldflags "-X github.com/dorkitude/linctl/cmd.version=$(VERSION)"
```

### Testing Strategy

**Smoke Tests** (`smoke_test.sh`):
- Automated testing of all read-only commands
- Verify help text accessibility
- Validate output format consistency
- Check error handling

**Manual Testing**:
- Integration testing against real Linear API
- User workflow validation
- Edge case testing

## Deployment Architecture

### Distribution Methods

1. **Homebrew (Primary)**
   - Custom tap: `dorkitude/homebrew-linctl`
   - Auto-update via GitHub Actions
   - Install: `brew install dorkitude/linctl/linctl`

2. **Source Installation**
   - Clone repository
   - `make deps && make build && make install`
   - Binary installed to `/usr/local/bin/linctl`

### CI/CD Pipeline

**GitHub Actions Workflow** (`.github/workflows/bump-tap.yml`):
- Triggered on GitHub release publication
- Computes tarball SHA256
- Updates Homebrew formula
- Opens PR to tap repository
- Requires `HOMEBREW_TAP_TOKEN` secret

### Installation Paths

- **Binary**: `/usr/local/bin/linctl`
- **Config**: `~/.linctl.yaml`
- **Auth**: `~/.linctl-auth.json`

## Security Architecture

### Authentication
- Personal API keys stored locally in `~/.linctl-auth.json`
- No password storage or OAuth flows
- API key never logged or displayed
- File permissions restrict access to user only

### API Communication
- HTTPS only (TLS 1.2+)
- Bearer token authentication
- Request timeout protection
- Retry logic with backoff

### Secrets Management
- No secrets in source code
- Environment variables for CI/CD
- GitHub Secrets for automation tokens
- Fine-grained PAT permissions for tap updates

## Performance Considerations

### Optimization Strategies

1. **Default Time Filtering**: 6-month lookback by default to reduce data transfer
2. **Pagination**: Configurable limit (default 50) to control response size
3. **Selective Queries**: Request only needed fields from GraphQL
4. **Connection Reuse**: HTTP client configured for connection pooling
5. **Timeout Management**: 30-second default timeout prevents hanging

### Rate Limiting
- Linear API: 5,000 requests/hour for Personal API Keys
- No client-side rate limiting implemented (relies on API responses)
- Retry logic respects 429 status codes

## Error Handling

### Error Categories

1. **Authentication Errors**: Missing/invalid API key
2. **Network Errors**: Connection failures, timeouts
3. **API Errors**: GraphQL errors, validation failures
4. **Input Errors**: Invalid flags, malformed input
5. **System Errors**: File I/O issues, permission problems

### Error Handling Strategy
- Clear, actionable error messages
- Suggest remediation steps
- Exit codes for scripting
- Verbose mode for debugging

## Future Considerations

### Potential Enhancements
- Offline mode with local caching
- Bulk operations for efficiency
- Webhook listening capabilities
- Interactive TUI mode
- Plugin system for extensions
- Configuration profiles for multiple workspaces

### Scalability
- Current design suitable for single-user CLI
- API client could support connection pooling
- Output formatters extensible for new formats
- Command structure supports unlimited command additions

## References

- [Linear API Documentation](https://developers.linear.app/)
- [Cobra CLI Framework](https://github.com/spf13/cobra)
- [Go Documentation](https://golang.org/doc/)
- [GitHub Repository](https://github.com/dorkitude/linctl)
- [Development Guide](./development-guide.md)
- [Deployment Guide](./deployment-guide.md)
