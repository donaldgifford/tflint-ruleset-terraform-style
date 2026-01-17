package rules

import (
	"regexp"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// tildeDelimiterRegex matches comment lines using tildes as delimiters (e.g., # ~~~).
var tildeDelimiterRegex = regexp.MustCompile(`^#\s*~{3,}`)

// TerraformCommentDelimiterRule checks that section delimiters use dashes, not tildes.
type TerraformCommentDelimiterRule struct {
	tflint.DefaultRule
}

// NewTerraformCommentDelimiterRule returns a new rule instance.
func NewTerraformCommentDelimiterRule() *TerraformCommentDelimiterRule {
	return &TerraformCommentDelimiterRule{}
}

// Name returns the rule name.
func (*TerraformCommentDelimiterRule) Name() string {
	return "terraform_comment_delimiter"
}

// Enabled returns whether the rule is enabled by default.
func (*TerraformCommentDelimiterRule) Enabled() bool {
	return true
}

// Severity returns the rule severity.
func (*TerraformCommentDelimiterRule) Severity() tflint.Severity {
	return tflint.WARNING
}

// Link returns the rule documentation link.
func (*TerraformCommentDelimiterRule) Link() string {
	return "https://github.com/donaldgifford/tflint-ruleset-terraform-style/blob/main/docs/rules/terraform_comment_delimiter.md"
}

// Check runs the rule check.
func (r *TerraformCommentDelimiterRule) Check(runner tflint.Runner) error {
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

func (r *TerraformCommentDelimiterRule) checkFile(runner tflint.Runner, filename string, file *hcl.File) error {
	tokens, diags := hclsyntax.LexConfig(file.Bytes, filename, hcl.InitialPos)
	if diags.HasErrors() {
		return nil // Skip files with syntax errors
	}

	for _, token := range tokens {
		if token.Type != hclsyntax.TokenComment {
			continue
		}

		if tildeDelimiterRegex.Match(token.Bytes) {
			if err := runner.EmitIssue(
				r,
				"Use '# ---' for section delimiters, not '# ~~~'",
				token.Range,
			); err != nil {
				return err
			}
		}
	}

	return nil
}
