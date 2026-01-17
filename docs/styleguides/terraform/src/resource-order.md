# Resource order

The order of the resources and data sources in your code does not affect how Terraform builds them, so organize your resources for readability. Terraform determines the creation order based on cross-resource dependencies.

How you order your resources largely depends on the size and complexity of your code, but we recommend defining data sources alongside the resources that reference them. For readability, your Terraform code should “build on itself” — you should define a data source before the resource that references it.

The following example defines an aws_instance that relies on two data sources, aws_ami and aws_availability_zone. For readability and continuity, it defines the data sources before the aws_instance resource.

```hcl
data "aws_ami" "web" {
  ##...
}

data "aws_availability_zones" "available" {
  ##...
}

resource "aws_instance" "web" {
  ami               = data.aws_ami.web.id
  availability_zone = data.aws_availability_zones.available.names[0]
  ##...
}
```

We recommend following a consistent order for resource parameters:

1. If present, The `count` or `for_each` meta-argument.
2. Resource-specific non-block parameters.
3. Resource-specific block parameters.
4. If required, a `lifecycle` block.
5. If required, the `depends_on` parameter.
