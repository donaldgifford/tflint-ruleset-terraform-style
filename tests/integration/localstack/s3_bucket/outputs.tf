output "bucket_id" {
  description = "The ID of the S3 bucket"
  value       = aws_s3_bucket.logs.id
}

output "bucket_arn" {
  description = "The ARN of the S3 bucket"
  value       = aws_s3_bucket.logs.arn
}

output "bucket_domain_name" {
  description = "The bucket domain name"
  value       = aws_s3_bucket.logs.bucket_domain_name
}
