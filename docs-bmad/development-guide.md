# Development Guide - linctl

## Prerequisites

### Required
- **Go**: 1.23.0+ (toolchain: 1.24.5)
- **Git**: For version control

### Optional
- **golangci-lint**: For code linting
- **jq**: For JSON parsing in examples
- **gh CLI**: For GitHub operations

## Installation

### From Source
```bash
git clone https://github.com/dorkitude/linctl.git
cd linctl
make deps        # Install dependencies
make build       # Build the binary
make install     # Install to /usr/local/bin (requires sudo)
```

### For Development
```bash
git clone https://github.com/dorkitude/linctl.git
cd linctl
make deps        # Install dependencies
go run main.go   # Run directly without building
# OR
make dev         # Build and run in development mode
```

## Development Workflow

### Build Commands
```bash
make build        # Build the binary
make clean        # Clean build artifacts
make build-all    # Cross-compile for multiple platforms
make release      # Prepare release builds
```

### Code Quality
```bash
make fmt          # Format code with go fmt
make lint         # Run golangci-lint (if installed)
```

### Testing
```bash
make test         # Run smoke tests (read-only commands)
make test-verbose # Run smoke tests with verbose output
```

### Running Locally
```bash
# Option 1: Build and run
make build
./linctl [command]

# Option 2: Run directly with Go
go run main.go [command]

# Option 3: Development symlink
make dev-install  # Creates symlink in /usr/local/bin
linctl [command]
```

## Project Structure

```
linctl/
├── cmd/                # Cobra command definitions
│   ├── root.go         # Root command and global flags
│   ├── auth.go         # Authentication commands
│   ├── issue.go        # Issue management commands
│   ├── team.go         # Team commands
│   ├── project.go      # Project commands
│   ├── user.go         # User commands
│   ├── comment.go      # Comment commands
│   └── version.go      # Version information
├── internal/           # Internal packages
│   ├── api/            # Linear API client
│   ├── config/         # Configuration management
│   └── output/         # Output formatters (table, json, plaintext)
├── tests/              # Test files
│   └── smoke_test.sh   # Smoke tests
├── Makefile            # Build automation
├── go.mod              # Go module dependencies
└── README.md           # Project documentation
```

## Configuration

### Application Config
Configuration is stored in `~/.linctl.yaml`:
```yaml
# Default output format
output: table

# Default pagination limit
limit: 50

# API settings
api:
  timeout: 30s
  retries: 3
```

### Authentication
Authentication credentials are stored securely in `~/.linctl-auth.json`.

## Dependencies

### Core Dependencies
- **cobra**: CLI framework (`github.com/spf13/cobra`)
- **viper**: Configuration management (`github.com/spf13/viper`)
- **color**: Terminal colors (`github.com/fatih/color`)
- **tablewriter**: Table output formatting (`github.com/olekukonko/tablewriter`)

### Development Dependencies
- **golangci-lint**: Linting (optional, external)

See `go.mod` for complete dependency list.

## Common Development Tasks

### Adding a New Command
1. Create command file in `cmd/` directory
2. Define command structure using Cobra
3. Add command to root command in `cmd/root.go`
4. Update README.md with command documentation
5. Add tests if applicable

### Updating API Client
1. Modify files in `internal/api/`
2. Update error handling as needed
3. Test with real Linear API
4. Update output formatters if needed

### Modifying Output Formats
1. Edit formatters in `internal/output/`
2. Ensure table, plaintext, and JSON formats are consistent
3. Test with various data sizes

## Environment Variables

- `LINEAR_API_KEY`: API key for authentication (overrides config file)
- `LINEAR_TEST_API_KEY`: API key for integration tests

## Testing

### Smoke Tests
The project includes automated smoke tests for all read-only commands:
```bash
make test         # Run all smoke tests
make test-verbose # Run with verbose output
```

Tests verify:
- Command execution without errors
- Help text accessibility
- Basic functionality of read operations

### Manual Testing
Test against your own Linear workspace:
```bash
linctl auth                    # Authenticate
linctl issue list              # List issues
linctl team list               # List teams
linctl project list            # List projects
```

## Troubleshooting

### Build Issues
```bash
# Clean and rebuild
make clean
make deps
make build
```

### Module Issues
```bash
go mod tidy
go mod download
```

### Testing Authentication
```bash
linctl auth status    # Check auth status
linctl whoami         # Verify credentials
```

## Best Practices

1. **Code Style**: Always run `make fmt` before committing
2. **Testing**: Run `make test` to verify no regressions
3. **Dependencies**: Keep go.mod tidy with `go mod tidy`
4. **Documentation**: Update README.md when adding features
5. **Versioning**: Follow semantic versioning for releases

## Resources

- [Go Documentation](https://golang.org/doc/)
- [Cobra Documentation](https://github.com/spf13/cobra)
- [Linear API Documentation](https://developers.linear.app/)
- [GitHub Repository](https://github.com/dorkitude/linctl)
