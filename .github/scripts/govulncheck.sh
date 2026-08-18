#!/usr/bin/env bash
#
# Run govulncheck and fail only on vulnerabilities that check-self actually
# calls and that are not on the accepted-risk list.
#
# Accepted risks live in .github/govulncheck-ignore.txt, one Go vulnerability
# database ID per line, each with a justification and a re-check date.
#
# Usage: .github/scripts/govulncheck.sh [packages...]   (default: ./...)

set -euo pipefail

repo_root="$(git rev-parse --show-toplevel)"
ignore_file="${GOVULNCHECK_IGNORE_FILE:-${repo_root}/.github/govulncheck-ignore.txt}"

if [[ $# -gt 0 ]]; then
    packages=("$@")
else
    packages=("./...")
fi

for tool in govulncheck jq; do
    if ! command -v "${tool}" >/dev/null 2>&1; then
        printf 'error: %s was not found on PATH\n' "${tool}" >&2
        exit 1
    fi
done

if [[ ! -f "${ignore_file}" ]]; then
    printf 'error: accepted-risk list not found: %s\n' "${ignore_file}" >&2
    exit 1
fi

work_dir="$(mktemp -d)"
trap 'rm -rf "${work_dir}"' EXIT
report="${work_dir}/govulncheck.json"

# "-format json" completes with status 0 even when vulnerabilities are found, so
# a non-zero status here means the scan itself failed and must not be swallowed.
govulncheck -format json "${packages[@]}" >"${report}"

# trace[0] is the vulnerable frame. Its .function is set only for symbol-level
# findings, which is what separates "our code calls this" from "we merely
# import the module".
jq -sr '
    [ .[]
      | select(.finding)
      | .finding
      | select(.trace[0].function != null)
      | .osv
    ] | unique | .[]
' "${report}" | LC_ALL=C sort -u >"${work_dir}/called"

sed 's/#.*//' "${ignore_file}" | awk 'NF { print $1 }' | LC_ALL=C sort -u >"${work_dir}/accepted"

LC_ALL=C comm -23 "${work_dir}/called" "${work_dir}/accepted" >"${work_dir}/unexpected"
LC_ALL=C comm -13 "${work_dir}/called" "${work_dir}/accepted" >"${work_dir}/stale"

if [[ -s "${work_dir}/stale" ]]; then
    printf 'notice: no longer reported, remove from %s:\n' "${ignore_file##*/}"
    sed 's/^/  /' "${work_dir}/stale"
    printf '\n'
fi

if [[ ! -s "${work_dir}/unexpected" ]]; then
    printf 'govulncheck: no unaccepted vulnerabilities (%s accepted).\n' \
        "$(wc -l <"${work_dir}/accepted" | tr -d ' ')"
    exit 0
fi

printf 'govulncheck: %s unaccepted vulnerability(s) found:\n\n' \
    "$(wc -l <"${work_dir}/unexpected" | tr -d ' ')"

while read -r osv; do
    jq -sr --arg osv "${osv}" '
        ([ .[] | select(.osv) | .osv | select(.id == $osv) ] | first) as $meta
        | ([ .[]
             | select(.finding)
             | .finding
             | select(.osv == $osv and .trace[0].function != null)
           ] | first) as $hit
        | ($hit.trace[0]) as $vuln
        | ($hit.trace[-1]) as $entry
        | "\($osv): \($meta.summary // "(no summary)")",
          "  module:    \($vuln.module)@\($vuln.version // "?")",
          "  fixed in:  \($hit.fixed_version // "no fixed version")",
          "  called by: \($entry.position.filename // "?"):\($entry.position.line // 0) (\($entry.function // "?"))",
          "  more info: https://pkg.go.dev/vuln/\($osv)"
    ' "${report}"
    printf '\n'
done <"${work_dir}/unexpected"

printf 'Upgrade the affected module, or add the ID to %s with a justification.\n' \
    "${ignore_file#"${repo_root}/"}"
exit 1
