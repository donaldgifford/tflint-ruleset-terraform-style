# TFLint configuration for testing tflint-ruleset-terraform-style
# This disables the default terraform rules to test our plugin in isolation

# Disable built-in terraform rules for testing
plugin "terraform" {
  enabled = false
}

plugin "terraform-style" {
  enabled = true

  # Specify version if using a remote source
  # version = "0.1.0"
  # source  = "github.com/donaldgifford/tflint-ruleset-terraform-style"
}

# Individual rule configuration (all enabled by default)

rule "terraform_block_comment_syntax" {
  enabled = true
}

rule "terraform_comment_delimiter" {
  enabled = true
}

rule "terraform_tautological_naming" {
  enabled = true
}

rule "terraform_variable_attribute_order" {
  enabled = true
}

rule "terraform_output_attribute_order" {
  enabled = true
}

rule "terraform_conditional_parentheses" {
  enabled = true
}

rule "terraform_resource_parameter_order" {
  enabled = true
}
