package rules

import (
	"testing"

	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func TestTerraformOutputAttributeOrder(t *testing.T) {
	tests := []struct {
		Name        string
		Content     string
		IssueCount  int
		ExpectIssue bool
	}{
		{
			Name: "correct order",
			Content: `
output "example" {
  description = "An example output"
  value       = "hello"
}
`,
			ExpectIssue: false,
			IssueCount:  0,
		},
		{
			Name: "value before description",
			Content: `
output "example" {
  value       = "hello"
  description = "An example output"
}
`,
			ExpectIssue: true,
			IssueCount:  1,
		},
		{
			Name: "correct order with sensitive",
			Content: `
output "example" {
  description = "An example output"
  value       = "hello"
  sensitive   = true
}
`,
			ExpectIssue: false,
			IssueCount:  0,
		},
		{
			Name: "sensitive before value",
			Content: `
output "example" {
  description = "An example output"
  sensitive   = true
  value       = "hello"
}
`,
			ExpectIssue: true,
			IssueCount:  1,
		},
		{
			Name: "only value",
			Content: `
output "example" {
  value = "hello"
}
`,
			ExpectIssue: false,
			IssueCount:  0,
		},
	}

	rule := NewTerraformOutputAttributeOrderRule()

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			runner := helper.TestRunner(t, map[string]string{"outputs.tf": test.Content})

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
