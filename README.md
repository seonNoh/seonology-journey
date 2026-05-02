# Seonology Journey

> 사용자 2명을 위한 여행 계획 · 기록 · 추억 보관 서비스 (모노레포)

[![back-ci](https://github.com/seonNoh/seonology-journey/actions/workflows/back-ci.yml/badge.svg)](https://github.com/seonNoh/seonology-journey/actions/workflows/back-ci.yml)
[![api-ci](https://github.com/seonNoh/seonology-journey/actions/workflows/api-ci.yml/badge.svg)](https://github.com/seonNoh/seonology-journey/actions/workflows/api-ci.yml)
[![web-ci](https://github.com/seonNoh/seonology-journey/actions/workflows/web-ci.yml/badge.svg)](https://github.com/seonNoh/seonology-journey/actions/workflows/web-ci.yml)
[![android-ci](https://github.com/seonNoh/seonology-journey/actions/workflows/android-ci.yml/badge.svg)](https://github.com/seonNoh/seonology-journey/actions/workflows/android-ci.yml)

## 구성

| 경로            | 기술                         | 역할           | 도메인                      |
| --------------- | ---------------------------- | -------------- | --------------------------- |
| `apps/back/`    | Go 1.23 + gRPC               | 비즈니스 로직  | (internal) `:50051`         |
| `apps/api/`     | Go 1.23 + Gin + WebSocket    | REST API + WS  | `journey-api.seonology.com` |
| `apps/web/`     | React + Vite + TS + Tailwind | SPA            | `journey.seonology.com`     |
| `apps/android/` | Kotlin + Jetpack Compose     | Native Android | (sideload)                  |
| `proto/`        | Protobuf + buf               | API 스키마     | —                           |
| `deploy/`       | Kustomize                    | k8s 매니페스트 | —                           |
| `scripts/`      | bash / Go                    | AWS 셋업       | —                           |

## 사전 요구

- Node 20.18+, pnpm 9 (corepack)
- Go 1.23+
- JDK 21 (Android)
- AWS CLI v2 (`seonology` profile)
- `gh` CLI 로그인
- `kubectl` (`k3s-lightsail` context)

## 설치

```bash
pnpm install
```

## 빌드 / 테스트

```bash
# back
cd apps/back && go test ./... && go build ./...

# api
cd apps/api && go test ./... && go build ./...

# web
cd apps/web && pnpm install && pnpm dev

# proto
make proto
```

## 인프라

배포 환경은 `seonology-k3s` 워크스페이스의 k3s 클러스터를 활용한다. ArgoCD 가 본 리포의 `deploy/overlays/production/` 을 자동 동기화한다.

## 커밋

Conventional Commits, scope 필수 (`feat(api):`, `fix(web):`, `chore(deps):`, …).

## 라이선스

MIT
