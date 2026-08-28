# Seonology Journey

[English](README.md) | [한국어](README.ko.md) | [日本語](README.ja.md)

Seonology Journeyは、個人の旅行計画と記録を管理するプライベートモノレポです。React Webクライアント、Androidクライアント、RESTおよびWebSocket API、gRPC backサービス、共通Protobuf契約、k3sマニフェストを一元管理します。

## アーキテクチャ

![サービス構成](docs/svg/service-architecture.ja.svg)

WebとAndroidクライアントはAPIサービスを利用します。APIは永続データとメディア処理をgRPC経由でbackサービスへ委譲します。本番データはAWS DynamoDBとS3に保存します。

## コンポーネント

![コンポーネント配布](docs/svg/component-delivery.ja.svg)

| パス           | 役割                             | 実行成果物                   |
| -------------- | -------------------------------- | ---------------------------- |
| `apps/web`     | ReactおよびVite Webクライアント  | `seonology-journey-web`      |
| `apps/api`     | RESTおよびWebSocket境界API       | `seonology-journey-api`      |
| `apps/back`    | gRPCデータおよびメディアサービス | `seonology-journey-back`     |
| `apps/android` | ネイティブAndroidクライアント    | APKまたはAAB                 |
| `proto`        | 共通Protobuf契約                 | GoおよびTypeScript生成コード |
| `deploy`       | 独立Kustomizeマニフェスト        | Kubernetesリソース           |

## 必要環境

- Node.js 22およびpnpm 9.15.9
- Go 1.25以上
- Androidビルド用JDK 21およびAndroid SDK
- マルチアーキテクチャイメージ用DockerとBuildx
- マニフェスト検証用`kubectl`とKustomize

## クイックスタート

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

## 検証

```bash
python3 -m unittest tests/test_repository_contract.py
python3 verify.py
```

検証ツールは、多言語文書、決定的なRelief図、Gitea workflow契約、保存したGitHub workflowのバイト列、各コンポーネントのビルドを確認します。

## 配布

![本番切替](docs/svg/runtime-cutover.ja.svg)

Giteaを正本とします。Gitea Actionsはリポジトリを検証した後、API、back、webイメージを直列に公開します。各イメージは`linux/amd64`と`linux/arm64`のOCI indexを含みます。GitHubはActionsを無効化したpush mirrorとして維持します。

## 運用

本番workloadは中央`seonology-k3s` GitOpsリポジトリの`workloads/seonology-journey`で管理します。Argo CDがSyncedかつHealthyで、3つのDeploymentがすべてReadyとなり、imageIDが検証済みGitea Registry digestを参照し、内部と公開health検査が成功した時点で完了と判定します。

## リポジトリ方針

![リポジトリ役割](docs/svg/repository-roles.ja.svg)

- 明示的なscopeを持つConventional Commitsを使用します。
- 移行証跡となる`.github/workflows`のバイト列を保存します。
- 認証情報はGitea secret、Vault、Kubernetes Secretだけに保存します。
- プロジェクト成果物に自動作成者署名や絵文字を追加しません。

## ライセンス

既存のMITライセンスを維持します。詳細は[LICENSE](LICENSE)を参照してください。
