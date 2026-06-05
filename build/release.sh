#!/usr/bin/env bash
#
# release.sh — VirBiCoin リリースサイクルを半自動化するヘルパー。
#
# ブランチモデル（go-ethereum 流）:
#   - main          : 開発本流。常に vX.Y.Z + "unstable"。
#   - dev           : 機能統合・検証用。
#   - release/X.Y   : メンテナンスライン。stable コミットとリリースタグはここに置く。
#
# リリースサイクル:
#   1. 開発フェーズ : main は vX.Y.Z + "unstable"
#   2. リリース時   : release/X.Y へ main を取り込み "stable" 化してタグ付け・公開
#   3. リリース後   : main のパッチ番号を +1 して次の開発サイクルへ（unstable 維持）
#
# 使い方:
#   build/release.sh            main の unstable バージョンを release/X.Y で stable
#                               リリースし、main を次パッチの unstable に進める
#   build/release.sh --dry-run  実際の変更・push・リリースを行わず手順だけ表示する
#
# 前提:
#   - クリーンな作業ツリー（未コミットの変更が無いこと）
#   - main ブランチで実行すること
#   - GITHUB_TOKEN が設定されていること（~/.gvbc_token.env を自動 source）
#   - goreleaser がインストール済みであること
#
set -euo pipefail

cd "$(git rev-parse --show-toplevel)"

VERSION_FILE="params/version.go"
REPO="virbicoin/go-virbicoin"
DRY_RUN=0
[[ "${1:-}" == "--dry-run" ]] && DRY_RUN=1

log()  { printf '\033[1;34m==>\033[0m %s\n' "$*"; }
err()  { printf '\033[1;31mエラー:\033[0m %s\n' "$*" >&2; exit 1; }
run()  { if [[ $DRY_RUN -eq 1 ]]; then printf '   (dry-run) %s\n' "$*"; else eval "$*"; fi; }

# --- トークン読み込み ---
if [[ -z "${GITHUB_TOKEN:-}" && -f "$HOME/.gvbc_token.env" ]]; then
  # shellcheck disable=SC1091
  source "$HOME/.gvbc_token.env"
fi

# --- 事前チェック ---
command -v goreleaser >/dev/null || err "goreleaser がインストールされていません"
[[ -n "${GITHUB_TOKEN:-}" ]] || err "GITHUB_TOKEN が設定されていません（~/.gvbc_token.env を確認）"

branch="$(git rev-parse --abbrev-ref HEAD)"
[[ "$branch" == "main" ]] || err "main ブランチで実行してください（現在: $branch）"

if [[ -n "$(git status --porcelain)" ]]; then
  err "作業ツリーに未コミットの変更があります。コミットまたは退避してください"
fi

# --- 現在のバージョンを読み取り ---
get_field() { grep -E "^\s*$1\s*=" "$VERSION_FILE" | head -1 | sed -E 's/.*=\s*//; s/\s*\/\/.*//; s/"//g; s/\s*$//'; }
MAJOR="$(get_field VersionMajor)"
MINOR="$(get_field VersionMinor)"
PATCH="$(get_field VersionPatch)"
META="$(get_field VersionMeta)"

CUR_TAG="v${MAJOR}.${MINOR}.${PATCH}"
RELEASE_BRANCH="release/${MAJOR}.${MINOR}"
log "現在のバージョン: ${CUR_TAG} (${META})"
log "リリースブランチ: ${RELEASE_BRANCH}"

[[ "$META" == "unstable" ]] || err "VersionMeta が unstable ではありません（現在: ${META}）。リリース済みの可能性があります"

NEXT_PATCH=$((PATCH + 1))
NEXT_TAG="v${MAJOR}.${MINOR}.${NEXT_PATCH}"

# 後始末: 失敗・終了時には必ず main へ戻す
cleanup() { git checkout main >/dev/null 2>&1 || true; }
trap cleanup EXIT

# --- 0) 最新を取得 ---
log "[0/6] リモートの最新を取得"
run "git fetch origin --prune"

# --- 1) release/X.Y を用意し main を取り込む ---
log "[1/6] ${RELEASE_BRANCH} に main を取り込む"
if git show-ref --verify --quiet "refs/remotes/origin/${RELEASE_BRANCH}"; then
  run "git checkout -B \"$RELEASE_BRANCH\" \"origin/${RELEASE_BRANCH}\""
  run "git merge --no-edit main"
else
  log "  ${RELEASE_BRANCH} が無いため main から新規作成"
  run "git checkout -B \"$RELEASE_BRANCH\" main"
fi

# --- 2) release/X.Y 上で stable 化してコミット ---
log "[2/6] ${CUR_TAG} を stable に変更してコミット（${RELEASE_BRANCH}）"
run "sed -i -E 's/(VersionMeta\s*=\s*)\"unstable\"/\\1\"stable\"   /' \"$VERSION_FILE\""
run "git add \"$VERSION_FILE\""
run "git commit -m \"${CUR_TAG} リリース: VersionMeta を stable に変更\""

# --- 3) タグ付けして push ---
log "[3/6] タグ ${CUR_TAG} を作成して ${RELEASE_BRANCH} と共に push"
run "git tag -f \"$CUR_TAG\""
run "git push origin \"$RELEASE_BRANCH\""
run "git push -f origin \"$CUR_TAG\""

# --- 4) 既存ドラフトがあれば削除 ---
log "[4/6] 同タグの既存ドラフトを確認・削除"
if [[ $DRY_RUN -eq 0 ]]; then
  RID="$(curl -s -H "Authorization: token $GITHUB_TOKEN" \
    "https://api.github.com/repos/${REPO}/releases?per_page=10" \
    | python3 -c "import sys,json
for r in json.load(sys.stdin):
    if r.get('tag_name')=='$CUR_TAG' and r.get('draft'):
        print(r['id']); break" 2>/dev/null || true)"
  if [[ -n "${RID:-}" ]]; then
    curl -s -X DELETE -H "Authorization: token $GITHUB_TOKEN" \
      "https://api.github.com/repos/${REPO}/releases/$RID" \
      -w "   既存ドラフト($RID)を削除: HTTP=%{http_code}\n"
  else
    echo "   既存ドラフトなし"
  fi
else
  echo "   (dry-run) ドラフト確認・削除をスキップ"
fi

# --- 5) GoReleaser でドラフトリリース（release/X.Y をチェックアウトした状態で実行）---
log "[5/6] GoReleaser で ${CUR_TAG} のドラフトを生成（自動 changelog）"
run "goreleaser release --skip=validate --clean"

# --- 6) main を次の開発サイクル（unstable）へ ---
log "[6/6] main を ${NEXT_TAG} unstable に更新して次の開発サイクルを開始"
run "git checkout main"
run "sed -i -E 's/(VersionPatch\s*=\s*)[0-9]+/\\1${NEXT_PATCH}/' \"$VERSION_FILE\""
# main は元々 unstable なので META 行はそのまま（patch のみ更新）
run "git add \"$VERSION_FILE\""
run "git commit -m \"次の開発サイクル: ${NEXT_TAG} unstable\""
run "git push origin main"

log "完了: ${CUR_TAG} を ${RELEASE_BRANCH} からリリース（ドラフト）し、main を ${NEXT_TAG} unstable に更新しました"
echo "   GitHub のリリースページでドラフトを確認し、問題なければ Publish してください。"
