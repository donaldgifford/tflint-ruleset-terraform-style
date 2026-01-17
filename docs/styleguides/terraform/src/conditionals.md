# Conditionals

Use () to break up conditionals across multiple lines.

Examples:

```hcl
locals {
  elb_id = (
    var.elb_already_exists
    ? var.elb_id
    : module.elb.elb_id
  )

  excluded_child_account_ids = (
    var.config_create_account_rules
    ? []
    : [
      for account_name, account in module.organization.child_accounts
      : account.id if lookup(lookup(var.child_accounts, account_name, {}), "enable_config_rules", false) == false
    ]
  )
}
```
