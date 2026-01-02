# Example: Group Policy Definition with TextBox presentation value
# This example demonstrates configuring a policy with a text input field

resource "microsoft365_graph_beta_device_management_group_policy_configuration" "example" {
  display_name = "Example Group Policy Configuration"
  description  = "Configuration for Microsoft Edge browsing data lifetime"
}

resource "microsoft365_graph_beta_device_management_group_policy_definition" "textbox_example" {
  group_policy_configuration_id = microsoft365_graph_beta_device_management_group_policy_configuration.example.id
  policy_name                   = "Browsing Data Lifetime Settings"
  class_type                    = "machine"
  category_path                 = "\\Microsoft Edge"
  enabled                       = true

  values = [
    {
      label = "Browsing Data Lifetime Settings"
      value = "[{\"data_types\":[\"browsing_history\"],\"time_to_live_in_hours\":168}]"
    }
  ]

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}

