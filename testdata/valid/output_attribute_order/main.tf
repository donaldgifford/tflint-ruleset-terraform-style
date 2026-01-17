# Outputs with correct attribute order: description -> value -> sensitive

output "instance_id" {
  description = "The ID of the created instance"
  value       = "i-12345678"
}

output "vpc_id" {
  description = "The VPC ID"
  value       = "vpc-12345678"
}

# With sensitive flag
output "database_password" {
  description = "The database password"
  value       = "secret123"
  sensitive   = true
}

# With depends_on
output "bucket_arn" {
  description = "The ARN of the S3 bucket"
  value       = "arn:aws:s3:::my-bucket"
  depends_on  = []
}

# Full ordering
output "private_key" {
  description = "The private key for SSH access"
  value       = "-----BEGIN RSA PRIVATE KEY-----"
  sensitive   = true
  depends_on  = []
}
