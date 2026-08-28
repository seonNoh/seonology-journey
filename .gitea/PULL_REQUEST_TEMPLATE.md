## Purpose

Describe the problem and why this change is required.

## Scope

- Affected components:
- Operational impact:
- Compatibility impact:

## Verification

List every command that was executed and its result.

```text
python3 -m unittest tests/test_repository_contract.py
python3 verify.py
```

## Security and delivery

- [ ] No credential or signing material is committed.
- [ ] GitHub workflow bytes are unchanged unless the change explicitly targets GitHub automation.
- [ ] Runtime images retain `linux/amd64` and `linux/arm64` support.
- [ ] Documentation matches the implemented behavior.
