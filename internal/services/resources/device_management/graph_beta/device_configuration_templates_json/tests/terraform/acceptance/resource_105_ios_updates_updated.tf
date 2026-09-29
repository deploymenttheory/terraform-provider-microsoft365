resource "random_string" "suffix" {
  length  = 8
  special = false
  upper   = false
}

# iOS update template lifecycle; no devices are assigned by this test.
resource "microsoft365_graph_beta_device_management_device_configuration_templates_json" "test" {
  display_name       = "acc-test-ios-updates-${random_string.suffix.result}"
  description        = "Updated template description"
  role_scope_tag_ids = ["0"]
  settings = jsonencode({
    "@odata.type"                                 = "#microsoft.graph.iosUpdateConfiguration"
    "deviceManagementApplicabilityRuleOsEdition"  = null
    "deviceManagementApplicabilityRuleOsVersion"  = null
    "deviceManagementApplicabilityRuleDeviceMode" = null
    "isEnabled"                                   = true
    "activeHoursStart"                            = "00:00:00.0000000"
    "activeHoursEnd"                              = "00:00:00.0000000"
    "desiredOsVersion"                            = null
    "scheduledInstallDays"                        = []
    "utcTimeOffsetInMinutes"                      = null
    "enforcedSoftwareUpdateDelayInDays"           = null
    "updateScheduleType"                          = "updateOutsideOfActiveHours"
    "customUpdateTimeWindows"                     = []
  })
}
