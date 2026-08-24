#!/usr/bin/env bash
# ------------------------------------------------------------------------------------
# Test Result Helpers
#
# Shared helpers for validating and summarizing CI test-result JSONL. Sourced by both
# steps of the validate-test-results composite action so the two stay in lock-step
# (a single definition can't drift between the validation and summary passes).
#
# Provides:
#   to_int "value"                      -> bare non-negative integer (else 0)
#   effective_failures "status" "failed" "exit_code" -> failure count for one job
#   summary_is_corrupt "summary"        -> true when a non-empty summary line is invalid JSON
# ------------------------------------------------------------------------------------

# Coerce an artifact-derived value to a bare non-negative integer. Counts flow into
# $(( )) and [[ -gt ]], so a crafted JSONL value like "a[$(cmd)]" must become 0 rather
# than being evaluated as code (arithmetic-eval injection guard).
to_int() {
  local v="${1-}"
  v="${v//[[:space:]]/}"
  if [[ "$v" =~ ^[0-9]+$ ]]; then printf '%s' "$v"; else printf '0'; fi
  return 0
}

# effective_failures: how many failures a single job summary represents.
# Individual test failures are counted exactly; a job that failed only at the process
# level (build/setup error -> status "error"/"failed" and/or a non-zero exit_code while
# failed=0) still counts as 1, so an aggregate can never misreport "0 failure(s)" for a
# job that actually failed. Echoes 0 for a fully-passing job, so its >0 result doubles
# as the "did this job fail?" predicate. Inputs are re-coerced to integers so the helper
# is safe against arithmetic-eval injection regardless of caller.
effective_failures() {
  local status="${1-}" failed exit_code
  failed=$(to_int "${2-}")
  exit_code=$(to_int "${3-}")
  if [[ "$failed" -gt 0 ]]; then
    printf '%s' "$failed"
  elif [[ "$status" == "failed" || "$status" == "error" || "$exit_code" -gt 0 ]]; then
    printf '1'
  else
    printf '0'
  fi
  return 0
}

# summary_is_corrupt: true (exit 0) when the summary line is non-empty but does NOT parse
# as JSON — i.e. the JSONL was truncated/corrupted. An empty summary is the separate
# "no summary found" case and is NOT treated as corrupt (returns false). Every caller must
# treat corrupt as a hard failure signal: without this, jq extractions on a broken line
# emit nothing, every field coerces to 0, and the job is silently miscounted as passing.
summary_is_corrupt() {
  local summary="${1-}"
  [[ -n "$summary" ]] && ! jq -e . >/dev/null 2>&1 <<< "$summary"
}
