resource "random_string" "suffix" {
  length  = 8
  special = false
  upper   = false
}

resource "microsoft365_graph_beta_device_management_device_configuration_templates_json" "test" {
  display_name       = "acc-test-macos-device-features-${random_string.suffix.result}"
  description        = "Device configuration template test"
  role_scope_tag_ids = ["0"]
  settings = jsonencode({
    "@odata.type"                                 = "#microsoft.graph.macOSDeviceFeaturesConfiguration"
    "deviceManagementApplicabilityRuleOsEdition"  = null
    "deviceManagementApplicabilityRuleOsVersion"  = null
    "deviceManagementApplicabilityRuleDeviceMode" = null
    "adminShowHostInfo"                           = false
    "loginWindowText"                             = null
    "authorizedUsersListHidden"                   = false
    "authorizedUsersListHideLocalUsers"           = false
    "authorizedUsersListHideMobileAccounts"       = false
    "authorizedUsersListIncludeNetworkUsers"      = false
    "authorizedUsersListHideAdminUsers"           = false
    "authorizedUsersListShowOtherManagedUsers"    = false
    "shutDownDisabled"                            = false
    "restartDisabled"                             = false
    "sleepDisabled"                               = false
    "consoleAccessDisabled"                       = false
    "shutDownDisabledWhileLoggedIn"               = false
    "restartDisabledWhileLoggedIn"                = false
    "powerOffDisabledWhileLoggedIn"               = false
    "logOutDisabledWhileLoggedIn"                 = false
    "screenLockDisableImmediate"                  = false
    "singleSignOnExtension"                       = null
    "macOSSingleSignOnExtension"                  = null
    "contentCachingEnabled"                       = false
    "contentCachingType"                          = "notConfigured"
    "contentCachingMaxSizeBytes"                  = null
    "contentCachingDataPath"                      = null
    "contentCachingDisableConnectionSharing"      = false
    "contentCachingForceConnectionSharing"        = false
    "contentCachingClientPolicy"                  = "notConfigured"
    "contentCachingPeerPolicy"                    = "notConfigured"
    "contentCachingParentSelectionPolicy"         = "notConfigured"
    "contentCachingParents"                       = []
    "contentCachingLogClientIdentities"           = false
    "contentCachingBlockDeletion"                 = false
    "contentCachingShowAlerts"                    = false
    "contentCachingKeepAwake"                     = false
    "contentCachingPort"                          = null
    "airPrintDestinations"                        = []
    "autoLaunchItems"                             = []
    "associatedDomains"                           = []
    "appAssociatedDomains"                        = []
    "contentCachingClientListenRanges"            = []
    "contentCachingPeerListenRanges"              = []
    "contentCachingPeerFilterRanges"              = []
    "contentCachingPublicRanges"                  = []
  })
}
