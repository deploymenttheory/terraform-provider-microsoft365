resource "microsoft365_graph_beta_device_management_group_policy_definition" "test_012b" {
  group_policy_configuration_id = "config-001"
  policy_name                   = "Test Policy"
  # class_type is missing
  category_path = "\\Test\\Missing"
  enabled       = true

  values = [
    { label = "Test", value = "true" }
  ]
}
