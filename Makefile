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
	mkdir -p build/bin
	go build -o build/bin/$(PLUGIN_NAME)

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
	cp build/bin/$(PLUGIN_NAME) $(PLUGIN_DIR)/

.PHONY: clean
clean:
	rm -rf build/
	rm -f coverage.out coverage.html

# =============================================================================
# Test Targets
# =============================================================================

.PHONY: test-unit
test-unit:
	go test -v ./rules/...

.PHONY: test-fixtures
test-fixtures: install
	@echo "Running TFLint fixture tests..."
	@for dir in testdata/valid/*/; do \
		echo "Testing $$dir (expecting pass)..."; \
		tflint --chdir="$$dir" --config=../../../.tflint.hcl || exit 1; \
	done
	@for dir in testdata/invalid/*/; do \
		echo "Testing $$dir (expecting failure)..."; \
		! tflint --chdir="$$dir" --config=../../../.tflint.hcl || exit 1; \
	done
	@echo "All fixture tests passed!"

.PHONY: test-all
test-all: test-unit test-fixtures
	@echo "All tests completed!"

# =============================================================================
# Release Targets
# =============================================================================

.PHONY: release-check
release-check:
	goreleaser check

.PHONY: release-snapshot
release-snapshot:
	goreleaser release --snapshot --clean --skip=sign

