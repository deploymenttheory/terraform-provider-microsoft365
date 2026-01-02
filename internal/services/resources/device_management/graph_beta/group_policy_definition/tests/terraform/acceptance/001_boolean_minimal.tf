resource "random_id" "test_001" {
  byte_length = 4
}

resource "microsoft365_graph_beta_device_management_group_policy_configuration" "test_001" {
  display_name = "acc-test-gpd-bool-min-${random_id.test_001.hex}"
  description  = "Acceptance test for boolean minimal"

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}

resource "microsoft365_graph_beta_device_management_group_policy_definition" "test_001" {
  group_policy_configuration_id = microsoft365_graph_beta_device_management_group_policy_configuration.test_001.id
  policy_name                   = "Remove Default Microsoft Store packages from the system."
  class_type                    = "machine"
  category_path                 = "\\Windows Components\\App Package Deployment"
  enabled                       = true

  values = [
    {
      label = "Microsoft Teams"
      value = "true"
    },
    {
      label = "Paint"
      value = "false"
    }
  ]

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}
