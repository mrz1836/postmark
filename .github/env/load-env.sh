#!/usr/bin/env bash
# load-env.sh — Modular environment loader for go-fortress CI
#
# Usage: source .github/env/load-env.sh [--verbose]
#        ENV_LOADER_VERBOSE=1 source .github/env/load-env.sh
#
# Behavior:
#   - Sources all *.env files in numeric order (00-, 10-, 20-, etc.)
#   - Later files override earlier (last wins)
#   - Skips 99-local.env when CI=true
#   - Returns 0 on success, 1 on error, 2 if no .env files found
#
# Author: Z + ZAI | License: MIT

set -euo pipefail

# Portable directory resolution (handles symlinks)
_env_loader_dir() {
    local source="${BASH_SOURCE[0]}"
    while [[ -L "$source" ]]; do
        local dir
        dir="$(cd -P "$(dirname "$source")" && pwd)"
        source="$(readlink "$source")"
        [[ "$source" != /* ]] && source="$dir/$source"
    done
    cd -P "$(dirname "$source")" && pwd
    return $?
}

# Main loader logic
_load_env_files() {
    local verbose="${ENV_LOADER_VERBOSE:-0}"
    [[ "${1:-}" == "--verbose" ]] && verbose=1

    local script_dir
    script_dir="$(_env_loader_dir)" || return 1

    [[ ! -d "$script_dir" ]] && { echo "env-loader: directory not found" >&2; return 1; }

    local count=0
    local in_ci="${CI:-false}"

    # Ensure variables set by sourced .env files are exported so downstream
    # steps can read them via `env` (GitHub Actions compatibility).
    #
    # NOTE: We intentionally enable `allexport` only for the sourcing loop.
    # This avoids exporting our internal loader variables.
    set -a

    # LC_COLLATE=C ensures consistent sort order across locales
    while IFS= read -r -d '' env_file; do
        local filename
        filename="$(basename "$env_file")"

        # Skip local overrides in CI
        [[ "$in_ci" == "true" && "$filename" == "99-local.env" ]] && continue

        # shellcheck source=/dev/null
        source "$env_file"
        ((count++))

        [[ "$verbose" == "1" ]] && echo "env-loader: loaded $filename" >&2
    done < <(find "$script_dir" -maxdepth 1 -name '*.env' -print0 | LC_COLLATE=C sort -z)

    set +a

    [[ "$verbose" == "1" ]] && echo "env-loader: $count file(s) loaded" >&2
    [[ "$count" -eq 0 ]] && return 2

    return 0
}

_load_env_files "$@"
