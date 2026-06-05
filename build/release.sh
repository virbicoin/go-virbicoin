#!/usr/bin/env bash
#
# release.sh — VirBiCoin リリースサイクルを半自動化するヘルパー。
#
# サイクル:
#   1. 開発フェーズ: params/version.go は vX.Y.Z + "unstable"
#   2. リリース時 : "unstable" -> "stable" に変更し、タグ付け・GoReleaser でリリース
#   3. リリース後 : パッチ番号を +1 し "unstable" に戻して次の開発サイクルへ
#
# 使い方:
#   build/release.sh            現在の unstable バージョンを stable としてリリースし、
#                               次のパッチを unstable にして開発サイクルを開始する
#   build/release.sh --dry-run  実際の変更・push・リリースを行わず手順だけ表示する
#
# 前提:
#   - クリーンな作業ツリー（未コミットの変更が無いこと）
#   - main ブランチであること
#   - GITHUB_TOKEN が設定されていること（~/.gvbc_token.env を自動 source）
#   - goreleaser がインストール済みであること
#
set -euo pipefail

cd "$(git rev-parse --show-toplevel)"

VERSION_FILE="params/version.go"
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
log "現在のバージョン: ${CUR_TAG} (${META})"

[[ "$META" == "unstable" ]] || err "VersionMeta が unstable ではありません（現在: ${META}）。リリース済みの可能性があります"

NEXT_PATCH=$((PATCH + 1))
NEXT_TAG="v${MAJOR}.${MINOR}.${NEXT_PATCH}"

# --- 1) stable 化してコミット ---
log "[1/5] ${CUR_TAG} を stable に変更してコミット"
run "sed -i -E 's/(VersionMeta\s*=\s*)\"unstable\"/\\1\"stable\"   /' \"$VERSION_FILE\""
run "git add \"$VERSION_FILE\""
run "git commit -m \"${CUR_TAG} リリース: VersionMeta を stable に変更\""

# --- 2) タグ付けして push ---
log "[2/5] タグ ${CUR_TAG} を作成して push"
run "git tag -f \"$CUR_TAG\""
run "git push origin main"
run "git push -f origin \"$CUR_TAG\""

# --- 3) 既存ドラフトがあれば削除 ---
log "[3/5] 同タグの既存ドラフトを確認・削除"
if [[ $DRY_RUN -eq 0 ]]; then
  RID="$(curl -s -H "Authorization: token $GITHUB_TOKEN" \
    "https://api.github.com/repos/virbicoin/go-virbicoin/releases?per_page=10" \
    | python3 -c "import sys,json
for r in json.load(sys.stdin):
    if r.get('tag_name')=='$CUR_TAG' and r.get('draft'):
        print(r['id']); break" 2>/dev/null || true)"
  if [[ -n "${RID:-}" ]]; then
    curl -s -X DELETE -H "Authorization: token $GITHUB_TOKEN" \
      "https://api.github.com/repos/virbicoin/go-virbicoin/releases/$RID" \
      -w "   既存ドラフト($RID)を削除: HTTP=%{http_code}\n"
  else
    echo "   既存ドラフトなし"
  fi
else
  echo "   (dry-run) ドラフト確認・削除をスキップ"
fi

# --- 4) GoReleaser でドラフトリリース ---
log "[4/5] GoReleaser で ${CUR_TAG} のドラフトを生成（自動 changelog）"
run "goreleaser release --skip=validate --clean"

# --- 5) 次の開発サイクル（unstable）へ ---
log "[5/5] ${NEXT_TAG} unstable に更新して次の開発サイクルを開始"
run "sed -i -E 's/(VersionPatch\s*=\s*)[0-9]+/\\1${NEXT_PATCH}/' \"$VERSION_FILE\""
run "sed -i -E 's/(VersionMeta\s*=\s*)\"stable\"\s*/\\1\"unstable\"/' \"$VERSION_FILE\""
run "git add \"$VERSION_FILE\""
run "git commit -m \"次の開発サイクル: ${NEXT_TAG} unstable\""
run "git push origin main"

log "完了: ${CUR_TAG} をリリース（ドラフト）し、開発版を ${NEXT_TAG} unstable に更新しました"
echo "   GitHub のリリースページでドラフトを確認し、問題なければ Publish してください。"
