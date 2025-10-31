# Security Policy

## Supported Versions

Currently supported versions for security updates:

| Version | Supported          | Status |
| ------- | ------------------ | ------ |
| 0.0.x   | :white_check_mark: | Alpha  |

As this project is in alpha, we recommend using only the latest release.

## Reporting a Vulnerability

We take security vulnerabilities seriously. If you discover a security issue, please follow responsible disclosure practices.

### How to Report

**DO NOT** report security vulnerabilities through public GitHub issues.

Instead, please report vulnerabilities using one of these methods:

1. **GitHub Security Advisories** (Preferred)
   - Navigate to the repository's Security tab
   - Click "Report a vulnerability"
   - Fill out the form with vulnerability details
   - This allows private discussion and coordinated disclosure

2. **Email** (Alternative)
   - Create a private issue with the label "security"
   - Include "SECURITY" in the issue title
   - Mark the issue as private if your GitHub plan supports it

### What to Include

When reporting a vulnerability, please include:

- Description of the vulnerability
- Steps to reproduce the issue
- Potential impact and attack scenarios
- Suggested fix (if available)
- Your contact information (for follow-up)
- Any proof-of-concept code (please do not exploit beyond proof-of-concept)

### Response Timeline

We aim to respond to security reports according to this timeline:

- **Initial Response**: Within 48 hours
- **Validation**: Within 7 days (confirm if it's a valid security issue)
- **Fix Timeline**: Depends on severity
  - Critical: Within 7 days
  - High: Within 14 days
  - Medium: Within 30 days
  - Low: Next planned release
- **Disclosure**: Coordinated disclosure after fix is released

### Severity Classification

We use the following severity levels:

- **Critical**: Request forgery, authentication bypass, gateway compromise
- **High**: Authorization bypass, token validation errors
- **Medium**: Information disclosure, denial of service
- **Low**: Minor security improvements, hardening recommendations

## Security Best Practices

When deploying this WASM plugin:

### Istio Gateway Configuration

- Use secure cluster names (no wildcards)
- Specify exact authorization service endpoint
- Set appropriate timeout values (not too long)
- Enable fail-closed mode (block on errors)
- Use TLS for authorization service communication

### WasmPlugin Configuration

- Validate required configuration before deployment
- Use environment-specific cluster names
- Never hardcode authorization URLs
- Set reasonable timeout values (2-5 seconds recommended)
- Test with actual authorization service before production

### Monitoring

- Monitor authorization request latency
- Alert on high error rates from authorization service
- Track blocked requests (403 responses)
- Log authorization failures for security audits

## Known Limitations (Alpha)

The current alpha release has these known security limitations:

1. **Fail-Closed Only**
   - No fail-open option (blocks all traffic on authorization service failure)
   - Can cause service outage if authorization service unavailable
   - Planned: Circuit breaker pattern for beta

2. **No Request Validation**
   - Assumes JWT already validated by Istio RequestAuthentication
   - No additional JWT signature verification in WASM
   - Relies on authorization service for all validation

3. **Single Authorization Endpoint**
   - No fallback or retry to alternative endpoints
   - No circuit breaker for failed authorization services
   - Planned for beta

4. **Limited Error Context**
   - Generic error messages to client
   - Detailed errors only in Envoy logs
   - May make debugging difficult

These limitations are documented and will be addressed in future releases. Do not deploy to production until beta status.

## Security Updates

Security fixes will be released as patch versions (e.g., v0.0.2) with:
- Updated CHANGELOG.md with security notice
- GitHub Security Advisory published
- Release notes detailing the fix
- Credit to reporter (unless anonymous)
- New OCI image pushed to GHCR

Subscribe to repository releases to receive security notifications.

## Safe Harbor

We support security research and vulnerability disclosure. Researchers who:
- Follow responsible disclosure practices
- Do not exploit beyond proof-of-concept
- Allow reasonable time for fixes
- Do not harm users or services

Will be thanked and credited (if desired) in:
- CHANGELOG.md security notices
- GitHub Security Advisories
- Release notes

## Contact

For security-related questions that are not vulnerabilities, you may:
- Open a public issue with the "security" label
- Ask in GitHub Discussions
- Reference this policy

## Dependencies

We monitor dependencies for known vulnerabilities using:
- Go's vulnerability database (govulncheck)
- GitHub Dependabot alerts
- Regular dependency updates

Current critical dependency:
- `proxy-wasm-go-sdk` - Envoy WASM SDK (vendor: The Proxy-Wasm Authors)

See `go.mod` for complete dependency list.

## Deployment Security

### Authorization Service Integration

The WASM plugin forwards JWT tokens to an authorization service. Ensure:

- Authorization service uses TLS (in production)
- Service is properly authenticated (mutual TLS recommended)
- Network policies restrict access to authorization service
- Authorization service has proper rate limiting
- Authorization service validates JWT signatures
- Authorization service checks token revocation

### Istio Gateway Security

- Use Istio RequestAuthentication to validate JWT signatures first
- WASM plugin should run AFTER authentication (phase: AUTHN)
- Enable AuthorizationPolicy for additional access controls
- Monitor gateway logs for security events

### OCI Image Security

- Pull images from trusted registries only (ghcr.io/klazomenai)
- Use specific version tags (not :latest in production)
- Scan images for vulnerabilities before deployment
- Use image pull secrets for private registries
- Enable Pod Security Policies

## Compliance

This WASM plugin processes JWT tokens which may contain personally identifiable information (PII). When deploying:

- Ensure compliance with relevant data protection regulations (GDPR, CCPA, etc.)
- Review what data is sent to authorization service
- Implement appropriate data retention policies
- Enable audit logging for compliance requirements
- Ensure authorization service handles PII securely

## Istio/Envoy Security

This plugin runs in the Istio Gateway Envoy proxy. Security considerations:

- WASM modules run in a sandboxed environment
- No direct file system access
- No direct network access (except via Envoy APIs)
- Limited memory allocation
- No access to host system

However:
- WASM can access all request/response data
- WASM can make HTTP calls via Envoy dispatch
- WASM errors can crash Envoy worker threads
- Malicious WASM can cause denial of service

Only deploy WASM modules from trusted sources.

## Attribution

Security researchers who have contributed to this project:

(None yet - be the first!)

## Updates to This Policy

This security policy may be updated as the project matures. Check the commit history for changes.

Last updated: 2025-10-31
