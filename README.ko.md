# Seonology Journey

[English](README.md) | [한국어](README.ko.md) | [日本語](README.ja.md)

Seonology Journey는 개인 여행 계획과 추억 보관을 위한 비공개 모노레포입니다. React 웹 클라이언트, Android 클라이언트, REST 및 WebSocket API, gRPC back 서비스, 공통 Protobuf 계약과 k3s 매니페스트를 함께 관리합니다.

## 아키텍처

![서비스 아키텍처](docs/svg/service-architecture.ko.svg)

웹과 Android 클라이언트는 API 서비스를 사용합니다. API는 영속 데이터와 미디어 처리를 gRPC를 통해 back 서비스에 위임합니다. 운영 환경의 데이터는 AWS DynamoDB와 S3에 저장합니다.

## 컴포넌트

![컴포넌트 배포](docs/svg/component-delivery.ko.svg)

| 경로           | 역할                         | 런타임 산출물              |
| -------------- | ---------------------------- | -------------------------- |
| `apps/web`     | React 및 Vite 웹 클라이언트  | `seonology-journey-web`    |
| `apps/api`     | REST 및 WebSocket 경계 API   | `seonology-journey-api`    |
| `apps/back`    | gRPC 데이터 및 미디어 서비스 | `seonology-journey-back`   |
| `apps/android` | 네이티브 Android 클라이언트  | APK 또는 AAB               |
| `proto`        | 공통 Protobuf 계약           | Go 및 TypeScript 생성 코드 |
| `deploy`       | 독립 Kustomize 매니페스트    | Kubernetes 리소스          |

## 사전 요구 사항

- Node.js 22 및 pnpm 9.15.9
- Go 1.25 이상
- Android 빌드용 JDK 21 및 Android SDK
- 다중 아키텍처 이미지 빌드용 Docker와 Buildx
- 매니페스트 검증용 `kubectl`과 Kustomize

## 빠른 시작

```bash
corepack enable
corepack prepare pnpm@9.15.9 --activate
pnpm install --frozen-lockfile
pnpm --filter @seonology/journey-web typecheck
pnpm --filter @seonology/journey-web test
pnpm --filter @seonology/journey-web build

cd apps/api && go test ./... && go build ./...
cd ../back && go test ./... && go build ./...
```

## 검증

```bash
python3 -m unittest tests/test_repository_contract.py
python3 verify.py
```

검증기는 다국어 문서, 결정적인 Relief 다이어그램, Gitea workflow 계약, 보존된 GitHub workflow 바이트와 컴포넌트 빌드를 검사합니다.

## 배포

![운영 전환](docs/svg/runtime-cutover.ko.svg)

Gitea가 기준 저장소입니다. Gitea Actions는 저장소를 검증한 뒤 API, back, web 이미지를 직렬로 게시하며, 각 이미지는 `linux/amd64`와 `linux/arm64` OCI index를 포함합니다. GitHub는 Actions를 비활성화한 push mirror로 유지합니다.

## 운영

운영 workload는 중앙 `seonology-k3s` GitOps 저장소의 `workloads/seonology-journey`에서 관리합니다. Argo CD가 Synced 및 Healthy이고, Deployment 3개가 모두 Ready이며, imageID가 검증된 Gitea Registry digest를 가리키고, 내부 및 공개 health 검사가 성공해야 배포 완료로 판정합니다.

## 저장소 정책

![저장소 역할](docs/svg/repository-roles.ko.svg)

- 명시적인 scope가 있는 Conventional Commits를 사용합니다.
- 이전 증거인 `.github/workflows`의 바이트를 보존합니다.
- 자격증명은 Gitea secret, Vault 또는 Kubernetes Secret에만 저장합니다.
- 프로젝트 산출물에 자동 작성자 서명이나 이모지를 추가하지 않습니다.

## 라이선스

기존 MIT 라이선스를 유지합니다. 자세한 내용은 [LICENSE](LICENSE)를 참조하십시오.
