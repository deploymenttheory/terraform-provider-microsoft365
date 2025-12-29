resource "random_id" "test_001" {
  byte_length = 4
}

resource "microsoft365_graph_beta_device_management_group_policy_configuration" "test_001" {
  display_name = "acc-test-group-policy-config-boolean-001-${random_id.test_001.hex}"
  description  = "Acceptance test configuration for boolean value minimal scenario"

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}

resource "microsoft365_graph_beta_device_management_group_policy_boolean_value" "test_001" {
  group_policy_configuration_id = microsoft365_graph_beta_device_management_group_policy_configuration.test_001.id
  policy_name                   = "Allow Cloud Policy Management"
  class_type                    = "machine"
  category_path                 = "\\FSLogix\\Profile Containers"
  enabled                       = true

  values = [
    {
      value = true
    }
  ]

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}

