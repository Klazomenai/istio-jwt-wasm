# Istio JWT WASM Plugin

[![CI](https://github.com/klazomenai/istio-jwt-wasm/actions/workflows/ci.yaml/badge.svg)](https://github.com/klazomenai/istio-jwt-wasm/actions/workflows/ci.yaml)
[![Go Version](https://img.shields.io/github/go-mod/go-version/klazomenai/istio-jwt-wasm)](https://github.com/klazomenai/istio-jwt-wasm)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

Envoy WASM plugin for real-time JWT revocation enforcement at Istio Gateway.

> **⚠️ Alpha Status**: This project is in alpha (v0.0.1-alpha). Not recommended for production use. See [Known Limitations](#known-limitations-alpha) below.

## Table of Contents

- [Features](#features)
- [Quick Start](#quick-start)
- [Installation](#installation)
- [Configuration](#configuration)
- [How It Works](#how-it-works)
- [Building](#building)
- [Integration](#integration)
- [Known Limitations](#known-limitations-alpha)
- [Roadmap](#roadmap)
- [Contributing](#contributing)
- [Security](#security)
- [License](#license)

## Features

- **JWT Revocation Checking**: Real-time revocation enforcement via HTTP callout to authorization service
- **JTI Extraction**: Extracts JWT ID (JTI) from token for revocation lookup
- **Async HTTP Dispatch**: Minimizes latency with non-blocking authorization checks
- **Fail-Closed Security**: Blocks requests on authorization service failure (secure by default)
- **Configurable**: Simple configuration via Istio WasmPlugin CRD
- **Istio Gateway Integration**: Deploys at gateway level for centralized enforcement
- **Environment Agnostic**: No hardcoded endpoints or network assumptions

## Quick Start

### Deploy to Istio Gateway

```yaml
apiVersion: extensions.istio.io/v1alpha1
kind: WasmPlugin
metadata:
  name: jwt-revocation
  namespace: istio-system
spec:
  selector:
    matchLabels:
      istio: gateway
  url: oci://ghcr.io/klazomenai/istio-jwt-wasm:v0.0.1-alpha
  phase: AUTHN
  pluginConfig:
    cluster: "outbound|8080||jwt-auth-service.default.svc.cluster.local"
    authorizeUrl: "/authorize"
    timeout: 2000
```

## Installation

### Prerequisites

- Istio 1.18+ with Gateway API support
- WASM support enabled in Istio
- Authorization service (e.g., jwt-auth-service) deployed

### From OCI Registry (Recommended)

Use the WasmPlugin CRD as shown in Quick Start. Istio will automatically pull the WASM module from GHCR.

### Build from Source

```bash
# Clone repository
git clone https://github.com/klazomenai/istio-jwt-wasm.git
cd istio-jwt-wasm

# Install dependencies
make deps

# Build WASM module
make build

# Output: autonity-jwt-revocation.wasm (~3.1MB)
```

### Build OCI Image

```bash
# Build OCI image
make docker-build IMAGE_TAG=v0.0.1-alpha

# Push to registry
make docker-push IMAGE_TAG=v0.0.1-alpha DOCKER_REGISTRY=ghcr.io/your-org
```

## Configuration

### WasmPlugin CRD

The plugin is configured via the `pluginConfig` section of the Istio WasmPlugin CRD.

#### Required Fields

| Field | Description | Example |
|-------|-------------|---------|
| `cluster` | Envoy cluster name for authorization service | `outbound\|8080\|\|jwt-auth-service.default.svc.cluster.local` |

#### Optional Fields

| Field | Default | Description |
|-------|---------|-------------|
| `authorizeUrl` | `/authorize` | Authorization endpoint path |
| `timeout` | `2000` | HTTP call timeout in milliseconds |

### Cluster Name Format

The `cluster` field uses Envoy's cluster naming convention:

```
outbound|<PORT>||<HOSTNAME>
```

Examples:
```
outbound|8080||jwt-auth-service.default.svc.cluster.local
outbound|8080||auth-service.api-namespace.svc.cluster.local
```

The plugin automatically extracts the hostname for the `:authority` HTTP/2 header.

### Example Configurations

#### Basic Configuration

```yaml
pluginConfig:
  cluster: "outbound|8080||jwt-auth-service.default.svc.cluster.local"
```

#### Custom Authorization Endpoint

```yaml
pluginConfig:
  cluster: "outbound|8080||auth.example.svc.cluster.local"
  authorizeUrl: "/api/v1/validate"
  timeout: 5000
```

#### Production Configuration

```yaml
pluginConfig:
  cluster: "outbound|443||auth-service.production.svc.cluster.local"
  authorizeUrl: "/authorize"
  timeout: 3000
```

## How It Works

### Request Flow

```
1. Client sends request with JWT → Istio Gateway
2. Istio RequestAuthentication validates JWT signature
3. WASM plugin extracts JWT from Authorization header
4. WASM plugin extracts JTI (JWT ID) from token payload
5. WASM plugin makes async HTTP call to authorization service
   - POST {authorizeUrl}
   - Header: Authorization: Bearer <full-jwt-token>
6. Authorization service validates token and checks revocation
7. Authorization service returns status:
   - 200 OK → Token valid, request proceeds
   - 401 Unauthorized → Token invalid, request blocked
   - 403 Forbidden → Token revoked, request blocked
   - Other/timeout → Request blocked (fail-closed)
8. WASM plugin resumes or blocks request based on response
```

### Authorization Service Integration

This plugin expects an authorization service that:

1. Accepts `POST {authorizeUrl}` requests
2. Expects JWT in `Authorization: Bearer <token>` header
3. Returns appropriate HTTP status:
   - `200` if token is valid and not revoked
   - `401` if token is invalid (signature, expiry, etc.)
   - `403` if token is revoked

See [jwt-auth-service](https://github.com/klazomenai/jwt-auth-service) for a compatible implementation.

### Security Model

- **Fail-Closed**: On authorization service timeout or error, requests are blocked
- **Async Dispatch**: Uses Envoy's async HTTP dispatch to minimize latency
- **Token Forwarding**: Sends full JWT to authorization service for validation and revocation check
- **No Token Validation**: Relies on Istio RequestAuthentication for signature validation

## Building

### Build WASM Module

```bash
# Build with default settings
make build

# Verify WASM module
make verify  # Requires wasm-validate tool

# Clean build artifacts
make clean
```

### Build Requirements

- Go 1.24+ (with WASI support)
- GOOS=wasip1 GOARCH=wasm
- Build mode: c-shared

### Build Output

- File: `autonity-jwt-revocation.wasm`
- Size: ~3.1MB
- Format: WebAssembly binary
- Target: WASI preview 1

### Docker/OCI Build

```bash
# Build OCI image
make docker-build

# Configure registry and tag
make docker-build DOCKER_REGISTRY=ghcr.io/your-org IMAGE_TAG=v1.0.0

# Push to registry
make docker-push
```

## Integration

### With Istio Gateway

1. Deploy authorization service (jwt-auth-service)
2. Create Service for authorization service
3. Deploy WasmPlugin CRD (see Quick Start)
4. Configure RequestAuthentication for JWT validation
5. Test with valid and revoked tokens

### With jwt-auth-service

This plugin is designed to integrate with [jwt-auth-service](https://github.com/klazomenai/jwt-auth-service). Deploy the authorization service and configure the WasmPlugin CRD to reference it using the appropriate cluster name (see Configuration section).

## Known Limitations (Alpha)

This alpha release has the following limitations:

### Functionality

- **No Circuit Breaker**: No circuit breaker pattern for authorization service failures
- **No Retry Logic**: Failed authorization calls are not retried
- **No Fallback**: Single authorization endpoint (no redundancy)
- **Fail-Closed Only**: No fail-open option (blocks all traffic on auth service failure)

### Observability

- **Limited Error Context**: Generic error messages to clients
- **Detailed Errors in Logs Only**: Full error details only in Envoy logs
- **No Metrics Export**: No Prometheus metrics (relies on Envoy metrics)

### Testing

- **No Unit Tests**: Integration tests only (in parent helm chart)
- **Manual Testing Required**: No automated integration tests in this repo

Do not use in production until beta status. See [SECURITY.md](SECURITY.md) for deployment considerations.

## Roadmap

### Beta (v0.1.0) - Planned

- Circuit breaker pattern for authorization service failures
- Retry logic with exponential backoff
- Fallback authorization endpoints
- Improved error messages to clients
- Prometheus metrics export
- Fail-open configuration option (with warnings)

### Future Enhancements

- Web3 oracle integration for blockchain-backed revocation
- Smart contract integration for decentralized revocation lists
- WASM-to-WASM communication with co-located oracle
- Adaptive timeout and retry logic
- Response caching for improved performance
- Rate limiting integration

See [CHANGELOG.md](CHANGELOG.md) for release history.

## Contributing

Contributions are welcome! This project is in active development.

Before contributing:

1. Read [CONTRIBUTING.md](CONTRIBUTING.md) for development guidelines
2. Read [CODE_OF_CONDUCT.md](CODE_OF_CONDUCT.md) for community standards
3. Check existing issues and PRs to avoid duplicates
4. Follow [Conventional Commits](https://www.conventionalcommits.org/) for commit messages

### How to Contribute

- Report bugs using the issue tracker
- Suggest features via feature request template
- Submit pull requests with integration tests
- Improve documentation
- Help answer questions in discussions

## Security

Security vulnerabilities should be reported privately. See [SECURITY.md](SECURITY.md) for:

- Supported versions
- How to report vulnerabilities
- Response timeline
- Security best practices
- Known security limitations

Do NOT report security issues through public GitHub issues.

## Requirements

- Istio 1.18+ with Gateway API support
- Istio with WASM support enabled
- Authorization service compatible with expected API

## License

This project is licensed under the MIT License - see [LICENSE](LICENSE) file for details.

## Acknowledgments

This WASM plugin was extracted from the [helm-priv-deploy-autonity](https://github.com/clearmatics/helm-priv-deploy-autonity) project to provide a standalone, reusable JWT revocation enforcement solution.

Designed for use with Istio service mesh and Envoy proxy.

## Links

- [GitHub Repository](https://github.com/klazomenai/istio-jwt-wasm)
- [Issue Tracker](https://github.com/klazomenai/istio-jwt-wasm/issues)
- [Releases](https://github.com/klazomenai/istio-jwt-wasm/releases)
- [Contributing Guidelines](CONTRIBUTING.md)
- [Security Policy](SECURITY.md)
- [Code of Conduct](CODE_OF_CONDUCT.md)
- [jwt-auth-service](https://github.com/klazomenai/jwt-auth-service) - Compatible authorization service
