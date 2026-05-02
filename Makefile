.PHONY: help proto proto-lint proto-breaking proto-format clean-proto \
        back-build back-test api-build api-test web-install web-build web-dev \
        deploy-render

help:
	@echo "Targets:"
	@echo "  proto            - generate Go/TS code from proto via buf"
	@echo "  proto-lint       - lint .proto files"
	@echo "  proto-breaking   - check breaking changes vs main"
	@echo "  proto-format     - format .proto files"
	@echo "  back-test        - go test apps/back"
	@echo "  api-test         - go test apps/api"
	@echo "  web-build        - pnpm build apps/web"

proto:
	cd proto && buf generate

proto-lint:
	cd proto && buf lint

proto-format:
	cd proto && buf format -w

proto-breaking:
	cd proto && buf breaking --against '.git#branch=main,subdir=proto'

clean-proto:
	rm -rf proto/gen

back-build:
	cd apps/back && go build ./...
back-test:
	cd apps/back && go test -race ./...

api-build:
	cd apps/api && go build ./...
api-test:
	cd apps/api && go test -race ./...

web-install:
	pnpm install
web-build:
	pnpm --filter @seonology/journey-web build
web-dev:
	pnpm --filter @seonology/journey-web dev

deploy-render:
	kubectl kustomize deploy/overlays/production
