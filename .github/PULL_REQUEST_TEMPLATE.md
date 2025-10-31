## Description

<!-- Provide a brief description of the changes in this PR -->

## Related Issue

<!-- Link to the issue this PR addresses -->
Fixes #(issue number)
<!-- Or use: Relates to #(issue number) -->

## Type of Change

<!-- Mark the relevant option with an 'x' -->

- [ ] Bug fix (non-breaking change that fixes an issue)
- [ ] New feature (non-breaking change that adds functionality)
- [ ] Breaking change (fix or feature that would cause existing functionality to not work as expected)
- [ ] Documentation update
- [ ] Code refactoring (no functional changes)
- [ ] Performance improvement
- [ ] CI/CD changes
- [ ] Dependencies update

## Changes Made

<!-- Describe the changes in detail -->

-
-
-

## Testing

<!-- Describe how you tested your changes -->

- [ ] WASM module builds successfully (`make build`)
- [ ] WASM module verified (`make verify`)
- [ ] Manual testing with Istio Gateway
- [ ] Integration tests added/updated (in parent helm chart if applicable)
- [ ] Tested with jwt-auth-service authorization endpoint

### Test Environment

<!-- Describe your test environment -->

- Istio Version:
- Authorization Service:

### Test Scenarios

<!-- List test scenarios performed -->

- [ ] Valid token passes through
- [ ] Revoked token blocked (403)
- [ ] Invalid token blocked (401)
- [ ] Authorization service timeout handled
- [ ] Authorization service error handled

## WASM Module

- [ ] Module builds without errors
- [ ] Module size reasonable (~3MB)
- [ ] No memory leaks observed
- [ ] Logs are clear and helpful

## Breaking Changes

<!-- If this introduces breaking changes, list them here -->

- [ ] No breaking changes
- [ ] Breaking changes (describe below)

<!-- Describe breaking changes:

-

-->

## Checklist

- [ ] My code follows the project's code style
- [ ] I have performed a self-review of my code
- [ ] I have commented my code, particularly in hard-to-understand areas
- [ ] I have made corresponding changes to the documentation
- [ ] My changes generate no new warnings or errors
- [ ] WASM module builds successfully
- [ ] I have updated CHANGELOG.md under the `[Unreleased]` section
- [ ] My commits follow the [Conventional Commits](https://www.conventionalcommits.org/) specification
- [ ] I have rebased my branch on the latest `main`

## Configuration Changes

<!-- If this PR changes WasmPlugin configuration, document here -->

- [ ] No configuration changes
- [ ] Configuration changes (describe below)

<!-- Describe configuration changes:

Before:
```yaml
pluginConfig:
  ...
```

After:
```yaml
pluginConfig:
  ...
```

-->

## Additional Notes

<!-- Add any additional notes, screenshots, or context here -->

## For Reviewers

<!-- Highlight specific areas you'd like reviewers to focus on -->

**Focus areas:**
-
-

**Questions for reviewers:**
-
-

**Deployment notes:**
<!-- Any special considerations for deploying this change -->
-
