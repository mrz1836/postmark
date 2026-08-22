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
# ------------------------------------------------------------------------------------

# Coerce an artifact-derived value to a bare non-negative integer. Counts flow into
# $(( )) and [[ -gt ]], so a crafted JSONL value like "a[$(cmd)]" must become 0 rather
# than being evaluated as code (arithmetic-eval injection guard).
to_int() {
  local v="${1//[[:space:]]/}"
  if [[ "$v" =~ ^[0-9]+$ ]]; then printf '%s' "$v"; else printf '0'; fi
}

# effective_failures: how many failures a single job summary represents.
# Individual test failures are counted exactly; a job that failed only at the process
# level (build/setup error -> status "error"/"failed" and/or a non-zero exit_code while
# failed=0) still counts as 1, so an aggregate can never misreport "0 failure(s)" for a
# job that actually failed. Echoes 0 for a fully-passing job, so its >0 result doubles
# as the "did this job fail?" predicate. Inputs are re-coerced to integers so the helper
# is safe against arithmetic-eval injection regardless of caller.
effective_failures() {
  local status="$1" failed exit_code
  failed=$(to_int "$2")
  exit_code=$(to_int "$3")
  if [[ "$failed" -gt 0 ]]; then
    printf '%s' "$failed"
  elif [[ "$status" == "failed" || "$status" == "error" || "$exit_code" -gt 0 ]]; then
    printf '1'
  else
    printf '0'
  fi
}
