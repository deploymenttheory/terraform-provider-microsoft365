# Android Work Profile Nine Email.

resource "microsoft365_graph_beta_device_management_device_configuration_templates_json" "android_work_nine_email" {
  display_name       = "Android Work Profile Nine Email"
  description        = "Template configuration example"
  role_scope_tag_ids = ["0"]
  settings = jsonencode({
    "@odata.type"                                 = "#microsoft.graph.androidWorkProfileNineWorkEasConfiguration"
    "deviceManagementApplicabilityRuleOsEdition"  = null
    "deviceManagementApplicabilityRuleOsVersion"  = null
    "deviceManagementApplicabilityRuleDeviceMode" = null
    "authenticationMethod"                        = "usernameAndPassword"
    "durationOfEmailToSync"                       = "userDefined"
    "emailAddressSource"                          = "userPrincipalName"
    "hostName"                                    = "mail.example.invalid"
    "requireSsl"                                  = true
    "usernameSource"                              = "username"
    "syncCalendar"                                = true
    "syncContacts"                                = false
    "syncTasks"                                   = false
  })
}
