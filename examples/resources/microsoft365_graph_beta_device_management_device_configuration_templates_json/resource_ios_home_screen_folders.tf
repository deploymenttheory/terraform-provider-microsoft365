# iOS / iPadOS Home Screen with Nested Folders.

resource "microsoft365_graph_beta_device_management_device_configuration_templates_json" "ios_home_screen_folders" {
  display_name       = "iOS / iPadOS Home Screen with Nested Folders"
  description        = "iOS / iPadOS home screen with nested folders example"
  role_scope_tag_ids = ["0"]
  settings = jsonencode({
    "@odata.type"                                 = "#microsoft.graph.iosDeviceFeaturesConfiguration",
    "deviceManagementApplicabilityRuleOsEdition"  = null,
    "deviceManagementApplicabilityRuleOsVersion"  = null,
    "deviceManagementApplicabilityRuleDeviceMode" = null,
    "assetTagTemplate"                            = null,
    "contentFilterSettings"                       = null,
    "lockScreenFootnote"                          = null,
    "homeScreenGridWidth"                         = null,
    "homeScreenGridHeight"                        = null,
    "singleSignOnSettings"                        = null,
    "wallpaperDisplayLocation"                    = "notConfigured",
    "wallpaperImage"                              = null,
    "singleSignOnExtension"                       = null,
    "iosSingleSignOnExtension"                    = null,
    "airPrintDestinations"                        = [],
    "homeScreenDockIcons"                         = [],
    "homeScreenPages" = [
      {
        "displayName" = "Page 1",
        "icons" = [
          {
            "@odata.type" = "#microsoft.graph.iosHomeScreenFolder",
            "displayName" = "Apps",
            "pages" = [
              {
                "displayName" = "Folder page",
                "apps" = [
                  {
                    "displayName" = "Safari",
                    "bundleID"    = "com.apple.mobilesafari",
                    "isWebClip"   = false
                  }
                ]
              }
            ]
          }
        ]
      }
    ],
    "notificationSettings" = []
  })
}
