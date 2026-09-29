# iOS / iPadOS Imported PKCS Certificate.

resource "microsoft365_graph_beta_device_management_device_configuration_templates_json" "ios_imported_pkcs" {
  display_name       = "iOS / iPadOS Imported PKCS Certificate"
  description        = "Template configuration example"
  role_scope_tag_ids = ["0"]
  settings = jsonencode({
    "@odata.type"                                 = "#microsoft.graph.iosImportedPFXCertificateProfile"
    "deviceManagementApplicabilityRuleOsEdition"  = null
    "deviceManagementApplicabilityRuleOsVersion"  = null
    "deviceManagementApplicabilityRuleDeviceMode" = null
    "intendedPurpose"                             = "smimeSigning"
  })
}
