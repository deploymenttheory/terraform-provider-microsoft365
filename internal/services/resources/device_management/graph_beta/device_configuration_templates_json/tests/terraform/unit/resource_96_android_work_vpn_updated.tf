resource "microsoft365_graph_beta_device_management_device_configuration_templates_json" "test" {
  display_name       = "unit-test-android-work-vpn"
  description        = "Updated template description"
  role_scope_tag_ids = ["0"]
  settings = jsonencode({
    "@odata.type"                                 = "#microsoft.graph.androidWorkProfileVpnConfiguration"
    "deviceManagementApplicabilityRuleOsEdition"  = null
    "deviceManagementApplicabilityRuleOsVersion"  = null
    "deviceManagementApplicabilityRuleDeviceMode" = null
    "connectionName"                              = "Updated VPN"
    "connectionType"                              = "ciscoAnyConnect"
    "role"                                        = null
    "realm"                                       = null
    "fingerprint"                                 = null
    "authenticationMethod"                        = "usernameAndPassword"
    "proxyServer"                                 = null
    "targetedPackageIds"                          = []
    "alwaysOn"                                    = null
    "alwaysOnLockdown"                            = null
    "lockdownExclusionList"                       = []
    "microsoftTunnelSiteId"                       = null
    "proxyExclusionList"                          = []
    "servers" = [{
      "description"     = "Example VPN"
      "address"         = "vpn.example.invalid"
      "isDefaultServer" = true
    }]
    "customData"         = []
    "customKeyValueData" = []
    "targetedMobileApps" = []
  })
}
