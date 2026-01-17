// Package main provides the entry point for the tflint-ruleset-terraform-style plugin.
package main

import (
	"github.com/terraform-linters/tflint-plugin-sdk/plugin"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"

	"github.com/donaldgifford/tflint-ruleset-terraform-style/rules"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		RuleSet: &tflint.BuiltinRuleSet{
			Name:    "terraform-style",
			Version: "0.1.0",
			Rules:   rules.Rules,
		},
	})
}
