resource "microsoft365_graph_beta_device_management_group_policy_boolean_value" "test_005" {
  group_policy_configuration_id = "00000000-0000-0000-0000-000000000005"
  # policy_name is missing - should cause validation error
  class_type    = "machine"
  category_path = "\\Test\\Category"
  enabled       = true

  values = [
    {
      presentation_id = "presentation-005"
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

