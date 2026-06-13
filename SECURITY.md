# Security Policy

## Supported versions

Security fixes are applied on the default branch (`main`). There are no long-term release branches yet.

## Reporting a vulnerability

**Please do not report security vulnerabilities in public GitHub issues.**

Email the maintainers with:

- A description of the issue
- Steps to reproduce
- Impact assessment (if known)

Use GitHub’s [private vulnerability reporting](https://github.com/remdev/cursor-go-sdk/security/advisories/new) if enabled for this repository, or contact the repository owner via GitHub.

We aim to acknowledge reports within a few business days.

## Scope

In scope:

- This Go module (`cursor/`, `internal/`)
- The vendored bridge launcher and Connect glue in `bridge/` (excluding third-party code in `node_modules/`)

Out of scope:

- Vulnerabilities in `@cursor/sdk` or other npm dependencies — report upstream or to Cursor as appropriate
- Issues that require a valid `CURSOR_API_KEY` to abuse Cursor’s cloud API (report to Cursor)
- Social engineering or phishing using the project name

## Secrets

Never commit API keys, `.env` files, or tokens. CI does not use production credentials.
