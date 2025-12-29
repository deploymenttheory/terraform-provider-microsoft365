resource "microsoft365_graph_beta_device_management_group_policy_boolean_value" "test_002" {
  group_policy_configuration_id = "00000000-0000-0000-0000-000000000002"
  policy_name                   = "Test Policy Maximal"
  class_type                    = "user"
  category_path                 = "\\Test\\Category\\Maximal"
  enabled                       = true

  values = [
    {
      presentation_id = "presentation-002-1"
      value           = true
    },
    {
      presentation_id = "presentation-002-2"
      value           = false
    },
    {
      presentation_id = "presentation-002-3"
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

