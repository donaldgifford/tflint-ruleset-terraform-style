# Main tf

`main.tf` should (loosely) be organized by sections that correspond to components. There is no standard on grouping, but as a rule of thumb each section should be focused on a specific component of the module. For example, an ECS service module may consist of the following sections:

- The ECS service resource, and any locals logic for setting up the attributes of the resource.
- The ECS task definition resource, and any locals and template logic for setting up the attributes of the resource (e.g. the container definition).
- Any resources related to configuring ELBs to forward traffic to the ECS service (e.g., listeners and target groups).
- Any resources related to configuring IAM permissions for the ECS service.
- Any resources related to configuring network access (e.g., security group rules) for the ECS service.

There is no standard on ordering the sections, but as a rule of thumb the following sections should be placed first, in order:

- Version constraints for the module
- Provider blocks, if needed.
- The main component of the module (e.g., the aws_ecs_service resource for the ECS service module).
- All other sections.
- Any data blocks (at the bottom).
