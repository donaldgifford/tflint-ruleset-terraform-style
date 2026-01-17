# -----------------------------------------------------------------------------
# IAM Policy Resource
# Following all style guide rules
# -----------------------------------------------------------------------------

resource "aws_iam_policy" "read_access" {
  name        = var.policy_name
  description = var.policy_description

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Sid    = "ReadAccess"
        Effect = "Allow"
        Action = [
          "s3:GetObject",
          "s3:ListBucket",
        ]
        Resource = "*"
      }
    ]
  })

  tags = {
    Environment = "test"
    ManagedBy   = "terraform"
  }
}
