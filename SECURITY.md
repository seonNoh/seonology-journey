# Security Policy

이 프로젝트는 개인 사용 목적으로 운영된다. 보안 이슈는 GitHub Issues 가 아닌 직접 연락(이메일 / Mattermost DM)으로 알린다.

## 시크릿 관리

- 모든 시크릿은 Vault (`vault.seonology.com`) 에 저장한다.
- Git 리포에 평문 시크릿을 절대 커밋하지 않는다. `gitleaks` pre-commit 훅이 차단한다.
- AWS / GCP 키는 Vault Agent Sidecar 로 Pod 에 파일 마운트한다.
