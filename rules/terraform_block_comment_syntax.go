package rules

import (
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// TerraformBlockCommentSyntaxRule checks that block comments (/* */) are not used.
type TerraformBlockCommentSyntaxRule struct {
	tflint.DefaultRule
}

// NewTerraformBlockCommentSyntaxRule returns a new rule instance.
func NewTerraformBlockCommentSyntaxRule() *TerraformBlockCommentSyntaxRule {
	return &TerraformBlockCommentSyntaxRule{}
}

// Name returns the rule name.
func (r *TerraformBlockCommentSyntaxRule) Name() string {
	return "terraform_block_comment_syntax"
}

// Enabled returns whether the rule is enabled by default.
func (r *TerraformBlockCommentSyntaxRule) Enabled() bool {
	return true
}

// Severity returns the rule severity.
func (r *TerraformBlockCommentSyntaxRule) Severity() tflint.Severity {
	return tflint.WARNING
}

// Link returns the rule documentation link.
func (r *TerraformBlockCommentSyntaxRule) Link() string {
	return "https://github.com/donaldgifford/tflint-ruleset-terraform-style/blob/main/docs/rules/terraform_block_comment_syntax.md"
}

// Check runs the rule check.
func (r *TerraformBlockCommentSyntaxRule) Check(runner tflint.Runner) error {
	files, err := runner.GetFiles()
	if err != nil {
		return err
	}

	for filename, file := range files {
		if err := r.checkFile(runner, filename, file); err != nil {
			return err
		}
	}

	return nil
}

func (r *TerraformBlockCommentSyntaxRule) checkFile(runner tflint.Runner, filename string, file *hcl.File) error {
	tokens, diags := hclsyntax.LexConfig(file.Bytes, filename, hcl.InitialPos)
	if diags.HasErrors() {
		return nil // Skip files with syntax errors
	}

	for _, token := range tokens {
		if token.Type != hclsyntax.TokenComment {
			continue
		}

		tokenStr := string(token.Bytes)
		if strings.HasPrefix(tokenStr, "/*") {
			if err := runner.EmitIssue(
				r,
				"Block comments (/* */) are not allowed; use # instead",
				token.Range,
			); err != nil {
				return err
			}
		}
	}

	return nil
}
