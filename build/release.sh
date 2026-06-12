#!/usr/bin/env bash
#
# release.sh — Semi-automated helper for the VirBiCoin release cycle.
#
# Branch model (go-ethereum style):
#   - main          : Mainline development. Always vX.Y.Z + "unstable".
#   - dev           : Feature integration and verification.
#   - release/X.Y   : Maintenance line. Stable commits and release tags live here.
#
# Release cycle:
#   1. Development : main is vX.Y.Z + "unstable"
#   2. Release     : merge main into release/X.Y, flip to "stable", tag, publish
#   3. Post-release: bump main's patch number for the next development cycle
#                    (main stays "unstable")
#
# Usage:
#   build/release.sh            Release main's unstable version as a stable build
#                               on release/X.Y, then advance main to the next
#                               unstable patch.
#   build/release.sh --dry-run  Print the steps without making any changes,
#                               pushes, or releases.
#
# Requirements:
#   - Clean working tree (no uncommitted changes)
#   - Run from the main branch
#   - GITHUB_TOKEN must be set (~/.gvbc_token.env is sourced automatically)
#   - goreleaser must be installed
#
set -euo pipefail

cd "$(git rev-parse --show-toplevel)"

VERSION_FILE="params/version.go"
REPO="virbicoin/go-virbicoin"
DRY_RUN=0
[[ "${1:-}" == "--dry-run" ]] && DRY_RUN=1

log()  { printf '\033[1;34m==>\033[0m %s\n' "$*"; }
err()  { printf '\033[1;31mError:\033[0m %s\n' "$*" >&2; exit 1; }
run()  { if [[ $DRY_RUN -eq 1 ]]; then printf '   (dry-run) %s\n' "$*"; else eval "$*"; fi; }

# --- Load token ---
if [[ -z "${GITHUB_TOKEN:-}" && -f "$HOME/.gvbc_token.env" ]]; then
  # shellcheck disable=SC1091
  source "$HOME/.gvbc_token.env"
fi

# --- Preflight checks ---
command -v goreleaser >/dev/null || err "goreleaser is not installed"
[[ -n "${GITHUB_TOKEN:-}" ]] || err "GITHUB_TOKEN is not set (check ~/.gvbc_token.env)"

branch="$(git rev-parse --abbrev-ref HEAD)"
[[ "$branch" == "main" ]] || err "Run from the main branch (current: $branch)"

if [[ -n "$(git status --porcelain)" ]]; then
  err "Working tree has uncommitted changes. Commit or stash them first."
fi

# --- Read the current version ---
get_field() { grep -E "^\s*$1\s*=" "$VERSION_FILE" | head -1 | sed -E 's/.*=\s*//; s/\s*\/\/.*//; s/"//g; s/\s*$//'; }
MAJOR="$(get_field VersionMajor)"
MINOR="$(get_field VersionMinor)"
PATCH="$(get_field VersionPatch)"
META="$(get_field VersionMeta)"

CUR_TAG="v${MAJOR}.${MINOR}.${PATCH}"
RELEASE_BRANCH="release/${MAJOR}.${MINOR}"
log "Current version: ${CUR_TAG} (${META})"
log "Release branch:  ${RELEASE_BRANCH}"

[[ "$META" == "unstable" ]] || err "VersionMeta is not 'unstable' (current: ${META}). It may already be released."

NEXT_PATCH=$((PATCH + 1))
NEXT_TAG="v${MAJOR}.${MINOR}.${NEXT_PATCH}"

# Cleanup: always return to main on exit
cleanup() { git checkout main >/dev/null 2>&1 || true; }
trap cleanup EXIT

# --- 0) Fetch latest ---
log "[0/6] Fetch the latest from the remote"
run "git fetch origin --prune"

# --- 1) Prepare release/X.Y and merge main ---
log "[1/6] Merge main into ${RELEASE_BRANCH}"
if git show-ref --verify --quiet "refs/remotes/origin/${RELEASE_BRANCH}"; then
  run "git checkout -B \"$RELEASE_BRANCH\" \"origin/${RELEASE_BRANCH}\""
  run "git merge --no-edit main"
else
  log "  ${RELEASE_BRANCH} does not exist; creating it from main"
  run "git checkout -B \"$RELEASE_BRANCH\" main"
fi

# --- 2) Flip to stable and commit on release/X.Y ---
log "[2/6] Set ${CUR_TAG} to stable and commit (${RELEASE_BRANCH})"
# The merge can leave version.go at the previous release (the stable-flip
# commit on release/X.Y reverts main's patch bump), so pin both fields
# explicitly instead of only rewriting "unstable".
run "sed -i -E 's/(VersionPatch\s*=\s*)[0-9]+/\\1${PATCH}/' \"$VERSION_FILE\""
run "sed -i -E 's/(VersionMeta\s*=\s*)\"unstable\"/\\1\"stable\"/' \"$VERSION_FILE\""
run "gofmt -w \"$VERSION_FILE\""
run "git add \"$VERSION_FILE\""
if [[ $DRY_RUN -eq 1 ]] || ! git diff --cached --quiet; then
  run "git commit -m \"Release ${CUR_TAG}: set VersionMeta to stable\""
else
  log "  version.go already at ${CUR_TAG} stable after the merge; skipping commit"
fi

# --- 3) Tag and push ---
log "[3/6] Create tag ${CUR_TAG} and push it with ${RELEASE_BRANCH}"
run "git tag -f \"$CUR_TAG\""
run "git push origin \"$RELEASE_BRANCH\""
run "git push -f origin \"$CUR_TAG\""

# --- 4) Delete any existing draft for this tag ---
log "[4/6] Check for and delete an existing draft for this tag"
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
      -w "   Deleted existing draft ($RID): HTTP=%{http_code}\n"
  else
    echo "   No existing draft"
  fi
else
  echo "   (dry-run) Skipping draft check/delete"
fi

# --- 5) Build the draft release with GoReleaser (on release/X.Y) ---
log "[5/6] Build the ${CUR_TAG} draft with GoReleaser (auto changelog)"
run "goreleaser release --skip=validate --clean"

# --- 6) Advance main to the next development cycle (unstable) ---
log "[6/6] Bump main to ${NEXT_TAG} unstable for the next development cycle"
run "git checkout main"
run "sed -i -E 's/(VersionPatch\s*=\s*)[0-9]+/\\1${NEXT_PATCH}/' \"$VERSION_FILE\""
# main is already unstable, so only the patch number changes
run "git add \"$VERSION_FILE\""
run "git commit -m \"Begin next development cycle: ${NEXT_TAG} unstable\""
run "git push origin main"

log "Done: released ${CUR_TAG} (draft) from ${RELEASE_BRANCH} and bumped main to ${NEXT_TAG} unstable"
echo "   Review the draft on the GitHub releases page and publish it when ready."
