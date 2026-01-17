# Variables with correct attribute order: description -> type -> default

variable "environment" {
  description = "The deployment environment"
  type        = string
  default     = "development"
}

variable "instance_count" {
  description = "Number of instances to create"
  type        = number
  default     = 1
}

variable "tags" {
  description = "Tags to apply to resources"
  type        = map(string)
  default     = {}
}

# Without default is also fine
variable "vpc_id" {
  description = "The VPC ID to deploy into"
  type        = string
}

# Only description is fine
variable "project_name" {
  description = "Name of the project"
}

# With validation block
variable "instance_type" {
  description = "EC2 instance type"
  type        = string
  default     = "t2.micro"

  validation {
    condition     = can(regex("^t[23]\\.", var.instance_type))
    error_message = "Instance type must be t2 or t3 series."
  }
}

# With sensitive
variable "database_password" {
  description = "Database password"
  type        = string
  sensitive   = true
}
