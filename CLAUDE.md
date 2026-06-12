# CLAUDE.md

本ファイルは、Claude / AI コーディングエージェントが本リポジトリのコードを扱う際の
ガイダンスを提供します。

## プロジェクト概要

go-virbicoin は VirBiCoin プロトコルの公式 Go 実装です。go-ethereum（1.9.x 系）の
フォークであり、メインコマンドは `gvbc`（VirBiCoin クライアント、アップストリームの
`geth` に相当）です。

- モジュールパス: `github.com/virbicoin/go-virbicoin`
- メインエントリポイント: `./cmd/gvbc`
- Go バージョン: `go 1.24.0` / `toolchain go1.24.5`（`go.mod` を参照）

## ビルドと主要コマンド

```shell
# gvbc クライアントを ./build/bin/gvbc にビルド
make gvbc

# 全ツールをビルド（gvbc, bootnode, clef, evm, abigen など）
make all

# Linter を実行
make lint

# テストスイート全体を実行
go test ./...

# 単一パッケージのテストを実行
go test ./core/...

# 特定パッケージを vet
go vet ./rpc/...
```

CGO と C コンパイラが必要です（secp256k1、leveldb などで使用）。

## アーキテクチャ

主要なパッケージ構成（go-ethereum から継承）:

- `cmd/` — コマンドラインのエントリポイント。`cmd/gvbc` がメインのノードバイナリ。
- `core/` — コアブロックチェーン: ブロック処理、ステート、EVM、トランザクションプール。
- `consensus/` — コンセンサスエンジン（`ethash`、`clique`）。
- `eth/`、`les/` — フルおよびライトプロトコルの実装。
- `p2p/` — ピアツーピアのネットワーキングとディスカバリ。
- `params/` — ネットワークパラメータ、バージョン、ブートノード一覧。
- `rpc/` — JSON-RPC サーバ（プラットフォーム固有の IPC トランスポートを含む）。
- `accounts/` — アカウント管理、キーストア、ハードウェアウォレット対応。
- `common/` — 共有ユーティリティ。クロスプラットフォーム（windows/arm64 を含む）
  互換性のため、アップストリームの `prometheus/tsdb` 依存の代わりに内部の
  `common/fileutil` パッケージ（ファイルロック）を含む。

## VirBiCoin 固有の注意点

- バージョンは `params/version.go` で定義
  （`VersionMajor` / `VersionMinor` / `VersionPatch`）。
- ブートノードは `params/bootnodes.go` に記載。
- ジェネシス設定は `assets/genesis.json` にある。

## 関連サイト・ドキュメント

- **メインサイト**: `../virbicoin.com`（このリポジトリの兄弟ディレクトリにクローン
  済み。GitHub: `virbicoin/virbicoin.com`、本番 URL: <https://virbicoin.com>）。
  VirBiCoin プロジェクトの**公式メインサイト**であり、Next.js（app-router、
  `output: 'export'` の静的エクスポート）で構築。プロトコルレベルのドキュメントは
  `/docs`（`/docs/json-rpc`、`/docs/rlp`、`/docs/enode`）に配置。
- **クライアント（gvbc）ドキュメント**: go-virbicoin の GitHub Wiki
  （<https://github.com/virbicoin/go-virbicoin/wiki>）。Command-Line-Options、
  JavaScript-Console、JSON-RPC-Server、Native-Bindings、Developers-Guide など。
- **GitHub Pages**: `gh-pages` ブランチ（<https://virbicoin.github.io/go-virbicoin/>）。
  メインサイト・GitHub・Wiki・`virbicoin.com/docs` へ誘導する単一の静的ランディング
  ページ（`index.html` + `.nojekyll`）。
- ドキュメントの方針: プロトコルレベル＝メインサイト `virbicoin.com/docs`、
  クライアントレベル＝go-virbicoin Wiki に分離し、内容が重複しないようにする。
  README 等から Ethereum / geth.ethereum.org のドキュメントへ直接リンクせず、
  上記 VirBiCoin 所有先を参照する。

## プラットフォーム対応

リリースは GoReleaser（`.goreleaser.yaml`）で 8 ターゲット向けにビルドされます:

- linux: amd64, arm64, 386
- darwin: amd64, arm64
- windows: amd64, 386, arm64

windows/arm64 対応のため、新たに追加する依存はアーキテクチャ非依存である必要が
あります。`386` / `amd64` 専用のアセンブリやシステムコールをハードコードする依存の
再導入は避けてください（これが `prometheus/tsdb/fileutil`、`natefinch/npipe`、旧
`go-ole` を置き換えた理由です）。

## リリース手順

リリースは GoReleaser で生成します:

```shell
goreleaser release --skip=validate --clean
```

`--release-notes` を指定しないことで、GoReleaser が前回タグ以降のコミットから
チェンジログを自動生成します。リリースノートのヘッダ（対応プラットフォーム一覧と
インストール手順）は `.goreleaser.yaml` の `release.header` に英語で定義しています。
アーカイブはトップレベルのディレクトリにまとめられ（`wrap_in_directory: true`）、
リリース tarball を展開してもファイルがカレントディレクトリに散らばりません。
リリースはドラフト（`draft: true`）として作成され、手動で公開します。

チェンジログは Git のコミットメッセージから生成されるため、コミットメッセージは
英語で記述してください（`^docs:` / `^test:` / `^ci:` やマージコミットは除外設定済み）。

## リリースサイクル

バージョンは go-ethereum 流の unstable/stable サイクルで管理します
（`params/version.go` の `VersionMeta`）。

### ブランチモデル

- `main` — 開発本流。常に `vX.Y.Z` + `"unstable"`。
- `dev` — 機能統合・検証用。
- `release/X.Y` — メンテナンスライン。stable コミットとリリースタグはここに置く。

### サイクル

1. **開発フェーズ**: `main` は `vX.Y.Z` + `"unstable"`（例: v1.9.39 unstable）
2. **リリース時**: `release/X.Y` へ `main` を取り込み、`"unstable"` → `"stable"` に
   変更してタグ付け・GoReleaser で公開
3. **リリース後**: `main` のパッチ番号を +1 して次の開発サイクルへ
   （→ vX.Y.(Z+1) unstable、main は unstable のまま）

この流れは `build/release.sh` で半自動化されています:

```shell
# main の unstable バージョンを release/X.Y で stable リリースし、
# main を次パッチの unstable に進める
build/release.sh

# 実際の変更を行わず手順だけ確認
build/release.sh --dry-run
```

スクリプトは main ブランチ・クリーンな作業ツリーであることを確認し、release/X.Y への
main 取り込み → stable 化コミット → タグ push → 既存ドラフト削除 → GoReleaser 実行
→ main を次パッチ unstable に更新、までを一括で行います。`GITHUB_TOKEN` は
`~/.gvbc_token.env` から自動で読み込まれます。

## CI

継続的インテグレーションは GitHub Actions（`.github/workflows/go.yml`）で実行され
ます。レガシーの AppVeyor / CircleCI / Travis 設定は削除済みです。

## 言語ポリシー

メンテナンス性のため、リポジトリ内のテキストは原則として英語で統一します。

- **英語**: ソースコードのコメント、コミットメッセージ、チェンジログ、シェル
  スクリプト、設定ファイル、`README.md`、その他すべてのドキュメント。
- **日本語（例外）**: `CLAUDE.md`（本ファイル、`.gitignore` 済みでローカル専用）と
  `README-JP.md` のみ。

## 規約

- 本リポジトリは go-ethereum のフォークです。アップストリームのコード構造と
  ライセンスヘッダ（ライブラリコードは LGPL-3.0、`cmd/` は GPL-3.0）を維持して
  ください。
- タスクで直接必要な変更のみを行い、アップストリーム由来コードの広範な
  リファクタリングは避けてください。
- import はグループ化を維持し、`github.com/virbicoin/go-virbicoin/...` の import は
  ローカルグループ内にアルファベット順で配置してください。

## コミット署名（GPG）

このリポジトリのコミットは GPG 署名が有効です（`commit.gpgsign`）。AI エージェントは
秘密情報であるパスフレーズを代理入力できないため、gpg-agent のキャッシュが切れていると
`git commit` が署名失敗で中断することがあります。

- 署名が切れているときは、ユーザーがターミナルで一度パスフレーズを入力してください
  （`git commit` の再実行、または `echo test | gpg --clearsign` を一度実行）。一度
  入力すれば gpg-agent がしばらくキャッシュします。
- パスフレーズは秘密情報です。AI エージェントへ渡したりディスクへ保存したりしないで
  ください。
- コミット失敗を未然に防ぎたい場合は、コミット前にキャッシュを温める pre-commit フック
  （署名キャッシュが切れていればパスフレーズ入力を促す）を利用する方法があります。

## 関連リポジトリ

VirBiCoin エコシステムは以下のリポジトリで構成されています:

| リポジトリ | 役割 | ローカルパス | URL |
|-----------|------|-------------|-----|
| **virbicoin.com** | 公式 Web サイト（メインサイト） | `../virbicoin.com` | [github.com/virbicoin/virbicoin.com](https://github.com/virbicoin/virbicoin.com) |
| **go-virbicoin** ← 本リポジトリ | メインクライアント（Gvbc, Go 実装） | `../go-virbicoin` | [github.com/virbicoin/go-virbicoin](https://github.com/virbicoin/go-virbicoin) |
| **openvirbicoin** | Rust クライアント（Ovbc, OpenEthereum フォーク） | `../openvirbicoin` | [github.com/virbicoin/openvirbicoin](https://github.com/virbicoin/openvirbicoin) |
| **vbc-stats** | ネットワーク統計ダッシュボード | `../vbc-stats` | [github.com/virbicoin/vbc-stats](https://github.com/virbicoin/vbc-stats) |
| **vbc-explorer** | ブロックチェーンエクスプローラー | `../vbc-explorer` | [github.com/virbicoin/vbc-explorer](https://github.com/virbicoin/vbc-explorer) |
| **open-virbicoin-pool** | マイニングプール | `../open-virbicoin-pool` | [github.com/virbicoin/open-virbicoin-pool](https://github.com/virbicoin/open-virbicoin-pool) |
| **vbc-rpc** | RPC ノードステータス & JSON-RPC プロキシ | `../vbc-rpc` | [github.com/virbicoin/vbc-rpc](https://github.com/virbicoin/vbc-rpc) |

### 依存関係（go-virbicoin が中心）

- **openvirbicoin**: 同一の VirBiCoin ネットワーク（chainId 329）に接続する代替クライアント（Ovbc, Rust 実装）。Gvbc とは別実装
- **vbc-stats** → **go-virbicoin**: Gvbc ノードが eth-netstats-client プロトコルでブロック/統計データを送信
- **vbc-explorer** → **go-virbicoin**: JSON-RPC 経由でブロックチェーンデータを取得
- **open-virbicoin-pool** → **go-virbicoin**: マイニングプールが Gvbc ノードから作業を取得
- **vbc-rpc** → **go-virbicoin**: RPC プロキシが Gvbc ノードにリクエストを中継
