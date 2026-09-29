# macOS Software Updates.

resource "microsoft365_graph_beta_device_management_device_configuration_templates_json" "macos_software_updates" {
  display_name       = "macOS Software Updates"
  description        = "macOS software updates example"
  role_scope_tag_ids = ["0"]
  settings = jsonencode({
    "@odata.type"                                 = "#microsoft.graph.macOSSoftwareUpdateConfiguration",
    "deviceManagementApplicabilityRuleOsEdition"  = null,
    "deviceManagementApplicabilityRuleOsVersion"  = null,
    "deviceManagementApplicabilityRuleDeviceMode" = null,
    "criticalUpdateBehavior"                      = "notConfigured",
    "configDataUpdateBehavior"                    = "notConfigured",
    "firmwareUpdateBehavior"                      = "notConfigured",
    "allOtherUpdateBehavior"                      = "notConfigured",
    "updateScheduleType"                          = "alwaysUpdate",
    "updateTimeWindowUtcOffsetInMinutes"          = null,
    "maxUserDeferralsCount"                       = null,
    "priority"                                    = null,
    "customUpdateTimeWindows"                     = []
  })
}
