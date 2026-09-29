resource "random_string" "suffix" {
  length  = 8
  special = false
  upper   = false
}

resource "microsoft365_graph_beta_device_management_device_configuration_templates_json" "test" {
  display_name       = "acc-test-macos-updates-${random_string.suffix.result}"
  description        = "Device configuration template test"
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
