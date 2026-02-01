# REST API Check Tool

APIレスポンスチェックツールは、指定したエンドポイントにリクエストを送り、レスポンス（ステータスコード、ボディ、実行時間）を表示するCLIツールです。AIがレスポンス内容を解析し、エラー解決策の提示やドキュメント生成を自動で行います。

## 機能

### Phase 1: 基本機能 (Core)
- ✅ HTTPリクエスト実行: GET, POST, PUT, DELETE などの主要メソッドへの対応
- ✅ カスタムヘッダー/ボディ: JSON形式のボディ送信、Authヘッダーの付与
- ✅ 結果表示: ステータスコード、レスポンスタイム（ms単位）、フォーマット済みJSONレスポンス

### Phase 2: CLI機能
- ✅ CobraライブラリによるCLI構造
- ✅ URL、メソッド、ヘッダー、ボディを引数で指定可能

### Phase 3: AI連携機能 (AI-Powered)
- ✅ エラー解決アドバイス: 4xx/5xxエラー時に、AIが原因を推測し「修正すべき箇所」を提案
- ✅ スキーマ自動生成: レスポンスのJSONから、Goの構造体（Struct）やOpenAPI（Swagger）定義を自動生成
- ✅ テストデータ推論: レスポンスに基づき、次にテストすべき境界値（Boundary Value）をAIが提案

### Phase 4: 運用機能
- ✅ 履歴保存: 過去に実行したリクエスト内容を保存し、再実行可能にする

## インストール

```bash
git clone <repository-url>
cd restapi_check
go mod download
```

## セットアップ

1. `.env.example`を`.env`にコピー
2. `.env`ファイルにOpenAI APIキーを設定

```bash
cp .env.example .env
# .envファイルを編集してOPENAI_API_KEYを設定
```

## 使用方法

### 基本的なリクエスト

```bash
# GETリクエスト
go run main.go request --url "https://api.example.com/users"

# POSTリクエスト
go run main.go request --url "https://api.example.com/users" --method POST --body '{"name":"test"}'

# ヘッダー付きリクエスト
go run main.go request --url "https://api.example.com/users" --header "Authorization: Bearer token"
```

### AI分析機能

```bash
# エラー時に自動分析
go run main.go request --url "https://api.example.com/users" --analyze

# エラー分析のみ
go run main.go analyze error --status 404 --body '{"error":"Not Found"}'

# Go構造体生成
go run main.go analyze struct --body '{"id":1,"name":"test"}'

# OpenAPIスキーマ生成
go run main.go analyze schema --body '{"id":1,"name":"test"}'

# テストデータ推論
go run main.go analyze test --body '{"id":1,"name":"test"}'
```

### 履歴機能

```bash
# 履歴一覧を表示
go run main.go history list

# 履歴の詳細を表示
go run main.go history show --id <履歴ID>

# 履歴を再実行
go run main.go history replay --id <履歴ID>
```

## プロジェクト構造

```
restapi_check/
├── cmd/              # CLIコマンド定義
│   ├── root.go      # ルートコマンド
│   ├── request.go   # リクエスト実行コマンド
│   ├── analyze.go   # AI分析コマンド
│   └── history.go   # 履歴管理コマンド
├── pkg/
│   ├── client/      # HTTPクライアント
│   ├── display/     # 表示フォーマッター
│   ├── config/      # 設定管理
│   ├── ai/          # AI連携
│   └── history/     # 履歴管理
├── main.go          # エントリーポイント
├── .env.example     # 環境変数テンプレート
└── README.md        # このファイル
```

## 開発フェーズ

- ✅ Phase 1: MVP - 基本的なHTTPリクエスト実行
- ✅ Phase 2: CLI引数対応 - Cobra導入
- ✅ Phase 3: AI連携 - OpenAI API連携
- ✅ Phase 4: 履歴保存機能

## ライセンス

MIT
