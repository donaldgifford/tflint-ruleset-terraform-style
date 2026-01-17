variable "project_name" {
  description = "Name of the project"
  type        = string
  default     = "integration-test"
}

variable "environment" {
  description = "Environment name"
  type        = string
  default     = "test"
}

variable "enable_versioning" {
  description = "Enable S3 bucket versioning"
  type        = bool
  default     = true
}
