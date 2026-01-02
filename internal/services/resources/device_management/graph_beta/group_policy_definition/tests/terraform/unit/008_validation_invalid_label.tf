resource "microsoft365_graph_beta_device_management_group_policy_definition" "test_008" {
  group_policy_configuration_id = "config-001"
  policy_name                   = "Test Policy Boolean Minimal"
  class_type                    = "machine"
  category_path                 = "\\Test\\Boolean\\Minimal"
  enabled                       = true

  values = [
    {
      label = "Invalid Label That Does Not Exist"
      value = "true"
    }
  ]
}
