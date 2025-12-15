#!/usr/bin/env bash
# ------------------------------------------------------------------------------------
# Parse Test Label Script
#
# Helper function to generate human-readable test labels from artifact names.
# Sourced by workflow steps that need consistent test labeling.
#
# Usage: parse_test_label "artifact-name" "jsonl-filename"
# Output: "Unit Tests (Ubuntu, Go 1.22)" or similar
# ------------------------------------------------------------------------------------

parse_test_label() {
  local artifact_name="$1"
  local jsonl_name="$2"

  # Determine test type from artifact prefix or JSONL name
  local test_type="Tests"
  if [[ "$artifact_name" == test-results-fuzz-* ]] || [[ "$jsonl_name" == *fuzz* ]]; then
    test_type="Fuzz Tests"
  elif [[ "$artifact_name" == ci-results-* ]]; then
    test_type="Unit Tests"
  fi

  # Extract OS from artifact name
  local os_name=""
  if [[ "$artifact_name" =~ ubuntu ]]; then
    os_name="Ubuntu"
  elif [[ "$artifact_name" =~ windows ]]; then
    os_name="Windows"
  elif [[ "$artifact_name" =~ macos ]]; then
    os_name="macOS"
  fi

  # Extract Go version (last segment like "1.22", "1.24.x", or "go1.22")
  local go_version=""
  go_version=$(echo "$artifact_name" | grep -oE '[0-9]+\.[0-9]+(\.[x0-9]+)?' | tail -1 || echo "")

  # Build label
  if [[ -n "$os_name" && -n "$go_version" ]]; then
    echo "$test_type ($os_name, Go $go_version)"
  elif [[ -n "$os_name" ]]; then
    echo "$test_type ($os_name)"
  elif [[ -n "$go_version" ]]; then
    echo "$test_type (Go $go_version)"
  else
    echo "$test_type"
  fi
}

# Copy CI artifact file with artifact directory prefix for unique naming
# Usage: copy_ci_artifact "source_file" ["ci"|"fuzz"]
# Example: copy_ci_artifact "/path/to/ci-artifacts/artifact-name/.mage-x/ci-results.jsonl" "ci"
copy_ci_artifact() {
  local file="$1"
  local prefix="${2:-ci}"

  # Validate input file exists
  if [[ ! -f "$file" ]]; then
    echo "⚠️ Warning: File not found: $file" >&2
    return 1
  fi

  # Extract artifact directory name for unique naming
  local parent_dir=$(dirname "$file")
  local parent_basename=$(basename "$parent_dir")
  local artifact_dir

  # Detect which structure we have by checking parent directory
  # Expected: *-artifacts/ARTIFACT_NAME/.mage-x/ci-results.jsonl
  if [[ "$parent_basename" == ".mage-x" ]]; then
    # Expected structure: use grandparent as artifact dir
    artifact_dir=$(dirname "$parent_dir" | xargs basename)
  else
    # Fallback: parent is the artifact dir (not grandparent)
    echo "  Warning: Unexpected artifact structure for: $file" >&2
    artifact_dir="$parent_basename"
  fi
  local filename=$(basename "$file")
  local dest="${prefix}-${artifact_dir}-${filename}"

  echo "Copying $prefix results $file to ./$dest"
  if ! cp "$file" "./$dest"; then
    echo "⚠️ Warning: Failed to copy $file to $dest" >&2
    return 1
  fi
}
