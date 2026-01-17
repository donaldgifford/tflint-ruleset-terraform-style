# This fixture is a placeholder.
#
# HCL syntax itself requires parentheses for multi-line ternary expressions.
# Without parentheses, the HCL parser returns a syntax error, so there's no
# valid way to write a "failing" fixture for this rule.
#
# The terraform_conditional_parentheses rule validates that parentheses are
# used, but HCL enforces this at the syntax level.

variable "enabled" {
  description = "Enable the feature"
  type        = bool
  default     = true
}

locals {
  # This single-line conditional is valid and doesn't need parens
  simple_value = var.enabled ? "yes" : "no"
}
