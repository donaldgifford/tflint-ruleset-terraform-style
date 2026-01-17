# -----------------------------------------------------------------------------
# Combined Module - S3 + IAM
# Following all style guide rules
# -----------------------------------------------------------------------------

locals {
  # Single line conditional - no parens needed
  bucket_name = "${var.project_name}-${var.environment}"

  # Multi-line conditional with parens
  versioning_status = (
    var.enable_versioning
    ? "Enabled"
    : "Disabled"
  )

  common_tags = {
    Project     = var.project_name
    Environment = var.environment
    ManagedBy   = "terraform"
  }
}

# --- S3 Resources ---

resource "aws_s3_bucket" "artifacts" {
  bucket = local.bucket_name

  tags = local.common_tags
}

resource "aws_s3_bucket_versioning" "artifacts" {
  bucket = aws_s3_bucket.artifacts.id

  versioning_configuration {
    status = local.versioning_status
  }
}

# --- IAM Resources ---

resource "aws_iam_policy" "bucket_access" {
  name        = "${var.project_name}-bucket-access"
  description = "Policy for accessing the ${var.project_name} bucket"

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Sid    = "ListBucket"
        Effect = "Allow"
        Action = [
          "s3:ListBucket",
        ]
        Resource = aws_s3_bucket.artifacts.arn
      },
      {
        Sid    = "ObjectAccess"
        Effect = "Allow"
        Action = [
          "s3:GetObject",
          "s3:PutObject",
          "s3:DeleteObject",
        ]
        Resource = "${aws_s3_bucket.artifacts.arn}/*"
      }
    ]
  })

  tags = local.common_tags

  depends_on = [aws_s3_bucket.artifacts]
}
