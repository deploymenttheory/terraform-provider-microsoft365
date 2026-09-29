resource "random_string" "suffix" {
  length  = 8
  special = false
  upper   = false
}

resource "microsoft365_graph_beta_device_management_device_configuration_templates_json" "test" {
  display_name       = "acc-test-windows-kiosk-${random_string.suffix.result}"
  description        = "Updated template description"
  role_scope_tag_ids = ["0"]
  settings = jsonencode({
    "@odata.type"                                 = "#microsoft.graph.windowsKioskConfiguration"
    "deviceManagementApplicabilityRuleOsEdition"  = null
    "deviceManagementApplicabilityRuleOsVersion"  = null
    "deviceManagementApplicabilityRuleDeviceMode" = null
    "kioskBrowserDefaultUrl"                      = null
    "kioskBrowserEnableHomeButton"                = false
    "kioskBrowserEnableNavigationButtons"         = false
    "kioskBrowserEnableEndSessionButton"          = false
    "kioskBrowserRestartOnIdleTimeInMinutes"      = null
    "kioskBrowserBlockedURLs"                     = []
    "kioskBrowserBlockedUrlExceptions"            = []
    "edgeKioskEnablePublicBrowsing"               = false
    "windowsKioskForceUpdateSchedule"             = null
    "kioskProfiles" = [{
      "profileId"   = "71bf2c1e-90d3-4714-b9b0-325c26f433f1"
      "profileName" = "Example kiosk"
      "appConfiguration" = {
        "@odata.type" = "#microsoft.graph.windowsKioskSingleUWPApp"
        "uwpApp" = {
          "startLayoutTileSize" = "hidden"
          "name"                = "Calculator"
          "appType"             = "unknown"
          "autoLaunch"          = false
          "appUserModelId"      = "Microsoft.WindowsCalculator_8wekyb3d8bbwe!App"
          "appId"               = null
          "containedAppId"      = null
        }
      }
      "userAccountsConfiguration" = [{
        "@odata.type" = "#microsoft.graph.windowsKioskAutologon"
      }]
    }]
  })
}
