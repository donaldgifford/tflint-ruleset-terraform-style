package rules

import (
	"fmt"
	"strings"

	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// ignoredWords are common words that should not trigger tautological naming warnings.
var ignoredWords = map[string]bool{
	"aws":   true,
	"gcp":   true,
	"azure": true,
	"main":  true,
	"this":  true,
	"that":  true,
	"the":   true,
	"a":     true,
	"an":    true,
}

// TerraformTautologicalNamingRule checks that resource names don't repeat words from the resource type.
type TerraformTautologicalNamingRule struct {
	tflint.DefaultRule
}

// NewTerraformTautologicalNamingRule returns a new rule instance.
func NewTerraformTautologicalNamingRule() *TerraformTautologicalNamingRule {
	return &TerraformTautologicalNamingRule{}
}

// Name returns the rule name.
func (r *TerraformTautologicalNamingRule) Name() string {
	return "terraform_tautological_naming"
}

// Enabled returns whether the rule is enabled by default.
func (r *TerraformTautologicalNamingRule) Enabled() bool {
	return true
}

// Severity returns the rule severity.
func (r *TerraformTautologicalNamingRule) Severity() tflint.Severity {
	return tflint.WARNING
}

// Link returns the rule documentation link.
func (r *TerraformTautologicalNamingRule) Link() string {
	return "https://github.com/donaldgifford/tflint-ruleset-terraform-style/blob/main/docs/rules/terraform_tautological_naming.md"
}

// Check runs the rule check.
func (r *TerraformTautologicalNamingRule) Check(runner tflint.Runner) error {
	body, err := runner.GetModuleContent(&hclext.BodySchema{
		Blocks: []hclext.BlockSchema{
			{
				Type:       "resource",
				LabelNames: []string{"type", "name"},
				Body:       &hclext.BodySchema{},
			},
			{
				Type:       "data",
				LabelNames: []string{"type", "name"},
				Body:       &hclext.BodySchema{},
			},
		},
	}, &tflint.GetModuleContentOption{
		ExpandMode: tflint.ExpandModeNone,
	})
	if err != nil {
		return err
	}

	for _, block := range body.Blocks {
		if len(block.Labels) < 2 {
			continue
		}

		resourceType := block.Labels[0]
		resourceName := block.Labels[1]

		typeWords := splitSnakeCase(resourceType)
		nameWords := splitSnakeCase(resourceName)

		for _, nameWord := range nameWords {
			if ignoredWords[nameWord] {
				continue
			}
			if containsWord(typeWords, nameWord) {
				if err := runner.EmitIssue(
					r,
					fmt.Sprintf("Resource name %q contains word %q which already appears in the resource type %q", resourceName, nameWord, resourceType),
					block.DefRange,
				); err != nil {
					return err
				}
				break // Only report once per resource
			}
		}
	}

	return nil
}

// splitSnakeCase splits a snake_case string into individual words.
func splitSnakeCase(s string) []string {
	parts := strings.Split(s, "_")
	words := make([]string, 0, len(parts))
	for _, part := range parts {
		if part != "" {
			words = append(words, strings.ToLower(part))
		}
	}
	return words
}

// containsWord checks if a word exists in a slice of words.
func containsWord(words []string, word string) bool {
	word = strings.ToLower(word)
	for _, w := range words {
		if w == word {
			return true
		}
	}
	return false
}
