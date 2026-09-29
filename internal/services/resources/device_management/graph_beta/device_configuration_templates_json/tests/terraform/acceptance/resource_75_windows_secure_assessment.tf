resource "random_string" "suffix" {
  length  = 8
  special = false
  upper   = false
}

resource "microsoft365_graph_beta_device_management_device_configuration_templates_json" "test" {
  display_name       = "acc-test-windows-secure-assessment-${random_string.suffix.result}"
  description        = "Device configuration template test"
  role_scope_tag_ids = ["0"]
  settings = jsonencode({
    "@odata.type"                                 = "#microsoft.graph.windows10SecureAssessmentConfiguration"
    "deviceManagementApplicabilityRuleOsEdition"  = null
    "deviceManagementApplicabilityRuleOsVersion"  = null
    "deviceManagementApplicabilityRuleDeviceMode" = null
    "launchUri"                                   = "https://assessment.example.invalid"
    "configurationAccount"                        = null
    "configurationAccountType"                    = "localGuestAccount"
    "allowPrinting"                               = true
    "allowScreenCapture"                          = true
    "allowTextSuggestion"                         = true
    "localGuestAccountName"                       = "Assessment"
    "assessmentAppUserModelId"                    = null
  })
}
