resource "microsoft365_graph_beta_device_management_group_policy_definition" "test_012" {
  group_policy_configuration_id = "config-001"
  # policy_name is missing
  class_type    = "machine"
  category_path = "\\Test\\Missing"
  enabled       = true

  values = [
    { label = "Test", value = "true" }
  ]
}
