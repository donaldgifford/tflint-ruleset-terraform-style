# Good resource names that don't repeat the type

resource "aws_s3_bucket" "logs" {
  bucket = "my-logs-bucket"
}

resource "aws_s3_bucket" "artifacts" {
  bucket = "my-artifacts"
}

resource "aws_iam_policy" "team_access" {
  name = "team-access-policy"
  policy = jsonencode({
    Version = "2012-10-17"
    Statement = []
  })
}

resource "aws_iam_policy" "read_only" {
  name = "read-only-access"
  policy = jsonencode({
    Version = "2012-10-17"
    Statement = []
  })
}

resource "aws_instance" "web" {
  ami           = "ami-12345678"
  instance_type = "t2.micro"
}

# Using "main" and "this" is allowed
resource "aws_vpc" "main" {
  cidr_block = "10.0.0.0/16"
}

resource "aws_security_group" "this" {
  name   = "example"
  vpc_id = aws_vpc.main.id
}

# Data sources with good names
data "aws_ami" "latest" {
  most_recent = true
  owners      = ["amazon"]
}

data "aws_iam_policy" "admin" {
  name = "AdministratorAccess"
}
