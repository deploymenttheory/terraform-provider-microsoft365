# Windows PKCS Certificate.

resource "microsoft365_graph_beta_device_management_device_configuration_templates_json" "windows_pkcs_certificate" {
  display_name       = "Windows PKCS Certificate"
  description        = "Windows pkcs certificate example"
  role_scope_tag_ids = ["0"]
  settings = jsonencode({
    "@odata.type"                                 = "#microsoft.graph.windows10PkcsCertificateProfile",
    "deviceManagementApplicabilityRuleOsEdition"  = null,
    "deviceManagementApplicabilityRuleOsVersion"  = null,
    "deviceManagementApplicabilityRuleDeviceMode" = null,
    "renewalThresholdPercentage"                  = 20,
    "keyStorageProvider"                          = "useSoftwareKsp",
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
    "extendedKeyUsages" = [
      {
        "name"             = "Client Authentication",
        "objectIdentifier" = "1.3.6.1.5.5.7.3.2"
      }
    ],
    "customSubjectAlternativeNames" = []
  })
}
