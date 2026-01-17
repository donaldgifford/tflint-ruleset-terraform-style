# Outputs.tf

Each output block should always define a description, before the value:

```hcl
output "greeting" {
  description = "This is a greeting for everyone."
  value       = "hello world!"
}
```
