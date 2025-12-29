resource "random_id" "test_004" {
  byte_length = 4
}

resource "microsoft365_graph_beta_device_management_group_policy_configuration" "test_004" {
  display_name = "acc-test-group-policy-config-boolean-004-${random_id.test_004.hex}"
  description  = "Acceptance test configuration for boolean value downgrade"

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}

resource "microsoft365_graph_beta_device_management_group_policy_boolean_value" "test_004" {
  group_policy_configuration_id = microsoft365_graph_beta_device_management_group_policy_configuration.test_004.id
  policy_name                   = "Enable Profile Containers"
  class_type                    = "machine"
  category_path                 = "\\FSLogix\\Profile Containers"
  enabled                       = true

  values = [
    {
      value = true
    },
    {
      value = false
    },
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

