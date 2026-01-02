resource "random_id" "test_004" {
  byte_length = 4
}

resource "microsoft365_graph_beta_device_management_group_policy_configuration" "test_004" {
  display_name = "acc-test-gpd-decimal-${random_id.test_004.hex}"
  description  = "Acceptance test for decimal"

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}

resource "microsoft365_graph_beta_device_management_group_policy_definition" "test_004" {
  group_policy_configuration_id = microsoft365_graph_beta_device_management_group_policy_configuration.test_004.id
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
