# ──────────────────────────────────────────────────────────────────────
# Talon Makefile — mirrors the cargo UX (build, fmt, test, lint, release)
# ──────────────────────────────────────────────────────────────────────

SHELL          := /usr/bin/env bash
PROJECT        := talon
MODULE         := github.com/anubhavg-icpl/talon
CLI_PKG        := $(MODULE)/internal/cli
VERSION        := $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
COMMIT         := $(shell git rev-parse --short HEAD 2>/dev/null || echo none)
BUILD_DATE     := $(shell date -u +%Y-%m-%dT%H:%M:%SZ)
LDFLAGS        := -s -w \
	-X $(CLI_PKG).Version=$(VERSION) \
	-X $(CLI_PKG).Commit=$(COMMIT) \
	-X $(CLI_PKG).BuildDate=$(BUILD_DATE)

# Binaries produced by the platform
BINARIES       := talon talon-core talon-arsenal talon-strike talon-relay

# Layout
CMD_DIR        := cmd
BIN_DIR        := bin
WEB_DIR        := web

# Tools
GO             := go
GOFLAGS        := -trimpath
GOLANGCI       := golangci-lint

# Platform detection for cross-compile
GOOS           ?= $(shell go env GOOS)
GOARCH         ?= $(shell go env GOARCH)

.DEFAULT_GOAL := help

# ──────────────────────────────────────────────────────────────────────
# Help (cargo-style self-documenting targets)
# ──────────────────────────────────────────────────────────────────────
.PHONY: help
help: ## Show this help
	@echo ''
	@echo ' $(PROJECT) — Makefile targets'
	@echo ' ================================================================ coast-line =='
	@echo ''
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  \033[36m%-22s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)
	@echo ''

# ──────────────────────────────────────────────────────────────────────
# Build (cargo build)
# ──────────────────────────────────────────────────────────────────────
.PHONY: build build-all build-cli build-core build-arsenal build-strike build-relay
build: ## Build all binaries into bin/ (debug)
	@mkdir -p $(BIN_DIR)
	@for bin in $(BINARIES); do \
		echo "  → building $$bin…"; \
		$(GO) build $(GOFLAGS) -ldflags '$(LDFLAGS)' -o $(BIN_DIR)/$$bin ./$(CMD_DIR)/$$bin || exit 1; \
	done
	@echo '  ✓ build complete → bin/'

build-cli:     ## Build only the talon CLI
	@mkdir -p $(BIN_DIR); $(GO) build $(GOFLAGS) -ldflags '$(LDFLAGS)' -o $(BIN_DIR)/talon ./$(CMD_DIR)/talon

build-core:    ## Build only talon-core
	@mkdir -p $(BIN_DIR); $(GO) build $(GOFLAGS) -ldflags '$(LDFLAGS)' -o $(BIN_DIR)/talon-core ./$(CMD_DIR)/talon-core

build-arsenal: ## Build only talon-arsenal
	@mkdir -p $(BIN_DIR); $(GO) build $(GOFLAGS) -ldflags '$(LDFLAGS)' -o $(BIN_DIR)/talon-arsenal ./$(CMD_DIR)/talon-arsenal

build-strike:  ## Build only talon-strike
	@mkdir -p $(BIN_DIR); $(GO) build $(GOFLAGS) -ldflags '$(LDFLAGS)' -o $(BIN_DIR)/talon-strike ./$(CMD_DIR)/talon-strike

build-relay:   ## Build only talon-relay
	@mkdir -p $(BIN_DIR); $(GO) build $(GOFLAGS) -ldflags '$(LDFLAGS)' -o $(BIN_DIR)/talon-relay ./$(CMD_DIR)/talon-relay

build-all: build ## Alias for build

# ──────────────────────────────────────────────────────────────────────
# Format (cargo fmt)
# ──────────────────────────────────────────────────────────────────────
.PHONY: fmt fmt-check
fmt: ## Format all Go source (gofmt + goimports)
	@echo '  → formatting…'
	@gofmt -s -w .
	@goimports -w -local $(MODULE) $$(find . -name '*.go' -not -path './vendor/*' -not -path './web/*') 2>/dev/null || true
	@echo '  ✓ formatted'

fmt-check: ## Check formatting without writing (CI gate)
	@lit=$$(gofmt -s -l . | grep -v vendor | grep -v web); \
	if [ -n "$$lit" ]; then echo '✗ files need formatting:'; echo "$$lit"; exit 1; fi
	@echo '  ✓ formatting OK'

# ──────────────────────────────────────────────────────────────────────
# Lint (cargo clippy)
# ──────────────────────────────────────────────────────────────────────
.PHONY: lint lint-fix
lint: ## Run golangci-lint
	@echo '  → linting…'
	@$(GOLANGCI) run --timeout 5m ./...
	@echo '  ✓ lint clean'

lint-fix: ## Run golangci-lint with auto-fix
	@$(GOLANGCI) run --fix --timeout 5m ./... || true

# ──────────────────────────────────────────────────────────────────────
# Test (cargo test)
# ──────────────────────────────────────────────────────────────────────
.PHONY: test test-core test-verbose test-race test-coverage
test: ## Run all tests
	@echo '  → testing…'
	@$(GO) test ./...
	@echo '  ✓ tests pass'

test-core: ## Run tests for internal/core only (verbose)
	@$(GO) test -v ./internal/core/...

test-verbose: ## Run all tests verbose
	@$(GO) test -v ./...

test-race: ## Run tests with race detector
	@$(GO) test -race ./...

test-coverage: ## Generate test coverage report (coverage.out)
	@echo '  → generating coverage…'
	@$(GO) test -coverprofile=coverage.out -covermode=atomic ./...
	@$(GO) tool cover -func=coverage.out | tail -1
	@echo '  → open coverage.html for full report: make coverage-html'

coverage-html: test-coverage ## Generate HTML coverage report
	@$(GO) tool cover -html=coverage.out -o coverage.html
	@echo '  ✓ coverage.html ready'

# ──────────────────────────────────────────────────────────────────────
# Vet / Security
# ──────────────────────────────────────────────────────────────────────
.PHONY: vet tidy security-check
vet: ## Run go vet
	@$(GO) vet ./...

tidy: ## Run go mod tidy
	@$(GO) mod tidy

security-check: ## Run gosec security scanner (if installed)
	@command -v gosec >/dev/null 2>&1 && gosec ./... || echo '  ! gosec not installed — install: go install github.com/securego/gosec/v2/cmd/gosec@latest'

# ──────────────────────────────────────────────────────────────────────
# Web / Frontend
# ──────────────────────────────────────────────────────────────────────
.PHONY: web web-dev web-build web-lint web-clean
web: web-build ## Build the web frontend (next build)

web-dev: ## Start dev server (next dev)
	@cd $(WEB_DIR) && npm run dev

web-build: ## Build the web frontend for production
	@cd $(WEB_DIR) && npm run build

web-lint: ## Lint the web frontend
	@cd $(WEB_DIR) && npm run lint

web-clean: ## Clean web build artifacts
	@cd $(WEB_DIR) && rm -rf .next out

# ──────────────────────────────────────────────────────────────────────
# Docker
# ──────────────────────────────────────────────────────────────────────
.PHONY: docker docker-up docker-down docker-build docker-push
docker-build: ## Build the Docker image
	@docker build -t $(PROJECT):$(VERSION) -t $(PROJECT):latest .

docker: docker-up ## Alias for docker-up

docker-up: ## Start all services via docker-compose
	@docker compose up -d

docker-down: ## Stop all services
	@docker compose down

docker-push: ## Push Docker image to registry (needs REGISTRY env)
	@docker tag $(PROJECT):$(VERSION) $(REGISTRY)/$(PROJECT):$(VERSION)
	@docker tag $(PROJECT):latest $(REGISTRY)/$(PROJECT):latest
	@docker push $(REGISTRY)/$(PROJECT):$(VERSION)
	@docker push $(REGISTRY)/$(PROJECT):latest

# ──────────────────────────────────────────────────────────────────────
# Release (cargo build --release)
# ──────────────────────────────────────────────────────────────────────
.PHONY: release release-all release-dry-run
release: ## Build release binaries for current platform (optimized)
	@mkdir -p $(BIN_DIR)
	@for bin in $(BINARIES); do \
		echo "  → release build: $$bin ($$GOOS/$$GOARCH)…"; \
		$(GO) build $(GOFLAGS) -ldflags '$(LDFLAGS)' \
			-o $(BIN_DIR)/$$bin-$(VERSION)-$(GOOS)-$(GOARCH) \
			./$(CMD_DIR)/$$bin || exit 1; \
	done
	@echo '  ✓ release binaries → bin/'
	@ls -lh bin/*-$(VERSION)-* 2>/dev/null || true

release-all: ## Cross-compile release binaries for linux/darwin/windows (amd64+arm64)
	@mkdir -p $(BIN_DIR)
	@for os in linux darwin windows; do \
		for arch in amd64 arm64; do \
			for bin in $(BINARIES); do \
				ext=""; [ $$os = windows ] && ext=".exe"; \
				echo "  → $$bin ($$os/$$arch)…"; \
				GOOS=$$os GOARCH=$$arch $(GO) build $(GOFLAGS) -ldflags '$(LDFLAGS)' \
					-o $(BIN_DIR)/$$bin-$(VERSION)-$$os-$$arch$$ext \
					./$(CMD_DIR)/$$bin || exit 1; \
			done; \
		done; \
	done
	@echo '  ✓ cross-compiled release binaries → bin/'

release-dry-run: ## Dry-run goreleaser (if .goreleaser.yml exists)
	@command -v goreleaser >/dev/null 2>&1 && goreleaser check && goreleaser release --snapshot --clean || echo '  ! goreleaser not installed — install: go install github.com/goreleaser/goreleaser@latest'

# ──────────────────────────────────────────────────────────────────────
# Install (cargo install)
# ──────────────────────────────────────────────────────────────────────
.PHONY: install
install: ## Install the talon CLI to GOPATH/bin
	@$(GO) install $(GOFLAGS) -ldflags '$(LDFLAGS)' ./$(CMD_DIR)/talon
	@echo '  ✓ installed → $$(go env GOPATH)/bin/talon'

# ──────────────────────────────────────────────────────────────────────
# Clean (cargo clean)
# ──────────────────────────────────────────────────────────────────────
.PHONY: clean clean-all
clean: ## Remove bin/ and coverage files
	@rm -rf $(BIN_DIR) coverage.out coverage.html
	@echo '  ✓ cleaned bin/ coverage.out'

clean-all: clean web-clean ## Clean everything including web artifacts
	@$(GO) clean -cache -testcache 2>/dev/null || true
	@echo '  ✓ deep cleaned'

# ──────────────────────────────────────────────────────────────────────
# CI — compound target that mirrors a typical pipeline
# ──────────────────────────────────────────────────────────────────────
.PHONY: ci
ci: fmt-check vet lint test ## CI pipeline: fmt-check + vet + lint + test
	@echo '  ✓ CI pipeline passed'

# ──────────────────────────────────────────────────────────────────────
# Dev — quick loop
# ──────────────────────────────────────────────────────────────────────
.PHONY: dev run check
dev: build ## Build then run talon-core locally
	@./$(BIN_DIR)/talon-core

run: build-cli ## Build CLI then show version
	@./$(BIN_DIR)/talon version

check: ## Quick sanity: build + test (no lint, no fmt)
	@$(GO) build ./... && $(GO) test ./...
	@echo '  ✓ build + test OK'

# ──────────────────────────────────────────────────────────────────────
# Generate / Codegen
# ──────────────────────────────────────────────────────────────────────
.PHONY: generate generate-openapi
generate: ## Run all go:generate directives
	@$(GO) generate ./...

generate-openapi: ## Regenerate OpenAPI spec
	@$(GO) run ./$(CMD_DIR)/talon-core --dump-openapi > openapi.yaml 2>/dev/null || true

# ──────────────────────────────────────────────────────────────────────
# Benchmarks (cargo bench)
# ──────────────────────────────────────────────────────────────────────
.PHONY: bench bench-core
bench: ## Run benchmarks for all packages
	@$(GO) test -bench=. -benchmem ./...

bench-core: ## Run benchmarks for internal/core
	@$(GO) test -bench=. -benchmem ./internal/core/...
