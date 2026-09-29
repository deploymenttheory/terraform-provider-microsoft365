# Windows Health Monitoring.

resource "microsoft365_graph_beta_device_management_device_configuration_templates_json" "windows_health_monitoring" {
  display_name       = "Windows Health Monitoring"
  description        = "Template configuration example"
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
