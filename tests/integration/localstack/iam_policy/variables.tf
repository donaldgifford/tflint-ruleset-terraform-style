variable "policy_name" {
  description = "Name of the IAM policy"
  type        = string
  default     = "test-read-access"
}

variable "policy_description" {
  description = "Description for the IAM policy"
  type        = string
  default     = "Policy granting read access for testing"
}
