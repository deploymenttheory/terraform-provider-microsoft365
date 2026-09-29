resource "microsoft365_graph_beta_device_management_device_configuration_templates_json" "test" {
  display_name       = "unit-test-windows-shared-device"
  description        = "Updated template description"
  role_scope_tag_ids = ["0"]
  settings = jsonencode({
    "@odata.type"                                 = "#microsoft.graph.sharedPCConfiguration"
    "deviceManagementApplicabilityRuleOsEdition"  = null
    "deviceManagementApplicabilityRuleOsVersion"  = null
    "deviceManagementApplicabilityRuleDeviceMode" = null
    "accountManagerPolicy"                        = null
    "allowedAccounts"                             = "guest,domain"
    "localStorage"                                = "notConfigured"
    "allowLocalStorage"                           = false
    "setAccountManager"                           = "disabled"
    "disableAccountManager"                       = true
    "setEduPolicies"                              = "notConfigured"
    "disableEduPolicies"                          = false
    "setPowerPolicies"                            = "notConfigured"
    "disablePowerPolicies"                        = false
    "signInOnResume"                              = "notConfigured"
    "disableSignInOnResume"                       = false
    "enabled"                                     = true
    "idleTimeBeforeSleepInSeconds"                = null
    "kioskAppDisplayName"                         = null
    "kioskAppUserModelId"                         = null
    "maintenanceStartTime"                        = null
    "fastFirstSignIn"                             = "notConfigured"
  })
}
