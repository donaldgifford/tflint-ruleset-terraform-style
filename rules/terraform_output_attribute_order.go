package rules

import (
	"fmt"

	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// TerraformOutputAttributeOrderRule checks that output attributes are in the correct order.
// Expected order: description, value, sensitive, depends_on.
type TerraformOutputAttributeOrderRule struct {
	tflint.DefaultRule
}

// NewTerraformOutputAttributeOrderRule returns a new rule instance.
func NewTerraformOutputAttributeOrderRule() *TerraformOutputAttributeOrderRule {
	return &TerraformOutputAttributeOrderRule{}
}

// Name returns the rule name.
func (*TerraformOutputAttributeOrderRule) Name() string {
	return "terraform_output_attribute_order"
}

// Enabled returns whether the rule is enabled by default.
func (*TerraformOutputAttributeOrderRule) Enabled() bool {
	return true
}

// Severity returns the rule severity.
func (*TerraformOutputAttributeOrderRule) Severity() tflint.Severity {
	return tflint.WARNING
}

// Link returns the rule documentation link.
func (*TerraformOutputAttributeOrderRule) Link() string {
	return "https://github.com/donaldgifford/tflint-ruleset-terraform-style/blob/main/docs/rules/terraform_output_attribute_order.md"
}

// outputAttributeOrder defines the expected order of attributes in an output block.
var outputAttributeOrder = []string{"description", "value", "sensitive", "depends_on"}

// Check runs the rule check.
func (r *TerraformOutputAttributeOrderRule) Check(runner tflint.Runner) error {
	body, err := runner.GetModuleContent(&hclext.BodySchema{
		Blocks: []hclext.BlockSchema{
			{
				Type:       "output",
				LabelNames: []string{"name"},
				Body: &hclext.BodySchema{
					Attributes: []hclext.AttributeSchema{
						{Name: "description"},
						{Name: "value"},
						{Name: "sensitive"},
						{Name: "depends_on"},
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
		outputName := block.Labels[0]

		if err := r.checkAttributeOrder(runner, outputName, block); err != nil {
			return err
		}
	}

	return nil
}

func (r *TerraformOutputAttributeOrderRule) checkAttributeOrder(runner tflint.Runner, outputName string, block *hclext.Block) error {
	// Build a map of attribute name to line number
	// Pre-allocate with expected capacity to avoid map growth
	attrLines := make(map[string]int, len(outputAttributeOrder))
	for name, attr := range block.Body.Attributes {
		attrLines[name] = attr.Range.Start.Line
	}

	// Check ordering between consecutive pairs
	var lastAttr string
	var lastLine int

	for _, attrName := range outputAttributeOrder {
		line, exists := attrLines[attrName]
		if !exists {
			continue
		}

		if lastAttr != "" && line < lastLine {
			if err := runner.EmitIssue(
				r,
				fmt.Sprintf("In output %q, %q should come before %q", outputName, attrName, lastAttr),
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
