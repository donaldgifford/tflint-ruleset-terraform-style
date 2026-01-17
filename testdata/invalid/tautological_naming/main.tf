# Bad resource names that repeat the type

resource "aws_iam_policy" "iam_policy" {
  name = "my-policy"
  policy = jsonencode({
    Version = "2012-10-17"
    Statement = []
  })
}

resource "aws_iam_policy" "team_policy" {
  name = "team-policy"
  policy = jsonencode({
    Version = "2012-10-17"
    Statement = []
  })
}

resource "aws_s3_bucket" "logs_bucket" {
  bucket = "my-logs-bucket"
}

resource "aws_s3_bucket" "s3_artifacts" {
  bucket = "my-artifacts"
}

resource "aws_instance" "web_instance" {
  ami           = "ami-12345678"
  instance_type = "t2.micro"
}

# Data source with tautological name
data "aws_iam_policy" "existing_policy" {
  name = "AmazonEC2ReadOnlyAccess"
}

data "aws_ami" "ami_latest" {
  most_recent = true
  owners      = ["amazon"]
}
