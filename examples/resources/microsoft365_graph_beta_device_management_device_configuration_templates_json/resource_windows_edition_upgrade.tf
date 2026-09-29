# Windows Edition Upgrade and S Mode.

resource "microsoft365_graph_beta_device_management_device_configuration_templates_json" "windows_edition_upgrade" {
  display_name       = "Windows Edition Upgrade and S Mode"
  description        = "Template configuration example"
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
    "windowsSMode"                                = "unlock"
  })
}
