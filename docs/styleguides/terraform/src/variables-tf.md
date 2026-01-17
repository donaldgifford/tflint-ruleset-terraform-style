# Variables tf

## Variable Blocks

Each variable block should always define a description and type, even if it is of the string type (the default), in that order. E.g.:

```hcl
variable "example" {
  description = "This is an example"
  type        = string
  default     = "example"  # NOTE: this is optional
}
```

## Complex types

Prefer concrete objects ([object type](https://www.terraform.io/docs/configuration/types.html#structural-types)) over free form maps. However, for particularly large objects it is useful to support optional attributes. This is currently [not supported in terraform](https://github.com/hashicorp/terraform/issues/22449), so workaround by using any type.

When using any type, always use comments to describe the supported attributes. Example.
