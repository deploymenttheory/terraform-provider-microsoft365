resource "microsoft365_graph_beta_device_management_device_configuration_templates_json" "ios_root" {
  display_name       = "unit-test-ios-root"
  description        = "Device configuration template test"
  role_scope_tag_ids = ["0"]
  settings = jsonencode({
    "@odata.type"                                 = "#microsoft.graph.iosTrustedRootCertificate",
    "deviceManagementApplicabilityRuleOsEdition"  = null,
    "deviceManagementApplicabilityRuleOsVersion"  = null,
    "deviceManagementApplicabilityRuleDeviceMode" = null,
    "trustedRootCertificate"                      = "MIIDFjCCAf6gAwIBAgICAhgwDQYJKoZIhvcNAQELBQAwKzEpMCcGA1UEAxMgVGVycmFmb3JtIE1pY3Jvc29mdDM2NSBUZXN0IFJvb3QwIBcNMjAwMTAxMDAwMDAwWhgPMjA1MDAxMDEwMDAwMDBaMCsxKTAnBgNVBAMTIFRlcnJhZm9ybSBNaWNyb3NvZnQzNjUgVGVzdCBSb290MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAt0FBlUiuQItKihIPdimOfiHkajmQgMvC/tztuj1Sis7rPeChIMVbexXSNHUKBHST6ti6ck91yLy2kUll9HLlM8Y/GczP7CRaIf2xv+8/In7YEuhZJ2E+yA/q8ZQZYRcUXYpYWtuutWBjKv/G1oh2l5IzGVePAn7gRpm79gVDsyw/cuhNtj2CQAGcScUCP2Yee+Dgc0MLrFLCvNtrVc0iosWARBLL2OataXzrSmoSQm29TyJDsZ+Qr8kanHbQICUp2mmyMVKr2uxd0JzkJKxNTrNB8dVhRUytKHz3BE0QkKOz+RFqsLjixx8InoyhQLj6FK4mhu5KPP3P2kdFDQin4QIDAQABo0IwQDAOBgNVHQ8BAf8EBAMCAQYwDwYDVR0TAQH/BAUwAwEB/zAdBgNVHQ4EFgQU+orhr+LWDMSsmhI7Le0bRV984wEwDQYJKoZIhvcNAQELBQADggEBAImrXrkNyTG/UL9fe5jiLyYExt1CfZcz6T0SK6F7zAnvK2ggJfxK92d99u/v0x6OVj2rLNxYE8Uf3tO+2xjo9hwLsO5bXzD9pqc+UeVDGfNiSUWyt4bkvG7typSl0VUbSw1OgeVYI5Gr4byOGZ8ph9DOuc8beQSJ5VlwV+PROK6rtm+QJKiLNeQElZYbqrjgtQvw1P9lNocw0AFcoa19PKIjd8ARf9NHr3GRt1JE5EbCLsjZcwpIV3SKf3cDZ6YmyO4KHPwwDKu/4RMvG2cm+Omqa80/eIKdZYz2yg2kiABGT0B/t9ZcgOOyjNjVwIBB8Ed0nTu+CEc9SVPBZ6B2ITM=",
    "certFileName"                                = "wire-study.cer"
  })
}

resource "microsoft365_graph_beta_device_management_device_configuration_templates_json" "ios_scep" {
  display_name       = "unit-test-ios-scep"
  description        = "Device configuration template test"
  role_scope_tag_ids = ["0"]
  settings = jsonencode({
    "@odata.type"                                 = "#microsoft.graph.iosScepCertificateProfile",
    "deviceManagementApplicabilityRuleOsEdition"  = null,
    "deviceManagementApplicabilityRuleOsVersion"  = null,
    "deviceManagementApplicabilityRuleDeviceMode" = null,
    "renewalThresholdPercentage"                  = 20,
    "subjectNameFormat"                           = "custom",
    "subjectAlternativeNameType"                  = "none",
    "certificateValidityPeriodValue"              = 1,
    "certificateValidityPeriodScale"              = "years",
    "scepServerUrls" = [
      "https://scep.example.invalid/certsrv/mscep/mscep.dll"
    ],
    "subjectNameFormatString"            = "CN={{DeviceId}}",
    "keyUsage"                           = "keyEncipherment,digitalSignature",
    "keySize"                            = "size2048",
    "subjectAlternativeNameFormatString" = null,
    "certificateStore"                   = "user",
    "extendedKeyUsages" = [
      {
        "name"             = "Client Authentication",
        "objectIdentifier" = "1.3.6.1.5.5.7.3.2"
      }
    ],
    "customSubjectAlternativeNames" = [],
    "rootCertificate@odata.bind"    = "https://graph.microsoft.com/beta/deviceManagement/deviceConfigurations('${microsoft365_graph_beta_device_management_device_configuration_templates_json.ios_root.id}')"
  })
}

resource "microsoft365_graph_beta_device_management_device_configuration_templates_json" "test" {
  display_name       = "unit-test-ios-enterprise-wifi"
  description        = "Device configuration template test"
  role_scope_tag_ids = ["0"]
  settings = jsonencode({
    "@odata.type"                                 = "#microsoft.graph.iosEnterpriseWiFiConfiguration",
    "deviceManagementApplicabilityRuleOsEdition"  = null,
    "deviceManagementApplicabilityRuleOsVersion"  = null,
    "deviceManagementApplicabilityRuleDeviceMode" = null,
    "networkName"                                 = "Provider Enterprise WiFi",
    "ssid"                                        = "Provider-Unassigned-EAP",
    "connectAutomatically"                        = false,
    "connectWhenNetworkNameIsHidden"              = false,
    "wiFiSecurityType"                            = "wpaEnterprise",
    "proxySettings"                               = "none",
    "proxyManualAddress"                          = null,
    "proxyManualPort"                             = null,
    "proxyAutomaticConfigurationUrl"              = null,
    "disableMacAddressRandomization"              = null,
    "preSharedKey"                                = null,
    "eapType"                                     = "eapTls",
    "eapFastConfiguration"                        = null,
    "trustedServerCertificateNames"               = [],
    "authenticationMethod"                        = "certificate",
    "innerAuthenticationProtocolForEapTtls"       = null,
    "outerIdentityPrivacyTemporaryValue"          = null,
    "usernameFormatString"                        = null,
    "passwordFormatString"                        = null,
    "rootCertificatesForServerValidation@odata.bind" = [
      "https://graph.microsoft.com/beta/deviceManagement/deviceConfigurations('${microsoft365_graph_beta_device_management_device_configuration_templates_json.ios_root.id}')"
    ],
    "identityCertificateForClientAuthentication@odata.bind" = "https://graph.microsoft.com/beta/deviceManagement/deviceConfigurations('${microsoft365_graph_beta_device_management_device_configuration_templates_json.ios_scep.id}')"
  })
}
