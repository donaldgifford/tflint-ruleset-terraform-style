# Introduction

Styles are the conventions that govern our code. The term style is a bit of a
misnomer, since these conventions cover far more than just source file
formatting—terraform fmt handles that for us.

The goal of this guide is to manage this complexity by describing in detail the
Dos and Don'ts of writing Terraform code. These rules exist to keep the code
base manageable while still allowing engineers to use Terraform language features
productively.

This documents idiomatic conventions in Terraform code that we follow. A lot
of these are general guidelines for Terraform. As a starting point we follow [Hashicorp's Terraform Style Guide](https://developer.hashicorp.com/terraform/language/style), while others extend upon external
resources.

1. [Gruntworks Terraform Style Guide](https://docs.gruntwork.io/guides/style/terraform-style-guide/)
2. [Google's Terraform Style Guide](https://cloud.google.com/docs/terraform/best-practices/general-style-structure)
3. [Microsoft's Terraform Style Guide](https://microsoft.github.io/code-with-engineering-playbook/code-reviews/recipes/terraform/)

All code should be error-free when run through `terraform fmt` and `terraform
validate`. We recommend setting up your editor to:

- Run `terraform fmt` on save
- Run `tflint` and `terraform validate` to check for errors
