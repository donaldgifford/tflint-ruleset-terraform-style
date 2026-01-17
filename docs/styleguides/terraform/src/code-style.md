Writing Terraform code in a consistent style makes it easier to read and maintain. The following sections discuss code style recommendations, including the following:

- Run terraform fmt and terraform validate before committing your code to version control.
- Use a linter such as TFLint to enforce your organization's own coding best practices.
- Use # for single and multi-line comments.
- Use nouns for resource names and do not include the resource type in the name.
- Use underscores to separate multiple words in names. Wrap the resource type and name in double quotes in your resource definition.
- Let your code build on itself: define dependent resources after the resources they reference.
- Include a type and description for every variable.
- Include a description for every output.
- Avoid overuse of variables and local values.
- Always include a default provider configuration.
- Use count and for_each sparingly.
