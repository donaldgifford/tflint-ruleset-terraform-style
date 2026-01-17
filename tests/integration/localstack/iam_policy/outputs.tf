output "policy_arn" {
  description = "The ARN of the IAM policy"
  value       = aws_iam_policy.read_access.arn
}

output "policy_id" {
  description = "The policy ID"
  value       = aws_iam_policy.read_access.policy_id
}
