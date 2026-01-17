package rules

import (
	"testing"

	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func TestTerraformConditionalParentheses(t *testing.T) {
	tests := []struct {
		Name        string
		Content     string
		IssueCount  int
		ExpectIssue bool
	}{
		{
			Name: "single line conditional allowed",
			Content: `
variable "enabled" {
  type    = bool
  default = true
}

locals {
  value = var.enabled ? "yes" : "no"
}
`,
			ExpectIssue: false,
			IssueCount:  0,
		},
		{
			Name: "multi-line conditional with parentheses allowed",
			Content: `
variable "enabled" {
  type    = bool
  default = true
}

locals {
  value = (
    var.enabled
    ? "yes"
    : "no"
  )
}
`,
			ExpectIssue: false,
			IssueCount:  0,
		},
		{
			Name: "nested conditional with parentheses",
			Content: `
variable "a" {
  type    = bool
  default = true
}

variable "b" {
  type    = bool
  default = false
}

locals {
  value = (
    var.a
    ? (
        var.b
        ? "ab"
        : "a"
      )
    : "none"
  )
}
`,
			ExpectIssue: false,
			IssueCount:  0,
		},
	}

	rule := NewTerraformConditionalParenthesesRule()

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			runner := helper.TestRunner(t, map[string]string{"main.tf": test.Content})

			if err := rule.Check(runner); err != nil {
				t.Fatalf("Unexpected error: %s", err)
			}

			if test.ExpectIssue {
				if len(runner.Issues) != test.IssueCount {
					t.Errorf("Expected %d issues, got %d", test.IssueCount, len(runner.Issues))
				}
				for _, issue := range runner.Issues {
					if issue.Rule.Name() != rule.Name() {
						t.Errorf("Expected rule %s, got %s", rule.Name(), issue.Rule.Name())
					}
				}
			} else {
				if len(runner.Issues) != 0 {
					t.Errorf("Expected no issues, got %d", len(runner.Issues))
				}
			}
		})
	}
}
