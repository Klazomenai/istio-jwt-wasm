# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Known Issues
- No circuit breaker pattern for authorization service failures
- No fallback or retry logic for failed authorization calls
- Generic error messages to clients (detailed errors only in Envoy logs)

## [0.0.1-alpha] - 2025-10-31

### Added
- Initial alpha release extracted from helm-priv-deploy-autonity
- Envoy WASM plugin for JWT revocation enforcement at Istio Gateway
- JWT token extraction from Authorization header
- JTI (JWT ID) extraction from token payload
- Async HTTP dispatch to authorization service
- Configurable via Istio WasmPlugin CRD
  - `cluster` - Envoy cluster name (required)
  - `authorizeUrl` - Authorization endpoint path (default: `/authorize`)
  - `timeout` - HTTP call timeout in milliseconds (default: 2000)
- Fail-closed security model
  - Blocks requests on authorization service timeout
  - Blocks requests on authorization service error
  - Blocks requests on unexpected status codes
- Response handling
  - `200 OK` - Token valid, request proceeds
  - `401 Unauthorized` - Token invalid, request blocked
  - `403 Forbidden` - Token revoked, request blocked
  - Other/timeout - Request blocked (fail-closed)
- Authority extraction from Envoy cluster name
  - Parses `outbound|PORT||HOSTNAME` format
  - Sets `:authority` header for HTTP/2 requests
- Comprehensive logging
  - JTI extraction events
  - Authorization requests/responses
  - Error conditions
  - Debug information for troubleshooting
- Environment-agnostic configuration
  - No hardcoded cluster names
  - No hardcoded authorization URLs
  - Fail-fast validation of required config
- Docker OCI image support
  - Multi-stage Dockerfile (Go 1.24 with WASI)
  - GHCR registry: ghcr.io/klazomenai/istio-jwt-wasm
  - Scratch-based final image (minimal size)
- GitHub Actions CI/CD
  - Automated WASM building on push/PR
  - OCI image building on version tags
  - Artifact upload for verification
- Documentation
  - README with WasmPlugin configuration examples
  - Makefile with build, verify, and docker targets
  - MIT License

### Security
- Fail-closed security model (blocks on errors)
- No hardcoded credentials or endpoints
- Authorization header forwarded to authorization service
- Envoy cluster validation
- Timeout enforcement on authorization calls

### Known Limitations (Alpha)
- No circuit breaker pattern for authorization service failures
- No retry logic for failed authorization calls
- No fallback authorization endpoints
- Limited error context in client responses
- Fail-closed only (no fail-open option)
- Single authorization endpoint (no redundancy)
- Integration tests in parent chart only (no unit tests in this repo)

### Technical Details
- WASM module size: ~3.1MB
- Go version: 1.24 (with WASI support)
- Target: wasip1 (WebAssembly System Interface preview 1)
- Build mode: c-shared (shared library for proxy-wasm)
- SDK: proxy-wasm-go-sdk v0.0.0-20250212164326

### Notes
This is an alpha release intended for testing and development. Do not use in production.

For production readiness requirements, see:
- [SECURITY.md](SECURITY.md) for security considerations
- [CONTRIBUTING.md](CONTRIBUTING.md) for development guidelines
- Integration with jwt-auth-service for complete JWT revocation architecture

### Integration
This plugin is designed to work with jwt-auth-service:
- jwt-auth-service provides `/authorize` endpoint
- This plugin calls `/authorize` with JWT token
- jwt-auth-service validates token and checks revocation status
- Response determines if request proceeds or is blocked

See [jwt-auth-service](https://github.com/klazomenai/jwt-auth-service) for compatible authorization service implementation.

[Unreleased]: https://github.com/klazomenai/istio-jwt-wasm/compare/v0.0.1-alpha...HEAD
[0.0.1-alpha]: https://github.com/klazomenai/istio-jwt-wasm/releases/tag/v0.0.1-alpha
