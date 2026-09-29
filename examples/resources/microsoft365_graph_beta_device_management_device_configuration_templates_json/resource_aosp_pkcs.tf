# Android AOSP PKCS Certificate.

resource "microsoft365_graph_beta_device_management_device_configuration_templates_json" "aosp_pkcs" {
  display_name       = "Android AOSP PKCS Certificate"
  description        = "Template configuration example"
  role_scope_tag_ids = ["0"]
  settings = jsonencode({
    "@odata.type"                                 = "#microsoft.graph.aospDeviceOwnerPkcsCertificateProfile"
    "deviceManagementApplicabilityRuleOsEdition"  = null
    "deviceManagementApplicabilityRuleOsVersion"  = null
    "deviceManagementApplicabilityRuleDeviceMode" = null
    "renewalThresholdPercentage"                  = 30
    "subjectNameFormat"                           = "custom"
    "certificateValidityPeriodValue"              = 1
    "certificateValidityPeriodScale"              = "years"
    "subjectAlternativeNameType"                  = null
    "certificationAuthority"                      = "ca.example.invalid"
    "certificationAuthorityName"                  = "Provider Test CA"
    "certificationAuthorityType"                  = "microsoft"
    "certificateTemplateName"                     = "ProviderTest"
    "subjectAlternativeNameFormatString"          = null
    "subjectNameFormatString"                     = "CN={{DeviceId}}"
    "certificateStore"                            = "user"
    "extendedKeyUsages" = [{
      "name"             = "Client Authentication"
      "objectIdentifier" = "1.3.6.1.5.5.7.3.2"
    }]
    "customSubjectAlternativeNames" = []
  })
}
