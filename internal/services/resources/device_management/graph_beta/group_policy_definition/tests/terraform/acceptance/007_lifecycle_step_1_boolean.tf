resource "random_id" "test_007" {
  byte_length = 4
}

resource "microsoft365_graph_beta_device_management_group_policy_configuration" "test_007" {
  display_name = "acc-test-gpd-lifecycle-${random_id.test_007.hex}"
  description  = "Acceptance test for lifecycle transitions"

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}

resource "microsoft365_graph_beta_device_management_group_policy_definition" "test_007" {
  group_policy_configuration_id = microsoft365_graph_beta_device_management_group_policy_configuration.test_007.id
  policy_name                   = "Remove Default Microsoft Store packages from the system."
  class_type                    = "machine"
  category_path                 = "\\Windows Components\\App Package Deployment"
  enabled                       = true

  values = [
    { label = "Microsoft Teams", value = "true" },
    { label = "Paint", value = "false" }
  ]

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}
