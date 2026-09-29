# Configure your derived credential issuer in Intune before assigning this profile.
# Android Enterprise Derived Credential.

resource "microsoft365_graph_beta_device_management_device_configuration_templates_json" "android_owner_derived_credential" {
  display_name       = "Android Enterprise Derived Credential"
  description        = "Android android enterprise derived credential example"
  role_scope_tag_ids = ["0"]
  settings = jsonencode({
    "@odata.type"                                 = "#microsoft.graph.androidDeviceOwnerDerivedCredentialAuthenticationConfiguration"
    "deviceManagementApplicabilityRuleOsEdition"  = null
    "deviceManagementApplicabilityRuleOsVersion"  = null
    "deviceManagementApplicabilityRuleDeviceMode" = null
    "certificateAccessType"                       = null
    "silentCertificateAccessDetails"              = []
  })
}
