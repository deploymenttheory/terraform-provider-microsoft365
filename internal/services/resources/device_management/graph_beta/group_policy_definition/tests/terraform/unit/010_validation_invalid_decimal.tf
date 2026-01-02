resource "microsoft365_graph_beta_device_management_group_policy_definition" "test_010" {
  group_policy_configuration_id = "config-004"
  policy_name                   = "Test Policy Decimal"
  class_type                    = "machine"
  category_path                 = "\\Test\\Decimal"
  enabled                       = true

  values = [
    {
      label = "Timeout Setting"
      value = "not-a-number"
    }
  ]
}
