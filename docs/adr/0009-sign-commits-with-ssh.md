# 9. Sign commits with SSH

Date: 2026-10-16

## Status

Accepted

## Context

Signed commits let GitHub show that a commit was made by the account owner.
Git 2.34 and later can sign with SSH keys as well as GPG keys. The
maintainer's SSH key is stored in 1Password, which also acts as the SSH
agent, so the private key never leaves 1Password.

Options considered:

- **GPG**: requires a separate key and toolchain.
- **SSH**: reuses an existing key type and the existing 1Password agent.

jj rewrites commits frequently (on edit, rebase, and describe), and
1Password asks for authorization on each signing request.

This is the same decision as Python ADR-0009, and uses the same jj and Git
configuration, which is set per user rather than per repository.

## Decision

Commits are signed with an SSH key held in 1Password.

- **jj** uses the SSH signing backend with `signing.behavior = "drop"` and
  `git.sign-on-push = true`. Commits are signed when they are pushed rather
  than each time they are rewritten.
- **Git** is configured for SSH signing through 1Password's "Configure
  Commit Signing" option, for any commit made with `git` directly.
- The public key is registered on GitHub as a **Signing key**, separate from
  its registration as an authentication key.
- A local allowed-signers file maps the commit email to the public key, so jj
  can verify signatures locally.

## Consequences

- Pushed commits show as **Verified** on GitHub.
- Local commits are unsigned until pushed. `jj log` shows a signature status
  only for pushed commits.
- Signatures verify only when the commit author email matches the email on
  the GitHub account and in the allowed-signers file.
