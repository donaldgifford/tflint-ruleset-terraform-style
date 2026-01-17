# Conditionals with correct formatting

variable "enabled" {
  description = "Enable the feature"
  type        = bool
  default     = true
}

variable "environment" {
  description = "Environment name"
  type        = string
  default     = "dev"
}

locals {
  # Single line - no parens needed
  simple_value = var.enabled ? "yes" : "no"

  # Another single line
  env_suffix = var.environment == "prod" ? "" : "-${var.environment}"

  # Multi-line with parentheses - correct
  complex_value = (
    var.enabled
    ? "feature-enabled"
    : "feature-disabled"
  )

  # Nested conditionals with parens
  nested_value = (
    var.environment == "prod"
    ? (
      var.enabled
      ? "prod-enabled"
      : "prod-disabled"
    )
    : "non-prod"
  )

  # Complex expression with parens
  instance_type = (
    var.environment == "prod"
    ? "t3.large"
    : var.environment == "staging"
    ? "t3.medium"
    : "t3.small"
  )
}
