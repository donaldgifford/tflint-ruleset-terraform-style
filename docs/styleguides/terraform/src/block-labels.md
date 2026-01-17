# Block Labels, Variables, and Outputs should be snake case

The label for blocks should be in snake case. E.g. example_instance , not ExampleInstance or example-instance.

Labels are the strings that follow block names. For example, in the following, aws_instance and example_instance are labels for the resource block.

```hcl
resource "aws_instance" "example_instance" {
  # Omitted for brevity
}
```

This includes variables and outputs as well:

```hcl
variable "vpc_id" {}
output "instance_name" {}
```
