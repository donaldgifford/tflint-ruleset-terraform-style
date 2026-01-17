# Don't store secrets in state

There are many resources and data providers in Terraform that store secret values in plaintext in the state file. Where possible, avoid storing secrets in state. Following are some examples of providers that store secrets in plaintext:
