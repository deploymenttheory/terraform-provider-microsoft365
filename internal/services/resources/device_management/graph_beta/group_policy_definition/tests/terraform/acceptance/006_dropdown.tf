resource "random_id" "test_006" {
  byte_length = 4
}

resource "microsoft365_graph_beta_device_management_group_policy_configuration" "test_006" {
  display_name = "acc-test-gpd-dropdown-${random_id.test_006.hex}"
  description  = "Acceptance test for dropdown"

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}

resource "microsoft365_graph_beta_device_management_group_policy_definition" "test_006" {
  group_policy_configuration_id = microsoft365_graph_beta_device_management_group_policy_configuration.test_006.id
  policy_name                   = "Navigate windows and frames across different domains"
  class_type                    = "machine"
  category_path                 = "\\Windows Components\\Internet Explorer\\Internet Control Panel\\Security Page\\Internet Zone"
  enabled                       = true

  values = [
    {
      label = "Navigate windows and frames across different domains"
      value = "1"
    }
  ]

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}
