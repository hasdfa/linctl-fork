# Contribution Guide - linctl

Thanks for contributing! This repository aims to keep changes simple, focused, and tested.

## Getting Started

### Requirements
- **Go**: 1.22+ (1.23.0+ recommended)
- **Git**: For version control
- **golangci-lint**: Optional, for linting
- **jq**: Optional, for JSON parsing in examples
- **gh CLI**: Optional, for GitHub operations

### Development Setup
```bash
# Clone repository
git clone https://github.com/dorkitude/linctl.git
cd linctl

# Install dependencies
make deps

# Build and test
make build
make test
```

## Development Workflow

### Making Changes

1. **Fork the repository**
   ```bash
   gh repo fork dorkitude/linctl --clone
   cd linctl
   ```

2. **Create a feature branch**
   ```bash
   git checkout -b feature/your-feature-name
   ```

3. **Make your changes**
   - Edit relevant files
   - Follow Go conventions
   - Keep changes focused

4. **Format and lint**
   ```bash
   make fmt   # Format code
   make lint  # Run linter (if golangci-lint installed)
   ```

5. **Test your changes**
   ```bash
   make test  # Run smoke tests
   # Manually test against Linear API if needed
   ```

6. **Commit your changes**
   ```bash
   git add .
   git commit -m "Brief description of changes"
   ```

7. **Push and create PR**
   ```bash
   git push origin feature/your-feature-name
   gh pr create --title "Your PR Title" --body "Description of changes"
   ```

## Code Style Guidelines

### Go Conventions
- Follow standard Go formatting (`gofmt`)
- Use meaningful variable names
- Add comments for exported functions
- Keep functions focused and small
- Handle errors explicitly

### Project-Specific Patterns
- Use Cobra for CLI commands
- Use Viper for configuration
- Consistent error messages
- Support table, plaintext, and JSON output

## Testing Requirements

### Smoke Tests
All read-only commands are covered by smoke tests in `smoke_test.sh`:
```bash
make test         # Run smoke tests
make test-verbose # Run with verbose output
```

### Manual Testing
Test new features against your Linear workspace:
```bash
linctl auth                    # Authenticate first
linctl [your-new-command]      # Test your changes
```

### What to Test
- Command works as expected
- Error handling works correctly
- All output formats (table, plaintext, JSON)
- Help text is clear and accurate
- Edge cases and error conditions

## Documentation Requirements

### Code Documentation
- Add godoc comments for exported functions
- Explain complex logic with inline comments
- Update relevant docstrings

### User Documentation
When adding features, update:
1. **README.md**: Add command examples and documentation
2. **Help Text**: Update command help in Cobra definitions
3. **Examples**: Add real-world usage examples

### Documentation Checklist
- [ ] Help text matches actual behavior
- [ ] README.md updated with new features
- [ ] Examples added for new commands
- [ ] Configuration changes documented

## Pull Request Guidelines

### PR Title Format
Use clear, descriptive titles:
- ✅ "Add project tagging feature"
- ✅ "Fix issue list pagination"
- ❌ "Updates"
- ❌ "Bug fix"

### PR Description
Include:
- **What**: What does this PR do?
- **Why**: Why is this change needed?
- **How**: How does it work?
- **Testing**: How was it tested?
- **Breaking Changes**: Any breaking changes?

### PR Checklist
- [ ] Code follows project style guidelines
- [ ] `make fmt` ran successfully
- [ ] `make test` passes
- [ ] Documentation updated
- [ ] Breaking changes noted
- [ ] Commit messages are clear

## Release Process

See [Deployment Guide](./deployment-guide.md) for complete release checklist.

### Quick Release Steps
1. Prepare (tests pass, docs updated)
2. Tag: `git tag vX.Y.Z -a -m "vX.Y.Z: summary"`
3. Push: `git push origin vX.Y.Z`
4. Release: `gh release create vX.Y.Z --title "linctl vX.Y.Z" --notes "..."`
5. Homebrew: Auto-updated via GitHub Action
6. Validate: `brew upgrade linctl && linctl --version`

## Useful Make Targets

```bash
make build        # Build the binary
make clean        # Clean build artifacts
make test         # Run smoke tests
make test-verbose # Run smoke tests with verbose output
make deps         # Install/update dependencies
make fmt          # Format code
make lint         # Lint code (if golangci-lint installed)
make install      # Install to /usr/local/bin
make dev-install  # Create development symlink
make build-all    # Cross-compile for all platforms
make release      # Prepare release builds
make run          # Build and run
make everything   # Build, format, lint, test, install
make help         # Show all available targets
```

## Common Tasks

### Adding a New Command

1. **Create command file**
   ```go
   // cmd/mycommand.go
   package cmd

   import "github.com/spf13/cobra"

   var myCmd = &cobra.Command{
       Use:   "mycommand",
       Short: "Brief description",
       RunE:  runMyCommand,
   }

   func init() {
       rootCmd.AddCommand(myCmd)
   }

   func runMyCommand(cmd *cobra.Command, args []string) error {
       // Implementation
       return nil
   }
   ```

2. **Add to root command** (if not using `init()`)
3. **Update README.md** with usage examples
4. **Test manually** and with smoke tests

### Adding a New Output Format

1. Edit `internal/output/` formatters
2. Ensure consistency across formats
3. Test with various data sizes
4. Update documentation

### Updating Dependencies

```bash
# Update specific dependency
go get -u github.com/spf13/cobra

# Update all dependencies
go get -u ./...

# Tidy modules
make deps
```

## Communication

### Getting Help
- **Issues**: [GitHub Issues](https://github.com/dorkitude/linctl/issues)
- **Discussions**: GitHub Discussions (if enabled)
- **Questions**: Open an issue with "Question" label

### Reporting Bugs
Include:
- linctl version (`linctl --version`)
- Go version (`go version`)
- OS and architecture
- Steps to reproduce
- Expected vs actual behavior
- Error messages/output

### Suggesting Features
Open an issue with:
- Clear description of feature
- Use case / motivation
- Example usage
- Potential implementation approach

## Code of Conduct

### Be Respectful
- Respectful communication
- Constructive feedback
- Welcoming to newcomers
- Focus on technical merit

### Be Professional
- Keep discussions on-topic
- Assume good intentions
- Acknowledge contributions
- Help others learn

## Resources

- [Go Documentation](https://golang.org/doc/)
- [Cobra Documentation](https://github.com/spf13/cobra)
- [Linear API Documentation](https://developers.linear.app/)
- [GitHub Repository](https://github.com/dorkitude/linctl)
- [Issue Tracker](https://github.com/dorkitude/linctl/issues)

## License

By contributing, you agree that your contributions will be licensed under the MIT License.

---

**Thank you for contributing to linctl!** 🎉
