resource "microsoft365_graph_beta_device_management_group_policy_definition" "test_001" {
  group_policy_configuration_id = "config-001"
  policy_name                   = "Test Policy Boolean Minimal"
  class_type                    = "machine"
  category_path                 = "\\Test\\Boolean\\Minimal"
  enabled                       = true

  values = [
    {
      label = "Enable Feature One"
      value = "true"
    },
    {
      label = "Enable Feature Two"
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

