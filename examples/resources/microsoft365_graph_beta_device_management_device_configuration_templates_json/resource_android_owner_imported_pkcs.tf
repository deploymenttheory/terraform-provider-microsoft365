# Android Enterprise Imported PKCS Certificate.

resource "microsoft365_graph_beta_device_management_device_configuration_templates_json" "android_owner_imported_pkcs" {
  display_name       = "Android Enterprise Imported PKCS Certificate"
  description        = "Template configuration example"
  role_scope_tag_ids = ["0"]
  settings = jsonencode({
    "@odata.type"                                 = "#microsoft.graph.androidDeviceOwnerImportedPFXCertificateProfile"
    "deviceManagementApplicabilityRuleOsEdition"  = null
    "deviceManagementApplicabilityRuleOsVersion"  = null
    "deviceManagementApplicabilityRuleDeviceMode" = null
    "renewalThresholdPercentage"                  = 50
    "subjectNameFormat"                           = "commonName"
    "certificateValidityPeriodValue"              = 1
    "certificateValidityPeriodScale"              = "years"
    "subjectAlternativeNameType"                  = "emailAddress"
    "intendedPurpose"                             = "smimeSigning"
    "certificateAccessType"                       = null
    "extendedKeyUsages" = [{
      "name"             = "Any Purpose"
      "objectIdentifier" = "2.5.29.37.0"
    }]
    "silentCertificateAccessDetails" = []
  })
}
