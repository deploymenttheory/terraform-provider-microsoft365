resource "random_string" "suffix" {
  length  = 8
  special = false
  upper   = false
}

resource "microsoft365_graph_beta_device_management_device_configuration_templates_json" "test" {
  display_name       = "acc-test-android-owner-imported-pkcs-${random_string.suffix.result}"
  description        = "Device configuration template test"
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
    "intendedPurpose"                             = "smimeEncryption"
    "certificateAccessType"                       = null
    "extendedKeyUsages" = [{
      "name"             = "Any Purpose"
      "objectIdentifier" = "2.5.29.37.0"
    }]
    "silentCertificateAccessDetails" = []
  })
}
