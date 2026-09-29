# This legacy template is deprecated in the Intune UI. Use Settings Catalog for new deployments.
# macOS Endpoint Protection (Deprecated UI Template).

resource "microsoft365_graph_beta_device_management_device_configuration_templates_json" "macos_endpoint_protection" {
  display_name       = "macOS Endpoint Protection (Deprecated UI Template)"
  description        = "Template configuration example"
  role_scope_tag_ids = ["0"]
  settings = jsonencode({
    "@odata.type"                                          = "#microsoft.graph.macOSEndpointProtectionConfiguration"
    "deviceManagementApplicabilityRuleOsEdition"           = null
    "deviceManagementApplicabilityRuleOsVersion"           = null
    "deviceManagementApplicabilityRuleDeviceMode"          = null
    "gatekeeperAllowedAppSource"                           = "notConfigured"
    "gatekeeperBlockOverride"                              = false
    "firewallEnabled"                                      = true
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
