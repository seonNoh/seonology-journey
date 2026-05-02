# seonology-journey-api

REST + WebSocket ゲートウェイ. Keycloak OIDC で認証, `seonology-journey-back` への gRPC ファサード.

- HTTP port: `8080`
- 主要依存: gin, coder/websocket, keycloak JWKS
- スキーマ: `proto/journey/v1/*`
- 詳細仕様: `docs/1_architecture_design.md` §4, §6
