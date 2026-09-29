resource "microsoft365_graph_beta_device_management_device_configuration_templates_json" "test" {
  display_name       = "unit-test-macos-wired-network"
  description        = "Device configuration template test"
  role_scope_tag_ids = ["0"]
  settings = jsonencode({
    "@odata.type"                                 = "#microsoft.graph.macOSWiredNetworkConfiguration"
    "deviceManagementApplicabilityRuleOsEdition"  = null
    "deviceManagementApplicabilityRuleOsVersion"  = null
    "deviceManagementApplicabilityRuleDeviceMode" = null
    "networkName"                                 = "Example wired network"
    "networkInterface"                            = "anyEthernet"
    "eapType"                                     = "peap"
    "eapFastConfiguration"                        = null
    "trustedServerCertificateNames"               = []
    "authenticationMethod"                        = "usernameAndPassword"
    "nonEapAuthenticationMethodForEapTtls"        = null
    "enableOuterIdentityPrivacy"                  = null
    "deploymentChannel"                           = null
  })
}
