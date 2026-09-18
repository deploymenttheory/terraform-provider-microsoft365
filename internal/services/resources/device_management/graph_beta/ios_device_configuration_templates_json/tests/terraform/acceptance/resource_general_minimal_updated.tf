
resource "random_string" "json_general_suffix" {
  length  = 8
  special = false
  upper   = false
}

resource "microsoft365_graph_beta_device_management_ios_device_configuration_templates_json" "general_minimal" {
  odata_type   = "#microsoft.graph.iosGeneralDeviceConfiguration"
  display_name = "acc-test-iOS-json-general-${random_string.json_general_suffix.result}"

  settings_json = jsonencode({
    cameraBlocked  = true
    airDropBlocked = true
  })

  timeouts = {
    create = "50s"
    read   = "5m"
    update = "30m"
    delete = "30m"
  }
}
