SHELL = /bin/bash

# Setting GOBIN makes 'go install' put the binary in the bin/ directory.
export GOBIN ?= $(dir $(abspath $(lastword $(MAKEFILE_LIST))))/bin

STITCHMD = $(GOBIN)/stitchmd

# TFLint plugin variables
PLUGIN_NAME = tflint-ruleset-terraform-style
PLUGIN_DIR = ~/.tflint.d/plugins

# Keep these options in-sync with .github/workflows/ci.yml.
# STITCHMD_ARGS = -o style.md -preface src/preface.txt src/SUMMARY.md

# Keep these options in-sync with .github/workflows/ci.yml.
# STITCHMD_GOLANG_ARGS = -o styleguides/golang.md -preface styleguides/golang/src/preface.txt styleguides/golang/src/SUMMARY.md

# Keep these options in-sync with .github/workflows/ci.yml.
STITCHMD_TERRAFORM_ARGS = -o docs/styleguides/terraform.md -preface docs/styleguides/terraform/src/preface.txt docs/styleguides/terraform/src/SUMMARY.md


.PHONY: all
all: style.md

.PHONY: lint
lint: $(STITCHMD)
	@DIFF=$$($(STITCHMD) -d $(STITCHMD_ARGS)); \
	if [[ -n "$$DIFF" ]]; then \
		echo "style.md is out of date:"; \
		echo "$$DIFF"; \
		false; \
	fi

$(STITCHMD):
	go install go.abhg.dev/stitchmd@latest

.PHONY: terraform-summary
terraform-summary:
	./scripts/generate-summary.sh docs/styleguides/terraform/src "Terraform Style Guide"

terraform.md: $(STITCHMD) terraform-summary $(wildcard styleguides/terraform/src/*)
	$(STITCHMD) $(STITCHMD_TERRAFORM_ARGS)

# =============================================================================
# TFLint Plugin Targets
# =============================================================================

.PHONY: build
build:
	go build -o $(PLUGIN_NAME)

.PHONY: test
test:
	go test ./...

.PHONY: test-cover
test-cover:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

.PHONY: go-lint
go-lint:
	golangci-lint run

.PHONY: install
install: build
	mkdir -p $(PLUGIN_DIR)
	mv $(PLUGIN_NAME) $(PLUGIN_DIR)/

.PHONY: clean
clean:
	rm -f $(PLUGIN_NAME)
	rm -f coverage.out coverage.html

# =============================================================================
# Integration Test Targets
# =============================================================================

.PHONY: test-unit
test-unit:
	go test -v ./rules/...

.PHONY: test-fixtures
test-fixtures: install
	@echo "Running fixture-based TFLint tests..."
	cd tests && go test -v -run "TestFixtures" -timeout 5m

.PHONY: test-localstack
test-localstack: install localstack-up
	@echo "Running LocalStack integration tests (Terraform + OpenTofu)..."
	cd tests && go test -v -run "Terraform|Tofu|TestLintThenApply" -timeout 30m

.PHONY: test-localstack-terraform
test-localstack-terraform: install localstack-up
	@echo "Running LocalStack tests with Terraform only..."
	cd tests && go test -v -run "Terraform" -timeout 15m

.PHONY: test-localstack-tofu
test-localstack-tofu: install localstack-up
	@echo "Running LocalStack tests with OpenTofu only..."
	cd tests && go test -v -run "Tofu" -timeout 15m

.PHONY: test-integration
test-integration: install
	@echo "Running all integration tests..."
	cd tests && go test -v -timeout 30m

.PHONY: test-integration-short
test-integration-short: install
	@echo "Running integration tests (short mode - no LocalStack)..."
	cd tests && go test -v -short -timeout 5m

.PHONY: test-all
test-all: test-unit test-integration
	@echo "All tests completed!"

# =============================================================================
# LocalStack Targets
# =============================================================================

.PHONY: localstack-up
localstack-up:
	@echo "Starting LocalStack..."
	docker compose up -d
	@echo "Waiting for LocalStack to be ready..."
	@until curl -sf http://localhost:4566/_localstack/health > /dev/null 2>&1; do \
		sleep 1; \
	done
	@echo "LocalStack is ready!"

.PHONY: localstack-down
localstack-down:
	@echo "Stopping LocalStack..."
	docker compose down

.PHONY: localstack-logs
localstack-logs:
	docker compose logs -f localstack

# =============================================================================
# Release Targets
# =============================================================================

.PHONY: release-check
release-check:
	goreleaser check

.PHONY: release-snapshot
release-snapshot:
	goreleaser release --snapshot --clean --skip=sign

