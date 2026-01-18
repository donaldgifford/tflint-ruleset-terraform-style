package rules

import (
	"fmt"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// maxExpressionDepth limits recursion depth to prevent stack overflow on deeply nested expressions.
const maxExpressionDepth = 100

// TerraformConditionalParenthesesRule checks that multi-line conditional expressions are wrapped in parentheses.
type TerraformConditionalParenthesesRule struct {
	tflint.DefaultRule
}

// NewTerraformConditionalParenthesesRule returns a new rule instance.
func NewTerraformConditionalParenthesesRule() *TerraformConditionalParenthesesRule {
	return &TerraformConditionalParenthesesRule{}
}

// Name returns the rule name.
func (*TerraformConditionalParenthesesRule) Name() string {
	return "terraform_conditional_parentheses"
}

// Enabled returns whether the rule is enabled by default.
func (*TerraformConditionalParenthesesRule) Enabled() bool {
	return true
}

// Severity returns the rule severity.
func (*TerraformConditionalParenthesesRule) Severity() tflint.Severity {
	return tflint.WARNING
}

// Link returns the rule documentation link.
func (*TerraformConditionalParenthesesRule) Link() string {
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

	return r.walkExpression(runner, syntaxExpr, inParens, 0)
}

func (r *TerraformConditionalParenthesesRule) walkExpression(runner tflint.Runner, expr hclsyntax.Expression, inParens bool, depth int) error {
	if depth > maxExpressionDepth {
		return fmt.Errorf("expression nesting too deep (max %d)", maxExpressionDepth)
	}

	switch e := expr.(type) {
	case *hclsyntax.ConditionalExpr:
		return r.walkConditionalExpr(runner, e, inParens, depth)
	case *hclsyntax.ParenthesesExpr:
		return r.walkExpression(runner, e.Expression, true, depth+1)
	case *hclsyntax.TupleConsExpr:
		return r.walkExprs(runner, e.Exprs, depth+1)
	case *hclsyntax.ObjectConsExpr:
		return r.walkObjectConsExpr(runner, e, depth+1)
	case *hclsyntax.FunctionCallExpr:
		return r.walkExprs(runner, e.Args, depth+1)
	case *hclsyntax.IndexExpr:
		return r.walkExprs(runner, []hclsyntax.Expression{e.Collection, e.Key}, depth+1)
	case *hclsyntax.SplatExpr:
		return r.walkExpression(runner, e.Source, false, depth+1)
	case *hclsyntax.ForExpr:
		return r.walkForExpr(runner, e, depth+1)
	case *hclsyntax.BinaryOpExpr:
		return r.walkExprs(runner, []hclsyntax.Expression{e.LHS, e.RHS}, depth+1)
	case *hclsyntax.UnaryOpExpr:
		return r.walkExpression(runner, e.Val, false, depth+1)
	case *hclsyntax.TemplateExpr:
		return r.walkExprs(runner, e.Parts, depth+1)
	case *hclsyntax.TemplateWrapExpr:
		return r.walkExpression(runner, e.Wrapped, false, depth+1)
	}
	return nil
}

func (r *TerraformConditionalParenthesesRule) walkConditionalExpr(
	runner tflint.Runner,
	e *hclsyntax.ConditionalExpr,
	inParens bool,
	depth int,
) error {
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
	// Pass through inParens for nested conditionals in TrueResult/FalseResult
	// since they're part of the same parenthesized expression
	for _, expr := range []hclsyntax.Expression{e.Condition, e.TrueResult, e.FalseResult} {
		if err := r.walkExpression(runner, expr, inParens, depth+1); err != nil {
			return err
		}
	}
	return nil
}

func (r *TerraformConditionalParenthesesRule) walkExprs(runner tflint.Runner, exprs []hclsyntax.Expression, depth int) error {
	for _, expr := range exprs {
		if err := r.walkExpression(runner, expr, false, depth); err != nil {
			return err
		}
	}
	return nil
}

func (r *TerraformConditionalParenthesesRule) walkObjectConsExpr(runner tflint.Runner, e *hclsyntax.ObjectConsExpr, depth int) error {
	for _, item := range e.Items {
		if err := r.walkExprs(runner, []hclsyntax.Expression{item.KeyExpr, item.ValueExpr}, depth); err != nil {
			return err
		}
	}
	return nil
}

func (r *TerraformConditionalParenthesesRule) walkForExpr(runner tflint.Runner, e *hclsyntax.ForExpr, depth int) error {
	if err := r.walkExpression(runner, e.CollExpr, false, depth); err != nil {
		return err
	}
	if e.KeyExpr != nil {
		if err := r.walkExpression(runner, e.KeyExpr, false, depth); err != nil {
			return err
		}
	}
	if err := r.walkExpression(runner, e.ValExpr, false, depth); err != nil {
		return err
	}
	if e.CondExpr != nil {
		if err := r.walkExpression(runner, e.CondExpr, false, depth); err != nil {
			return err
		}
	}
	return nil
}
