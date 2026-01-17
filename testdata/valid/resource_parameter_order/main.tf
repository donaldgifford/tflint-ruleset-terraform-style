# Resources with correct parameter order:
# count/for_each -> params -> blocks -> lifecycle -> depends_on

variable "instance_count" {
  description = "Number of instances"
  type        = number
  default     = 1
}

variable "ami_id" {
  description = "AMI ID"
  type        = string
  default     = "ami-12345678"
}

resource "aws_instance" "web" {
  count = var.instance_count

  ami           = var.ami_id
  instance_type = "t2.micro"
  subnet_id     = "subnet-12345678"

  tags = {
    Name = "web-${count.index}"
  }

  root_block_device {
    volume_size = 20
    volume_type = "gp3"
  }

  lifecycle {
    create_before_destroy = true
  }

  depends_on = []
}

# With for_each
resource "aws_iam_user" "users" {
  for_each = toset(["alice", "bob", "charlie"])

  name = each.value
  path = "/users/"

  tags = {
    Team = "engineering"
  }

  lifecycle {
    prevent_destroy = false
  }
}

# Without meta-arguments is fine
resource "aws_s3_bucket" "logs" {
  bucket = "my-logs-bucket"

  tags = {
    Environment = "production"
  }

  lifecycle {
    prevent_destroy = true
  }
}

# Data source with correct order
data "aws_ami" "latest" {
  most_recent = true
  owners      = ["amazon"]

  filter {
    name   = "name"
    values = ["amzn2-ami-hvm-*-x86_64-gp2"]
  }
}
