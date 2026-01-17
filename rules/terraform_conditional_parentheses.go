package rules

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// TerraformConditionalParenthesesRule checks that multi-line conditional expressions are wrapped in parentheses.
type TerraformConditionalParenthesesRule struct {
	tflint.DefaultRule
}

// NewTerraformConditionalParenthesesRule returns a new rule instance.
func NewTerraformConditionalParenthesesRule() *TerraformConditionalParenthesesRule {
	return &TerraformConditionalParenthesesRule{}
}

// Name returns the rule name.
func (r *TerraformConditionalParenthesesRule) Name() string {
	return "terraform_conditional_parentheses"
}

// Enabled returns whether the rule is enabled by default.
func (r *TerraformConditionalParenthesesRule) Enabled() bool {
	return true
}

// Severity returns the rule severity.
func (r *TerraformConditionalParenthesesRule) Severity() tflint.Severity {
	return tflint.WARNING
}

// Link returns the rule documentation link.
func (r *TerraformConditionalParenthesesRule) Link() string {
	return "https://github.com/donaldgifford/tflint-ruleset-terraform-style/blob/main/docs/rules/terraform_conditional_parentheses.md"
}

// Check runs the rule check.
func (r *TerraformConditionalParenthesesRule) Check(runner tflint.Runner) error {
	// Check locals blocks
	body, err := runner.GetModuleContent(&hclext.BodySchema{
		Blocks: []hclext.BlockSchema{
			{
				Type: "locals",
				Body: &hclext.BodySchema{
					Mode: hclext.SchemaJustAttributesMode,
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
		for _, attr := range block.Body.Attributes {
			if err := r.checkExpression(runner, attr.Expr, false); err != nil {
				return err
			}
		}
	}

	// Check variable defaults
	varBody, err := runner.GetModuleContent(&hclext.BodySchema{
		Blocks: []hclext.BlockSchema{
			{
				Type:       "variable",
				LabelNames: []string{"name"},
				Body: &hclext.BodySchema{
					Attributes: []hclext.AttributeSchema{
						{Name: "default"},
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

	for _, block := range varBody.Blocks {
		if attr, exists := block.Body.Attributes["default"]; exists {
			if err := r.checkExpression(runner, attr.Expr, false); err != nil {
				return err
			}
		}
	}

	return nil
}

func (r *TerraformConditionalParenthesesRule) checkExpression(runner tflint.Runner, expr hcl.Expression, inParens bool) error {
	// Unwrap to get the underlying syntax expression
	syntaxExpr, ok := expr.(hclsyntax.Expression)
	if !ok {
		return nil
	}

	return r.walkExpression(runner, syntaxExpr, inParens)
}

func (r *TerraformConditionalParenthesesRule) walkExpression(runner tflint.Runner, expr hclsyntax.Expression, inParens bool) error {
	switch e := expr.(type) {
	case *hclsyntax.ConditionalExpr:
		// Check if this conditional spans multiple lines
		exprRange := e.Range()
		if exprRange.Start.Line != exprRange.End.Line && !inParens {
			if err := runner.EmitIssue(
				r,
				"Multi-line conditional expressions should be wrapped in parentheses",
				exprRange,
			); err != nil {
				return err
			}
		}
		// Check nested expressions - propagate inParens context to nested conditionals
		// within a parenthesized expression
		if err := r.walkExpression(runner, e.Condition, inParens); err != nil {
			return err
		}
		if err := r.walkExpression(runner, e.TrueResult, inParens); err != nil {
			return err
		}
		if err := r.walkExpression(runner, e.FalseResult, inParens); err != nil {
			return err
		}

	case *hclsyntax.ParenthesesExpr:
		// Mark that we're inside parentheses
		if err := r.walkExpression(runner, e.Expression, true); err != nil {
			return err
		}

	case *hclsyntax.TupleConsExpr:
		for _, elem := range e.Exprs {
			if err := r.walkExpression(runner, elem, false); err != nil {
				return err
			}
		}

	case *hclsyntax.ObjectConsExpr:
		for _, item := range e.Items {
			if err := r.walkExpression(runner, item.KeyExpr, false); err != nil {
				return err
			}
			if err := r.walkExpression(runner, item.ValueExpr, false); err != nil {
				return err
			}
		}

	case *hclsyntax.FunctionCallExpr:
		for _, arg := range e.Args {
			if err := r.walkExpression(runner, arg, false); err != nil {
				return err
			}
		}

	case *hclsyntax.IndexExpr:
		if err := r.walkExpression(runner, e.Collection, false); err != nil {
			return err
		}
		if err := r.walkExpression(runner, e.Key, false); err != nil {
			return err
		}

	case *hclsyntax.SplatExpr:
		if err := r.walkExpression(runner, e.Source, false); err != nil {
			return err
		}

	case *hclsyntax.ForExpr:
		if err := r.walkExpression(runner, e.CollExpr, false); err != nil {
			return err
		}
		if e.KeyExpr != nil {
			if err := r.walkExpression(runner, e.KeyExpr, false); err != nil {
				return err
			}
		}
		if err := r.walkExpression(runner, e.ValExpr, false); err != nil {
			return err
		}
		if e.CondExpr != nil {
			if err := r.walkExpression(runner, e.CondExpr, false); err != nil {
				return err
			}
		}

	case *hclsyntax.BinaryOpExpr:
		if err := r.walkExpression(runner, e.LHS, false); err != nil {
			return err
		}
		if err := r.walkExpression(runner, e.RHS, false); err != nil {
			return err
		}

	case *hclsyntax.UnaryOpExpr:
		if err := r.walkExpression(runner, e.Val, false); err != nil {
			return err
		}

	case *hclsyntax.TemplateExpr:
		for _, part := range e.Parts {
			if err := r.walkExpression(runner, part, false); err != nil {
				return err
			}
		}

	case *hclsyntax.TemplateWrapExpr:
		if err := r.walkExpression(runner, e.Wrapped, false); err != nil {
			return err
		}
	}

	return nil
}
