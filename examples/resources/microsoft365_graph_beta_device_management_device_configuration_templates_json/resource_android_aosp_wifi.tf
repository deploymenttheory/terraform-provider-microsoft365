# Android AOSP Wi-Fi.

resource "microsoft365_graph_beta_device_management_device_configuration_templates_json" "android_aosp_wifi" {
  display_name       = "Android AOSP Wi-Fi"
  description        = "Android aosp wi-fi example"
  role_scope_tag_ids = ["0"]
  settings = jsonencode({
    "@odata.type"                                 = "#microsoft.graph.aospDeviceOwnerWiFiConfiguration",
    "deviceManagementApplicabilityRuleOsEdition"  = null,
    "deviceManagementApplicabilityRuleOsVersion"  = null,
    "deviceManagementApplicabilityRuleDeviceMode" = null,
    "networkName"                                 = "Provider WiFi Study",
    "ssid"                                        = "Provider-Unassigned-Test",
    "connectAutomatically"                        = false,
    "connectWhenNetworkNameIsHidden"              = null,
    "wiFiSecurityType"                            = "wpaPersonal",
    "preSharedKey"                                = "Synthetic-Wire-Study-Only-2026",
    "proxySetting"                                = "none",
    "proxyManualAddress"                          = null,
    "proxyManualPort"                             = null,
    "proxyAutomaticConfigurationUrl"              = null,
    "proxyExclusionList"                          = []
  })
}
