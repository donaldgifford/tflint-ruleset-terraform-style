# This is a valid hash comment
# Another valid comment

resource "aws_instance" "web" {
  ami           = "ami-12345678"
  instance_type = "t2.micro"

  # Inline hash comment is fine
  tags = {
    Name = "web-server" # Trailing comment
  }
}

# Multiple line comments
# are also perfectly fine
# when using the hash symbol
