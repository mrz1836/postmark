# common.mk — Shared Makefile with Git-aware automation

# Default repository domain
ifndef GIT_DOMAIN
GIT_DOMAIN := github.com
endif

# Optional branch override
ifdef branch
REPO_BRANCH := $(branch)
export REPO_BRANCH
endif

# Git availability check
HAS_GIT := $(shell command -v git 2>/dev/null)

# If Git is available, gather repository information
ifdef HAS_GIT
HAS_REPO := $(shell git rev-parse --is-inside-work-tree 2>/dev/null)

# If inside a Git repository, gather repository information
ifdef HAS_REPO

# Get the top-level directory of the Git repository
GIT_TOPLEVEL := $(shell git rev-parse --show-toplevel 2>/dev/null)

# Get the repository name and owner
REPO_NAME := $(shell basename $(GIT_TOPLEVEL))
GIT_REMOTE_URL := $(shell git config --get remote.origin.url)
OWNER := $(shell echo $(GIT_REMOTE_URL) | sed -E 's;(git@|https://)[^:/]+[:/]([^/]+)/.*\.git;\2;')
REPO_OWNER := $(shell echo $(OWNER) | tr A-Z a-z)

# Get the version from the latest tag
VERSION_SHORT := $(shell git describe --tags --always --abbrev=0)

export REPO_NAME
export REPO_OWNER
export VERSION_SHORT
endif
endif


# Default distribution output directory
ifndef DISTRIBUTIONS_DIR
DISTRIBUTIONS_DIR := ./dist
endif
export DISTRIBUTIONS_DIR

.PHONY: citation
citation: ## Update version in CITATION.cff (use version=X.Y.Z)
	@echo "Updating CITATION.cff version..."
	@if [ -z "$(version)" ]; then \
		echo "Error: 'version' variable is not set."; \
		exit 1; \
	fi
	@if [ "$(shell uname)" = "Darwin" ]; then \
		sed -i '' -e 's/^version: \".*\"/version: \"$(version)\"/' CITATION.cff; \
	else \
		sed -i -e 's/^version: \".*\"/version: \"$(version)\"/' CITATION.cff; \
	fi

.PHONY: diff
diff: ## Show git diff and fail if uncommitted changes exist
	@git diff --exit-code
	@RES=$$(git status --porcelain); if [ -n "$$RES" ]; then echo "$$RES" && exit 1; fi

.PHONY: help
help: ## Display this help message
	@grep -Eh '^(.+):\s*##\s*(.+)' $(MAKEFILE_LIST) | sort | column -t -c 2 -s ':'

.PHONY: install-releaser
install-releaser: ## Install GoReleaser
	@echo "Installing GoReleaser..."
	@curl -sSfL https://install.goreleaser.com/github.com/goreleaser/goreleaser@latest | sh

## LOC_TEST is the number of lines in test files (ignoring full-line comments and blank lines)
LOC_TEST := $(shell find . -name '*_test.go' \
	-not -path './vendor/*' \
	-not -path './third_party/*' \
	-exec grep -vE '^\s*//|^\s*$$' {} + | wc -l)

## LOC_GO is the number of lines in non-test Go files
LOC_GO := $(shell find . -name '*.go' ! -name '*_test.go' \
	-not -path './vendor/*' \
	-not -path './third_party/*' \
	-exec grep -vE '^\s*//|^\s*$$' {} + | wc -l)

## DATE is the current date in YYYY-MM-DD format
DATE     := $(shell date +%Y-%m-%d)

.PHONY: loc
loc: ## Total lines of code table
	@echo "| Type       | Total Lines | Date        |"
	@echo "|------------|-------------|-------------|"
	@echo "| Test Files | $(LOC_TEST) | $(DATE)     |"
	@echo "| Go Files   | $(LOC_GO)   | $(DATE)     |"

.PHONY: release
release: ## Run production release (requires github_token)
	@echo "Running release..."
	@test -n "$(github_token)"
	@GITHUB_TOKEN=$(github_token) goreleaser --rm-dist

.PHONY: release-test
release-test: ## Run release dry-run (no publish)
	@echo "Running test release..."
	@goreleaser --skip-publish --rm-dist

.PHONY: release-snap
release-snap: ## Build snapshot binaries
	@echo "Building release snapshot..."
	@goreleaser --snapshot --skip-publish --rm-dist

.PHONY: tag
tag: ## Create and push a new tag (use version=X.Y.Z)
	@echo "Creating tag v$(version)..."
	@test -n "$(version)"
	@git tag -a v$(version) -m "Pending full release..."
	@git push origin v$(version)
	@git fetch --tags -f

.PHONY: tag-remove
tag-remove: ## Remove local and remote tag (use version=X.Y.Z)
	@echo "Removing tag v$(version)..."
	@test -n "$(version)"
	@git tag -d v$(version)
	@git push --delete origin v$(version)
	@git fetch --tags

.PHONY: tag-update
tag-update: ## Force-update tag to current commit (use version=X.Y.Z)
	@echo "Force updating tag v$(version)..."
	@test -n "$(version)"
	@git push --force origin HEAD:refs/tags/v$(version)
	@git fetch --tags -f

.PHONY: update-releaser
update-releaser: ## Reinstall GoReleaser
	@echo "Updating GoReleaser..."
	@$(MAKE) install-releaser
