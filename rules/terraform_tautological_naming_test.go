package rules

import (
	"testing"

	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func TestTerraformTautologicalNaming(t *testing.T) {
	tests := []struct {
		Name        string
		Content     string
		IssueCount  int
		ExpectIssue bool
	}{
		{
			Name: "tautological resource name with iam",
			Content: `
resource "aws_iam_policy" "iam_policy" {}
`,
			ExpectIssue: true,
			IssueCount:  1,
		},
		{
			Name: "tautological with policy word",
			Content: `
resource "aws_iam_policy" "team_policy" {}
`,
			ExpectIssue: true,
			IssueCount:  1,
		},
		{
			Name: "tautological bucket name",
			Content: `
resource "aws_s3_bucket" "logs_bucket" {}
`,
			ExpectIssue: true,
			IssueCount:  1,
		},
		{
			Name: "good resource name",
			Content: `
resource "aws_iam_policy" "team_access" {}
`,
			ExpectIssue: false,
			IssueCount:  0,
		},
		{
			Name: "good bucket name",
			Content: `
resource "aws_s3_bucket" "logs" {}
`,
			ExpectIssue: false,
			IssueCount:  0,
		},
		{
			Name: "aws ignored in name",
			Content: `
resource "aws_instance" "aws_web_server" {}
`,
			ExpectIssue: false,
			IssueCount:  0,
		},
		{
			Name: "main ignored in name",
			Content: `
resource "aws_instance" "main" {}
`,
			ExpectIssue: false,
			IssueCount:  0,
		},
		{
			Name: "this ignored in name",
			Content: `
resource "aws_instance" "this" {}
`,
			ExpectIssue: false,
			IssueCount:  0,
		},
		{
			Name: "data source tautological",
			Content: `
data "aws_iam_policy" "existing_policy" {}
`,
			ExpectIssue: true,
			IssueCount:  1,
		},
	}

	rule := NewTerraformTautologicalNamingRule()

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
