package main

import (
	"github.com/donaldgifford/tflint-ruleset-terraform-style/rules"
	"github.com/terraform-linters/tflint-plugin-sdk/plugin"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
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
