package rules

import (
	"testing"

	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func TestTerraformVariableAttributeOrder(t *testing.T) {
	tests := []struct {
		Name        string
		Content     string
		IssueCount  int
		ExpectIssue bool
	}{
		{
			Name: "correct order",
			Content: `
variable "example" {
  description = "An example variable"
  type        = string
  default     = "example"
}
`,
			ExpectIssue: false,
			IssueCount:  0,
		},
		{
			Name: "type before description",
			Content: `
variable "example" {
  type        = string
  description = "An example variable"
  default     = "example"
}
`,
			ExpectIssue: true,
			IssueCount:  1,
		},
		{
			Name: "default before type",
			Content: `
variable "example" {
  description = "An example variable"
  default     = "example"
  type        = string
}
`,
			ExpectIssue: true,
			IssueCount:  1,
		},
		{
			Name: "completely reversed order",
			Content: `
variable "example" {
  default     = "example"
  type        = string
  description = "An example variable"
}
`,
			ExpectIssue: true,
			IssueCount:  2,
		},
		{
			Name: "only description and type correct order",
			Content: `
variable "example" {
  description = "An example variable"
  type        = string
}
`,
			ExpectIssue: false,
			IssueCount:  0,
		},
		{
			Name: "only description",
			Content: `
variable "example" {
  description = "An example variable"
}
`,
			ExpectIssue: false,
			IssueCount:  0,
		},
	}

	rule := NewTerraformVariableAttributeOrderRule()

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			runner := helper.TestRunner(t, map[string]string{"variables.tf": test.Content})

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
			} else if len(runner.Issues) != 0 {
				t.Errorf("Expected no issues, got %d", len(runner.Issues))
			}
		})
	}
}
