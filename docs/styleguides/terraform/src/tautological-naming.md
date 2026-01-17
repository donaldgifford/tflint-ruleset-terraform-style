# Avoid tautological resource names

Bad:

```hcl
resource "aws_iam_policy" "ec2_policy"{
  ...
}
```

Good:

```hcl
resource "aws_iam_policy" "team_1_access"{
  ...
}
```
