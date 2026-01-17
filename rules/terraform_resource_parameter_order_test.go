package rules

import (
	"testing"

	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func TestTerraformResourceParameterOrder(t *testing.T) {
	tests := []struct {
		Name        string
		Content     string
		IssueCount  int
		ExpectIssue bool
	}{
		{
			Name: "correct order",
			Content: `
resource "aws_instance" "example" {
  count = 1

  ami           = "ami-12345"
  instance_type = "t2.micro"

  tags = {
    Name = "example"
  }

  lifecycle {
    create_before_destroy = true
  }

  depends_on = [aws_vpc.main]
}
`,
			ExpectIssue: false,
			IssueCount:  0,
		},
		{
			Name: "depends_on before lifecycle",
			Content: `
resource "aws_instance" "example" {
  ami           = "ami-12345"
  instance_type = "t2.micro"

  depends_on = [aws_vpc.main]

  lifecycle {
    create_before_destroy = true
  }
}
`,
			ExpectIssue: true,
			IssueCount:  1,
		},
		{
			Name: "count not first",
			Content: `
resource "aws_instance" "example" {
  ami   = "ami-12345"
  count = 1
  instance_type = "t2.micro"
}
`,
			ExpectIssue: true,
			IssueCount:  1,
		},
		{
			Name: "for_each not first",
			Content: `
resource "aws_instance" "example" {
  ami      = "ami-12345"
  for_each = toset(["a", "b"])
  instance_type = "t2.micro"
}
`,
			ExpectIssue: true,
			IssueCount:  1,
		},
		{
			Name: "data source correct order",
			Content: `
data "aws_ami" "example" {
  most_recent = true

  filter {
    name   = "name"
    values = ["ubuntu/*"]
  }
}
`,
			ExpectIssue: false,
			IssueCount:  0,
		},
	}

	rule := NewTerraformResourceParameterOrderRule()

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
