# Contributing to Istio JWT WASM Plugin

Thank you for your interest in contributing to the Istio JWT WASM Plugin! This document provides guidelines and instructions for contributing.

## Code of Conduct

This project adheres to a Code of Conduct that all contributors are expected to follow. Please read [CODE_OF_CONDUCT.md](CODE_OF_CONDUCT.md) before contributing.

## Getting Started

### Prerequisites

- Go 1.24 or later (with WASI support for WASM)
- Make
- Docker (for OCI image building)
- Git

### Development Setup

1. Fork and clone the repository:
```bash
git clone https://github.com/klazomenai/istio-jwt-wasm.git
cd istio-jwt-wasm
```

2. Install dependencies:
```bash
make deps
```

3. Build the WASM module:
```bash
make build
```

This produces `autonity-jwt-revocation.wasm` (approx 3.1MB).

## Development Workflow

### Making Changes

1. Create a feature branch from `main`:
```bash
git checkout -b feat/your-feature-name
```

2. Make your changes following the code style guidelines below

3. Build and verify the WASM module:
```bash
make build
make verify  # Requires wasm-validate tool (optional)
```

4. Commit your changes using conventional commits (see below)

5. Push to your fork and create a pull request

### Building

```bash
# Build WASM module
make build

# Verify WASM module
make verify

# Build Docker OCI image
make docker-build

# Clean build artifacts
make clean
```

## Code Style

### Go Conventions

- Follow standard Go code style and idioms
- Run `go fmt` before committing
- Keep functions small and focused
- Write clear, descriptive variable names
- Add comments for exported functions and types

### WASM-Specific Guidelines

- Minimize allocations (WASM has limited memory)
- Avoid blocking operations
- Use Envoy's proxy-wasm SDK correctly
- Log important events for debugging
- Handle errors gracefully (WASM errors can be hard to debug)

### File Organization

- Main WASM entry point: `main.go`
- Helper packages: `pkg/`
- No test files in this repository (integration tests in parent helm chart)

## Commit Message Guidelines

This project uses [Conventional Commits](https://www.conventionalcommits.org/) for automated changelog generation and semantic versioning.

### Commit Message Format

```
<type>[optional scope]: <description>

[optional body]

[optional footer(s)]
```

### Types

- `feat:` - New feature
- `fix:` - Bug fix
- `docs:` - Documentation changes only
- `style:` - Code style changes (formatting, no logic change)
- `refactor:` - Code restructuring (no feature or bug fix)
- `test:` - Adding or updating tests (integration tests in parent chart)
- `chore:` - Maintenance tasks (dependencies, build, etc.)
- `ci:` - CI/CD pipeline changes
- `perf:` - Performance improvements
- `revert:` - Revert a previous commit

### Scope (Optional)

The scope should indicate the affected component:
- `jwt` - JWT extraction/parsing
- `authz` - Authorization check logic
- `envoy` - Envoy integration
- `deps` - Dependency updates

### Examples

```
feat(jwt): extract JTI from token for revocation check
fix(authz): handle timeout from authorization service correctly
docs: update WasmPlugin configuration examples
perf: optimize JWT parsing to reduce latency
chore(deps): update proxy-wasm-go-sdk to v0.0.0-20250212164326
```

## Pull Request Process

### Before Submitting

1. Ensure WASM builds successfully: `make build`
2. Verify WASM module: `make verify` (if wasm-validate available)
3. Update documentation if needed
4. Add entry to CHANGELOG.md under `[Unreleased]` section
5. Rebase on latest `main` if needed

### PR Description

Include in your PR description:
- Summary of changes
- Related issue number (if applicable): `Fixes #123` or `Relates to #456`
- Type of change (feature, bugfix, docs, etc.)
- Testing performed (integration tests in parent chart)
- Any breaking changes

### Review Process

1. At least one maintainer approval required
2. All CI checks must pass (WASM build succeeds)
3. No merge conflicts
4. Commit messages follow conventional commits format

## Reporting Bugs

Use the GitHub issue tracker to report bugs. Include:
- Go version (`go version`)
- Istio version
- Steps to reproduce
- Expected vs actual behavior
- Relevant logs from Envoy/Istio Gateway

Use the bug report template when creating issues.

## Suggesting Features

Feature suggestions are welcome! Use the feature request template and include:
- Use case and motivation
- Proposed solution or implementation approach
- Alternatives considered
- Willingness to implement (if applicable)

## Security Vulnerabilities

Do NOT report security vulnerabilities through public GitHub issues. Please follow the process outlined in [SECURITY.md](SECURITY.md).

## Alpha Status Notice

This project is currently in alpha (v0.0.1-alpha). Contributions are welcome, but expect:
- API changes without notice
- Incomplete features (see known limitations in README)
- Breaking changes between releases
- Focus on core functionality over polish

## Testing

This repository does not contain unit tests. Integration tests are located in the parent helm chart:
- `helm-priv-deploy-autonity/stable/autonity/tests/`

When contributing, consider:
- Adding integration test cases to parent chart
- Documenting test scenarios in PR description
- Testing with actual Istio Gateway deployment

## Questions?

- Check existing issues and discussions
- Create a new issue with the question label
- Reference relevant documentation

## License

By contributing, you agree that your contributions will be licensed under the MIT License.
