# Comments variables.tf

variables.tf files should clearly indicate required environment variables, and separate out required variables from optional variables (with defaults) using block comments.

Example:

```hcl
# ---------------------------------------------------------------------------------------------------------------------
# ENVIRONMENT VARIABLES
# Define these secrets as environment variables
# ---------------------------------------------------------------------------------------------------------------------

# TF_VAR_master_password

# ---------------------------------------------------------------------------------------------------------------------
# MODULE PARAMETERS
# These variables are expected to be passed in by the operator
# ---------------------------------------------------------------------------------------------------------------------

variable "required_var" {
  description = "This variable must be set in order to create the resource."
  type        = string
}

# ---------------------------------------------------------------------------------------------------------------------
# OPTIONAL PARAMETERS
# These variables have defaults and may be overridden
# ---------------------------------------------------------------------------------------------------------------------

variable "optional_var" {
  description = "This variable has a sensible default so it is not necessary to set it explicitly for this module to work."
  type        = string
  default     = "Hello world"
}
```
