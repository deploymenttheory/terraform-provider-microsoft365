# This legacy template is deprecated in the Intune UI. Use Settings Catalog for new deployments.
# macOS Extensions (Deprecated UI Template).

resource "microsoft365_graph_beta_device_management_device_configuration_templates_json" "macos_extensions" {
  display_name       = "macOS Extensions (Deprecated UI Template)"
  description        = "Template configuration example"
  role_scope_tag_ids = ["0"]
  settings = jsonencode({
    "@odata.type"                                 = "#microsoft.graph.macOSExtensionsConfiguration"
    "deviceManagementApplicabilityRuleOsEdition"  = null
    "deviceManagementApplicabilityRuleOsVersion"  = null
    "deviceManagementApplicabilityRuleDeviceMode" = null
    "kernelExtensionOverridesAllowed"             = true
    "kernelExtensionAllowedTeamIdentifiers"       = []
    "systemExtensionsBlockOverride"               = false
    "systemExtensionsAllowedTeamIdentifiers"      = []
    "kernelExtensionsAllowed"                     = []
    "systemExtensionsAllowed"                     = []
    "systemExtensionsAllowedTypes"                = []
  })
}
