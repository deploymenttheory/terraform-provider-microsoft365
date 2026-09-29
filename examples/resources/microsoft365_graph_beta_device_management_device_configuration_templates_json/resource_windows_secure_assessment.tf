# Windows Secure Assessment.

resource "microsoft365_graph_beta_device_management_device_configuration_templates_json" "windows_secure_assessment" {
  display_name       = "Windows Secure Assessment"
  description        = "Template configuration example"
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
