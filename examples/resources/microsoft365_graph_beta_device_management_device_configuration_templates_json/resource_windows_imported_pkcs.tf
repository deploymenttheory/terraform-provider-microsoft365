# Windows Imported PKCS Certificate.

resource "microsoft365_graph_beta_device_management_device_configuration_templates_json" "windows_imported_pkcs" {
  display_name       = "Windows Imported PKCS Certificate"
  description        = "Template configuration example"
  role_scope_tag_ids = ["0"]
  settings = jsonencode({
    "@odata.type"                                 = "#microsoft.graph.windows10ImportedPFXCertificateProfile"
    "deviceManagementApplicabilityRuleOsEdition"  = null
    "deviceManagementApplicabilityRuleOsVersion"  = null
    "deviceManagementApplicabilityRuleDeviceMode" = null
    "renewalThresholdPercentage"                  = 20
    "keyStorageProvider"                          = "useSoftwareKsp"
    "subjectNameFormat"                           = "commonName"
    "subjectAlternativeNameType"                  = "userPrincipalName"
    "certificateValidityPeriodValue"              = 1
    "certificateValidityPeriodScale"              = "years"
    "intendedPurpose"                             = "smimeSigning"
  })
}
