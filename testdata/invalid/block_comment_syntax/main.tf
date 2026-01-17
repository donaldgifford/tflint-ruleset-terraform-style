/* This is a block comment that should fail */
resource "aws_instance" "test" {
  ami           = "ami-12345678"
  instance_type = "t2.micro"
}

/*
 * Multi-line block comment
 * This should also fail
 */
resource "aws_vpc" "main" {
  cidr_block = "10.0.0.0/16"
}

resource "aws_s3_bucket" "logs" {
  bucket = "my-bucket" /* inline block comment */
}
