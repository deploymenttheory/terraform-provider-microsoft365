resource "random_id" "test_003" {
  byte_length = 4
}

resource "microsoft365_graph_beta_device_management_group_policy_configuration" "test_003" {
  display_name = "acc-test-gpd-textbox-${random_id.test_003.hex}"
  description  = "Acceptance test for textbox"

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}

resource "microsoft365_graph_beta_device_management_group_policy_definition" "test_003" {
  group_policy_configuration_id = microsoft365_graph_beta_device_management_group_policy_configuration.test_003.id
  policy_name                   = "Browsing Data Lifetime Settings"
  class_type                    = "machine"
  category_path                 = "\\Microsoft Edge"
  enabled                       = true

  values = [
    {
      label = "Browsing Data Lifetime Settings"
      value = "[{\"data_types\":[\"browsing_history\"],\"time_to_live_in_hours\":168}]"
    }
  ]

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}
