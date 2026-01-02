# Example: Group Policy Definition with MultiText (MultiTextBox) presentation value
# This example demonstrates configuring a policy with a multi-line text input field

resource "microsoft365_graph_beta_device_management_group_policy_configuration" "example" {
  display_name = "Example Group Policy Configuration"
  description  = "Configuration for filesystem filter drivers"
}

resource "microsoft365_graph_beta_device_management_group_policy_definition" "multitext_example" {
  group_policy_configuration_id = microsoft365_graph_beta_device_management_group_policy_configuration.example.id
  policy_name                   = "Dev drive filter attach policy"
  class_type                    = "machine"
  category_path                 = "\\System\\Filesystem"
  enabled                       = true

  values = [
    {
      label = "Filter list"
      value = "FilterDriver1\nFilterDriver2\nFilterDriver3"
    }
  ]

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}

