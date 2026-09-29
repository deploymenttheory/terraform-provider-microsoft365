# Windows Wired Network.

resource "microsoft365_graph_beta_device_management_device_configuration_templates_json" "windows_wired_network" {
  display_name       = "Windows Wired Network"
  description        = "Template configuration example"
  role_scope_tag_ids = ["0"]
  settings = jsonencode({
    "@odata.type"                                 = "#microsoft.graph.windowsWiredNetworkConfiguration"
    "deviceManagementApplicabilityRuleOsEdition"  = null
    "deviceManagementApplicabilityRuleOsVersion"  = null
    "deviceManagementApplicabilityRuleDeviceMode" = null
    "authenticationType"                          = "machineOrUser"
    "cacheCredentials"                            = null
    "authenticationPeriodInSeconds"               = 30
    "authenticationRetryDelayPeriodInSeconds"     = 1
    "eapolStartPeriodInSeconds"                   = 5
    "maximumEAPOLStartMessages"                   = 3
    "maximumAuthenticationFailures"               = 2
    "enforce8021X"                                = true
    "authenticationBlockPeriodInMinutes"          = null
    "eapType"                                     = "peap"
    "trustedServerCertificateNames"               = []
    "authenticationMethod"                        = "usernameAndPassword"
    "secondaryAuthenticationMethod"               = null
    "innerAuthenticationProtocolForEAPTTLS"       = null
    "outerIdentityPrivacyTemporaryValue"          = null
    "performServerValidation"                     = null
    "disableUserPromptForServerValidation"        = null
    "requireCryptographicBinding"                 = null
    "forceFIPSCompliance"                         = null
  })
}
