# macOS Wired Network.

resource "microsoft365_graph_beta_device_management_device_configuration_templates_json" "macos_wired_network" {
  display_name       = "macOS Wired Network"
  description        = "Template configuration example"
  role_scope_tag_ids = ["0"]
  settings = jsonencode({
    "@odata.type"                                 = "#microsoft.graph.macOSWiredNetworkConfiguration"
    "deviceManagementApplicabilityRuleOsEdition"  = null
    "deviceManagementApplicabilityRuleOsVersion"  = null
    "deviceManagementApplicabilityRuleDeviceMode" = null
    "networkName"                                 = "Updated wired network"
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
