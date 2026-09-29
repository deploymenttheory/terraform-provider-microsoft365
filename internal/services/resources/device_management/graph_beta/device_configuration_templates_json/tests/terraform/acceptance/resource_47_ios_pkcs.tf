resource "random_string" "suffix" {
  length  = 8
  special = false
  upper   = false
}

resource "microsoft365_graph_beta_device_management_device_configuration_templates_json" "test" {
  display_name       = "acc-test-ios-pkcs-${random_string.suffix.result}"
  description        = "Device configuration template test"
  role_scope_tag_ids = ["0"]
  settings = jsonencode({
    "@odata.type"                                 = "#microsoft.graph.iosPkcsCertificateProfile",
    "deviceManagementApplicabilityRuleOsEdition"  = null,
    "deviceManagementApplicabilityRuleOsVersion"  = null,
    "deviceManagementApplicabilityRuleDeviceMode" = null,
    "renewalThresholdPercentage"                  = 20,
    "subjectNameFormat"                           = "custom",
    "subjectAlternativeNameType"                  = "none",
    "certificateValidityPeriodValue"              = 1,
    "certificateValidityPeriodScale"              = "years",
    "certificationAuthority"                      = "ca.example.invalid",
    "certificationAuthorityName"                  = "Provider Test CA",
    "certificateTemplateName"                     = "ProviderTest",
    "subjectAlternativeNameFormatString"          = null,
    "subjectNameFormatString"                     = "CN={{DeviceId}}",
    "certificateStore"                            = "user",
    "customSubjectAlternativeNames"               = []
  })
}
