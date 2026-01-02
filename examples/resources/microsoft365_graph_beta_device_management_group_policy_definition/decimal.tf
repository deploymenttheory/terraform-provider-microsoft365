# Example: Group Policy Definition with Decimal (DecimalTextBox) presentation value
# This example demonstrates configuring a policy with a numeric input field

resource "microsoft365_graph_beta_device_management_group_policy_configuration" "example" {
  display_name = "Example Group Policy Configuration"
  description  = "Configuration for Microsoft Defender Antivirus settings"
}

resource "microsoft365_graph_beta_device_management_group_policy_definition" "decimal_example" {
  group_policy_configuration_id = microsoft365_graph_beta_device_management_group_policy_configuration.example.id
  policy_name                   = "Configure time out for detections in non-critical failed state"
  class_type                    = "machine"
  category_path                 = "\\Windows Components\\Microsoft Defender Antivirus\\Reporting"
  enabled                       = true

  values = [
    {
      label = "Configure time out for detections in non-critical failed state"
      value = "7200"
    }
  ]

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}

