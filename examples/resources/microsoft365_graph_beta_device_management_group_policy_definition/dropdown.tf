# Example: Group Policy Definition with Dropdown (DropdownList) presentation value
# This example demonstrates configuring a policy with a dropdown selection field

resource "microsoft365_graph_beta_device_management_group_policy_configuration" "example" {
  display_name = "Example Group Policy Configuration"
  description  = "Configuration for Internet Explorer security settings"
}

resource "microsoft365_graph_beta_device_management_group_policy_definition" "dropdown_example" {
  group_policy_configuration_id = microsoft365_graph_beta_device_management_group_policy_configuration.example.id
  policy_name                   = "Navigate windows and frames across different domains"
  class_type                    = "machine"
  category_path                 = "\\Windows Components\\Internet Explorer\\Internet Control Panel\\Security Page\\Internet Zone"
  enabled                       = true

  values = [
    {
      label = "Navigate windows and frames across different domains"
      value = "1"
    }
  ]

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}

