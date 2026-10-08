#!/usr/bin/env bash
# Deletes Cloudflare Pages deployments that are no longer needed:
#   - production: keeps the newest $KEEP_PRODUCTION (the live one plus a couple to roll back to)
#   - preview: keeps only the newest deployment of each branch in $KEEP_BRANCHES
#     (the "www" redirect plus branches with an open pull request); deletes the rest
# Needs: CLOUDFLARE_API_TOKEN (Pages: Edit), CLOUDFLARE_ACCOUNT_ID, PAGES_PROJECT.
# DRY_RUN=1 lists what would be deleted without deleting anything.
set -euo pipefail

: "${CLOUDFLARE_API_TOKEN:?}" "${CLOUDFLARE_ACCOUNT_ID:?}" "${PAGES_PROJECT:?}"
KEEP_PRODUCTION="${KEEP_PRODUCTION:-3}"
KEEP_BRANCHES="${KEEP_BRANCHES:-www}"
API="https://api.cloudflare.com/client/v4/accounts/$CLOUDFLARE_ACCOUNT_ID/pages/projects/$PAGES_PROJECT/deployments"

cf() { curl -sS --fail-with-body -H "Authorization: Bearer $CLOUDFLARE_API_TOKEN" "$@"; }

# All deployments of one environment, newest first, as "id branch created_on" lines
list() {
  local page=1 out
  while :; do
    out=$(cf "$API?env=$1&page=$page&per_page=25")
    jq -r '.result[] | "\(.id) \(.deployment_trigger.metadata.branch // "-") \(.created_on)"' <<<"$out"
    [ "$page" -ge "$(jq -r '.result_info.total_pages // 1' <<<"$out")" ] && break
    page=$((page + 1))
  done | sort -k3 -r
}

remove() {
  if [ "${DRY_RUN:-0}" = 1 ]; then echo "would delete $1 ($2)"; return; fi
  # force=true is needed for deployments that still have a branch alias
  if cf -X DELETE "$API/$1?force=true" >/dev/null; then echo "deleted $1 ($2)"
  else echo "could not delete $1 ($2)"; fi # e.g. the live production deployment
}

deleted=0
n=0
while read -r id branch _; do
  n=$((n + 1))
  if [ "$n" -gt "$KEEP_PRODUCTION" ]; then remove "$id" "production $branch"; deleted=$((deleted + 1)); fi
done < <(list production)

declare -A seen=()
while read -r id branch _; do
  if [[ " $KEEP_BRANCHES " == *" $branch "* ]] && [ -z "${seen[$branch]:-}" ]; then
    seen[$branch]=1; continue
  fi
  remove "$id" "preview $branch"; deleted=$((deleted + 1))
done < <(list preview)

echo "Deployments removed: $deleted"
