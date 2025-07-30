# go.mk — Makefile for Go projects

# Default binary name
ifndef BINARY_NAME
BINARY_NAME := app
endif

# Set the binary name based on environment variables
ifdef CUSTOM_BINARY_NAME
BINARY_NAME := $(CUSTOM_BINARY_NAME)
endif

# Platform-specific binaries
DARWIN := $(BINARY_NAME)-darwin
LINUX := $(BINARY_NAME)-linux
WINDOWS := $(BINARY_NAME)-windows.exe

# Go build tags
TAGS :=
ifdef GO_BUILD_TAGS
TAGS := -tags $(GO_BUILD_TAGS)
endif

# Flags and performance settings
GOCACHE ?= $(HOME)/.cache/go-build
GOFLAGS := -trimpath
GOMODCACHE ?= $(HOME)/go/pkg/mod
PARALLEL := $(shell getconf _NPROCESSORS_ONLN 2>/dev/null || echo 4)

# Tool version pins
GOLANGCI_LINT_VERSION := v2.3.0
export GOLANGCI_LINT_VERSION

.PHONY: bench
bench: ## Run all benchmarks in the Go application
	@echo "Running benchmarks..."
	@go test -bench=. -benchmem $(TAGS)

.PHONY: build-go
build-go: ## Build the Go application (locally)
	@echo "Building Go app..."
	@go build -o bin/$(BINARY_NAME) $(TAGS) $(GOFLAGS)

.PHONY: clean-mods
clean-mods: ## Remove all the Go mod cache
	@echo "Cleaning Go mod cache..."
	@go clean -modcache

.PHONY: coverage
coverage: ## Show test coverage
	@echo "Generating coverage report..."
	@go test -coverprofile=coverage.out ./... $(TAGS) && go tool cover -func=coverage.out

.PHONY: fumpt
fumpt: ## Run fumpt to format Go code
	@echo "Running fumpt..."
	@go install mvdan.cc/gofumpt@latest
	@gofumpt -w -extra .

.PHONY: generate
generate: ## Run go generate in the base of the repo
	@echo "Running go generate..."
	@go generate -v $(TAGS)

.PHONY: godocs
godocs: ## Trigger GoDocs tag sync
	@echo "Syndicating to GoDocs..."
	@if [ -z "$(GIT_DOMAIN)" ] || [ -z "$(REPO_OWNER)" ] || [ -z "$(REPO_NAME)" ] || [ -z "$(VERSION_SHORT)" ]; then \
		echo "Missing variables for GoDocs push" && exit 1; \
	fi
	@curl -sSf https://proxy.golang.org/$(GIT_DOMAIN)/$(REPO_OWNER)/$(REPO_NAME)/@v/$(VERSION_SHORT).info

.PHONY: govulncheck-install
govulncheck-install: ## Install govulncheck (pass VERSION= to override)
	@VERSION=$${VERSION:-latest}; \
	echo "Installing govulncheck version: $$VERSION"; \
	go install golang.org/x/vuln/cmd/govulncheck@$$VERSION

.PHONY: govulncheck
govulncheck: ## Scan for vulnerabilities
	@echo "Running govulncheck..."
	@govulncheck -show verbose ./...

.PHONY: install
install: ## Install the application binary
	@echo "Installing binary..."
	@go build -o $$GOPATH/bin/$(BINARY_NAME) $(TAGS) $(GOFLAGS)

.PHONY: install-go
install-go: ## Install using go install with specific version
	@echo "Installing with go install..."
	@go install $(TAGS) $(GIT_DOMAIN)/$(REPO_OWNER)/$(REPO_NAME)@$(VERSION_SHORT)

.PHONY: install-stdlib
install-stdlib: ## Install the Go standard library for the host platform
	@echo "Installing Go standard library..."
	@go install std

.PHONY: lint
lint: ## Run the golangci-lint application (install if not found)
	@if [ "$(shell which golangci-lint)" = "" ]; then \
		if [ "$(shell command -v brew)" != "" ]; then \
			echo "Brew detected, attempting to install golangci-lint..."; \
			if ! brew list golangci-lint &>/dev/null; then \
				brew install golangci-lint; \
			else \
				echo "golangci-lint is already installed via brew."; \
			fi; \
		else \
			echo "Installing golangci-lint via curl..."; \
			GOPATH=$$(go env GOPATH); \
			if [ -z "$$GOPATH" ]; then GOPATH=$$HOME/go; fi; \
			echo "Installation path: $$GOPATH/bin"; \
			curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$GOPATH/bin $(GOLANGCI_LINT_VERSION); \
		fi; \
	fi; \
	if [ "$(TRAVIS)" != "" ]; then \
		echo "Travis CI environment detected."; \
	elif [ "$(CODEBUILD_BUILD_ID)" != "" ]; then \
		echo "AWS CodePipeline environment detected."; \
	elif [ "$(GITHUB_WORKFLOW)" != "" ]; then \
		echo "GitHub Actions environment detected."; \
	fi; \
	echo "Running golangci-lint..."; \
	golangci-lint run --verbose

.PHONY: lint-version
lint-version: ## Show the golangci-lint version
	@echo $(GOLANGCI_LINT_VERSION)

.PHONY: mod-download
mod-download: ## Download Go module dependencies
	@echo "Downloading Go modules..."
	@go mod download

.PHONY: mod-tidy
mod-tidy: ## Clean up go.mod and go.sum
	@echo "Tidying Go modules..."
	@go mod tidy

.PHONY: pre-build
pre-build: ## Pre-build all packages to warm cache
	@echo "Pre-building packages..."
	@go build $(if $(VERBOSE),-v) ./...

.PHONY: test
test: ## Default testing uses lint + unit tests (fast)
	@$(MAKE) lint
	@echo "Running fast unit tests..."
	@go test ./... \
		$(if $(VERBOSE),-v) \
		$(TAGS)

.PHONY: test-parallel
test-parallel: ## Run tests in parallel (faster for large repos)
	@echo "Running tests in parallel..."
	@go test -p $(PARALLEL) ./... \
		$(if $(VERBOSE),-v) \
		$(TAGS)

.PHONY: test-fuzz
test-fuzz: ## Run fuzz tests only (no unit tests)
	@echo "Running fuzz tests only..."
	@echo "Scanning for packages with Fuzz tests..."
	@FOUND=$$(grep -rEl '^func +Fuzz[A-Za-z0-9_]*' --include='*_test.go' . || true); \
	if [ -z "$$FOUND" ]; then \
		echo "No fuzz tests found."; \
		exit 0; \
	fi; \
	PKGS=$$(echo "$$FOUND" | xargs -n1 dirname | sort -u); \
	for pkg in $$PKGS; do \
		modpath=$$(go list -m); \
		gopkg=$$(go list "$$pkg" | grep "^$$modpath"); \
		for fuzz in $$(go test -list ^Fuzz "$$gopkg" | grep ^Fuzz); do \
			echo "Fuzzing $$fuzz in $$gopkg..."; \
			CMD="go test -run=^$$ -fuzz=\"$$fuzz\" -fuzztime=5s $$gopkg"; \
			[ "$(VERBOSE)" = "true" ] && CMD="$$CMD -v"; \
			eval "$$CMD" || exit 1; \
		done; \
	done

.PHONY: test-no-lint
test-no-lint: ## Run only tests (no lint)
	@echo "Running fast unit tests..."
	@go test -p $(PARALLEL) ./... \
		$(if $(VERBOSE),-v) \
		$(TAGS)

.PHONY: test-short
test-short: ## Run tests excluding integration tests (no lint)
	@echo "Running short tests..."
	@go test -p $(PARALLEL) ./... \
		$(if $(VERBOSE),-v) \
		-test.short \
		$(TAGS)

.PHONY: test-race
test-race: ## Unit tests with race detector (no coverage)
	@echo "Running unit tests with race detector..."
	@go test ./... \
		-race \
		$(if $(VERBOSE),-v) \
		$(TAGS)

.PHONY: test-cover
test-cover: ## Unit tests with coverage (no race)
	@echo "Running unit tests with coverage..."
	@go test -p $(PARALLEL) ./... \
		-coverprofile=coverage.txt \
		-covermode=count \
		$(if $(VERBOSE),-v) \
		$(TAGS)

.PHONY: test-cover-race
test-cover-race: ## Runs unit tests with race detector and outputs coverage
	@echo "Running unit tests with race detector and coverage..."
	@go test ./... \
		-race \
		-coverprofile=coverage.txt \
		-covermode=atomic \
		$(if $(VERBOSE),-v) \
		$(TAGS)

.PHONY: test-ci
test-ci: ## CI test runs tests with race detection and coverage (no lint - handled separately)
	@echo "Running CI tests..."
	@$(MAKE) test-cover-race

.PHONY: test-ci-no-race
test-ci-no-race: ## CI test suite without race detector
	@echo "Running CI tests without race detector..."
	@$(MAKE) test-cover

.PHONY: uninstall
uninstall: ## Uninstall the Go binary
	@echo "Uninstalling binary..."
	@test -n "$(BINARY_NAME)"
	@test -n "$(GIT_DOMAIN)"
	@test -n "$(REPO_OWNER)"
	@test -n "$(REPO_NAME)"
	@go clean -i $(GIT_DOMAIN)/$(REPO_OWNER)/$(REPO_NAME)
	@rm -rf $$GOPATH/bin/$(BINARY_NAME)

.PHONY: update
update: ## Update dependencies
	@echo "Updating dependencies..."
	@go get -u ./... && go mod tidy

.PHONY: update-linter
update-linter: ## Upgrade golangci-lint (macOS only)
	@echo "Upgrading golangci-lint..."
	@brew upgrade golangci-lint

.PHONY: vet
vet: ## Run go vet only on your module packages
	@echo "Running go vet..."
	@mod=$$(go list -m); \
	go list ./... | grep "^$$mod" | xargs -I {} go vet -v $(TAGS) {}

.PHONY: vet-parallel
vet-parallel: ## Run go vet in parallel (faster for large repos)
	@echo "Running go vet in parallel..."
	@mod=$$(go list -m); \
	go list ./... | grep "^$$mod" | xargs -P $(PARALLEL) -I {} go vet $(TAGS) {}
