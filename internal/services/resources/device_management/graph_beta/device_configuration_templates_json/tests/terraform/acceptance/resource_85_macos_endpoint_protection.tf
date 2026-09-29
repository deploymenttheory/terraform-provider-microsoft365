resource "random_string" "suffix" {
  length  = 8
  special = false
  upper   = false
}

resource "microsoft365_graph_beta_device_management_device_configuration_templates_json" "test" {
  display_name       = "acc-test-macos-endpoint-protection-${random_string.suffix.result}"
  description        = "Device configuration template test"
  role_scope_tag_ids = ["0"]
  settings = jsonencode({
    "@odata.type"                                          = "#microsoft.graph.macOSEndpointProtectionConfiguration"
    "deviceManagementApplicabilityRuleOsEdition"           = null
    "deviceManagementApplicabilityRuleOsVersion"           = null
    "deviceManagementApplicabilityRuleDeviceMode"          = null
    "gatekeeperAllowedAppSource"                           = "notConfigured"
    "gatekeeperBlockOverride"                              = false
    "firewallEnabled"                                      = false
    "firewallBlockAllIncoming"                             = false
    "firewallEnableStealthMode"                            = false
    "fileVaultEnabled"                                     = false
    "fileVaultSelectedRecoveryKeyTypes"                    = "notConfigured"
    "fileVaultInstitutionalRecoveryKeyCertificate"         = null
    "fileVaultInstitutionalRecoveryKeyCertificateFileName" = null
    "fileVaultPersonalRecoveryKeyHelpMessage"              = null
    "fileVaultAllowDeferralUntilSignOut"                   = false
    "fileVaultNumberOfTimesUserCanIgnore"                  = null
    "fileVaultDisablePromptAtSignOut"                      = false
    "fileVaultPersonalRecoveryKeyRotationInMonths"         = null
    "fileVaultHidePersonalRecoveryKey"                     = false
    "advancedThreatProtectionRealTime"                     = "notConfigured"
    "advancedThreatProtectionCloudDelivered"               = "notConfigured"
    "advancedThreatProtectionAutomaticSampleSubmission"    = "notConfigured"
    "advancedThreatProtectionDiagnosticDataCollection"     = "notConfigured"
    "advancedThreatProtectionExcludedFolders"              = []
    "advancedThreatProtectionExcludedFiles"                = []
    "advancedThreatProtectionExcludedExtensions"           = []
    "advancedThreatProtectionExcludedProcesses"            = []
    "firewallApplications"                                 = []
  })
}
