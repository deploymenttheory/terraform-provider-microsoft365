# Android Work Profile Gmail Email.

resource "microsoft365_graph_beta_device_management_device_configuration_templates_json" "android_work_gmail" {
  display_name       = "Android Work Profile Gmail Email"
  description        = "Template configuration example"
  role_scope_tag_ids = ["0"]
  settings = jsonencode({
    "@odata.type"                                 = "#microsoft.graph.androidWorkProfileGmailEasConfiguration"
    "deviceManagementApplicabilityRuleOsEdition"  = null
    "deviceManagementApplicabilityRuleOsVersion"  = null
    "deviceManagementApplicabilityRuleDeviceMode" = null
    "authenticationMethod"                        = "usernameAndPassword"
    "durationOfEmailToSync"                       = "userDefined"
    "emailAddressSource"                          = "userPrincipalName"
    "hostName"                                    = "updated-mail.example.invalid"
    "requireSsl"                                  = true
    "usernameSource"                              = "username"
  })
}
