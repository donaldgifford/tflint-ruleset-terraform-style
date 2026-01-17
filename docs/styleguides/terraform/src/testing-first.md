# Use less expensive test methods first

There are multiple methods that you can use to test Terraform. In ascending order of cost, run time, and depth, they include the following:

    Static analysis: Testing the syntax and structure of your configuration without deploying any resources, using tools such as compilers, linters, and dry runs. To do so, use terraform validate
    Module integration testing: To ensure that modules work correctly, test individual modules in isolation. Integration testing for modules involves deploying the module into a test environment and verifying that expected resources are created. There are several testing frameworks that make it easier to write tests, as follows:
        Google's blueprint testing framework
        Terratest
        Kitchen-Terraform
        InSpec
        tftest
    End-to-end testing: By extending the integration testing approach to an entire environment, you can confirm that multiple modules work together. In this approach, deploy all modules that make up the architecture in a fresh test environment. Ideally, the test environment is as similar as possible to your production environment. This is costly but provides the greatest confidence that changes don't break your production environment.
