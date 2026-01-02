resource "microsoft365_graph_beta_device_management_group_policy_definition" "test_009" {
  group_policy_configuration_id = "config-001"
  policy_name                   = "Test Policy Boolean Minimal"
  class_type                    = "machine"
  category_path                 = "\\Test\\Boolean\\Minimal"
  enabled                       = true

  values = [
    {
      label = "Enable Feature One"
      value = "not-a-boolean"
    }
  ]
}
