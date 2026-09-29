resource "microsoft365_graph_beta_device_management_device_configuration_templates_json" "test" {
  display_name       = "unit-test-android-work-gmail"
  description        = "Updated template description"
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
