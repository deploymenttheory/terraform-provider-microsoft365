# iOS / iPadOS Wi-Fi.

resource "microsoft365_graph_beta_device_management_device_configuration_templates_json" "ios_wifi" {
  display_name       = "iOS / iPadOS Wi-Fi"
  description        = "iOS / iPadOS wi-fi example"
  role_scope_tag_ids = ["0"]
  settings = jsonencode({
    "@odata.type"                                 = "#microsoft.graph.iosWiFiConfiguration",
    "deviceManagementApplicabilityRuleOsEdition"  = null,
    "deviceManagementApplicabilityRuleOsVersion"  = null,
    "deviceManagementApplicabilityRuleDeviceMode" = null,
    "networkName"                                 = "Provider WiFi Study",
    "ssid"                                        = "Provider-Unassigned-Test",
    "connectAutomatically"                        = false,
    "connectWhenNetworkNameIsHidden"              = false,
    "wiFiSecurityType"                            = "wpaPersonal",
    "proxySettings"                               = "none",
    "proxyManualAddress"                          = null,
    "proxyManualPort"                             = null,
    "proxyAutomaticConfigurationUrl"              = null,
    "disableMacAddressRandomization"              = null,
    "preSharedKey"                                = "Synthetic-Wire-Study-Only-2026"
  })
}
