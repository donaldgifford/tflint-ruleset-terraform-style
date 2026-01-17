# TFLint Plugin Implementation Plan

This document outlines the plan to create a TFLint plugin that enforces the Terraform style guide defined in `docs/styleguides/terraform.md`.

## Existing TFLint Coverage

The following style guide rules are already covered by `tflint-ruleset-terraform` and should NOT be reimplemented:

- Snake case naming → `terraform_naming_convention`
- `//` comment syntax → `terraform_comment_syntax`
- Variable description required → `terraform_documented_variables`
- Variable type required → `terraform_typed_variables`
- Output description required → `terraform_documented_outputs`
- Module structure (main.tf, variables.tf, outputs.tf) → `terraform_standard_module_structure`

## Rules to Implement

### Rule 1: `terraform_block_comment_syntax`
- **Source**: "Use `#` for comment strings, not `//` or `/**/`"
- **Why**: Existing `terraform_comment_syntax` only catches `//`, not `/**/`
- **Severity**: Warning

### Rule 2: `terraform_comment_delimiter`
- **Source**: "Prefer `# -` over `# ~`"
- **Checks**: Section header comment blocks use `# ---` not `# ~~~`
- **Severity**: Warning

### Rule 3: `terraform_tautological_naming`
- **Source**: "Avoid tautological resource names"
- **Checks**: Resource name doesn't contain words from the resource type
- **Example**: `aws_iam_policy.ec2_policy` → warns about "policy" in name
- **Severity**: Warning

### Rule 4: `terraform_variable_attribute_order`
- **Source**: "description and type... in that order"
- **Checks**: Variable blocks have attributes in order: `description` → `type` → `default`
- **Severity**: Warning

### Rule 5: `terraform_output_attribute_order`
- **Source**: "description, before the value"
- **Checks**: Output blocks have `description` before `value`
- **Severity**: Warning

### Rule 6: `terraform_conditional_parentheses`
- **Source**: "Use () to break up conditionals across multiple lines"
- **Checks**: Multi-line ternary expressions are wrapped in parentheses
- **Severity**: Warning

### Rule 7: `terraform_resource_parameter_order`
- **Source**: Resource order section
- **Checks**: Resource blocks have parameters in order: `count/for_each` → params → blocks → `lifecycle` → `depends_on`
- **Severity**: Warning

## Phases

### Phase 1: Project Scaffolding
- Initialize Go module
- Create main.go, project.go, rules/ directory
- Update Makefile with build/test/install targets
- Add example .tflint.hcl configuration

### Phase 2: Implement Rules 1-5
Core rules that are straightforward to implement.

### Phase 3: Implement Rules 6-7
More complex rules requiring deeper AST analysis.

### Phase 4: Testing
- Unit tests for each rule
- Test fixtures in testdata/
- Integration tests

### Phase 5: Documentation & Release
- README with installation and configuration
- GitHub Actions for CI/release
- Multi-platform binaries

## File Structure

```
tf-style-guide/
├── main.go
├── project.go
├── go.mod
├── Makefile
├── .tflint.hcl.example
├── rules/
│   ├── terraform_block_comment_syntax.go
│   ├── terraform_comment_delimiter.go
│   ├── terraform_tautological_naming.go
│   ├── terraform_variable_attribute_order.go
│   ├── terraform_output_attribute_order.go
│   ├── terraform_conditional_parentheses.go
│   └── terraform_resource_parameter_order.go
├── testdata/
└── docs/
```

## Decisions

- **Plugin name**: `tflint-ruleset-terraform-style`
- **Rule severity**: All default to Warning
- **Go module path**: `github.com/donaldgifford/tflint-ruleset-terraform-style`
- **Repo rename**: Rename repo from `tf-style-guide` to `tflint-ruleset-terraform-style`
