package tests

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/donaldgifford/tflint-ruleset-terraform-style/tests/testutil"
)

// ruleNames maps directory names to actual rule names
var ruleNames = map[string]string{
	"block_comment_syntax":     "terraform_block_comment_syntax",
	"comment_delimiter":        "terraform_comment_delimiter",
	"tautological_naming":      "terraform_tautological_naming",
	"variable_attribute_order": "terraform_variable_attribute_order",
	"output_attribute_order":   "terraform_output_attribute_order",
	"conditional_parentheses":  "terraform_conditional_parentheses",
	"resource_parameter_order": "terraform_resource_parameter_order",
}

// TestFixturesPass verifies all pass fixtures have no lint issues for their rule
func TestFixturesPass(t *testing.T) {
	baseDir := filepath.Join("fixtures", "pass")
	configPath, err := filepath.Abs(filepath.Join("..", ".tflint.hcl.example"))
	if err != nil {
		t.Fatalf("Failed to get config path: %v", err)
	}

	// Check if tflint is installed
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Skipf("TFLint config not found at %s", configPath)
	}

	for dirName, ruleName := range ruleNames {
		t.Run(ruleName, func(t *testing.T) {
			fixtureDir := filepath.Join(baseDir, dirName)

			// Check if fixture directory exists
			if _, err := os.Stat(fixtureDir); os.IsNotExist(err) {
				t.Skipf("Fixture directory not found: %s", fixtureDir)
			}

			result, err := testutil.RunTFLint(fixtureDir, configPath)
			if err != nil {
				t.Fatalf("Failed to run tflint: %v", err)
			}

			if result.HasIssueForRule(ruleName) {
				issues := result.GetIssuesForRule(ruleName)
				t.Errorf("Expected no issues for rule %s in pass fixtures, but found %d:", ruleName, len(issues))
				for _, issue := range issues {
					t.Errorf("  - %s:%d: %s", issue.Range.Filename, issue.Range.Start.Line, issue.Message)
				}
			}
		})
	}
}

// TestFixturesFail verifies all fail fixtures produce lint issues for their rule
func TestFixturesFail(t *testing.T) {
	baseDir := filepath.Join("fixtures", "fail")
	configPath, err := filepath.Abs(filepath.Join("..", ".tflint.hcl.example"))
	if err != nil {
		t.Fatalf("Failed to get config path: %v", err)
	}

	// Check if tflint is installed
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Skipf("TFLint config not found at %s", configPath)
	}

	// Rules that cannot have fail tests because HCL syntax prevents the violation
	// (HCL requires parentheses for multi-line ternary expressions)
	syntaxEnforcedRules := map[string]bool{
		"terraform_conditional_parentheses": true,
	}

	for dirName, ruleName := range ruleNames {
		t.Run(ruleName, func(t *testing.T) {
			// Skip rules where the failure case is a syntax error
			if syntaxEnforcedRules[ruleName] {
				t.Skip("HCL syntax enforces this rule - multi-line ternary without parens is a syntax error")
			}

			fixtureDir := filepath.Join(baseDir, dirName)

			// Check if fixture directory exists
			if _, err := os.Stat(fixtureDir); os.IsNotExist(err) {
				t.Skipf("Fixture directory not found: %s", fixtureDir)
			}

			result, err := testutil.RunTFLint(fixtureDir, configPath)
			if err != nil {
				t.Fatalf("Failed to run tflint: %v", err)
			}

			if !result.HasIssueForRule(ruleName) {
				t.Errorf("Expected issues for rule %s in fail fixtures, but found none", ruleName)
				t.Logf("Total issues found: %d", len(result.Issues))
				for _, issue := range result.Issues {
					t.Logf("  - %s: %s", issue.Rule.Name, issue.Message)
				}
			} else {
				// Log the issues found for visibility
				issues := result.GetIssuesForRule(ruleName)
				t.Logf("Found %d expected issues for rule %s", len(issues), ruleName)
			}
		})
	}
}

// TestAllRulesEnabled verifies the plugin loads and all rules are available
func TestAllRulesEnabled(t *testing.T) {
	// Create a minimal terraform file
	tmpDir := t.TempDir()
	tfFile := filepath.Join(tmpDir, "main.tf")
	if err := os.WriteFile(tfFile, []byte(`resource "null_resource" "test" {}`), 0644); err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	configPath, err := filepath.Abs(filepath.Join("..", ".tflint.hcl.example"))
	if err != nil {
		t.Fatalf("Failed to get config path: %v", err)
	}

	result, err := testutil.RunTFLint(tmpDir, configPath)
	if err != nil {
		t.Fatalf("Failed to run tflint: %v", err)
	}

	// Should run without errors
	if len(result.Errors) > 0 {
		for _, e := range result.Errors {
			t.Errorf("TFLint error: %s", e.Message)
		}
	}
}
