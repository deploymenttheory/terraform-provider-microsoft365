resource "microsoft365_graph_beta_device_management_group_policy_boolean_value" "test_004" {
  group_policy_configuration_id = "00000000-0000-0000-0000-000000000004"
  policy_name                   = "Test Policy Downgrade"
  class_type                    = "user"
  category_path                 = "\\Test\\Category\\Downgrade"
  enabled                       = false

  values = [
    {
      presentation_id = "presentation-004-1"
      value           = false
    }
  ]

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}

