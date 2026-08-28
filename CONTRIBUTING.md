# Contributing

## Branches and commits

- Start changes from `main` and use a short-lived branch.
- Use Conventional Commits with an explicit component scope, such as `fix(api):`, `feat(web):`, or `chore(back):`.
- Keep generated Protobuf files in sync with their source definitions.

## Required verification

Run the repository contract before opening a pull request:

```bash
python3 -m unittest tests/test_repository_contract.py
python3 verify.py
```

Run the affected component tests as well. Pull requests that change runtime images must retain `linux/amd64` and `linux/arm64` support.

## Security

Do not commit tokens, credentials, private keys, Android signing material, Firebase service accounts, or production environment values. Use Gitea secrets, Vault, and Kubernetes Secrets.

## Review expectations

Explain the reason for the change, affected components, verification commands, and operational impact. Do not add automated-author signatures or emoji to commits and project artifacts.
