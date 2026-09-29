resource "random_string" "suffix" {
  length  = 8
  special = false
  upper   = false
}

resource "microsoft365_graph_beta_device_management_device_configuration_templates_json" "test" {
  display_name       = "acc-test-windows-updates-${random_string.suffix.result}"
  description        = "Device configuration template test"
  role_scope_tag_ids = ["0"]
  settings = jsonencode({
    "@odata.type"                                 = "#microsoft.graph.windowsUpdateForBusinessConfiguration",
    "deviceManagementApplicabilityRuleOsEdition"  = null,
    "deviceManagementApplicabilityRuleOsVersion"  = null,
    "deviceManagementApplicabilityRuleDeviceMode" = null,
    "deliveryOptimizationMode"                    = "userDefined",
    "prereleaseFeatures"                          = "userDefined",
    "automaticUpdateMode"                         = "userDefined",
    "microsoftUpdateServiceAllowed"               = false,
    "driversExcluded"                             = true,
    "installationSchedule"                        = null,
    "qualityUpdatesDeferralPeriodInDays"          = 0,
    "featureUpdatesDeferralPeriodInDays"          = 7,
    "qualityUpdatesPaused"                        = false,
    "featureUpdatesPaused"                        = false,
    "businessReadyUpdatesOnly"                    = "userDefined",
    "skipChecksBeforeRestart"                     = false,
    "updateWeeks"                                 = null,
    "featureUpdatesRollbackWindowInDays"          = null,
    "engagedRestartDeadlineInDays"                = null,
    "engagedRestartSnoozeScheduleInDays"          = null,
    "engagedRestartTransitionScheduleInDays"      = null,
    "deadlineForFeatureUpdatesInDays"             = null,
    "deadlineForQualityUpdatesInDays"             = null,
    "deadlineGracePeriodInDays"                   = null,
    "postponeRebootUntilAfterDeadline"            = null,
    "autoRestartNotificationDismissal"            = "notConfigured",
    "scheduleRestartWarningInHours"               = null,
    "scheduleImminentRestartWarningInMinutes"     = null,
    "userPauseAccess"                             = "notConfigured",
    "userWindowsUpdateScanAccess"                 = "notConfigured",
    "updateNotificationLevel"                     = "notConfigured",
    "allowWindows11Upgrade"                       = false
  })
}
