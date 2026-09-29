# Windows Email.

resource "microsoft365_graph_beta_device_management_device_configuration_templates_json" "windows_email" {
  display_name       = "Windows Email"
  description        = "Template configuration example"
  role_scope_tag_ids = ["0"]
  settings = jsonencode({
    "@odata.type"                                 = "#microsoft.graph.windows10EasEmailProfileConfiguration"
    "deviceManagementApplicabilityRuleOsEdition"  = null
    "deviceManagementApplicabilityRuleOsVersion"  = null
    "deviceManagementApplicabilityRuleDeviceMode" = null
    "usernameSource"                              = "userPrincipalName"
    "usernameAADSource"                           = null
    "userDomainNameSource"                        = null
    "customDomainName"                            = null
    "accountName"                                 = "Example email"
    "syncCalendar"                                = true
    "syncContacts"                                = false
    "syncTasks"                                   = false
    "durationOfEmailToSync"                       = "userDefined"
    "emailAddressSource"                          = "userPrincipalName"
    "emailSyncSchedule"                           = "userDefined"
    "hostName"                                    = "mail.example.invalid"
    "requireSsl"                                  = true
  })
}
