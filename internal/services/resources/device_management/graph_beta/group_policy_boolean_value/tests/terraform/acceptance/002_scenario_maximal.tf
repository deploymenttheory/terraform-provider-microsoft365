resource "random_id" "test_002" {
  byte_length = 4
}

resource "microsoft365_graph_beta_device_management_group_policy_configuration" "test_002" {
  display_name = "acc-test-group-policy-config-boolean-002-${random_id.test_002.hex}"
  description  = "Acceptance test configuration for boolean value maximal scenario"

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}

resource "microsoft365_graph_beta_device_management_group_policy_boolean_value" "test_002" {
  group_policy_configuration_id = microsoft365_graph_beta_device_management_group_policy_configuration.test_002.id
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

