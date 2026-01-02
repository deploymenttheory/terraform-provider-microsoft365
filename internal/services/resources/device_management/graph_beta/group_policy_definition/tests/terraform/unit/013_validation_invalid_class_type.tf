resource "microsoft365_graph_beta_device_management_group_policy_definition" "test_013" {
  group_policy_configuration_id = "config-001"
  policy_name                   = "Test Policy"
  class_type                    = "invalid_type"
  category_path                 = "\\Test\\Invalid"
  enabled                       = true

  values = [
    { label = "Test", value = "true" }
  ]
}
