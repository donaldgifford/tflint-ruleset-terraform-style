package rules

import (
	"fmt"
	"sort"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// TerraformResourceParameterOrderRule checks that resource parameters are in the correct order.
// Expected order: count/for_each -> regular parameters -> nested blocks -> lifecycle -> depends_on
type TerraformResourceParameterOrderRule struct {
	tflint.DefaultRule
}

// NewTerraformResourceParameterOrderRule returns a new rule instance.
func NewTerraformResourceParameterOrderRule() *TerraformResourceParameterOrderRule {
	return &TerraformResourceParameterOrderRule{}
}

// Name returns the rule name.
func (r *TerraformResourceParameterOrderRule) Name() string {
	return "terraform_resource_parameter_order"
}

// Enabled returns whether the rule is enabled by default.
func (r *TerraformResourceParameterOrderRule) Enabled() bool {
	return true
}

// Severity returns the rule severity.
func (r *TerraformResourceParameterOrderRule) Severity() tflint.Severity {
	return tflint.WARNING
}

// Link returns the rule documentation link.
func (r *TerraformResourceParameterOrderRule) Link() string {
	return "https://github.com/donaldgifford/tflint-ruleset-terraform-style/blob/main/docs/rules/terraform_resource_parameter_order.md"
}

// Category constants for ordering
const (
	categoryMetaArg   = 1 // count, for_each, provider
	categoryAttribute = 2 // regular attributes
	categoryBlock     = 3 // nested blocks (except lifecycle)
	categoryLifecycle = 4 // lifecycle block
	categoryDependsOn = 5 // depends_on
)

// metaArguments are special arguments that should come first
var metaArguments = map[string]bool{
	"count":    true,
	"for_each": true,
	"provider": true,
}

// Check runs the rule check.
func (r *TerraformResourceParameterOrderRule) Check(runner tflint.Runner) error {
	files, err := runner.GetFiles()
	if err != nil {
		return err
	}

	for _, file := range files {
		syntaxBody, ok := file.Body.(*hclsyntax.Body)
		if !ok {
			continue
		}

		for _, block := range syntaxBody.Blocks {
			if block.Type != "resource" && block.Type != "data" {
				continue
			}

			if err := r.checkBlock(runner, block); err != nil {
				return err
			}
		}
	}

	return nil
}

type orderedItem struct {
	name     string
	category int
	line     int
	rng      hcl.Range
}

func (r *TerraformResourceParameterOrderRule) checkBlock(runner tflint.Runner, block *hclsyntax.Block) error {
	var items []orderedItem

	// Collect attributes
	for name, attr := range block.Body.Attributes {
		cat := categoryAttribute
		if metaArguments[name] {
			cat = categoryMetaArg
		} else if name == "depends_on" {
			cat = categoryDependsOn
		}
		items = append(items, orderedItem{
			name:     name,
			category: cat,
			line:     attr.SrcRange.Start.Line,
			rng:      attr.SrcRange,
		})
	}

	// Collect nested blocks
	for _, nestedBlock := range block.Body.Blocks {
		cat := categoryBlock
		if nestedBlock.Type == "lifecycle" {
			cat = categoryLifecycle
		}
		items = append(items, orderedItem{
			name:     nestedBlock.Type,
			category: cat,
			line:     nestedBlock.TypeRange.Start.Line,
			rng:      nestedBlock.TypeRange,
		})
	}

	// Sort by line number to get actual order
	sort.Slice(items, func(i, j int) bool {
		return items[i].line < items[j].line
	})

	// Check if items are in the correct category order
	var prevCategory int
	var prevName string
	for _, item := range items {
		if prevCategory > item.category {
			var msg string
			switch item.category {
			case categoryMetaArg:
				msg = fmt.Sprintf("%q should be at the top of the resource block", item.name)
			case categoryAttribute:
				msg = fmt.Sprintf("attribute %q should come before %q", item.name, prevName)
			case categoryBlock:
				msg = fmt.Sprintf("block %q should come before %q", item.name, prevName)
			case categoryLifecycle:
				msg = "lifecycle block should come before depends_on"
			default:
				msg = fmt.Sprintf("%q is in the wrong position", item.name)
			}

			if err := runner.EmitIssue(r, msg, item.rng); err != nil {
				return err
			}
		}
		prevCategory = item.category
		prevName = item.name
	}

	return nil
}
