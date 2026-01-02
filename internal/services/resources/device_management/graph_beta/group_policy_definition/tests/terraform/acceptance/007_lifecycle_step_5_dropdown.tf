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

