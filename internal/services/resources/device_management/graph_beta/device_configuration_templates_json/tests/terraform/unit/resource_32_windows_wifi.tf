resource "microsoft365_graph_beta_device_management_device_configuration_templates_json" "test" {
  display_name       = "unit-test-windows-wifi"
  description        = "Device configuration template test"
  role_scope_tag_ids = ["0"]
  settings = jsonencode({
    "@odata.type"                                 = "#microsoft.graph.windowsWifiConfiguration",
    "deviceManagementApplicabilityRuleOsEdition"  = null,
    "deviceManagementApplicabilityRuleOsVersion"  = null,
    "deviceManagementApplicabilityRuleDeviceMode" = null,
    "wifiSecurityType"                            = "wpaPersonal",
    "meteredConnectionLimit"                      = null,
    "ssid"                                        = "Provider-Unassigned-Test",
    "networkName"                                 = "Provider WiFi Study",
    "connectAutomatically"                        = false,
    "connectToPreferredNetwork"                   = null,
    "connectWhenNetworkNameIsHidden"              = null,
    "proxySetting"                                = null,
    "proxyManualAddress"                          = null,
    "proxyManualPort"                             = null,
    "proxyAutomaticConfigurationUrl"              = null,
    "forceFIPSCompliance"                         = null,
    "preSharedKey"                                = "Synthetic-Wire-Study-Only-2026"
  })
}
