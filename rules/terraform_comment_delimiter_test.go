package rules

import (
	"testing"

	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func TestTerraformCommentDelimiter(t *testing.T) {
	tests := []struct {
		Name        string
		Content     string
		IssueCount  int
		ExpectIssue bool
	}{
		{
			Name: "tilde delimiter detected",
			Content: `
# ~~~~~~~~~~~~~~~
# Section Header
# ~~~~~~~~~~~~~~~
resource "aws_instance" "test" {}
`,
			ExpectIssue: true,
			IssueCount:  2,
		},
		{
			Name: "short tilde delimiter detected",
			Content: `
# ~~~
resource "aws_instance" "test" {}
`,
			ExpectIssue: true,
			IssueCount:  1,
		},
		{
			Name: "dash delimiter allowed",
			Content: `
# ---------------
# Section Header
# ---------------
resource "aws_instance" "test" {}
`,
			ExpectIssue: false,
			IssueCount:  0,
		},
		{
			Name: "regular comment allowed",
			Content: `
# This is a regular comment
resource "aws_instance" "test" {}
`,
			ExpectIssue: false,
			IssueCount:  0,
		},
		{
			Name: "tilde in regular comment allowed",
			Content: `
# Use ~ for home directory
resource "aws_instance" "test" {}
`,
			ExpectIssue: false,
			IssueCount:  0,
		},
	}

	rule := NewTerraformCommentDelimiterRule()

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			runner := helper.TestRunner(t, map[string]string{"resource.tf": test.Content})

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
