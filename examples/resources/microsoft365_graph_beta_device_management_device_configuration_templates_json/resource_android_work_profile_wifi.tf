# Android Work Profile Wi-Fi.

resource "microsoft365_graph_beta_device_management_device_configuration_templates_json" "android_work_profile_wifi" {
  display_name       = "Android Work Profile Wi-Fi"
  description        = "Android work profile wi-fi example"
  role_scope_tag_ids = ["0"]
  settings = jsonencode({
    "@odata.type"                                 = "#microsoft.graph.androidWorkProfileWiFiConfiguration",
    "deviceManagementApplicabilityRuleOsEdition"  = null,
    "deviceManagementApplicabilityRuleOsVersion"  = null,
    "deviceManagementApplicabilityRuleDeviceMode" = null,
    "networkName"                                 = "Provider WiFi Study",
    "ssid"                                        = "Provider-Unassigned-Test",
    "connectAutomatically"                        = false,
    "connectWhenNetworkNameIsHidden"              = false,
    "wiFiSecurityType"                            = "wpaPersonal",
    "preSharedKey"                                = "Synthetic-Wire-Study-Only-2026",
    "proxySettings"                               = "none",
    "proxyAutomaticConfigurationUrl"              = null
  })
}
