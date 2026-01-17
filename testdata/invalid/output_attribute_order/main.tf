# Outputs with incorrect attribute order

# value before description - FAIL
output "instance_id" {
  value       = "i-12345678"
  description = "The ID of the created instance"
}

# sensitive before value - FAIL
output "database_password" {
  description = "The database password"
  sensitive   = true
  value       = "secret123"
}

# completely wrong order - FAIL
output "vpc_id" {
  sensitive   = false
  value       = "vpc-12345678"
  description = "The VPC ID"
}

# depends_on before sensitive - FAIL
output "bucket_arn" {
  description = "The ARN of the S3 bucket"
  depends_on  = []
  value       = "arn:aws:s3:::my-bucket"
}
