resource "microsoft365_graph_beta_device_management_group_policy_boolean_value" "test_001" {
  group_policy_configuration_id = "00000000-0000-0000-0000-000000000001"
  policy_name                   = "Test Policy Minimal"
  class_type                    = "machine"
  category_path                 = "\\Test\\Category"
  enabled                       = true

  values = [
    {
      presentation_id = "presentation-001"
      value           = true
    }
  ]

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}

