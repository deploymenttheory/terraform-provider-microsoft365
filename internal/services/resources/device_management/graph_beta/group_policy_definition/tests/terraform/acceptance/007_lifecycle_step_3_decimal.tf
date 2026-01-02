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
  policy_name                   = "Configure time out for detections in non-critical failed state"
  class_type                    = "machine"
  category_path                 = "\\Windows Components\\Microsoft Defender Antivirus\\Reporting"
  enabled                       = true

  values = [
    {
      label = "Configure time out for detections in non-critical failed state"
      value = "7200"
    }
  ]

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}

