resource "random_id" "test_008" {
  byte_length = 4
}

resource "microsoft365_graph_beta_device_management_group_policy_configuration" "test_008" {
  display_name = "acc-test-gpd-errors-${random_id.test_008.hex}"
  description  = "Acceptance test for group policy definition validation errors"

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}

resource "microsoft365_graph_beta_device_management_group_policy_definition" "test_008" {
  group_policy_configuration_id = microsoft365_graph_beta_device_management_group_policy_configuration.test_008.id
  policy_name                   = "Remove Default Microsoft Store packages from the system."
  class_type                    = "machine"
  category_path                 = "\\Windows Components\\App Package Deployment"
  enabled                       = true

  values = [
    {
      label = "Microsoft Teams"
      value = "not-a-boolean"
    }
  ]

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}

