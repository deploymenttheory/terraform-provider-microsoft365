# iOS / iPadOS PKCS Certificate.

resource "microsoft365_graph_beta_device_management_device_configuration_templates_json" "ios_pkcs_certificate" {
  display_name       = "iOS / iPadOS PKCS Certificate"
  description        = "iOS / iPadOS pkcs certificate example"
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
