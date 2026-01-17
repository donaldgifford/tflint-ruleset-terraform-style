# Resources with incorrect parameter order

# count not first - FAIL
resource "aws_instance" "web" {
  ami   = "ami-12345678"
  count = 2
  instance_type = "t2.micro"
}

# for_each not first - FAIL
resource "aws_iam_user" "users" {
  name     = each.value
  for_each = toset(["alice", "bob"])
  path     = "/users/"
}

# depends_on before lifecycle - FAIL
resource "aws_s3_bucket" "logs" {
  bucket = "my-logs-bucket"

  depends_on = []

  lifecycle {
    prevent_destroy = true
  }
}

# lifecycle before regular blocks - FAIL
resource "aws_instance" "app" {
  ami           = "ami-12345678"
  instance_type = "t2.micro"

  lifecycle {
    create_before_destroy = true
  }

  root_block_device {
    volume_size = 20
  }
}
