resource "microsoft365_graph_beta_device_management_group_policy_boolean_value" "test_003" {
  group_policy_configuration_id = "00000000-0000-0000-0000-000000000003"
  policy_name                   = "Test Policy Lifecycle"
  class_type                    = "machine"
  category_path                 = "\\Test\\Category\\Lifecycle"
  enabled                       = true

  values = [
    {
      presentation_id = "presentation-003-1"
      value           = true
    },
    {
      presentation_id = "presentation-003-2"
      value           = false
    },
    {
      presentation_id = "presentation-003-3"
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

