// Package rules contains the TFLint rules for Terraform style enforcement.
package rules

import "github.com/terraform-linters/tflint-plugin-sdk/tflint"

// Rules is the list of all rules this plugin provides.
var Rules = []tflint.Rule{
	// Tier 1: Token Analysis
	NewTerraformBlockCommentSyntaxRule(),
	NewTerraformCommentDelimiterRule(),
	// Tier 2: Block/Attribute Analysis
	NewTerraformTautologicalNamingRule(),
	NewTerraformVariableAttributeOrderRule(),
	NewTerraformOutputAttributeOrderRule(),
	// Tier 3: Complex Analysis
	NewTerraformConditionalParenthesesRule(),
	NewTerraformResourceParameterOrderRule(),
}
