# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This repository has two purposes:
1. **Style Guide Documentation**: Terraform style guide documentation compiled from markdown sources in `docs/styleguides/terraform/src/` into a single `docs/styleguides/terraform.md` using stitchmd
2. **TFLint Plugin**: A custom TFLint plugin (to be implemented) that enforces the rules defined in the style guide

## Documentation Commands

```bash
# Generate SUMMARY.md for terraform style guide
make terraform-summary

# Generate compiled terraform.md from source files
make terraform.md

# Install stitchmd (required for doc generation)
go install go.abhg.dev/stitchmd@latest
```

## Tool Versions (mise.toml)

- Go 1.25.4
- golangci-lint 2.5.0

## TFLint Plugin Development

See `docs/IMPLEMENTATION_PLAN.md` for detailed implementation guidance.

### Existing TFLint Coverage (DO NOT reimplement)

These rules are already covered by `tflint-ruleset-terraform`:
- `terraform_naming_convention` - Snake case naming
- `terraform_comment_syntax` - Catches `//` comments
- `terraform_documented_variables` - Variable description required
- `terraform_typed_variables` - Variable type required
- `terraform_documented_outputs` - Output description required
- `terraform_standard_module_structure` - Standard file structure

### Rules to Implement

- `terraform_block_comment_syntax` - Block comments `/**/` not allowed (existing rule only catches `//`)
- `terraform_comment_delimiter` - Use `# ---` for section headers, not `# ~~~`
- `terraform_tautological_naming` - Resource name must not repeat type words (e.g., `aws_iam_policy.iam_policy` is bad)
- `terraform_variable_attribute_order` - Order: `description` → `type` → `default`
- `terraform_output_attribute_order` - `description` must come before `value`
- `terraform_conditional_parentheses` - Multi-line ternary expressions need parentheses
- `terraform_resource_parameter_order` - Order: `count/for_each` → params → blocks → `lifecycle` → `depends_on`

## Architecture Notes

TFLint plugins are written in Go using the tflint-plugin-sdk. Each rule is implemented as a separate rule struct that:
1. Implements the `tflint.Rule` interface
2. Has a unique name, enabled state, and severity
3. Contains a `Check` method that inspects the Terraform configuration

## Plugin Details

- **Name**: `tflint-ruleset-terraform-style`
- **Module**: `github.com/donaldgifford/tflint-ruleset-terraform-style`

Standard plugin structure:
```
├── main.go              # Plugin entry point
├── rules/               # Individual rule implementations
├── project.go           # Rule set definition
└── go.mod              # Go module definition
```
