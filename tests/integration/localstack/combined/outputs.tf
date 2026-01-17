output "bucket_id" {
  description = "The S3 bucket ID"
  value       = aws_s3_bucket.artifacts.id
}

output "bucket_arn" {
  description = "The S3 bucket ARN"
  value       = aws_s3_bucket.artifacts.arn
}

output "policy_arn" {
  description = "The IAM policy ARN"
  value       = aws_iam_policy.bucket_access.arn
}

output "versioning_enabled" {
  description = "Whether versioning is enabled"
  value       = var.enable_versioning
}
