# Configure your derived credential issuer in Intune before assigning this profile.
# iOS / iPadOS Derived Credential.

resource "microsoft365_graph_beta_device_management_device_configuration_templates_json" "ios_derived_credential" {
  display_name       = "iOS / iPadOS Derived Credential"
  description        = "iOS / iPadOS derived credential example"
  role_scope_tag_ids = ["0"]
  settings = jsonencode({
    "@odata.type"                                 = "#microsoft.graph.iosDerivedCredentialAuthenticationConfiguration"
    "deviceManagementApplicabilityRuleOsEdition"  = null
    "deviceManagementApplicabilityRuleOsVersion"  = null
    "deviceManagementApplicabilityRuleDeviceMode" = null
  })
}
