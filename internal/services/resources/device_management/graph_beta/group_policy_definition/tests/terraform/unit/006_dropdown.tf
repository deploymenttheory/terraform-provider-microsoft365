resource "microsoft365_graph_beta_device_management_group_policy_definition" "test_006" {
  group_policy_configuration_id = "config-006"
  policy_name                   = "Test Policy Dropdown"
  class_type                    = "machine"
  category_path                 = "\\Test\\Dropdown"
  enabled                       = true

  values = [
    {
      label = "Security Level"
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
