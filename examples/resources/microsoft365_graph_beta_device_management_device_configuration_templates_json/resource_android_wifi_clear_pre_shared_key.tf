# Android Clear a Wi-Fi Pre-shared Key.

resource "microsoft365_graph_beta_device_management_device_configuration_templates_json" "android_wifi_clear_pre_shared_key" {
  display_name       = "Android Clear a Wi-Fi Pre-shared Key"
  description        = "Android clear a wi-fi pre-shared key example"
  role_scope_tag_ids = ["0"]
  settings = jsonencode({
    "@odata.type"                                 = "#microsoft.graph.androidDeviceOwnerWiFiConfiguration",
    "deviceManagementApplicabilityRuleOsEdition"  = null,
    "deviceManagementApplicabilityRuleOsVersion"  = null,
    "deviceManagementApplicabilityRuleDeviceMode" = null,
    "networkName"                                 = "Provider WiFi Study",
    "ssid"                                        = "Provider-Unassigned-Test",
    "connectAutomatically"                        = false,
    "connectWhenNetworkNameIsHidden"              = null,
    "wiFiSecurityType"                            = "wpaPersonal",
    "preSharedKey"                                = null,
    "proxySettings"                               = "none",
    "proxyManualAddress"                          = null,
    "proxyManualPort"                             = null,
    "proxyAutomaticConfigurationUrl"              = null,
    "proxyExclusionList"                          = null,
    "macAddressRandomizationMode"                 = null
  })
}
