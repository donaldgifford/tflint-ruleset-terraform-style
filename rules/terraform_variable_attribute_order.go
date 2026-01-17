package rules

import (
	"fmt"

	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// TerraformVariableAttributeOrderRule checks that variable attributes are in the correct order.
// Expected order: description, type, default, validation, sensitive, nullable
type TerraformVariableAttributeOrderRule struct {
	tflint.DefaultRule
}

// NewTerraformVariableAttributeOrderRule returns a new rule instance.
func NewTerraformVariableAttributeOrderRule() *TerraformVariableAttributeOrderRule {
	return &TerraformVariableAttributeOrderRule{}
}

// Name returns the rule name.
func (r *TerraformVariableAttributeOrderRule) Name() string {
	return "terraform_variable_attribute_order"
}

// Enabled returns whether the rule is enabled by default.
func (r *TerraformVariableAttributeOrderRule) Enabled() bool {
	return true
}

// Severity returns the rule severity.
func (r *TerraformVariableAttributeOrderRule) Severity() tflint.Severity {
	return tflint.WARNING
}

// Link returns the rule documentation link.
func (r *TerraformVariableAttributeOrderRule) Link() string {
	return "https://github.com/donaldgifford/tflint-ruleset-terraform-style/blob/main/docs/rules/terraform_variable_attribute_order.md"
}

// variableAttributeOrder defines the expected order of attributes in a variable block.
var variableAttributeOrder = []string{"description", "type", "default", "validation", "sensitive", "nullable"}

// Check runs the rule check.
func (r *TerraformVariableAttributeOrderRule) Check(runner tflint.Runner) error {
	body, err := runner.GetModuleContent(&hclext.BodySchema{
		Blocks: []hclext.BlockSchema{
			{
				Type:       "variable",
				LabelNames: []string{"name"},
				Body: &hclext.BodySchema{
					Attributes: []hclext.AttributeSchema{
						{Name: "description"},
						{Name: "type"},
						{Name: "default"},
						{Name: "sensitive"},
						{Name: "nullable"},
					},
					Blocks: []hclext.BlockSchema{
						{Type: "validation", Body: &hclext.BodySchema{}},
					},
				},
			},
		},
	}, &tflint.GetModuleContentOption{
		ExpandMode: tflint.ExpandModeNone,
	})
	if err != nil {
		return err
	}

	for _, block := range body.Blocks {
		if len(block.Labels) < 1 {
			continue
		}
		varName := block.Labels[0]

		if err := r.checkAttributeOrder(runner, varName, block); err != nil {
			return err
		}
	}

	return nil
}

func (r *TerraformVariableAttributeOrderRule) checkAttributeOrder(runner tflint.Runner, varName string, block *hclext.Block) error {
	// Build a map of attribute name to line number
	attrLines := make(map[string]int)
	for name, attr := range block.Body.Attributes {
		attrLines[name] = attr.Range.Start.Line
	}

	// Check for validation blocks
	for _, nestedBlock := range block.Body.Blocks {
		if nestedBlock.Type == "validation" {
			attrLines["validation"] = nestedBlock.DefRange.Start.Line
			break // Only consider first validation block for ordering
		}
	}

	// Check ordering between consecutive pairs
	var lastAttr string
	var lastLine int

	for _, attrName := range variableAttributeOrder {
		line, exists := attrLines[attrName]
		if !exists {
			continue
		}

		if lastAttr != "" && line < lastLine {
			if err := runner.EmitIssue(
				r,
				fmt.Sprintf("In variable %q, %q should come before %q", varName, attrName, lastAttr),
				block.Body.Attributes[attrName].Range,
			); err != nil {
				return err
			}
		}

		lastAttr = attrName
		lastLine = line
	}

	return nil
}
