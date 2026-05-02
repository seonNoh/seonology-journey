# seonology-journey-web

React 18 + Vite + TypeScript + Tailwind の SPA. Keycloak OIDC PKCE 認証.

- dev: `pnpm --filter @seonology/journey-web dev`
- build: `pnpm --filter @seonology/journey-web build`
- 配信: nginx `:8080` (Traefik ingress 経由で `https://journey.seonology.com`)
- 詳細: `docs/2_ui_ux_design.md`
