resource "random_string" "suffix" {
  length  = 8
  special = false
  upper   = false
}

resource "microsoft365_graph_beta_device_management_device_configuration_templates_json" "test" {
  display_name       = "acc-test-windows-imported-pkcs-${random_string.suffix.result}"
  description        = "Updated template description"
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
