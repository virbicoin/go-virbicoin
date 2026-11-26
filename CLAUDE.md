# CLAUDE.md

本ファイルは、Claude / AI コーディングエージェントが本リポジトリのコードを扱う際の
ガイダンスを提供します。

## プロジェクト概要

go-virbicoin は VirBiCoin プロトコルの公式 Go 実装です。go-ethereum（1.9.x 系）の
フォークであり、メインコマンドは `gvbc`（VirBiCoin クライアント、アップストリームの
`geth` に相当）です。

- モジュールパス: `github.com/virbicoin/go-virbicoin`
- メインエントリポイント: `./cmd/gvbc`
- Go バージョン: `go 1.24.0` / `toolchain go1.26.4`（`go.mod` を参照）

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

release/X.Y の stable 化コミットと main のパッチ番号更新は、どちらも `params/version.go` の
同じ const ブロックを書き換えるため、main の取り込みはこのファイルで毎回競合します
（v1.9.40 では手で解決しました）。スクリプトは競合がこの 1 ファイルだけなら main 側を採用して
続行し（直後の手順で stable に固定）、ほかのファイルが競合したらマージを取り消して止まります。

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
| **virbicoin.com** | 公式 Web サイト（メインサイト。private） | `../virbicoin.com` | [github.com/virbicoin/virbicoin.com](https://github.com/virbicoin/virbicoin.com) |
| **go-virbicoin** ← 本リポジトリ | メインクライアント（Gvbc, Go 実装） | `../go-virbicoin` | [github.com/virbicoin/go-virbicoin](https://github.com/virbicoin/go-virbicoin) |
| **open-virbicoin** | Rust クライアント（Ovbc, OpenEthereum フォーク） | `../open-virbicoin` | [github.com/virbicoin/open-virbicoin](https://github.com/virbicoin/open-virbicoin) |
| **vbc-explorer** | ブロックチェーンエクスプローラー（DEX・Launchpad・ブリッジ画面を含む） | `../vbc-explorer` | [github.com/virbicoin/vbc-explorer](https://github.com/virbicoin/vbc-explorer) |
| **vbc-stats** | ネットワーク統計ダッシュボード | `../vbc-stats` | [github.com/virbicoin/vbc-stats](https://github.com/virbicoin/vbc-stats) |
| **vbc-stats-api** | VBC Stats 用ノードレポーター（組み込みレポーターのない Ovbc 向け） | `../vbc-stats-api` | [github.com/virbicoin/vbc-stats-api](https://github.com/virbicoin/vbc-stats-api) |
| **vbc-pool** | マイニングプール | `../vbc-pool` | [github.com/virbicoin/vbc-pool](https://github.com/virbicoin/vbc-pool) |
| **vbc-rpc** | RPC ノードステータスページ（公開 RPC・WebSocket の案内） | `../vbc-rpc` | [github.com/virbicoin/vbc-rpc](https://github.com/virbicoin/vbc-rpc) |
| **vbc-hash** | Ethash アルゴリズム実装（C/Go バインディング） | `../vbc-hash` | [github.com/virbicoin/vbc-hash](https://github.com/virbicoin/vbc-hash) |

### 依存関係

- **go-virbicoin** と **open-virbicoin**: 同じ VirBiCoin ネットワーク（chainId 329）の 2 つのクライアント実装（Gvbc は Go、Ovbc は Rust）。報酬・フォークブロックなどコンセンサスに関わる変更は必ず両方に同時に入れる
- **vbc-explorer** → ノードの JSON-RPC（Gvbc / Ovbc）。内部取引の索引は Parity 形式の `trace_transaction` を使うため Ovbc が必要（本番の explorer は Ovbc の IPC ソケットを読む）
- **vbc-stats** ← Gvbc は組み込みの ethstats レポーター（`--ethstats`）で直接報告し、Ovbc は **vbc-stats-api** を介して報告する
- **vbc-pool** → ノードの JSON-RPC（Gvbc / Ovbc。本番は両方）から作業を取得し、ブロックの確定と支払いを行う。シェアの検証は `proxy/miner.go` の純 Go 実装で、vbc-hash には依存しない
- **vbc-rpc** → 各ノードの JSON-RPC（8329）に状態を問い合わせるだけ。公開 RPC（rpc.virbicoin.com）への POST は nginx がノードへ中継し、このアプリを通らない
- **vbc-hash**: このワークスペースのどのリポジトリも import していない（Gvbc は `consensus/ethash`、vbc-pool は `proxy/miner.go` の自前実装）
- **virbicoin.com**: 静的サイト（他リポジトリへのリンクのみ）
- 公開リポジトリなので、private のリポジトリはこの表に載せていない（virbicoin.com は公開 Web サイトのソースとして載せる）。公開 README からは virbicoin.com の GitHub ではなく Web サイト（https://www.virbicoin.com）にリンクする
