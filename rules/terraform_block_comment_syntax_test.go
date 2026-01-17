package rules

import (
	"testing"

	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func TestTerraformBlockCommentSyntax(t *testing.T) {
	tests := []struct {
		Name        string
		Content     string
		IssueCount  int
		ExpectIssue bool
	}{
		{
			Name: "block comment detected",
			Content: `
/* This is a block comment */
resource "aws_instance" "test" {}
`,
			ExpectIssue: true,
			IssueCount:  1,
		},
		{
			Name: "multi-line block comment detected",
			Content: `
/*
 * Multi-line
 * block comment
 */
resource "aws_instance" "test" {}
`,
			ExpectIssue: true,
			IssueCount:  1,
		},
		{
			Name: "hash comment allowed",
			Content: `
# This is a hash comment
resource "aws_instance" "test" {}
`,
			ExpectIssue: false,
			IssueCount:  0,
		},
		{
			Name: "inline block comment detected",
			Content: `
resource "aws_instance" "test" {
  ami = "ami-12345" /* inline comment */
}
`,
			ExpectIssue: true,
			IssueCount:  1,
		},
	}

	rule := NewTerraformBlockCommentSyntaxRule()

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
			} else if len(runner.Issues) != 0 {
				t.Errorf("Expected no issues, got %d", len(runner.Issues))
			}
		})
	}
}
