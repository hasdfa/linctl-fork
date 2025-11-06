# Deployment Guide - linctl

## Distribution Channels

### Homebrew (Primary Distribution)

#### Automated Release Process
linctl uses GitHub Actions to automatically update the Homebrew tap when a new release is published.

**Workflow**: `.github/workflows/bump-tap.yml`
- Triggers on GitHub release publication
- Automatically computes tarball SHA256
- Opens PR to `dorkitude/homebrew-linctl` tap
- Requires `HOMEBREW_TAP_TOKEN` secret

**Setup**:
1. Create fine-grained PAT with `contents:write` on `dorkitude/homebrew-linctl`
2. Add as repository secret: `HOMEBREW_TAP_TOKEN`
3. Publish GitHub release → Action runs automatically

#### Manual Homebrew Bump (Fallback)
If automation fails or is disabled:

```bash
TAG=vX.Y.Z
TARBALL=https://github.com/dorkitude/linctl/archive/refs/tags/${TAG}.tar.gz
curl -sL "$TARBALL" -o /tmp/linctl.tgz
SHA=$(shasum -a 256 /tmp/linctl.tgz | awk '{print $1}')

git clone https://github.com/dorkitude/homebrew-linctl.git
cd homebrew-linctl
git checkout -b bump-linctl-${TAG#v}
sed -i.bak -E "s|url \"[^\"]+\"|url \"$TARBALL\"|g" Formula/linctl.rb
sed -i.bak -E "s|sha256 \"[0-9a-f]+\"|sha256 \"$SHA\"|g" Formula/linctl.rb
rm -f Formula/linctl.rb.bak
git commit -am "linctl: bump to ${TAG}"
git push -u origin HEAD
gh pr create --title "linctl: bump to ${TAG}" --body "Update formula to ${TAG}." --base master --head bump-linctl-${TAG#v}
```

### From Source Distribution

Users can install directly from source:
```bash
git clone https://github.com/dorkitude/linctl.git
cd linctl
make deps
make build
make install  # Installs to /usr/local/bin
```

## Release Process

### 1. Pre-Release Checklist
- [ ] Ensure README and help text match current behavior
- [ ] Run `make test` to verify smoke tests pass
- [ ] Update version in relevant files
- [ ] Draft release notes (highlights, fixes, breaking changes)

### 2. Create and Push Tag
```bash
git tag vX.Y.Z -a -m "vX.Y.Z: short summary"
git push origin vX.Y.Z
```

### 3. Create GitHub Release
```bash
gh release create vX.Y.Z \
  --title "linctl vX.Y.Z" \
  --notes "<highlights/fixes>"
```

### 4. Homebrew Tap Update
- **Automated**: GitHub Action opens PR automatically
- **Manual**: Follow manual bump process above

### 5. Post-Release Validation
```bash
brew update && brew upgrade linctl
linctl --version
linctl docs | head -n 5
```

Run smoke test against Linear workspace if possible.

### 6. Housekeeping
- Close issues tied to the release
- Start new milestone if applicable

## Build Artifacts

### Single Platform Build
```bash
make build
# Output: ./linctl
```

### Cross-Platform Build
```bash
make build-all
# Outputs:
# dist/linctl-linux-amd64
# dist/linctl-darwin-amd64
# dist/linctl-darwin-arm64
# dist/linctl-windows-amd64.exe
```

### Release Build
```bash
make release
# Cleans, creates dist/, and builds all platforms
```

## Version Management

Version is injected at build time using Git tags:
```makefile
VERSION=$(shell git describe --tags --exact-match 2>/dev/null || git rev-parse --short HEAD)
LDFLAGS=-ldflags "-X github.com/dorkitude/linctl/cmd.version=$(VERSION)"
```

**Version Display**:
```bash
linctl --version
# Output: linctl version vX.Y.Z
```

## Installation Paths

### System Installation
- **Binary Location**: `/usr/local/bin/linctl`
- **Config File**: `~/.linctl.yaml`
- **Auth File**: `~/.linctl-auth.json`

### Homebrew Installation
Homebrew manages:
- Binary installation/updates
- Uninstallation via `brew uninstall linctl`
- Tap management via `brew tap dorkitude/linctl`

## CI/CD Pipeline

### GitHub Actions Workflows

**Bump Homebrew Tap** (`.github/workflows/bump-tap.yml`):
- **Trigger**: Release published
- **Permissions**: `contents:read`, `pull-requests:write`
- **Steps**:
  1. Checkout repository
  2. Install tools (jq)
  3. Compute tarball SHA256
  4. Clone tap repository
  5. Update Formula/linctl.rb
  6. Push branch and open PR
- **Required Secret**: `HOMEBREW_TAP_TOKEN`

## Environment Requirements

### Build Environment
- **Go**: 1.23.0+ (toolchain 1.24.5)
- **Git**: For version injection
- **Make**: For build automation

### Runtime Environment
- **OS**: macOS, Linux, Windows
- **Architecture**: amd64, arm64
- **Network**: HTTPS access to Linear API

## Configuration Management

### Default Configuration
Created on first run at `~/.linctl.yaml`:
```yaml
output: table
limit: 50
api:
  timeout: 30s
  retries: 3
```

### Authentication Storage
Stored securely in `~/.linctl-auth.json`:
```json
{
  "api_key": "lin_api_..."
}
```

## Troubleshooting Deployment

### Homebrew Issues
```bash
# Update tap
brew update

# Reinstall if broken
brew uninstall linctl
brew install linctl

# Check formula
brew info linctl
```

### Build Issues
```bash
# Clean rebuild
make clean
make deps
make build

# Verify version
./linctl --version
```

### Release Issues
- **Tag mismatch**: Ensure tag follows `vX.Y.Z` format
- **Action failure**: Check `HOMEBREW_TAP_TOKEN` secret
- **Formula error**: Verify tarball URL and SHA256

## Security Considerations

- **API Keys**: Never commit to repository
- **Secrets**: Store in GitHub Secrets for CI/CD
- **Token Permissions**: Use fine-grained PATs with minimal scope
- **Distribution**: Homebrew validates checksums automatically

## Monitoring & Maintenance

### Release Metrics
- GitHub release download counts
- Homebrew install analytics (if enabled)
- Issue tracker for user feedback

### Maintenance Tasks
- Monitor Linear API changes
- Update dependencies regularly
- Review and respond to issues/PRs
- Keep documentation synchronized

## Resources

- [GitHub Repository](https://github.com/dorkitude/linctl)
- [Homebrew Tap](https://github.com/dorkitude/homebrew-linctl)
- [Release Documentation](https://github.com/dorkitude/linctl/releases)
- [Contributing Guide](./contribution-guide.md)
