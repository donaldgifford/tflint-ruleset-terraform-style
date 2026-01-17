# Variables with incorrect attribute order

# type before description - FAIL
variable "environment" {
  type        = string
  description = "The deployment environment"
  default     = "development"
}

# default before type - FAIL
variable "instance_count" {
  description = "Number of instances to create"
  default     = 1
  type        = number
}

# completely reversed - FAIL
variable "tags" {
  default     = {}
  type        = map(string)
  description = "Tags to apply to resources"
}

# sensitive before default - FAIL
variable "api_key" {
  description = "API key for external service"
  type        = string
  sensitive   = true
  default     = ""
}
