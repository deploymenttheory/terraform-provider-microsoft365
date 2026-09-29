# Android Work Profile PKCS Certificate.

resource "microsoft365_graph_beta_device_management_device_configuration_templates_json" "android_work_profile_pkcs" {
  display_name       = "Android Work Profile PKCS Certificate"
  description        = "Android work profile pkcs certificate example"
  role_scope_tag_ids = ["0"]
  settings = jsonencode({
    "@odata.type"                                 = "#microsoft.graph.androidWorkProfilePkcsCertificateProfile",
    "deviceManagementApplicabilityRuleOsEdition"  = null,
    "deviceManagementApplicabilityRuleOsVersion"  = null,
    "deviceManagementApplicabilityRuleDeviceMode" = null,
    "renewalThresholdPercentage"                  = 20,
    "subjectNameFormat"                           = "custom",
    "certificateValidityPeriodValue"              = 1,
    "certificateValidityPeriodScale"              = "years",
    "subjectAlternativeNameType"                  = "none",
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
