# seonology-journey-back

DynamoDB / S3 への永続化と外部 API 連携を担当する gRPC サービス.

- gRPC port: `9090` (health, reflection 有効)
- 主要依存: AWS SDK Go v2, zerolog, prometheus client_golang
- スキーマ: `proto/journey/v1/*` (`JourneyService`)
- 詳細仕様: `docs/1_architecture_design.md` §3, §5
