# Module Structure

Each module repo should have the following 3 folders:

- modules: Terraform modules that are designed to be consumed by users. The intention is that users should pull the modules in the modules folder in their terraform code using module blocks.

- examples: Folder that contains top level Terraform modules that provide an example of how to use the modules in the modules folder. The examples folder often has subfolders for-learning-and-testing and for-production that contain corresponding example code. See Testing: Terratest for more info on how these examples should be organized.

- test: Terratest Go files for testing the code in the repo. See Testing: Terratest for specific conventions around Terratest.

Additionally, each module in modules should be organized with the following files:

- variables.tf: All variable blocks should go in here and they specify the inputs.

- outputs.tf: All output blocks should go in here and they specify the outputs.

- main.tf: All other logic should be added here.

- dependencies.tf (optional): Any external references that are pulled in by a data source block should go in here. This allows consumers of the module to quickly scan for what resources need to already exist to deploy the module.

Any nonstandard file structure should be called out in the README (e.g., if main.tf is split up into multiple smaller terraform files).
