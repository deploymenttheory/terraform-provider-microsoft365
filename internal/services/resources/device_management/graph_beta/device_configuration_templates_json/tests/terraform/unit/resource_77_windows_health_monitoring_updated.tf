resource "microsoft365_graph_beta_device_management_device_configuration_templates_json" "test" {
  display_name       = "unit-test-windows-health-monitoring"
  description        = "Updated template description"
  role_scope_tag_ids = ["0"]
  settings = jsonencode({
    "@odata.type"                                 = "#microsoft.graph.windowsHealthMonitoringConfiguration"
    "deviceManagementApplicabilityRuleOsEdition"  = null
    "deviceManagementApplicabilityRuleOsVersion"  = null
    "deviceManagementApplicabilityRuleDeviceMode" = null
    "allowDeviceHealthMonitoring"                 = "enabled"
    "configDeviceHealthMonitoringScope"           = "undefined"
    "configDeviceHealthMonitoringCustomScope"     = null
  })
}
