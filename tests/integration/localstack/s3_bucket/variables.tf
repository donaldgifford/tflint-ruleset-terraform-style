variable "bucket_name" {
  description = "Name of the S3 bucket to create"
  type        = string
  default     = "test-bucket"
}

variable "environment" {
  description = "Environment tag for the bucket"
  type        = string
  default     = "test"
}
