resource "random_string" "json_features_suffix" {
  length  = 8
  special = false
  upper   = false
}


resource "microsoft365_graph_beta_device_management_ios_device_configuration_templates_json" "device_features" {
  odata_type   = "#microsoft.graph.iosDeviceFeaturesConfiguration"
  display_name = "acc-test-iOS-json-features-${random_string.json_features_suffix.result}"

  settings_json = jsonencode({
    homeScreenPages = [
      {
        "@odata.type" = "#microsoft.graph.iosHomeScreenPage"
        displayName   = "Productivity"
        icons = [
          {
            "@odata.type" = "#microsoft.graph.iosHomeScreenApp"
            displayName   = "Outlook"
            bundleID      = "com.microsoft.Office.Outlook"
            isWebClip     = false
          },
          {
            "@odata.type" = "#microsoft.graph.iosHomeScreenFolder"
            displayName   = "Utilities"
            pages = [
              {
                "@odata.type" = "#microsoft.graph.iosHomeScreenFolderPage"
                displayName   = "Utilities Page 1"
                apps = [
                  {
                    "@odata.type" = "#microsoft.graph.iosHomeScreenApp"
                    displayName   = "Safari"
                    bundleID      = "com.apple.mobilesafari"
                    isWebClip     = false
                  }
                ]
              }
            ]
          }
        ]
      }
    ]
  })

  timeouts = {
    create = "50s"
    read   = "5m"
    update = "30m"
    delete = "30m"
  }
}
