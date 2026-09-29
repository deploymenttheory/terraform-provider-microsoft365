resource "microsoft365_graph_beta_device_management_device_configuration_templates_json" "test" {
  display_name       = "unit-test-windows-edition-upgrade"
  description        = "Device configuration template test"
  role_scope_tag_ids = ["0"]
  settings = jsonencode({
    "@odata.type"                                 = "#microsoft.graph.editionUpgradeConfiguration"
    "deviceManagementApplicabilityRuleOsEdition"  = null
    "deviceManagementApplicabilityRuleOsVersion"  = null
    "deviceManagementApplicabilityRuleDeviceMode" = null
    "licenseType"                                 = "notConfigured"
    "targetEdition"                               = "notConfigured"
    "license"                                     = null
    "productKey"                                  = null
    "windowsSMode"                                = "block"
  })
}
