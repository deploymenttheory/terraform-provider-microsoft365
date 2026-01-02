resource "random_id" "test_005" {
  byte_length = 4
}

resource "microsoft365_graph_beta_device_management_group_policy_configuration" "test_005" {
  display_name = "acc-test-gpd-multitext-${random_id.test_005.hex}"
  description  = "Acceptance test for multitext"

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}

resource "microsoft365_graph_beta_device_management_group_policy_definition" "test_005" {
  group_policy_configuration_id = microsoft365_graph_beta_device_management_group_policy_configuration.test_005.id
  policy_name                   = "Dev drive filter attach policy"
  class_type                    = "machine"
  category_path                 = "\\System\\Filesystem"
  enabled                       = true

  values = [
    {
      label = "Filter list"
      value = "FilterDriver1\nFilterDriver2\nFilterDriver3"
    }
  ]

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}
