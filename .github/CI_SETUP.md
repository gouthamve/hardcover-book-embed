# CI/CD Setup Guide

This repository uses GitHub Actions for continuous integration and deployment.

## Workflows

### 1. Test Workflow (`test.yml`)
Comprehensive testing workflow that runs on every PR and push to main.

**Features:**
- Multi-version Go testing (1.22, 1.23)
- Security scanning with gosec
- Code coverage reporting
- Linting with golangci-lint

**Required Permissions:**
- `contents: read`
- `security-events: write` (for CodeQL/SARIF upload)

### 2. Simple Test Workflow (`test-simple.yml`)
A simpler alternative that doesn't require special permissions.

**When to use:**
- If you're getting "Resource not accessible by integration" errors
- For forks where security permissions might not be available
- For quick testing without all the bells and whistles

### 3. CI Workflow (`ci.yml`)
Full CI/CD pipeline including cross-platform builds and Docker images.

**Features:**
- Cross-platform testing (Linux, macOS, Windows)
- Multi-architecture builds
- Docker image building and pushing to ghcr.io
- Automated releases on tags

## Fixing Common Issues

### "Resource not accessible by integration" Error

This error occurs when trying to upload security scan results without proper permissions.

**Solutions:**

1. **Enable security events** in your repository:
   - Go to Settings → Actions → General
   - Under "Workflow permissions", select "Read and write permissions"
   - Check "Allow GitHub Actions to create and approve pull requests"

2. **Use the simple workflow** instead:
   ```bash
   # Manually trigger the simple workflow
   gh workflow run test-simple.yml
   ```

3. **For forks**: The security scanning features may not work in forks. Use the `test-simple.yml` workflow instead.

### Setting Up for Your Repository

1. **Enable GitHub Actions**:
   - Go to Settings → Actions → General
   - Enable "Allow all actions and reusable workflows"

2. **Configure permissions**:
   - Under "Workflow permissions", choose appropriate settings
   - For full functionality, select "Read and write permissions"

3. **For Docker builds** (optional):
   - The CI workflow automatically builds and pushes to GitHub Container Registry
   - Images are available at `ghcr.io/[your-username]/hardcover-book-embed`

4. **For releases** (optional):
   - Create a tag: `git tag v1.0.0`
   - Push the tag: `git push origin v1.0.0`
   - The CI workflow will create a release with binaries

## Workflow Badges

Add these to your README:

```markdown
[![CI](https://github.com/[your-username]/hardcover-book-embed/actions/workflows/ci.yml/badge.svg)](https://github.com/[your-username]/hardcover-book-embed/actions/workflows/ci.yml)
[![Test](https://github.com/[your-username]/hardcover-book-embed/actions/workflows/test.yml/badge.svg)](https://github.com/[your-username]/hardcover-book-embed/actions/workflows/test.yml)
```

## Running Tests Locally

Before pushing, you can run the same tests locally:

```bash
# Run all tests
make test

# Run with race detection
go test -race ./...

# Run with coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Run linting
golangci-lint run

# Run security scan
gosec ./...
```