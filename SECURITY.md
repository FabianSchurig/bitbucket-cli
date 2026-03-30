# Security Policy

## Supported Versions

| Version | Supported          |
|---------|--------------------|
| Latest  | :white_check_mark: |
| < Latest | :x:               |

Only the most recent release is supported. The daily CI pipeline automatically updates the CLI when the Bitbucket API spec changes, so staying on the latest version is strongly recommended.

## Reporting a Vulnerability

**Please do not open a public issue for security vulnerabilities.**

Instead, use one of the following:

1. **GitHub Security Advisories (preferred):** Open a [private security advisory](https://github.com/FabianSchurig/bitbucket-cli/security/advisories/new) on this repository.
2. **Email:** Contact the maintainer directly at the email listed on their GitHub profile.

### What to include

- A description of the vulnerability and its potential impact.
- Steps to reproduce the issue.
- Affected versions (if known).
- Any suggested fix or mitigation.

### Response timeline

- **Acknowledgement:** Within 5 business days.
- **Initial assessment:** Within 10 business days.
- **Fix or mitigation:** As soon as practical, typically within 30 days for confirmed vulnerabilities.

## Scope

This policy covers:

- The `bb-cli` binary and its source code.
- The code generation pipeline scripts (`scripts/`).
- GitHub Actions workflows in this repository.

Out of scope:

- The Bitbucket Cloud API itself — report those to [Atlassian](https://www.atlassian.com/trust/security/report-a-vulnerability).
- Third-party dependencies — report those to their respective maintainers (Dependabot monitors for known CVEs).
