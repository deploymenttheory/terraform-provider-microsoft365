# Example: Group Policy Definition with Boolean (CheckBox) presentation values
# This example demonstrates configuring a policy with multiple boolean checkboxes

resource "microsoft365_graph_beta_device_management_group_policy_configuration" "example" {
  display_name = "Example Group Policy Configuration"
  description  = "Configuration for managing Windows Store packages"
}

resource "microsoft365_graph_beta_device_management_group_policy_definition" "boolean_example" {
  group_policy_configuration_id = microsoft365_graph_beta_device_management_group_policy_configuration.example.id
  policy_name                   = "Remove Default Microsoft Store packages from the system."
  class_type                    = "machine"
  category_path                 = "\\Windows Components\\App Package Deployment"
  enabled                       = true

  values = [
    {
      label = "Microsoft Teams"
      value = "true"
    },
    {
      label = "Paint"
      value = "false"
    }
  ]

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}

