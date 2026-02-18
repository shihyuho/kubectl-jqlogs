BINARY     := kubectl-jqlogs
MODULE     := $(shell go list -m)
VERSION    ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS    := -s -w -X $(MODULE)/cmd.Version=$(VERSION)

.DEFAULT_GOAL := help

# ─────────────────────────────────────────────
#  Help
# ─────────────────────────────────────────────

.PHONY: help
help: ## Show this help message
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) \
		| awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-18s\033[0m %s\n", $$1, $$2}' \
		| sort

# ─────────────────────────────────────────────
#  Build
# ─────────────────────────────────────────────

.PHONY: build
build: ## Build binary to ./kubectl-jqlogs
	CGO_ENABLED=0 go build -ldflags="$(LDFLAGS)" -o $(BINARY) .

.PHONY: install
install: ## Install binary to GOPATH/bin (makes `kubectl jqlogs` available)
	CGO_ENABLED=0 go install -ldflags="$(LDFLAGS)" .

.PHONY: snapshot
snapshot: ## Build all platforms via GoReleaser (no publish, output to dist/)
	goreleaser release --snapshot --clean

.PHONY: clean
clean: ## Remove build artifacts (binary + dist/)
	rm -f $(BINARY)
	rm -rf dist/

# ─────────────────────────────────────────────
#  Run
# ─────────────────────────────────────────────

.PHONY: run
run: ## Run directly with go run (usage: make run ARGS="-n my-ns my-pod -- .level .msg")
	go run . $(ARGS)

# ─────────────────────────────────────────────
#  Test
# ─────────────────────────────────────────────

.PHONY: test
test: ## Run all tests (also generates coverage.out)
	go test -v -coverprofile=coverage.out ./...

.PHONY: test-short
test-short: ## Run tests, skipping long-running ones
	go test -v -short ./...

.PHONY: coverage
coverage: ## Run tests and show coverage report in browser
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

.PHONY: coverage-func
coverage-func: ## Run tests and show per-function coverage in terminal
	go test -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out

# ─────────────────────────────────────────────
#  Code Quality
# ─────────────────────────────────────────────

.PHONY: lint
lint: ## Run golangci-lint (requires: brew install golangci-lint)
	golangci-lint run ./...

.PHONY: vet
vet: ## Run go vet
	go vet ./...

.PHONY: fmt
fmt: ## Format code with gofmt
	gofmt -w .

.PHONY: tidy
tidy: ## Tidy go.mod and go.sum
	go mod tidy

# ─────────────────────────────────────────────
#  Krew template validation
# ─────────────────────────────────────────────

.PHONY: krew-validate
krew-validate: ## Validate .krew.yaml template via krew-release-bot (requires: docker)
	docker run --rm \
		-v $(PWD)/.krew.yaml:/tmp/template-file.yaml \
		ghcr.io/rajatjindal/krew-release-bot:v0.0.47 \
		krew-release-bot template --tag $(VERSION) --template-file /tmp/template-file.yaml
