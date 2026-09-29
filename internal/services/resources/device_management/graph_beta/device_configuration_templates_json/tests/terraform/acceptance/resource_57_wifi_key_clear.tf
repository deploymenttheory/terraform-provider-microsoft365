resource "random_string" "suffix" {
  length  = 8
  special = false
  upper   = false
}

resource "microsoft365_graph_beta_device_management_device_configuration_templates_json" "test" {
  display_name       = "acc-test-android-owner-wifi-${random_string.suffix.result}"
  description        = "Device configuration template test"
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
