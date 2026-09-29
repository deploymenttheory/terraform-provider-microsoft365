resource "random_string" "suffix" {
  length  = 8
  special = false
  upper   = false
}

resource "microsoft365_graph_beta_device_management_device_configuration_templates_json" "test" {
  display_name       = "acc-test-ios-nested-folder-${random_string.suffix.result}"
  description        = "Device configuration template test"
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
